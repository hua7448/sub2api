import { Clock, ExternalLink, MessageCircle } from 'lucide-react'
import { Section } from './section'
import { SITE } from '@/lib/docs-data'

export function Contact() {
  return (
    <Section
      id="contact"
      eyebrow="Support"
      title="联系支持"
      description="充值、配置或账号问题，欢迎随时联系客服。"
    >
      <div className="grid gap-4 sm:grid-cols-3">
        <div className="rounded-xl border border-border bg-card/50 p-5 transition-all duration-300 hover:border-primary/30">
          <MessageCircle className="size-5 text-primary" />
          <p className="mt-3 font-mono text-[11px] uppercase tracking-wider text-muted-foreground">
            QQ 群
          </p>
          <p className="mt-1 text-lg font-semibold text-foreground">{SITE.qqGroup}</p>
        </div>
        <div className="rounded-xl border border-border bg-card/50 p-5 transition-all duration-300 hover:border-primary/30">
          <Clock className="size-5 text-primary" />
          <p className="mt-3 font-mono text-[11px] uppercase tracking-wider text-muted-foreground">
            工作时间
          </p>
          <p className="mt-1 text-lg font-semibold text-foreground">
            {SITE.hours}
          </p>
        </div>
        <a
          href={SITE.home}
          target="_blank"
          rel="noreferrer"
          className="group rounded-xl border border-border bg-card/50 p-5 transition-all duration-300 hover:border-primary/30"
        >
          <ExternalLink className="size-5 text-primary" />
          <p className="mt-3 font-mono text-[11px] uppercase tracking-wider text-muted-foreground">
            主站
          </p>
          <p className="mt-1 truncate text-sm font-medium text-primary underline-offset-4 group-hover:underline">
            {SITE.home}
          </p>
        </a>
      </div>

      <footer className="mt-12 border-t border-border/60 pt-6 text-center">
        <p className="font-mono text-xs text-muted-foreground">
          <span className="text-primary">SmartQ</span> · 纯血 Claude + Codex ·
          面向开发者的 AI 编码额度平台
        </p>
      </footer>
    </Section>
  )
}
