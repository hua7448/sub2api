#!/usr/bin/env python3
"""Monitor usable OpenAI OAuth accounts and notify WeCom when inventory is low."""

from __future__ import annotations

import argparse
import json
import os
import shutil
import subprocess
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path
from typing import Any


DEFAULT_THRESHOLD = 20
DEFAULT_INTERVAL_SECONDS = 60
DEFAULT_MESSAGE_TEMPLATE = "账号目前数量为{count}，请及时补充"
DEFAULT_WEBHOOK_URL = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=27012de3-051a-45ba-b75c-acb59599e0a0"


BASE_COUNT_SQL = """
SELECT COUNT(*)
FROM accounts
WHERE deleted_at IS NULL
  AND lower(platform) = lower(%s)
  AND lower(type) = lower(%s)
  AND status = 'active'
  AND schedulable IS TRUE
  AND (auto_pause_on_expired IS NOT TRUE OR expires_at IS NULL OR expires_at > NOW())
  AND (rate_limit_reset_at IS NULL OR rate_limit_reset_at <= NOW())
  AND (overload_until IS NULL OR overload_until <= NOW())
  AND (temp_unschedulable_until IS NULL OR temp_unschedulable_until <= NOW())
"""


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Notify WeCom when usable OpenAI OAuth accounts are below threshold."
    )
    parser.add_argument("--config", help="Sub2API config.yaml path. Defaults to common locations.")
    parser.add_argument(
        "--webhook-url",
        default=os.getenv("WECHAT_WEBHOOK_URL", DEFAULT_WEBHOOK_URL),
        help="WeCom webhook URL.",
    )
    parser.add_argument("--threshold", type=int, default=int(os.getenv("OPENAI_OAUTH_ACCOUNT_THRESHOLD", DEFAULT_THRESHOLD)))
    parser.add_argument("--platform", default=os.getenv("MONITOR_ACCOUNT_PLATFORM", "openai"))
    parser.add_argument("--account-type", default=os.getenv("MONITOR_ACCOUNT_TYPE", "oauth"))
    parser.add_argument("--plan-type", default=os.getenv("MONITOR_PLAN_TYPE"), help="Optional credentials.plan_type filter.")
    parser.add_argument("--message-template", default=os.getenv("MONITOR_MESSAGE_TEMPLATE", DEFAULT_MESSAGE_TEMPLATE))
    parser.add_argument("--dry-run", action="store_true", help="Print the decision without sending WeCom notification.")
    parser.add_argument("--connect-timeout", type=int, default=int(os.getenv("MONITOR_DB_TIMEOUT", "15")))
    parser.add_argument(
        "--interval-seconds",
        type=int,
        default=int(os.getenv("MONITOR_INTERVAL_SECONDS", DEFAULT_INTERVAL_SECONDS)),
        help="Loop interval in seconds. Default: 60.",
    )
    parser.add_argument("--once", action="store_true", help="Run one check and exit.")
    return parser.parse_args()


def count_sql_and_params(args: argparse.Namespace) -> tuple[str, tuple[str, ...]]:
    params = [args.platform, args.account_type]
    sql = BASE_COUNT_SQL
    if args.plan_type:
        sql += "  AND lower(trim(coalesce(credentials->>'plan_type', ''))) = lower(%s)\n"
        params.append(args.plan_type)
    return sql, tuple(params)


def common_config_paths() -> list[Path]:
    repo_root = Path(__file__).resolve().parents[1]
    paths: list[Path] = []
    for value in (os.getenv("SUB2API_CONFIG"), os.getenv("CONFIG_FILE")):
        if value:
            paths.append(Path(value))
    data_dir = os.getenv("DATA_DIR")
    if data_dir:
        paths.append(Path(data_dir) / "config.yaml")
    paths.extend(
        [
            Path("/etc/sub2api/config.yaml"),
            Path("/opt/sub2api/config.yaml"),
            Path("/opt/sub2api/config/config.yaml"),
            Path("/root/sub2api/backend/config.yaml"),
            Path("/root/sub2api/config.yaml"),
            repo_root / "config.yaml",
            repo_root / "backend" / "config.yaml",
            Path("config.yaml"),
            Path("backend/config.yaml"),
            Path("/app/data/config.yaml"),
        ]
    )
    return paths


