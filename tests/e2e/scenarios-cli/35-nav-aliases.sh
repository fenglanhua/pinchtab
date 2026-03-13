#!/bin/bash
# 35-nav-aliases.sh — CLI nav aliases (goto, navigate)

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab nav <url>"

pt_ok nav "${FIXTURES_URL}/index.html"
assert_output_json
assert_output_contains "tabId" "returns tab ID"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab nav --new-tab <url>"

pt_ok nav --new-tab "${FIXTURES_URL}/form.html"
assert_output_json
assert_output_contains "tabId" "opens in new tab"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab goto <url> (alias for nav)"

pt_ok goto "${FIXTURES_URL}/index.html"
assert_output_json
assert_output_contains "tabId" "goto works as alias"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab navigate <url> (alias for nav)"

pt_ok navigate "${FIXTURES_URL}/index.html"
assert_output_json
assert_output_contains "tabId" "navigate works as alias"

end_test
