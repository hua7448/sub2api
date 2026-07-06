import { ClipboardCheck, KeyRound, UserPlus, Wallet } from 'lucide-react'
import { Section } from './section'
import { SITE } from '@/lib/docs-data'

const STEPS = [
  {
    icon: UserPlus,
    title: '注册 SmartQ',
    desc: 'SmartQ 目前开放注册，前往主站创建账号。',
  },
  {
    icon: Wallet,
    title: '联系客服充值',
    desc: `试运行阶段充值请联系客服 QQ：${SITE.qqGroup}。`,
  },
  {
    icon: KeyRound,
    title: '创建 API 密钥',
    desc: '在"API 密钥"页面创建密钥，并选择对应分组。',
  },
  {
    icon: ClipboardCheck,
    title: '复制使用配置',
    desc: '点击"使用密钥"，按客户端与系统复制配置。',
  },
]

export function QuickStart() {
  return (
    <Section
      id="quick-start"
      eyebrow="Quick Start"
      title="快速开始"
      description="四步即可开始使用。SmartQ 目前开放注册，试运行阶段充值请联系客服。"
    >
      <ol className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {STEPS.map((step, i) => {
          const Icon = step.icon
          return (
            <li
              key={step.title}
              className="group relative rounded-xl border border-border bg-card/50 p-5 transition-all duration-300 hover:-translate-y-1 hover:border-primary/30 hover:bg-card"
            >
              <div className="flex items-center justify-between">
                <span className="flex size-10 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary transition-transform duration-300 group-hover:scale-110">
                  <Icon className="size-5" />
                </span>
                <span className="font-pixel text-xs text-muted-foreground/50">
                  0{i + 1}
                </span>
              </div>
              <h3 className="mt-4 text-base font-semibold text-foreground">
                {step.title}
              </h3>
              <p className="mt-1.5 text-sm leading-relaxed text-muted-foreground">
                {step.desc}
              </p>
            </li>
          )
        })}
      </ol>
    </Section>
  )
}
