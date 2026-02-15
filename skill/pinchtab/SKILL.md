---
name: pinchtab
description: >
  Control a headless or headed Chrome browser via Pinchtab's HTTP API. Use for web automation,
  scraping, form filling, navigation, and multi-tab workflows. Pinchtab exposes the accessibility
  tree as flat JSON with stable refs — optimized for AI agents (low token cost, fast).
  Use when the task involves: browsing websites, filling forms, clicking buttons, extracting
  page text, taking screenshots, or any browser-based automation. Requires a running Pinchtab
  instance (Go binary).
homepage: https://pinchtab.com
metadata:
  openclaw:
    emoji: "🦀"
    requires:
      bins: ["pinchtab"]
---

# Pinchtab

Fast, lightweight browser control for AI agents via HTTP + accessibility tree.

## Setup

Ensure Pinchtab is running:

```bash
# Headless (default for automation)
BRIDGE_HEADLESS=true pinchtab &

# With UI (debugging)
pinchtab &
```

Default port: `18800`. Override with `BRIDGE_PORT=18801`.
Auth: set `BRIDGE_TOKEN=<secret>` and pass `Authorization: Bearer <secret>`.

Base URL for all examples: `http://localhost:18800`

## Core Workflow

The typical agent loop:

1. **Navigate** to a URL
2. **Snapshot** the accessibility tree (get refs)
3. **Act** on refs (click, type, press)
4. **Snapshot** again to see results

Refs (e.g. `e0`, `e5`, `e12`) are cached per tab after each snapshot — no need to re-snapshot before every action unless the page changed significantly.

## API Reference

### Navigate

```bash
curl -X POST http://localhost:18800/navigate \
  -H 'Content-Type: application/json' \
  -d '{"url": "https://example.com"}'
```

### Snapshot (accessibility tree)

```bash
# Full tree
curl http://localhost:18800/snapshot

# Interactive elements only (buttons, links, inputs) — much smaller
curl "http://localhost:18800/snapshot?filter=interactive"

# Limit depth
curl "http://localhost:18800/snapshot?depth=5"
```

Returns flat JSON array of nodes with `ref`, `role`, `name`, `depth`, `value`, `nodeId`.

**Token optimization**: Use `?filter=interactive` for action-oriented tasks (~75% fewer tokens). Use full snapshot only when you need to read page content.

### Act on elements

```bash
# Click by ref
curl -X POST http://localhost:18800/action \
  -H 'Content-Type: application/json' \
  -d '{"kind": "click", "ref": "e5"}'

# Type into focused element (click first, then type)
curl -X POST http://localhost:18800/action \
  -H 'Content-Type: application/json' \
  -d '{"kind": "click", "ref": "e12"}'
curl -X POST http://localhost:18800/action \
  -H 'Content-Type: application/json' \
  -d '{"kind": "type", "ref": "e12", "text": "hello world"}'

# Press a key
curl -X POST http://localhost:18800/action \
  -H 'Content-Type: application/json' \
  -d '{"kind": "press", "key": "Enter"}'

# Focus an element
curl -X POST http://localhost:18800/action \
  -H 'Content-Type: application/json' \
  -d '{"kind": "focus", "ref": "e3"}'

# Fill (set value directly, no keystrokes)
curl -X POST http://localhost:18800/action \
  -H 'Content-Type: application/json' \
  -d '{"kind": "fill", "selector": "#email", "text": "user@example.com"}'
```

### Extract text

```bash
curl http://localhost:18800/text
```

Returns `{url, title, text}`. Cheapest option (~1K tokens for most pages).

### Screenshot

```bash
# Raw JPEG bytes
curl "http://localhost:18800/screenshot?raw=true" -o screenshot.jpg

# With quality setting (default 80)
curl "http://localhost:18800/screenshot?raw=true&quality=50" -o screenshot.jpg
```

### Evaluate JavaScript

```bash
curl -X POST http://localhost:18800/evaluate \
  -H 'Content-Type: application/json' \
  -d '{"expression": "document.title"}'
```

### Tab management

```bash
# List tabs
curl http://localhost:18800/tabs

# Open new tab
curl -X POST http://localhost:18800/tab \
  -H 'Content-Type: application/json' \
  -d '{"action": "new", "url": "https://example.com"}'

# Close tab
curl -X POST http://localhost:18800/tab \
  -H 'Content-Type: application/json' \
  -d '{"action": "close", "tabId": "TARGET_ID"}'
```

Multi-tab: pass `?tabId=TARGET_ID` to snapshot/screenshot/text, or `"tabId"` in POST body.

### Health check

```bash
curl http://localhost:18800/health
```

## Token Cost Guide

| Method | Typical tokens | When to use |
|---|---|---|
| `/text` | ~1K | Reading page content |
| `/snapshot?filter=interactive` | ~5K | Finding buttons/links to click |
| `/snapshot` | ~20K | Full page understanding |
| `/screenshot` | ~2K (vision) | Visual verification |

**Strategy**: Start with `/snapshot?filter=interactive`. Use full `/snapshot` only when you need to read static text or understand layout. Use `/text` when you only need the readable content.

## Environment Variables

| Var | Default | Description |
|---|---|---|
| `BRIDGE_PORT` | `18800` | HTTP port |
| `BRIDGE_HEADLESS` | `false` | Run Chrome headless |
| `BRIDGE_TOKEN` | (none) | Bearer auth token |
| `BRIDGE_PROFILE` | `~/.browser-bridge/chrome-profile` | Chrome profile dir |
| `BRIDGE_NO_RESTORE` | `false` | Skip tab restore on startup |
| `CDP_URL` | (none) | Connect to existing Chrome DevTools |

## Tips

- Refs are stable between snapshot and actions — no need to re-snapshot before clicking
- After navigation or major page changes, take a new snapshot to get fresh refs
- Use `filter=interactive` by default, fall back to full snapshot when needed
- Pinchtab persists sessions — tabs survive restarts (disable with `BRIDGE_NO_RESTORE=true`)
- Chrome profile is persistent — cookies/logins carry over between runs
