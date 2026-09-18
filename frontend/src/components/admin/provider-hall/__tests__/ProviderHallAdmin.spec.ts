import { beforeEach, describe, expect, it, vi } from 'vitest'
import { defineComponent } from 'vue'
import { flushPromises, mount } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import ProviderHallConfigForm from '../ProviderHallConfigForm.vue'
import ProviderHallGroups from '../ProviderHallGroups.vue'
import ProviderHallProfiles from '../ProviderHallProfiles.vue'
import * as hall from '@/api/admin/providerHall'
import * as groupsAPI from '@/api/admin/groups'
import * as usersAPI from '@/api/admin/users'
import zh from '@/i18n/locales/zh/admin/providerHall'
import en from '@/i18n/locales/en/admin/providerHall'

// The test build uses the runtime-only vue-i18n bundle, which cannot compile
// plain strings; wrap every string (recursively) as a message function.
function asMessageFunctions(value: unknown): unknown {
  if (typeof value === 'string') return () => value
  if (value && typeof value === 'object') return Object.fromEntries(Object.entries(value as Record<string, unknown>).map(([k, v]) => [k, asMessageFunctions(v)]))
  return value
}

vi.mock('@/api/admin/providerHall', () => ({ getHealth: vi.fn(async () => ({ collection: { nodes: [] } })), preflightConfig: vi.fn(async () => ({ valid: true, disabled_targets: [] })), listGroups: vi.fn(), listModels: vi.fn(async () => []), listProbeKeys: vi.fn(async () => []), saveSettings: vi.fn(), updateConfig: vi.fn(), getGroup: vi.fn(), getTargets: vi.fn(), updateGroup: vi.fn(), updateTargets: vi.fn(), createProfile: vi.fn(), updateProfile: vi.fn() }))
vi.mock('@/api/admin/groups', () => ({ getAllIncludingInactive: vi.fn() }))
vi.mock('@/api/admin/users', () => ({ list: vi.fn(), getById: vi.fn(), getUserApiKeys: vi.fn() }))
vi.mock('@/stores/app', () => ({ useAppStore: () => ({ showSuccess: vi.fn() }) }))

const SelectStub = defineComponent({
  props: ['modelValue', 'options'],
  emits: ['update:modelValue'],
  template: '<select :value="modelValue" @change="$emit(\'update:modelValue\', $event.target.value ? Number($event.target.value) : null)"><option v-for="o in options" :key="o.value" :value="o.value">{{ o.label }}</option></select>',
})
const globalOptions = () => ({ plugins: [createI18n({ legacy: false, locale: 'zh', missingWarn: false, fallbackWarn: false,
  messages: { zh: { admin: { providerHall: asMessageFunctions(zh.providerHall) } } } })],
  stubs: { Select: SelectStub, BaseDialog: { props: ['show'], template: '<div v-if="show"><slot /><slot name="footer" /></div>' } } })
const config: hall.ProviderHallConfig = { version: 3, updated_at: '2026-09-10T00:00:00Z', updated_by: 1,
  collection_enabled: false, display_enabled: false, tasks_enabled: false, default_model: '', default_protocol: 'responses',
  default_range: '6h', gateway_origin: '', operator_user_id: null, daily_budget: '0.00000000', expected_nodes: [] }
const profile: hall.ProviderHallProfile = { id: 7, version: 1, updated_at: config.updated_at, updated_by: 1, model: 'gpt-test', protocol: 'responses',
  supports_tools: false, output_limit: 256, model_aliases: [], reference_input_price: null, reference_cache_price: null,
  reference_cache_rate: null, reference_confirmed_at: null }
const group = (id: number): hall.ProviderHallGroup => ({ group_id: id, listed: false, display_name: `Group ${id}`, description: '', display_order: 0, version: 4, updated_at: config.updated_at, updated_by: 1 })
const targets = (id: number): hall.ProviderHallTargetSet => ({ group_id: id, version: 4, updated_at: config.updated_at, updated_by: 1,
  items: [{ id: 8, group_id: id, version: 1, updated_at: config.updated_at, updated_by: 1, profile_id: 7, probe_key_id: 5,
    enabled: true, probe_interval_seconds: 300, verification_interval_seconds: 86400 }] })

