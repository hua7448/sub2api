import { Section } from './section'
import { CodeBlock } from './code-block'
import { SimpleTabs } from './simple-tabs'
import { MessageSquare } from 'lucide-react'

const WINDOWS_NODE_INSTALL = 'winget install OpenJS.NodeJS.LTS'
const MACOS_NODE_INSTALL = 'brew install node'
const CLAUDE_CODE_INSTALL =
  'npm install -g @anthropic-ai/claude-code@latest'

const UNIX = `export ANTHROPIC_BASE_URL="https://smartapi.mynatapp.cc"
export ANTHROPIC_AUTH_TOKEN="********"
export CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1

claude`

const CMD = `set ANTHROPIC_BASE_URL=https://smartapi.mynatapp.cc
set ANTHROPIC_AUTH_TOKEN=********
set CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1

claude`

const POWERSHELL = `$env:ANTHROPIC_BASE_URL="https://smartapi.mynatapp.cc"
$env:ANTHROPIC_AUTH_TOKEN="********"
$env:CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1

claude`

const SETTINGS = `{
  "env": {
    "ANTHROPIC_BASE_URL": "https://smartapi.mynatapp.cc",
    "ANTHROPIC_AUTH_TOKEN": "********",
    "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": "1",
    "CLAUDE_CODE_ATTRIBUTION_HEADER": "0"
  }
}`

const GLM_SETTINGS = `{
  "env": {
    "ANTHROPIC_BASE_URL": "https://smartapi.mynatapp.cc",
    "ANTHROPIC_AUTH_TOKEN": "********",
    "CLAUDE_CODE_USE_VERTEX": "0",
    "CLAUDE_CODE_ATTRIBUTION_HEADER": "0"
  }
}`

const GLM_NOTE = `在 Claude Code 中使用 GLM 模型时，保持 ANTHROPIC_BASE_URL 指向 SmartQ，并将账号或分组映射到 GLM 模型即可。Claude Code 会通过 SmartQ 的路由层转发到对应上游。`

export function ClaudeConfig() {
  return (
    <Section
      id="claude-config"
      eyebrow="Claude Code"
      title="Claude Code 安装与配置"
      description="通过 npm 安装 Claude Code，然后按操作系统完成环境变量配置。"
    >
      <div className="mb-8 rounded-xl border border-border bg-card/50 p-5">
        <p className="font-mono text-xs uppercase tracking-[0.16em] text-primary">
          Claude Code CLI
        </p>
        <h3 className="mt-2 text-lg font-semibold text-foreground">
          通过 npm 安装
        </h3>
        <p className="mt-1.5 text-sm leading-6 text-muted-foreground">
          先安装 Node.js，再全局安装最新版 Claude Code。
        </p>

        <div className="mt-4">
          <SimpleTabs
            tabs={[
              {
                value: 'macos-install',
                label: 'macOS',
                content: (
                  <div className="grid gap-4 md:grid-cols-2">
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
                        2. 安装 Claude Code
                      </h4>
                      <CodeBlock
                        filename="Terminal"
                        code={CLAUDE_CODE_INSTALL}
                      />
                    </div>
                  </div>
                ),
              },
              {
                value: 'windows-install',
                label: 'Windows',
                content: (
                  <div className="grid gap-4 md:grid-cols-2">
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
                        2. 安装 Claude Code
                      </h4>
                      <CodeBlock
                        filename="CMD（管理员）"
                        code={CLAUDE_CODE_INSTALL}
                      />
                    </div>
                  </div>
                ),
              },
            ]}
          />
        </div>
      </div>

      <h3 className="mb-4 text-lg font-semibold text-foreground">
        环境变量配置
      </h3>
      <SimpleTabs
        tabs={[
          {
            value: 'unix',
            label: 'macOS / Linux',
            content: <CodeBlock filename="terminal" code={UNIX} />,
          },
          {
            value: 'cmd',
            label: 'Windows CMD',
            content: <CodeBlock filename="cmd.exe" code={CMD} />,
          },
          {
            value: 'powershell',
            label: 'PowerShell',
            content: <CodeBlock filename="PowerShell" code={POWERSHELL} />,
          },
        ]}
      />

      <div className="mt-6">
        <h3 className="mb-3 font-mono text-sm text-muted-foreground">
          settings.json 示例
        </h3>
        <CodeBlock filename="settings.json" code={SETTINGS} />
      </div>

      <div className="mt-10 rounded-xl border border-border bg-card/50 p-5">
        <p className="font-mono text-xs uppercase tracking-[0.16em] text-primary">
          GLM
        </p>
        <h3 className="mt-2 text-lg font-semibold text-foreground">
          在 Claude Code 中使用 GLM
        </h3>
        <p className="mt-1.5 text-sm leading-6 text-muted-foreground">
          {GLM_NOTE}
        </p>

        <div className="mt-4 flex items-start gap-3 rounded-lg border border-amber/25 bg-amber/[0.08] px-4 py-3 text-sm text-amber">
          <MessageSquare className="mt-0.5 size-4 shrink-0" />
          <span>
            请把示例中的 <code className="font-mono text-foreground">ANTHROPIC_AUTH_TOKEN</code> 替换为你自己的 SmartQ API Key。
          </span>
        </div>

        <div className="mt-4">
          <h4 className="mb-2 text-sm font-medium text-foreground">
            settings.json（GLM 模式）
          </h4>
          <CodeBlock filename="settings.json" code={GLM_SETTINGS} />
        </div>
      </div>
    </Section>
  )
}
