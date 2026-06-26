# Agent Rules

This repository is the `hua7448/sub2api` fork of `Wei-Shaw/sub2api`.

Hard rules for all automation and agent work:

- Follow `docs/FORK_WORKFLOW_CN.md` for branch naming, upstream sync, release tags, server trial deployment, production rollout, rollback, and data safety.
- Keep `origin` pointed at `git@github.com:hua7448/sub2api.git`.
- Keep `upstream` pointed at `https://github.com/Wei-Shaw/sub2api.git`.
- Treat upstream sync as needed only when the official project publishes a new baseline version tag; same-baseline upstream commits are not merged by default unless explicitly approved as an urgent exception.
- Do not commit `.env`, database dumps, `data/`, `postgres_data/`, `redis_data/`, user records, or private reports.
- Do not suggest or run production-destructive commands such as `docker compose down -v`, `rm -rf data postgres_data redis_data`, or production `docker volume rm`.
- Production deployments must use explicit version tags, not `latest`.
- Before replacing the production service, use a separate trial container and non-production port, usually `4146`.
- If a change includes database migrations, call out rollback risk before production rollout.
