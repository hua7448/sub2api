import { beforeEach, describe, expect, it, vi } from 'vitest'
import { defineComponent } from 'vue'
import { flushPromises, mount } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import ProviderHallGroups from '../ProviderHallGroups.vue'
import ProviderHallJobs from '../ProviderHallJobs.vue'
import ProviderHallJobDetailDialog from '../ProviderHallJobDetailDialog.vue'
import ProviderHallHealth from '../ProviderHallHealth.vue'
import * as hall from '@/api/admin/providerHall'
import * as groupsAPI from '@/api/admin/groups'
import * as usersAPI from '@/api/admin/users'
import zh from '@/i18n/locales/zh/admin/providerHall'
import en from '@/i18n/locales/en/admin/providerHall'

const showSuccess = vi.fn()
vi.mock('@/api/admin/providerHall', () => ({ getGroup: vi.fn(), getTargets: vi.fn(), updateGroup: vi.fn(), updateTargets: vi.fn(), enqueueProbe: vi.fn(), enqueueVerification: vi.fn(), listJobs: vi.fn(), getJob: vi.fn(), cancelJob: vi.fn(), getHealth: vi.fn() }))
vi.mock('@/api/admin/groups', () => ({ getAllIncludingInactive: vi.fn() }))
vi.mock('@/api/admin/users', () => ({ list: vi.fn(), getById: vi.fn(), getUserApiKeys: vi.fn() }))
vi.mock('@/stores/app', () => ({ useAppStore: () => ({ showSuccess }) }))

const SelectStub = defineComponent({
  props: ['modelValue', 'options'],
  emits: ['update:modelValue'],
  template: '<select :value="modelValue" @change="$emit(\'update:modelValue\', $event.target.value === \'\' ? \'\' : isNaN(Number($event.target.value)) ? $event.target.value : Number($event.target.value))"><option v-for="o in options" :key="o.value" :value="o.value">{{ o.label }}</option></select>',
})
type MessageTree = { [key: string]: string | MessageTree }
// The test build has no message compiler; wrap leaves into interpolating functions.
function runtime(tree: MessageTree): Record<string, unknown> {
  return Object.fromEntries(Object.entries(tree).map(([key, value]) => [key, typeof value === 'string'
    ? (ctx: { named: (name: string) => unknown }) => value.replace(/\{(\w+)\}/g, (_, name: string) => String(ctx.named(name) ?? ''))
    : runtime(value)]))
}
const messages = { zh: runtime({ admin: { providerHall: zh.providerHall }, common: { loading: '加载中', view: '查看', cancel: '取消', close: '关闭', delete: '删除', save: '保存', confirm: '确认', autoRefresh: { title: '自动刷新', enable: '启用', countdown: '{seconds}s', seconds: '{n}s' } } }) }
const globalOptions = () => ({
  plugins: [createI18n({ legacy: false, locale: 'zh', missingWarn: false, fallbackWarn: false, messages })],
  stubs: { Select: SelectStub, BaseDialog: { props: ['show', 'title'], template: '<div v-if="show" role="dialog"><h2>{{ title }}</h2><slot /><slot name="footer" /></div>' }, StatCard: { props: ['title', 'value'], template: '<div class="stat"><span>{{ title }}</span><b>{{ value }}</b></div>' }, AutoRefreshButton: true, Pagination: { props: ['total', 'page', 'pageSize'], template: '<nav class="pagination">{{ total }}</nav>' } },
})
const stamp = '2026-09-12T00:00:00Z'
const profile: hall.ProviderHallProfile = { id: 7, version: 1, updated_at: stamp, updated_by: 1, model: 'gpt-test', protocol: 'responses', supports_tools: false, output_limit: 256, model_aliases: [], reference_input_price: null, reference_cache_price: null, reference_cache_rate: null, reference_confirmed_at: null }
const group = (id: number): hall.ProviderHallGroup => ({ group_id: id, listed: true, display_name: `Group ${id}`, description: '', display_order: 0, version: 4, updated_at: stamp, updated_by: 1 })
const targets = (id: number, enabled = true): hall.ProviderHallTargetSet => ({ group_id: id, version: 4, updated_at: stamp, updated_by: 1,
  items: [{ id: 8, group_id: id, version: 1, updated_at: stamp, updated_by: 1, profile_id: 7, probe_key_id: 5, enabled, probe_interval_seconds: 300, verification_interval_seconds: 86400 }] })
