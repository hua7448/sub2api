import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { enableAutoUnmount, flushPromises, mount } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import ProviderHallView from '../ProviderHallView.vue'
import ProviderHallConfigForm from '@/components/admin/provider-hall/ProviderHallConfigForm.vue'
import * as hall from '@/api/admin/providerHall'
import zh from '@/i18n/locales/zh/admin/providerHall'

vi.mock('vue-router', () => ({ onBeforeRouteLeave: vi.fn() }))
vi.mock('@/components/layout/AppLayout.vue', () => ({ default: { template: '<main><slot /></main>' } }))
vi.mock('@/api/admin/providerHall', () => ({ getConfig: vi.fn(), listProfiles: vi.fn(), updateConfig: vi.fn() }))
vi.mock('@/api/admin/users', () => ({ list: vi.fn(async () => ({ items: [] })) }))
vi.mock('@/stores/app', () => ({ useAppStore: () => ({ showSuccess: vi.fn() }) }))

enableAutoUnmount(afterEach)
afterEach(() => { vi.restoreAllMocks() })

const config: hall.ProviderHallConfig = {
  version: 3, updated_at: '2026-09-11T00:00:00Z', updated_by: 1,
  collection_enabled: false, display_enabled: false, tasks_enabled: false,
  default_model: '', default_protocol: 'responses', default_range: '6h',
  gateway_origin: '', operator_user_id: null, daily_budget: '1', expected_nodes: [],
}
function deferred<T>() {
  let resolve!: (value: T) => void
  let reject!: (reason: unknown) => void
  const promise = new Promise<T>((yes, no) => { resolve = yes; reject = no })
  return { promise, resolve, reject }
}
function render() {
  return mount(ProviderHallView, {
    global: {
      plugins: [createI18n({ legacy: false, locale: 'zh', missingWarn: false, fallbackWarn: false,
        messages: { zh: { admin: { providerHall: Object.fromEntries(Object.entries(zh.providerHall).map(([key, value]) => [key, () => value])) } } } })],
      stubs: { AppLayout: { template: '<main><slot /></main>' }, ProviderHallGroups: true, ProviderHallProfiles: true, Select: true },
    },
  })
}

beforeEach(() => {
  vi.clearAllMocks()
  vi.mocked(hall.getConfig).mockResolvedValue({ ...config })
  vi.mocked(hall.listProfiles).mockResolvedValue([])
})

describe('Provider Hall reload ordering', () => {
  it('blocks editing and saving until a delayed config reload finishes', async () => {
    const wrapper = render()
    await flushPromises()
    const form = wrapper.getComponent(ProviderHallConfigForm)
    const budget = form.get('input[inputmode="decimal"]')
    const pending = deferred<hall.ProviderHallConfig>()
    vi.mocked(hall.getConfig).mockReturnValueOnce(pending.promise)
    await form.get('button[type="button"]').trigger('click')
    expect(budget.element.matches(':disabled')).toBe(true)
    expect(form.get('button[type="submit"]').attributes('disabled')).toBeDefined()
    expect(form.get('button[type="button"]').attributes('disabled')).toBeDefined()
    // A synthetic submit also has to obey the loading guard.
    await form.get('form').trigger('submit')
    expect(hall.updateConfig).not.toHaveBeenCalled()
    pending.resolve({ ...config })
    await flushPromises()
    expect(budget.element.matches(':disabled')).toBe(false)
    await budget.setValue('9')
    const event = new Event('beforeunload', { cancelable: true })
    window.dispatchEvent(event)
    expect(event.defaultPrevented).toBe(true)
    expect((budget.element as HTMLInputElement).value).toBe('9')
    wrapper.unmount()
  })

  it('prevents reload and save overlap so a saved version cannot regress', async () => {
    const wrapper = render()
    await flushPromises()
    const form = wrapper.getComponent(ProviderHallConfigForm)
    const pending = deferred<hall.ProviderHallConfig>()
    vi.mocked(hall.getConfig).mockReturnValueOnce(pending.promise)
    await form.get('button[type="button"]').trigger('click')
    await form.get('form').trigger('submit')
    expect(hall.updateConfig).not.toHaveBeenCalled()
    pending.resolve({ ...config })
    await flushPromises()
    await form.get('input[inputmode="decimal"]').setValue('9')
    const saving = deferred<hall.ProviderHallConfig>()
    vi.mocked(hall.updateConfig).mockReturnValueOnce(saving.promise)
    await form.get('form').trigger('submit')
    await form.get('button[type="button"]').trigger('click')
    expect(hall.getConfig).toHaveBeenCalledTimes(2)
    saving.resolve({ ...config, version: 4, daily_budget: '9' })
    await flushPromises()
    expect(form.props('config')).toMatchObject({ version: 4, daily_budget: '9' })
    const event = new Event('beforeunload', { cancelable: true })
    window.dispatchEvent(event)
    expect(event.defaultPrevented).toBe(false)
    wrapper.unmount()
  })

  it('unlocks failed reloads while preserving the old unsaved input and dirty state', async () => {
    const wrapper = render()
    await flushPromises()
    const form = wrapper.getComponent(ProviderHallConfigForm)
    await form.get('input[inputmode="decimal"]').setValue('9')
    const confirm = vi.spyOn(window, 'confirm').mockReturnValue(true)
    const pending = deferred<hall.ProviderHallConfig>()
    vi.mocked(hall.getConfig).mockReturnValueOnce(pending.promise)
    await form.get('button[type="button"]').trigger('click')
    expect(form.get('input[inputmode="decimal"]').element.matches(':disabled')).toBe(true)
    pending.reject(new Error('offline'))
    await flushPromises()
    const budget = form.get('input[inputmode="decimal"]')
    expect(budget.element.matches(':disabled')).toBe(false)
    expect((budget.element as HTMLInputElement).value).toBe('9')
    const event = new Event('beforeunload', { cancelable: true })
    window.dispatchEvent(event)
    expect(event.defaultPrevented).toBe(true)
    expect(wrapper.get('[role="alert"]').text()).toContain('offline')
    confirm.mockRestore()
    wrapper.unmount()
  })

  it('does not replace freshly saved profiles when a config reload finishes', async () => {
    const wrapper = render()
    await flushPromises()
    const pending = deferred<hall.ProviderHallConfig>()
    vi.mocked(hall.getConfig).mockReturnValueOnce(pending.promise)
    await wrapper.getComponent(ProviderHallConfigForm).get('button[type="button"]').trigger('click')
    const profile: hall.ProviderHallProfile = { id: 7, version: 4, updated_at: config.updated_at, updated_by: 1,
      model: 'gpt-test', protocol: 'responses', supports_tools: false, output_limit: 256, model_aliases: [],
      reference_input_price: null, reference_cache_price: null, reference_cache_rate: null, reference_confirmed_at: null }
    wrapper.getComponent({ name: 'ProviderHallProfiles' }).vm.$emit('saved', profile)
    pending.resolve({ ...config })
    await flushPromises()
    expect(wrapper.getComponent(ProviderHallConfigForm).props('profiles')).toEqual([profile])
    expect(hall.listProfiles).toHaveBeenCalledTimes(1)
    wrapper.unmount()
  })
})
