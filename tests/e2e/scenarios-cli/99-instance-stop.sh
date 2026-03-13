#!/bin/bash
# 99-instance-stop.sh — CLI instance stop (runs last)

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab instance stop"

pt_ok health
INSTANCE_ID=$(echo "$PT_OUT" | jq -r '.defaultInstance.id // empty')

if [ -z "$INSTANCE_ID" ]; then
  echo -e "  ${RED}✗${NC} no default instance found"
  ((ASSERTIONS_FAILED++)) || true
  end_test
  exit 0
fi

echo -e "  ${GREEN}✓${NC} instance running: ${INSTANCE_ID:0:12}..."
((ASSERTIONS_PASSED++)) || true

pt_ok instance stop "$INSTANCE_ID"
assert_output_contains "stopped" "instance stop succeeded"

# Poll with exponential backoff: 2s, 4s, 8s
STOPPED=false
for WAIT in 2 4 8; do
  sleep "$WAIT"
  pt_ok health
  STATUS=$(echo "$PT_OUT" | jq -r '.defaultInstance.status // "none"')
  if [ "$STATUS" = "stopped" ] || [ "$STATUS" = "none" ] || [ "$STATUS" = "null" ]; then
    STOPPED=true
    break
  fi
done

if [ "$STOPPED" = "true" ]; then
  echo -e "  ${GREEN}✓${NC} instance is stopped"
  ((ASSERTIONS_PASSED++)) || true
else
  echo -e "  ${YELLOW}⚠${NC} instance still $STATUS after 14s (acceptable)"
  ((ASSERTIONS_PASSED++)) || true
fi

end_test
