# Sub2API Provider Hall Test

Independent Docker test deployment, initialized 2026-09-12.

The operational source of truth is [the deployment README](/opt/sub2api-provider-hall-test/README.md). It records release hashes, versions, ports, services, data volumes, credential references, verification, backups and recovery commands.

- App: `http://127.0.0.1:18082/login`, via SSH tunnel from the operator computer.
- Deployment root: `/opt/sub2api-provider-hall-test`.
- Current login record: `credentials/login.json` under that root, root-only.
- App, PostgreSQL and Redis use an independent Docker network and persistent storage.
- Production configuration files, PostgreSQL/Redis tunnels and credentials are not used by this stack.
- History: [DEPLOYMENT-RECORD.md](DEPLOYMENT-RECORD.md).
