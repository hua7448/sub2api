# Server Inventory

Last scanned: 2026-09-09 (UTC)

## Identity

- Hostname: `node4`
- Location: California, USA (provided by the operator; not independently geolocated)
- Operating system: Ubuntu 26.04 LTS (`resolute`)
- Kernel: `7.0.0-27-generic` on `x86_64`
- Uptime at scan: 52 days, 9 hours

## Capacity

- CPU: 2 vCPU, AMD EPYC-Milan
- Memory: 1.9 GiB total, 1.4 GiB available at scan
- Swap: 4.0 GiB, 3.8 GiB available at scan
- Root disk: 40 GiB ext4, 16 GiB used (41%), 23 GiB available

## Installed tooling detected

- `bubblewrap` 0.11.1-1ubuntu0.1 (`/usr/bin/bwrap`)
- Nginx 1.28.3-2ubuntu1.6
- Node.js 22.23.2 (NodeSource)
- npm present on PATH
- Python 3.14.3
- Certbot 4.0.0-4
- Git, curl, wget

## Access and security observations

- SSH effective configuration: port 22 on IPv4 and IPv6; public-key authentication enabled; password authentication disabled; root login is `prohibit-password`.
- `/root/.ssh/authorized_keys` exists with mode `600`.
- UFW/firewalld status was not available from the current execution context.
- Automatic-updates status was not available from the current execution context.
- No secrets were copied into this inventory.

## Scan limitations

This scan ran inside the Codex execution environment. Network netlink access, systemd service enumeration, listener inspection, and container enumeration were restricted, so the host's complete service and port inventory still needs a host-level scan.

## Secret handling rule

Do not store passwords, private keys, API tokens, or database credentials in this file, shell history, Git, or ordinary Codex configuration. For future deployments, store secrets in an operator-approved secret manager or an encrypted file (for example, SOPS/age), and keep only the secret name, owner, service, and rotation metadata here. Deployment manifests should reference secret names or environment-file paths rather than embedding values.

## Follow-up inventory

Host-level follow-up during Dujiao-Next deployment (2026-09-09 local):

- Public IPv4: `179.255.156.111` (interface inspection).
- Existing Nginx sites: `api.sharesai.xyz` and `cn2-api.sharesai.xyz` on 80/443.
- Existing `sub2api-node4.service`: loopback 8080; state tunnel ports 25432 and 26379.
- Dujiao-Next v1.4.7 now runs on public 18080 via Nginx, app loopback 18081, dedicated Redis loopback 16379.
- 2026-09-09 PDT / 2026-09-10 UTC: Dujiao-Next also serves `https://fk.baibaihub.com` through Cloudflare and a dedicated Nginx 80/443 vhost. Let's Encrypt certificate and Certbot renewal are configured; legacy port 18080 remains available.
- Deployment and encrypted credential references: `deployments/dujiao-next/README.md`.
- New app, Redis, daily encrypted backup timer and existing API services were verified.

- Confirm host-level listening ports and active services.
- Confirm firewall policy and automatic security updates.
- Record public IP/DNS only if the operator supplies or explicitly requests that lookup.
- Record each deployed service with its domain, ports, data path, backup path, and secret references.

## 2026-09-12: Isolated Provider Hall Test

- Docker Engine 29.1.3 and Compose 2.40.3 installed; `docker.service` and `containerd.service` enabled.
- Sub2API provider-hall 0.1.179 (20260912 binary package) runs in Compose project `provider-hall-test`, app loopback `18082`, no public domain/vhost.
- Dedicated PostgreSQL 18.6 and Redis 7.4.11 containers have no published host ports and do not use production state tunnels.
- Deployment root `/opt/sub2api-provider-hall-test`; data/config in `data/`, persistent Docker volumes `sub2api-provider-hall-test-pgdata` and `sub2api-provider-hall-test-redisdata`.
- Root-only administrator login reference: `/opt/sub2api-provider-hall-test/credentials/login.json`; runtime secrets in adjacent `runtime.env` and `db-password`.
- Root-only local backups under `/opt/sub2api-provider-hall-test/backups`; initialized dump restored and verified on 2026-09-12.
- Documentation: `deployments/sub2api-provider-hall-test/README.md` and `DEPLOYMENT-RECORD.md`.
- Production Nginx virtual hosts, Sub2API systemd service and Dujiao-Next remain in place.

- 2026-09-13: Test app publishes `0.0.0.0:18082` for public testing; PostgreSQL and Redis remain internal-only.

- 2026-09-17: Provider Hall test deployment record consolidated with detailed pitfalls in `/opt/sub2api-provider-hall-test/README.md` and `/root/.codex/deployments/sub2api-provider-hall-test/DEPLOYMENT-RECORD.md`; no runtime or production service changes were made.
