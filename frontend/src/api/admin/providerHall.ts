import { apiClient } from '../client'

export type ProviderHallProtocol = 'responses' | 'chat_completions' | 'messages'
export type ProviderHallRange = '6h' | '24h' | '7d' | '30d'

export interface ProviderHallVersion {
  version: number
  updated_at: string
  updated_by: number | null
}

export interface ProviderHallReadiness {
  collection: boolean
  tasks: boolean
  display: boolean
}

export interface ProviderHallConfig extends ProviderHallVersion {
  /** Which runtime switches this deployment can turn on (absent on old backends). */
  readiness?: ProviderHallReadiness
  collection_enabled: boolean
  display_enabled: boolean
  tasks_enabled: boolean
  default_model: string
  default_protocol: ProviderHallProtocol
  default_range: ProviderHallRange
  gateway_origin: string
  operator_user_id: number | null
  daily_budget: string
  expected_nodes: string[]
}

export interface ProviderHallGroup extends ProviderHallVersion {
  group_id: number
  listed: boolean
  display_name: string
  description: string
  display_order: number
}

export interface ProviderHallProfile extends ProviderHallVersion {
  id: number
  model: string
  protocol: ProviderHallProtocol
  supports_tools: boolean
  output_limit: number
  model_aliases: string[]
  reference_input_price: string | null
  reference_cache_price: string | null
  reference_cache_rate: string | null
  reference_confirmed_at: string | null
}

export type { ProviderHallMetric } from '@/types/providerHall'

export interface ProviderHallTargetInput {
  profile_id: number
  probe_key_id: number | null
  enabled: boolean
  probe_interval_seconds?: number
  verification_interval_seconds?: number
}

export interface ProviderHallTarget extends ProviderHallVersion, Required<ProviderHallTargetInput> {
  id: number
  group_id: number
}

export interface ProviderHallTargetSet extends ProviderHallVersion {
  group_id: number
  items: ProviderHallTarget[]
}

export interface ProviderHallTargetSetInput {
  version: number
  items: ProviderHallTargetInput[]
}

export type ProviderHallConfigInput = Omit<ProviderHallConfig, 'updated_at' | 'updated_by' | 'readiness'>
export type ProviderHallGroupInput = Omit<ProviderHallGroup, 'group_id' | 'updated_at' | 'updated_by'>
export type ProviderHallProfileInput = Omit<ProviderHallProfile, 'id' | 'updated_at' | 'updated_by'>

const base = '/admin/provider-hall'

export async function getConfig(signal?: AbortSignal): Promise<ProviderHallConfig> {
  return (await apiClient.get<ProviderHallConfig>(`${base}/config`, { signal })).data
}

export async function updateConfig(input: ProviderHallConfigInput): Promise<ProviderHallConfig> {
  return (await apiClient.put<ProviderHallConfig>(`${base}/config`, input)).data
}

export async function getGroup(id: number, signal?: AbortSignal): Promise<ProviderHallGroup> {
  return (await apiClient.get<ProviderHallGroup>(`${base}/groups/${id}`, { signal })).data
}

export async function updateGroup(id: number, input: ProviderHallGroupInput): Promise<ProviderHallGroup> {
  return (await apiClient.put<ProviderHallGroup>(`${base}/groups/${id}`, input)).data
}

export async function getTargets(id: number, signal?: AbortSignal): Promise<ProviderHallTargetSet> {
  return (await apiClient.get<ProviderHallTargetSet>(`${base}/groups/${id}/targets`, { signal })).data
}

export async function updateTargets(id: number, input: ProviderHallTargetSetInput): Promise<ProviderHallTargetSet> {
  const items = input.items.map(item => ({
    profile_id: item.profile_id,
    probe_key_id: item.probe_key_id,
    enabled: item.enabled,
    probe_interval_seconds: item.probe_interval_seconds,
    verification_interval_seconds: item.verification_interval_seconds
  }))
  return (await apiClient.put<ProviderHallTargetSet>(`${base}/groups/${id}/targets`, { version: input.version, items })).data
}

export async function listProfiles(signal?: AbortSignal): Promise<ProviderHallProfile[]> {
  return (await apiClient.get<ProviderHallProfile[]>(`${base}/profiles`, { signal })).data
}

export async function createProfile(input: ProviderHallProfileInput): Promise<ProviderHallProfile> {
  return (await apiClient.post<ProviderHallProfile>(`${base}/profiles`, input)).data
}

export async function updateProfile(id: number, input: ProviderHallProfileInput): Promise<ProviderHallProfile> {
  return (await apiClient.put<ProviderHallProfile>(`${base}/profiles/${id}`, input)).data
}

// ---- Tasks and operations (batch B6) ----

export type ProviderHallJobKind = 'probe' | 'verification'
export type ProviderHallJobStatus = 'queued' | 'running' | 'succeeded' | 'failed' | 'cancelled' | 'unknown'
export type ProviderHallSampleStatus = 'prepared' | 'dispatched' | 'received' | 'uncertain'
export type ProviderHallVerdict = 'passed' | 'failed' | 'suspected' | 'insufficient'

export interface ProviderHallJobSnapshot {
  profile: { id: number; version: number; model: string; protocol: ProviderHallProtocol; supports_tools: boolean; output_limit: number; model_aliases: string[] }
  target_version: number
  probe_key_id: number
  operator_user_id: number
  gateway_origin: string
}

export interface ProviderHallJob {
  id: number
  kind: ProviderHallJobKind
  target_id: number
  group_id: number
  profile_id: number
  config_snapshot: ProviderHallJobSnapshot
  slot_at: string | null
  not_before: string | null
  idempotency_key: string | null
  requested_by: number | null
  status: ProviderHallJobStatus
  lease_owner: string | null
  lease_until: string | null
  attempts: number
  budget_day: string | null
  error_code: string
  error_message: string
  created_at: string
  started_at: string | null
  finished_at: string | null
  updated_at: string
}