def strip_yaml_value(value: str) -> str:
    value = value.split(" #", 1)[0].strip()
    if len(value) >= 2 and value[0] == value[-1] and value[0] in {"'", '"'}:
        value = value[1:-1]
    return value


def load_database_block(path: Path | None) -> dict[str, Any]:
    if path is not None and not path.exists():
        raise FileNotFoundError(f"Config file not found: {path}")

    if path is None:
        for candidate in common_config_paths():
            if candidate.exists():
                path = candidate
                break
    if path is None or not path.exists():
        return {}

    text = path.read_text(encoding="utf-8")

    try:
        import yaml  # type: ignore

        data = yaml.safe_load(text) or {}
        database = data.get("database") or {}
        if isinstance(database, dict):
            return dict(database)
    except Exception:
        pass

    database: dict[str, Any] = {}
    in_database = False
    for raw_line in text.splitlines():
        if not raw_line.strip() or raw_line.lstrip().startswith("#"):
            continue
        if raw_line.startswith("database:"):
            in_database = True
            continue
        if in_database and raw_line and not raw_line.startswith((" ", "\t")):
            break
        if in_database and ":" in raw_line:
            key, value = raw_line.split(":", 1)
            database[key.strip()] = strip_yaml_value(value)
    return database


def parse_database_url(url: str) -> dict[str, Any]:
    parsed = urllib.parse.urlparse(url)
    query = urllib.parse.parse_qs(parsed.query)
    return {
        "host": parsed.hostname or "",
        "port": parsed.port or 5432,
        "user": urllib.parse.unquote(parsed.username or ""),
        "password": urllib.parse.unquote(parsed.password or ""),
        "dbname": parsed.path.lstrip("/"),
        "sslmode": (query.get("sslmode") or ["prefer"])[0],
    }


def database_config(config_path: str | None) -> dict[str, Any]:
    config = {
        "host": "localhost",
        "port": 5432,
        "user": "postgres",
        "password": "postgres",
        "dbname": "sub2api",
        "sslmode": "prefer",
    }

    database_url = os.getenv("DATABASE_URL")
    if database_url:
        config.update({k: v for k, v in parse_database_url(database_url).items() if v not in ("", None)})

    config.update({k: v for k, v in load_database_block(Path(config_path) if config_path else None).items() if v not in ("", None)})

    env_map = {
        "host": "DATABASE_HOST",
        "port": "DATABASE_PORT",
        "user": "DATABASE_USER",
        "password": "DATABASE_PASSWORD",
        "dbname": "DATABASE_DBNAME",
        "sslmode": "DATABASE_SSLMODE",
    }
    for key, env_name in env_map.items():
        value = os.getenv(env_name)
        if value not in (None, ""):
            config[key] = value

    config["port"] = int(config.get("port") or 5432)
    return config


def count_with_psycopg(config: dict[str, Any], args: argparse.Namespace) -> int | None:
    try:
        try:
            import psycopg  # type: ignore

            with psycopg.connect(
                host=config["host"],
                port=config["port"],
                user=config["user"],
                password=config.get("password") or None,
                dbname=config["dbname"],
                sslmode=config.get("sslmode") or "prefer",
                connect_timeout=args.connect_timeout,
            ) as conn:
                with conn.cursor() as cur:
                    sql, params = count_sql_and_params(args)
                    cur.execute(sql, params)
                    row = cur.fetchone()
                    return int(row[0])
        except ModuleNotFoundError:
            import psycopg2  # type: ignore

            with psycopg2.connect(
                host=config["host"],
                port=config["port"],
                user=config["user"],
                password=config.get("password") or None,
                dbname=config["dbname"],
                sslmode=config.get("sslmode") or "prefer",
                connect_timeout=args.connect_timeout,
            ) as conn:
                with conn.cursor() as cur:
                    sql, params = count_sql_and_params(args)
                    cur.execute(sql, params)
                    row = cur.fetchone()
                    return int(row[0])
    except ModuleNotFoundError:
        return None


