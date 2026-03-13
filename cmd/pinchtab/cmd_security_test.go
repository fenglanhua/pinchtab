package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pinchtab/pinchtab/internal/config"
)

func TestHandleSecurityCommandDefaultConfigSkipsEmptySections(t *testing.T) {
	cfg := testRuntimeConfig()

	output := captureStdout(t, func() {
		handleSecurityCommand(cfg)
	})

	required := []string{
		"Security",
		"All recommended security defaults are active.",
	}
	for _, needle := range required {
		if !strings.Contains(output, needle) {
			t.Fatalf("expected output to contain %q\n%s", needle, output)
		}
	}

	unwanted := []string{
		"Security posture",
		"Warnings",
		"Recommended security defaults",
		"Recommended defaults",
		"Restore recommended security defaults in config?",
		"Interactive restore skipped because stdin/stdout is not a terminal.",
	}
	for _, needle := range unwanted {
		if strings.Contains(output, needle) {
			t.Fatalf("expected output to skip %q\n%s", needle, output)
		}
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = w
	defer func() {
		os.Stdout = orig
	}()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("close writer error = %v", err)
	}
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if err := r.Close(); err != nil {
		t.Fatalf("close reader error = %v", err)
	}
	return string(data)
}

func TestApplyGuardsDownPreset(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "pinchtab", "config.json")
	t.Setenv("PINCHTAB_CONFIG", configPath)

	fc := config.DefaultFileConfig()
	fc.Server.Token = "guarded-token"
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := config.SaveFileConfig(&fc, configPath); err != nil {
		t.Fatalf("SaveFileConfig() error = %v", err)
	}

	cfg, gotPath, changed, err := applyGuardsDownPreset()
	if err != nil {
		t.Fatalf("applyGuardsDownPreset() error = %v", err)
	}
	if !changed {
		t.Fatal("expected guards down preset to change config")
	}
	if gotPath != configPath {
		t.Fatalf("config path = %q, want %q", gotPath, configPath)
	}

	if cfg.Bind != "127.0.0.1" {
		t.Fatalf("Bind = %q, want 127.0.0.1", cfg.Bind)
	}
	if cfg.Token != "guarded-token" {
		t.Fatalf("Token = %q, want existing token to remain", cfg.Token)
	}
	if !cfg.AllowEvaluate || !cfg.AllowMacro || !cfg.AllowScreencast || !cfg.AllowDownload || !cfg.AllowUpload {
		t.Fatalf("expected sensitive endpoints enabled, got %+v", cfg)
	}
	if !cfg.AttachEnabled {
		t.Fatal("expected attach endpoint enabled")
	}
	if got := strings.Join(cfg.AttachAllowHosts, ","); got != "127.0.0.1,localhost,::1" {
		t.Fatalf("AttachAllowHosts = %q", got)
	}
	if got := strings.Join(cfg.AttachAllowSchemes, ","); got != "ws,wss" {
		t.Fatalf("AttachAllowSchemes = %q", got)
	}
	if cfg.IDPI.Enabled || cfg.IDPI.StrictMode || cfg.IDPI.ScanContent || cfg.IDPI.WrapContent {
		t.Fatalf("expected IDPI protections disabled, got %+v", cfg.IDPI)
	}
}

func testRuntimeConfig() *config.RuntimeConfig {
	return &config.RuntimeConfig{
		Bind:               "127.0.0.1",
		Token:              "abcd1234efgh5678",
		AllowEvaluate:      false,
		AllowMacro:         false,
		AllowScreencast:    false,
		AllowDownload:      false,
		AllowUpload:        false,
		AttachEnabled:      false,
		AttachAllowHosts:   []string{"127.0.0.1", "localhost", "::1"},
		AttachAllowSchemes: []string{"ws", "wss"},
		IDPI: config.IDPIConfig{
			Enabled:        true,
			AllowedDomains: []string{"127.0.0.1", "localhost", "::1"},
			StrictMode:     true,
			ScanContent:    true,
			WrapContent:    true,
		},
	}
}
