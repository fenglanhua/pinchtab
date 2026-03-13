#!/bin/bash
# 13-pdf.sh — CLI PDF export command

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab pdf"

pt_ok nav "${FIXTURES_URL}/form.html"

# Just verify the command succeeds (binary output)
pt pdf
if [ "$PT_CODE" -eq 0 ]; then
  echo -e "  ${GREEN}✓${NC} pdf export succeeded"
  ((ASSERTIONS_PASSED++)) || true
else
  echo -e "  ${RED}✗${NC} pdf export failed"
  ((ASSERTIONS_FAILED++)) || true
fi

end_test

# SKIP: pdf -o flag not yet in cobra refactor
