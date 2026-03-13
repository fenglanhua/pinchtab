package config

import "testing"

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"0.8.0", "0.8.0", 0},
		{"0.7.0", "0.8.0", -1},
		{"0.8.0", "0.7.0", 1},
		{"1.0.0", "0.9.9", 1},
		{"0.8.1", "0.8.0", 1},
		{"0.8.0", "0.8.1", -1},
		{"1.0.0", "1.0.0", 0},
	}
	for _, tt := range tests {
		t.Run(tt.a+"_vs_"+tt.b, func(t *testing.T) {
			got := CompareVersions(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("CompareVersions(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestNeedsWizard(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    bool
	}{
		{"empty version", "", true},
		{"old version", "0.7.0", true},
		{"current version", CurrentConfigVersion, false},
		{"future version", "1.0.0", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &FileConfig{ConfigVersion: tt.version}
			if got := NeedsWizard(cfg); got != tt.want {
				t.Errorf("NeedsWizard(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

func TestIsFirstRun(t *testing.T) {
	if !IsFirstRun(&FileConfig{}) {
		t.Error("expected IsFirstRun for empty config")
	}
	if IsFirstRun(&FileConfig{ConfigVersion: "0.8.0"}) {
		t.Error("expected not IsFirstRun for versioned config")
	}
}
