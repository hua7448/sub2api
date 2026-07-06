'use client'

import { useState } from 'react'
import { Check, Copy } from 'lucide-react'
import { cn } from '@/lib/utils'
import { copyText } from '@/lib/copy'

type CodeBlockProps = {
  code: string
  filename?: string
  language?: string
  className?: string
}

export function CodeBlock({ code, filename, className }: CodeBlockProps) {
  const [copied, setCopied] = useState(false)

  async function handleCopy() {
    const ok = await copyText(code)
    if (ok) {
      setCopied(true)
      setTimeout(() => setCopied(false), 1600)
    }
  }

  return (
    <div
      className={cn(
        'group relative overflow-hidden rounded-xl border border-border bg-[oklch(0.13_0.006_230)] transition-all duration-300 hover:border-primary/30',
        className,
      )}
    >
      <div className="flex items-center justify-between border-b border-border/70 bg-card/60 px-4 py-2.5">
        <div className="flex items-center gap-2">
          <span className="size-2.5 rounded-full bg-[oklch(0.6_0.2_22)]/70" />
          <span className="size-2.5 rounded-full bg-amber/70" />
          <span className="size-2.5 rounded-full bg-primary/70" />
          {filename ? (
            <span className="ml-2 font-mono text-xs text-muted-foreground">
              {filename}
            </span>
          ) : null}
        </div>
        <button
          type="button"
          onClick={handleCopy}
          aria-label="复制代码"
          className="flex items-center gap-1.5 rounded-md border border-border/60 bg-secondary/60 px-2.5 py-1 font-mono text-[11px] text-muted-foreground transition-all duration-200 hover:border-primary/40 hover:text-primary active:scale-95"
        >
          {copied ? (
            <>
              <Check className="size-3.5 text-primary" />
              <span className="text-primary">已复制</span>
            </>
          ) : (
            <>
              <Copy className="size-3.5" />
              <span>复制</span>
            </>
          )}
        </button>
      </div>
      <pre className="overflow-x-auto p-4 text-[13px] leading-relaxed">
        <code className="font-mono text-foreground/90">{code}</code>
      </pre>
    </div>
  )
}
