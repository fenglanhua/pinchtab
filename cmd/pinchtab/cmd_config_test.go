package main

import (
	"strings"
	"testing"

	"github.com/pinchtab/pinchtab/internal/config"
)

func TestRenderConfigOverview(t *testing.T) {
	cfg := &config.RuntimeConfig{
		Port:              "9867",
		Strategy:          "simple",
		AllocationPolicy:  "fcfs",
		StealthLevel:      "light",
		TabEvictionPolicy: "close_lru",
		Token:             "very-long-token-secret",
	}
	output := renderConfigOverview(cfg, "/tmp/pinchtab/config.json", "http://localhost:9867", false)

	required := []string{
		"Config",
		"Strategy",
		"Allocation policy",
		"Stealth level",
		"Tab eviction",
		"Copy token",
		"More",
		"/tmp/pinchtab/config.json",
		"very...cret",
		"Dashboard:",
	}
	for _, needle := range required {
		if !strings.Contains(output, needle) {
			t.Fatalf("expected config overview to contain %q\n%s", needle, output)
		}
	}
}

func TestClipboardCommands(t *testing.T) {
	commands := clipboardCommands()
	if len(commands) == 0 {
		t.Fatal("expected clipboard commands")
	}
	for _, command := range commands {
		if command.name == "" {
			t.Fatalf("clipboard command missing name: %+v", command)
		}
	}
}
