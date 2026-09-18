import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import UseGroupDialog from '../UseGroupDialog.vue'
import { createHallI18n, makeRow } from './testUtils'
import { keysAPI } from '@/api/keys'
import type { ApiKey } from '@/types'

const showError = vi.fn()
const showSuccess = vi.fn()
vi.mock('@/stores/app', () => ({ useAppStore: () => ({ showError, showSuccess }) }))
vi.mock('@/api/keys', () => ({ keysAPI: { list: vi.fn(), update: vi.fn(), create: vi.fn() } }))

const key = (id: number, groupId: number | null): ApiKey =>
  ({ id, name: `key-${id}`, group_id: groupId, group: groupId ? ({ id: groupId, name: `G${groupId}` } as ApiKey['group']) : undefined } as ApiKey)

function mountDialog(group = makeRow({ group_id: 9, name: 'Nine' })) {
  return mount(UseGroupDialog, {
    props: { show: true, group },
    global: {
      plugins: [createHallI18n()],
      stubs: { BaseDialog: { props: ['show'], template: '<div v-if="show"><slot /></div>' }, LoadingSpinner: true },
    },
  })
}

beforeEach(() => {
  vi.clearAllMocks()
  vi.mocked(keysAPI.list).mockResolvedValue({ items: [key(1, 3), key(2, 9)], total: 2, page: 1, page_size: 20, pages: 1 })
})

describe('UseGroupDialog', () => {
  it('switches the selected key to the group with only group_id', async () => {
    vi.mocked(keysAPI.update).mockResolvedValue(key(1, 9))
    const wrapper = mountDialog()
    await flushPromises()
    expect(keysAPI.list).toHaveBeenCalledWith(1, 20)
    const radios = wrapper.findAll('input[type="radio"]')
    expect(radios[1].attributes('disabled')).toBeDefined() // already in group 9
    expect(wrapper.text()).toContain('G3 → Nine')
    await radios[0].setValue(true)
    await wrapper.get('form').trigger('submit')
    await flushPromises()
    expect(keysAPI.update).toHaveBeenCalledWith(1, { group_id: 9 })
    expect(showSuccess).toHaveBeenCalled()
    expect(wrapper.emitted('done')?.[0][0]).toMatchObject({ action: 'switch' })
    expect(wrapper.emitted('close')).toHaveLength(1)
  })

  it('creates a key bound to the group and keeps the name when the request fails', async () => {
    vi.mocked(keysAPI.create).mockRejectedValueOnce({ message: 'quota exceeded' }).mockResolvedValueOnce(key(5, 9))
    const wrapper = mountDialog()
    await flushPromises()
    await wrapper.get('[data-testid="tab-create"]').trigger('click')
    const input = wrapper.get('[data-testid="create-name"]')
    await input.setValue('  my key ')
    await wrapper.get('form').trigger('submit')
    await flushPromises()
    expect(keysAPI.create).toHaveBeenCalledWith('my key', 9)
    expect(showError).toHaveBeenCalledWith('quota exceeded')
    expect((input.element as HTMLInputElement).value).toBe('  my key ')
    expect(wrapper.emitted('close')).toBeUndefined()
    await wrapper.get('form').trigger('submit')
    await flushPromises()
    expect(keysAPI.create).toHaveBeenCalledTimes(2)
    expect(wrapper.emitted('done')?.[0][0]).toMatchObject({ action: 'create' })
    expect(wrapper.emitted('close')).toHaveLength(1)
  })

  it('keeps the selection and reports the error when switching fails', async () => {
    vi.mocked(keysAPI.update).mockRejectedValue({ message: 'forbidden' })
    const wrapper = mountDialog()
    await flushPromises()
    const radio = wrapper.findAll('input[type="radio"]')[0]
    await radio.setValue(true)
    await wrapper.get('form').trigger('submit')
    await flushPromises()
    expect(showError).toHaveBeenCalledWith('forbidden')
    expect((radio.element as HTMLInputElement).checked).toBe(true)
    expect(wrapper.emitted('close')).toBeUndefined()
  })

  it('shows an empty hint when the user has no keys', async () => {
    vi.mocked(keysAPI.list).mockResolvedValue({ items: [], total: 0, page: 1, page_size: 20, pages: 0 })
    const wrapper = mountDialog()
    await flushPromises()
    expect(wrapper.find('[data-testid="no-keys"]').exists()).toBe(true)
    expect(wrapper.get('[data-testid="switch-submit"]').attributes('disabled')).toBeDefined()
  })
})
