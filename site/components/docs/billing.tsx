import { Code2, Info, Terminal } from 'lucide-react'
import { Section } from './section'

export function Billing() {
  return (
    <Section
      id="billing"
      eyebrow="Billing"
      title="计费说明"
      description="不同客户端的额度消耗方式不同，请按实际使用场景预估消耗。"
    >
      <div className="grid gap-4 md:grid-cols-2">
        <article className="group rounded-xl border border-border bg-card/50 p-6 transition-all duration-300 hover:border-primary/30">
          <div className="flex items-center gap-3">
            <span className="flex size-10 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
              <Terminal className="size-5" />
            </span>
            <h3 className="text-lg font-semibold text-foreground">Claude Code</h3>
            <span className="ml-auto rounded-md bg-amber/15 px-2 py-1 font-mono text-xs text-amber">
              1.8X 倍率
            </span>
          </div>
          <ul className="mt-4 space-y-2.5 text-sm leading-relaxed text-muted-foreground">
            <li className="flex gap-2">
              <span className="mt-2 size-1 shrink-0 rounded-full bg-primary" />
              Claude 按 <span className="text-foreground">1.8X 倍率</span> 扣除额度。
            </li>
            <li className="flex gap-2">
              <span className="mt-2 size-1 shrink-0 rounded-full bg-primary" />
              长上下文、大量工具调用、长时间连续任务都会增加消耗。
            </li>
          </ul>
        </article>

        <article className="group rounded-xl border border-border bg-card/50 p-6 transition-all duration-300 hover:border-primary/30">
          <div className="flex items-center gap-3">
            <span className="flex size-10 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
              <Code2 className="size-5" />
            </span>
            <h3 className="text-lg font-semibold text-foreground">Codex</h3>
            <span className="ml-auto rounded-md bg-primary/15 px-2 py-1 font-mono text-xs text-primary">
              量大管饱
            </span>
          </div>
          <ul className="mt-4 space-y-2.5 text-sm leading-relaxed text-muted-foreground">
            <li className="flex gap-2">
              <span className="mt-2 size-1 shrink-0 rounded-full bg-primary" />
              试运行阶段主打 <span className="text-foreground">量大管饱</span>。
            </li>
            <li className="flex gap-2">
              <span className="mt-2 size-1 shrink-0 rounded-full bg-primary" />
              适合高频编码、审查、调试等场景。
            </li>
          </ul>
        </article>
      </div>

      <p className="mt-4 flex items-center gap-2 rounded-lg border border-amber/20 bg-amber/[0.06] px-4 py-3 text-sm text-amber">
        <Info className="size-4 shrink-0" />
        实际扣费和额度变化以 SmartQ 站内显示为准。
      </p>
    </Section>
  )
}
