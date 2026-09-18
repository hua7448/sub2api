import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import HallSparkline from '../HallSparkline.vue'
import { buildSparkline } from '../hallGeometry'
import { createHallI18n } from './testUtils'
import type { ProviderHallSparkPoint } from '@/types/providerHall'

const point = (i: number, real: number | null, probe: number | null, gap = false): ProviderHallSparkPoint => ({
  t: `2026-09-12T00:${String(i * 5).padStart(2, '0')}:00Z`,
  real_ttft_ms: real,
  probe_ms: probe,
  gap,
})

describe('buildSparkline', () => {
  it('merges consecutive gap buckets into one rectangle', () => {
    const points = [point(0, 1000, 500), point(1, null, null, true), point(2, null, null, true), point(3, 1200, 600), point(4, null, null, true)]
    const geometry = buildSparkline(points, 200, 40)
    expect(geometry.gaps).toHaveLength(2)
    const step = 200 / 4
    expect(geometry.gaps[0].x).toBeCloseTo(step * 1 - step / 2, 1)
    expect(geometry.gaps[0].width).toBeCloseTo(step * 2, 1)
    // Trailing gap is clipped to the right edge.
    expect(geometry.gaps[1].x + geometry.gaps[1].width).toBeCloseTo(200, 1)
  })

  it('breaks the polyline at null values instead of bridging them', () => {
    const points = [point(0, 1000, null), point(1, null, null, true), point(2, 1200, null)]
    const geometry = buildSparkline(points, 200, 40)
    expect(geometry.hasReal).toBe(true)
    expect(geometry.hasProbe).toBe(false)
    expect(geometry.real.match(/M/g)).toHaveLength(2)
    expect(geometry.real).not.toContain('L')
  })

  it('shares one y-scale between the two series', () => {
    const geometry = buildSparkline([point(0, 4000, 1000), point(1, 4000, 1000)], 100, 40, 0)
    // real is the max → y=0; probe is the min → y=height
    expect(geometry.real).toContain(',0.00')
    expect(geometry.probe).toContain(',40.00')
  })

  it('handles empty input', () => {
    const geometry = buildSparkline([], 200, 40)
    expect(geometry.gaps).toEqual([])
    expect(geometry.hasReal).toBe(false)
  })
})

describe('HallSparkline component', () => {
  it('renders one rect per merged gap and both series paths', () => {
    const points = [point(0, 1000, 500), point(1, null, null, true), point(2, 1200, 600)]
    const wrapper = mount(HallSparkline, { props: { points, name: 'A' }, global: { plugins: [createHallI18n()] } })
    expect(wrapper.findAll('[data-testid="sparkline-gap"]')).toHaveLength(1)
    expect(wrapper.find('[data-testid="sparkline-real"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="sparkline-probe"]').exists()).toBe(true)
    expect(wrapper.get('svg').attributes('aria-label')).toContain('A')
    expect(wrapper.get('svg').attributes('aria-label')).toContain('3')
  })

  it('renders a dash when there is no data', () => {
    const wrapper = mount(HallSparkline, { props: { points: [] }, global: { plugins: [createHallI18n()] } })
    expect(wrapper.find('[data-testid="sparkline-real"]').exists()).toBe(false)
    expect(wrapper.get('svg').text()).toBe('-')
  })
})
