import { SidebarNav } from '@/components/docs/sidebar-nav'
import { MobileNav } from '@/components/docs/mobile-nav'
import { Hero } from '@/components/docs/hero'
import { QuickStart } from '@/components/docs/quick-start'
import { Pricing } from '@/components/docs/pricing'
import { Billing } from '@/components/docs/billing'
import { ApiKeyGuide } from '@/components/docs/api-key-guide'
import { CodexConfig } from '@/components/docs/codex-config'
import { ClaudeConfig } from '@/components/docs/claude-config'
import { InviteReward } from '@/components/docs/invite-reward'
import { FaqSection } from '@/components/docs/faq-section'
import { Contact } from '@/components/docs/contact'
import { Terminal } from 'lucide-react'

import { KimiConfig } from '@/components/docs/kimi-config'

export default function Page() {
  return (
    <div className="relative min-h-screen">
      {/* background grid */}
      <div
        aria-hidden
        className="bg-grid bg-grid-fade pointer-events-none fixed inset-0 -z-10"
      />

      <MobileNav />

      <div className="mx-auto flex w-full max-w-7xl gap-10 px-4 sm:px-6 lg:px-8">
        {/* Sidebar */}
        <aside className="sticky top-0 hidden h-screen w-64 shrink-0 flex-col overflow-y-auto border-r border-border/60 py-8 pr-4 lg:flex">
          <a href="#overview" className="mb-8 flex items-center gap-2.5 px-3">
            <span className="flex size-8 items-center justify-center rounded-lg border border-primary/30 bg-primary/10">
              <Terminal className="size-4 text-primary" />
            </span>
            <span className="font-mono text-sm font-semibold tracking-tight">
              SmartQ<span className="text-primary"> docs</span>
            </span>
          </a>
          <SidebarNav />
          <div className="mt-auto px-3 pt-8">
            <p className="font-mono text-[10px] leading-relaxed text-muted-foreground/60">
              Claude · Codex · Kimi · GLM
            </p>
          </div>
        </aside>

        {/* Main */}
        <main className="min-w-0 flex-1 pb-24">
          <Hero />
          <QuickStart />
          <Pricing />
          <Billing />
          <ApiKeyGuide />
          <CodexConfig />
          <ClaudeConfig />
          <KimiConfig />
          <InviteReward />
          <FaqSection />
          <Contact />
        </main>
      </div>
    </div>
  )
}
