#!/bin/bash
# 07-screenshot.sh — CLI screenshot command

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab screenshot"

pt_ok nav "${FIXTURES_URL}/buttons.html"

# Just verify the command succeeds (binary output)
pt screenshot
if [ "$PT_CODE" -eq 0 ]; then
  echo -e "  ${GREEN}✓${NC} screenshot succeeded"
  ((ASSERTIONS_PASSED++)) || true
else
  echo -e "  ${RED}✗${NC} screenshot failed"
  ((ASSERTIONS_FAILED++)) || true
fi

end_test

# SKIP: screenshot -o flag not yet in cobra refactor
