package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"sync"
	"time"
)

// DashboardConfig holds tunable timeouts for agent status transitions.
type DashboardConfig struct {
	IdleTimeout       time.Duration // time before agent marked idle (default 30s)
	DisconnectTimeout time.Duration // time before agent marked disconnected (default 5m)
	ReaperInterval    time.Duration // how often to check agent status (default 10s)
	SSEBufferSize     int           // per-client SSE channel buffer (default 64)
}

//go:embed dashboard
var dashboardFS embed.FS

// ---------------------------------------------------------------------------
// Agent Activity — tracks what each agent is doing in real time
// ---------------------------------------------------------------------------

type AgentActivity struct {
	AgentID    string    `json:"agentId"`
	Profile    string    `json:"profile,omitempty"`
	CurrentURL string    `json:"currentUrl,omitempty"`
	CurrentTab string    `json:"currentTab,omitempty"`
	LastAction string    `json:"lastAction,omitempty"`
	LastSeen   time.Time `json:"lastSeen"`
	Status     string    `json:"status"` // "active", "idle", "disconnected"
	ActionCount int      `json:"actionCount"`
}

type AgentEvent struct {
	AgentID   string `json:"agentId"`
	Profile   string `json:"profile,omitempty"`
	Action    string `json:"action"`
	URL       string `json:"url,omitempty"`
	TabID     string `json:"tabId,omitempty"`
	Detail    string `json:"detail,omitempty"`
	Status    int    `json:"status"`
	DurationMs int64 `json:"durationMs"`
	Timestamp time.Time `json:"timestamp"`
}

type Dashboard struct {
	cfg      DashboardConfig
	agents   map[string]*AgentActivity
	sseConns map[chan AgentEvent]struct{}
	cancel   context.CancelFunc
	mu       sync.RWMutex
}

func NewDashboard(cfg *DashboardConfig) *Dashboard {
	c := DashboardConfig{
		IdleTimeout:       30 * time.Second,
		DisconnectTimeout: 5 * time.Minute,
		ReaperInterval:    10 * time.Second,
		SSEBufferSize:     64,
	}
	if cfg != nil {
		if cfg.IdleTimeout > 0 {
			c.IdleTimeout = cfg.IdleTimeout
		}
		if cfg.DisconnectTimeout > 0 {
			c.DisconnectTimeout = cfg.DisconnectTimeout
		}
		if cfg.ReaperInterval > 0 {
			c.ReaperInterval = cfg.ReaperInterval
		}
		if cfg.SSEBufferSize > 0 {
			c.SSEBufferSize = cfg.SSEBufferSize
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	d := &Dashboard{
		cfg:      c,
		agents:   make(map[string]*AgentActivity),
		sseConns: make(map[chan AgentEvent]struct{}),
		cancel:   cancel,
	}
	go d.reaper(ctx)
	return d
}

// Shutdown stops the reaper goroutine.
func (d *Dashboard) Shutdown() { d.cancel() }

func (d *Dashboard) reaper(ctx context.Context) {
	ticker := time.NewTicker(d.cfg.ReaperInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			d.mu.Lock()
			now := time.Now()
			for id, a := range d.agents {
				if a.Status == "disconnected" {
					continue
				}
				if now.Sub(a.LastSeen) > d.cfg.DisconnectTimeout {
					d.agents[id].Status = "disconnected"
				} else if now.Sub(a.LastSeen) > d.cfg.IdleTimeout {
					d.agents[id].Status = "idle"
				}
			}
			d.mu.Unlock()
		}
	}
}

// RecordEvent processes an agent action and broadcasts to SSE subscribers.
func (d *Dashboard) RecordEvent(evt AgentEvent) {
	d.mu.Lock()

	a, ok := d.agents[evt.AgentID]
	if !ok {
		a = &AgentActivity{AgentID: evt.AgentID}
		d.agents[evt.AgentID] = a
	}
	a.LastSeen = evt.Timestamp
	a.LastAction = evt.Action
	a.Status = "active"
	a.ActionCount++
	a.Profile = evt.Profile
	if evt.URL != "" {
		a.CurrentURL = evt.URL
	}
	if evt.TabID != "" {
		a.CurrentTab = evt.TabID
	}

	// Copy SSE channels
	chans := make([]chan AgentEvent, 0, len(d.sseConns))
	for ch := range d.sseConns {
		chans = append(chans, ch)
	}
	d.mu.Unlock()

	// Non-blocking broadcast
	for _, ch := range chans {
		select {
		case ch <- evt:
		default: // drop if slow
		}
	}
}

// GetAgents returns current state of all agents.
func (d *Dashboard) GetAgents() []AgentActivity {
	d.mu.RLock()
	defer d.mu.RUnlock()

	agents := make([]AgentActivity, 0, len(d.agents))
	for _, a := range d.agents {
		agents = append(agents, *a)
	}
	return agents
}

// ---------------------------------------------------------------------------
// HTTP Handlers
// ---------------------------------------------------------------------------

func (d *Dashboard) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /dashboard", d.handleDashboardUI)
	mux.HandleFunc("GET /dashboard/agents", d.handleAgents)
	mux.HandleFunc("GET /dashboard/events", d.handleSSE)

	// Serve static assets (CSS, JS) from embedded filesystem
	sub, _ := fs.Sub(dashboardFS, "dashboard")
	mux.Handle("GET /dashboard/", http.StripPrefix("/dashboard/", http.FileServer(http.FS(sub))))
}

