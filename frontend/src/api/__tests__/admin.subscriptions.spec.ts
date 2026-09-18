import { beforeEach, describe, expect, it, vi } from 'vitest'

const { post } = vi.hoisted(() => ({
  post: vi.fn()
}))

vi.mock('@/api/client', () => ({
  apiClient: {
    post
  }
}))

import { resetTodayUsage, type ResetTodayUsageResult } from '@/api/admin/subscriptions'

describe('admin subscriptions API', () => {
  beforeEach(() => {
    post.mockReset()
  })

  it('sends explicit confirmation and the idempotency key when resetting today usage', async () => {
    const response: ResetTodayUsageResult = {
      cutoff_at: '2026-08-01T12:50:00Z',
      captured_at: '2026-08-01T12:50:01Z',
      subscriptions: 637,
      today_window_subscriptions: 600,
      reset_subscriptions: 590,
      reset_and_deducted_usd: 7290.01929621,
      stale_daily_cleared_usd: 12.5,
      cache_invalidations: 1274
    }
    post.mockResolvedValue({ data: response })

    const result = await resetTodayUsage('reset-today-key')

    expect(post).toHaveBeenCalledWith(
      '/admin/subscriptions/reset-today-usage',
      { confirm: true },
      {
        headers: {
          'Idempotency-Key': 'reset-today-key'
        }
      }
    )
    expect(result).toEqual(response)
  })
})
