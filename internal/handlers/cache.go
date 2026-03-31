package handlers

import (
	"fmt"
	"net/http"

	"github.com/pinchtab/pinchtab/internal/httpx"
)

// HandleCacheClear clears the browser's HTTP disk cache.
// POST /cache/clear
func (h *Handlers) HandleCacheClear(w http.ResponseWriter, r *http.Request) {
	if err := h.ensureChrome(); err != nil {
		httpx.Error(w, http.StatusServiceUnavailable, err)
		return
	}

	if err := h.Bridge.ClearCache(r.Context()); err != nil {
		if h.writeBridgeUnavailable(w, err) {
			return
		}
		httpx.Error(w, http.StatusInternalServerError, fmt.Errorf("clear cache: %w", err))
		return
	}

	httpx.JSON(w, http.StatusOK, map[string]any{"status": "cleared"})
}

// HandleCacheStatus checks if the browser cache can be cleared.
// GET /cache/status
func (h *Handlers) HandleCacheStatus(w http.ResponseWriter, r *http.Request) {
	if err := h.ensureChrome(); err != nil {
		httpx.Error(w, http.StatusServiceUnavailable, err)
		return
	}

	canClear, err := h.Bridge.CanClearCache(r.Context())
	if err != nil {
		if h.writeBridgeUnavailable(w, err) {
			return
		}
		httpx.Error(w, http.StatusInternalServerError, fmt.Errorf("cache status: %w", err))
		return
	}

	httpx.JSON(w, http.StatusOK, map[string]any{"canClear": canClear})
}
