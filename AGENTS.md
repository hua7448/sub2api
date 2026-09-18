# Sub2API Working Guide

## Scope and Evidence

- This is a heavily customized Sub2API source snapshot, not a pristine upstream checkout. Preserve the existing product behavior when fixing or importing code.
- The initial inspection on 2026-09-10 found no `.git` directory. Do not infer a current branch, remote, clean worktree, or exact upstream diff from historical documentation. Do not initialize Git merely to bypass provenance checks.
- `backend/cmd/server/VERSION` contains `0.1.179`; this is an embedded version fallback, not proof of the deployed version or upstream ancestry.
- Read `docs/PROJECT_MAP.md` for the architecture and extension map. `CUSTOM_CHANGELOG.md` records local production changes through July 2026; `docs/reviews/CPA_ACCOUNT_IDENTITY_20260907.md` describes later source changes.
- `dev-memory/`, `DEV_GUIDE.md`, and `LOCAL_DEV_GUIDE.md` contain historical Windows environments and older branches. Verify their commands and assumptions against current code. Their recorded owner, credentials, and deployment state are not current-session facts.
- Communicate with the user in Chinese unless requested otherwise.

## Code Layout

- `backend/cmd/server/`: startup, build metadata, Wire composition and shutdown.
- `backend/internal/server/router.go` and `routes/`: panel API under `/api/v1`, provider gateways under `/v1`, `/v1beta`, and compatibility aliases.
- `backend/internal/server/middleware/`: authentication, API key billing gates, audit, request metadata, rate limits. A separate `backend/internal/middleware/` also exists; follow actual imports.
- `backend/internal/handler/`: HTTP handlers and `dto/`; administrator handlers live in `handler/admin/`.
- `backend/internal/service/`: domain services, upstream protocols, account selection, billing, subscriptions and workers.
- `backend/internal/repository/`: Ent/SQL data access, Redis adapters, migration runner and cache invalidation.
- `backend/ent/schema/`: ORM definitions. Other Ent directories are predominantly generated code.
- `backend/internal/securityaudit/`: prompt scanning, synchronous guard, async jobs, configuration and admin API; keep its existing coordinator integration.
- `backend/internal/payment/`: payment registry and provider implementations.
- `frontend/src/`: Vue 3 + TypeScript, Pinia, Vue Router, Tailwind and vue-i18n. Pages are in `views/`, shared widgets in `components/`, HTTP clients in `api/`.
- `frontend/src/features/`: prompt audit and channel monitor V2 feature modules.
- `frontend/public/image-studio/`: checked-in embedded subapplication assets. Establish their source/build provenance before replacing bundles.
- `deploy/`: container, native service and installation tooling. `docs/operations/`: historical one-time operations, not automatic migrations.

## Development and Verification

Use pnpm for the frontend and retain `frontend/pnpm-lock.yaml`. CI uses Node 20 and pnpm 9. Do not regenerate the lockfile with another major version as an incidental change.

Toolchain conflict to resolve before claiming CI compatibility: `backend/go.mod` declares Go 1.26.5, CI reads that file but asserts 1.26.6, Dockerfiles use 1.26.6, and the production contract script pins Linux Go 1.26.5 at an external absolute path. Do not silently normalize production pins during unrelated work.

Commands from the repository root, once appropriate dependencies are available:

```sh
pnpm --dir frontend install --frozen-lockfile
pnpm --dir frontend run dev
pnpm --dir frontend run lint:check
pnpm --dir frontend run typecheck
pnpm --dir frontend run test:run
make test-frontend
make -C backend test-unit
make -C backend test-integration
make -C backend build
```

