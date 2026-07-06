import { Download, ExternalLink } from 'lucide-react'
import { Section } from './section'
import { CodeBlock } from './code-block'
import { SimpleTabs } from './simple-tabs'

const CODEX_APP_URL = 'https://openai.com/codex/'
const WINDOWS_NODE_INSTALL = 'winget install OpenJS.NodeJS.LTS'
const MACOS_NODE_INSTALL = 'brew install node'
const CODEX_CLI_INSTALL = 'npm install -g @openai/codex@latest'

const CONFIG_TOML = `model_provider = "OpenAI"
model = "gpt-5.5"
review_model = "gpt-5.5"
model_reasoning_effort = "xhigh"
disable_response_storage = true
network_access = "enabled"
windows_wsl_setup_acknowledged = true

[model_providers.OpenAI]
name = "OpenAI"
base_url = "https://smartapi.mynatapp.cc"
wire_api = "responses"
requires_openai_auth = true

[features]
goals = true`

const AUTH_JSON = `{
  "OPENAI_API_KEY": "********"
}`

export function CodexConfig() {
  return (
    <Section
      id="codex-config"
      eyebrow="Codex"
      title="Codex 安装与配置"
      description="下载 Codex App，或通过 npm 安装 Codex CLI，然后编辑配置文件。"
    >
      <div className="mb-8 grid gap-4 lg:grid-cols-[1.2fr_0.8fr]">
        <div className="rounded-xl border border-border bg-card/50 p-5">
          <p className="font-mono text-xs uppercase tracking-[0.16em] text-primary">
            Codex CLI
          </p>
          <h3 className="mt-2 text-lg font-semibold text-foreground">
            通过 npm 安装
          </h3>
          <p className="mt-1.5 text-sm leading-6 text-muted-foreground">
            先安装 Node.js，再全局安装最新版 Codex CLI。
          </p>

          <div className="mt-4">
            <SimpleTabs
              tabs={[
                {
                  value: 'macos-install',
                  label: 'macOS',
                  content: (
                    <div className="grid gap-4">
                      <div>
                        <h4 className="mb-2 text-sm font-medium text-foreground">
                          1. 安装 Node.js
                        </h4>
                        <CodeBlock
                          filename="Terminal"
                          code={MACOS_NODE_INSTALL}
                        />
                      </div>
                      <div>
                        <h4 className="mb-2 text-sm font-medium text-foreground">
                          2. 安装 Codex CLI
                        </h4>
                        <CodeBlock
                          filename="Terminal"
                          code={CODEX_CLI_INSTALL}
                        />
                      </div>
                    </div>
                  ),
                },
                {
                  value: 'windows-install',
                  label: 'Windows',
                  content: (
                    <div className="grid gap-4">
                      <div>
                        <h4 className="mb-2 text-sm font-medium text-foreground">
                          1. 安装 Node.js
                        </h4>
                        <CodeBlock
                          filename="PowerShell"
                          code={WINDOWS_NODE_INSTALL}
                        />
                      </div>
                      <div>
                        <h4 className="mb-2 text-sm font-medium text-foreground">
                          2. 安装 Codex CLI
                        </h4>
                        <CodeBlock
                          filename="CMD（管理员）"
                          code={CODEX_CLI_INSTALL}
                        />
                      </div>
                    </div>
                  ),
                },
              ]}
            />
          </div>
        </div>

        <div className="rounded-xl border border-border bg-card/50 p-5">
          <Download className="size-5 text-primary" />
          <h3 className="mt-3 text-lg font-semibold text-foreground">
            Codex App
          </h3>
          <p className="mt-1.5 text-sm leading-6 text-muted-foreground">
            前往 OpenAI 官方页面下载并安装 Codex App。
          </p>
          <a
            href={CODEX_APP_URL}
            target="_blank"
            rel="noreferrer"
            className="mt-4 inline-flex items-center gap-2 rounded-lg bg-primary px-3.5 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/80"
          >
            官方下载
            <ExternalLink className="size-4" />
          </a>
        </div>
      </div>

      <h3 className="mb-4 text-lg font-semibold text-foreground">
        配置文件
      </h3>
      <SimpleTabs
        tabs={[
          {
            value: 'unix',
            label: 'macOS / Linux',
            content: (
              <>
                <div className="grid gap-4 lg:grid-cols-[1.2fr_1fr]">
                  <CodeBlock
                    filename="~/.codex/config.toml"
                    code={CONFIG_TOML}
                  />
                  <CodeBlock filename="~/.codex/auth.json" code={AUTH_JSON} />
                </div>
                <p className="mt-4 text-sm leading-6 text-muted-foreground">
                  请确保配置目录存在。macOS/Linux 用户可运行{' '}
                  <code className="font-mono text-foreground">
                    mkdir -p ~/.codex
                  </code>{' '}
                  创建目录。
                </p>
              </>
            ),
          },
          {
            value: 'windows',
            label: 'Windows',
            content: (
              <>
                <div className="grid gap-4 lg:grid-cols-[1.2fr_1fr]">
                  <CodeBlock
                    filename="%userprofile%\.codex/config.toml"
                    code={CONFIG_TOML}
                  />
                  <CodeBlock
                    filename="%userprofile%\.codex/auth.json"
                    code={AUTH_JSON}
                  />
                </div>
                <p className="mt-4 text-sm leading-6 text-muted-foreground">
                  按 Win+R，输入 %userprofile%\.codex
                  打开配置目录。如目录不存在，请先手动创建。
                </p>
              </>
            ),
          },
        ]}
      />
    </Section>
  )
}