def sql_literal(value: str) -> str:
    return "'" + value.replace("'", "''") + "'"


def count_with_psql(config: dict[str, Any], args: argparse.Namespace) -> int:
    if not shutil.which("psql"):
        raise RuntimeError("Neither psycopg/psycopg2 nor psql is available.")

    sql, params = count_sql_and_params(args)
    query = sql
    for param in params:
        query = query.replace("%s", sql_literal(param), 1)
    env = os.environ.copy()
    if config.get("password"):
        env["PGPASSWORD"] = str(config["password"])
    if config.get("sslmode"):
        env["PGSSLMODE"] = str(config["sslmode"])
    env["PGCONNECT_TIMEOUT"] = str(args.connect_timeout)

    command = [
        "psql",
        "-X",
        "-A",
        "-t",
        "-v",
        "ON_ERROR_STOP=1",
        "-h",
        str(config["host"]),
        "-p",
        str(config["port"]),
        "-U",
        str(config["user"]),
        "-d",
        str(config["dbname"]),
        "-c",
        query,
    ]
    result = subprocess.run(
        command,
        env=env,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        timeout=args.connect_timeout + 5,
        check=False,
    )
    if result.returncode != 0:
        target = f"{config['user']}@{config['host']}:{config['port']}/{config['dbname']}"
        detail = result.stderr.strip() or result.stdout.strip() or "no error output"
        raise RuntimeError(f"psql failed for {target}: {detail}")
    return int(result.stdout.strip())


def count_accounts(config: dict[str, Any], args: argparse.Namespace) -> int:
    count = count_with_psycopg(config, args)
    if count is not None:
        return count
    return count_with_psql(config, args)


def send_wecom_text(webhook_url: str, content: str) -> None:
    payload = json.dumps(
        {"msgtype": "text", "text": {"content": content}},
        ensure_ascii=False,
    ).encode("utf-8")
    request = urllib.request.Request(
        webhook_url,
        data=payload,
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    try:
        with urllib.request.urlopen(request, timeout=10) as response:
            body = response.read().decode("utf-8", errors="replace")
    except urllib.error.URLError as exc:
        raise RuntimeError(f"WeCom webhook request failed: {exc}") from exc

    try:
        result = json.loads(body)
    except json.JSONDecodeError as exc:
        raise RuntimeError(f"WeCom webhook returned non-JSON response: {body}") from exc

    if result.get("errcode") != 0:
        raise RuntimeError(f"WeCom webhook returned error: {result}")


def run_once(args: argparse.Namespace) -> int:
    config = database_config(args.config)
    count = count_accounts(config, args)

    if count >= args.threshold:
        filter_label = f"{args.platform}/{args.account_type}"
        if args.plan_type:
            filter_label += f"/{args.plan_type}"
        print(f"OK: {filter_label} usable account count is {count}, threshold is {args.threshold}.")
        return 0

    message = args.message_template.format(count=count, threshold=args.threshold)
    if args.dry_run:
        print(f"DRY RUN: would send notification: {message}")
        return 0
    if not args.webhook_url:
        print("WECHAT_WEBHOOK_URL or --webhook-url is required when notification is needed.", file=sys.stderr)
        return 2

    send_wecom_text(args.webhook_url, message)
    print(f"ALERT SENT: {message}")
    return 0


def main() -> int:
    args = parse_args()
    if args.threshold <= 0:
        print("--threshold must be greater than 0", file=sys.stderr)
        return 2
    if args.interval_seconds <= 0:
        print("--interval-seconds must be greater than 0", file=sys.stderr)
        return 2

    if args.once:
        return run_once(args)

    print(f"Monitoring every {args.interval_seconds}s. Press Ctrl+C to stop.")
    while True:
        try:
            run_once(args)
        except Exception as exc:
            print(f"ERROR: {exc}", file=sys.stderr)
        time.sleep(args.interval_seconds)


if __name__ == "__main__":
    raise SystemExit(main())
