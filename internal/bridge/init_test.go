package bridge

import (
	"slices"
	"testing"

	"github.com/pinchtab/pinchtab/internal/config"
)

func TestBuildChromeArgsSuppressesCrashDialogs(t *testing.T) {
	args := buildChromeArgs(&config.RuntimeConfig{}, 9222)

	for _, want := range []string{
		"--disable-session-crashed-bubble",
		"--hide-crash-restore-bubble",
		"--noerrdialogs",
	} {
		if !slices.Contains(args, want) {
			t.Fatalf("missing chrome arg %q in %v", want, args)
		}
	}
}

func TestDefaultChromeFlagArgsDisablesMetricsReporting(t *testing.T) {
	args := defaultChromeFlagArgs()
	for _, want := range []string{"--disable-metrics-reporting", "--metrics-recording-only"} {
		found := false
		for _, arg := range args {
			if arg == want {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected %s in args, got %v", want, args)
		}
	}
}
