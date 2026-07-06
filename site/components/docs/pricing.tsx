import { Sparkles, ExternalLink } from 'lucide-react'
import { Section } from './section'
import { PRICING } from '@/lib/docs-data'
import { cn } from '@/lib/utils'

const PAY_URL = 'https://pay.ldxp.cn/shop/PBQVKSJX'

export function Pricing() {
  return (
    <Section
      id="pricing"
      eyebrow="Pricing"
      title="充值套餐"
      description="站内额度以 USD 展示，充值使用人民币支付。点击下方按钮前往站内充值页面完成支付。"
    >
      <div className="mb-6 flex flex-col gap-3 rounded-lg border border-primary/20 bg-primary/[0.06] px-4 py-3 sm:flex-row sm:items-center sm:justify-between">
        <p className="text-sm text-foreground">
          <span className="font-semibold">准备充值？</span>
          <span className="ml-1 text-muted-foreground">
            点击下方按钮前往充值页面，完成后额度会自动到账。
          </span>
        </p>
        <a
          href={PAY_URL}
          target="_blank"
          rel="noreferrer"
          className="inline-flex items-center justify-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/80"
        >
          前往充值
          <ExternalLink className="size-4" />
        </a>
      </div>

      <div className="overflow-hidden rounded-xl border border-border bg-card/40">
        {/* header row */}
        <div className="hidden grid-cols-[1fr_1fr_1fr] gap-4 border-b border-border bg-card/60 px-5 py-3 font-mono text-[11px] uppercase tracking-wider text-muted-foreground sm:grid">
          <span>站内额度</span>
          <span>充值价格</span>
          <span className="text-right">折扣</span>
        </div>
        <ul className="divide-y divide-border">
          {PRICING.map((tier) => (
            <li
              key={tier.usd}
              className={cn(
                'grid grid-cols-2 items-center gap-4 px-5 py-4 transition-colors sm:grid-cols-[1fr_1fr_1fr]',
                tier.highlight
                  ? 'bg-primary/[0.06] hover:bg-primary/10'
                  : 'hover:bg-secondary/40',
              )}
            >
              <div className="flex items-center gap-2">
                <span className="text-lg font-semibold text-foreground">
                  {tier.usd}
                  <span className="ml-1 text-sm font-normal text-muted-foreground">
                    USD
                  </span>
                </span>
                {tier.highlight ? (
                  <span className="inline-flex items-center gap-1 rounded-full border border-primary/30 bg-primary/10 px-2 py-0.5 font-mono text-[10px] text-primary">
                    <Sparkles className="size-3" />
                    超值
                  </span>
                ) : null}
              </div>
              <div className="font-mono text-sm text-foreground sm:text-base">
                {tier.cny}
              </div>
              <div className="col-span-2 sm:col-span-1 sm:text-right">
                <span
                  className={cn(
                    'inline-flex rounded-md px-2.5 py-1 font-mono text-xs',
                    tier.discount === '不打折'
                      ? 'bg-secondary/60 text-muted-foreground'
                      : 'bg-amber/15 text-amber',
                  )}
                >
                  {tier.discount}
                </span>
              </div>
            </li>
          ))}
        </ul>
      </div>
    </Section>
  )
}