const job = (id: number, status: hall.ProviderHallJobStatus): hall.ProviderHallJob => ({ id, kind: 'probe', target_id: 8, group_id: 1, profile_id: 7, status, attempts: 1, error_code: '', error_message: '', created_at: stamp, updated_at: stamp,
  config_snapshot: { profile: { id: 7, version: 1, model: 'gpt-test', protocol: 'responses', supports_tools: false, output_limit: 256, model_aliases: [] }, target_version: 1, probe_key_id: 5, operator_user_id: 2, gateway_origin: 'https://gw.example' },
  slot_at: null, not_before: null, idempotency_key: null, requested_by: null, lease_owner: null, lease_until: null, budget_day: null, started_at: null, finished_at: null })

beforeEach(() => {
  vi.clearAllMocks()
  vi.mocked(usersAPI.list).mockResolvedValue({ items: [], total: 0, page: 1, page_size: 50, pages: 0 })
  vi.mocked(groupsAPI.getAllIncludingInactive).mockResolvedValue([{ id: 1, name: 'One', platform: 'openai' }] as Awaited<ReturnType<typeof groupsAPI.getAllIncludingInactive>>)
  vi.mocked(usersAPI.getUserApiKeys).mockResolvedValue({ items: [], total: 0, page: 1, page_size: 100, pages: 0 })
  vi.mocked(hall.getGroup).mockImplementation(async id => group(id))
  vi.mocked(hall.getTargets).mockImplementation(async id => targets(id))
})

describe('Provider Hall task buttons', () => {
  it('reuses one idempotency key across a double click and reports reuse', async () => {
    let resolveFirst!: (value: hall.ProviderHallEnqueueResult) => void
    vi.mocked(hall.enqueueProbe).mockImplementationOnce(() => new Promise(resolve => { resolveFirst = resolve }))
    const wrapper = mount(ProviderHallGroups, { props: { operatorId: 2, profiles: [profile] }, global: globalOptions() })
    await flushPromises()
    const probe = wrapper.get('button[aria-label="立即探测 gpt-test"]')
    await probe.trigger('click')
    await probe.trigger('click')
    expect(hall.enqueueProbe).toHaveBeenCalledTimes(1)
    expect(probe.attributes('disabled')).toBeDefined()
    resolveFirst({ job_id: 42, status: 'queued', reused: false })
    await flushPromises()
    expect(showSuccess).toHaveBeenCalledWith('已加入队列：任务 #42')
    const firstKey = vi.mocked(hall.enqueueProbe).mock.calls[0][2]
    expect(firstKey).toMatch(/^admin-probe-1-7-/)
    expect(firstKey.length).toBeLessThanOrEqual(128)
    vi.mocked(hall.enqueueProbe).mockResolvedValueOnce({ job_id: 42, status: 'queued', reused: true })
    await probe.trigger('click')
    await flushPromises()
    expect(vi.mocked(hall.enqueueProbe).mock.calls[1][2]).not.toBe(firstKey)
    expect(showSuccess).toHaveBeenLastCalledWith('已存在相同任务 #42，未重复创建')
    wrapper.unmount()
  })

  it('keeps the key for a transport retry but drops it after a business rejection', async () => {
    vi.mocked(hall.enqueueVerification).mockRejectedValueOnce(new Error('network'))
    const wrapper = mount(ProviderHallGroups, { props: { operatorId: 2, profiles: [profile] }, global: globalOptions() })
    await flushPromises()
    const verify = wrapper.get('button[aria-label="立即检测 gpt-test"]')
    await verify.trigger('click')
    await flushPromises()
    vi.mocked(hall.enqueueVerification).mockRejectedValueOnce({ reason: 'PROVIDER_HALL_BUDGET_EXHAUSTED' })
    await verify.trigger('click')
    await flushPromises()
    const calls = vi.mocked(hall.enqueueVerification).mock.calls
    expect(calls[1][2]).toBe(calls[0][2])
    expect(wrapper.get('[role="alert"]').text()).toContain('今日预算已用尽')
    vi.mocked(hall.enqueueVerification).mockResolvedValueOnce({ job_id: 9, status: 'queued', reused: false })
    await verify.trigger('click')
    await flushPromises()
    expect(calls[2][2]).not.toBe(calls[0][2])
    wrapper.unmount()
  })

  it('disables the buttons for disabled or unsaved targets', async () => {
    vi.mocked(hall.getTargets).mockImplementation(async id => targets(id, false))
    const wrapper = mount(ProviderHallGroups, { props: { operatorId: 2, profiles: [profile] }, global: globalOptions() })
    await flushPromises()
    expect(wrapper.get('button[aria-label="立即探测 gpt-test"]').attributes('disabled')).toBeDefined()
    wrapper.unmount()
  })
})

