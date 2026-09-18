import type { Account } from '@/types'

const WINDOW_RESET_TOLERANCE_MS = 2 * 60 * 1000

function parseNumber(value: unknown): number | null {
  if (value == null || value === '') return null
  const parsed = Number(value)
  return Number.isFinite(parsed) ? parsed : null
}

function parseTimestamp(value: unknown): number | null {
  if (typeof value !== 'string' || value.trim() === '') return null
  const parsed = Date.parse(value)
  return Number.isFinite(parsed) ? parsed : null
}

function codexWindowResetAt(extra: Record<string, unknown>, window: '5h' | '7d'): number | null {
  const resetAt = parseTimestamp(extra[`codex_${window}_reset_at`])
  if (resetAt !== null) return resetAt

  const resetAfterSeconds = parseNumber(extra[`codex_${window}_reset_after_seconds`])
  if (resetAfterSeconds === null || resetAfterSeconds <= 0) return null

  const updatedAt = parseTimestamp(extra.codex_usage_updated_at) ?? Date.now()
  return updatedAt + resetAfterSeconds * 1000
}

export function isOpenAICodexWindowResetRateLimit(account: Account, now = Date.now()): boolean {
  if (account.platform !== 'openai' || !account.rate_limit_reset_at) return false

  const rateLimitResetAt = Date.parse(account.rate_limit_reset_at)
  if (!Number.isFinite(rateLimitResetAt) || rateLimitResetAt <= now) return false

  const extra = account.extra as Record<string, unknown> | undefined
  if (!extra) return false

  const used5h = parseNumber(extra.codex_5h_used_percent)
  const used7d = parseNumber(extra.codex_7d_used_percent)
  if (used5h === null || used7d === null || used5h >= 100 || used7d >= 100) return false

  return (['5h', '7d'] as const).some((window) => {
    const resetAt = codexWindowResetAt(extra, window)
    return resetAt !== null && Math.abs(rateLimitResetAt - resetAt) <= WINDOW_RESET_TOLERANCE_MS
  })
}

export function isAccountRuntimeRateLimited(account: Account, now = Date.now()): boolean {
  if (!account.rate_limit_reset_at) return false
  const resetAt = Date.parse(account.rate_limit_reset_at)
  if (!Number.isFinite(resetAt) || resetAt <= now) return false
  return !isOpenAICodexWindowResetRateLimit(account, now)
}
