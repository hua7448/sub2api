# 邀请返利历史补发归档

本目录保存 2026-07-22 邀请返利历史数据核对与一次性补发脚本。它们是生产快照的操作留档，不是应用启动时自动执行的数据库迁移，也不能代替 `backend/migrations/185_affiliate_redeem_code_sources.sql`。

## 文件

- `affiliate_rebate_backfill_audit.sql`：只读计算历史兑换事件，汇总已返利、缺失返利并列出缺失明细。
- `affiliate_rebate_backfill_apply.sql`：只补订阅天数返利的一次性事务脚本，包含前置断言、插入断言和事务提交。

## 重要限制

- `apply` 脚本绑定当时的生产快照，硬性要求候选记录为 65 条、合计 195 天。数据不一致时会抛错并终止。
- 成功补发后，原前置条件自然不再成立，因此不能直接重复执行。
- 后续环境或新数据必须先运行只读 `audit`，重新审核结果并生成新的、带独立断言的操作脚本。
- 两份脚本不含固定邮箱、密码、令牌、API Key 或服务器地址；执行 `audit` 时查询结果会显示相关用户邮箱，结果不得提交到仓库。
- 通用防重复能力由 migration 185 的 `source_redeem_code_id` 和唯一索引提供。

## 生产执行记录

- 执行时间：2026-07-22 10:25:13 UTC。
- 写入结果：65 条 `accrue_days` 账本记录，合计 195 个订阅返利天数，涉及 24 位邀请人。
- migration 185 在 `schema_migrations` 中仅有 1 条执行记录。
- 10:49 UTC 后续在线流程又正常产生 3 条、每条 3 天的返利记录，与历史补发批次可区分。

## 校验值

```text
a79a0f5ea74c6d949479d5f1f36ff4aa9217baa65b24793fca2d88331a7b8145  affiliate_rebate_backfill_audit.sql
6dc06c95a70f9cbdf9d6a5096aeaca87e72c58f7395dd760cd5b1462fd1a5a31  affiliate_rebate_backfill_apply.sql
```

归档状态：已执行，仅供审计和追溯，禁止加入自动部署或 migration 扫描目录。
