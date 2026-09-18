import { createI18n } from 'vue-i18n'
import zh from '@/i18n/locales/zh/providerHall'
import type {
  ProviderHallDetailResponse,
  ProviderHallListResponse,
  ProviderHallMetric,
  ProviderHallRow,
} from '@/types/providerHall'

type MessageTree = { [key: string]: string | MessageTree }

/**
 * The test build uses the vue-i18n runtime bundle (no message compiler), so
 * every string leaf is wrapped into a message function that interpolates
 * `{named}` placeholders by hand.
 */
export function toRuntimeMessages(tree: MessageTree): Record<string, unknown> {
  const out: Record<string, unknown> = {}
  for (const [key, value] of Object.entries(tree)) {
    if (typeof value === 'string') {
      out[key] = (ctx: { named: (name: string) => unknown }) =>
        value.replace(/\{(\w+)\}/g, (_, name: string) => String(ctx.named(name) ?? ''))
    } else {
      out[key] = toRuntimeMessages(value)
    }
  }
  return out
}

export function createHallI18n() {
  return createI18n({
    legacy: false,
    locale: 'zh',
    missingWarn: false,
    fallbackWarn: false,
    messages: { zh: toRuntimeMessages(zh as unknown as MessageTree) },
  })
}

export function metric<T>(value: T | null, state: ProviderHallMetric<T>['state'] = 'ok', reason = ''): ProviderHallMetric<T> {
  return {
    value,
    state,
    reason_code: reason,
    window_start: '2026-09-12T01:00:00Z',
    window_end: '2026-09-12T02:00:00Z',
    computed_at: '2026-09-12T02:00:45Z',
  }
}

export function makeRow(overrides: Partial<ProviderHallRow> = {}): ProviderHallRow {
  const id = overrides.group_id ?? 1
  return {
    group_id: id,
    name: `Group ${id}`,
    description: 'Public route',
    display_order: id,
    rate: '0.22',
    quote: { input_price: '1.1', cache_price: '0.11', unit: 'usd_per_million', applicable: true, reason_code: '' },
    historical_price: metric('0.35'),
    predicted_rate: metric('0.33'),
    health: {
      status: 'up',
      checked_at: '2026-09-12T02:00:00Z',
      models: [
        { profile_id: 1, model: 'astra', protocol: 'responses', status: 'up', checked_at: '2026-09-12T02:00:00Z' },
        { profile_id: 2, model: 'terra', protocol: 'chat_completions', status: 'down', checked_at: null },
      ],
    },
    metrics: {
      ttft_fast95_ms: metric(8733),
      ttft_p90_ms: metric(14330),
      cache_rate: metric(0.774),
      success_rate: metric(0.991),
    },
    sparkline: {
      range: '6h',
      bucket_seconds: 300,
      points: [
        { t: '2026-09-12T00:00:00Z', real_ttft_ms: 9000, probe_ms: 2000, gap: false },
        { t: '2026-09-12T00:05:00Z', real_ttft_ms: null, probe_ms: null, gap: true },
        { t: '2026-09-12T00:10:00Z', real_ttft_ms: 8000, probe_ms: 2100, gap: false },
      ],
    },
    verification: { verdict: 'passed', reason_code: '', completed_at: '2026-09-12T01:30:00Z', expired: false, report_id: 77 },
    default_profile: { profile_id: 1, model: 'astra', protocol: 'responses' },
    ...overrides,
  }
}

export function makeList(rows: ProviderHallRow[], overrides: Partial<ProviderHallListResponse> = {}): ProviderHallListResponse {
  return {
    snapshot_id: '2026-09-12T02:00:00Z',
    data_through: '2026-09-12T02:00:00Z',
    metric_version: 1,
    pricing_at: '2026-09-12T02:01:00Z',
    range: '6h',
    catalog: {
      models: [
        { profile_id: 1, model: 'astra', protocol: 'responses' },
        { profile_id: 2, model: 'terra', protocol: 'chat_completions' },
      ],
      default: { profile_id: 1, model: 'astra', protocol: 'responses' },
    },
    summary: { available: rows.filter((r) => r.health.status === 'up').length, listed: rows.length, abnormal: rows.filter((r) => r.health.status === 'down').length, verified: rows.filter((r) => r.verification.verdict === 'passed' && !r.verification.expired).length },
    items: rows,
    pagination: { page: 1, page_size: 50, total: rows.length },
    ...overrides,
  }
}

export function makeDetail(row: ProviderHallRow): ProviderHallDetailResponse {
  return {
    ...row,
    detail: {
      e2e_availability_6h: metric(0.88),
      tps: metric(27.59),
      probe_ttft_ms: metric(2058),
      probe_total_ms: metric(2064),
      p90_ms: metric(14330),
      probe_tokens: { input: 76, output: 5, at: '2026-09-12T01:55:53Z' },
      last_probe_at: '2026-09-12T01:55:53Z',
    },
    trend: {
      range: '6h',
      bucket_seconds: 300,
      points: [
        { t: '2026-09-12T00:00:00Z', real_fast95_ms: 9000, probe_total_ms: 2000, probe_ttft_ms: 1800, gap: false, probe_failed: 0 },
        { t: '2026-09-12T00:05:00Z', real_fast95_ms: null, probe_total_ms: null, probe_ttft_ms: null, gap: true, probe_failed: 1 },
      ],
    },
    profiles: row.health.models.map((m) => ({
      profile_id: m.profile_id,
      model: m.model,
      protocol: m.protocol,
      health: m,
      verification: row.verification,
    })),
  }
}
