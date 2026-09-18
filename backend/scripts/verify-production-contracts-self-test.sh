#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)
REPO_ROOT=$(git -C "$SCRIPT_DIR" rev-parse --show-toplevel)
SOURCE_VERIFIER="$SCRIPT_DIR/verify-production-contracts.sh"
SOURCE_SELF_TEST="$SCRIPT_DIR/verify-production-contracts-self-test.sh"
TMP_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/sub2api-contract-self-test.XXXXXX")
REPO_CLONE="$TMP_ROOT/repo"
WORKTREE="$TMP_ROOT/worktree"
STALE_WORKTREE="$TMP_ROOT/stale-worktree"
LEGACY_COMPACT_DUP_WORKTREE="$TMP_ROOT/legacy-compact-dup-worktree"
MIGRATION_186_WORKTREE="$TMP_ROOT/migration-186-worktree"
MIGRATION_194_WORKTREE="$TMP_ROOT/migration-194-worktree"
MIGRATION_186="backend/migrations/186_subscription_quota_reset_operations.sql"
MIGRATION_194="backend/migrations/194_user_subscription_daily_quota_weekly_limit.sql"

cleanup() {
  rm -rf "$TMP_ROOT"
}
trap cleanup EXIT INT TERM

fail() {
  printf 'SELF-TEST FAIL: %s\n' "$1" >&2
  exit 1
}

expect_pass() {
  local label=$1
  shift
  local output
  output=$("$@") || fail "$label unexpectedly failed"
  [[ "$output" == *'"status":"pass"'* ]] || fail "$label did not emit pass provenance"
  printf 'SELF-TEST PASS [%s]\n' "$label" >&2
}

expect_failure() {
  local label=$1
  local expected_check=$2
  shift 2
  local output
  if output=$("$@" 2>/dev/null); then
    fail "$label unexpectedly passed"
  fi
  [[ "$output" == *"\"failed_check\":\"$expected_check\""* ]] || fail "$label failed for the wrong reason: $output"
  printf 'SELF-TEST PASS [%s]\n' "$label" >&2
}

# Commit the scripts under test only in an isolated clone. Its canonical branch
# can advance for clean-tree negative cases without moving any source-repo ref.
git clone --quiet --no-hardlinks --branch production/backend-current "$REPO_ROOT" "$REPO_CLONE"
cp "$SOURCE_VERIFIER" "$REPO_CLONE/backend/scripts/verify-production-contracts.sh"
cp "$SOURCE_SELF_TEST" "$REPO_CLONE/backend/scripts/verify-production-contracts-self-test.sh"
chmod +x "$REPO_CLONE/backend/scripts/verify-production-contracts.sh"
git -C "$REPO_CLONE" add \
  backend/scripts/verify-production-contracts.sh \
  backend/scripts/verify-production-contracts-self-test.sh
git -C "$REPO_CLONE" \
  -c user.name='Production Contract Self-Test' \
  -c user.email='production-contract-self-test@invalid' \
  commit --allow-empty -m 'self-test: verifier candidate' >/dev/null
CANDIDATE_HEAD=$(git -C "$REPO_CLONE" rev-parse HEAD)

git -C "$REPO_CLONE" worktree add --detach "$WORKTREE" "$CANDIDATE_HEAD" >/dev/null
VERIFIER="$WORKTREE/backend/scripts/verify-production-contracts.sh"

expect_pass "clean_source" "$VERIFIER" --source-only
expect_pass "default_guard" "$VERIFIER"

printf 'self-test UI dirtiness\n' >"$WORKTREE/frontend/.production-contract-self-test"
expect_pass "ui_dirtiness_allowed" "$VERIFIER" --source-only
rm "$WORKTREE/frontend/.production-contract-self-test"

printf 'self-test non-UI dirtiness\n' >"$WORKTREE/backend/.production-contract-self-test"
expect_failure "non_ui_dirtiness_rejected" "dirty_non_ui_source" "$VERIFIER" --source-only
rm "$WORKTREE/backend/.production-contract-self-test"

# A tracked production route deletion must be rejected before literal checks.
rm "$WORKTREE/backend/internal/server/routes/user.go"
expect_failure "tracked_route_removal_rejected" "dirty_non_ui_source" "$VERIFIER" --source-only
git -C "$WORKTREE" restore -- backend/internal/server/routes/user.go

