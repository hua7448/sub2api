/**
 * Provider hall shared types. Mirrors backend/internal/handler/dto/provider_hall.go
 * (ProviderHallMetric) and dto/provider_hall_user.go (user-facing contract).
 * JSON field names are the contract; do not rename.
 */

export type ProviderHallProtocol = 'responses' | 'chat_completions' | 'messages'
export type ProviderHallRange = '6h' | '24h' | '7d' | '30d'

export type ProviderHallMetricState =
  | 'ok'
  | 'insufficient'
  | 'incomplete'
  | 'stale'
  | 'disabled'
  | 'not_applicable'

export interface ProviderHallMetric<T> {
  value: T | null
  state: ProviderHallMetricState
  reason_code: string
  window_start: string | null
  window_end: string | null
  computed_at: string | null
}

export interface ProviderHallProfileRef {
  profile_id: number
  model: string
  protocol: ProviderHallProtocol
}

export interface ProviderHallQuote {
  input_price: string
  cache_price: string
  unit: 'usd_per_million' | 'quota_per_million' | string
  applicable: boolean
  reason_code: string
}

export type ProviderHallHealthStatus = 'up' | 'down' | 'unknown'

export interface ProviderHallModelHealth {
  profile_id: number
  model: string
  protocol: ProviderHallProtocol
  status: ProviderHallHealthStatus
  checked_at: string | null
}

export interface ProviderHallHealth {
  status: ProviderHallHealthStatus
  checked_at: string | null
  models: ProviderHallModelHealth[]
}

export interface ProviderHallRowMetrics {
  ttft_fast95_ms: ProviderHallMetric<number>
  ttft_p90_ms: ProviderHallMetric<number>
  cache_rate: ProviderHallMetric<number>
  success_rate: ProviderHallMetric<number>
}

export interface ProviderHallSparkPoint {
  t: string
  real_ttft_ms: number | null
  probe_ms: number | null
  gap: boolean
}

export interface ProviderHallSparkline {
  range: ProviderHallRange | string
  bucket_seconds: number
  points: ProviderHallSparkPoint[]
}

export type ProviderHallVerdict = 'passed' | 'failed' | 'suspected' | 'insufficient'

export interface ProviderHallVerificationBadge {
  verdict: ProviderHallVerdict | null
  reason_code: string
  completed_at: string | null
  expired: boolean
  report_id: number | null
}

export interface ProviderHallRow {
  group_id: number
  name: string
  description: string
  display_order: number
  rate: string
  quote: ProviderHallQuote
  historical_price: ProviderHallMetric<string>
  predicted_rate: ProviderHallMetric<string>
  health: ProviderHallHealth
  metrics: ProviderHallRowMetrics
  sparkline: ProviderHallSparkline
  verification: ProviderHallVerificationBadge
  default_profile: ProviderHallProfileRef
}

export interface ProviderHallCatalog {
  models: ProviderHallProfileRef[]
  default: ProviderHallProfileRef
}

export interface ProviderHallSummary {
  available: number
  listed: number
  abnormal: number
  verified: number
}

export interface ProviderHallPagination {
  page: number
  page_size: number
  total: number
}

export interface ProviderHallListResponse {
  snapshot_id: string
  data_through: string | null
  metric_version: number
  pricing_at: string
  range: ProviderHallRange | string
  catalog: ProviderHallCatalog
  summary: ProviderHallSummary
  items: ProviderHallRow[]
  pagination: ProviderHallPagination
}

export interface ProviderHallProbeTokens {
  input: number
  output: number
  at: string
}

export interface ProviderHallDetailMetrics {
  e2e_availability_6h: ProviderHallMetric<number>
  tps: ProviderHallMetric<number>
  probe_ttft_ms: ProviderHallMetric<number>
  probe_total_ms: ProviderHallMetric<number>
  p90_ms: ProviderHallMetric<number>
  probe_tokens: ProviderHallProbeTokens | null
  last_probe_at: string | null
}

export interface ProviderHallTrendPoint {
  t: string
  real_fast95_ms: number | null
  probe_total_ms: number | null
  probe_ttft_ms: number | null
  gap: boolean
  probe_failed: number
}

export interface ProviderHallTrend {
  range: ProviderHallRange | string
  bucket_seconds: number
  points: ProviderHallTrendPoint[]
}

export interface ProviderHallProfileStatus extends ProviderHallProfileRef {
  health: ProviderHallModelHealth
  verification: ProviderHallVerificationBadge
}

export interface ProviderHallDetailResponse extends ProviderHallRow {
  detail: ProviderHallDetailMetrics
  trend: ProviderHallTrend
  profiles: ProviderHallProfileStatus[]
}

export interface ProviderHallVerificationReport {
  report_id: number
  profile_id: number
  model: string
  protocol: ProviderHallProtocol
  verdict: ProviderHallVerdict | string
  execution_status: 'completed' | 'partial' | 'error' | string
  reason_code: string
  completed_at: string
  expires_at: string
  expired: boolean
  summary: ProviderHallVerificationSummary | Record<string, unknown> | null
}

export interface ProviderHallSuiteSummary {
  passed: number
  failed: number
  error: number
}

export interface ProviderHallVerificationSummary {
  arithmetic?: ProviderHallSuiteSummary
  json?: ProviderHallSuiteSummary
  tool?: ProviderHallSuiteSummary | null
  model?: { matched: number; mismatched: number; missing: number; seen: string[] }
}

export interface ProviderHallVerificationListResponse {
  items: ProviderHallVerificationReport[]
  pagination: ProviderHallPagination
}

export const PROVIDER_HALL_SORT_FIELDS = [
  'display_order',
  'rate',
  'historical_price',
  'predicted_rate',
  'ttft_fast95',
  'cache_rate',
  'success_rate',
] as const

export type ProviderHallSortField = (typeof PROVIDER_HALL_SORT_FIELDS)[number]
export type ProviderHallSortDirection = 'asc' | 'desc'

export interface ProviderHallSortRule {
  field: ProviderHallSortField
  direction: ProviderHallSortDirection
}

export const PROVIDER_HALL_RANGES: readonly ProviderHallRange[] = ['6h', '24h', '7d', '30d']
