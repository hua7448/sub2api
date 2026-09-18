#!/usr/bin/env python3
"""
渠道 Azure Foundry 检测工具 — 非流式响应头硬指纹捕获 + 模型替换检测

与生产环境 channel-monitor 同构：
- Azure 后端的硬指纹只在**非流式**响应头里 (Apim-Request-Id / Azureai-* / Azureml-Served-By-Cluster / X-Ms-Region)
- new-api / one-api 之类代理在 stream 模式会剥光这些 header
- body 完全和真 Anthropic Max/官Key 一致 (msg_*, stop_details: null, inference_geo: "not_available")，
  只能靠 header 区分

用法:
    python3 detect_azure.py <base_url> <api_key> [model]
示例:
    python3 detect_azure.py http://15.204.212.18:3095 sk-xxx
    python3 detect_azure.py http://15.204.212.18:3095 sk-xxx claude-haiku-4-5-20251001
"""

import json
import secrets
import sys
import uuid
import urllib.error
import urllib.request
from urllib.parse import urlparse

DEFAULT_MODELS = [
    "claude-haiku-4-5-20251001",
    "claude-sonnet-4-6",
    "claude-opus-4-6",
    "claude-opus-4-7",
]

# Azure 硬指纹 — 任一存在即视为 Azure 后端
AZURE_HEADER_KEYS = [
    "apim-request-id",
    "azureai-processed-tier",
    "azureai-requested-tier",
    "azureml-served-by-cluster",
    "x-ms-region",
]

# 与 channel-monitor 同款的 per-model input_tokens 基线（仅 Anthropic 直连/Azure 适用）
EXPECTED_INPUT_TOKENS = {
    "claude-haiku-4-5-20251001": 24,
    "claude-sonnet-4-6": 25,
    "claude-opus-4-6": 25,
    "claude-opus-4-7": 38,
}


def build_headers(api_key: str, host: str) -> dict:
    """同 channel-monitor src/lib/api-tester.ts buildClaudeCodeHeaders。"""
    return {
        "Host": host,
        "User-Agent": "claude-cli/2.0.41 (external, cli)",
        "Connection": "keep-alive",
        "Accept": "application/json",
        "Accept-Encoding": "identity",
        "Content-Type": "application/json",
        "X-Stainless-Retry-Count": "0",
        "X-Stainless-Timeout": "600",
        "X-Stainless-Lang": "js",
        "X-Stainless-Package-Version": "0.60.0",
        "X-Stainless-OS": "Windows",
        "X-Stainless-Arch": "x64",
        "X-Stainless-Runtime": "node",
        "X-Stainless-Runtime-Version": "v22.14.0",
        "anthropic-dangerous-direct-browser-access": "true",
        "anthropic-version": "2023-06-01",
        "x-app": "cli",
        "Authorization": f"Bearer {api_key}",
        "anthropic-beta": "claude-code-20250219,interleaved-thinking-2025-05-14,fine-grained-tool-streaming-2025-05-14",
        "accept-language": "*",
        "sec-fetch-mode": "cors",
    }


def build_payload(model: str) -> dict:
    return {
        "model": model,
        "messages": [{"role": "user", "content": "你好"}],
        "temperature": 1,
        "system": [{
            "type": "text",
            "text": "You are Claude Code, Anthropic's official CLI for Claude.",
            "cache_control": {"type": "ephemeral"},
        }],
        "metadata": {"user_id": f"user_{secrets.token_hex(32)}_account__session_{uuid.uuid4()}"},
        "max_tokens": 1,
    }


def detect_azure_from_headers(headers) -> tuple:
    """返回 (is_azure: bool, hits: list[str])"""
    hits = []
    for key in AZURE_HEADER_KEYS:
        val = headers.get(key)
        if val:
            hits.append(f"{key}: {val}")
    return (len(hits) > 0, hits)


def normalize_model(m: str) -> str:
    """同 api-tester.ts: anthropic[/.] 前缀 / -v\\d+(:\\d+)? / -\\d{8} 后缀去除。"""
    import re
    m = re.sub(r"^anthropic[/.]", "", m)
    m = re.sub(r"-v\d+(?::\d+)?$", "", m)
    m = re.sub(r"-\d{8}$", "", m)
    return m


