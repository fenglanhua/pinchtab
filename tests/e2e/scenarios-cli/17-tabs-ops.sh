#!/bin/bash
# 17-tabs-ops.sh — Tab-specific operations

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab tab (list)"

pt_ok nav "${FIXTURES_URL}/form.html"
pt_ok tab
assert_output_json

end_test

# SKIP: tabs snapshot/screenshot/eval/cookies/text per-tab operations
# not yet available in cobra refactor
