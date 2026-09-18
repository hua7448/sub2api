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
  auto_schedule_enabled?: boolean
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
  groups?: { id: number; name: string }[]
  is_default?: boolean
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
  auto_schedule_enabled?: boolean
  probe_interval_seconds?: number
  verification_interval_seconds?: number
}

export interface ProviderHallTarget extends ProviderHallVersion, Required<Omit<ProviderHallTargetInput, 'auto_schedule_enabled'>> {
  auto_schedule_enabled?: boolean
  id: number
  group_id: number
}

export interface ProviderHallTargetSet extends ProviderHallVersion {
  listing?: ProviderHallGroup
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
    auto_schedule_enabled: item.auto_schedule_enabled,
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

/** Removes an unused profile. A profile any group still targets is refused with 409. */
export async function deleteProfile(id: number): Promise<{ profile_id: number }> {
  return (await apiClient.delete<{ profile_id: number }>(`${base}/profiles/${id}`)).data
}

/**
 * One merged candidate per model and protocol across every group, for the
 * profile editor. Unlike `listModels`, it is not scoped to a single group.
 */
export async function listProfileCandidates(signal?: AbortSignal): Promise<ProviderHallProfileCandidate[]> {
  return (await apiClient.get<ProviderHallProfileCandidate[]>(`${base}/profile-candidates`, { signal })).data
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
  group_name?: string
  source?: 'manual' | 'scheduled'
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
  group_name?: string
  model?: string
  source?: string
  profile_id?: number | null
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
  tasks_enabled?: boolean
  auto_schedule_enabled?: boolean
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

export async function enqueueProbe(groupId: number, profileId: number, idempotencyKey: string, versions?: { target_version: number; profile_version: number }): Promise<ProviderHallEnqueueResult> {
  return (await apiClient.post<ProviderHallEnqueueResult>(`${base}/groups/${groupId}/probes`, { ...enqueueBody(profileId, idempotencyKey), ...versions })).data
}

export async function enqueueVerification(groupId: number, profileId: number, idempotencyKey: string, versions?: { target_version: number; profile_version: number }): Promise<ProviderHallEnqueueResult> {
  return (await apiClient.post<ProviderHallEnqueueResult>(`${base}/groups/${groupId}/verifications`, { ...enqueueBody(profileId, idempotencyKey), ...versions })).data
}

export interface ProviderHallGroupSummary extends ProviderHallGroup {
  key_statuses?: Record<string, number>
  latest_probe?: ProviderHallAdminResult | null
  latest_verification?: ProviderHallAdminResult | null
  name: string
  platform: string
  status: string
  supported: boolean
  reason: string
  targets: ProviderHallTarget[]
  models: string[]
  effective_profile_id: number | null
}
export interface ProviderHallAdminResult { job_id: number; status: string; verdict: string; at: string }
export interface ProviderHallModelCandidate {
  model: string
  protocol: ProviderHallProtocol
  source: string
  /** Every source that offers this model, when merged across groups. */
  sources?: string[]
  account_id: number
  upstream_model: string
  available: boolean
  reason: string
  /** Groups this model/protocol can be probed in. Only the merged listing fills it. */
  groups?: number[]
}

/** A merged candidate from `GET /profile-candidates`, keyed by model + protocol. */
export type ProviderHallProfileCandidate = ProviderHallModelCandidate & { groups: number[] }
export interface ProviderHallProbeKeyOption { id: number; name: string; status: string; registered: boolean; group_id: number; operator_user_id: number }
export type ProviderHallSettingsInput = ProviderHallGroupInput & { items: ProviderHallTargetInput[] }
export async function listGroups(params: { search?: string; platform?: string; listed?: string; sort?: string; page?: number; page_size?: number } = {}, signal?: AbortSignal) {
  return (await apiClient.get<{ items: ProviderHallGroupSummary[]; total: number; page: number; page_size: number }>(`${base}/groups`, { params, signal })).data
}
export async function listModels(id: number, signal?: AbortSignal) {
  return (await apiClient.get<ProviderHallModelCandidate[]>(`${base}/groups/${id}/models`, { signal })).data
}
export async function refreshModels(id: number) {
  return (await apiClient.post<{ account_id: number; success: boolean; error?: string; models: ProviderHallModelCandidate[] }[]>(`${base}/groups/${id}/models/refresh`)).data
}
export async function listProbeKeys(id: number, signal?: AbortSignal) {
  return (await apiClient.get<ProviderHallProbeKeyOption[]>(`${base}/groups/${id}/probe-keys`, { signal })).data
}
export async function ensureProbeKey(id: number, profile_id: number) {
  return (await apiClient.post<ProviderHallProbeKeyOption>(`${base}/groups/${id}/probe-keys`, { profile_id })).data
}
export async function saveSettings(id: number, input: ProviderHallSettingsInput) {
  return (await apiClient.put<ProviderHallTargetSet>(`${base}/groups/${id}/settings`, input)).data
}
export async function preflightGroup(id: number, input: ProviderHallSettingsInput) {
  return (await apiClient.post<ProviderHallTargetSet>(`${base}/groups/${id}/preflight`, input)).data
}
export async function preflightConfig(input: ProviderHallConfigInput) {
  return (await apiClient.post<{ valid: boolean; disabled_targets: ProviderHallTarget[] }>(`${base}/config/preflight`, input)).data
}
export async function checkGateway(origin: string) {
  return (await apiClient.post<{ status: number; healthy: boolean }>(`${base}/config/check-gateway`, { origin })).data
}
export interface ProviderHallBatchResult { id: number; success: boolean; error?: string; reason?: string; metadata?: Record<string, string>; settings?: ProviderHallTargetSet }
export async function batchGroups(groups: (ProviderHallSettingsInput & { id: number })[], preview: boolean) {
  return (await apiClient.post<ProviderHallBatchResult[]>(`${base}/groups/batch`, { groups, preview })).data
}

export async function listJobs(filter: ProviderHallJobFilter = {}, signal?: AbortSignal): Promise<ProviderHallPaged<ProviderHallJob>> {
  const params: Record<string, string | number> = {}
  if (filter.group_name) params.group_name = filter.group_name
  if (filter.model) params.model = filter.model
  if (filter.source) params.source = filter.source
  if (filter.profile_id) params.profile_id = filter.profile_id
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
