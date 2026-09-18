import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { createMemoryHistory, createRouter } from 'vue-router'
import ProviderHallView from '../ProviderHallView.vue'
import * as api from '@/api/providerHall'
import { createHallI18n, makeDetail, makeList, makeRow, metric } from '@/components/provider-hall/__tests__/testUtils'

vi.mock('@/components/layout/AppLayout.vue', () => ({ default: { template: '<main><slot /></main>' } }))
vi.mock('@/stores/auth', () => ({ useAuthStore: () => ({ user: { id: 42 } }) }))
vi.mock('@/stores/app', () => ({ useAppStore: () => ({ showError: vi.fn(), showSuccess: vi.fn() }) }))
vi.mock('@/api/keys', () => ({ keysAPI: { list: vi.fn(async () => ({ items: [], total: 0, page: 1, page_size: 20, pages: 0 })), update: vi.fn(), create: vi.fn() } }))
vi.mock('vue-chartjs', () => ({ Line: { name: 'Line', template: '<canvas data-testid="line-chart"></canvas>' } }))
vi.mock('@/api/providerHall', async () => {
  const actual = await vi.importActual<typeof import('@/api/providerHall')>('@/api/providerHall')
  return { ...actual, list: vi.fn(), getGroup: vi.fn(), listVerifications: vi.fn() }
})

async function render(path = '/providers') {
  const router = createRouter({ history: createMemoryHistory(), routes: [{ path: '/providers', component: ProviderHallView }] })
  await router.push(path)
  await router.isReady()
  const wrapper = mount(ProviderHallView, {
    attachTo: document.body,
    global: {
      plugins: [router, createHallI18n()],
      stubs: {
        Teleport: true,
        Select: { props: ['modelValue', 'options'], template: '<select><option v-for="o in options" :key="o.value" :value="o.value">{{ o.label }}</option></select>' },
        AutoRefreshButton: true,
      },
    },
  })
  await flushPromises()
  return { wrapper, router }
}

const rows = [
  makeRow({ group_id: 1, name: 'A015-Pro' }),
  makeRow({
    group_id: 2,
    name: 'A018-Plus',
    health: { status: 'down', checked_at: null, models: [] },
    metrics: { ...makeRow().metrics, ttft_fast95_ms: metric<number>(null, 'insufficient', 'samples_below_20') },
    verification: { verdict: 'suspected', reason_code: 'model_mismatch', completed_at: '2026-09-12T01:00:00Z', expired: false, report_id: 5 },
  }),
]

beforeEach(() => {
  vi.clearAllMocks()
  localStorage.clear()
  vi.mocked(api.list).mockResolvedValue(makeList(rows))
  vi.mocked(api.getGroup).mockImplementation(async (id) => makeDetail(rows.find((r) => r.group_id === id)!))
  vi.mocked(api.listVerifications).mockResolvedValue({
    items: [{ report_id: 5, profile_id: 1, model: 'astra', protocol: 'responses', verdict: 'suspected', execution_status: 'completed', reason_code: 'model_mismatch', completed_at: '2026-09-12T01:00:00Z', expires_at: '2026-09-14T01:00:00Z', expired: false, summary: { arithmetic: { passed: 3, failed: 0, error: 0 }, json: { passed: 3, failed: 0, error: 0 }, tool: null, model: { matched: 1, mismatched: 2, missing: 0, seen: ['other'] } } }],
    pagination: { page: 1, page_size: 20, total: 1 },
  })
})

