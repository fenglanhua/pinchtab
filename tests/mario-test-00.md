# Pinchtab Test Report — 00:00 UTC, 2026-02-17

**Branch:** autorun
**Agent:** Mario (OpenClaw)
**Hour:** 00 (even — test run)

---

## Unit Tests (`go test ./... -v -count=1`)

| Metric | Value |
|--------|-------|
| Total pass | 77 (including sub-tests) |
| Total fail | 0 |
| Total skip | 0 |
| Duration | ~0.33s |
| Result | ✅ **ALL PASS** |

---

## Integration Tests (`go test -tags integration -v -count=1`)

| Metric | Value |
|--------|-------|
| Total pass | 77 (including sub-tests) |
| Total fail | 0 |
| Total skip | 1 |
| Duration | ~3.3s |
| Result | ✅ **ALL PASS** (1 skip) |

### Skipped Tests
- **TestWebGLVendorSpoofed** — Expected skip in headless mode (no GPU). Corresponds to SI4.

### Integration Test Mapping (Section 4)
| Test | Plan ID | Status |
|------|---------|--------|
| TestStealthScriptInjected | SI1 | ✅ Pass |
| TestCanvasNoiseApplied | SI2 | ✅ Pass |
| TestFontMetricsNoise | SI3 | ✅ Pass |
| TestWebGLVendorSpoofed | SI4 | ⏭️ Skip (headless) |
| TestPluginsPresent | SI5 | ✅ Pass |
| TestFingerprintRotation | SI6 | ✅ Pass |
| TestCDPTimezoneOverride | SI7 | ✅ Pass |
| TestStealthStatusEndpoint | SI8 | ✅ Pass |

---

## Performance Metrics

| Metric | Value |
|--------|-------|
| Build time | 0.45s (0.60s user, 0.43s sys) |
| Binary size | 12M |
| Unit test duration | 0.33s |
| Integration test duration | 3.3s |
| Benchmarks | No bench functions defined (ran 0 benchmarks) |

---

## TEST-PLAN.md Scenario Coverage

### Automated (via `go test`)
- **Section 1.2** N5 (invalid URL), N6 (missing URL), N7 (bad JSON) — via TestHandleNavigate_* ✅
- **Section 1.3** S10 (no tab) — via TestHandleSnapshot_NoTab ✅
- **Section 1.4** T4 (no tab) — via TestHandleText_NoTab ✅
- **Section 1.5** A9 (unknown kind), A10 (missing kind), A11 (ref not found), A13 (no tab), A15 (empty batch) — via TestHandleAction_* ✅
- **Section 1.6** TB4 (close no tabId), TB5 (bad action) — via TestHandleTab_* ✅
- **Section 1.7** SS3 (no tab) — via TestHandleScreenshot_NoTab ✅
- **Section 1.8** E3 (missing expression), E4 (bad JSON), E5 (no tab) — via TestHandleEvaluate_* ✅
- **Section 1.9** C3 (no tab), C4 (bad JSON), C5 (empty cookies) — via TestHandleCookies_* ✅
- **Section 1.10** ST1 (stealth status) — via TestHandleStealthStatus ✅
- **Section 4** SI1-SI8 — integration tests ✅ (SI4 skip expected)

### Not Automated (require running instance + curl)
- Sections 1.1 (H1-H7), 1.2 (N1-N4, N8), 1.3 (S1-S9, S11-S12), 1.4 (T1-T3, T5)
- Section 1.5 (A1-A8, A12, A14, A16-A17)
- Section 1.6 (TB1-TB3, TB6), 1.7 (SS1-SS2), 1.8 (E1-E2)
- Section 1.9 (C1-C2), 1.10 (ST2-ST8)
- Section 2 (headed mode), Section 3 (multi-agent), Section 5 (Docker)

---

## Known Issues (Section 8) Status

| # | Issue | Status |
|---|-------|--------|
| K1 | Active tab tracking unreliable | 🔴 OPEN — still P0 |
| K2 | Tab close hangs | 🟡 OPEN — still P1 |
| K3 | x.com title empty | 🟢 OPEN — P2, SPA limitation |
| K4 | Chrome flag warning | 🟢 OPEN — P2 |
| K5 | Stealth PRNG weak | ✅ FIXED |
| K6 | Chrome UA hardcoded | ✅ FIXED |
| K7 | Fingerprint rotation JS-only | ✅ FIXED |
| K8 | Timezone hardcoded | ✅ FIXED |
| K9 | Stealth status hardcoded | ✅ FIXED |

---

## Release Criteria (Section 9) Progress

### Must Pass (P0)
- ✅ `go test ./...` 100% pass (77 tests)
- ✅ `go test -tags integration` pass (7 pass, 1 skip headless, 0 fail)
- ❌ K1 (active tab tracking) — still open
- ❌ K2 (tab close hangs) — still open
- ⚠️ Zero crashes — no crashes observed this run
- ⚠️ Section 1 curl scenarios — not yet automated in CI

### Should Pass (P1)
- ⚠️ Multi-agent scenarios — not automated
- ⚠️ Stealth bot.sannysoft.com — manual only
- ⚠️ Session persistence — not automated

### Nice to Have (P2)
- ⚠️ Coverage > 30% — not measured
- ✅ K5-K9 fixed
- ⚠️ K3-K4 documented but open
- ✅ Performance baselined (build: 0.45s, binary: 12M)
