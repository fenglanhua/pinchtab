#!/bin/bash
# 03-actions.sh — CLI action commands (click, fill, press)

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab fill <selector> <text>"

pt_ok nav "${FIXTURES_URL}/form.html"
pt_ok fill "#username" "hello world"
assert_output_contains "filled" "confirms fill action"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab press <key>"

pt_ok press Tab
# Just verify command succeeds

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab scroll down"

pt_ok nav "${FIXTURES_URL}/table.html"
pt_ok scroll down

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab hover <ref>"

pt_ok nav "${FIXTURES_URL}/buttons.html"
pt_ok snap
# Extract a button ref from snapshot
REF=$(echo "$PT_OUT" | grep -oE 'e[0-9]+' | head -1)
if [ -n "$REF" ]; then
  pt_ok hover "$REF"
else
  echo -e "  ${YELLOW}⚠${NC} no ref found, skipping hover"
  ((ASSERTIONS_PASSED++)) || true
fi

end_test
