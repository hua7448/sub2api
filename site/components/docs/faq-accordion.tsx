'use client'

import { useState } from 'react'
import { Plus } from 'lucide-react'
import { cn } from '@/lib/utils'
import type { FaqItem } from '@/lib/docs-data'

export function FaqAccordion({ items }: { items: FaqItem[] }) {
  const [open, setOpen] = useState<number | null>(0)

  return (
    <div className="divide-y divide-border overflow-hidden rounded-xl border border-border bg-card/40">
      {items.map((item, index) => {
        const isOpen = open === index
        return (
          <div key={item.q} className="transition-colors duration-200">
            <button
              type="button"
              onClick={() => setOpen(isOpen ? null : index)}
              aria-expanded={isOpen}
              className="flex w-full items-center justify-between gap-4 px-5 py-4 text-left transition-colors hover:bg-secondary/40"
            >
              <span className="flex items-start gap-3 text-sm font-medium text-foreground">
                <span className="mt-0.5 font-mono text-xs text-primary">
                  {String(index + 1).padStart(2, '0')}
                </span>
                <span className="text-pretty">{item.q}</span>
              </span>
              <Plus
                className={cn(
                  'size-4 shrink-0 text-muted-foreground transition-transform duration-300 ease-out',
                  isOpen && 'rotate-45 text-primary',
                )}
              />
            </button>
            <div
              className={cn(
                'grid transition-all duration-300 ease-out',
                isOpen
                  ? 'grid-rows-[1fr] opacity-100'
                  : 'grid-rows-[0fr] opacity-0',
              )}
            >
              <div className="overflow-hidden">
                <p className="px-5 pb-5 pl-12 text-sm leading-relaxed text-muted-foreground">
                  {item.a}
                </p>
              </div>
            </div>
          </div>
        )
      })}
    </div>
  )
}
