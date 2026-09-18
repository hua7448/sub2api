# Project Memory

> Updated: 2026-09-18

## ⚠️ 当前项目:供应商大厅版本(provider-hall)

- **分支**: `release/v0.1.179-provider-hall`(孤儿分支,基线上游 v0.1.179 + Provider Hall B2–B8),仓库 `/Users/simple/smartq/sub2api`(origin: hua7448/sub2api)
- **权威背景文档**: [provider-hall.md](./provider-hall.md) — 分支来源、功能状态、产品契约、本机测试环境、服务器部署(179.255.156.111)
- **开发记录**: `docs/reviews/PROVIDER_HALL_*.md` + 总计划 `docs/reviews/PROVIDER_HALL_PLAN_20260911.md`
- **与 main 不合并**(无共同历史,基线已岔开)
- **下方 bmai 内容是早期项目(2026-03-30)遗留,已过时**,仅供历史参考;其中的 owner、凭据、部署状态与本分支无关。

---

# sub2api-bmai Project Memory(历史遗留,已过时)

> Updated: 2026-03-30

## Quick Reference
- **Repo**: bayma888/sub2api-bmai (fork of Wei-Shaw/sub2api)
- **Stack**: Go 1.26.1 (Ent ORM + Gin) + Vue3 (pnpm) + PostgreSQL 16 + Redis
- **Owner**: Master (bayma888 / Xiaobei / bma-Beima-北妈)
- **Language**: Master is learning Chinese, speak Chinese to him. Convert Chinese instructions to English internally.
- **Details**: See [project-info.md](./project-info.md), [features.md](./features.md), [pitfalls.md](./pitfalls.md), [deployment.md](./deployment.md)

## Completed Features (merged to upstream main)
1. **Group drag-sort** (PR #519)
2. **User search by notes**
3. **Show admin adjustment notes**
4. **Admin user balance history**
5. **API key quota & expiration**

## Pending PRs (pushed, awaiting merge)
6. **Group selector UX** (`feature/group-display-fix`)
7. **Usage page user balance popup** (`feature/usage-user-balance-popup`)

## Private Production Branch
- **Branch**: `bmai-sub2api-my-production20260329`
- **Purpose**: 私有功能 + 自编译二进制，不提 PR 到上游
- **Current version**: `bmai-v0.1.105.1` (based on upstream v0.1.105)
- **VPS**: `<PRODUCTION_HOST>:8080`, service name `bmai-sub2api`
- **See**: [deployment.md](./deployment.md) for full CI/CD & deployment details

## Planned Features (private, not for upstream PR)
8. **排行榜游戏化系统 (Gamification Leaderboard)** — see [features.md](./features.md)
   - **Status**: ✅ 全栈代码开发完成，待视觉验收
   - **Branch**: `feature/usage-user-balance-popup` (当前分支，与 #7 共用)
   - **需求文档**: `LEADERBOARD_GAMIFICATION_PLAN.md`、`LEADERBOARD_FRONTEND_PLAN.md`
   - **后端**: types + repository SQL + service + handler + route 全部完成
   - **前端**: LeaderboardView.vue + API + i18n(zh/en) + router + sidebar 全部完成
   - **检查**: go build ✅ | go vet ✅ | gofmt ✅ | unit test ✅ | vite build ✅
   - **待做**: 本地运行视觉验收、排名变动箭头(P3)、Redis缓存(P4)
   - **部署tag**: `bmai-v0.1.105.2`

## Upstream Status
- **Latest upstream version**: v0.1.105
- **Local main synced**: ✅ up to date with upstream v0.1.105

## Key Patterns
- CSS: **camelCase** class names (not kebab-case)
- Frontend pkg manager: **pnpm** only (never npm)
- Commit messages: **Chinese** descriptions (Master prefers)
- Go interface changes → must update ALL test stubs
- Ent schema changes → must `go generate ./ent`
- Always use `127.0.0.1` not `localhost` for psql
- gh CLI not authenticated locally — use browser for PR creation