describe('ProviderHallView', () => {
  it('loads the list, renders summary, rows and toolbar state', async () => {
    const { wrapper } = await render()
    expect(api.list).toHaveBeenCalledTimes(1)
    expect(wrapper.get('[data-testid="summary-available"]').text()).toBe('1')
    expect(wrapper.get('[data-testid="summary-listed"]').text()).toBe('2')
    expect(wrapper.get('[data-testid="summary-abnormal"]').text()).toBe('1')
    expect(wrapper.get('[data-testid="summary-verified"]').text()).toBe('1')
    expect(wrapper.findAll('[data-testid="hall-row"]')).toHaveLength(2)
    expect(wrapper.get('[data-testid="hall-updated"]').text()).toContain('2026')
    expect(wrapper.get('[data-range="6h"]').classes()).toContain('tab-active')
    expect(wrapper.text()).toContain('样本不足')
    wrapper.unmount()
  })

  it('changes range and sort through the toolbar and refetches', async () => {
    const { wrapper, router } = await render()
    await wrapper.get('[data-range="24h"]').trigger('click')
    await flushPromises()
    expect(router.currentRoute.value.query.range).toBe('24h')
    expect(api.list).toHaveBeenLastCalledWith(expect.objectContaining({ range: '24h' }), expect.any(AbortSignal))
    await wrapper.get('[data-sort="ttft_fast95"]').trigger('click')
    await flushPromises()
    expect(wrapper.get('[data-sort="ttft_fast95"]').classes()).toContain('tab-active')
    expect(wrapper.get('[data-testid="sort-arrow"]').text()).toBe('↓')
    expect(router.currentRoute.value.query.sort).toBe('ttft_fast95:desc')
    expect(localStorage.getItem('providerHall.sort.42')).toBe('ttft_fast95:desc')
    wrapper.unmount()
  })

  it('expands a row to load details and opens the report dialog', async () => {
    const { wrapper } = await render()
    await wrapper.findAll('[data-testid="hall-row"]')[0].trigger('click')
    await flushPromises()
    expect(api.getGroup).toHaveBeenCalledWith(1, '6h', expect.any(AbortSignal))
    const detail = wrapper.get('[data-testid="hall-detail-row"]')
    expect(detail.text()).toContain('27.59/s')
    expect(detail.find('[data-testid="line-chart"]').exists()).toBe(true)
    await wrapper.findAll('[data-testid="hall-row"]')[1].get('[data-testid="verification-report-link"]').trigger('click')
    await flushPromises()
    expect(api.listVerifications).toHaveBeenCalledWith(2, expect.objectContaining({ page: 1 }), expect.any(AbortSignal))
    const dialog = wrapper.get('[data-testid="verification-dialog"]')
    expect(dialog.find('[data-testid="verification-report"]').exists()).toBe(true)
    expect(dialog.text()).toContain('疑似')
    expect(dialog.text()).toContain('不符 2')
    wrapper.unmount()
  })

  it('opens the use-group dialog from a row', async () => {
    const { wrapper } = await render()
    await wrapper.findAll('[data-testid="hall-row"]')[0].get('[data-testid="use-group"]').trigger('click')
    await flushPromises()
    expect(wrapper.get('[data-testid="use-group-dialog"]').exists()).toBe(true)
    expect(wrapper.get('.modal-title').text()).toBe('使用分组 A015-Pro')
    expect(wrapper.findAll('[data-testid="hall-detail-row"]')).toHaveLength(0)
    wrapper.unmount()
  })

  it('shows the stale banner after a failed refresh and the error state when nothing loaded', async () => {
    const { wrapper } = await render()
    vi.mocked(api.list).mockRejectedValueOnce({ message: 'down' })
    await wrapper.get('[data-testid="hall-refresh"]').trigger('click')
    await flushPromises()
    expect(wrapper.get('[data-testid="hall-stale"]').text()).toContain('down')
    expect(wrapper.findAll('[data-testid="hall-row"]')).toHaveLength(2)
    wrapper.unmount()

    vi.mocked(api.list).mockRejectedValueOnce({ message: 'first failure' })
    const failed = await render()
    expect(failed.wrapper.get('[data-testid="hall-error"]').text()).toContain('first failure')
    expect(failed.wrapper.find('[data-testid="hall-stale"]').exists()).toBe(false)
    failed.wrapper.unmount()
  })
})