def probe_model(base_url: str, api_key: str, model: str, timeout: int = 15) -> dict:
    """非流式 probe，返回详细诊断字段。"""
    parsed = urlparse(base_url)
    host = parsed.netloc
    url = f"{base_url.rstrip('/')}/v1/messages?beta=true"

    req = urllib.request.Request(
        url,
        data=json.dumps(build_payload(model)).encode("utf-8"),
        headers=build_headers(api_key, host),
        method="POST",
    )

    result = {
        "model": model,
        "http_status": None,
        "is_azure": False,
        "azure_hits": [],
        "body_id": None,
        "body_model": None,
        "model_matches": None,
        "input_tokens": None,
        "input_tokens_expected": EXPECTED_INPUT_TOKENS.get(model),
        "input_tokens_valid": None,
        "stop_details_present": None,
        "inference_geo": None,
        "error": None,
        "raw_headers": {},
    }

    try:
        with urllib.request.urlopen(req, timeout=timeout) as resp:
            result["http_status"] = resp.status
            result["raw_headers"] = dict(resp.headers.items())

            is_azure, hits = detect_azure_from_headers(resp.headers)
            result["is_azure"] = is_azure
            result["azure_hits"] = hits

            body = resp.read().decode("utf-8", errors="replace")
            try:
                data = json.loads(body)
            except json.JSONDecodeError:
                result["error"] = f"Invalid JSON body (first 200 chars): {body[:200]}"
                return result

            result["body_id"] = data.get("id")
            result["body_model"] = data.get("model")
            if isinstance(result["body_model"], str):
                result["model_matches"] = normalize_model(result["body_model"]) == normalize_model(model)

            usage = data.get("usage", {}) if isinstance(data.get("usage"), dict) else {}
            result["input_tokens"] = usage.get("input_tokens")
            if result["input_tokens_expected"] is not None and isinstance(result["input_tokens"], int):
                result["input_tokens_valid"] = result["input_tokens"] == result["input_tokens_expected"]

            result["stop_details_present"] = "stop_details" in data
            result["inference_geo"] = usage.get("inference_geo")

    except urllib.error.HTTPError as e:
        result["http_status"] = e.code
        result["raw_headers"] = dict(e.headers.items()) if e.headers else {}
        is_azure, hits = detect_azure_from_headers(e.headers or {})
        result["is_azure"] = is_azure
        result["azure_hits"] = hits
        try:
            result["error"] = e.read().decode("utf-8", errors="replace")[:300]
        except Exception:
            result["error"] = str(e)
    except Exception as e:
        result["error"] = str(e)

    return result


def print_result(r: dict) -> None:
    print(f"\n  ═══ {r['model']} ═══")
    print(f"  HTTP {r['http_status']}")

    if r["error"]:
        print(f"  ⚠️  {r['error']}")
        return

    # Azure 硬指纹
    if r["is_azure"]:
        print(f"  ✅ Azure 后端 — 命中 {len(r['azure_hits'])} 个硬指纹:")
        for hit in r["azure_hits"]:
            print(f"     • {hit}")
    else:
        print(f"  ❌ 未检测到 Azure header")
        # 列出出现的代理头方便排查
        proxy_headers = {k: v for k, v in r["raw_headers"].items()
                         if k.lower().startswith(("x-oneapi", "x-new-api", "server"))}
        if proxy_headers:
            print(f"     代理标识:")
            for k, v in proxy_headers.items():
                print(f"     • {k}: {v}")

    # Body validation
    print(f"  Body 校验:")
    print(f"     id: {r['body_id']}")
    print(f"     model: {r['body_model']}")
    if r["model_matches"] is not None:
        mark = "✓" if r["model_matches"] else "✗ 模型替换"
        print(f"     model_matches: {mark}")
    if r["input_tokens"] is not None and r["input_tokens_expected"] is not None:
        mark = "✓" if r["input_tokens_valid"] else "✗ 异常 token"
        print(f"     input_tokens: {r['input_tokens']} / 基线 {r['input_tokens_expected']} {mark}")
    print(f"     stop_details: {'✓' if r['stop_details_present'] else '✗'}")
    print(f"     inference_geo: {r['inference_geo']}")


def main():
    if len(sys.argv) < 3:
        print("用法: python3 detect_azure.py <base_url> <api_key> [model]")
        sys.exit(1)

    base_url = sys.argv[1].rstrip("/")
    api_key = sys.argv[2]
    models = [sys.argv[3]] if len(sys.argv) > 3 else DEFAULT_MODELS

    print(f"目标: {base_url}")
    print(f"模型: {', '.join(models)}")
    print("=" * 60)

    azure_count = 0
    for m in models:
        r = probe_model(base_url, api_key, m)
        print_result(r)
        if r["is_azure"]:
            azure_count += 1

    print("\n" + "=" * 60)
    print(f"总结: {azure_count}/{len(models)} 模型命中 Azure 硬指纹")
    if azure_count == len(models):
        print("  ✅ 该渠道为 **Azure Foundry Claude**（建议设 azure_check=1）")
    elif azure_count == 0:
        print("  ❌ 未检测到 Azure 后端 — 可能是 Max/官Key/Bedrock/OpenRouter")
    else:
        print(f"  ⚠️  混合后端 — 部分模型走 Azure 部分不走（中转站负载均衡）")


if __name__ == "__main__":
    main()
