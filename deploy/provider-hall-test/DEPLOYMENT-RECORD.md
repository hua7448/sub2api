# Deployment Record

## 2026-09-12: Initial Attempt (Retrospective)

- Operator requested independent Docker testing of the supplied `sub2api-provider-hall-0.1.179-20260912.tar.gz` binary release.
- Docker Engine 29.1.3 installed; amd64 binary packaged in image `sub2api-provider-hall:0.1.179-20260912`.
- Initial app used production config and PostgreSQL/Redis tunnels in `api` role. A production backup attempt failed authentication and produced an empty file; it was not a successful backup. Health-check success did not verify initialization or provider-hall readiness.
- Operator then authorized a new isolated database and randomized credentials. PostgreSQL was created with database/user `sub2api_test` and volume `sub2api-provider-hall-test-pgdata`.
- PostgreSQL 18 initially failed due to its volume being mounted at the old `/var/lib/postgresql/data` path. Corrected to `/var/lib/postgresql`.
- The app remained uninitialized: copied config/installation markers bypassed setup, `DATABASE_NAME` was incorrect (required `DATABASE_DBNAME`), and modifying the env file followed only by restart did not update the container environment. Redis still referenced production. Pre-generated administrator fields were not valid verified login credentials.

## 2026-09-12 07:58-08:02 PDT / 14:58-15:02 UTC: Completed Isolated Initialization

- Read inventory, current deployment files and setup implementation. Stopped the failing test container.
- Installed Compose 2.40.3 and replaced the stale Compose manifest with app, PostgreSQL and dedicated Redis services on `provider-hall-test_default`.
- Preserved the isolated PostgreSQL volume and took `backups/test-before-initialize-20260912.dump`; confirmed the database initially had no public tables.
- Archived earlier app config, invalid install marker, unused administrator files and failed empty production dump under root-only `backups/failed-bootstrap-20260912/`.
- Generated a random administrator email identifier and strong password, plus independent stable JWT/TOTP secrets. Current credentials are `credentials/login.json` and `credentials/runtime.env`; no secret values are recorded here.
- Removed production config mounts and tunnel endpoints. Recreated containers so corrected environment variables took effect. App role is `primary`, node ID `provider-hall-test`.
- Real `AUTO_SETUP` completed database/Redis checks, migrations, administrator creation, config writing and installation lock creation at approximately 07:59 PDT.
- Confirmed `231_provider_hall_facts.sql`, `232_provider_hall_metrics.sql`, `233_provider_hall_jobs.sql` applied at 14:59:06 UTC; verified 17 provider-hall tables and one administrator.
- Verified login credentials through `POST /api/v1/auth/login` and admin role through `GET /api/v1/auth/me`; fetched login HTML and all six referenced JS/CSS assets successfully. Tokens/passwords were not printed.
- `GET /setup/status` returned completed with `needs_setup=false`; all three containers became healthy.
- Created `backups/20260912T150044Z/`: database dump, config/credential archive and SHA-256 checksums. Restored dump to temporary isolated database `provider_hall_restore_check_20260912`; confirmed three provider-hall migrations and one administrator, then dropped the temporary database.
- Restarted app, PostgreSQL and Redis; repeated health and login verification successfully at 15:01:19 UTC. No ERROR/FATAL/panic or auto-generated TOTP-key warnings in the checked startup logs.
- Production Sub2API, state tunnel, Nginx and Dujiao-Next units remained active; production Sub2API health returned OK. This completion operation did not connect to the production databases or edit production services/vhosts.
- Only the test app publishes a host port: `127.0.0.1:18082` -> container `8080`. PostgreSQL and Redis expose no host ports; old test PostgreSQL port `15432` is no longer published.
- Provider-hall business switches remain disabled: supplier channels/keys, operator funding, HTTPS origin and budget require actual test business configuration. No loopback override was enabled and no chargeable probes were launched.
- Added backup and credential-suppressing login verification scripts, protected build context with `.dockerignore`, and updated operational documentation/inventory.

Image identities at completion:

- App: `sha256:9e9415d6e84be06b407895654f0aa04db15c059c40b374522d14c87316c998e0`.
- PostgreSQL 18.6: `sha256:d3e1620b530c944afa6e887d22eb899824da68e19c52024bf98f5220c88a65b2`.
- Redis 7.4.11: `sha256:ff02b58f971e7d7d156a1267e283fcbbeee91773b6aa36c49dac28ecfe28eadf`.

## 2026-09-13: Public Test Port

- Changed only the test app mapping to `0.0.0.0:18082:8080` at operator request. Verified `/health` through loopback and `179.255.156.111`. PostgreSQL and Redis remain unpublished; production services and Nginx vhosts were unchanged. Public URL: `http://179.255.156.111:18082/login` (HTTP only, no domain/TLS).

## 2026-09-15 23:21 PDT: Loopback Gateway Origin Override

