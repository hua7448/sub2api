#!/usr/bin/env bash

set -euo pipefail

readonly SCHEMA="sub2api.production-contracts.v1"
readonly CANONICAL_REF="refs/heads/production/backend-current"
readonly REQUIRED_ANCESTOR="992223cbb7add9034a87d4497491cd25c0dbe577"
readonly LEGACY_COMPACT_COMMIT="46a237d8dd6b30187330a958975b5b68a0ba390e"
readonly IMAGE_TASK_CORS_COMMIT="18cc7858baa6eab016475466f4e880d3ac49e8ef"
readonly REQUIRED_GO_BIN="/home/ubuntu/.local/go1.26.5/bin/go"
readonly REQUIRED_GO_VERSION="go version go1.26.5 linux/amd64"
readonly MIGRATION_186="backend/migrations/186_subscription_quota_reset_operations.sql"
readonly MIGRATION_186_SHA256="98ea27a13a0f82e8b85a7d890659618b0c2dcb7e4886ae26076ba7defb3da780"
readonly MIGRATION_194="backend/migrations/194_user_subscription_daily_quota_weekly_limit.sql"
readonly MIGRATION_194_SHA256="6ea2ac059abdecc693d9d7c2dfb2456590ed11724635c6cbd29d577410fa9190"

MODE="tests"
CHECK_COUNT=0
TEST_COUNT=0
DIRTY_NON_UI_COUNT=0
HEAD_COMMIT="unknown"
HEAD_TREE="unknown"
CANONICAL_TIP="unknown"
GO_VERSION="not-run"

usage() {
  cat >&2 <<'EOF'
Usage: backend/scripts/verify-production-contracts.sh [--source-only]

Verifies the production backend's immutable source contracts. By default it
also runs the focused Go regression tests. --source-only skips Go execution.

Stdout contains one JSON provenance record. Diagnostics and test output go to
stderr.
EOF
}

