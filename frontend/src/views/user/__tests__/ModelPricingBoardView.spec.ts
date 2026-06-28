import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { defineComponent, h } from 'vue'

import ModelPricingBoardView from '../ModelPricingBoardView.vue'

const { getBoardMock, showErrorMock } = vi.hoisted(() => ({
  getBoardMock: vi.fn(),
  showErrorMock: vi.fn(),
}))

vi.mock('@/api/modelPricing', () => ({
  __esModule: true,
  default: {
    getBoard: getBoardMock,
  },
  getBoard: getBoardMock,
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    showError: showErrorMock,
  }),
}))

vi.mock('@/utils/apiError', () => ({
  extractApiErrorMessage: (_error: unknown, fallback: string) => fallback,
}))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  const messages: Record<string, string> = {
    'common.loading': 'Loading',
    'common.refresh': 'Refresh',
    'modelPricing.searchPlaceholder': 'Search models, platforms, groups or channels...',
    'modelPricing.loadError': 'Failed to load model pricing',
    'modelPricing.unitPerMillion': '/ 1M tokens',
    'modelPricing.userRate': 'User rate',
    'modelPricing.tabs.claude': 'Claude',
    'modelPricing.tabs.codex': 'Codex',
    'modelPricing.tabs.domestic': 'Domestic',
    'modelPricing.emptyState.claude': 'No Claude discount models available',
    'modelPricing.emptyState.codex': 'No Codex discount models available',
    'modelPricing.emptyState.domestic': 'No domestic model prices available',
    'modelPricing.emptyStateDescription.claude': 'No Anthropic pricing cards match the current filters.',
    'modelPricing.emptyStateDescription.codex': 'No Codex pricing cards match the current filters.',
    'modelPricing.emptyStateDescription.domestic': 'Configure DeepSeek, GLM, Kimi, or similar models in channel pricing to show them here.',
    'modelPricing.card.sitePrice': 'Our Price',
    'modelPricing.card.officialPrice': 'Official Price',
    'modelPricing.card.priceFormula': 'Formula: official price * group rate x{rate}, shown at 1:1 USD parity',
    'modelPricing.card.savings': 'Save {percent}',
    'modelPricing.card.group': 'Group',
    'modelPricing.card.channel': 'Channel',
    'modelPricing.card.rateLabel': 'Rate',
    'modelPricing.card.rate': 'x{rate} rate',
    'modelPricing.card.source': 'Source',
    'modelPricing.card.input': 'Input',
    'modelPricing.card.output': 'Output',
    'modelPricing.card.cacheRead': 'Cache Read',
    'monitorCommon.providers.anthropic': 'Anthropic',
    'monitorCommon.providers.openai': 'OpenAI',
  }

  return {
    ...actual,
    useI18n: () => ({
      t: (key: string, params?: Record<string, string>) =>
        (messages[key] ?? key).replace(/\{(\w+)\}/g, (_, token) => params?.[token] ?? `{${token}}`),
    }),
  }
})

const AppLayoutStub = { template: '<div><slot /></div>' }
const IconStub = defineComponent({
  props: {
    name: {
      type: String,
      default: '',
    },
  },
  setup(props) {
    return () => h('span', { 'data-icon': props.name })
  },
})
const ModelIconStub = defineComponent({
  props: {
    model: {
      type: String,
      default: '',
    },
  },
  setup(props) {
    return () => h('span', { 'data-model-icon': props.model })
  },
})

