'use client'

import { useState } from 'react'
import { ArrowRight, Check, Copy, Key, Terminal } from 'lucide-react'
import { SITE } from '@/lib/docs-data'
import { copyText } from '@/lib/copy'
import { scrollToSection } from '@/lib/scroll'

type MetaItem = {
  label: string
  value: string
  href?: string
  mono?: boolean
  copy?: boolean
}

const META: MetaItem[] = [
  { label: '主站', value: SITE.home, href: SITE.home },
  { label: 'API Base URL', value: SITE.baseUrl, mono: true, copy: true },
  { label: '支持客户端', value: SITE.clients.join(' · ') },
  { label: '汇率', value: '全站 1 元人民币 = 1 美元额度' },
  { label: '客服 QQ', value: SITE.qqGroup },
  { label: '工作时间', value: SITE.hours },
]

export function Hero() {
  const [copied, setCopied] = useState(false)

  async function handleCopyBaseUrl() {
    const ok = await copyText(SITE.baseUrl)
    if (ok) {
      setCopied(true)
      setTimeout(() => setCopied(false), 1600)
    }
  }

  return (
    <section id="overview" className="relative scroll-mt-20 pb-12 pt-10 md:pt-16">
      <div className="grid items-center gap-10 lg:grid-cols-[1.05fr_0.95fr]">
        {/* Left: copy */}
        <div>
          <span className="inline-flex items-center gap-2 rounded-full border border-primary/25 bg-primary/10 px-3 py-1 font-mono text-[11px] uppercase tracking-[0.18em] text-primary">
            <span className="size-1.5 animate-pulse rounded-full bg-primary" />
            Claude · Codex · Kimi · GLM
          </span>

          <h1 className="mt-5 text-balance text-4xl font-semibold leading-tight tracking-tight text-foreground md:text-5xl">
            SmartQ 客户使用文档
          </h1>
          <p className="mt-4 max-w-xl text-pretty text-base leading-relaxed text-muted-foreground md:text-lg">
            面向开发者的 AI 编码额度平台。支持 Claude Code、Codex、Kimi CLI
            与 GLM，配置清晰、可复制、可快速查找。
          </p>

          <dl className="mt-8 grid max-w-xl grid-cols-1 gap-px overflow-hidden rounded-xl border border-border bg-border sm:grid-cols-2">
            {META.map((item) => (
              <div
                key={item.label}
                className="group/meta bg-card px-4 py-3 transition-colors duration-200 hover:bg-card/70"
              >
                <dt className="font-mono text-[11px] uppercase tracking-wider text-muted-foreground">
                  {item.label}
                </dt>
                <dd className="mt-1 flex items-center gap-2 text-sm text-foreground">
                  {item.href ? (
                    <a
                      href={item.href}
                      target="_blank"
                      rel="noreferrer"
                      className="truncate text-primary underline-offset-4 hover:underline"
                    >
                      {item.value}
                    </a>
                  ) : (
                    <span
                      className={
                        item.mono
                          ? 'truncate font-mono text-primary'
                          : 'truncate'
                      }
                    >
                      {item.value}
                    </span>
                  )}
                  {item.copy ? (
                    <button
                      type="button"
                      onClick={handleCopyBaseUrl}
                      aria-label="复制 API Base URL"
                      className="ml-auto flex shrink-0 items-center gap-1 rounded-md border border-border/60 bg-secondary/50 px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground transition-all duration-200 hover:border-primary/40 hover:text-primary active:scale-95"
                    >
                      {copied ? (
                        <>
                          <Check className="size-3 text-primary" />
                          <span className="text-primary">已复制</span>
                        </>
                      ) : (
                        <>
                          <Copy className="size-3" />
                          <span>复制</span>
                        </>
                      )}
                    </button>
                  ) : null}
                </dd>
              </div>
            ))}
          </dl>

          <div className="mt-8 flex flex-wrap items-center gap-3">
            <a
              href={SITE.home}
              target="_blank"
              rel="noreferrer"
              className="group inline-flex items-center gap-2 rounded-lg bg-primary px-5 py-2.5 text-sm font-medium text-primary-foreground transition-all duration-200 hover:shadow-[0_0_28px_-6px_oklch(0.82_0.14_172_/_0.6)] active:scale-[0.98]"
            >
              立即访问 SmartQ
              <ArrowRight className="size-4 transition-transform duration-200 group-hover:translate-x-0.5" />
            </a>
            <button
              type="button"
              onClick={() => scrollToSection('api-key')}
              className="inline-flex items-center gap-2 rounded-lg border border-border bg-card/60 px-5 py-2.5 text-sm font-medium text-foreground backdrop-blur transition-all duration-200 hover:border-primary/40 hover:text-primary active:scale-[0.98]"
            >
              <Key className="size-4" />
              查看使用密钥教程
            </button>
          </div>
        </div>

        {/* Right: terminal window */}
        <div className="lg:animate-float">
          <div className="overflow-hidden rounded-2xl border border-border bg-[oklch(0.12_0.006_230)] border-glow">
            <div className="flex items-center gap-2 border-b border-border/70 bg-card/60 px-4 py-3">
              <span className="size-3 rounded-full bg-[oklch(0.6_0.2_22)]/70" />
              <span className="size-3 rounded-full bg-amber/70" />
              <span className="size-3 rounded-full bg-primary/70" />
              <span className="ml-2 flex items-center gap-1.5 font-mono text-xs text-muted-foreground">
                <Terminal className="size-3.5" />
                smartq — zsh
              </span>
            </div>
            <div className="scanlines p-5 font-mono text-[13px] leading-relaxed md:text-sm">
              <p className="text-foreground/90">
                <span className="text-primary">$</span> smartq connect
              </p>
              <p className="mt-2 text-muted-foreground">
                provider:{' '}
                <span className="text-foreground">Claude · Codex · Kimi · GLM</span>
              </p>
              <p className="text-muted-foreground">
                status: <span className="text-primary text-glow">ready</span>
              </p>
              <p className="text-muted-foreground">
                base_url:{' '}
                <span className="text-foreground">{SITE.baseUrl}</span>
              </p>
              <p className="mt-3 text-foreground/90">
                <span className="text-primary">$</span>{' '}
                <span className="animate-caret border-r-2 border-primary pr-0.5">
                  {'\u00A0'}
                </span>
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
