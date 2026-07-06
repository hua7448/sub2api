import { Gift } from 'lucide-react'
import { Section } from './section'

export function InviteReward() {
  return (
    <Section
      id="invite"
      eyebrow="Referral"
      title="邀请奖励"
      description="邀请一个新用户完成充值后，邀请人可获得 10 USD 站内额度奖励。邀请奖励没有次数上限。"
    >
      <div className="relative overflow-hidden rounded-2xl border border-primary/20 bg-card/50 p-8 border-glow">
        <div className="bg-grid bg-grid-fade pointer-events-none absolute inset-0 opacity-60" />
        <div className="relative flex flex-col items-center gap-6 sm:flex-row sm:gap-8">
          <span className="flex size-20 items-center justify-center rounded-2xl border border-primary/25 bg-primary/10 text-primary animate-float">
            <Gift className="size-9" />
          </span>
          <div className="text-center sm:text-left">
            <p className="font-mono text-xs uppercase tracking-[0.2em] text-muted-foreground">
              每邀请一位完成充值
            </p>
            <p className="mt-2 text-4xl font-semibold tracking-tight text-foreground md:text-5xl">
              <span className="text-primary text-glow">10</span>{' '}
              <span className="text-2xl text-muted-foreground md:text-3xl">
                USD
              </span>
            </p>
            <p className="mt-2 text-sm text-muted-foreground">
              站内额度奖励 · 无次数上限
            </p>
          </div>
        </div>
      </div>
    </Section>
  )
}