describe('ModelPricingBoardView', () => {
  beforeEach(() => {
    getBoardMock.mockReset()
    showErrorMock.mockReset()
  })

  it('filters cards by Claude and Codex tabs and searches within the active tab', async () => {
    getBoardMock.mockResolvedValue({
      items: [
        {
          model_id: 'claude-sonnet-4',
          platform: 'anthropic',
          group_id: 1,
          group_name: 'Anthropic Pro',
          channel_id: 10,
          channel_name: 'Claude Alpha',
          rate_multiplier: 1,
          user_rate_overridden: false,
          site_input_price: 0.000003,
          site_output_price: 0.000015,
          site_cache_read_price: null,
          official_input_price: 0.000004,
          official_output_price: 0.00002,
          official_cache_read_price: null,
          input_savings: 0.25,
          output_savings: 0.25,
          cache_read_savings: null,
        },
        {
          model_id: 'codex-mini-latest',
          platform: 'openai',
          group_id: 2,
          group_name: 'Codex Internal',
          channel_id: 20,
          channel_name: 'Codex Beta',
          rate_multiplier: 1.2,
          user_rate_overridden: true,
          site_input_price: 0.000002,
          site_output_price: null,
          site_cache_read_price: 0.0000005,
          official_input_price: 0.000003,
          official_output_price: null,
          official_cache_read_price: null,
          input_savings: 0.333,
          output_savings: null,
          cache_read_savings: null,
        },
        {
          model_id: 'gpt-5.5',
          platform: 'openai',
          group_id: 3,
          group_name: 'Codex GPT',
          channel_id: 30,
          channel_name: 'OpenAI Codex',
          rate_multiplier: 1,
          user_rate_overridden: false,
          site_input_price: 0.000001,
          site_output_price: 0.000002,
          site_cache_read_price: null,
          official_input_price: 0.0000012,
          official_output_price: 0.0000024,
          official_cache_read_price: null,
          input_savings: 0.166,
          output_savings: 0.166,
          cache_read_savings: null,
        },
        {
          model_id: 'kimi-k2.7-code',
          platform: 'kimi',
          group_id: 4,
          group_name: 'Domestic Anthropic',
          channel_id: 40,
          channel_name: 'Kimi Channel',
          rate_multiplier: 1,
          user_rate_overridden: false,
          site_input_price: 0.00000095,
          site_output_price: 0.000004,
          site_cache_read_price: 0.00000019,
          official_input_price: 0.00000095,
          official_output_price: 0.000004,
          official_cache_read_price: 0.00000019,
          input_savings: 0,
          output_savings: 0,
          cache_read_savings: 0,
        },
      ],
    })

    const wrapper = mount(ModelPricingBoardView, {
      global: {
        stubs: {
          AppLayout: AppLayoutStub,
          Icon: IconStub,
          ModelIcon: ModelIconStub,
        },
      },
    })

    await flushPromises()

    const tabs = wrapper.findAll('[data-testid^="model-pricing-tab-"]')
    expect(tabs.map((tab) => tab.text())).toEqual(['Codex', 'Claude', 'Domestic'])

    expect(wrapper.text()).toContain('codex-mini-latest')
    expect(wrapper.text()).toContain('gpt-5.5')
    expect(wrapper.text()).not.toContain('claude-sonnet-4')
    expect(wrapper.text()).not.toContain('kimi-k2.7-code')
    expect(wrapper.text()).toContain('Formula: official price * group rate x1.2, shown at 1:1 USD parity')

    await wrapper.get('input[type="text"]').setValue('internal')
    await flushPromises()
    expect(wrapper.text()).toContain('codex-mini-latest')
    expect(wrapper.text()).not.toContain('gpt-5.5')

    await wrapper.get('input[type="text"]').setValue('claude alpha')
    await flushPromises()
    expect(wrapper.text()).toContain('No Codex discount models available')
    expect(wrapper.text()).not.toContain('claude-sonnet-4')

    await wrapper.get('input[type="text"]').setValue('')
    await wrapper.get('[data-testid="model-pricing-tab-claude"]').trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('claude-sonnet-4')
    expect(wrapper.text()).not.toContain('codex-mini-latest')

    await wrapper.get('[data-testid="model-pricing-tab-domestic"]').trigger('click')
    await flushPromises()
    expect(wrapper.text()).toContain('kimi-k2.7-code')
    expect(wrapper.text()).toContain('Kimi')
    expect(wrapper.text()).not.toContain('claude-sonnet-4')
  })

  it('hides savings badge when official price is missing for a pricing row', async () => {
    getBoardMock.mockResolvedValue({
      items: [
        {
          model_id: 'codex-lite',
          platform: 'openai',
          group_id: 2,
          group_name: 'Codex Internal',
          channel_id: 20,
          channel_name: 'Codex Beta',
          rate_multiplier: 1.2,
          user_rate_overridden: false,
          site_input_price: 0.000002,
          site_output_price: null,
          site_cache_read_price: 0.0000005,
          official_input_price: 0.000003,
          official_output_price: null,
          official_cache_read_price: null,
          input_savings: 0.333,
          output_savings: null,
          cache_read_savings: 0.5,
        },
      ],
    })

    const wrapper = mount(ModelPricingBoardView, {
      global: {
        stubs: {
          AppLayout: AppLayoutStub,
          Icon: IconStub,
          ModelIcon: ModelIconStub,
        },
      },
    })

    await flushPromises()
    await wrapper.get('[data-testid="model-pricing-tab-codex"]').trigger('click')
    await flushPromises()

    const text = wrapper.text()
    expect(text).toContain('Save 33.3%')
    expect(text).not.toContain('Save 50%')
    expect(text).not.toContain('Official Price$0.5 / 1M tokens')
  })
})
