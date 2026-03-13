#!/bin/bash
# 28-pdf-flags.sh — CLI pdf flags

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab pdf -o custom.pdf"

pt_ok nav "${FIXTURES_URL}/index.html"
pt_ok pdf -o /tmp/e2e-custom.pdf

if [ -f /tmp/e2e-custom.pdf ]; then
  echo -e "  ${GREEN}✓${NC} file created"
  ((ASSERTIONS_PASSED++)) || true
  rm -f /tmp/e2e-custom.pdf
else
  echo -e "  ${RED}✗${NC} file not created"
  ((ASSERTIONS_FAILED++)) || true
fi

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab pdf --landscape"

pt_ok pdf --landscape -o /tmp/e2e-landscape.pdf
rm -f /tmp/e2e-landscape.pdf

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab pdf --scale 0.5"

pt_ok pdf --scale 0.5 -o /tmp/e2e-scaled.pdf
rm -f /tmp/e2e-scaled.pdf

end_test
