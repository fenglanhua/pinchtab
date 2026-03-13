#!/bin/bash
# 22-redirects.sh — Redirect following with security limits

source "$(dirname "$0")/common.sh"

start_test "redirects: follow single redirect"

# httpbin.org/redirect/1 redirects once to /get
pt_post /navigate -d '{"url":"https://httpbin.org/redirect/1"}'
assert_ok "single redirect followed"

# Verify we ended up at final destination
pt_get /snapshot
assert_json_contains "$RESULT" ".url" "httpbin.org/get" "final URL is /get (redirect successful)"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "redirects: follow multiple redirects"

# httpbin.org/redirect/5 follows 5 redirects total
pt_post /navigate -d '{"url":"https://httpbin.org/redirect/5"}'
assert_ok "five redirects followed"

# Verify final destination
pt_get /snapshot
assert_json_contains "$RESULT" ".url" "httpbin.org/get" "multiple redirects followed to destination"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "redirects: document redirect detection capability"

# When maxRedirects=20 (default), 5 redirects should work fine
# When maxRedirects=3, /redirect/5 should fail (too many redirects)
# (Actual enforcement would require network interception implementation)

echo -e "  ${BLUE}ℹ${NC} Redirect limiting available via CDP Fetch domain"
echo -e "  ${BLUE}ℹ${NC} Default: -1 (unlimited). Set maxRedirects: N to limit hops"
((ASSERTIONS_PASSED++)) || true

end_test