export interface ProviderHallSample {
  id: number
  job_id: number
  test_id: string
  seq: number
  trace_id: string
  status: ProviderHallSampleStatus
  prepared_at: string
  dispatched_at: string | null
  received_at: string | null
  client_request_id: string | null
  http_status: number | null
  ttft_ms: number | null
  total_ms: number | null
  generation_ms: number | null
  input_tokens: number | null
  output_tokens: number | null
  response_model: string | null
  result: 'passed' | 'failed' | 'error' | 'model_mismatch' | null
  error_code: string
  detail: Record<string, unknown>
  billing_status: string
  actual_cost: string | null
}

export interface ProviderHallVerificationSuite { passed: number; failed: number; error: number }
export interface ProviderHallVerificationSummary {
  arithmetic?: ProviderHallVerificationSuite
  json?: ProviderHallVerificationSuite
  tool?: ProviderHallVerificationSuite | null
  model?: { matched: number; mismatched: number; missing: number; seen: string[] }
  [key: string]: unknown
}

export interface ProviderHallVerification {
  job_id: number
  group_id: number
  profile_id: number
  target_id: number
  verdict: ProviderHallVerdict
  execution_status: 'completed' | 'partial' | 'error'
  reason_code: string
  summary: ProviderHallVerificationSummary
  profile_version: number
  target_version: number
  completed_at: string
  expires_at: string
  stale: boolean
  created_at: string
}

export interface ProviderHallJobDetail {
  job: ProviderHallJob
  samples: ProviderHallSample[]
  verification: ProviderHallVerification | null
}

export interface ProviderHallEnqueueResult { job_id: number; status: ProviderHallJobStatus; reused: boolean }

export interface ProviderHallPaged<T> { items: T[]; total: number; page: number; page_size: number; pages: number }

export interface ProviderHallJobFilter {
  status?: ProviderHallJobStatus | ''
  kind?: ProviderHallJobKind | ''
  group_id?: number | null
  page?: number
  page_size?: number
}

export interface ProviderHallHealthNode {
  node_id: string
  epoch_id: number
  version: string
  heartbeat_at: string
  confirmed_at: string | null
  persisted_seq: number
  overflowed: boolean
  lost: boolean
}

export interface ProviderHallHealthGap { id: number; node_id: string; epoch_id: number; scope: string; started_at: string; reason: string }

export interface ProviderHallHealth {
  generated_at: string
  collection: {
    enabled: boolean
    nodes: ProviderHallHealthNode[]
    missing_expected: string[]
    open_gaps: ProviderHallHealthGap[]
    local_queue: { depth: number; capacity: number; dropped: number }
  }
  aggregator: { watermark: string | null; lag_seconds: number | null; dirty_count: number; last_run_at: string | null; last_error: string }
  reconciliation: { pending: number; uncertain: number; failed_24h: number }
  budget: { day: string; budget: string; confirmed_spend: string; uncertain_spend: string; in_flight: number; paused_reason: string }
  jobs: { queued: number; running: number; unknown: number; failed_24h: number }
}

function enqueueBody(profileId: number, idempotencyKey: string) {
  return { profile_id: profileId, idempotency_key: idempotencyKey }
}

export async function enqueueProbe(groupId: number, profileId: number, idempotencyKey: string): Promise<ProviderHallEnqueueResult> {
  return (await apiClient.post<ProviderHallEnqueueResult>(`${base}/groups/${groupId}/probes`, enqueueBody(profileId, idempotencyKey))).data
}

export async function enqueueVerification(groupId: number, profileId: number, idempotencyKey: string): Promise<ProviderHallEnqueueResult> {
  return (await apiClient.post<ProviderHallEnqueueResult>(`${base}/groups/${groupId}/verifications`, enqueueBody(profileId, idempotencyKey))).data
}

export async function listJobs(filter: ProviderHallJobFilter = {}, signal?: AbortSignal): Promise<ProviderHallPaged<ProviderHallJob>> {
  const params: Record<string, string | number> = {}
  if (filter.status) params.status = filter.status
  if (filter.kind) params.kind = filter.kind
  if (filter.group_id) params.group_id = filter.group_id
  if (filter.page) params.page = filter.page
  if (filter.page_size) params.page_size = filter.page_size
  return (await apiClient.get<ProviderHallPaged<ProviderHallJob>>(`${base}/jobs`, { params, signal })).data
}

export async function getJob(id: number, signal?: AbortSignal): Promise<ProviderHallJobDetail> {
  return (await apiClient.get<ProviderHallJobDetail>(`${base}/jobs/${id}`, { signal })).data
}

export async function cancelJob(id: number): Promise<{ job_id: number; status: ProviderHallJobStatus }> {
  return (await apiClient.post<{ job_id: number; status: ProviderHallJobStatus }>(`${base}/jobs/${id}/cancel`)).data
}

export async function getHealth(signal?: AbortSignal): Promise<ProviderHallHealth> {
  return (await apiClient.get<ProviderHallHealth>(`${base}/health`, { signal })).data
}

export async function listGroupVerifications(groupId: number, options: { profile_id?: number | null; page?: number; page_size?: number } = {}, signal?: AbortSignal): Promise<ProviderHallPaged<ProviderHallVerification>> {
  const params: Record<string, number> = {}
  if (options.profile_id) params.profile_id = options.profile_id
  if (options.page) params.page = options.page
  if (options.page_size) params.page_size = options.page_size
  return (await apiClient.get<ProviderHallPaged<ProviderHallVerification>>(`${base}/groups/${groupId}/verifications`, { params, signal })).data
}
