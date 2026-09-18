import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import HallTable from '../HallTable.vue'
import { createHallI18n, makeDetail, makeRow, metric } from './testUtils'

vi.mock('vue-chartjs', () => ({ Line: { name: 'Line', template: '<canvas data-testid="line-chart"></canvas>' } }))

function mountTable(rows = [makeRow({ group_id: 1 }), makeRow({ group_id: 2 })], extra: Record<string, unknown> = {}) {
  return mount(HallTable, {
    attachTo: document.body,
    props: {
      items: rows,
      expanded: new Set<number>(),
      details: new Map(),
      detailLoading: new Set<number>(),
      detailErrors: new Map(),
      range: '6h',
      ...extra,
    } as never,
    global: { plugins: [createHallI18n()], stubs: { Teleport: true } },
  })
}

describe('HallTable', () => {
  it('renders values only from backend state and placeholders otherwise', () => {
    const rows = [
      makeRow({ group_id: 1 }),
      makeRow({
        group_id: 2,
        historical_price: metric<string>(null, 'insufficient', 'bills_below_200'),
        predicted_rate: metric<string>(null, 'not_applicable', 'tiered_pricing'),
        metrics: {
          ttft_fast95_ms: metric<number>(null, 'incomplete', 'collection_gap'),
          ttft_p90_ms: metric<number>(null, 'insufficient'),
          cache_rate: metric(0.5, 'stale'),
          success_rate: metric<number>(null, 'disabled', 'target_disabled'),
        },
        health: { status: 'down', checked_at: null, models: [] },
        verification: { verdict: null, reason_code: '', completed_at: null, expired: false, report_id: null },
      }),
    ]
    const wrapper = mountTable(rows)
    const trs = wrapper.findAll('[data-testid="hall-row"]')
    expect(trs).toHaveLength(2)
    expect(trs[0].get('[data-testid="row-price"]').text()).toContain('0.35')
    expect(trs[0].get('[data-testid="row-rate"]').text()).toBe('0.22x')
    expect(trs[1].get('[data-testid="row-price-placeholder"]').text()).toBe('账单不足 200 条')
    expect(trs[1].get('[data-testid="row-predicted-placeholder"]').text()).toBe('阶梯定价')
    expect(trs[1].get('[data-testid="ttft-value"]').text()).toBe('采集缺口')
    expect(trs[1].find('[data-testid="ttft-bar"]').exists()).toBe(false)
    expect(trs[1].get('[data-testid="row-health"]').text()).toContain('异常')
    expect(trs[1].get('[data-testid="verification-pill"]').text()).toBe('未检测')
    expect(trs[1].find('[data-testid="verification-report-link"]').exists()).toBe(false)
    // stale cache ring still renders its value plus a stale marker
    expect(trs[1].text()).toContain('50%')
    expect(trs[1].find('[data-testid="ring-stale"]').exists()).toBe(true)
    wrapper.unmount()
  })

  it('scales the TTFT bar relative to the page maximum', () => {
    const rows = [
      makeRow({ group_id: 1, metrics: { ...makeRow().metrics, ttft_fast95_ms: metric(5000) } }),
      makeRow({ group_id: 2, metrics: { ...makeRow().metrics, ttft_fast95_ms: metric(20000) } }),
    ]
    const wrapper = mountTable(rows)
    const bars = wrapper.findAll('[data-testid="hall-row"] [data-testid="ttft-bar"]')
    expect(bars[0].attributes('style')).toContain('width: 25%')
    expect(bars[1].attributes('style')).toContain('width: 100%')
    expect(wrapper.findAll('[data-testid="hall-row"] [data-testid="ttft-speed"]')[1].attributes('data-speed')).toBe('medium')
    wrapper.unmount()
  })

  it('renders mobile cards alongside the desktop table', () => {
    const wrapper = mountTable()
    const cards = wrapper.findAll('[data-testid="hall-card"]')
    expect(cards).toHaveLength(2)
    expect(cards[0].text()).toContain('Group 1')
    expect(cards[0].find('[data-testid="use-group"]').exists()).toBe(true)
    expect(wrapper.get('[data-testid="hall-cards"]').classes()).toContain('md:hidden')
    wrapper.unmount()
  })

  it('toggles with Enter/Space and moves focus with arrow keys', async () => {
    const wrapper = mountTable()
    const trs = wrapper.findAll('[data-testid="hall-row"]')
    await trs[0].trigger('keydown', { key: 'Enter' })
    await trs[1].trigger('keydown', { key: ' ' })
    expect(wrapper.emitted('toggle')).toEqual([[1], [2]])
    ;(trs[0].element as HTMLElement).focus()
    await trs[0].trigger('keydown', { key: 'ArrowDown' })
    expect(document.activeElement).toBe(trs[1].element)
    await trs[1].trigger('keydown', { key: 'ArrowUp' })
    expect(document.activeElement).toBe(trs[0].element)
    wrapper.unmount()
  })

  it('emits use-group and view-report without toggling the row', async () => {
    const wrapper = mountTable()
    const tr = wrapper.findAll('[data-testid="hall-row"]')[0]
    await tr.get('[data-testid="use-group"]').trigger('click')
    await tr.get('[data-testid="verification-report-link"]').trigger('click')
    expect(wrapper.emitted('toggle')).toBeUndefined()
    expect(wrapper.emitted('use-group')?.[0][0]).toMatchObject({ group_id: 1 })
    expect(wrapper.emitted('view-report')?.[0][0]).toEqual({ groupId: 1, reportId: 77 })
    wrapper.unmount()
  })

  it('renders the expanded detail panel for the active range only', () => {
    const row = makeRow({ group_id: 1 })
    const details = new Map([[`1:6h`, makeDetail(row)]])
    const wrapper = mountTable([row], { expanded: new Set([1]), details })
    const detail = wrapper.get('[data-testid="hall-detail-row"]')
    expect(detail.get('[data-testid="detail-metrics"]').text()).toContain('27.59/s')
    expect(detail.get('[data-testid="detail-tokens"]').text()).toContain('76')
    expect(detail.find('[data-testid="line-chart"]').exists()).toBe(true)
    wrapper.unmount()

    const other = mountTable([row], { expanded: new Set([1]), details, range: '24h', detailLoading: new Set([1]) })
    expect(other.get('[data-testid="hall-detail-row"]').find('[data-testid="detail-loading"]').exists()).toBe(true)
    other.unmount()
  })
})
