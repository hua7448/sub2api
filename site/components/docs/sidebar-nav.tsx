'use client'

import { useEffect, useRef, useState } from 'react'
import { cn } from '@/lib/utils'
import { NAV_ITEMS } from '@/lib/docs-data'
import { scrollToSection, SCROLL_OFFSET } from '@/lib/scroll'

export function SidebarNav() {
  const [active, setActive] = useState(NAV_ITEMS[0].id)
  // While a click-triggered smooth scroll is in flight, ignore scroll-spy
  // updates so the highlight doesn't jump to sections we pass through.
  const lockRef = useRef(false)
  const lockTimer = useRef<ReturnType<typeof setTimeout> | null>(null)

  useEffect(() => {
    // Pick the last section whose top has crossed the reference line near the
    // top of the viewport. This is deterministic even for short sections.
    function computeActive() {
      if (lockRef.current) return
      const line = SCROLL_OFFSET + 80
      let current = NAV_ITEMS[0].id
      for (const item of NAV_ITEMS) {
        const el = document.getElementById(item.id)
        if (!el) continue
        if (el.getBoundingClientRect().top - line <= 0) {
          current = item.id
        }
      }
      // Snap to the last item when scrolled to the very bottom.
      if (window.innerHeight + window.scrollY >= document.body.scrollHeight - 2) {
        current = NAV_ITEMS[NAV_ITEMS.length - 1].id
      }
      setActive(current)
    }

    computeActive()
    window.addEventListener('scroll', computeActive, { passive: true })
    window.addEventListener('resize', computeActive)
    return () => {
      window.removeEventListener('scroll', computeActive)
      window.removeEventListener('resize', computeActive)
    }
  }, [])

  function handleClick(e: React.MouseEvent, id: string) {
    e.preventDefault()
    setActive(id)
    lockRef.current = true
    if (lockTimer.current) clearTimeout(lockTimer.current)
    // Release the lock after the smooth scroll has settled.
    lockTimer.current = setTimeout(() => {
      lockRef.current = false
    }, 800)
    scrollToSection(id)
  }

  return (
    <nav aria-label="文档导航" className="flex flex-col gap-1">
      <p className="mb-3 px-3 font-mono text-[11px] uppercase tracking-[0.2em] text-muted-foreground">
        文档目录
      </p>
      {NAV_ITEMS.map((item) => {
        const Icon = item.icon
        const isActive = active === item.id
        return (
          <a
            key={item.id}
            href={`#${item.id}`}
            onClick={(e) => handleClick(e, item.id)}
            className={cn(
              'group flex items-center gap-2.5 rounded-lg px-3 py-2 text-sm transition-all duration-200',
              isActive
                ? 'bg-primary/10 text-primary'
                : 'text-muted-foreground hover:bg-secondary/50 hover:text-foreground',
            )}
          >
            <Icon
              className={cn(
                'size-4 shrink-0 transition-transform duration-200 group-hover:scale-110',
                isActive ? 'text-primary' : 'text-muted-foreground/70',
              )}
            />
            <span className="truncate">{item.label}</span>
            {isActive ? (
              <span className="ml-auto size-1.5 rounded-full bg-primary" />
            ) : null}
          </a>
        )
      })}
    </nav>
  )
}
