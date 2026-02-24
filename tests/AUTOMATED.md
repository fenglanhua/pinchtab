# Automated Integration Tests

This document tracks which scenarios from the test plan are now covered by automated CI tests in `integration/`.

**CI Workflow:** `.github/workflows/integration.yml` — runs on PRs and main branch pushes.

**Run locally:** `go test -tags integration -v -timeout 10m -count=1 ./tests/integration/`

---

## Test Coverage (Automated)

### Health & Startup
- ✅ **H1** — Health check (`GET /health` returns 200 with status=ok)

### Navigation
- ✅ **N1** — Basic navigate to example.com
- ✅ **N2** — Navigate returns title
- ✅ **N4** — Navigate with newTab flag
- ✅ **N5** — Navigate invalid URL returns error
- ✅ **N6** — Navigate missing URL returns 400
- ✅ **N7** — Navigate bad JSON returns 400

### Snapshot (Accessibility Tree)
- ✅ **S1** — Basic snapshot returns nodes/tree
- ✅ **S2** — Interactive filter works
- ✅ **S3** — Depth filter works
- ✅ **S4** — Text format output
- ✅ **S5** — YAML format output
- ✅ **S5** (variant) — maxTokens parameter

### Text Extraction
- ✅ **T1** — Readability mode (`GET /text`)
- ✅ **T2** — Raw mode (`GET /text?mode=raw`)

### Actions
- ✅ **A1** — Click by ref
- ✅ **A2** — Type by ref
- ✅ **A3** — Fill by ref
- ✅ **A4** — Press key
- ✅ **A5** — Focus element
- ✅ **A8** — Scroll page
- ✅ **A9** — Unknown kind returns 400
- ✅ **A10** — Missing kind returns 400
- ✅ **A11** — Ref not found error
- ✅ **A12** — CSS selector click
- ✅ **A14** — Batch actions
- ✅ **A15** — Batch empty returns 400

### Tabs
- ✅ **TB1** — List tabs
- ✅ **TB2** — New tab
- ✅ **TB3** — Close tab
- ✅ **TB5** — Bad action returns 400

### Screenshots
- ✅ **SS1** — Basic screenshot (base64)
- ✅ **SS2** — Raw screenshot (JPEG bytes)

### JavaScript Evaluation
- ✅ **E1** — Simple eval (1+1)
- ✅ **E2** — DOM eval (document.title)
- ✅ **E3** — Missing expression returns 400
- ✅ **E4** — Bad JSON returns 400

### PDF Export
- ✅ **PD1** — PDF base64 output
- ✅ **PD2** — PDF raw bytes
- ✅ **PD3** — PDF save to file
- ✅ **PD5** — PDF landscape mode
- ✅ **PD6** — PDF scale parameter

### Cookies
- ✅ **C1** — Get cookies
- ✅ **C2** — Set cookies
- ✅ **C3** — Get cookies no tab (error)
- ✅ **C4** — Set cookies bad JSON (400)
- ✅ **C5** — Set cookies empty (400)

### Stealth & Fingerprinting
- ✅ **ST1** — navigator.webdriver undefined
- ✅ **ST3** — navigator.plugins present
- ✅ **ST4** — window.chrome.runtime present
- ✅ **ST5** — Fingerprint rotation with OS specified
- ✅ **ST6** — Fingerprint rotation random (no OS)
- ✅ **ST8** — Stealth status endpoint

*Note: ST2 (canvas noise) skipped — unreliable in headless CI. ST7 replaced with specific tab rotation test.*

### Configuration Extended
- ✅ **CF7** — Chrome version default in UA
- ✅ **CF8** — Chrome version persists after fingerprint rotate
- ✅ **CF6** (variant) — Chrome version override via TEST_CHROME_VERSION

---

## Manual Test Coverage

The following scenarios require manual testing or deployment-specific setups:

### Manual Verification (Fix Verified in Code)
- ✅ **CF3-Extended** — CDP_URL mode (fix verified, needs manual test to confirm: `manual/cf3-cdp-create-tab-repro.md`)

### Not Yet Automated
- **N3** — SPA title handling (x.com heavy SPA)
- **N8** — Navigation timeout behavior
- **S6** — Snapshot diff mode
- **S7** — Snapshot diff first call
- **S8** — Snapshot file output
- **S9** — Snapshot with tabId filter
- **S10** — Snapshot no tab error
- **S11** — Large page snapshot (20K+ tokens)
- **S12** — Ref stability across snapshots
- **T3** — Text with tabId
- **T4** — Text no tab
- **T5** — Token efficiency
- **A6** — Hover action
- **A7** — Select option
- **A13** — Action no tab error
- **A16-A17** — Human click/type (bezier movement)
- **TB4** — Close without tabId (400)
- **TB6** — Max tabs limit
- **UP1-UP12** — File upload (requires test assets)
- **CF1-CF5** — Config file precedence (requires file setup)
- **CF3** — CDP_URL external Chrome
- **CF4** — Custom profile directory
- **CF5** — No restore flag
- **SP1-SP3** — Session persistence (requires restart)
- **HM1-HM3** — Headed mode (requires display)
- **MA1-MA8** — Multi-agent scenarios
- **ER1-ER8** — Error handling & edge cases
- **Docker (D1-D7)** — Requires Docker, deployment testing
- **Dashboard** — Requires manual profile management

See `manual/` directory for detailed test plans.

---

## Performance Testing

Token usage, speed benchmarks, and Chrome startup metrics tracked separately in `performance/`.

---

## Statistics

**Automated:** 48 scenarios  
**Manual/Future:** ~50+ scenarios  
**Total Coverage:** 98+ test scenarios across health, nav, snapshot, text, actions, tabs, screenshots, eval, PDF, cookies, stealth, and config

---

*Last updated: 2026-02-24*