json_escape() {
  local value=${1-}
  value=${value//\\/\\\\}
  value=${value//\"/\\\"}
  value=${value//$'\n'/\\n}
  value=${value//$'\r'/\\r}
  value=${value//$'\t'/\\t}
  printf '%s' "$value"
}

emit_provenance() {
  local status=$1
  local failed_check=${2-}
  local detail=${3-}
  local verified_at
  verified_at=$(date -u +'%Y-%m-%dT%H:%M:%SZ')

  printf '{'
  printf '"schema":"%s"' "$(json_escape "$SCHEMA")"
  printf ',"status":"%s"' "$(json_escape "$status")"
  printf ',"mode":"%s"' "$(json_escape "$MODE")"
  printf ',"head":"%s"' "$(json_escape "$HEAD_COMMIT")"
  printf ',"tree":"%s"' "$(json_escape "$HEAD_TREE")"
  printf ',"canonical_ref":"%s"' "$CANONICAL_REF"
  printf ',"canonical_tip":"%s"' "$(json_escape "$CANONICAL_TIP")"
  printf ',"required_ancestor":"%s"' "$REQUIRED_ANCESTOR"
  printf ',"feature_commits":{"legacy_compact_v2":"%s","image_task_cors":"%s"}' "$LEGACY_COMPACT_COMMIT" "$IMAGE_TASK_CORS_COMMIT"
  printf ',"migrations":{"186":{"path":"%s","sha256":"%s"},"194":{"path":"%s","sha256":"%s"}}' \
    "$MIGRATION_186" "$MIGRATION_186_SHA256" "$MIGRATION_194" "$MIGRATION_194_SHA256"
  printf ',"dirty_non_ui_count":%d' "$DIRTY_NON_UI_COUNT"
  printf ',"checks_passed":%d' "$CHECK_COUNT"
  printf ',"go_tests":%d' "$TEST_COUNT"
  printf ',"go_version":"%s"' "$(json_escape "$GO_VERSION")"
  printf ',"verified_at_utc":"%s"' "$verified_at"
  if [[ -n "$failed_check" ]]; then
    printf ',"failed_check":"%s"' "$(json_escape "$failed_check")"
  fi
  if [[ -n "$detail" ]]; then
    printf ',"detail":"%s"' "$(json_escape "$detail")"
  fi
  printf '}\n'
}

fail() {
  local check=$1
  local detail=$2
  printf 'FAIL [%s] %s\n' "$check" "$detail" >&2
  emit_provenance "fail" "$check" "$detail"
  exit 1
}

pass() {
  local check=$1
  CHECK_COUNT=$((CHECK_COUNT + 1))
  printf 'PASS [%s]\n' "$check" >&2
}

require_file() {
  local check=$1
  local path=$2
  [[ -f "$REPO_ROOT/$path" ]] || fail "$check" "missing file: $path"
  pass "$check"
}

require_literal() {
  local check=$1
  local path=$2
  local literal=$3
  [[ -f "$REPO_ROOT/$path" ]] || fail "$check" "missing file: $path"
  grep -Fq -- "$literal" "$REPO_ROOT/$path" || fail "$check" "missing required source token in $path"
  pass "$check"
}

reject_literal() {
  local check=$1
  local path=$2
  local literal=$3
  [[ -f "$REPO_ROOT/$path" ]] || fail "$check" "missing file: $path"
  if grep -Fq -- "$literal" "$REPO_ROOT/$path"; then
    fail "$check" "forbidden source token in $path"
  fi
  pass "$check"
}

reject_literal_in_block() {
  local check=$1
  local path=$2
  local start_literal=$3
  local end_literal=$4
  local forbidden_literal=$5
  local block
  [[ -f "$REPO_ROOT/$path" ]] || fail "$check" "missing file: $path"
  block=$(awk -v start="$start_literal" -v end="$end_literal" '
    index($0, start) { active=1; found_start=1 }
    active && index($0, end) && index($0, start) == 0 { found_end=1; exit }
    active { print }
    END { if (!found_start || !found_end) exit 2 }
  ' "$REPO_ROOT/$path") || fail "$check" "cannot isolate required source block in $path"
  if [[ "$block" == *"$forbidden_literal"* ]]; then
    fail "$check" "forbidden source token in scoped block in $path"
  fi
  pass "$check"
}

require_literal_in_block() {
  local check=$1
  local path=$2
  local start_literal=$3
  local end_literal=$4
  local required_literal=$5
  local block
  [[ -f "$REPO_ROOT/$path" ]] || fail "$check" "missing file: $path"
  block=$(awk -v start="$start_literal" -v end="$end_literal" '
    index($0, start) { active=1; found_start=1 }
    active && index($0, end) && index($0, start) == 0 { found_end=1; exit }
    active { print }
    END { if (!found_start || !found_end) exit 2 }
  ' "$REPO_ROOT/$path") || fail "$check" "cannot isolate required source block in $path"
  if [[ "$block" != *"$required_literal"* ]]; then
    fail "$check" "missing required source token in scoped block in $path"
  fi
  pass "$check"
}

require_test() {
  local check=$1
  local path=$2
  local test_name=$3
  require_literal "$check" "$path" "func $test_name("
}

sha256_file() {
  local path=$1
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$path" | awk '{print $1}'
    return
  fi
  if command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$path" | awk '{print $1}'
    return
  fi
  fail "sha256_tool" "neither sha256sum nor shasum is available"
}

is_ui_path() {
  case "$1" in
    frontend/*|backend/internal/web/dist/*)
      return 0
      ;;
  esac
  return 1
}

check_dirty_non_ui() {
  local dirty_paths_file
  local path
  local -A seen=()
  local -a dirty=()

  dirty_paths_file=$(mktemp "${TMPDIR:-/tmp}/sub2api-contract-dirty.XXXXXX") || fail "dirty_source_scan" "cannot create temporary dirty-path inventory"
  if ! git -C "$REPO_ROOT" diff --name-only -z -- >"$dirty_paths_file"; then
    rm -f "$dirty_paths_file"
    fail "dirty_source_scan" "cannot enumerate unstaged source changes"
  fi
  if ! git -C "$REPO_ROOT" diff --cached --name-only -z -- >>"$dirty_paths_file"; then
    rm -f "$dirty_paths_file"
    fail "dirty_source_scan" "cannot enumerate staged source changes"
  fi
  if ! git -C "$REPO_ROOT" ls-files --others --exclude-standard -z -- >>"$dirty_paths_file"; then
    rm -f "$dirty_paths_file"
    fail "dirty_source_scan" "cannot enumerate untracked source paths"
  fi

  while IFS= read -r -d '' path; do
    [[ -n "$path" ]] || continue
    [[ -z "${seen[$path]+x}" ]] || continue
    seen[$path]=1
    if is_ui_path "$path"; then
      continue
    fi
    dirty+=("$path")
  done <"$dirty_paths_file"
  rm -f "$dirty_paths_file"

  if ((${#dirty[@]} > 0)); then
    DIRTY_NON_UI_COUNT=${#dirty[@]}
    fail "dirty_non_ui_source" "uncommitted non-UI path: ${dirty[0]} (total ${#dirty[@]})"
  fi
  pass "dirty_non_ui_source"
}

find_go() {
  local candidate=${GO_BIN:-$REQUIRED_GO_BIN}
  local candidate_real
  local required_real

  [[ -x "$REQUIRED_GO_BIN" && -x "$candidate" ]] || return 1
  required_real=$(readlink -f -- "$REQUIRED_GO_BIN") || return 1
  candidate_real=$(readlink -f -- "$candidate") || return 1
  [[ "$candidate_real" == "$required_real" ]] || return 1
  printf '%s' "$candidate_real"
}

json_test_action_seen() {
  local action=$1
  local test_name=$2
  local line

  while IFS= read -r line; do
    if [[ "$line" == *"\"Action\":\"$action\""* && "$line" == *"\"Test\":\"$test_name\""* ]]; then
      return 0
    fi
  done
  return 1
}

run_go_test() {
  local check=$1
  local package=$2
  local go_bin=$3
  shift 3
  local -a expected_tests=("$@")
  local joined_tests
  local output
  local pattern
  local test_name

  ((${#expected_tests[@]} > 0)) || fail "$check" "no expected Go tests were declared"
  joined_tests=$(IFS='|'; printf '%s' "${expected_tests[*]}")
  pattern="^(${joined_tests})$"

  printf 'TEST [%s] %s\n' "$check" "$package" >&2
  if output=$(cd "$REPO_ROOT/backend" && GOENV=off GOFLAGS= GOWORK=off GOTOOLCHAIN=local "$go_bin" test -json -tags unit "$package" -run "$pattern" -count=1 2>&1); then
    printf '%s\n' "$output" >&2
  else
    printf '%s\n' "$output" >&2
    fail "$check" "focused Go regression test failed"
  fi

  for test_name in "${expected_tests[@]}"; do
    if ! json_test_action_seen "run" "$test_name" <<<"$output"; then
      fail "$check" "expected Go test did not run: $test_name"
    fi
    if ! json_test_action_seen "pass" "$test_name" <<<"$output"; then
      fail "$check" "expected Go test did not pass: $test_name"
    fi
  done

  TEST_COUNT=$((TEST_COUNT + ${#expected_tests[@]}))
  pass "$check"
}

case "${1-}" in
  "")
    ;;
  --source-only)
    MODE="source-only"
    ;;
  -h|--help)
    usage
    exit 0
    ;;
  *)
    usage
    exit 2
    ;;
esac
if (($# > 1)); then
  usage
  exit 2
fi

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P) || exit 1
REPO_ROOT=$(git -C "$SCRIPT_DIR" rev-parse --show-toplevel 2>/dev/null) || {
  printf 'FAIL [git_repository] verifier is not inside a Git worktree\n' >&2
  exit 1
}
readonly REPO_ROOT

HEAD_COMMIT=$(git -C "$REPO_ROOT" rev-parse HEAD 2>/dev/null) || fail "git_head" "cannot resolve HEAD"
HEAD_TREE=$(git -C "$REPO_ROOT" rev-parse 'HEAD^{tree}' 2>/dev/null) || fail "git_tree" "cannot resolve HEAD tree"

CANONICAL_TIP=$(git -C "$REPO_ROOT" rev-parse --verify "${CANONICAL_REF}^{commit}" 2>/dev/null) || fail "canonical_tip" "cannot resolve $CANONICAL_REF"
[[ "$HEAD_COMMIT" == "$CANONICAL_TIP" ]] || \
  fail "canonical_tip" "HEAD $HEAD_COMMIT must exactly equal $CANONICAL_REF tip $CANONICAL_TIP"
pass "canonical_tip"

git -C "$REPO_ROOT" cat-file -e "${REQUIRED_ANCESTOR}^{commit}" 2>/dev/null || fail "required_ancestor_exists" "required commit is unavailable"
git -C "$REPO_ROOT" merge-base --is-ancestor "$REQUIRED_ANCESTOR" HEAD || fail "required_ancestor" "HEAD does not inherit 992223cbb"
pass "required_ancestor"

git -C "$REPO_ROOT" merge-base --is-ancestor "$LEGACY_COMPACT_COMMIT" HEAD || fail "legacy_compact_commit" "legacy compact v2 fix is not inherited"
pass "legacy_compact_commit"
git -C "$REPO_ROOT" merge-base --is-ancestor "$IMAGE_TASK_CORS_COMMIT" HEAD || fail "image_task_cors_commit" "Image task CORS fix is not inherited"
pass "image_task_cors_commit"

check_dirty_non_ui

require_file "migration_186_exists" "$MIGRATION_186"
actual_sha=$(sha256_file "$REPO_ROOT/$MIGRATION_186")
[[ "$actual_sha" == "$MIGRATION_186_SHA256" ]] || fail "migration_186_sha256" "expected $MIGRATION_186_SHA256, got $actual_sha"
pass "migration_186_sha256"
require_file "migration_194_exists" "$MIGRATION_194"
actual_sha=$(sha256_file "$REPO_ROOT/$MIGRATION_194")
[[ "$actual_sha" == "$MIGRATION_194_SHA256" ]] || fail "migration_194_sha256" "expected $MIGRATION_194_SHA256, got $actual_sha"
pass "migration_194_sha256"

# User daily reset route, request DTO, handler, service, repository port, and result DTO.
require_literal "user_reset_route" "backend/internal/server/routes/user.go" 'subscriptions.POST("/:id/reset-daily-quota", panelRateLimiter.DailyQuotaReset(), h.Subscription.ResetDailyQuota)'
require_literal "user_reset_request_dto" "backend/internal/handler/subscription_handler.go" 'type ResetDailyQuotaRequest struct {'
require_literal "user_reset_confirmation" "backend/internal/handler/subscription_handler.go" 'Confirm bool `json:"confirm"`'
require_literal "user_reset_handler" "backend/internal/handler/subscription_handler.go" 'func (h *SubscriptionHandler) ResetDailyQuota(c *gin.Context) {'
require_literal "user_reset_service" "backend/internal/service/user_daily_quota_reset.go" 'func (s *SubscriptionService) ResetUserDailyQuota(ctx context.Context, userID, subscriptionID int64, operationKeyHash string) (*UserDailyQuotaResetResult, error) {'
require_literal "user_reset_repository_port" "backend/internal/service/user_subscription_port.go" 'ResetDailyQuotaForUser(ctx context.Context, userID, subscriptionID int64, operationKeyHash string) (*UserDailyQuotaResetResult, error)'
require_literal "user_reset_result_dto" "backend/internal/service/user_subscription.go" 'type UserDailyQuotaResetResult struct {'
require_literal "user_reset_result_availability" "backend/internal/service/user_subscription.go" 'DailyQuotaResetAvailable   bool       `json:"daily_quota_reset_available"`'
require_literal "user_reset_native_cleared_window" "backend/internal/repository/user_subscription_repo.go" "'cleared_daily_window_start', COALESCE(us.daily_window_start, date_trunc('day', candidate.captured_at AT TIME ZONE 'Asia/Shanghai') AT TIME ZONE 'Asia/Shanghai')"
require_literal "user_reset_native_cleared_usage" "backend/internal/repository/user_subscription_repo.go" "'cleared_daily_usage_usd', GREATEST(us.daily_usage_usd, 0)"
reject_literal_in_block "user_reset_nullable_cleared_window_rejected" \
  "backend/internal/repository/user_subscription_repo.go" \
  'func (r *userSubscriptionRepository) ResetDailyQuotaForUser' \
  'func (r *userSubscriptionRepository) ResetUsageWindows' \
  "'cleared_daily_window_start', us.daily_window_start,"
reject_literal_in_block "user_reset_unclamped_cleared_usage_rejected" \
  "backend/internal/repository/user_subscription_repo.go" \
  'func (r *userSubscriptionRepository) ResetDailyQuotaForUser' \
  'func (r *userSubscriptionRepository) ResetUsageWindows' \
  "'cleared_daily_usage_usd', us.daily_usage_usd"
require_test "user_reset_native_ledger_sql_test" "backend/internal/repository/user_subscription_daily_quota_reset_test.go" "TestResetDailyQuotaForUserUsesAtomicOwnedEffectiveUpdate"
require_test "user_reset_null_window_native_ledger_test" "backend/internal/repository/user_subscription_daily_quota_reset_test.go" "TestResetDailyQuotaForUserNullWindowWritesShanghaiNativeLedgerWindow"
require_test "reset_today_null_window_native_operation_test" "backend/internal/service/subscription_reset_today_test.go" "TestResetTodayUsageTransactionAcceptsNativeSelfResetFromInitiallyNullDailyWindow"

# Admin reset-today route, request/result DTOs, handler, and transactional service.
require_literal "admin_reset_today_route" "backend/internal/server/routes/admin.go" 'subscriptions.POST("/reset-today-usage", h.Admin.Subscription.ResetTodayUsage)'
require_literal "admin_reset_today_request_dto" "backend/internal/handler/admin/subscription_handler.go" 'type ResetTodayUsageRequest struct {'
require_literal "admin_reset_today_handler" "backend/internal/handler/admin/subscription_handler.go" 'func (h *SubscriptionHandler) ResetTodayUsage(c *gin.Context) {'
require_literal "admin_reset_today_service" "backend/internal/service/subscription_reset_today.go" 'func (s *SubscriptionService) AdminResetTodayUsage(ctx context.Context, actorUserID int64, idempotencyKeyHash string) (*ResetTodayUsageResult, error) {'
require_literal "admin_reset_today_result_dto" "backend/internal/service/subscription_reset_today.go" 'type ResetTodayUsageResult struct {'
require_literal "admin_reset_today_timezone" "backend/internal/service/subscription_reset_today.go" 'resetTodayUsageTimezone       = subscriptionDailyCalendarTimezone'
require_literal "admin_reset_today_calendar_start" "backend/internal/service/subscription_reset_today.go" 'return subscriptionDailyCalendarStart(cutoff), nil'

# Admin per-subscription reset route, selectors, handler, and service must remain available.
require_literal "admin_reset_quota_route" "backend/internal/server/routes/admin.go" 'subscriptions.POST("/:id/reset-quota", h.Admin.Subscription.ResetQuota)'
require_literal "admin_reset_quota_request_dto" "backend/internal/handler/admin/subscription_handler.go" 'type ResetSubscriptionQuotaRequest struct {'
require_literal "admin_reset_quota_daily_selector" "backend/internal/handler/admin/subscription_handler.go" 'Daily   bool `json:"daily"`'
require_literal "admin_reset_quota_weekly_selector" "backend/internal/handler/admin/subscription_handler.go" 'Weekly  bool `json:"weekly"`'
require_literal "admin_reset_quota_monthly_selector" "backend/internal/handler/admin/subscription_handler.go" 'Monthly bool `json:"monthly"`'
require_literal "admin_reset_quota_handler" "backend/internal/handler/admin/subscription_handler.go" 'func (h *SubscriptionHandler) ResetQuota(c *gin.Context) {'
require_literal "admin_reset_quota_selector_validation" "backend/internal/handler/admin/subscription_handler.go" 'if !req.Daily && !req.Weekly && !req.Monthly {'
require_literal "admin_reset_quota_service" "backend/internal/service/subscription_service.go" 'func (s *SubscriptionService) AdminResetQuota(ctx context.Context, subscriptionID int64, resetDaily, resetWeekly, resetMonthly bool) (*UserSubscription, error) {'
require_literal "admin_reset_quota_repository_cutoff_port" "backend/internal/service/user_subscription_port.go" 'ResetUsageWindows(ctx context.Context, id int64, resetDaily, resetWeekly, resetMonthly bool, dailyStart, periodicStart time.Time) (time.Time, error)'
require_literal_in_block "admin_reset_quota_database_cutoff_sql" \
  "backend/internal/repository/user_subscription_repo.go" \
  'func (r *userSubscriptionRepository) ResetUsageWindows' \
  'func (r *userSubscriptionRepository) ResetDailyUsage' \
  "updated_at = GREATEST(clock_timestamp(), us.updated_at + interval '1 microsecond')"
require_literal_in_block "admin_reset_quota_database_cutoff_returning" \
  "backend/internal/repository/user_subscription_repo.go" \
  'func (r *userSubscriptionRepository) ResetUsageWindows' \
  'func (r *userSubscriptionRepository) ResetDailyUsage' \
  'RETURNING us.updated_at'
reject_literal_in_block "admin_reset_quota_application_cutoff_rejected" \
  "backend/internal/repository/user_subscription_repo.go" \
  'func (r *userSubscriptionRepository) ResetUsageWindows' \
  'func (r *userSubscriptionRepository) ResetDailyUsage' \
  'updated_at = $'
require_literal_in_block "admin_reset_quota_barrier_path" \
  "backend/internal/service/subscription_service.go" \
  'func (s *SubscriptionService) AdminResetQuota' \
  'func (s *SubscriptionService) CheckAndResetWindows' \
  'cacheErr := s.invalidateSubscriptionResetCaches(cacheCtx, sub.UserID, sub.GroupID, resetCutoff)'
require_literal_in_block "admin_reset_quota_retry_path" \
  "backend/internal/service/subscription_service.go" \
  'func (s *SubscriptionService) AdminResetQuota' \
  'func (s *SubscriptionService) CheckAndResetWindows' \
  's.scheduleSubscriptionResetCacheRetries(sub.UserID, sub.GroupID, resetCutoff)'
require_literal "subscription_reset_cache_timeout" "backend/internal/service/user_daily_quota_reset.go" 'const subscriptionQuotaCacheTimeout = 5 * time.Second'
require_literal "subscription_reset_retry_schedule" "backend/internal/service/user_daily_quota_reset.go" 'var subscriptionQuotaCacheRetryDelays = []time.Duration{'
require_literal "subscription_reset_retry_5s" "backend/internal/service/user_daily_quota_reset.go" '5 * time.Second,'
require_literal "subscription_reset_retry_12s" "backend/internal/service/user_daily_quota_reset.go" '12 * time.Second,'
require_literal "subscription_reset_retry_30s" "backend/internal/service/user_daily_quota_reset.go" '30 * time.Second,'
require_literal "subscription_reset_retry_1m" "backend/internal/service/user_daily_quota_reset.go" '1 * time.Minute,'
require_literal "subscription_reset_retry_2m" "backend/internal/service/user_daily_quota_reset.go" '2 * time.Minute,'
require_literal "subscription_reset_retry_4m" "backend/internal/service/user_daily_quota_reset.go" '4 * time.Minute,'
require_literal "subscription_reset_retry_6m" "backend/internal/service/user_daily_quota_reset.go" '6 * time.Minute,'
require_literal "subscription_reset_retry_generation" "backend/internal/service/user_daily_quota_reset.go" 'generation uint64'
require_literal "subscription_reset_retry_generation_increment" "backend/internal/service/user_daily_quota_reset.go" 'state.generation++'
require_literal "subscription_reset_retry_generation_cycle" "backend/internal/service/user_daily_quota_reset.go" '_, cycleGeneration := s.userDailyQuotaCacheRetrySnapshot(state)'
require_literal "subscription_reset_retry_generation_finish" "backend/internal/service/user_daily_quota_reset.go" 'if state.generation != cycleGeneration || (!lastAttemptedCutoff.IsZero() && state.resetAt.After(lastAttemptedCutoff)) {'
require_literal_in_block "admin_reset_quota_local_l1_sync" \
  "backend/internal/service/user_daily_quota_reset.go" \
  'func (s *SubscriptionService) invalidateSubscriptionResetCaches' \
  'func (s *SubscriptionService) invalidateUserDailyQuotaCaches' \
  's.InvalidateSubCacheSync(userID, groupID)'
require_literal_in_block "admin_reset_quota_redis_barrier" \
  "backend/internal/service/user_daily_quota_reset.go" \
  'func (s *SubscriptionService) invalidateSubscriptionResetCaches' \
  'func (s *SubscriptionService) invalidateUserDailyQuotaCaches' \
  's.billingCacheService.InvalidateSubscriptionAt(ctx, userID, groupID, resetAt)'
require_literal_in_block "admin_reset_quota_cross_node_publish" \
  "backend/internal/service/user_daily_quota_reset.go" \
  'func (s *SubscriptionService) invalidateSubscriptionResetCaches' \
  'func (s *SubscriptionService) invalidateUserDailyQuotaCaches' \
  's.billingCacheService.PublishSubscriptionCacheInvalidation(ctx, subCacheKey(userID, groupID))'
reject_literal_in_block "admin_reset_quota_helper_unbarriered_delete_rejected" \
  "backend/internal/service/user_daily_quota_reset.go" \
  'func (s *SubscriptionService) invalidateSubscriptionResetCaches' \
  'func (s *SubscriptionService) invalidateUserDailyQuotaCaches' \
  'InvalidateSubscription(ctx'
reject_literal_in_block "admin_reset_quota_unbarriered_delete_rejected" \
  "backend/internal/service/subscription_service.go" \
  'func (s *SubscriptionService) AdminResetQuota' \
  'func (s *SubscriptionService) CheckAndResetWindows' \
  'InvalidateSubscription(ctx'
require_literal "reset_today_cache_workers" "backend/internal/service/subscription_reset_today.go" 'resetTodayUsageCacheWorkers   = 16'
require_literal "reset_today_cache_batch_timeout" "backend/internal/service/subscription_reset_today.go" 'resetTodayUsageCacheBatchTTL  = 15 * time.Second'
require_literal_in_block "reset_today_cache_l1_all_del" \
  "backend/internal/service/subscription_reset_today.go" \
  'func (s *SubscriptionService) invalidateResetTodayUsageTargetBatchWithTimeouts' \
  'func (s *SubscriptionService) scheduleResetTodayUsageCacheRetries' \
  's.InvalidateSubCache(target.UserID, target.GroupID)'
require_literal_in_block "subscription_cache_l1_del_primitive" \
  "backend/internal/service/subscription_service.go" \
  'func (s *SubscriptionService) InvalidateSubCache(userID, groupID int64) {' \
  'func (s *SubscriptionService) InvalidateSubCacheSync' \
  's.subCacheL1.Del(subCacheKey(userID, groupID))'
reject_literal_in_block "subscription_cache_l1_del_primitive_wait_rejected" \
  "backend/internal/service/subscription_service.go" \
  'func (s *SubscriptionService) InvalidateSubCache(userID, groupID int64) {' \
  'func (s *SubscriptionService) InvalidateSubCacheSync' \
  '.Wait()'
require_literal_in_block "reset_today_cache_l1_single_wait" \
  "backend/internal/service/subscription_reset_today.go" \
  'func (s *SubscriptionService) invalidateResetTodayUsageTargetBatchWithTimeouts' \
  'func (s *SubscriptionService) scheduleResetTodayUsageCacheRetries' \
  's.subCacheL1.Wait()'
reject_literal_in_block "reset_today_cache_per_key_wait_rejected" \
  "backend/internal/service/subscription_reset_today.go" \
  'func (s *SubscriptionService) invalidateResetTodayUsageTargetBatchWithTimeouts' \
  'func (s *SubscriptionService) scheduleResetTodayUsageCacheRetries' \
  'InvalidateSubCacheSync('
require_literal_in_block "reset_today_cache_bounded_workers" \
  "backend/internal/service/subscription_reset_today.go" \
  'func (s *SubscriptionService) invalidateResetTodayUsageTargetBatchWithTimeouts' \
  'func (s *SubscriptionService) scheduleResetTodayUsageCacheRetries' \
  'workerCount := min(resetTodayUsageCacheWorkers, len(targets))'
require_literal_in_block "reset_today_cache_shared_batch_context" \
  "backend/internal/service/subscription_reset_today.go" \
  'func (s *SubscriptionService) invalidateResetTodayUsageTargetBatchWithTimeouts' \
  'func (s *SubscriptionService) scheduleResetTodayUsageCacheRetries' \
  'batchCtx, batchCancel := context.WithTimeout(context.Background(), batchTimeout)'
require_literal_in_block "reset_today_cache_per_target_context" \
  "backend/internal/service/subscription_reset_today.go" \
  'func (s *SubscriptionService) invalidateResetTodayUsageTargetBatchWithTimeouts' \
  'func (s *SubscriptionService) scheduleResetTodayUsageCacheRetries' \
  'ctx, cancel := context.WithTimeout(batchCtx, targetTimeout)'
require_literal_in_block "reset_today_cache_deadline_fast_failure" \
  "backend/internal/service/subscription_reset_today.go" \
  'func (s *SubscriptionService) invalidateResetTodayUsageTargetBatchWithTimeouts' \
  'func (s *SubscriptionService) scheduleResetTodayUsageCacheRetries' \
  'if err := batchCtx.Err(); err != nil {'
require_literal_in_block "reset_today_cache_failed_subset_retry" \
  "backend/internal/service/subscription_reset_today.go" \
  'func (s *SubscriptionService) AdminResetTodayUsage' \
  'func validResetTodayUsageKeyHash' \
  's.scheduleResetTodayUsageCacheRetries(failedTargets, op.Result.CutoffAt)'
reject_literal_in_block "reset_today_cache_all_targets_retry_rejected" \
  "backend/internal/service/subscription_reset_today.go" \
  'func (s *SubscriptionService) AdminResetTodayUsage' \
  'func validResetTodayUsageKeyHash' \
  's.scheduleResetTodayUsageCacheRetries(op.Targets, op.Result.CutoffAt)'
require_test "admin_reset_quota_route_test" "backend/internal/server/routes/admin_subscription_reset_today_route_test.go" "TestRegisterSubscriptionRoutesIncludesResetQuota"
require_test "admin_reset_quota_dto_test" "backend/internal/handler/admin/subscription_reset_quota_handler_test.go" "TestResetSubscriptionQuotaRequestDecodesDailyWeeklyMonthlySelectors"
require_test "admin_reset_quota_handler_selector_test" "backend/internal/handler/admin/subscription_reset_quota_handler_test.go" "TestResetQuotaForwardsDailyWeeklyMonthlySelectors"
require_test "admin_reset_quota_handler_empty_test" "backend/internal/handler/admin/subscription_reset_quota_handler_test.go" "TestResetQuotaRequiresAtLeastOneSelectedWindow"
require_test "admin_reset_quota_daily_service_test" "backend/internal/service/subscription_reset_quota_test.go" "TestAdminResetQuota_ResetDailyOnly"
require_test "admin_reset_quota_daily_calendar_service_test" "backend/internal/service/subscription_reset_quota_test.go" "TestAdminResetQuota_ResetBoth"
require_test "admin_reset_quota_weekly_service_test" "backend/internal/service/subscription_reset_quota_test.go" "TestAdminResetQuota_ResetWeeklyOnly"
require_test "admin_reset_quota_monthly_service_test" "backend/internal/service/subscription_reset_quota_test.go" "TestAdminResetQuota_ResetMonthlyOnly"
require_test "admin_reset_quota_database_cutoff_service_test" "backend/internal/service/subscription_reset_quota_test.go" "TestAdminResetQuotaUsesDatabaseCutoffBarrierPublishesAndInvalidatesLocalL1"
require_test "admin_reset_quota_cache_retry_test" "backend/internal/service/subscription_reset_quota_test.go" "TestAdminResetQuotaCacheFailureDoesNotFailCommittedResetAndRetries"
require_test "subscription_reset_cache_generation_test" "backend/internal/service/user_daily_quota_reset_test.go" "TestSubscriptionQuotaCacheRetryNewGenerationGetsFullCycleAfterLastFailure"
require_test "reset_today_cache_failed_subset_ttl_test" "backend/internal/service/subscription_reset_today_test.go" "TestResetTodayCacheRetryUsesOnlyFailedSubsetAndRunsPastCacheTTL"
require_test "reset_today_cache_slow_target_test" "backend/internal/service/subscription_reset_today_test.go" "TestResetTodayCacheBatchClearsAllL1BeforeSlowRemoteAndDoesNotStarveFastTarget"
require_test "reset_today_cache_batch_deadline_test" "backend/internal/service/subscription_reset_today_test.go" "TestResetTodayCacheBatchHasBoundedConcurrencyAndSharedDeadline"
require_test "admin_reset_quota_postcommit_refresh_test" "backend/internal/service/subscription_reset_quota_test.go" "TestAdminResetQuotaRefreshFailureReturnsCommittedSnapshot"
require_test "admin_reset_quota_database_cutoff_repository_test" "backend/internal/repository/user_subscription_reset_usage_test.go" "TestResetUsageWindowsReturnsStrictlyMonotonicDatabaseCutoff"

# Natural-day and one-day-card behavior must remain explicit in implementation and tests.
require_literal "subscription_calendar_timezone" "backend/internal/service/subscription_calendar.go" 'const subscriptionDailyCalendarTimezone = "Asia/Shanghai"'
require_literal "subscription_calendar_location" "backend/internal/service/subscription_calendar.go" 'var subscriptionDailyCalendarLocation = mustLoadSubscriptionDailyCalendarLocation()'
require_literal "calendar_day_implementation" "backend/internal/service/user_subscription.go" 'windowStart := subscriptionDailyCalendarStart(now)'
require_literal "one_day_card_implementation" "backend/internal/service/user_subscription.go" 's.DailyWindowStart == nil || s.HasOneTimeDailyQuota()'
require_literal "calendar_activation_path" "backend/internal/service/subscription_service.go" 'ActivateWindows(ctx, sub.ID, subscriptionDailyCalendarStart(now), now)'
require_literal "calendar_admin_reset_path" "backend/internal/service/subscription_service.go" 'dailyStart := subscriptionDailyCalendarStart(windowStart)'
require_literal_in_block "calendar_check_and_reset_path" \
  "backend/internal/service/subscription_service.go" \
  'func (s *SubscriptionService) CheckAndResetWindows' \
  'func (s *SubscriptionService) EnsureWindowMaintenance' \
  'if windowStart, ok := sub.automaticDailyWindowStartAt(now); ok {'
require_literal_in_block "calendar_normalized_copy_path" \
  "backend/internal/service/subscription_service.go" \
  'func normalizedDailyWindowCopy' \
  'func (s *SubscriptionService) GetActiveSubscription' \
  'if windowStart, ok := normalized.automaticDailyWindowStartAt(now); ok {'
require_literal_in_block "calendar_normalize_list_path" \
  "backend/internal/service/subscription_service.go" \
  'func normalizeExpiredWindowsAt' \
  'func normalizeSubscriptionStatus' \
  'if windowStart, ok := sub.automaticDailyWindowStartAt(now); ok {'
require_literal_in_block "calendar_validate_path" \
  "backend/internal/service/subscription_service.go" \
  'func (s *SubscriptionService) ValidateAndCheckLimits' \
  'func (s *SubscriptionService) DoWindowMaintenance' \
  'if sub.canAutomaticallyResetDailyAt(now) {'
require_literal_in_block "calendar_get_by_id_path" \
  "backend/internal/service/subscription_service.go" \
  'func (s *SubscriptionService) GetByID' \
  'func subscriptionServiceNow' \
  'return normalizedDailyWindowCopy(sub, subscriptionServiceNow(s)), nil'
require_literal_in_block "calendar_progress_path" \
  "backend/internal/service/subscription_service.go" \
  'func (s *SubscriptionService) calculateProgress' \
  'func (s *SubscriptionService) GetUserSubscriptionsWithProgress' \
  'sub = normalizedDailyWindowCopy(sub, subscriptionServiceNow(s))'
require_literal_in_block "calendar_subscription_cache_path" \
  "backend/internal/service/billing_cache_service.go" \
  'func (s *BillingCacheService) normalizeSubscriptionCacheForWindow' \
  'func (d *subscriptionCacheData) hasWindowMetadata' \
  'windowStart := subscriptionDailyCalendarStart(now)'
require_literal "calendar_response_path" "backend/internal/handler/dto/mappers.go" 'DailyWindowResetsAt:        sub.DailyResetTime(),'
require_literal "calendar_utc_global_test_setup" "backend/internal/service/subscription_daily_calendar_window_test.go" 'require.Equal(t, "UTC", timezone.Name())'
require_literal "calendar_response_utc_global_test_setup" "backend/internal/handler/dto/user_subscription_mapper_reset_time_test.go" 'require.NoError(t, timezone.Init("UTC"))'
reject_literal "calendar_model_global_start_rejected" "backend/internal/service/user_subscription.go" 'timezone.StartOfDay'
reject_literal "calendar_model_global_location_rejected" "backend/internal/service/user_subscription.go" 'timezone.Location'
reject_literal "calendar_model_time_local_rejected" "backend/internal/service/user_subscription.go" 'time.Local'
reject_literal "calendar_service_global_start_rejected" "backend/internal/service/subscription_service.go" 'timezone.StartOfDay'
reject_literal "calendar_service_global_location_rejected" "backend/internal/service/subscription_service.go" 'timezone.Location'
reject_literal "calendar_service_time_local_rejected" "backend/internal/service/subscription_service.go" 'time.Local'
reject_literal "calendar_reset_today_global_start_rejected" "backend/internal/service/subscription_reset_today.go" 'timezone.StartOfDay'
reject_literal "calendar_reset_today_global_location_rejected" "backend/internal/service/subscription_reset_today.go" 'timezone.Location'
reject_literal "calendar_reset_today_time_local_rejected" "backend/internal/service/subscription_reset_today.go" 'time.Local'
reject_literal "calendar_reset_today_location_reload_rejected" "backend/internal/service/subscription_reset_today.go" 'time.LoadLocation(resetTodayUsageTimezone)'
reject_literal_in_block "calendar_subscription_cache_global_start_rejected" \
  "backend/internal/service/billing_cache_service.go" \
  'func (s *BillingCacheService) normalizeSubscriptionCacheForWindow' \
  'func (d *subscriptionCacheData) hasWindowMetadata' \
  'timezone.StartOfDay'
require_test "calendar_midnight_test" "backend/internal/service/subscription_daily_calendar_window_test.go" "TestSubscriptionDailyCalendarBoundaryUsesShanghaiMidnight"
require_test "manual_reset_midnight_test" "backend/internal/service/subscription_daily_calendar_window_test.go" "TestSubscriptionDailyResetTimeUsesNextCalendarMidnightAfterManualReset"
require_test "calendar_check_and_reset_test" "backend/internal/service/subscription_daily_calendar_window_test.go" "TestCheckAndResetWindowsAdvancesDailyAcrossMultipleDaysOnly"
require_test "calendar_normalize_test" "backend/internal/service/subscription_daily_calendar_window_test.go" "TestNormalizeExpiredWindowsAdvancesDailyResponseCopy"
require_test "calendar_validate_test" "backend/internal/service/subscription_daily_calendar_window_test.go" "TestValidateAndCheckLimitsUsesDailyCalendarBoundary"
require_test "calendar_progress_test" "backend/internal/service/subscription_daily_calendar_window_test.go" "TestCalculateProgressUsesDailyCalendarMidnight"
require_test "calendar_get_by_id_test" "backend/internal/service/subscription_daily_calendar_window_test.go" "TestGetByIDNormalizesDailyResponseCopyWithoutMutatingRepositorySnapshot"
require_test "dto_calendar_midnight_test" "backend/internal/handler/dto/user_subscription_mapper_reset_time_test.go" "TestUserSubscriptionMappersExposeEffectiveWindowResetTimes"
require_test "one_day_card_test" "backend/internal/service/user_subscription_daily_quota_test.go" "TestUserSubscriptionNeedsDailyReset_DailyCardKeepsOneTimeQuota"
require_test "multi_day_reset_test" "backend/internal/service/user_subscription_daily_quota_test.go" "TestUserSubscriptionNeedsDailyReset_MultiDaySubscriptionStillRefreshes"
require_test "admin_midnight_boundary_test" "backend/internal/service/subscription_reset_today_boundary_test.go" "TestResetTodayUsageWindowClassificationAtShanghaiMidnight"

# Cache reset barriers and database-derived usage write timestamps prevent stale replay.
require_literal "reset_barrier_port" "backend/internal/service/billing_cache_service.go" 'SetSubscriptionResetBarrierAndInvalidate(ctx context.Context, userID, groupID int64, cutoff time.Time) error'
require_literal "reset_barrier_implementation" "backend/internal/repository/billing_cache.go" 'func (c *billingCache) SetSubscriptionResetBarrierAndInvalidate(ctx context.Context, userID, groupID int64, cutoff time.Time) error {'
require_literal "usage_write_timestamp_sql" "backend/internal/repository/usage_billing_repo.go" "updated_at = GREATEST(clock_timestamp(), us.updated_at + interval '1 microsecond')"
require_literal "usage_write_timestamp_forwarding" "backend/internal/service/gateway_usage_billing.go" '*result.SubscriptionUsageWriteAt,'
require_test "reset_barrier_test" "backend/internal/repository/billing_cache_subscription_barrier_test.go" "TestBillingCacheSubscriptionBarrierRejectsOldIncrementAndAllowsNewIncrement"
require_test "usage_write_timestamp_test" "backend/internal/repository/usage_billing_subscription_timestamp_unit_test.go" "TestApplyUsageBillingEffectsUsesStrictlyMonotonicSubscriptionDatabaseWriteTime"

# Image Studio download and Image2 behavior are production contracts. The download
# route must authenticate before its own 8 KiB parser and must not pass through a
# middleware that can pre-read the larger gateway body allowance.
require_file "image_download_service_file" "backend/internal/service/image_download.go"
require_file "image_download_network_file" "backend/internal/service/image_download_network.go"
require_file "image_download_handler_file" "backend/internal/handler/image_download_handler.go"
require_literal "image_download_route" "backend/internal/server/routes/gateway.go" 'gateway.POST("/images/download", imageDownloadHandler.Download)'
reject_literal "image_download_root_alias_rejected" "backend/internal/server/routes/gateway.go" 'r.POST("/images/download",'
require_literal_in_block "image_download_pre_composite_route" \
  "backend/internal/server/routes/gateway.go" \
  'gateway.Use(gin.HandlerFunc(apiKeyAuth))' \
  'gateway.Use(compositeTarget)' \
  'gateway.POST("/images/download", imageDownloadHandler.Download)'
require_literal "image_download_allowed_host" "backend/internal/service/image_download.go" 'pre-signed-firefly-prod.s3-accelerate.amazonaws.com'
require_literal "image_download_max_bytes" "backend/internal/service/image_download.go" 'ImageDownloadMaxBytes    = 64 << 20'
require_literal "image_download_five_minute_client" "backend/internal/service/image_download.go" 'client = newImageDownloadHTTPClient(5 * time.Minute)'
require_literal "image_download_redirect_validation" "backend/internal/service/image_download.go" 'return validateImageDownloadURL(req.URL)'
require_literal "image_download_stream_limit" "backend/internal/service/image_download.go" 'io.LimitReader(resp.Body, remaining)'
require_literal "image_download_strict_dial" "backend/internal/service/image_download_network.go" 'DialContext:           imageDownloadSafeDialContext,'
require_literal "image_download_reject_mixed_dns" "backend/internal/service/image_download_network.go" 'Reject the whole resolution set if any answer is special-use.'
require_literal "image_download_benchmark_prefix" "backend/internal/service/image_download_network.go" 'netip.MustParsePrefix("198.18.0.0/15")'
require_literal "image_download_multicast_prefix" "backend/internal/service/image_download_network.go" 'netip.MustParsePrefix("224.0.0.0/4")'
require_literal "image_download_reserved_prefix" "backend/internal/service/image_download_network.go" 'netip.MustParsePrefix("240.0.0.0/4")'
require_literal "image_download_ipv6_multicast_prefix" "backend/internal/service/image_download_network.go" 'netip.MustParsePrefix("ff00::/8")'
require_literal "image_download_ipv6_dummy_prefix" "backend/internal/service/image_download_network.go" 'netip.MustParsePrefix("100:0:0:1::/64")'
require_literal "image_download_request_max_bytes" "backend/internal/handler/image_download_handler.go" 'imageDownloadRequestMaxBytes = 8 << 10'
require_literal "image_download_copy_buffer" "backend/internal/handler/image_download_handler.go" 'imageDownloadCopyBufferBytes = 32 << 10'
require_literal "image_download_process_concurrency" "backend/internal/handler/image_download_handler.go" 'imageDownloadMaxConcurrent   = 16'
require_literal "image_download_key_concurrency" "backend/internal/handler/image_download_handler.go" 'imageDownloadMaxPerAPIKey    = 2'
require_literal "image_download_key_rate" "backend/internal/handler/image_download_handler.go" 'imageDownloadMaxPerKeyMinute = 30'
require_literal "image_download_process_rate" "backend/internal/handler/image_download_handler.go" 'imageDownloadMaxTotalMinute  = 240'
require_literal "image_download_skip_billing" "backend/internal/server/middleware/api_key_auth.go" 'isImageDownloadRead(c.Request.Method, c.Request.URL.Path)'
require_literal "image_download_exact_billing_predicate" "backend/internal/server/middleware/api_key_auth.go" 'return method == http.MethodPost && path == "/v1/images/download"'
require_literal "image2_exact_size_guard" "backend/internal/service/openai_images.go" 'if strings.EqualFold(strings.TrimSpace(model), "gpt-image-2") {'
require_literal "image2_model_aware_rewrite" "backend/internal/service/openai_images.go" 'size = normalizeOpenAIImageSizeForUpstream(size, model)'
require_literal "image2_type_url_alias" "backend/internal/service/openai_images.go" 'if typeValue == "url" || typeValue == "b64_json" {'
require_literal "image2_exact_size_capability" "backend/internal/service/openai_images.go" 'OpenAIImagesCapabilityExactSize OpenAIImagesCapability = "images-exact-size"'
require_literal "image2_exact_size_classifier" "backend/internal/service/openai_images.go" 'if requiresOpenAIImagesExactSize(req) {'
require_literal "image2_exact_size_api_key_capability" "backend/internal/service/account.go" 'case OpenAIImagesCapabilityExactSize:'
require_literal "image2_exact_size_opt_in_key" "backend/internal/service/account.go" 'const OpenAIImageExactSizeSupportedExtraKey = "openai_image_exact_size_supported"'
require_literal "image2_exact_size_opt_in_check" "backend/internal/service/account.go" 'supported, ok := a.Extra[OpenAIImageExactSizeSupportedExtraKey].(bool)'
require_literal "image2_exact_size_account_mapping_filter" "backend/internal/service/account.go" 'a.GetMappedModel("gpt-image-2")'
require_literal "image2_exact_size_account_rejection" "backend/internal/service/openai_images.go" 'return nil, fmt.Errorf("exact image size requires an opted-in API key image account")'
require_literal "image2_exact_size_forward_opt_in" "backend/internal/service/openai_images.go" '!account.SupportsOpenAIImageCapability(OpenAIImagesCapabilityExactSize)'
require_literal "image2_exact_size_mapping_rejection" "backend/internal/service/openai_images.go" 'return nil, fmt.Errorf("exact image size requires upstream model gpt-image-2")'
require_literal "image2_effective_model_capability" "backend/internal/service/openai_images.go" 'func (r *OpenAIImagesRequest) ResolveRequiredCapabilityForEffectiveModel(effectiveModel string) error {'
require_literal "image2_handler_effective_model_capability" "backend/internal/handler/openai_images.go" 'parsed.ResolveRequiredCapabilityForEffectiveModel(effectiveImageModel)'
require_literal "image2_handler_effective_model_scheduler" "backend/internal/handler/openai_images.go" 'effectiveImageModel,'
require_literal "image2_handler_effective_model_forward" "backend/internal/handler/openai_images.go" 'body, parsed, forwardMappedModel)'
require_literal "image2_scheduler_metadata_projection" "backend/internal/repository/scheduler_cache.go" 'service.OpenAIImageExactSizeSupportedExtraKey,'
require_test "image_download_route_test" "backend/internal/server/routes/image_download_route_test.go" "TestGatewayRoutesImageDownloadPathIsRegisteredWithoutRootAlias"
require_test "image_download_auth_route_test" "backend/internal/server/routes/image_download_route_test.go" "TestGatewayRoutesImageDownloadRequiresAPIKeyAuthentication"
require_test "image_download_route_order_test" "backend/internal/server/routes/image_download_route_order_test.go" "TestGatewayImageDownloadRouteRunsBeforeCompositeBodyParsing"
require_literal "image_download_auth_test_unit_tag" "backend/internal/server/middleware/api_key_auth_image_download_test.go" "//go:build unit"
require_test "image_download_skip_billing_test" "backend/internal/server/middleware/api_key_auth_image_download_test.go" "TestAPIKeyAuthImageDownloadSkipsBillingChecks"
require_test "image_download_exact_predicate_test" "backend/internal/server/middleware/api_key_auth_image_download_test.go" "TestIsImageDownloadReadRequiresExactMethodAndPath"
require_test "image_download_cors_test" "backend/internal/server/middleware/image_download_cors_test.go" "TestCORSPreflightAllowsImageDownloadAuthorization"
require_test "image_download_stream_test" "backend/internal/handler/image_download_handler_test.go" "TestImageDownloadHandlerStreamsAttachment"
require_test "image_download_request_limit_test" "backend/internal/handler/image_download_request_limit_test.go" "TestImageDownloadHandlerRejectsOversizedRequestBeforeFetch"
require_test "image_download_parallel_limit_test" "backend/internal/handler/image_download_handler_test.go" "TestImageDownloadConcurrencyGateParallelAcquireAndIdempotentRelease"
require_test "image_download_truncated_stream_test" "backend/internal/handler/image_download_handler_test.go" "TestImageDownloadHandlerRecordsTruncatedStream"
require_test "image_download_url_test" "backend/internal/service/image_download_test.go" "TestImageDownloadServiceRejectsUntrustedURLsBeforeRequest"
require_test "image_download_redirect_test" "backend/internal/service/image_download_test.go" "TestImageDownloadServiceRejectsRedirectToUntrustedHost"
require_test "image_download_length_test" "backend/internal/service/image_download_test.go" "TestImageDownloadServiceRejectsMissingAndOversizedContentLength"
require_test "image_download_type_test" "backend/internal/service/image_download_test.go" "TestImageDownloadServiceRejectsInvalidOrMismatchedContentType"
require_test "image_download_private_dns_test" "backend/internal/service/image_download_network_test.go" "TestImageDownloadSafeDialRejectsPrivateDNSAnswerBeforeDial"
require_test "image_download_mixed_dns_test" "backend/internal/service/image_download_network_test.go" "TestImageDownloadSafeDialRejectsMixedPublicPrivateAnswers"
require_test "image_download_special_ip_test" "backend/internal/service/image_download_network_test.go" "TestImageDownloadPublicIPPolicyRejectsPrivateAndSpecialUse"
require_test "image_download_redirect_limit_test" "backend/internal/service/image_download_ssrf_test.go" "TestImageDownloadServiceStopsAfterRedirectLimit"
require_test "image2_exact_size_test" "backend/internal/service/openai_images_image2_upstream_size_test.go" "TestRewriteOpenAIImagesRequest_GPTImage2PreservesExactSizeAndImage2Fields"
require_test "image2_legacy_size_test" "backend/internal/service/openai_images_image2_upstream_size_test.go" "TestRewriteOpenAIImagesRequest_LegacyModelStillNormalizesCustomSize"
require_test "image2_multipart_size_test" "backend/internal/service/openai_images_image2_multipart_test.go" "TestRewriteOpenAIImagesMultipartRequest_GPTImage2PreservesExactSize"
require_test "image2_json_edit_type_url_test" "backend/internal/service/openai_images_image2_edits_test.go" "TestParseOpenAIImagesRequest_JSONEditTypeURLAliasesResponseFormat"
require_test "image2_json_edit_apikey_exact_size_test" "backend/internal/service/openai_images_image2_edits_test.go" "TestOpenAIGatewayServiceForwardImages_APIKeyJSONEditPreservesExactSizeAndNestedImage"
require_test "image2_json_edit_oauth_exact_size_rejection_test" "backend/internal/service/openai_images_image2_edits_test.go" "TestOpenAIGatewayServiceForwardImages_OAuthJSONEditExactSizeFailsBeforeUpstream"
require_test "image2_json_edit_oauth_standard_size_test" "backend/internal/service/openai_images_image2_edits_test.go" "TestOpenAIGatewayServiceForwardImages_OAuthJSONEditStandardSizeUsesTypeURLAlias"
require_test "image2_json_edit_response_format_precedence_test" "backend/internal/service/openai_images_image2_edits_test.go" "TestParseOpenAIImagesRequest_ResponseFormatWinsTypeAlias"
require_test "image2_json_edit_mapping_rejection_test" "backend/internal/service/openai_images_image2_edits_test.go" "TestOpenAIGatewayServiceForwardImages_ExactSizeRejectsLegacyModelMapping"
require_test "image2_unmarked_apikey_rejection_test" "backend/internal/service/openai_images_image2_edits_test.go" "TestOpenAIGatewayServiceForwardImages_UnmarkedAPIKeyExactSizeFailsBeforeUpstream"
require_test "image2_exact_size_stream_preservation_test" "backend/internal/service/openai_images_image2_edits_test.go" "TestOpenAIGatewayServiceForwardImages_APIKeyExactSizeStreamingPreservesRequest"
require_test "image2_effective_model_capability_test" "backend/internal/service/openai_images_image2_edits_test.go" "TestOpenAIImagesRequestResolveRequiredCapabilityForEffectiveModel"
require_test "image2_exact_size_account_capability_test" "backend/internal/service/openai_images_exact_size_scheduler_test.go" "TestAccountSupportsOpenAIImageCapability_ExactSizeRequiresAPIKey"
require_test "image2_exact_size_scheduler_test" "backend/internal/service/openai_images_exact_size_scheduler_test.go" "TestOpenAIGatewayService_SelectAccountForExactImageSizeNeverFallsBackToOAuth"
require_test "image2_oauth_actual_size_test" "backend/internal/service/openai_images_actual_size_test.go" "TestOpenAIGatewayServiceForwardImages_OAuthUsesDecodedOutputDimensions"
require_test "image2_oauth_stream_actual_size_test" "backend/internal/service/openai_images_actual_size_test.go" "TestOpenAIGatewayServiceForwardImages_OAuthStreamingUsesDecodedOutputDimensions"
require_test "image2_scheduler_metadata_round_trip_test" "backend/internal/repository/scheduler_cache_test.go" "TestSchedulerMetadataAccountKeepsOpenAIImageExactSizeCapability"
require_test "image2_scheduler_outbox_test" "backend/internal/repository/account_repo_image_exact_size_extra_test.go" "TestShouldEnqueueSchedulerOutboxForExtraUpdates_ImageExactSizeCapabilityIsRelevant"

# Keep compatibility fixes that were historically lost in UI-oriented releases.
require_literal "legacy_compact_adapter" "backend/internal/service/openai_gateway_request_body.go" 'func normalizeOpenAILegacyCompactRemoteV2Body(body []byte) ([]byte, bool, error) {'
require_literal "legacy_compact_feature_definition" "backend/internal/service/openai_compact_body_signal.go" 'const openAIRemoteCompactionV2Feature = "remote_compaction_v2"'
require_literal "legacy_compact_forward_feature" "backend/internal/service/openai_gateway_forward.go" 'ensureOpenAIBetaFeature(req.Header.Values("x-codex-beta-features"), openAIRemoteCompactionV2Feature)'
require_literal "legacy_compact_passthrough_feature" "backend/internal/service/openai_gateway_passthrough.go" 'ensureOpenAIBetaFeature(req.Header.Values("x-codex-beta-features"), openAIRemoteCompactionV2Feature)'
legacy_compact_feature_count=$(
  grep -RhoF -- 'const openAIRemoteCompactionV2Feature = "remote_compaction_v2"' \
    "$REPO_ROOT/backend/internal/service" | wc -l | tr -d '[:space:]'
)
[[ "$legacy_compact_feature_count" == "1" ]] ||
  fail "legacy_compact_feature_unique" "expected exactly one remote compaction v2 feature declaration, got $legacy_compact_feature_count"
pass "legacy_compact_feature_unique"
require_test "legacy_compact_test" "backend/internal/service/openai_legacy_compact_remote_v2_test.go" "TestOpenAIGatewayServiceForwardOAuthLegacyCompactUsesRemoteV2AndReturnsUnaryJSON"
require_literal "image_task_cors_header" "backend/internal/server/middleware/cors.go" '"X-Sub2API-Image-Task-Key",'
require_literal "image_account_policy_cors_header" "backend/internal/server/middleware/cors.go" '"X-Studio-Account-Policy",'
require_test "image_task_cors_test" "backend/internal/server/middleware/cors_test.go" "TestCORS_PreflightAllowsStudioImageTaskHeader"
require_literal "image_cors_test_requests_both_headers" "backend/internal/server/middleware/cors_test.go" '"authorization,content-type,x-studio-account-policy,x-sub2api-image-task-key"'
require_literal "image_cors_test_asserts_account_policy" "backend/internal/server/middleware/cors_test.go" 'assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "X-Studio-Account-Policy")'

# Codex account identity and upstream stability release contracts.
# These extend the historical contracts; source-only checks never replace the
# positive repository and gateway regressions executed below.
require_literal oauth_shared_client_cookie_isolation backend/internal/repository/req_client_pool.go 'client := req.C().SetTimeout(opts.Timeout).SetCookieJar(nil)'
require_literal oauth_refresh_no_redirect backend/internal/repository/openai_oauth_service.go 'client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }'
require_literal oauth_refresh_release_compare_owner backend/internal/repository/gemini_token_cache.go 'if redis.call("GET", KEYS[1]) == ARGV[1] then'
require_literal oauth_refresh_unique_acquisition_owner backend/internal/service/oauth_refresh_api.go 'lockOwner := newOAuthRefreshLockOwner()'
require_literal startup_slots_redis_clock backend/internal/repository/concurrency_cache.go 'local now = tonumber(redis.call('"'"'TIME'"'"')[1])'
require_literal startup_slots_only_expired backend/internal/repository/concurrency_cache.go 'local removed = redis.call('"'"'ZREMRANGEBYSCORE'"'"', key, '"'"'-inf'"'"', now - slotTTL)'
require_literal openai_bound_proxy_unavailable_error backend/internal/service/openai_account_proxy_binding.go '"OPENAI_ACCOUNT_PROXY_UNAVAILABLE"'
require_literal openai_gateway_bound_proxy_guard backend/internal/service/openai_gateway_service.go 'if _, err := resolveOpenAIAccountProxyURL(ctx, account, nil); err != nil {'
require_literal openai_retry_after_quota_recovery_barrier backend/internal/service/account_usage_service.go 'if hasActiveOpenAIOAuth429RetryAfter(account, time.Now()) {'
require_literal scheduler_fingerprint_mode_projection backend/internal/repository/scheduler_cache.go '"codex_fingerprint_mode",'
require_literal scheduler_fingerprint_seed_projection backend/internal/repository/scheduler_cache.go '"codex_fingerprint_seed",'
require_literal scheduler_apikey_passthrough_projection backend/internal/repository/scheduler_cache.go '"openai_passthrough",'
require_literal scheduler_oauth_passthrough_projection backend/internal/repository/scheduler_cache.go '"openai_oauth_passthrough",'
require_literal openai_ws_unsupported_idle_recycle backend/internal/service/openai_ws_pool.go 'conn.idleDuration(now) >= openAIWSConnIdleRecycleAfter {'
require_literal openai_continuation_unavailable_recognition backend/internal/service/openai_messages_continuation.go 'lower == "previous_response_id is not available for this user" ||'
require_literal openai_sticky_spillover_keeps_binding backend/internal/service/openai_gateway_scheduling.go 'sessionHash != "" && !stickySpillover && !gatewayProfitControlGateActive(ctx)'
require_test codex_repository_regression_1 backend/internal/repository/req_client_pool_cookie_test.go TestGetSharedReqClient_DoesNotShareResponseCookiesBetweenAccounts
require_test codex_repository_regression_2 backend/internal/repository/openai_oauth_refresh_security_test.go TestOpenAIRefreshRejectsCredentialRedirects
require_test codex_repository_regression_3 backend/internal/repository/gemini_token_cache_lock_test.go TestOAuthRefreshLockExpiredOwnerCannotReleaseNewLease
require_test codex_repository_regression_4 backend/internal/repository/gemini_token_cache_lock_test.go TestOAuthRefreshLockOnlyOwnerCanRelease
require_test codex_repository_regression_5 backend/internal/repository/concurrency_cache_startup_multinode_test.go TestConcurrencyStartupCleanupPreservesOtherNodesActiveSlots
require_test codex_repository_regression_6 backend/internal/repository/concurrency_cache_startup_multinode_test.go TestConcurrencyStartupCleanupReapsOnlyExpiredSlots
require_test codex_repository_regression_7 backend/internal/repository/account_repo_openai_oauth_429_test.go TestSetOpenAIOAuth429RateLimitedAtomicallyBindsDeadlineAndSource
require_test codex_repository_regression_8 backend/internal/repository/scheduler_cache_unit_test.go TestBuildSchedulerMetadataAccount_KeepsOpenAIPassthroughForModelGate
require_test codex_repository_regression_9 backend/internal/repository/account_repo_codex_fingerprint_seed_test.go TestBulkUpdateEnsuresCodexFingerprintSeedWithPerRowSQL
require_test codex_service_regression_1 backend/internal/service/oauth_refresh_api_test.go TestRefreshIfNeeded_RefreshLockOwnerIsUniquePerAcquisition
require_test codex_service_regression_2 backend/internal/service/openai_oauth_429_retry_after_test.go TestOpenAIOAuth429RetryAfterPersistsAndBlocksUntilDeadline
require_test codex_service_regression_3 backend/internal/service/openai_oauth_429_retry_after_test.go TestOpenAIOAuth429RetryAfterSurvivesNonExhaustedQuotaSnapshot
require_test codex_service_regression_4 backend/internal/service/openai_oauth_429_retry_after_test.go TestOpenAIOAuth429RetryAfterSourceCannotProtectUnrelatedOrExpiredLimits
require_test codex_service_regression_5 backend/internal/service/openai_account_proxy_binding_test.go TestOpenAIAccountProxyBinding_RefreshRejectsUnavailableBeforeAnyUpstream
require_test codex_service_regression_6 backend/internal/service/openai_account_proxy_binding_test.go TestOpenAIAccountProxyBinding_QuotaRejectsUnavailableBeforeTokenOrClient
require_test codex_service_regression_7 backend/internal/service/openai_gateway_bound_egress_test.go TestOpenAIBoundEgressRejectsMissingProxyBeforeHTTP
require_test codex_service_regression_8 backend/internal/service/openai_gateway_bound_egress_test.go TestOpenAIBoundEgressRejectsMissingProxyAtWebSocketIngress
require_test codex_service_regression_9 backend/internal/service/openai_codex_account_profile_forward_test.go TestOpenAIGatewayService_Forward_AccountProfileAndFullAcrossHTTPAndWS
require_test codex_service_regression_10 backend/internal/service/openai_codex_missing_turn_metadata_test.go TestCodexMissingTurnMetadataForwardPaths
require_test codex_service_regression_11 backend/internal/service/openai_codex_missing_turn_metadata_test.go TestCodexMissingTurnMetadataWebSocketIngress
require_test codex_service_regression_12 backend/internal/service/openai_oauth_passthrough_test.go TestOpenAIGatewayService_CodexFingerprintCompactDoesNotRewriteBodyCacheKeyOrMetadata
require_test codex_service_regression_13 backend/internal/service/openai_codex_full_mode_test.go TestCodexFullModeCreatedAccountsHaveIndependentStableIdentities
require_test codex_service_regression_14 backend/internal/service/openai_ws_pool_test.go TestOpenAIWSConnPool_RecyclesUnsupportedIdlePingConnection
require_test codex_service_regression_15 backend/internal/service/openai_compat_model_test.go TestOpenAICompatPreviousResponseUnavailableRecognitionIsStrict
require_test codex_service_regression_16 backend/internal/service/openai_compat_model_test.go TestForwardAsAnthropic_PreviousResponseUnavailableRetryFailureDoesNotLoop
require_test codex_service_regression_17 backend/internal/service/openai_gateway_service_test.go TestOpenAISelectAccountWithLoadAwareness_StickyCapacitySpilloverKeepsBinding
require_test codex_service_regression_18 backend/internal/service/openai_account_scheduler_test.go TestOpenAIGatewayService_SelectAccountWithScheduler_DefaultDisabled_LoadBatchReportsFilterReasons


# Latest CPA account-identity parity additions; retain all historical contracts.
require_test codex_identity_20260907_1 backend/internal/service/openai_codex_body_session_test.go TestCodexFingerprintBodySessionPrecedence
require_test codex_identity_20260907_2 backend/internal/service/openai_codex_body_session_test.go TestCodexBodyOnlySessionsAcrossForwardTransports
require_test codex_identity_20260907_3 backend/internal/service/openai_codex_metadata_merge_test.go TestCodexMetadataPreservesBothCarriersAndPrecision
require_test codex_identity_20260907_4 backend/internal/service/openai_codex_metadata_merge_test.go TestCodexMetadataMergeModeAndTurnIsolation
require_test codex_identity_20260907_5 backend/internal/service/openai_codex_metadata_merge_test.go TestCodexMetadataCacheAliasesAndMalformedFallback
require_test codex_identity_20260907_6 backend/internal/service/openai_codex_metadata_merge_test.go TestCodexMetadataFallbackCacheBelongsToOriginalCarrier

if [[ "$MODE" == "tests" ]]; then
  go_bin=$(find_go) || fail "go_toolchain" "required Go binary is $REQUIRED_GO_BIN; rejected ${GO_BIN:-not-set}"
  GO_VERSION=$(GOENV=off GOTOOLCHAIN=local "$go_bin" version 2>/dev/null) || fail "go_toolchain" "cannot execute $go_bin"
  [[ "$GO_VERSION" == "$REQUIRED_GO_VERSION" ]] || fail "go_toolchain" "expected $REQUIRED_GO_VERSION, got $GO_VERSION"
  pass "go_toolchain"

  run_go_test "go_routes" "./internal/server/routes" "$go_bin" \
    "TestRegisterUserSubscriptionRoutesIncludesResetDailyQuota" \
    "TestRegisterSubscriptionRoutesIncludesResetTodayUsage" \
    "TestRegisterSubscriptionRoutesIncludesResetQuota"
  run_go_test "go_handlers" "./internal/handler/..." "$go_bin" \
    "TestResetDailyQuotaRequiresConfirmationAndIdempotencyKey" \
    "TestResetDailyQuotaUsesJWTUserAndReplaysIdempotentResponse" \
    "TestResetTodayUsageRequiresExplicitConfirmation" \
    "TestResetTodayUsageRequiresIdempotencyKeyEvenInObserveOnlyMode" \
    "TestResetSubscriptionQuotaRequestDecodesDailyWeeklyMonthlySelectors" \
    "TestResetQuotaForwardsDailyWeeklyMonthlySelectors" \
    "TestResetQuotaRequiresAtLeastOneSelectedWindow" \
    "TestUserSubscriptionMappersExposeEffectiveWindowResetTimes"
  run_go_test "go_subscription_services" "./internal/service" "$go_bin" \
    "TestSubscriptionDailyCalendarBoundaryUsesShanghaiMidnight" \
    "TestSubscriptionDailyResetTimeUsesNextCalendarMidnightAfterManualReset" \
    "TestCheckAndResetWindowsAdvancesDailyAcrossMultipleDaysOnly" \
    "TestNormalizeExpiredWindowsAdvancesDailyResponseCopy" \
    "TestValidateAndCheckLimitsUsesDailyCalendarBoundary" \
    "TestCalculateProgressUsesDailyCalendarMidnight" \
    "TestGetByIDNormalizesDailyResponseCopyWithoutMutatingRepositorySnapshot" \
    "TestUserSubscriptionNeedsDailyReset_DailyCardKeepsOneTimeQuota" \
    "TestUserSubscriptionNeedsDailyReset_MultiDaySubscriptionStillRefreshes" \
    "TestAdminResetQuota_ResetBoth" \
    "TestAdminResetQuota_ResetDailyOnly" \
    "TestAdminResetQuota_ResetWeeklyOnly" \
    "TestAdminResetQuota_ResetMonthlyOnly" \
    "TestAdminResetQuotaUsesDatabaseCutoffBarrierPublishesAndInvalidatesLocalL1" \
    "TestAdminResetQuotaCacheFailureDoesNotFailCommittedResetAndRetries" \
    "TestSubscriptionQuotaCacheRetryNewGenerationGetsFullCycleAfterLastFailure" \
    "TestResetTodayCacheRetryUsesOnlyFailedSubsetAndRunsPastCacheTTL" \
    "TestResetTodayCacheBatchClearsAllL1BeforeSlowRemoteAndDoesNotStarveFastTarget" \
    "TestResetTodayCacheBatchHasBoundedConcurrencyAndSharedDeadline" \
    "TestAdminResetQuotaRefreshFailureReturnsCommittedSnapshot" \
    "TestResetTodayUsageTransactionAcceptsNativeSelfResetFromInitiallyNullDailyWindow" \
    "TestResetTodayUsageWindowClassificationAtShanghaiMidnight" \
    "TestResetTodayUsageSQLPreservesWindowStartsAndClampsUsage" \
    "TestAdminResetTodayUsage_AtomicSummaryAndCacheFailureStillSucceeds" \
    "TestDailyQuotaResetAvailabilityUsesAsiaShanghaiLegacyAnchor" \
    "TestResetUserDailyQuotaCacheFailureDoesNotFailAndRetries" \
    "TestOpenAIGatewayServiceForwardOAuthLegacyCompactUsesRemoteV2AndReturnsUnaryJSON" \
    "TestNormalizeOpenAILegacyCompactRemoteV2Body" \
    "TestOpenAIGatewayService_OAuthPassthrough_CodexMissingInstructionsGetsDefault"
  run_go_test "go_subscription_repository" "./internal/repository" "$go_bin" \
    "TestResetDailyQuotaForUserUsesAtomicOwnedEffectiveUpdate" \
    "TestResetDailyQuotaForUserNullWindowWritesShanghaiNativeLedgerWindow" \
    "TestResetUsageWindowsReturnsStrictlyMonotonicDatabaseCutoff" \
    "TestBillingCacheSubscriptionBarrierRejectsStaleSnapshotAndAllowsFresh" \
    "TestBillingCacheSubscriptionBarrierRejectsOldIncrementAndAllowsNewIncrement" \
    "TestApplyUsageBillingEffectsUsesStrictlyMonotonicSubscriptionDatabaseWriteTime" \
    "TestSchedulerMetadataAccountKeepsOpenAIImageExactSizeCapability" \
    "TestShouldEnqueueSchedulerOutboxForExtraUpdates_ImageExactSizeCapabilityIsRelevant"
  run_go_test "go_image_download_routes" "./internal/server/routes" "$go_bin" \
    "TestGatewayRoutesImageDownloadPathIsRegisteredWithoutRootAlias" \
    "TestGatewayRoutesImageDownloadRequiresAPIKeyAuthentication" \
    "TestGatewayImageDownloadRouteRunsBeforeCompositeBodyParsing"
  run_go_test "go_image_download_middleware" "./internal/server/middleware" "$go_bin" \
    "TestIsImageDownloadReadRequiresExactMethodAndPath" \
    "TestAPIKeyAuthImageDownloadSkipsBillingChecks" \
    "TestCORSPreflightAllowsImageDownloadAuthorization"
  run_go_test "go_image_download_handlers" "./internal/handler" "$go_bin" \
    "TestImageDownloadHandlerStreamsAttachment" \
    "TestImageDownloadHandlerRejectsOversizedRequestBeforeFetch" \
    "TestImageDownloadConcurrencyGateRateLimitsAndResets" \
    "TestImageDownloadConcurrencyGateParallelAcquireAndIdempotentRelease" \
    "TestImageDownloadHandlerRecordsTruncatedStream"
  run_go_test "go_image_download_services" "./internal/service" "$go_bin" \
    "TestImageDownloadServiceStreamsAllowedImage" \
    "TestImageDownloadServiceRejectsUntrustedURLsBeforeRequest" \
    "TestImageDownloadServiceRejectsRedirectToUntrustedHost" \
    "TestImageDownloadServiceRejectsMissingAndOversizedContentLength" \
    "TestImageDownloadServiceRejectsInvalidOrMismatchedContentType" \
    "TestImageDownloadServiceDefaultClientBlocksPrivateDialAddress" \
    "TestImageDownloadServiceStopsAfterRedirectLimit" \
    "TestImageDownloadPublicIPPolicyRejectsPrivateAndSpecialUse" \
    "TestImageDownloadSafeDialRejectsPrivateDNSAnswerBeforeDial" \
    "TestImageDownloadSafeDialRejectsMixedPublicPrivateAnswers" \
    "TestRewriteOpenAIImagesRequest_GPTImage2PreservesExactSizeAndImage2Fields" \
    "TestRewriteOpenAIImagesRequest_LegacyModelStillNormalizesCustomSize" \
    "TestRewriteOpenAIImagesMultipartRequest_GPTImage2PreservesExactSize" \
    "TestParseOpenAIImagesRequest_JSONEditTypeURLAliasesResponseFormat" \
    "TestOpenAIGatewayServiceForwardImages_APIKeyJSONEditPreservesExactSizeAndNestedImage" \
    "TestOpenAIGatewayServiceForwardImages_OAuthJSONEditExactSizeFailsBeforeUpstream" \
    "TestOpenAIGatewayServiceForwardImages_OAuthJSONEditStandardSizeUsesTypeURLAlias" \
    "TestParseOpenAIImagesRequest_ResponseFormatWinsTypeAlias" \
    "TestOpenAIGatewayServiceForwardImages_ExactSizeRejectsLegacyModelMapping" \
    "TestOpenAIGatewayServiceForwardImages_UnmarkedAPIKeyExactSizeFailsBeforeUpstream" \
    "TestOpenAIGatewayServiceForwardImages_APIKeyExactSizeStreamingPreservesRequest" \
    "TestOpenAIImagesRequestResolveRequiredCapabilityForEffectiveModel" \
    "TestAccountSupportsOpenAIImageCapability_ExactSizeRequiresAPIKey" \
    "TestOpenAIGatewayService_SelectAccountForExactImageSizeNeverFallsBackToOAuth" \
    "TestOpenAIGatewayServiceForwardImages_OAuthUsesDecodedOutputDimensions" \
    "TestOpenAIGatewayServiceForwardImages_OAuthStreamingUsesDecodedOutputDimensions"
  run_go_test "go_image_task_cors" "./internal/server/middleware" "$go_bin" \
    "TestCORS_PreflightAllowsStudioImageTaskHeader"

  run_go_test "go_codex_repository_contracts" "./internal/repository" "$go_bin" \
    "TestGetSharedReqClient_DoesNotShareResponseCookiesBetweenAccounts" \
    "TestOpenAIRefreshRejectsCredentialRedirects" \
    "TestOAuthRefreshLockExpiredOwnerCannotReleaseNewLease" \
    "TestOAuthRefreshLockOnlyOwnerCanRelease" \
    "TestConcurrencyStartupCleanupPreservesOtherNodesActiveSlots" \
    "TestConcurrencyStartupCleanupReapsOnlyExpiredSlots" \
    "TestSetOpenAIOAuth429RateLimitedAtomicallyBindsDeadlineAndSource" \
    "TestBuildSchedulerMetadataAccount_KeepsOpenAIPassthroughForModelGate" \
    "TestBulkUpdateEnsuresCodexFingerprintSeedWithPerRowSQL"
  run_go_test "go_codex_service_contracts" "./internal/service" "$go_bin" \
    "TestRefreshIfNeeded_RefreshLockOwnerIsUniquePerAcquisition" \
    "TestOpenAIOAuth429RetryAfterPersistsAndBlocksUntilDeadline" \
    "TestOpenAIOAuth429RetryAfterSurvivesNonExhaustedQuotaSnapshot" \
    "TestOpenAIOAuth429RetryAfterSourceCannotProtectUnrelatedOrExpiredLimits" \
    "TestOpenAIAccountProxyBinding_RefreshRejectsUnavailableBeforeAnyUpstream" \
    "TestOpenAIAccountProxyBinding_QuotaRejectsUnavailableBeforeTokenOrClient" \
    "TestOpenAIBoundEgressRejectsMissingProxyBeforeHTTP" \
    "TestOpenAIBoundEgressRejectsMissingProxyAtWebSocketIngress" \
    "TestOpenAIGatewayService_Forward_AccountProfileAndFullAcrossHTTPAndWS" \
    "TestCodexMissingTurnMetadataForwardPaths" \
    "TestCodexMissingTurnMetadataWebSocketIngress" \
    "TestOpenAIGatewayService_CodexFingerprintCompactDoesNotRewriteBodyCacheKeyOrMetadata" \
    "TestCodexFullModeCreatedAccountsHaveIndependentStableIdentities" \
    "TestOpenAIWSConnPool_RecyclesUnsupportedIdlePingConnection" \
    "TestOpenAICompatPreviousResponseUnavailableRecognitionIsStrict" \
    "TestForwardAsAnthropic_PreviousResponseUnavailableRetryFailureDoesNotLoop" \
    "TestOpenAISelectAccountWithLoadAwareness_StickyCapacitySpilloverKeepsBinding" \
    "TestOpenAIGatewayService_SelectAccountWithScheduler_DefaultDisabled_LoadBatchReportsFilterReasons"
  run_go_test "go_codex_identity_20260907" "./internal/service" "$go_bin" \
    "TestCodexFingerprintBodySessionPrecedence" \
    "TestCodexBodyOnlySessionsAcrossForwardTransports" \
    "TestCodexMetadataPreservesBothCarriersAndPrecision" \
    "TestCodexMetadataMergeModeAndTurnIsolation" \
    "TestCodexMetadataCacheAliasesAndMalformedFallback" \
    "TestCodexMetadataFallbackCacheBelongsToOriginalCarrier"
fi

emit_provenance "pass"
