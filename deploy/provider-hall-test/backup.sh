#!/bin/sh
set -eu
umask 077
base=/opt/sub2api-provider-hall-test
stamp=$(date -u +%Y%m%dT%H%M%SZ)
dest="$base/backups/$stamp"
mkdir -m 700 "$dest"
docker exec sub2api-provider-hall-test-postgres pg_dump -U sub2api_test -d sub2api_test -Fc > "$dest/database.dump.partial"
test -s "$dest/database.dump.partial"
docker exec -i sub2api-provider-hall-test-postgres pg_restore --list < "$dest/database.dump.partial" > "$dest/database.toc"
mv "$dest/database.dump.partial" "$dest/database.dump"
tar -czf "$dest/config-and-credentials.tar.gz" -C "$base" compose.yaml credentials data/config.yaml data/.installed
sha256sum "$dest/database.dump" "$dest/config-and-credentials.tar.gz" > "$dest/SHA256SUMS"
printf 'Backup completed: %s\n' "$dest"
