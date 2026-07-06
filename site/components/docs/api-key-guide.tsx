import { TriangleAlert } from 'lucide-react'
import { Section } from './section'

const STEPS = [
  '进入"API 密钥"页面',
  '找到要使用的密钥',
  '点击右侧"使用密钥"',
  'Codex 分组选择 Codex CLI，Claude 分组选择 Claude Code',
  '按系统选择 macOS / Linux、Windows CMD 或 PowerShell',
  '点击复制，把配置粘贴到对应文件或终端',
]

export function ApiKeyGuide() {
  return (
    <Section
      id="api-key"
      eyebrow="Core · API Key"
      title="使用 API 密钥教程"
      description="这是页面核心内容。推荐直接使用站内的「使用密钥」功能生成配置，省去手动编辑文件。"
    >
      <div className="rounded-xl border border-primary/20 bg-card/50 p-6 border-glow">
        <h3 className="flex items-center gap-2 text-base font-semibold text-foreground">
          <span className="font-mono text-xs text-primary">推荐方式</span>
          <span className="text-muted-foreground">/</span>
          使用站内"使用密钥"
        </h3>

        <ol className="mt-5 space-y-3">
          {STEPS.map((step, i) => (
            <li key={step} className="flex items-start gap-3.5">
              <span className="flex size-7 shrink-0 items-center justify-center rounded-md border border-primary/25 bg-primary/10 font-mono text-xs text-primary">
                {i + 1}
              </span>
              <span className="pt-0.5 text-sm leading-relaxed text-foreground/90">
                {step}
              </span>
            </li>
          ))}
        </ol>
      </div>

      <p className="mt-4 flex items-start gap-2.5 rounded-lg border border-amber/25 bg-amber/[0.08] px-4 py-3 text-sm text-amber">
        <TriangleAlert className="mt-0.5 size-4 shrink-0" />
        <span>
          如果站内"使用密钥"弹窗和本文档示例不一致，
          <span className="font-medium">请以站内弹窗为准</span>。
        </span>
      </p>
    </Section>
  )
}
