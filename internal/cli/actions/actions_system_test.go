package actions

import (
	"testing"
)

func TestHealth(t *testing.T) {
	m := newMockServer()
	m.response = `{"status":"ok","version":"dev"}`
	defer m.close()
	client := m.server.Client()

	Health(client, m.base(), "")
	if m.lastPath != "/health" {
		t.Errorf("expected /health, got %s", m.lastPath)
	}
}

func TestAuthHeader(t *testing.T) {
	m := newMockServer()
	defer m.close()
	client := m.server.Client()

	Health(client, m.base(), "my-secret-token")
	auth := m.lastHeaders.Get("Authorization")
	if auth != "Bearer my-secret-token" {
		t.Errorf("expected 'Bearer my-secret-token', got %q", auth)
	}
}

func TestNoAuthHeader(t *testing.T) {
	m := newMockServer()
	defer m.close()
	client := m.server.Client()

	Health(client, m.base(), "")
	auth := m.lastHeaders.Get("Authorization")
	if auth != "" {
		t.Errorf("expected no auth header, got %q", auth)
	}
}
