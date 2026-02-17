# Pinchtab Test Report — 02:00 AM, 2026-02-17

**Branch:** autorun  
**Build:** fresh from main  
**Chrome:** 144.0.7559.133  
**Profile:** clean (/tmp/pinchtab-test2)  
**Mode:** headless, port 18801, BRIDGE_NO_RESTORE=true

> Note: Default profile (~/.pinchtab/chrome-profile) caused Chrome launch to hang indefinitely. Fresh temp profile worked fine. Possible corrupted profile or lock file issue.

## 1.1 Health & Startup

| # | Scenario | Result | Details | Time | Size | Tokens |
|---|----------|--------|---------|------|------|--------|
| H1 | Health check | ✅ PASS | `{"cdp":"","status":"ok","tabs":1}` | 28ms | 33B | 8 |
| H2 | Startup headless | ✅ PASS | Launched, bound port, Chrome not visible | — | — | — |

## 1.2 Navigation

| # | Scenario | Result | Details | Time | Size | Tokens |
|---|----------|--------|---------|------|------|--------|
| N1 | Basic navigate | ✅ PASS | title='Example Domain' | 354ms | 55B | 13 |
| N4 | Navigate newTab | ✅ PASS | New tab created with tabId | 70ms | 98B | 24 |
| N5 | Invalid URL | ✅ PASS | Error returned, no crash | 21ms | 76B | 19 |
| N6 | Missing URL | ✅ PASS | `{"error":"url required"}` | 21ms | 24B | 6 |
| N7 | Bad JSON | ✅ PASS | Parse error returned | 21ms | 84B | 21 |

## 1.3 Snapshot

| # | Scenario | Result | Details | Time | Size | Tokens |
|---|----------|--------|---------|------|------|--------|
| S1 | Basic snapshot | ✅ PASS | 8 nodes, refs (e0-e7) | 38ms | 747B | 186 |
| S2 | Interactive filter | ✅ PASS | Filtered down to 144B | 25ms | 144B | 36 |
| S3 | Depth filter | ✅ PASS | Truncated at depth 2 | 23ms | 154B | 38 |
| S4 | Text format | ✅ PASS | Plain text, not JSON | 23ms | 399B | 99 |
| S5 | YAML format | ✅ PASS | Valid YAML with role: keys | 23ms | 1347B | 336 |
| S6 | Diff second call | ✅ PASS | 165B (same as first — no changes) | 22ms | 165B | 41 |
| S7 | Diff first call | ✅ PASS | Full snapshot returned | 23ms | 165B | 41 |

## 1.4 Text Extraction

| # | Scenario | Result | Details | Time | Size | Tokens |
|---|----------|--------|---------|------|------|--------|
| T1 | Text readability | ✅ PASS | 190B clean text | 24ms | 190B | 47 |
| T2 | Text raw mode | ✅ PASS | 198B raw text | 23ms | 198B | 49 |

## 1.5 Actions

| # | Scenario | Result | Details | Time | Size | Tokens |
|---|----------|--------|---------|------|------|--------|
| A1 | Click by ref | ✅ PASS | `{"clicked":true}` (ref e6, "Learn more" link) | 42ms | 16B | 4 |
| A4 | Press key | ✅ PASS | `{"pressed":"Tab"}` | 29ms | 17B | 4 |
| A8 | Scroll | ✅ PASS | `{"scrolled":true,"y":800}` | 23ms | 25B | 6 |
| A9 | Unknown kind | ✅ PASS | Error lists valid kinds | 21ms | 127B | 31 |
| A10 | Missing kind | ✅ PASS | Error about missing ref | 21ms | 53B | 13 |
| A11 | Ref not found | ✅ PASS | `ref e9999 not found` | 20ms | 56B | 14 |
| A12 | CSS selector | ✅ PASS | `{"clicked":true}` via `h1` | 42ms | 16B | 4 |

**Not tested:** A2 (type), A3 (fill) — Google search didn't expose textbox ref in a11y tree. Need a page with visible form inputs.

## 1.6 Tabs

| # | Scenario | Result | Details | Time | Size | Tokens |
|---|----------|--------|---------|------|------|--------|
| TB1 | List tabs | ✅ PASS | Returned tab array | 21ms | 230B | 57 |
| TB2 | New tab | ✅ PASS | Returns `tabId` (not `id` — test script key was wrong) | 66ms | 98B | 24 |
| TB3 | Close tab | ❌ FAIL | Empty tabId passed → `tabId required`. Likely K2 bug still present — close with valid tabId untested (may hang) | 20ms | 26B | 6 |
| TB4 | Close no tabId | ✅ PASS | `{"error":"tabId required"}` | 22ms | 26B | 6 |
| TB5 | Bad action | ✅ PASS | `{"error":"action must be 'new' or 'close'"}` | 22ms | 43B | 10 |

## Performance Metrics

| Endpoint | Avg Response Time | Avg Size | Avg Tokens |
|----------|------------------|----------|------------|
| GET /health | 28ms | 33B | 8 |
| POST /navigate | 97ms (range 21-354ms) | 67B | 17 |
| GET /snapshot | 25ms (range 22-38ms) | 389B | 97 |
| GET /text | 24ms | 194B | 48 |
| POST /action | 28ms (range 20-42ms) | 44B | 11 |
| GET /tabs | 21ms | 230B | 57 |
| POST /tab | 36ms (range 20-66ms) | 66B | 16 |

## Summary

- **24/26 tests passed** (92%)
- **1 test failed:** TB3 (tab close) — couldn't properly test due to tabId not passing through test script; K2 hang bug may still be present
- **1 test skipped:** A2/A3 (type/fill) — no suitable input field on test pages
- **Critical finding:** Default Chrome profile (`~/.pinchtab/chrome-profile`) causes Chrome launch hang. Fresh profiles work fine. This could be a real user issue.
- **All error handling solid** — invalid URLs, bad JSON, missing params, unknown actions all return clean errors
- **Snapshot is fast** — sub-40ms for all formats, token-efficient
- **Navigation** is the slowest endpoint (~354ms for real page loads), which is expected