describe('Provider Hall jobs tab', () => {
  it('lists jobs, filters, and cancels through the confirm dialog', async () => {
    vi.mocked(hall.listJobs).mockResolvedValue({ items: [job(1, 'queued'), job(2, 'succeeded')], total: 2, page: 1, page_size: 20, pages: 1 })
    vi.mocked(hall.cancelJob).mockResolvedValue({ job_id: 1, status: 'cancelled' })
    const wrapper = mount(ProviderHallJobs, { global: globalOptions() })
    await flushPromises()
    expect(wrapper.findAll('tbody tr')).toHaveLength(2)
    expect(wrapper.text()).toContain('排队中')
    expect(wrapper.find('button[aria-label="取消任务 #2"]').exists()).toBe(false)
    await wrapper.get('#hall-job-status').setValue('running')
    await flushPromises()
    expect(hall.listJobs).toHaveBeenLastCalledWith(expect.objectContaining({ status: 'running', page: 1 }), expect.any(AbortSignal))
    await wrapper.get('button[aria-label="取消任务 #1"]').trigger('click')
    expect(wrapper.get('[role="dialog"]').text()).toContain('确定取消任务 #1')
    await wrapper.findAll('[role="dialog"] button').find(b => b.text() === '确认')!.trigger('click')
    await flushPromises()
    expect(hall.cancelJob).toHaveBeenCalledWith(1)
    expect(showSuccess).toHaveBeenCalledWith('任务已取消')
    expect(hall.listJobs).toHaveBeenCalledTimes(3)
    await wrapper.get('button[aria-label="任务详情 #1"]').trigger('click')
    expect(wrapper.emitted('open')).toEqual([[1]])
    wrapper.unmount()
  })

  it('shows the in-flight conflict clearly and keeps the list', async () => {
    vi.mocked(hall.listJobs).mockResolvedValue({ items: [job(3, 'running')], total: 1, page: 1, page_size: 20, pages: 1 })
    vi.mocked(hall.cancelJob).mockRejectedValue({ reason: 'PROVIDER_HALL_JOB_IN_FLIGHT' })
    const wrapper = mount(ProviderHallJobs, { global: globalOptions() })
    await flushPromises()
    await wrapper.get('button[aria-label="取消任务 #3"]').trigger('click')
    await wrapper.findAll('[role="dialog"] button').find(b => b.text() === '确认')!.trigger('click')
    await flushPromises()
    expect(wrapper.get('[role="alert"]').text()).toContain('任务已有样本发出')
    expect(wrapper.findAll('tbody tr')).toHaveLength(1)
    expect(showSuccess).not.toHaveBeenCalled()
    wrapper.unmount()
  })

  it('renders job detail with samples and the report summary', async () => {
    vi.mocked(hall.getJob).mockResolvedValue({ job: { ...job(5, 'succeeded'), kind: 'verification' },
      samples: [{ id: 1, job_id: 5, test_id: 'arith', seq: 1, trace_id: 't', status: 'received', prepared_at: stamp, dispatched_at: stamp, received_at: stamp, client_request_id: 'c', http_status: 200, ttft_ms: 120, total_ms: 800, generation_ms: 680, input_tokens: 20, output_tokens: 3, response_model: 'gpt-test', result: 'passed', error_code: '', detail: {}, billing_status: 'confirmed', actual_cost: '0.00001' }],
      verification: { job_id: 5, group_id: 1, profile_id: 7, target_id: 8, verdict: 'passed', execution_status: 'completed', reason_code: '', summary: { arithmetic: { passed: 3, failed: 0, error: 0 }, json: { passed: 3, failed: 0, error: 0 }, tool: null, model: { matched: 6, mismatched: 0, missing: 0, seen: ['gpt-test'] } }, profile_version: 1, target_version: 1, completed_at: stamp, expires_at: stamp, stale: false, created_at: stamp } })
    const wrapper = mount(ProviderHallJobDetailDialog, { props: { jobId: 5 }, global: globalOptions() })
    await flushPromises()
    const text = wrapper.text()
    expect(text).toContain('通过')
    expect(text).toContain('3 通过 / 0 失败 / 0 异常')
    expect(text).toContain('6 匹配 / 0 不符 / 0 缺失')
    expect(text).toContain('120 ms / 800 ms')
    expect(text).not.toContain('工具调用')
    expect(text).not.toContain('gateway_origin')
    wrapper.unmount()
  })
})