- `make test-frontend` runs lint, typecheck and a selected critical Vitest suite, not the entire frontend suite.
- `make -C backend test` runs untagged tests and golangci-lint. It does not replace the separate `unit` and `integration` build-tag suites.
- For focused Go checks, run `go test -tags=unit ./internal/<package> -run '<pattern>'` from `backend/`. Inspect test setup for Docker, PostgreSQL, Redis or environment-variable requirements; report skipped tests.
- Frontend `lint` auto-fixes files; use `lint:check` for read-only verification.
- Backend development entry point: `go run ./cmd/server` from `backend/`, with an isolated local `DATA_DIR` and PostgreSQL/Redis configuration. Setup uses `config.yaml` and `.installed`; complete setup rather than blindly copying historical credentials or bypassing it with a lock file.
- Vite defaults to port 3000 and proxies to port 8080; override using `VITE_DEV_PORT` and `VITE_DEV_PROXY_TARGET`.
- `pnpm --dir frontend run build` empties and writes `backend/internal/web/dist`. Embedded Go builds require that directory and `-tags embed`. The ordinary backend Makefile build does not embed the UI.
- Production contract verification is in `backend/scripts/verify-production-contracts.sh`. Even `--source-only` requires Git ancestry/canonical refs; this export cannot satisfy that provenance check. Do not weaken the script to make the export pass.

## Change Conventions

- Follow route -> handler/DTO -> service -> repository boundaries and existing provider sets. Update frontend API types, callers, settings projections and translations together when changing a contract.
- Reuse `frontend/src/api/client.ts`, shared components and current styles. Preserve dark mode and mobile behavior. Route visibility is not backend authorization.
- Keep Chinese and English locale entries synchronized under `frontend/src/i18n/locales/`.
- Update affected test stubs when changing a Go interface. Exercise relevant HTTP, streaming, WebSocket and failover paths when changing gateway behavior.
- Edit ORM schemas in `backend/ent/schema/`, then run `go generate ./ent` from `backend/`. Update Wire source/provider sets and run `go generate ./cmd/server` when wiring changes. Do not patch generated outputs as the sole fix.
- Add forward-only SQL migrations. Never rename, renumber, delete or rewrite an already-applied migration. The runner uses complete filenames and checksums; repeated numeric prefixes already exist and are not sufficient evidence of a collision.
- `_notx.sql` is reserved for allowed concurrent index operations. Do not add executable Down SQL to an Up file; the runner executes the whole file. The migration README's `make migrate-up/down` examples have no matching targets here.
- Preserve precise historical checksum compatibility rules; do not add broad checksum bypasses. Keep one-time data repair outside migrations with explicit snapshot assumptions and verification.

## Product Contracts to Preserve

- Instance roles: `primary` owns migrations, seeds and designated global workers; `api` skips those. Process-local services and separately coordinated workers may still run. Check each provider's lifecycle instead of disabling every background service.
- `RUN_MODE=simple` disables billing/quota checks; do not use it to validate standard billing behavior.
- Affiliate rewards include balance and subscription days, redeem-code/order provenance, qualification thresholds, expiry and idempotency. Threshold counting excludes the currently rewarded customer.
- Provider-specific email aliases and `+tag` aliases are independent registration policies. Changes must cover registration, verification, binding and OAuth completion while preserving existing login behavior.
- Composite groups resolve to concrete platforms. Preserve model mapping, endpoint matching and actual-provider billing/quota attribution, including domestic providers already wired into the gateway.
- Profit control participates in account selection and admission/rechecks; preserve downstream/upstream price semantics and failover behavior.
- Subscription daily resets include timezone boundaries, transactional operation records, idempotency and cache barriers. Read the production contract script before changing these paths.
- Codex identity handling shares account-scoped IDs across HTTP/WS paths. Preserve persistent seeds, off/device/session/full modes, session-source precedence, per-turn metadata isolation and opaque metadata values. Limit transformations to structured identity fields.
- Async OpenAI/Grok image tasks and Gemini/Vertex batch image jobs are different systems. Preserve task ownership, feature gates, storage, billing settlement and recovery semantics.
- Channel monitor V1 active probes and V2 passive aggregation have separate behavior; the current fallback mode is V1. Respect V2 backfill limits and privacy settings.
- Prompt audit has `off`, `async_audit` and `blocking` modes. Preserve per-request/WS-turn enforcement, configuration snapshots, queue behavior and the configured failure policy.
- Preserve embedded UI revision/cache handling. Historical production UI locking forbids replacing a production UI as an incidental backend update; verify the actual target's release procedure when deployment is requested.
- Record subsequent custom implementation changes in `CUSTOM_CHANGELOG.md` with behavior, source scope, migrations/configuration, validation and actual deployment status. Never present historical tests or deployments as fresh evidence.
