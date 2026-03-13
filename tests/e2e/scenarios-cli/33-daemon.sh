#!/bin/bash
# 33-daemon.sh — CLI daemon command

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab daemon (non-interactive shows status)"

pt daemon
assert_exit_code 0 "daemon status displayed"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab daemon install (fails without systemd)"

pt daemon install
if [ "$PT_CODE" -ne 0 ]; then
  echo -e "  ${GREEN}✓${NC} fails gracefully without systemd (exit $PT_CODE)"
  ((ASSERTIONS_PASSED++)) || true
else
  echo -e "  ${RED}✗${NC} daemon install unexpectedly succeeded"
  ((ASSERTIONS_FAILED++)) || true
fi

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab daemon unknown-subcommand → exit 2"

pt daemon bogus-command
assert_exit_code 2 "unknown subcommand rejected"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab daemon start (fails without service manager)"

pt daemon start
if [ "$PT_CODE" -ne 0 ]; then
  echo -e "  ${GREEN}✓${NC} start fails gracefully without service manager (exit $PT_CODE)"
  ((ASSERTIONS_PASSED++)) || true
else
  echo -e "  ${RED}✗${NC} daemon start unexpectedly succeeded"
  ((ASSERTIONS_FAILED++)) || true
fi

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab daemon stop (fails without service manager)"

pt daemon stop
if [ "$PT_CODE" -ne 0 ]; then
  echo -e "  ${GREEN}✓${NC} stop fails gracefully without service manager (exit $PT_CODE)"
  ((ASSERTIONS_PASSED++)) || true
else
  echo -e "  ${RED}✗${NC} daemon stop unexpectedly succeeded"
  ((ASSERTIONS_FAILED++)) || true
fi

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab daemon restart (fails without service manager)"

pt daemon restart
if [ "$PT_CODE" -ne 0 ]; then
  echo -e "  ${GREEN}✓${NC} restart fails gracefully without service manager (exit $PT_CODE)"
  ((ASSERTIONS_PASSED++)) || true
else
  echo -e "  ${RED}✗${NC} daemon restart unexpectedly succeeded"
  ((ASSERTIONS_FAILED++)) || true
fi

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab daemon uninstall (graceful when not installed)"

pt daemon uninstall
assert_exit_code_lte 1 "uninstall handled gracefully"

end_test
