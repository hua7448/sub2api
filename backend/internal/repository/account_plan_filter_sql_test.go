package repository

import (
	"strings"
	"testing"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
)

// TestAccountPlanFilterPredicateSQL 验证 plan_type 过滤谓词在 Postgres 方言下
// 生成的占位符是 $N 而不是裸 ?，避免与查询里其它条件的 $N 混用导致
// "syntax error at or near ,"（曾导致按 pro 过滤账号始终 500）。
func TestAccountPlanFilterPredicateSQL(t *testing.T) {
	pred := accountPlanFilterPredicate("pro")

	s := entsql.Dialect(dialect.Postgres).
		Select("id").
		From(entsql.Table("accounts"))
	pred(s)

	query, args := s.Query()

	if strings.Contains(query, "?") {
		t.Fatalf("生成的 SQL 仍含有裸 ? 占位符（Postgres 不接受）: %s", query)
	}
	if !strings.Contains(query, "$1") {
		t.Fatalf("期望出现 $N 占位符，实际: %s", query)
	}
	// pro: plan_type(2) + tier_id(3) + paidTier(1) + currentTier(1) = 7 个参数
	if len(args) != 7 {
		t.Fatalf("期望 7 个参数，实际 %d 个: %v", len(args), args)
	}
	t.Logf("query=%s args=%v", query, args)
}
