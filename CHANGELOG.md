# Changelog

## Unreleased (v0.4.0)

### Added
- **Tab locking** — `POST /tab/lock`, `POST /tab/unlock` with timeout-based deadlock prevention for multi-agent coordination
- **Tab ownership** — `/tabs` shows `owner` and `lockedUntil` on locked tabs
- **Token optimization** — `maxTokens`, `selector`, `format=compact` params on `/snapshot`
- **CSS animation disabling** — `BRIDGE_NO_ANIMATIONS` env var + `?noAnimations=true` per-request
- **Stealth levels** — `BRIDGE_STEALTH=light` (default) vs `full`; light mode works with X.com and most sites
- **Welcome page** — headed mode shows 🦀 branding on startup
- **`CHROME_BINARY`** — custom Chrome/Chromium path support
- **`CHROME_FLAGS`** — extra Chrome flags (space-separated)
- **`BRIDGE_BLOCK_MEDIA`** — block all media (images + fonts + CSS + video)
- **`/welcome` endpoint** — serves embedded welcome page

### Fixed
- **K10** — Profile hang on startup (lock file cleanup, unclean exit detection, 15s timeout, auto-retry)
- **K11** — File output `path` param now honored, auto-creates parent dirs
- **`blockImages` on `CreateTab`** — global image/media blocking applied to new tabs
- **`Date.getTimezoneOffset` infinite recursion** — stealth script was calling itself; saved original method reference
- **UA mismatch detection** — removed hardcoded User-Agent override, Chrome uses native UA

### Changed
- Default stealth level changed from `full` to `light` (compatibility over fingerprint resistance)
- Dockerfile Go version updated to 1.26
- Coverage improved from 28.9% to 36%+

## v0.3.0

- Stealth suite (navigator, WebGL, canvas noise, font metrics, WebRTC, timezone, plugins, Chrome flags)
- Human interaction (bezier mouse, typing simulation)
- Fingerprint rotation via CDP
- Image blocking (`BRIDGE_BLOCK_IMAGES`)
- Stealth injection on all tabs (CreateTab + TabContext)
- Multi-agent concurrency verified (MA1-MA8)
- 92 unit tests + ~100 integration tests

## v0.2.0

- Session persistence (cookies, tabs survive restarts)
- Config file support
- Readability `/text` endpoint
- Smart diff (`?diff=true`)
- YAML/file output formats
- Handler split (4 files)

## v0.1.0

- Initial release
- Core HTTP API (16 endpoints)
- Accessibility tree snapshots with stable refs
- Chrome DevTools Protocol bridge
- Tab management
- Basic stealth (webdriver, chrome.runtime, plugins)
