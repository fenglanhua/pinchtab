package report

import "testing"

func TestDaemonStatusLooksRunning(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{
			name:   "launchd running",
			status: "state = running",
			want:   true,
		},
		{
			name:   "systemd running",
			status: "Active: active (running) since Thu 2026-03-12",
			want:   true,
		},
		{
			name:   "stopped",
			status: "Active: inactive (dead)",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := daemonStatusLooksRunning(tt.status); got != tt.want {
				t.Fatalf("daemonStatusLooksRunning(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}
