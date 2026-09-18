import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'

import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import SubscriptionsView from '../SubscriptionsView.vue'

const {
  listSubscriptions,
  getAllGroups,
  resetTodayUsage,
  showSuccess,
  showError
} = vi.hoisted(() => ({
  listSubscriptions: vi.fn(),
  getAllGroups: vi.fn(),
  resetTodayUsage: vi.fn(),
  showSuccess: vi.fn(),
  showError: vi.fn()
}))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    subscriptions: {
      list: listSubscriptions,
      resetTodayUsage
    },
    groups: {
      getAll: getAllGroups
    },
    usage: {
      searchUsers: vi.fn()
    }
  }
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({ showSuccess, showError })
}))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string, params?: Record<string, unknown>) =>
        params ? `${key}:${JSON.stringify(params)}` : key
    })
  }
})

const resetResult = {
  cutoff_at: '2026-08-01T12:50:00Z',
  captured_at: '2026-08-01T12:50:01Z',
  subscriptions: 637,
  today_window_subscriptions: 600,
  reset_subscriptions: 590,
  reset_and_deducted_usd: 7290.01929621,
  stale_daily_cleared_usd: 12.5,
  cache_invalidations: 1274
}

const mountView = () =>
  mount(SubscriptionsView, {
    global: {
      stubs: {
        AppLayout: { template: '<div><slot /></div>' },
        TablePageLayout: {
          template: '<div><slot name="filters" /><slot name="table" /><slot name="pagination" /></div>'
        },
        DataTable: true,
        Pagination: true,
        BaseDialog: true,
        ConfirmDialog: true,
        EmptyState: true,
        Select: true,
        GroupBadge: true,
        GroupOptionItem: true,
        Icon: true,
        RouterLink: true,
        Teleport: true
      }
    }
  })

const activeConfirmDialog = (wrapper: ReturnType<typeof mountView>) =>
  wrapper.findAllComponents(ConfirmDialog).find((dialog) => dialog.props('show'))

describe('admin SubscriptionsView reset today usage', () => {
  beforeEach(() => {
    localStorage.clear()
    listSubscriptions.mockReset()
    getAllGroups.mockReset()
    resetTodayUsage.mockReset()
    showSuccess.mockReset()
    showError.mockReset()

    listSubscriptions.mockResolvedValue({
      items: [],
      total: 0,
      page: 1,
      page_size: 20,
      pages: 0
    })
    getAllGroups.mockResolvedValue([])
    vi.stubGlobal('crypto', {
      randomUUID: vi.fn(() => '11111111-1111-4111-8111-111111111111')
    })
  })

  afterEach(() => {
    vi.useRealTimers()
    vi.unstubAllGlobals()
    vi.restoreAllMocks()
  })

  it('requires confirmation, disables while pending, and reports the reset result', async () => {
    let resolveReset!: (value: typeof resetResult) => void
    resetTodayUsage.mockReturnValue(new Promise((resolve) => {
      resolveReset = resolve
    }))

    const wrapper = mountView()
    await flushPromises()

    const resetButton = wrapper.get('[data-test="reset-today-usage"]')
    await resetButton.trigger('click')

    expect(resetTodayUsage).not.toHaveBeenCalled()
    const dialog = activeConfirmDialog(wrapper)
    expect(dialog).toBeTruthy()

    dialog?.vm.$emit('confirm')
    await flushPromises()

    expect(resetTodayUsage).toHaveBeenCalledTimes(1)
    expect(resetTodayUsage).toHaveBeenCalledWith('11111111-1111-4111-8111-111111111111')
    expect(resetButton.attributes('disabled')).toBeDefined()

    resolveReset(resetResult)
    await flushPromises()

    expect(showSuccess).toHaveBeenCalledWith(
      'admin.subscriptions.resetTodayUsageSuccess:{"subscriptions":637,"amount":"7290.01929621"}'
    )
    expect(listSubscriptions).toHaveBeenCalledTimes(2)
    expect(resetButton.attributes('disabled')).toBeUndefined()
    wrapper.unmount()
  })

  it('reuses the same idempotency key when a failed request is confirmed again', async () => {
    vi.spyOn(console, 'error').mockImplementation(() => undefined)
    resetTodayUsage
      .mockRejectedValueOnce(new Error('network response lost'))
      .mockResolvedValueOnce(resetResult)

    const wrapper = mountView()
    await flushPromises()

    await wrapper.get('[data-test="reset-today-usage"]').trigger('click')
    activeConfirmDialog(wrapper)?.vm.$emit('confirm')
    await flushPromises()

    expect(showError).toHaveBeenCalledWith('network response lost')

    await wrapper.get('[data-test="reset-today-usage"]').trigger('click')
    activeConfirmDialog(wrapper)?.vm.$emit('confirm')
    await flushPromises()

    expect(resetTodayUsage).toHaveBeenCalledTimes(2)
    expect(resetTodayUsage.mock.calls[0][0]).toBe('11111111-1111-4111-8111-111111111111')
    expect(resetTodayUsage.mock.calls[1][0]).toBe(resetTodayUsage.mock.calls[0][0])
    expect(showSuccess).toHaveBeenCalledTimes(1)
    wrapper.unmount()
  })

  it('rotates the idempotency key after the Shanghai calendar date changes', async () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date('2026-08-01T15:59:59Z'))
    const randomUUID = vi.fn()
      .mockReturnValueOnce('11111111-1111-4111-8111-111111111111')
      .mockReturnValueOnce('22222222-2222-4222-8222-222222222222')
    vi.stubGlobal('crypto', { randomUUID })
    vi.spyOn(console, 'error').mockImplementation(() => undefined)
    resetTodayUsage.mockRejectedValue(new Error('network response lost'))

    const wrapper = mountView()
    await flushPromises()

    await wrapper.get('[data-test="reset-today-usage"]').trigger('click')
    activeConfirmDialog(wrapper)?.vm.$emit('confirm')
    await flushPromises()

    vi.setSystemTime(new Date('2026-08-01T16:00:01Z'))
    await wrapper.get('[data-test="reset-today-usage"]').trigger('click')
    activeConfirmDialog(wrapper)?.vm.$emit('confirm')
    await flushPromises()

    expect(resetTodayUsage.mock.calls[0][0]).toBe('11111111-1111-4111-8111-111111111111')
    expect(resetTodayUsage.mock.calls[1][0]).toBe('22222222-2222-4222-8222-222222222222')
    wrapper.unmount()
  })

  it('shows the normalized backend error message', async () => {
    resetTodayUsage.mockRejectedValue({
      status: 409,
      code: 'IDEMPOTENCY_IN_PROGRESS',
      message: 'The reset is still processing',
      metadata: { retry_after: 2 }
    })
    vi.spyOn(console, 'error').mockImplementation(() => undefined)

    const wrapper = mountView()
    await flushPromises()

    await wrapper.get('[data-test="reset-today-usage"]').trigger('click')
    activeConfirmDialog(wrapper)?.vm.$emit('confirm')
    await flushPromises()

    expect(showError).toHaveBeenCalledWith('The reset is still processing')
    wrapper.unmount()
  })

  it('does not call the API when confirmation is cancelled', async () => {
    const wrapper = mountView()
    await flushPromises()

    await wrapper.get('[data-test="reset-today-usage"]').trigger('click')
    activeConfirmDialog(wrapper)?.vm.$emit('cancel')
    await flushPromises()

    expect(resetTodayUsage).not.toHaveBeenCalled()
    wrapper.unmount()
  })
})