beforeEach(() => {
  vi.clearAllMocks()
  vi.mocked(usersAPI.list).mockResolvedValue({ items: [], total: 0, page: 1, page_size: 50, pages: 0 })
  vi.mocked(groupsAPI.getAllIncludingInactive).mockResolvedValue([{ id: 1, name: 'One', platform: 'openai' }, { id: 2, name: 'Two', platform: 'composite' }] as Awaited<ReturnType<typeof groupsAPI.getAllIncludingInactive>>)
  vi.mocked(usersAPI.getUserApiKeys).mockResolvedValue({ items: [], total: 0, page: 1, page_size: 100, pages: 0 })
  vi.mocked(hall.listGroups).mockResolvedValue({ items: [1, 2].map(id => ({ ...group(id), name: `Group ${id}`, platform: 'openai', status: 'active', supported: true, reason: '', targets: targets(id).items, models: ['gpt-test'], effective_profile_id: 7 })), total: 2, page: 1, page_size: 20 })
  vi.mocked(hall.getGroup).mockImplementation(async id => group(id))
  vi.mocked(hall.getTargets).mockImplementation(async id => targets(id))
})

describe('Provider Hall admin configuration', () => {
  it('keeps decimal precision, strips audit fields and leaves runtime switches disabled when the build is not ready', async () => {
    const wrapper = mount(ProviderHallConfigForm, { props: { config, profiles: [profile] }, global: globalOptions() })
    await flushPromises()
    await wrapper.get('input[inputmode="decimal"]').setValue('123456789012.12345678')
    await wrapper.get('textarea').setValue(' node-a\n\nnode-b ')
    vi.mocked(hall.updateConfig).mockResolvedValue({ ...config, version: 4 })
    await wrapper.get('form').trigger('submit')
    await flushPromises()
    expect(hall.updateConfig).toHaveBeenCalledWith(expect.objectContaining({ daily_budget: '123456789012.12345678', expected_nodes: ['node-a', 'node-b'], collection_enabled: false, display_enabled: false, tasks_enabled: false }))
    expect(hall.updateConfig).not.toHaveBeenCalledWith(expect.objectContaining({ updated_by: expect.anything() }))
    expect(wrapper.findAll('input[type="checkbox"]').slice(0, 3).every(input => input.attributes('disabled') !== undefined)).toBe(true)
    wrapper.unmount()
  })

  it('names the offending field when the backend rejects the config', async () => {
    const ready = { ...config, readiness: { collection: true, tasks: true, display: true } }
    const wrapper = mount(ProviderHallConfigForm, { props: { config: ready, profiles: [profile] }, global: globalOptions() })
    await flushPromises()
    await wrapper.findAll('input[type="checkbox"]')[2].setValue(true)
    vi.mocked(hall.updateConfig).mockRejectedValue({ status: 400, code: 'PROVIDER_HALL_INVALID_CONFIG', message: 'invalid', metadata: { field: 'display_enabled' } })
    await wrapper.get('form').trigger('submit')
    await flushPromises()
    expect(wrapper.get('[role="alert"]').text()).toContain('用户大厅需要先开启真实请求采集')
    expect(wrapper.get('input[type="url"]').exists()).toBe(true)
    wrapper.unmount()
  })

  it('enables exactly the ready switches and submits their values', async () => {
    const ready = { ...config, readiness: { collection: true, tasks: true, display: false } }
    const wrapper = mount(ProviderHallConfigForm, { props: { config: ready, profiles: [profile] }, global: globalOptions() })
    await flushPromises()
    const boxes = wrapper.findAll('input[type="checkbox"]')
    expect(boxes.map(b => b.attributes('disabled') !== undefined)).toEqual([false, false, true, false])
    await boxes[0].setValue(true)
    vi.mocked(hall.updateConfig).mockResolvedValue({ ...ready, collection_enabled: true, version: 4 })
    await wrapper.get('form').trigger('submit')
    await flushPromises()
    expect(hall.updateConfig).toHaveBeenCalledWith(expect.objectContaining({ collection_enabled: true, tasks_enabled: false, display_enabled: false }))
    expect(hall.updateConfig).not.toHaveBeenCalledWith(expect.objectContaining({ readiness: expect.anything() }))
    wrapper.unmount()
  })

  it('keeps config input after a version conflict', async () => {
    const wrapper = mount(ProviderHallConfigForm, { props: { config, profiles: [] }, global: globalOptions() })
    await flushPromises()
    await wrapper.get('input[type="url"]').setValue('https://gateway.example.test')
    vi.mocked(hall.updateConfig).mockRejectedValue({ reason: 'PROVIDER_HALL_VERSION_CONFLICT' })
    await wrapper.get('form').trigger('submit')
    await flushPromises()
    expect(wrapper.get('[role="alert"]').text()).toContain('当前输入已保留')
    expect((wrapper.get('input[type="url"]').element as HTMLInputElement).value).toBe('https://gateway.example.test')
    expect(wrapper.emitted('saved')).toBeUndefined()
    wrapper.unmount()
  })

  it('ignores a delayed group response after switching groups', async () => {
    let resolveFirst!: (value: hall.ProviderHallGroup) => void
    let firstSignal: AbortSignal | undefined
    vi.mocked(hall.getGroup).mockImplementation((id, signal) => {
      if (id === 1) { firstSignal = signal; return new Promise(resolve => { resolveFirst = resolve }) }
      return Promise.resolve(group(id))
    })
    const wrapper = mount(ProviderHallGroups, { props: { operatorId: null, profiles: [profile] }, global: globalOptions() })
    await flushPromises()
    await wrapper.findAll('button').find(b => b.text().includes('管理'))!.trigger('click')
    await flushPromises()
    await wrapper.findAll('button').filter(b => b.text().includes('管理'))[1].trigger('click')
    await flushPromises()
    expect(firstSignal?.aborted).toBe(true)
    expect((wrapper.get('input[maxlength="100"]').element as HTMLInputElement).value).toBe('Group 2')
    resolveFirst(group(1))
    await flushPromises()
    expect((wrapper.get('input[maxlength="100"]').element as HTMLInputElement).value).toBe('Group 2')
    wrapper.unmount()
  })

  it('saves listing and targets together and retains the draft on conflict', async () => {
    const wrapper = mount(ProviderHallGroups, { props: { operatorId: null, profiles: [profile] }, global: globalOptions() })
    await flushPromises()
    await wrapper.findAll('button').find(b => b.text().includes('管理'))!.trigger('click')
    await flushPromises()
    await wrapper.get('input[max="1440"]').setValue('10')
    vi.mocked(hall.saveSettings).mockRejectedValue({ reason: 'PROVIDER_HALL_VERSION_CONFLICT' })
    await wrapper.get('form').trigger('submit')
    await flushPromises()
    expect(hall.saveSettings).toHaveBeenCalledWith(1, expect.objectContaining({ version: 4, items: [expect.objectContaining({ profile_id: 7, probe_interval_seconds: 600 })] }))
    expect((wrapper.get('input[max="1440"]').element as HTMLInputElement).value).toBe('10')
    expect(wrapper.get('[role="alert"]').text()).toContain('当前输入已保留')
    wrapper.unmount()
  })

  it('rejects an incoherent group and target snapshot', async () => {
    vi.mocked(hall.getTargets).mockResolvedValue({ ...targets(1), version: 5 })
    const wrapper = mount(ProviderHallGroups, { props: { operatorId: null, profiles: [] }, global: globalOptions() })
    await flushPromises()
    await wrapper.findAll('button').find(b => b.text().includes('管理'))!.trigger('click')
    await flushPromises()
    expect(wrapper.get('[role="alert"]').text()).toContain('配置已被修改')
    expect(wrapper.find('form').exists()).toBe(false)
    wrapper.unmount()
  })

  it('requires explicit reference confirmation and preserves legitimate zero prices', async () => {
    const wrapper = mount(ProviderHallProfiles, { props: { profiles: [] }, global: globalOptions() })
    await wrapper.findAll('button').find(b => b.text().includes('新增档案'))!.trigger('click')
    await wrapper.get('input[maxlength="200"]').setValue('gpt-zero')
    await wrapper.findAll('input[type="checkbox"]')[1].setValue(true)
    for (const input of wrapper.findAll('input[inputmode="decimal"]')) await input.setValue('0')
    await wrapper.get('input[max="100"]').setValue('0')
    await wrapper.get('form').trigger('submit')
    expect(hall.createProfile).not.toHaveBeenCalled()
    await wrapper.findAll('button').find(b => b.text().includes('确认当前基准'))!.trigger('click')
    vi.mocked(hall.createProfile).mockResolvedValue(profile)
    await wrapper.get('form').trigger('submit')
    await flushPromises()
    expect(hall.createProfile).toHaveBeenCalledWith(expect.objectContaining({ reference_input_price: '0', reference_cache_price: '0', reference_cache_rate: '0', reference_confirmed_at: expect.any(String) }))
    wrapper.unmount()
  })

  it('has matching Chinese and English locale keys', () => {
    expect(Object.keys(zh.providerHall).sort()).toEqual(Object.keys(en.providerHall).sort())
  })
})