- Enabled `PROVIDER_HALL_ALLOW_LOOPBACK=1` only on the isolated provider-hall test app container in `/opt/sub2api-provider-hall-test/compose.yaml`.
- Recreated only `sub2api-provider-hall-test` with Docker Compose. PostgreSQL and Redis containers, volumes and credentials were not changed.
- Verified the recreated app container has the loopback override in its environment and reports Docker health `healthy`; `GET /health` on `127.0.0.1:18082` succeeds.
- Compose backup before this edit: `/opt/sub2api-provider-hall-test/compose.yaml.before-loopback-20260915T232129-0700`.

## 2026-09-15 23:36 PDT: Provider Hall Probe Failure Diagnosis

- Investigated provider-hall probe failures after operator configured real upstreams.
- Database task records showed failed probe and verification jobs with `transport_refused` / `sample_error`, and sample records had HTTP status `0` with `transport_refused` details.
- Confirmed from inside the app container that `http://127.0.0.1:18082/health` is refused, while `http://127.0.0.1:8080/health` succeeds. The earlier gateway origin used the host-published port, which is not listening inside the container.
- Updated provider-hall config `gateway_origin` from `http://127.0.0.1:18082` to `http://127.0.0.1:8080`; config version advanced to `4`. No credentials, supplier targets, PostgreSQL/Redis containers or production services were changed.

## 2026-09-15 23:49 PDT: Provider Hall Probe Funding Fix

- Investigated operator-reported failed jobs `#101` and `#112`. After the gateway origin correction, both jobs were reaching `/v1/responses` but returned HTTP `403`.
- Replayed a minimal request through the test gateway with the registered probe API key. The gateway response was `INSUFFICIENT_BALANCE`, confirming the failures were caused by the provider-hall operator account having zero balance and the probe key having zero quota.
- Added test-only funding to operator user `1` and probe API key `1` in the isolated test database: user balance `100`, key quota `100`. No secret values were printed or recorded.
- Verified a direct `/v1/responses` call returned HTTP `200` for `gpt-5.5`.
- Triggered provider-hall verification job `#175` and probe job `#176` with profile `1`; both completed successfully. Verification `#175` passed all six samples, and probe `#176` returned HTTP `200` with result `passed`.
- Production services, Nginx, PostgreSQL/Redis production tunnels, container image, Compose manifest and supplier credentials were not changed.

## 2026-09-16 00:33 PDT: Provider Hall Display Sample Generation

- Generated additional provider-hall test data for the operator to inspect the user-facing supplier hall.
- Triggered manual probe jobs `#235` through `#240`; all six succeeded. Triggered verification job `#241`; it completed with verdict `passed` across arithmetic and JSON suites. A later verification batch hit upstream instability/rate limits and was not used as the final healthy sample.
- Added test quota to ordinary API key `2` and sent ordinary group `2` user traffic so hall real-user metrics could populate. The public hall snapshot then showed group `2` / `OAI 001` with `success_rate=1.0` and `cache_rate=0.8741179149`.
- Cleaned up queued test jobs created during manual sampling, including prepared samples that never dispatched, by marking non-dispatched pending samples unbilled. This removed the temporary `billing_backlog` condition caused by interrupted/cancelled test jobs.
- Ran final probe job `#692`; it succeeded with HTTP `200`, model `gpt-5.5`, result `passed`, and refreshed group `2` health to `up` at `2026-09-16T07:32:07Z`.
- Final hall API check showed summary `available=1`, `listed=2`, `abnormal=0`, `verified=1`; group `2` health `up`, verification report `#286` passed, success rate OK and cache rate OK. TTFT cards still report insufficient TTFT sample count because many non-stream ordinary requests do not record TTFT.
- Left `collection_enabled=true`, `display_enabled=true`, and `tasks_enabled=false` to stop continuous scheduled probing and avoid further test spend. Manual probes/verifications require temporarily enabling tasks again.
- No production services, Nginx vhosts, Compose image, supplier credentials, PostgreSQL/Redis containers or persistent volumes were changed.

## 2026-09-16 01:17 PDT: Provider Hall TTFT Sample Completion

- Operator asked whether missing average/TTFT data comes from real user traffic. Confirmed the hall's TTFT cards require ordinary user gateway requests with recorded `ttft_ms`, while probe/verification jobs mainly drive health and verification state.
- Used ordinary API key `2` on group `2` to send streaming `/v1/responses` requests against the isolated test gateway. No API key value was printed or recorded.
- Confirmed the test API base is `http://179.255.156.111:18082`; the local execution used the loopback equivalent `http://127.0.0.1:18082/v1/responses`.
- After aggregation, group `2` / `OAI 001` had `21` ordinary user TTFT samples in the one-hour window. Hall metrics showed `ttft_fast95_ms=36170.8`, `ttft_p90_ms=64413`, `success_rate=0.86`, and `cache_rate=0.8734432123`.
- The same final check showed group `2` health `up` for model `gpt-5.5`. A later automatic verification run had verdict `insufficient` / `sample_error`, so `verified` was temporarily `0`; this reflected latest verification sampling, not lack of real-user TTFT metrics.
- Final observed balances: operator user balance approximately `99.364932`, ordinary API key `2` quota `100` with used quota approximately `0.215958`.
- Documentation updated. Production services, Nginx vhosts, Compose image, supplier credentials, PostgreSQL/Redis containers and persistent volumes were not changed.
