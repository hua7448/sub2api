'use client'

import { useState, type ReactNode } from 'react'
import { cn } from '@/lib/utils'

export type Tab = {
  value: string
  label: string
  content: ReactNode
}

export function SimpleTabs({ tabs }: { tabs: Tab[] }) {
  const [active, setActive] = useState(tabs[0]?.value)

  return (
    <div>
      <div
        role="tablist"
        className="inline-flex flex-wrap gap-1 rounded-lg border border-border bg-card/60 p-1"
      >
        {tabs.map((tab) => {
          const isActive = tab.value === active
          return (
            <button
              key={tab.value}
              role="tab"
              aria-selected={isActive}
              type="button"
              onClick={() => setActive(tab.value)}
              className={cn(
                'relative rounded-md px-3.5 py-1.5 font-mono text-xs transition-all duration-200',
                isActive
                  ? 'bg-primary/15 text-primary'
                  : 'text-muted-foreground hover:text-foreground',
              )}
            >
              {tab.label}
            </button>
          )
        })}
      </div>
      <div className="mt-4">
        {tabs.map((tab) =>
          tab.value === active ? (
            <div
              key={tab.value}
              role="tabpanel"
              className="duration-300 animate-in fade-in-50 slide-in-from-bottom-1"
            >
              {tab.content}
            </div>
          ) : null,
        )}
      </div>
    </div>
  )
}
