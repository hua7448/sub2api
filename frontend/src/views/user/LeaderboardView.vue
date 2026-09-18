<template>
  <!-- Embedded mode intentionally omits AppLayout so an iframe does not get a second header/sidebar. -->
  <component :is="isEmbedded ? 'div' : AppLayout">
    <div class="mx-auto space-y-5 px-4 py-6">
      <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <div class="flex items-center gap-3">
          <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-amber-400 to-orange-500 shadow-lg">
            <svg class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M16.5 18.75h-9m9 0a3 3 0 013 3h-15a3 3 0 013-3m9 0v-3.375c0-.621-.503-1.125-1.125-1.125h-.871M7.5 18.75v-3.375c0-.621.504-1.125 1.125-1.125h.872m5.007 0H9.497m5.007 0a7.454 7.454 0 01-.982-3.172M9.497 14.25a7.454 7.454 0 00.981-3.172M5.25 4.236c-.982.143-1.954.317-2.916.52A6.003 6.003 0 007.73 9.728M5.25 4.236V4.5c0 2.108.966 3.99 2.48 5.228M5.25 4.236V2.721C7.456 2.41 9.71 2.25 12 2.25c2.291 0 4.545.16 6.75.47v1.516M18.75 4.236c.982.143 1.954.317 2.916.52A6.003 6.003 0 0016.27 9.728M18.75 4.236V4.5c0 2.108-.966 3.99-2.48 5.228m0 0a6.003 6.003 0 01-3.77 1.522m3.77-1.522a6.003 6.003 0 00-.34-6.478M9.73 9.728a6.003 6.003 0 003.77 1.522m-3.77-1.522a6.003 6.003 0 01.34-6.478m6.66 0a6.003 6.003 0 00-7 0" />
            </svg>
          </div>
          <div>
            <h1 class="text-2xl font-bold text-txt-primary">{{ t('leaderboard.title') }}</h1>
            <p class="text-sm text-txt-secondary">{{ t('leaderboard.description') }}</p>
          </div>
        </div>

        <div class="flex max-w-full gap-1 overflow-x-auto rounded-xl bg-cyber-elevated p-1.5">
          <button
            v-for="period in periods"
            :key="period.value"
            type="button"
            class="whitespace-nowrap rounded-lg px-4 py-2 text-sm font-medium transition-all duration-200"
            :class="currentPeriod === period.value
              ? 'bg-primary-500 text-white shadow-sm'
              : 'text-txt-secondary hover:text-txt-primary'"
            @click="selectPeriod(period.value)"
          >
            {{ period.label }}
          </button>
        </div>
      </div>

      <div v-if="isLoggedIn && bestRankData" class="card p-6">
        <div class="flex flex-col items-center gap-6 md:flex-row">
          <div class="flex flex-col items-center gap-1.5">
            <span class="text-6xl">{{ myTier.icon }}</span>
            <span class="text-base font-bold" :class="myTier.textClass">
              {{ t('leaderboard.tier.' + myTier.key) }}
            </span>
          </div>

          <div class="min-w-0 flex-1">
            <div v-if="tierProgress" class="mb-4">
              <div class="mb-2 flex items-center justify-between gap-3">
                <span class="text-sm font-semibold text-txt-primary">{{ t('leaderboard.tier.' + myTier.key) }}</span>
                <span class="text-sm font-medium text-txt-secondary">{{ tierProgress.gapText }}</span>
              </div>
              <div class="h-3 overflow-hidden rounded-full bg-cyber-elevated">
                <div
                  class="h-full rounded-full transition-all duration-700 ease-out"
                  :class="myTier.barClass"
                  :style="{ width: tierProgress.percent + '%' }"
                ></div>
              </div>
            </div>

            <div class="flex flex-wrap gap-3">
              <div
                v-for="board in visibleBoards"
                :key="board.type"
                class="flex items-center gap-2 rounded-xl bg-cyber-bg px-4 py-2.5"
              >
                <span class="text-base">{{ board.emoji }}</span>
                <span class="text-sm text-txt-secondary">{{ board.label }}</span>
                <span class="text-lg font-bold text-txt-primary">#{{ getMyRankForType(board.type) || '—' }}</span>
                <span class="text-sm text-txt-dim">{{ getMyValueForType(board.type) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="loadingAll && !hasData" class="flex h-64 items-center justify-center">
        <div class="flex flex-col items-center gap-3">
          <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
          <span class="text-sm text-txt-dim">{{ t('leaderboard.loading') }}</span>
        </div>
      </div>

      <div v-if="hasData || !loadingAll" class="relative">
        <div v-if="loadingAll && hasData" class="absolute -top-2 right-0 z-10">
          <div class="h-5 w-5 animate-spin rounded-full border-2 border-primary-400 border-t-transparent"></div>
        </div>

        <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 xl:grid-cols-4" :class="{ 'opacity-60 transition-opacity': loadingAll && hasData }">
          <div
            v-for="(board, boardIndex) in visibleBoards"
            :key="board.type"
            class="card overflow-hidden"
            :style="{ animationDelay: boardIndex * 80 + 'ms' }"
            style="animation: fadeInUp 0.4s ease-out both"
          >
            <div class="flex items-center gap-2.5 border-b border-cyber-border px-5 py-4">
              <span class="text-xl">{{ board.emoji }}</span>
              <span class="text-base font-bold text-txt-primary">{{ board.label }}</span>
            </div>

            <div>
              <div v-if="getBoardItems(board.type).length === 0" class="px-5 py-12 text-center text-base text-txt-dim">
                {{ t('leaderboard.noData') }}
              </div>

              <div
                v-for="entry in getBoardItems(board.type)"
                :key="`${board.type}-${entry.user_id}-${entry.rank}`"
                class="flex items-center gap-3 border-l-[3px] px-5 py-3.5 transition-colors hover:bg-cyber-elevated/50"
                :class="entry.rank <= 3 ? rankBorderClass(entry.rank) : 'border-transparent'"
              >
                <div class="flex w-9 flex-shrink-0 items-center justify-center">
                  <span v-if="entry.rank <= 3" class="text-xl">{{ rankEmoji(entry.rank) }}</span>
                  <span v-else class="text-xl font-bold text-txt-dim">{{ entry.rank }}</span>
                </div>

                <div class="min-w-0 flex-1">
                  <div class="flex items-center gap-1.5">
                    <span class="text-sm">{{ getTierIcon(entry.rank, getBoardTotal(board.type)) }}</span>
                    <button
                      v-if="isAdmin && entry.user_id > 0"
                      type="button"
                      class="truncate text-sm font-medium text-primary-500 hover:text-primary-600 hover:underline"
                      @click="showUserInfo(entry.user_id)"
                    >
                      {{ entry.masked_email }}
                    </button>
                    <span v-else class="truncate text-sm font-medium text-txt-primary">{{ entry.masked_email }}</span>
                  </div>
                  <div v-if="entry.title" class="mt-0.5 text-xs text-txt-dim">{{ entry.title }}</div>
                </div>

                <div class="flex-shrink-0 text-right">
                  <span class="text-base font-bold text-txt-primary">{{ formatValue(entry.value, board.type) }}</span>
                </div>
              </div>
            </div>

            <div
              v-if="isLoggedIn && getMyRankInfo(board.type)"
              class="border-t border-cyber-border bg-cyber-bg/50 px-5 py-4"
            >
              <div class="flex items-center gap-2">
                <span class="text-sm text-txt-secondary">📍 {{ t('leaderboard.myRank') }}</span>
                <span class="text-lg font-bold text-primary-500">#{{ getMyRankInfo(board.type)?.rank }}</span>
                <span class="text-sm font-medium text-txt-secondary">{{ formatValue(getMyRankInfo(board.type)?.value || 0, board.type) }}</span>
              </div>
              <div v-if="getMyRankInfo(board.type)?.next_rank_gap" class="mt-2 text-sm font-medium text-amber-500">
                💡 {{ getGapTextForBoard(board.type) }}
              </div>
              <div v-if="getTop3GapText(board.type)" class="mt-1.5 text-sm text-amber-500/80">
                🏆 {{ getTop3GapText(board.type) }}
              </div>
              <div
                v-if="getMyRankInfo(board.type)?.rank === 1 && getBoardItems(board.type).length > 1"
                class="mt-1.5 text-sm font-medium text-emerald-500"
              >
                👑 {{ t('leaderboard.leading', {
                  amount: formatValue((getMyRankInfo(board.type)?.value || 0) - getBoardItems(board.type)[1].value, board.type)
                }) }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="isLoggedIn && bestGapText && !loadingAll" class="card sticky bottom-4 z-10 px-6 py-5">
        <div class="flex items-center gap-3">
          <span class="text-xl">💡</span>
          <span class="text-base font-semibold text-amber-500">{{ bestGapText }}</span>
        </div>
      </div>

      <div v-if="!isLoggedIn && !loadingAll" class="card sticky bottom-4 z-10 px-6 py-5">
        <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div class="flex items-center gap-3">
            <span class="text-2xl">🔥</span>
            <div>
              <p class="text-base font-bold text-txt-primary">{{ t('leaderboard.ctaTitle') }}</p>
              <p class="text-sm text-txt-secondary">{{ t('leaderboard.ctaDesc') }}</p>
            </div>
          </div>
          <a
            :href="mainSiteLoginUrl"
            target="_top"
            class="flex-shrink-0 rounded-xl bg-gradient-to-r from-amber-500 to-orange-500 px-6 py-3 text-center text-sm font-bold text-white shadow-lg transition-all hover:brightness-110 hover:shadow-xl"
          >
            🏆 {{ t('leaderboard.ctaButton') }}
          </a>
        </div>
      </div>
    </div>

    <UserBalanceHistoryModal
      v-if="isAdmin"
      :show="showBalanceModal"
      :user="selectedUser"
      :hide-actions="true"
      @close="closeUserInfo"
    />
  </component>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { AppLayout } from '@/components/layout'
import UserBalanceHistoryModal from '@/components/admin/user/UserBalanceHistoryModal.vue'
import { adminAPI } from '@/api/admin'
import { usageAPI } from '@/api/usage'
import type { LeaderboardPeriod, LeaderboardResponse, LeaderboardType } from '@/api/usage'
import { useAuthStore } from '@/stores/auth'
import type { AdminUser } from '@/types'

interface TierInfo {
  key: 'champion' | 'diamond' | 'platinum' | 'gold' | 'silver' | 'bronze'
  icon: string
  textClass: string
  barClass: string
}

interface BoardConfig {
  type: LeaderboardType
  emoji: string
  label: string
}

const { t } = useI18n()
const authStore = useAuthStore()
const isLoggedIn = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)

const selectedUser = ref<AdminUser | null>(null)
const showBalanceModal = ref(false)

async function showUserInfo(userId: number) {
  if (!isAdmin.value || userId <= 0) return
  try {
    selectedUser.value = await adminAPI.users.getById(userId)
    showBalanceModal.value = true
  } catch (error) {
    console.error('Failed to load leaderboard user info:', error)
  }
}

function closeUserInfo() {
  showBalanceModal.value = false
  selectedUser.value = null
}

const pageParams = new URLSearchParams(window.location.search)
const isEmbedded = pageParams.get('ui_mode') === 'embedded'

const mainSiteLoginUrl = computed(() => {
  const sourceHost = pageParams.get('src_host')
  if (!sourceHost) return '/login'
  try {
    const url = new URL(sourceHost, window.location.origin)
    if (url.protocol !== 'http:' && url.protocol !== 'https:') return '/login'
    return new URL('/login', url).toString()
  } catch {
    return '/login'
  }
})

const currentPeriod = ref<LeaderboardPeriod>('last24h')
const boardDataMap = ref<Partial<Record<LeaderboardType, LeaderboardResponse>>>({})
const loadingAll = ref(false)
const hasData = computed(() => Object.keys(boardDataMap.value).length > 0)

const allBoards: BoardConfig[] = [
  { type: 'cost', emoji: '💰', label: '' },
  { type: 'recharge', emoji: '💎', label: '' },
  { type: 'tokens', emoji: '🔥', label: '' },
  { type: 'active_days', emoji: '💪', label: '' }
]

const visibleBoards = computed<BoardConfig[]>(() =>
  allBoards.map((board) => ({
    ...board,
    label: t(`leaderboard.type${capitalize(board.type)}`)
  }))
)

const periods = computed<Array<{ value: LeaderboardPeriod; label: string }>>(() => [
  { value: 'last24h', label: t('leaderboard.periodLast24h') },
  { value: 'today', label: t('leaderboard.periodToday') },
  { value: 'yesterday', label: t('leaderboard.periodYesterday') },
  { value: 'last7d', label: t('leaderboard.periodLast7d') },
  { value: 'month', label: t('leaderboard.periodMonth') },
  { value: 'last_month', label: t('leaderboard.periodLastMonth') }
])

const TIERS: TierInfo[] = [
  { key: 'champion', icon: '👑', textClass: 'text-red-500', barClass: 'bg-red-500' },
  { key: 'diamond', icon: '💎', textClass: 'text-cyan-500', barClass: 'bg-cyan-500' },
  { key: 'platinum', icon: '💠', textClass: 'text-txt-secondary', barClass: 'bg-gray-400' },
  { key: 'gold', icon: '⭐', textClass: 'text-amber-500', barClass: 'bg-amber-500' },
  { key: 'silver', icon: '⚪', textClass: 'text-txt-dim', barClass: 'bg-gray-300' },
  { key: 'bronze', icon: '🔰', textClass: 'text-orange-400', barClass: 'bg-orange-400' }
]

function getTier(rank: number, total: number): TierInfo {
  if (total <= 0 || rank <= 0) return TIERS[5]
  if (rank === 1) return TIERS[0]
  if (rank <= 3) return TIERS[1]
  if (rank <= Math.ceil(total * 0.1)) return TIERS[2]
  if (rank <= Math.ceil(total * 0.3)) return TIERS[3]
  if (rank <= Math.ceil(total * 0.6)) return TIERS[4]
  return TIERS[5]
}

function getTierIcon(rank: number, total: number): string {
  return getTier(rank, total).icon
}

const bestRankData = computed(() => {
  let bestRank = Number.POSITIVE_INFINITY
  let bestTotal = 0
  for (const board of allBoards) {
    const data = boardDataMap.value[board.type]
    if (!data?.my_rank || data.my_rank.rank <= 0) continue
    if (data.my_rank.rank < bestRank) {
      bestRank = data.my_rank.rank
      bestTotal = data.items.length
    }
  }
  return Number.isFinite(bestRank) ? { rank: bestRank, total: bestTotal } : null
})

const myTier = computed(() => {
  if (!bestRankData.value) return TIERS[5]
  return getTier(bestRankData.value.rank, bestRankData.value.total)
})

const tierProgress = computed(() => {
  const best = bestRankData.value
  if (!best || best.total <= 1) return null

  const boundaries = [
    { tier: 'champion', maxRank: 1 },
    { tier: 'diamond', maxRank: 3 },
    { tier: 'platinum', maxRank: Math.ceil(best.total * 0.1) },
    { tier: 'gold', maxRank: Math.ceil(best.total * 0.3) },
    { tier: 'silver', maxRank: Math.ceil(best.total * 0.6) },
    { tier: 'bronze', maxRank: best.total }
  ]

  let currentIndex = boundaries.length - 1
  for (let index = 0; index < boundaries.length; index += 1) {
    if (best.rank <= boundaries[index].maxRank) {
      currentIndex = index
      break
    }
  }
  if (currentIndex === 0) return { percent: 100, gapText: t('leaderboard.tierMax') }

  const nextBoundary = boundaries[currentIndex - 1].maxRank
  const currentBoundary = boundaries[currentIndex].maxRank
  const ranksInTier = currentBoundary - nextBoundary
  const ranksAboveNext = best.rank - nextBoundary
  const percent = ranksInTier > 0
    ? Math.round(((ranksInTier - ranksAboveNext) / ranksInTier) * 100)
    : 0
  const nextTier = t(`leaderboard.tier.${boundaries[currentIndex - 1].tier}`)
  return {
    percent: Math.max(5, Math.min(percent, 100)),
    gapText: t('leaderboard.tierGap', { count: ranksAboveNext, tier: nextTier })
  }
})

function getMyRankForType(type: LeaderboardType): number | null {
  const rank = boardDataMap.value[type]?.my_rank?.rank
  return rank && rank > 0 ? rank : null
}

function getMyValueForType(type: LeaderboardType): string {
  const rank = boardDataMap.value[type]?.my_rank
  return rank && rank.rank > 0 ? formatValue(rank.value, type) : ''
}

function getMyRankInfo(type: LeaderboardType) {
  const rank = boardDataMap.value[type]?.my_rank
  return rank && rank.rank > 0 ? rank : null
}

function actionFor(type: LeaderboardType): string {
  const keys: Record<LeaderboardType, string> = {
    cost: 'leaderboard.gapCost',
    recharge: 'leaderboard.gapRecharge',
    tokens: 'leaderboard.gapTokens',
    requests: 'leaderboard.gapRequests',
    active_days: 'leaderboard.gapActiveDays'
  }
  return t(keys[type])
}

function getGapTextForBoard(type: LeaderboardType): string {
  const rank = getMyRankInfo(type)
  if (!rank || rank.rank <= 1 || !rank.next_rank_gap) return ''
  return t('leaderboard.gapNext', {
    action: actionFor(type),
    amount: formatValue(rank.next_rank_gap, type),
    rank: rank.rank - 1,
    email: rank.next_rank_email || ''
  })
}

function getTop3GapText(type: LeaderboardType): string {
  const rank = getMyRankInfo(type)
  const items = boardDataMap.value[type]?.items
  if (!rank || rank.rank <= 3 || !items || items.length < 3) return ''
  const gap = items[2].value - rank.value
  if (gap <= 0) return ''
  return t('leaderboard.gapTop3', {
    action: actionFor(type),
    amount: formatValue(gap, type)
  })
}

function getBoardItems(type: LeaderboardType) {
  return boardDataMap.value[type]?.items.slice(0, 6) || []
}

function getBoardTotal(type: LeaderboardType): number {
  return boardDataMap.value[type]?.items.length || 0
}

const bestGapText = computed(() => {
  for (const board of allBoards) {
    const rank = boardDataMap.value[board.type]?.my_rank
    if (!rank || rank.rank <= 1 || !rank.next_rank_gap) continue
    const text = t('leaderboard.gapInBoard', {
      board: board.emoji + t(`leaderboard.type${capitalize(board.type)}`),
      rank: rank.rank,
      action: actionFor(board.type),
      amount: formatValue(rank.next_rank_gap, board.type)
    })
    return rank.next_rank_email ? `${text} (${rank.next_rank_email})` : text
  }
  return ''
})

function capitalize(value: string): string {
  return value.replace(/(^|_)(\w)/g, (_match, _separator, character: string) => character.toUpperCase())
}

function formatValue(value: number, type?: LeaderboardType | string): string {
  if (type === 'cost' || type === 'recharge') return `$${value.toFixed(2)}`
  if (type === 'active_days') return `${value}${t('leaderboard.days')}`
  if (value >= 1_000_000) return `${(value / 1_000_000).toFixed(1)}M`
  if (value >= 1_000) return `${(value / 1_000).toFixed(1)}K`
  return value.toLocaleString()
}

function rankEmoji(rank: number): string {
  if (rank === 1) return '🥇'
  if (rank === 2) return '🥈'
  if (rank === 3) return '🥉'
  return ''
}

function rankBorderClass(rank: number): string {
  if (rank === 1) return 'border-amber-400'
  if (rank === 2) return 'border-gray-300'
  if (rank === 3) return 'border-orange-300'
  return 'border-transparent'
}

const FRONTEND_CACHE_TTL = 10 * 60 * 1000
const frontendCache = new Map<string, { data: LeaderboardResponse; time: number }>()

async function selectPeriod(period: LeaderboardPeriod) {
  if (period === currentPeriod.value && hasData.value) return
  currentPeriod.value = period
  await fetchAllBoards()
}

async function fetchAllBoards() {
  loadingAll.value = true
  try {
    const now = Date.now()
    const results: LeaderboardResponse[] = []
    const toFetch: Array<{ index: number; board: BoardConfig }> = []

    allBoards.forEach((board, index) => {
      const cacheKey = `${board.type}:${currentPeriod.value}`
      const cached = frontendCache.get(cacheKey)
      if (cached && now - cached.time < FRONTEND_CACHE_TTL) {
        results[index] = cached.data
      } else {
        toFetch.push({ index, board })
      }
    })

    if (toFetch.length > 0) {
      const request = isLoggedIn.value ? usageAPI.getLeaderboard : usageAPI.getLeaderboardPublic
      const fetched = await Promise.all(
        toFetch.map(({ board }) => request(board.type, currentPeriod.value, 20))
      )
      toFetch.forEach(({ index, board }, fetchedIndex) => {
        const data = fetched[fetchedIndex]
        results[index] = data
        frontendCache.set(`${board.type}:${currentPeriod.value}`, { data, time: now })
      })
    }

    const next: Partial<Record<LeaderboardType, LeaderboardResponse>> = {}
    allBoards.forEach((board, index) => {
      const data = results[index]
      if (data) next[board.type] = padWithSimulatedUsers(data, board.type)
    })
    boardDataMap.value = next
  } catch (error) {
    console.error('Failed to fetch leaderboards:', error)
  } finally {
    loadingAll.value = false
  }
}

const SIMULATED_EMAILS = [
  'vi**@gmail.com', 'so**@qq.com', 'al**@outlook.com', 'ja**@163.com',
  'mi**@icloud.com', 'zh**@gmail.com', 'li**@hotmail.com', 'wa**@126.com',
  'ch**@proton.me', 'yu**@yahoo.com', 'ke**@gmail.com', 'to**@live.com'
]

const SIMULATED_TITLES = ['', '', '', '活跃用户', '', '资深玩家', '', '', '老用户', '']

function getSimulatedValueRange(type: LeaderboardType, period: LeaderboardPeriod) {
  const periodMultipliers: Record<LeaderboardPeriod, number> = {
    last24h: 1,
    today: 1,
    yesterday: 1,
    last7d: 4,
    month: 8,
    last30d: 8,
    last_month: 8
  }
  const multiplier = periodMultipliers[period]
  switch (type) {
    case 'cost': return { base: 15 * multiplier, variance: 80 * multiplier }
    case 'recharge': return { base: 25 * multiplier, variance: 120 * multiplier }
    case 'tokens': return { base: 8_000 * multiplier, variance: 50_000 * multiplier }
    case 'active_days': return { base: Math.min(multiplier, 28), variance: Math.min(3 * multiplier, 28) }
    default: return { base: 10 * multiplier, variance: 50 * multiplier }
  }
}

function seededRandom(seed: number): number {
  const value = Math.sin(seed) * 10_000
  return value - Math.floor(value)
}

function padWithSimulatedUsers(data: LeaderboardResponse, type: LeaderboardType): LeaderboardResponse {
  const minimumUsers = 10
  const realCount = data.items.length
  if (realCount >= minimumUsers) return data

  const maximumRealValue = realCount > 0 ? data.items[0].value : 0
  const periodSeeds: Record<LeaderboardPeriod, number> = {
    last24h: 1,
    today: 2,
    yesterday: 3,
    last7d: 4,
    month: 5,
    last30d: 6,
    last_month: 7
  }
  const seedBase = type.charCodeAt(0) * 10_000 + periodSeeds[data.period] * 1_000
  const range = getSimulatedValueRange(type, data.period)
  const simulated: LeaderboardResponse['items'] = []
  const usedEmails = new Set(data.items.map((item) => item.masked_email))

  for (let index = 0; index < minimumUsers - realCount; index += 1) {
    let email = SIMULATED_EMAILS[index % SIMULATED_EMAILS.length]
    if (usedEmails.has(email)) email = SIMULATED_EMAILS[(index + 5) % SIMULATED_EMAILS.length]
    usedEmails.add(email)

    const random = seededRandom(seedBase + index * 7)
    let value = maximumRealValue > range.base * 2
      ? maximumRealValue * (0.3 + random * 1.5)
      : range.base + random * range.variance
    value = Math.max(value, range.base * 0.8)
    if (type === 'active_days') value = Math.round(value)
    if (type === 'cost' || type === 'recharge') value = Math.round(value * 100) / 100

    simulated.push({
      rank: 0,
      user_id: -(index + 1),
      masked_email: email,
      value,
      requests: Math.round(value * (2 + random * 5)),
      tokens: Math.round(value * (100 + random * 500)),
      title: SIMULATED_TITLES[index % SIMULATED_TITLES.length] || undefined
    })
  }

  const items = [...data.items, ...simulated]
    .sort((left, right) => right.value - left.value)
    .map((item, index) => ({ ...item, rank: index + 1 }))

  let myRank = data.my_rank
  if (myRank && myRank.rank > 0) {
    const myEntry = items.find((item) => item.user_id > 0 && item.value === myRank?.value)
    if (myEntry) {
      const nextEntry = items.find((item) => item.rank === myEntry.rank - 1)
      myRank = {
        ...myRank,
        rank: myEntry.rank,
        next_rank_gap: nextEntry ? nextEntry.value - myRank.value : myRank.next_rank_gap,
        next_rank_email: nextEntry ? nextEntry.masked_email : myRank.next_rank_email
      }
    }
  } else if (isLoggedIn.value) {
    const nextEntry = items[items.length - 1]
    myRank = {
      rank: items.length + 1,
      value: 0,
      requests: 0,
      tokens: 0,
      next_rank_gap: nextEntry?.value,
      next_rank_email: nextEntry?.masked_email
    }
  }

  return { ...data, items, my_rank: myRank }
}

onMounted(async () => {
  const urlToken = pageParams.get('token')
  if (urlToken && !isLoggedIn.value) {
    try {
      await authStore.setToken(urlToken)
      const cleanUrl = new URL(window.location.href)
      cleanUrl.searchParams.delete('token')
      cleanUrl.searchParams.delete('user_id')
      window.history.replaceState({}, '', cleanUrl.toString())
    } catch {
      // An invalid embedded token is treated as an anonymous public view.
    }
  }
  await fetchAllBoards()
})
</script>

<style scoped>
@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