expect_failure "fake_go_bin_rejected" "go_toolchain" env GO_BIN=/bin/true "$VERIFIER"

# The candidate's parent still has all fixed feature ancestors, so this must
# fail specifically at the canonical-tip gate rather than at ancestry checks.
git -C "$REPO_CLONE" worktree add --detach "$STALE_WORKTREE" "${CANDIDATE_HEAD}^" >/dev/null
cp "$SOURCE_VERIFIER" "$STALE_WORKTREE/backend/scripts/verify-production-contracts.sh"
chmod +x "$STALE_WORKTREE/backend/scripts/verify-production-contracts.sh"
expect_failure "stale_head_rejected" "canonical_tip" \
  "$STALE_WORKTREE/backend/scripts/verify-production-contracts.sh" --source-only

# The official v0.1.179 refactor moved this declaration. Preserve the local
# compatibility path while rejecting accidental duplicate declarations.
printf '\nconst openAIRemoteCompactionV2Feature = "remote_compaction_v2"\n' >> \
  "$REPO_CLONE/backend/internal/service/openai_gateway_request_body.go"
git -C "$REPO_CLONE" add backend/internal/service/openai_gateway_request_body.go
git -C "$REPO_CLONE" \
  -c user.name='Production Contract Self-Test' \
  -c user.email='production-contract-self-test@invalid' \
  commit -m 'self-test: duplicate legacy compact feature' >/dev/null
LEGACY_COMPACT_DUP_HEAD=$(git -C "$REPO_CLONE" rev-parse HEAD)
git -C "$REPO_CLONE" worktree add --detach "$LEGACY_COMPACT_DUP_WORKTREE" "$LEGACY_COMPACT_DUP_HEAD" >/dev/null
expect_failure "legacy_compact_duplicate_rejected" "legacy_compact_feature_unique" \
  "$LEGACY_COMPACT_DUP_WORKTREE/backend/scripts/verify-production-contracts.sh" --source-only

# Advance only the isolated clone's canonical branch for clean-tree hash cases.
git -C "$REPO_CLONE" restore --source="$CANDIDATE_HEAD" -- backend/internal/service/openai_gateway_request_body.go
printf '\n-- self-test migration 186 hash mutation\n' >>"$REPO_CLONE/$MIGRATION_186"
git -C "$REPO_CLONE" add "$MIGRATION_186" backend/internal/service/openai_gateway_request_body.go
git -C "$REPO_CLONE" \
  -c user.name='Production Contract Self-Test' \
  -c user.email='production-contract-self-test@invalid' \
  commit -m 'self-test: mutate pinned migration 186' >/dev/null
MIGRATION_186_HEAD=$(git -C "$REPO_CLONE" rev-parse HEAD)
git -C "$REPO_CLONE" worktree add --detach "$MIGRATION_186_WORKTREE" "$MIGRATION_186_HEAD" >/dev/null
expect_failure "migration_186_hash_rejected" "migration_186_sha256" \
  "$MIGRATION_186_WORKTREE/backend/scripts/verify-production-contracts.sh" --source-only

git -C "$REPO_CLONE" restore --source="$CANDIDATE_HEAD" -- "$MIGRATION_186"
printf '\n-- self-test migration 194 hash mutation\n' >>"$REPO_CLONE/$MIGRATION_194"
git -C "$REPO_CLONE" add "$MIGRATION_186" "$MIGRATION_194"
git -C "$REPO_CLONE" \
  -c user.name='Production Contract Self-Test' \
  -c user.email='production-contract-self-test@invalid' \
  commit -m 'self-test: mutate pinned migration 194' >/dev/null
MIGRATION_194_HEAD=$(git -C "$REPO_CLONE" rev-parse HEAD)
git -C "$REPO_CLONE" worktree add --detach "$MIGRATION_194_WORKTREE" "$MIGRATION_194_HEAD" >/dev/null
expect_failure "migration_194_hash_rejected" "migration_194_sha256" \
  "$MIGRATION_194_WORKTREE/backend/scripts/verify-production-contracts.sh" --source-only

printf '{"schema":"sub2api.production-contracts-self-test.v1","status":"pass","cases":10}\n'
