import { mount } from '@vue/test-utils'
import { defineComponent } from 'vue'
import { describe, expect, it, vi } from 'vitest'
import AccountTableFilters from '../AccountTableFilters.vue'

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({ t: (key: string) => key })
  }
})

describe('AccountTableFilters', () => {
  it('includes K12 in the plan type options', () => {
    const SelectStub = defineComponent({
      name: 'SelectStub',
      props: ['modelValue', 'options'],
      template: '<div class="select-stub" />'
    })
    const wrapper = mount(AccountTableFilters, {
      props: {
        searchQuery: '',
        filters: {}
      },
      global: {
        stubs: {
          SearchInput: true,
          Select: SelectStub
        }
      }
    })

    const planSelect = wrapper.findAllComponents(SelectStub)[2]
    expect(planSelect.props('options')).toContainEqual({ value: 'k12', label: 'K12' })
  })
})
