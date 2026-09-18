"""
curl 命令
curl -X POST "https://api.sharesai.xyz/v1/images/generations" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer 你的API_KEY" \
  -d '{
    "model": "gpt-image-2",
    "prompt": "一个美丽的小女孩在和妈妈跳舞，现实主义风格."
  }'

"""

import base64
import os

import requests

BASE_URL = os.getenv("OPENAI_BASE_URL", "https://api.sharesai.xyz").rstrip("/")
API_KEY = "sk-fd01043b565c19c408fa3d5e9d10aa780562d34549bf8f0aab6cbb70fe29d85a"

if not API_KEY:
    raise ValueError("Set OPENAI_API_KEY before running this script.")

if BASE_URL.endswith("/v1"):
    api_url = f"{BASE_URL}/images/generations"
else:
    api_url = f"{BASE_URL}/v1/images/generations"

prompt = """
生成一场直播pk赛，两个主播分别是丁真珍珠在卖锐刻五代，王源在卖芙蓉王
""".strip()

response = requests.post(
    api_url,
    headers={
        "Content-Type": "application/json",
        "Authorization": f"Bearer {API_KEY}",
    },
    json={
        "model": "gpt-image-2",
        "prompt": prompt,
    },
    timeout=120,
)
response.raise_for_status()

image_base64 = response.json()["data"][0]["b64_json"]
image_bytes = base64.b64decode(image_base64)

with open("otter.png", "wb") as f:
    f.write(image_bytes)

print("Image saved to otter.png")
