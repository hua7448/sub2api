/**
 * Provider hall value formatting. Prices and multipliers arrive as decimal
 * strings and are re-rendered with string arithmetic only; the frontend never
 * multiplies or divides money.
 */
import type { ProviderHallMetric, ProviderHallMetricState } from '@/types/providerHall'

/** Milliseconds → "8,733 ms"; null/invalid → "-". */
export function formatMs(value: number | null | undefined, locale?: string): string {
  if (value == null || !Number.isFinite(value)) return '-'
  return `${Math.round(value).toLocaleString(locale || undefined)} ms`
}

/** 0..1 ratio → "99.1%"; null/invalid → "-". */
export function formatRate(value: number | null | undefined, digits = 1): string {
  if (value == null || !Number.isFinite(value)) return '-'
  const pct = value * 100
  const rounded = Math.round(pct * 10 ** digits) / 10 ** digits
  return `${trimTrailingZeros(rounded.toFixed(digits))}%`
}

/** Tokens per second → "27.59/s". */
export function formatTps(value: number | null | undefined): string {
  if (value == null || !Number.isFinite(value)) return '-'
  return `${trimTrailingZeros(value.toFixed(2))}/s`
}

/**
 * Decimal string → display string with at most `maxScale` fractional digits.
 * Pure string manipulation: no float parsing so 0.1 + 0.2 style drift cannot
 * appear. Non-decimal input is returned unchanged.
 */
export function formatDecimal(value: string | null | undefined, maxScale = 4, minScale = 0): string {
  if (value == null) return '-'
  const trimmed = String(value).trim()
  if (!/^-?(0|[1-9][0-9]*)(\.[0-9]+)?$/.test(trimmed)) return trimmed || '-'
  const negative = trimmed.startsWith('-')
  const [intPart, fracRaw = ''] = (negative ? trimmed.slice(1) : trimmed).split('.')
  let frac = fracRaw.slice(0, maxScale).replace(/0+$/, '')
  while (frac.length < minScale) frac += '0'
  const grouped = intPart.replace(/\B(?=(\d{3})+(?!\d))/g, ',')
  return `${negative ? '-' : ''}${grouped}${frac ? `.${frac}` : ''}`
}

/** Price per million tokens: "0.35" with unit suffix supplied by caller. */
export function formatPrice(value: string | null | undefined): string {
  return formatDecimal(value, 4, 2)
}

/** Multiplier "0.22" → "0.22x". */
export function formatMultiplier(value: string | null | undefined): string {
  const text = formatDecimal(value, 4, 0)
  return text === '-' ? '-' : `${text}x`
}

export function formatPriceUnit(unit: string | null | undefined): string {
  if (unit === 'quota_per_million') return '/M'
  if (unit === 'usd_per_million') return '$/M'
  return ''
}

/** Whether a metric carries a renderable value (ok/stale/incomplete-with-value). */
export function metricHasValue<T>(metric: ProviderHallMetric<T> | null | undefined): metric is ProviderHallMetric<T> & { value: T } {
  if (!metric) return false
  if (metric.value == null) return false
  return metric.state === 'ok' || metric.state === 'stale' || metric.state === 'incomplete'
}

/** True when the value is missing and a placeholder should be shown. */
export function metricIsPlaceholder<T>(metric: ProviderHallMetric<T> | null | undefined): boolean {
  return !metricHasValue(metric)
}

/**
 * i18n key (under providerHall.state) for a metric's state/reason. The
 * reason code takes precedence when a translation exists for it; callers use
 * `te` to check and fall back to the state key.
 */
export function stateLabelKeys(state: ProviderHallMetricState | string, reasonCode?: string): string[] {
  const keys: string[] = []
  if (reasonCode) keys.push(`providerHall.reason.${reasonCode}`)
  keys.push(`providerHall.state.${state}`)
  return keys
}

/** TTFT speed class thresholds shared by bar + tooltip. */
export type TtftSpeed = 'fast' | 'medium' | 'slow'
export function ttftSpeed(ms: number | null | undefined): TtftSpeed | null {
  if (ms == null || !Number.isFinite(ms)) return null
  if (ms <= 15000) return 'fast'
  if (ms <= 30000) return 'medium'
  return 'slow'
}

/** Success-rate score for the ring: clamp((rate-0.5)/0.5, 0, 1). */
export function successScore(rate: number | null | undefined): number {
  if (rate == null || !Number.isFinite(rate)) return 0
  return Math.min(1, Math.max(0, (rate - 0.5) / 0.5))
}

export function formatDateTime(value: string | null | undefined, locale?: string): string {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return new Intl.DateTimeFormat(locale || undefined, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
  }).format(date)
}

export function formatShortTime(value: string | null | undefined, locale?: string): string {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return new Intl.DateTimeFormat(locale || undefined, {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  }).format(date)
}

function trimTrailingZeros(text: string): string {
  if (!text.includes('.')) return text
  return text.replace(/\.?0+$/, '')
}
