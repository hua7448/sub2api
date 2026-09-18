/**
 * Pure geometry helpers for the provider hall SVG marks. Kept out of the
 * components so they can be unit-tested without mounting.
 */
import type { ProviderHallSparkPoint } from '@/types/providerHall'

export const RING_RADIUS = 18
export const RING_CIRCUMFERENCE = 2 * Math.PI * RING_RADIUS

/** Arc length for a 0..1 ratio, clamped. */
export function ringArc(ratio: number | null | undefined): number {
  if (ratio == null || !Number.isFinite(ratio)) return 0
  const clamped = Math.min(1, Math.max(0, ratio))
  return Math.round(clamped * RING_CIRCUMFERENCE * 1000) / 1000
}

export type RingTone = 'good' | 'series' | 'warning' | 'critical' | 'muted'

/**
 * Success-rate tone. >95% bright green (status good), >=85% series blue,
 * >=70% warning, below critical. Rates are 0..1.
 */
export function successTone(rate: number | null | undefined): RingTone {
  if (rate == null || !Number.isFinite(rate)) return 'muted'
  if (rate > 0.95) return 'good'
  if (rate >= 0.85) return 'series'
  if (rate >= 0.7) return 'warning'
  return 'critical'
}

export interface SparklineGeometry {
  real: string
  probe: string
  gaps: Array<{ x: number; width: number }>
  hasReal: boolean
  hasProbe: boolean
}

function scaleY(value: number, min: number, max: number, height: number, pad: number): number {
  if (max <= min) return height / 2
  const inner = height - pad * 2
  return pad + (1 - (value - min) / (max - min)) * inner
}

/**
 * Build the two polylines (null breaks the line) and merged gap rectangles.
 * Both series share one y-scale because both are milliseconds.
 */
export function buildSparkline(points: ProviderHallSparkPoint[], width = 200, height = 40, pad = 3): SparklineGeometry {
  const n = points.length
  const out: SparklineGeometry = { real: '', probe: '', gaps: [], hasReal: false, hasProbe: false }
  if (n === 0) return out
  const values: number[] = []
  for (const p of points) {
    if (p.real_ttft_ms != null && Number.isFinite(p.real_ttft_ms)) values.push(p.real_ttft_ms)
    if (p.probe_ms != null && Number.isFinite(p.probe_ms)) values.push(p.probe_ms)
  }
  const min = values.length ? Math.min(...values) : 0
  const max = values.length ? Math.max(...values) : 0
  const step = n > 1 ? width / (n - 1) : 0
  const xAt = (i: number) => (n > 1 ? i * step : width / 2)

  const line = (pick: (p: ProviderHallSparkPoint) => number | null) => {
    let d = ''
    let open = false
    let any = false
    points.forEach((p, i) => {
      const v = pick(p)
      if (v == null || !Number.isFinite(v)) {
        open = false
        return
      }
      any = true
      const x = xAt(i).toFixed(2)
      const y = scaleY(v, min, max, height, pad).toFixed(2)
      d += `${open ? 'L' : 'M'}${x},${y} `
      open = true
    })
    return { d: d.trim(), any }
  }
  const real = line((p) => p.real_ttft_ms)
  const probe = line((p) => p.probe_ms)
  out.real = real.d
  out.probe = probe.d
  out.hasReal = real.any
  out.hasProbe = probe.any

  // Gap rectangles: each gap bucket spans half a step either side; merge runs.
  let runStart = -1
  const flush = (endIndex: number) => {
    if (runStart < 0) return
    const x0 = Math.max(0, xAt(runStart) - step / 2)
    const x1 = Math.min(width, xAt(endIndex) + step / 2)
    out.gaps.push({ x: Math.round(x0 * 100) / 100, width: Math.round(Math.max(1, x1 - x0) * 100) / 100 })
    runStart = -1
  }
  points.forEach((p, i) => {
    if (p.gap) {
      if (runStart < 0) runStart = i
    } else {
      flush(i - 1)
    }
  })
  flush(n - 1)
  return out
}
