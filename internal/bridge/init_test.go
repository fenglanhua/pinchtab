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