func (d *Dashboard) handleAgents(w http.ResponseWriter, r *http.Request) {
	jsonResp(w, 200, d.GetAgents())
}

func (d *Dashboard) handleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := make(chan AgentEvent, d.cfg.SSEBufferSize)
	d.mu.Lock()
	d.sseConns[ch] = struct{}{}
	d.mu.Unlock()

	defer func() {
		d.mu.Lock()
		delete(d.sseConns, ch)
		d.mu.Unlock()
	}()

	// Send current agent state as initial event
	agents := d.GetAgents()
	data, _ := json.Marshal(agents)
	fmt.Fprintf(w, "event: init\ndata: %s\n\n", data)
	flusher.Flush()

	for {
		select {
		case evt := <-ch:
			data, _ := json.Marshal(evt)
			fmt.Fprintf(w, "event: action\ndata: %s\n\n", data)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (d *Dashboard) handleDashboardUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	data, _ := dashboardFS.ReadFile("dashboard/dashboard.html")
	w.Write(data)
}

// ---------------------------------------------------------------------------
// Tracking Middleware — extracts agent ID from header or query
// ---------------------------------------------------------------------------

// EventObserver receives agent events for additional processing (e.g. profile tracking).
type EventObserver func(evt AgentEvent)

// extractAgentID reads the agent identifier from X-Agent-Id header or agentId query param.
func extractAgentID(r *http.Request) string {
	if id := r.Header.Get("X-Agent-Id"); id != "" {
		return id
	}
	if id := r.URL.Query().Get("agentId"); id != "" {
		return id
	}
	return "anonymous"
}

// extractProfile reads the profile name from X-Profile header or profile query param.
func extractProfile(r *http.Request) string {
	if p := r.Header.Get("X-Profile"); p != "" {
		return p
	}
	return r.URL.Query().Get("profile")
}

// isManagementRoute returns true for routes that shouldn't be tracked in the activity feed.
func isManagementRoute(path string) bool {
	return strings.HasPrefix(path, "/dashboard") ||
		strings.HasPrefix(path, "/profiles") ||
		strings.HasPrefix(path, "/instances") ||
		strings.HasPrefix(path, "/screencast/tabs") ||
		path == "/welcome" || path == "/favicon.ico" || path == "/health"
}

// actionDetail extracts a human-readable detail string from the request.
func actionDetail(r *http.Request) string {
	switch r.URL.Path {
	case "/navigate":
		return r.URL.Query().Get("url")
	case "/actions":
		return "batch action"
	case "/snapshot":
		if sel := r.URL.Query().Get("selector"); sel != "" {
			return "selector=" + sel
		}
	}
	return ""
}

func (d *Dashboard) TrackingMiddleware(observers []EventObserver, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		if isManagementRoute(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		sw := &statusWriter{ResponseWriter: w, code: 200}
		next.ServeHTTP(sw, r)

		evt := AgentEvent{
			AgentID:    extractAgentID(r),
			Profile:    extractProfile(r),
			Action:     r.Method + " " + r.URL.Path,
			URL:        r.URL.Query().Get("url"),
			TabID:      r.URL.Query().Get("tabId"),
			Detail:     actionDetail(r),
			Status:     sw.code,
			DurationMs: time.Since(start).Milliseconds(),
			Timestamp:  start,
		}

		d.RecordEvent(evt)

		for _, obs := range observers {
			obs(evt)
		}
	})
}

// ---------------------------------------------------------------------------