describe('Provider Hall health tab', () => {
  it('renders nodes, gaps and paused reason', async () => {
    vi.mocked(hall.getHealth).mockResolvedValue({ generated_at: stamp,
      collection: { enabled: true, nodes: [{ node_id: 'node-a', epoch_id: 1, version: 'v1', heartbeat_at: stamp, confirmed_at: stamp, persisted_seq: 10, overflowed: false, lost: false }, { node_id: 'node-b', epoch_id: 2, version: 'v1', heartbeat_at: stamp, confirmed_at: null, persisted_seq: 0, overflowed: true, lost: true }], missing_expected: ['node-z'], open_gaps: [{ id: 1, node_id: 'node-b', epoch_id: 2, scope: 'collection', started_at: stamp, reason: 'queue_overflow' }], local_queue: { depth: 3, capacity: 8192, dropped: 0 } },
      aggregator: { watermark: stamp, lag_seconds: 90, dirty_count: 2, last_run_at: stamp, last_error: '' },
      reconciliation: { pending: 1, uncertain: 0, failed_24h: 0 },
      budget: { day: '2026-09-12', budget: '10', confirmed_spend: '10.5', uncertain_spend: '0', in_flight: 0, paused_reason: 'budget_exhausted' },
      jobs: { queued: 0, running: 0, unknown: 0, failed_24h: 0 } })
    const wrapper = mount(ProviderHallHealth, { global: globalOptions() })
    await flushPromises()
    const text = wrapper.text()
    expect(text).toContain('node-a')
    expect(text).toContain('失联')
    expect(text).toContain('溢出')
    expect(text).toContain('缺席的预期节点: node-z')
    expect(text).toContain('queue_overflow')
    expect(text).toContain('预算已用尽')
    expect(text).toContain('3 / 8192，已丢弃 0')
    wrapper.unmount()
  })

  it('has matching Chinese and English locale keys', () => {
    expect(Object.keys(zh.providerHall).sort()).toEqual(Object.keys(en.providerHall).sort())
  })
})
