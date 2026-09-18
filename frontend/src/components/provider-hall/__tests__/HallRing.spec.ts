import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import HallRing from '../HallRing.vue'
import { RING_CIRCUMFERENCE, ringArc, successTone } from '../hallGeometry'
import { createHallI18n, metric } from './testUtils'

describe('HallRing geometry', () => {
  it('maps ratios to arc lengths and clamps outside 0..1', () => {
    expect(ringArc(0)).toBe(0)
    expect(ringArc(0.5)).toBeCloseTo(RING_CIRCUMFERENCE / 2, 2)
    expect(ringArc(1)).toBeCloseTo(RING_CIRCUMFERENCE, 2)
    expect(ringArc(1.7)).toBeCloseTo(RING_CIRCUMFERENCE, 2)
    expect(ringArc(-1)).toBe(0)
    expect(ringArc(null)).toBe(0)
  })

  it('colours the success ring by score thresholds', () => {
    expect(successTone(0.99)).toBe('good')
    expect(successTone(0.951)).toBe('good')
    expect(successTone(0.95)).toBe('series')
    expect(successTone(0.85)).toBe('series')
    expect(successTone(0.849)).toBe('warning')
    expect(successTone(0.7)).toBe('warning')
    expect(successTone(0.41)).toBe('critical')
    expect(successTone(null)).toBe('muted')
  })
})

describe('HallRing component', () => {
  const mountRing = (props: Record<string, unknown>) =>
    mount(HallRing, { props: { label: '1h', ...props } as never, global: { plugins: [createHallI18n()] } })

  it('renders the arc and value for an ok metric', () => {
    const wrapper = mountRing({ metric: metric(0.991), mode: 'score' })
    const arc = wrapper.get('.hall-ring-arc')
    expect(arc.attributes('data-tone')).toBe('good')
    expect(arc.attributes('stroke-dasharray')?.startsWith(ringArc(0.991).toString())).toBe(true)
    expect(wrapper.text()).toContain('99.1%')
    expect(wrapper.find('[data-testid="ring-placeholder"]').exists()).toBe(false)
  })

  it('uses the series hue in plain mode regardless of the value', () => {
    const wrapper = mountRing({ metric: metric(0.41), mode: 'plain' })
    expect(wrapper.get('.hall-ring-arc').attributes('data-tone')).toBe('series')
  })

  it('shows a placeholder and no arc when samples are insufficient', () => {
    const wrapper = mountRing({ metric: metric<number>(null, 'insufficient', 'samples_below_20'), mode: 'score' })
    expect(wrapper.find('.hall-ring-arc').exists()).toBe(false)
    expect(wrapper.get('[data-testid="ring-placeholder"]').text()).toBe('样本不足')
    expect(wrapper.get('text').text()).toBe('-')
  })

  it('shows a value with a stale marker for stale metrics', () => {
    const wrapper = mountRing({ metric: metric(0.8, 'stale'), mode: 'score' })
    expect(wrapper.get('.hall-ring-arc').attributes('data-tone')).toBe('warning')
    expect(wrapper.get('[data-testid="ring-stale"]').text()).toBe('数据过期')
  })

  it('falls back to the state label when the reason code has no translation', () => {
    const wrapper = mountRing({ metric: metric<number>(null, 'incomplete', 'something_new') })
    expect(wrapper.get('[data-testid="ring-placeholder"]').text()).toBe('数据不完整')
  })
})
