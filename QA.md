# Pinchtab QA Report — 2026-02-15

**Testers:** Bosch, Mario

---

## Round 1 — Build 4fc2a3e

### Bugs Found

| # | Severity | Bug | Found by |
|---|----------|-----|----------|
| 1 | 🔴 P0 | Active tab not tracked after `/navigate` — snapshot/text return stale tab | Bosch |
| 2 | 🔴 P0 | Invalid JSON — unescaped control chars in snapshot/text (Yahoo Finance, StackOverflow) | Bosch |
| 3 | 🟡 P1 | `newTab:true` silently ignored on `/navigate` | Mario |
| 4 | 🟡 P1 | Tab close (`POST /tab`) returns 400 — `tabId` field not recognized | Bosch |
| 5 | 🟡 P1 | No tab switch/focus API — no recovery when tracking drifts | Bosch |
| 6 | 🟡 P1 | `/action` missing `kind` returns unhelpful `"unknown action: "` | Mario |
| 7 | 🟢 P2 | `/navigate` returns empty title on some sites (BBC, x.com) — race condition | Mario |
| 8 | 🟢 P2 | `/text` on google.com includes full language picker blob | Mario |
| 9 | 🟢 P2 | Chrome flag warning: `--disable-blink-features=AutomationControlled` | Bosch |

---

## Round 2 — Build 723c927 (Mario fixes)

### Re-test Results

| # | Test | Result | Notes |
|---|------|--------|-------|
| 1 | Active tab tracking (1st navigate after startup) | ❌ FAIL | First navigate still returns stale profile tab |
| 2 | Active tab tracking (2nd navigate) | ✅ PASS | Tracks correctly after first |
| 3 | Active tab tracking (3rd navigate) | ✅ PASS | Keeps tracking |
| 4 | JSON validity — Yahoo Finance `/snapshot` | ✅ PASS | **Fixed!** 707 nodes, valid JSON |
| 5 | JSON validity — Yahoo Finance `/text` | ✅ PASS | **Fixed!** |
| 6 | JSON validity — StackOverflow | ✅ PASS | **Fixed!** |
| 7 | Tab close | ❌ FAIL | Now hangs indefinitely instead of 400 (regression) |
| 8 | `/action` error message | ✅ PASS | **Fixed!** Lists valid `kind` values |

### Bug Status After Round 2

| # | Bug | Status |
|---|-----|--------|
| 1 | Active tab tracking | ⚠️ Partial — works after 1st navigate, fails on initial |
| 2 | Invalid JSON | ✅ Fixed |
| 3 | `newTab:true` broken | ❓ Not retested |
| 4 | Tab close | ❌ Regressed — now hangs instead of 400 |
| 5 | No tab switch API | ❌ Still missing |
| 6 | `/action` error message | ✅ Fixed |
| 7 | Empty title on navigate | ❓ Not retested |
| 8 | Google language blob in `/text` | ❓ Not retested |
| 9 | Chrome flag warning | ❌ Still present |

### Remaining Issues (priority order)

1. **Tab close hangs** — regression, was 400 now hangs forever
2. **First navigate doesn't set active tab** — stale profile tabs confuse initial tracking
3. **No tab switch API** — needed for recovery
4. **10 tabs accumulated** during testing with no way to clean up

---

## Performance — Token Usage

### Pinchtab `/snapshot` vs `/text` (Mario)

| Site | /snapshot tokens | /text tokens | Savings |
|------|-----------------|-------------|---------|
| Google | ~2K | ~764 | 2.5× |
| GitHub | ~9.8K | ~1.2K | 7.8× |
| x.com | ~2K | ~121 | 17× |
| BBC | ~26.7K | ~3.5K | 7.7× |
| Wikipedia | ~20.5K | ~3.5K | 5.8× |
| LinkedIn | ~7.5K | ~6.1K | 1.2× |

### Pinchtab vs OpenClaw Browser (Bosch)

| Site | Pinchtab snapshot | Pinchtab /text | OpenClaw aria tree |
|------|-------------------|---------------|-------------------|
| Yahoo Finance | ~16K tokens | ~1.4K tokens | ~3.5K tokens |
| Google Finance | ~12K tokens | ~1.1K tokens | ~3.7K tokens |
| Hacker News | ~24K tokens | ~875 tokens | — |

### Key Findings

- **Pinchtab snapshots are 3–4× larger** than OpenClaw aria trees (verbose JSON per node)
- **Pinchtab `/text` is the most token-efficient** (~1K tokens for complex finance pages)
- **OpenClaw aria tree** is the best balance for interactive browsing (~3.5K tokens)
- **Recommendation:** Add a compact text-based snapshot format to close the gap with OpenClaw

---

## Retest Results (post-fix, 2026-02-15)

| Bug | Status | Notes |
|-----|--------|-------|
| 🟡 `newTab:true` broken | ✅ FIXED | Creates new CDP tab, returns new tabId |
| 🟡 `/action` unhelpful error | ✅ FIXED | Lists valid `kind` values |
| 🟢 `/navigate` empty title | ✅ PARTIAL | BBC works ("BBC - Home"), x.com still empty (SPA >2s) |
| 🟢 `/text` Google blob | ✅ FIXED | Tokens dropped ~764 → ~143 |
| 🔴 Active tab tracking | ❌ STILL BROKEN | Navigate→read returns stale tab content |

**Active tab tracking remains the critical P0.** After navigating to x.com, `/text` returned Google's content. Sequential navigate→read is unreliable without explicit `tabId` targeting.

---

## Sites Tested

**Mario:** Google, GitHub, BBC, Wikipedia, x.com, LinkedIn  
**Bosch:** HN, Example.com, Yahoo Finance, Google Finance, Bloomberg, StackOverflow

All loaded fine, no bot detection, zero crashes. ✅

## What Works Well

- ✅ `/navigate` — fast, returns title+URL
- ✅ `/snapshot` — comprehensive a11y tree
- ✅ `/snapshot?filter=interactive` — filters to actionable elements
- ✅ `/text` — clean, compact, token-efficient
- ✅ `/action` with `click` — reliable
- ✅ `/tabs` — accurate listing
- ✅ Fast startup (~3s), headless works, no bot detection
