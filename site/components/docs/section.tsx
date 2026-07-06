import type { ReactNode } from 'react'
import { cn } from '@/lib/utils'

type SectionProps = {
  id: string
  eyebrow: string
  title: string
  description?: string
  children: ReactNode
  className?: string
}

export function Section({
  id,
  eyebrow,
  title,
  description,
  children,
  className,
}: SectionProps) {
  return (
    <section
      id={id}
      className={cn(
        'scroll-mt-20 border-t border-border/60 pb-14 pt-12 md:pb-20 md:pt-14',
        className,
      )}
    >
      <div className="mb-8">
        <span className="inline-flex items-center gap-2 font-mono text-xs uppercase tracking-[0.2em] text-primary">
          <span className="size-1.5 rounded-full bg-primary" />
          {eyebrow}
        </span>
        <h2 className="mt-3 text-balance text-2xl font-semibold tracking-tight text-foreground md:text-3xl">
          {title}
        </h2>
        {description ? (
          <p className="mt-3 max-w-2xl text-pretty leading-relaxed text-muted-foreground">
            {description}
          </p>
        ) : null}
      </div>
      {children}
    </section>
  )
}
