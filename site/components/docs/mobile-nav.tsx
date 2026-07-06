'use client'

import { useEffect, useState } from 'react'
import { Menu, Terminal, X } from 'lucide-react'
import { cn } from '@/lib/utils'
import { NAV_ITEMS } from '@/lib/docs-data'
import { scrollToSection } from '@/lib/scroll'

export function MobileNav() {
  const [open, setOpen] = useState(false)

  useEffect(() => {
    document.body.style.overflow = open ? 'hidden' : ''
    return () => {
      document.body.style.overflow = ''
    }
  }, [open])

  function handleClick(e: React.MouseEvent, id: string) {
    e.preventDefault()
    setOpen(false)
    // wait for the drawer close + body overflow reset, then scroll
    // clear the sticky mobile header (~56px) with a larger offset
    setTimeout(() => scrollToSection(id, 68), 240)
  }

  return (
    <div className="lg:hidden">
      <header className="sticky top-0 z-40 flex items-center justify-between border-b border-border/70 bg-background/80 px-4 py-3 backdrop-blur-xl">
        <div className="flex items-center gap-2">
          <span className="flex size-7 items-center justify-center rounded-md border border-primary/30 bg-primary/10">
            <Terminal className="size-4 text-primary" />
          </span>
          <span className="font-mono text-sm font-semibold tracking-tight">
            SmartQ<span className="text-primary"> docs</span>
          </span>
        </div>
        <button
          type="button"
          onClick={() => setOpen(true)}
          aria-label="打开菜单"
          className="flex size-9 items-center justify-center rounded-lg border border-border bg-card/60 text-foreground active:scale-95"
        >
          <Menu className="size-5" />
        </button>
      </header>

      <div
        className={cn(
          'fixed inset-0 z-50 transition-opacity duration-300',
          open ? 'pointer-events-auto opacity-100' : 'pointer-events-none opacity-0',
        )}
      >
        <div
          className="absolute inset-0 bg-background/70 backdrop-blur-sm"
          onClick={() => setOpen(false)}
        />
        <div
          className={cn(
            'absolute right-0 top-0 h-full w-72 max-w-[85%] border-l border-border bg-card p-4 shadow-2xl transition-transform duration-300 ease-out',
            open ? 'translate-x-0' : 'translate-x-full',
          )}
        >
          <div className="mb-4 flex items-center justify-between">
            <span className="font-mono text-xs uppercase tracking-[0.2em] text-muted-foreground">
              文档目录
            </span>
            <button
              type="button"
              onClick={() => setOpen(false)}
              aria-label="关闭菜单"
              className="flex size-8 items-center justify-center rounded-lg border border-border text-muted-foreground active:scale-95"
            >
              <X className="size-4" />
            </button>
          </div>
          <nav className="flex flex-col gap-1">
            {NAV_ITEMS.map((item) => {
              const Icon = item.icon
              return (
                <a
                  key={item.id}
                  href={`#${item.id}`}
                  onClick={(e) => handleClick(e, item.id)}
                  className="flex items-center gap-2.5 rounded-lg px-3 py-2.5 text-sm text-muted-foreground transition-colors hover:bg-secondary/50 hover:text-foreground"
                >
                  <Icon className="size-4 text-muted-foreground/70" />
                  {item.label}
                </a>
              )
            })}
          </nav>
        </div>
      </div>
    </div>
  )
}
