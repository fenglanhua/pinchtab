# Third-Party Licenses

Pinchtab depends on the following open-source packages. All are compatible with MIT licensing.

## Direct Dependencies

### chromedp/chromedp
- **License:** MIT
- **Copyright:** (c) 2016-2025 Kenneth Shaw
- **URL:** https://github.com/chromedp/chromedp
- **Purpose:** Chrome DevTools Protocol driver — launches and controls Chrome

### chromedp/cdproto
- **License:** MIT
- **Copyright:** (c) 2016-2025 Kenneth Shaw
- **URL:** https://github.com/chromedp/cdproto
- **Purpose:** Generated Go types for the Chrome DevTools Protocol

## Transitive Dependencies

### chromedp/sysutil
- **License:** MIT
- **Copyright:** (c) 2016-2017 Kenneth Shaw
- **URL:** https://github.com/chromedp/sysutil
- **Purpose:** System utilities for chromedp (finding Chrome binary)

### go-json-experiment/json
- **License:** BSD 3-Clause
- **Copyright:** (c) 2020 The Go Authors
- **URL:** https://github.com/go-json-experiment/json
- **Purpose:** Experimental JSON library used by cdproto

### gobwas/ws
- **License:** MIT
- **Copyright:** (c) 2017-2021 Sergey Kamardin
- **URL:** https://github.com/gobwas/ws
- **Purpose:** WebSocket implementation for CDP communication

### gobwas/httphead
- **License:** MIT
- **Copyright:** (c) 2017 Sergey Kamardin
- **URL:** https://github.com/gobwas/httphead
- **Purpose:** HTTP header parsing (ws dependency)

### gobwas/pool
- **License:** MIT
- **Copyright:** (c) 2017-2019 Sergey Kamardin
- **URL:** https://github.com/gobwas/pool
- **Purpose:** Pool utilities (ws dependency)

### golang.org/x/sys
- **License:** BSD 3-Clause
- **Copyright:** (c) 2009 The Go Authors
- **URL:** https://github.com/golang/sys
- **Purpose:** Go system call wrappers

## Summary

| Package | License | Compatible |
|---------|---------|------------|
| chromedp/chromedp | MIT | ✅ |
| chromedp/cdproto | MIT | ✅ |
| chromedp/sysutil | MIT | ✅ |
| go-json-experiment/json | BSD 3-Clause | ✅ |
| gobwas/ws | MIT | ✅ |
| gobwas/httphead | MIT | ✅ |
| gobwas/pool | MIT | ✅ |
| golang.org/x/sys | BSD 3-Clause | ✅ |

All dependencies are MIT or BSD 3-Clause licensed, fully compatible with Pinchtab's MIT license.
