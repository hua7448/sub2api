import { ExternalLink, MessageSquare } from 'lucide-react'
import { Section } from './section'
import { CodeBlock } from './code-block'
import { SimpleTabs } from './simple-tabs'
import { SITE } from '@/lib/docs-data'

const KIMI_DOWNLOAD_URL = 'https://www.moonshot.cn/kimi-cli'

const MACOS_INSTALL = `brew install node
npm install -g kimi-cli@latest`

const WINDOWS_INSTALL = `winget install OpenJS.NodeJS.LTS
npm install -g kimi-cli@latest`

const KIMI_CONFIG_TOML = `default_model = "my-proxy/kimi-for-coding"
default_permission_mode = "manual"
default_plan_mode = false
merge_all_available_skills = true
telemetry = true

[providers."my-proxy"]
type = "anthropic"
base_url = "${SITE.baseUrl}/"
api_key = "sk-********************************"

[models."my-proxy/kimi-for-coding"]
provider = "my-proxy"
model = "kimi-for-coding"
max_context_size = 262144
capabilities = ["thinking", "tool_use"]

[thinking]
enabled = true
effort = "high"

[loop_control]
max_retries_per_step = 3
reserved_context_size = 50000

[background]
max_running_tasks = 4
keep_alive_on_exit = false

[[permission.rules]]
decision = "allow"
pattern = "Read"

[[permission.rules]]
decision = "deny"
pattern = "Bash(rm -rf*)"`

export function KimiConfig() {
  return (
    <Section
      id="kimi-config"
      eyebrow="Kimi"
      title="Kimi CLI 安装与配置"
      description="在 Kimi Code 中通过 OpenAI 兼容协议接入 SmartQ，即可使用 kimi-for-coding 等模型。"
    >
      <div className="mb-8 rounded-xl border border-border bg-card/50 p-5">
        <p className="font-mono text-xs uppercase tracking-[0.16em] text-primary">
          Kimi CLI
        </p>
        <h3 className="mt-2 text-lg font-semibold text-foreground">
          通过 npm 安装
        </h3>
        <p className="mt-1.5 text-sm leading-6 text-muted-foreground">
          先安装 Node.js，再全局安装 Kimi CLI。也可以直接下载官方客户端。
        </p>

        <div className="mt-4">
          <SimpleTabs
            tabs={[
              {
                value: 'macos-install',
                label: 'macOS',
                content: (
                  <CodeBlock
                    filename="Terminal"
                    code={MACOS_INSTALL}
                  />
                ),
              },
              {
                value: 'windows-install',
                label: 'Windows',
                content: (
                  <CodeBlock
                    filename="PowerShell"
                    code={WINDOWS_INSTALL}
                  />
                ),
              },
            ]}
          />
        </div>

        <a
          href={KIMI_DOWNLOAD_URL}
          target="_blank"
          rel="noreferrer"
          className="mt-4 inline-flex items-center gap-2 rounded-lg bg-primary px-3.5 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/80"
        >
          官方下载
          <ExternalLink className="size-4" />
        </a>
      </div>

      <h3 className="mb-4 text-lg font-semibold text-foreground">
        配置文件
      </h3>
      <p className="mb-4 text-sm leading-relaxed text-muted-foreground">
        将以下内容写入 <code className="font-mono text-foreground">~/.kimi-code/config.toml</code>。
        请把 <code className="font-mono text-foreground">api_key</code> 替换为你自己的 SmartQ API Key。
      </p>

      <div className="flex items-start gap-3 rounded-lg border border-amber/25 bg-amber/[0.08] px-4 py-3 text-sm text-amber">
        <MessageSquare className="mt-0.5 size-4 shrink-0" />
        <span>
          当前示例中 API Key 已脱敏。配置时务必使用自己在 SmartQ 创建的密钥，不要分享或提交到公开仓库。
        </span>
      </div>

      <div className="mt-4">
        <CodeBlock
          filename="~/.kimi-code/config.toml"
          code={KIMI_CONFIG_TOML}
        />
      </div>
    </Section>
  )
}
