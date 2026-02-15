package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// TabState represents a saved tab for session persistence.
type TabState struct {
	ID    string `json:"id"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

// SessionState is the on-disk format for saved sessions.
type SessionState struct {
	Tabs    []TabState `json:"tabs"`
	SavedAt string     `json:"savedAt"`
}

// markCleanExit patches Chrome's preferences to prevent "didn't shut down correctly" bar.
func markCleanExit() {
	prefsPath := filepath.Join(profileDir, "Default", "Preferences")
	data, err := os.ReadFile(prefsPath)
	if err != nil {
		return
	}
	patched := strings.ReplaceAll(string(data), `"exit_type":"Crashed"`, `"exit_type":"Normal"`)
	patched = strings.ReplaceAll(patched, `"exited_cleanly":false`, `"exited_cleanly":true`)
	if patched != string(data) {
		if err := os.WriteFile(prefsPath, []byte(patched), 0644); err != nil {
			slog.Error("patch prefs", "err", err)
		}
	}
}

// SaveState writes all open tab URLs to sessions.json.
func (b *Bridge) SaveState() {
	targets, err := b.ListTargets()
	if err != nil {
		slog.Error("save state: list targets", "err", err)
		return
	}

	tabs := make([]TabState, 0, len(targets))
	for _, t := range targets {
		if t.URL != "" && t.URL != "about:blank" && t.URL != "chrome://newtab/" {
			tabs = append(tabs, TabState{
				ID:    string(t.TargetID),
				URL:   t.URL,
				Title: t.Title,
			})
		}
	}

	state := SessionState{
		Tabs:    tabs,
		SavedAt: time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		slog.Error("save state: marshal", "err", err)
		return
	}
	path := filepath.Join(stateDir, "sessions.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		slog.Error("save state: write", "err", err)
	} else {
		slog.Info("saved tabs", "count", len(tabs), "path", path)
	}
}

// RestoreState reopens tabs from the last saved session.
func (b *Bridge) RestoreState() {
	path := filepath.Join(stateDir, "sessions.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var state SessionState
	if err := json.Unmarshal(data, &state); err != nil {
		return
	}

	restored := 0
	for _, tab := range state.Tabs {
		if strings.Contains(tab.URL, "/sorry/") || strings.Contains(tab.URL, "about:blank") {
			continue
		}
		ctx, cancel := chromedp.NewContext(b.browserCtx)
		tCtx, tCancel := context.WithTimeout(ctx, 10*time.Second)
		if err := chromedp.Run(tCtx, chromedp.Navigate(tab.URL)); err != nil {
			tCancel()
			cancel()
			slog.Warn("restore tab failed", "url", tab.URL, "err", err)
			continue
		}
		tCancel()
		newID := string(chromedp.FromContext(ctx).Target.TargetID)
		b.tabs[newID] = &TabEntry{ctx: ctx, cancel: cancel}
		restored++
	}
	if restored > 0 {
		slog.Info("restored tabs", "count", restored)
	}
}
