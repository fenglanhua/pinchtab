# Pinchtab QA Retest — 2026-02-15 16:27

**Build:** latest main (no new commits since last QA)  
**Tester:** Bosch (automated retest)

---

## Bug Fix Results

### 🟡 P1 — `newTab:true` → ✅ PASS
Tabs went from 3→4 after `navigate` with `newTab:true`. New tab opened correctly.

### 🟡 P1 — `/action` Unhelpful Error → ✅ PASS
Response: `{"error":"missing required field 'kind' — valid values: click, type, fill, press, focus, hover, select, scroll"}`
Lists valid values as expected.

### 🟢 P2 — `/navigate` Empty Title → ⚠️ PARTIAL
- **BBC:** ✅ returns `"BBC - Home"`
- **x.com:** ❌ still returns `"title":""`
- x.com uses heavy JS/SPA hydration — title may need longer wait or DOMContentLoaded fallback

### 🟢 P2 — `/text` Google Language Blob → ✅ PASS
Language picker blob (Afrikaans, azərbaycanca, etc.) is gone. Google /text now returns clean content.

---

## 🔴 P0 — Active Tab Tracking — STILL BROKEN

This is the **critical outstanding bug**. After `/navigate`, `/text` and `/snapshot` frequently return data from a **stale tab** instead of the most recently navigated one.

**Evidence from retest:**
| Navigated to | /text returned | /snapshot returned |
|---|---|---|
| x.com | Google data (stale) | Google snapshot (stale) |
| LinkedIn | GitHub data (stale) | GitHub snapshot (stale) |

The active tab pointer drifts, especially when navigating rapidly between sites. Sometimes it works (Google→BBC was fine), sometimes it doesn't.

---

## Performance — Token Counts (Corrected for Active Sites Only)

Only showing results where active tab was correctly tracked:

| Site | /snapshot tokens | /text tokens | Previous /snapshot | Previous /text |
|------|-----------------|-------------|-------------------|---------------|
| Google | ~848 | ~143 | ~2K | ~764 |
| GitHub | ~9,835 | ~1,249 | ~9.8K | ~1.2K |
| BBC | ~26,598 | ~3,479 | ~26.7K | ~3.5K |
| Wikipedia | ~20,474 | ~3,306 | ~20.5K | ~3.5K |
| x.com | ❌ stale | ❌ stale | ~2K | ~121 |
| LinkedIn | ❌ stale | ❌ stale | ~7.5K | ~6.1K |

**Changes from previous benchmarks:**
- Google /text dropped significantly (~764→~143 tokens) — likely the language picker fix removing bloat ✅
- Google /snapshot also smaller (~2K→~848) — same reason
- All other tracked sites are essentially unchanged
- x.com and LinkedIn couldn't be measured due to active tab bug

---

## Summary

| Fix | Status |
|-----|--------|
| `newTab:true` opens new tab | ✅ PASS |
| `/action` missing kind lists valid values | ✅ PASS |
| `/navigate` BBC returns title | ✅ PASS |
| `/navigate` x.com returns title | ❌ FAIL (still empty) |
| `/text` Google language blob removed | ✅ PASS |
| Active tab tracking (P0) | ❌ STILL BROKEN |

**3 of 4 targeted fixes confirmed working.** x.com title remains empty.

**The P0 active tab tracking bug is the main blocker** — it makes sequential navigate→read workflows unreliable. This must be fixed before Pinchtab can be used for multi-site automation.

### Go Tests
`go test ./...` → ✅ all pass (cached)
