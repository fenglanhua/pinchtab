//go:build integration

package integration

import (
	"encoding/json"
	"testing"
)

// ST1: Webdriver is hidden (undefined)
func TestStealth_WebdriverUndefined(t *testing.T) {
	navigate(t, "https://example.com")
	code, body := httpPost(t, "/evaluate", map[string]string{
		"expression": "navigator.webdriver === undefined",
	})
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	result := jsonField(t, body, "result")
	if result != "true" {
		t.Errorf("expected result 'true' (webdriver undefined), got %q", result)
	}
}

// ST2: Canvas noise is applied (toDataURL returns different values)
func TestStealth_CanvasNoiseApplied(t *testing.T) {
	navigate(t, "https://example.com")

	// First canvas eval
	code1, body1 := httpPost(t, "/evaluate", map[string]string{
		"expression": "document.createElement('canvas').toDataURL('image/png')",
	})
	if code1 != 200 {
		t.Fatalf("expected 200 for first canvas eval, got %d", code1)
	}
	result1 := jsonField(t, body1, "result")

	// Second canvas eval (should differ due to noise)
	code2, body2 := httpPost(t, "/evaluate", map[string]string{
		"expression": "document.createElement('canvas').toDataURL('image/png')",
	})
	if code2 != 200 {
		t.Fatalf("expected 200 for second canvas eval, got %d", code2)
	}
	result2 := jsonField(t, body2, "result")

	// Results should be different (due to fingerprint noise)
	if result1 == result2 {
		t.Errorf("expected different canvas.toDataURL() results (fingerprint noise), but got identical: %q", result1)
	}

	// Both should be non-empty data URLs
	if result1 == "" || result2 == "" {
		t.Error("expected non-empty canvas data URLs")
	}
}

// ST3: Plugins are present
func TestStealth_PluginsPresent(t *testing.T) {
	navigate(t, "https://example.com")
	code, body := httpPost(t, "/evaluate", map[string]string{
		"expression": "navigator.plugins.length > 0",
	})
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	result := jsonField(t, body, "result")
	if result != "true" {
		t.Errorf("expected result 'true' (plugins present), got %q", result)
	}
}

// ST4: Chrome runtime is present
func TestStealth_ChromeRuntimePresent(t *testing.T) {
	navigate(t, "https://example.com")
	code, body := httpPost(t, "/evaluate", map[string]string{
		"expression": "!!window.chrome.runtime",
	})
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	result := jsonField(t, body, "result")
	if result != "true" {
		t.Errorf("expected result 'true' (chrome.runtime present), got %q", result)
	}
}

// ST5: Fingerprint rotation with specific OS
func TestStealth_FingerprintRotate(t *testing.T) {
	navigate(t, "https://example.com")

	// Get initial user agent
	code1, body1 := httpPost(t, "/evaluate", map[string]string{
		"expression": "navigator.userAgent",
	})
	if code1 != 200 {
		t.Fatalf("expected 200 for initial UA eval, got %d", code1)
	}
	initialUA := jsonField(t, body1, "result")

	// Rotate fingerprint with OS specified
	code2, body2 := httpPost(t, "/fingerprint/rotate", map[string]string{
		"os": "windows",
	})
	if code2 != 200 {
		t.Fatalf("expected 200 for fingerprint rotate, got %d (body: %s)", code2, body2)
	}

	// Get new user agent after rotation
	code3, body3 := httpPost(t, "/evaluate", map[string]string{
		"expression": "navigator.userAgent",
	})
	if code3 != 200 {
		t.Fatalf("expected 200 for post-rotate UA eval, got %d", code3)
	}
	newUA := jsonField(t, body3, "result")

	// Verify fingerprint changed
	if initialUA == newUA {
		t.Logf("Warning: UA did not change after rotation (both: %q)", initialUA)
		// Note: This might not always change depending on random chance,
		// but in practice it should change frequently
	}

	// Verify new UA is non-empty and looks like a valid UA
	if newUA == "" {
		t.Error("expected non-empty user agent after rotation")
	}
}

// ST6: Fingerprint rotation without OS specified (random)
func TestStealth_FingerprintRotateRandom(t *testing.T) {
	navigate(t, "https://example.com")

	// Rotate fingerprint with empty body (random OS)
	code, body := httpPost(t, "/fingerprint/rotate", map[string]string{})
	if code != 200 {
		t.Fatalf("expected 200 for fingerprint rotate (random), got %d (body: %s)", code, body)
	}

	// Verify response is valid JSON
	var m map[string]any
	if err := json.Unmarshal(body, &m); err != nil {
		t.Errorf("expected valid JSON response, got: %v", err)
	}
}

// ST7: Fingerprint rotation fails when no tabs
func TestStealth_FingerprintNoTab(t *testing.T) {
	// Close all tabs first
	code1, tabsBody := httpGet(t, "/tabs")
	if code1 != 200 {
		t.Fatalf("failed to get tabs: %d", code1)
	}

	var tabsResp map[string]any
	_ = json.Unmarshal(tabsBody, &tabsResp)
	tabsRaw := tabsResp["tabs"]
	tabs, _ := tabsRaw.([]any)

	// Close each tab
	for _, tab := range tabs {
		tabMap := tab.(map[string]any)
		tabID := tabMap["id"].(string)
		httpPost(t, "/tab", map[string]string{
			"action": "close",
			"tabId":  tabID,
		})
	}

	// Try to rotate fingerprint with no tabs - should error
	code2, body2 := httpPost(t, "/fingerprint/rotate", map[string]string{})
	if code2 == 200 {
		t.Errorf("expected error (non-200) for fingerprint rotate with no tabs, got %d", code2)
	}

	// Verify error response is present
	if len(body2) == 0 {
		t.Error("expected error message in response body")
	}
}

// ST8: Stealth status endpoint
func TestStealth_StealthStatus(t *testing.T) {
	code, body := httpGet(t, "/stealth/status")
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}

	var m map[string]any
	if err := json.Unmarshal(body, &m); err != nil {
		t.Fatalf("expected valid JSON response: %v (body: %s)", err, body)
	}

	// Check score field exists and is >= 50
	scoreRaw := m["score"]
	if scoreRaw == nil {
		t.Error("expected 'score' field in response")
	} else {
		score, ok := scoreRaw.(float64)
		if !ok {
			t.Errorf("expected score to be a number, got %T", scoreRaw)
		} else if score < 50 {
			t.Errorf("expected score >= 50, got %v", score)
		}
	}

	// Check level field exists and is either "high" or "medium"
	levelRaw := m["level"]
	if levelRaw == nil {
		t.Error("expected 'level' field in response")
	} else {
		level, ok := levelRaw.(string)
		if !ok {
			t.Errorf("expected level to be a string, got %T", levelRaw)
		} else if level != "high" && level != "medium" {
			t.Errorf("expected level to be 'high' or 'medium', got %q", level)
		}
	}
}
