//go:build integration

package integration

import (
	"encoding/json"
	"testing"
)

// ER5: Unicode content handling
func TestError_UnicodeContent(t *testing.T) {
	// Navigate to a page with Unicode (CJK/emoji/RTL)
	navigate(t, "https://www.wikipedia.org")

	// Test snapshot
	code, snapBody := httpGet(t, "/snapshot")
	if code != 200 {
		t.Errorf("snapshot failed with %d", code)
	}

	// Verify response is valid JSON
	var snapData map[string]any
	if err := json.Unmarshal(snapBody, &snapData); err != nil {
		t.Errorf("snapshot response is not valid JSON: %v", err)
	}

	// Test text
	code, textBody := httpGet(t, "/text")
	if code != 200 {
		t.Errorf("text failed with %d", code)
	}

	// Verify response is valid JSON
	var textData map[string]any
	if err := json.Unmarshal(textBody, &textData); err != nil {
		t.Errorf("text response is not valid JSON: %v", err)
	}
}

// ER6: Empty page handling
func TestError_EmptyPage(t *testing.T) {
	// Navigate to about:blank (empty page)
	navigate(t, "about:blank")

	// Test snapshot
	code, snapBody := httpGet(t, "/snapshot")
	if code != 200 {
		t.Errorf("snapshot on empty page failed with %d", code)
	}

	// Verify response is valid JSON
	var snapData map[string]any
	if err := json.Unmarshal(snapBody, &snapData); err != nil {
		t.Errorf("snapshot response is not valid JSON: %v", err)
	}

	// Test text
	code, textBody := httpGet(t, "/text")
	if code != 200 {
		t.Errorf("text on empty page failed with %d", code)
	}

	// Verify response is valid JSON
	var textData map[string]any
	if err := json.Unmarshal(textBody, &textData); err != nil {
		t.Errorf("text response is not valid JSON: %v", err)
	}
}
