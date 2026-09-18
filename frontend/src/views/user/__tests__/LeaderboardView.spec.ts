import { readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

import { describe, expect, it } from 'vitest'

const testDir = dirname(fileURLToPath(import.meta.url))
const viewSource = readFileSync(resolve(testDir, '../LeaderboardView.vue'), 'utf8')
const usageApiSource = readFileSync(resolve(testDir, '../../../api/usage.ts'), 'utf8')
const routerSource = readFileSync(resolve(testDir, '../../../router/index.ts'), 'utf8')
const sidebarSource = readFileSync(resolve(testDir, '../../../components/layout/AppSidebar.vue'), 'utf8')
const enCommonSource = readFileSync(resolve(testDir, '../../../i18n/locales/en/common.ts'), 'utf8')
const zhCommonSource = readFileSync(resolve(testDir, '../../../i18n/locales/zh/common.ts'), 'utf8')
const enMiscSource = readFileSync(resolve(testDir, '../../../i18n/locales/en/misc.ts'), 'utf8')
const zhMiscSource = readFileSync(resolve(testDir, '../../../i18n/locales/zh/misc.ts'), 'utf8')

describe('Leaderboard frontend integration', () => {
  it('keeps the final four-board, tier and time-period experience', () => {
    for (const type of ['cost', 'recharge', 'tokens', 'active_days']) {
      expect(viewSource).toContain(`{ type: '${type}'`)
    }
    expect(viewSource).not.toContain("{ type: 'requests', emoji:")
    for (const period of ['last24h', 'today', 'yesterday', 'last7d', 'month', 'last_month']) {
      expect(viewSource).toContain(`value: '${period}'`)
    }
    expect(viewSource).toContain('const TIERS: TierInfo[]')
    expect(viewSource).toContain('const tierProgress = computed')
    expect(viewSource).toContain('getTop3GapText')
  })

  it('uses authenticated/public endpoints and the ten-minute frontend cache', () => {
    expect(usageApiSource).toContain("'/usage/leaderboard'")
    expect(usageApiSource).toContain("'/public/leaderboard'")
    expect(viewSource).toContain('isLoggedIn.value ? usageAPI.getLeaderboard : usageAPI.getLeaderboardPublic')
    expect(viewSource).toContain('const FRONTEND_CACHE_TTL = 10 * 60 * 1000')
    expect(viewSource).toContain('const frontendCache = new Map')
    expect(viewSource).toContain('padWithSimulatedUsers')
  })

  it('supports public embedding, token adoption, CTA and admin inspection', () => {
    expect(viewSource).toContain(":is=\"isEmbedded ? 'div' : AppLayout\"")
    expect(viewSource).toContain("pageParams.get('ui_mode') === 'embedded'")
    expect(viewSource).toContain('await authStore.setToken(urlToken)')
    expect(viewSource).toContain("cleanUrl.searchParams.delete('token')")
    expect(viewSource).toContain("cleanUrl.searchParams.delete('user_id')")
    expect(viewSource).toContain('target="_top"')
    expect(viewSource).toContain("pageParams.get('src_host')")
    expect(viewSource).toContain('adminAPI.users.getById(userId)')
    expect(viewSource).toContain(':hide-actions="true"')
  })

  it('registers a public route and preserves existing local navigation constraints', () => {
    expect(routerSource).toMatch(/path: '\/leaderboard'[\s\S]*?requiresAuth: false/)
    expect(sidebarSource).toContain("{ path: '/leaderboard', label: t('nav.leaderboard')")
    expect(sidebarSource).toContain("{ path: '/batch-image', label: t('nav.batchImage')")
    expect(sidebarSource).toContain('const flagAdminPayment = () => false')
  })

  it('provides split English and Chinese leaderboard locale keys', () => {
    expect(enCommonSource).toContain("leaderboard: 'Leaderboard'")
    expect(zhCommonSource).toContain("leaderboard: '排行榜'")
    expect(enMiscSource).toContain("title: 'Leaderboard'")
    expect(enMiscSource).toContain("typeActiveDays: 'Active Days'")
    expect(zhMiscSource).toContain("title: '排行榜'")
    expect(zhMiscSource).toContain("typeActiveDays: '劳模榜'")
  })
})
