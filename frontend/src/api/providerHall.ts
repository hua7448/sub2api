import { apiClient } from './client'
import type {
  ProviderHallDetailResponse,
  ProviderHallListResponse,
  ProviderHallRange,
  ProviderHallSortRule,
  ProviderHallVerificationListResponse,
} from '@/types/providerHall'

export type {
  ProviderHallDetailResponse,
  ProviderHallListResponse,
  ProviderHallRow,
  ProviderHallVerificationListResponse,
  ProviderHallVerificationReport,
} from '@/types/providerHall'

export interface ProviderHallListParams {
  range?: ProviderHallRange
  model?: string
  protocol?: string
  search?: string
  sort?: ProviderHallSortRule[]
  page?: number
  page_size?: number
}

export interface ProviderHallVerificationParams {
  profile_id?: number
  page?: number
  page_size?: number
}

const base = '/provider-hall'

/** `field[:asc|desc]` joined by commas; at most three rules are sent. */
export function serializeSort(rules: ProviderHallSortRule[] | undefined): string | undefined {
  if (!rules || rules.length === 0) return undefined
  return rules
    .slice(0, 3)
    .map((rule) => `${rule.field}:${rule.direction}`)
    .join(',')
}

export async function list(params: ProviderHallListParams = {}, signal?: AbortSignal): Promise<ProviderHallListResponse> {
  const query: Record<string, string | number> = {}
  if (params.range) query.range = params.range
  if (params.model) query.model = params.model
  if (params.protocol) query.protocol = params.protocol
  if (params.search) query.search = params.search
  const sort = serializeSort(params.sort)
  if (sort) query.sort = sort
  if (params.page) query.page = params.page
  if (params.page_size) query.page_size = params.page_size
  const { data } = await apiClient.get<ProviderHallListResponse>(base, { params: query, signal })
  return data
}

export async function getGroup(id: number, range?: ProviderHallRange, signal?: AbortSignal): Promise<ProviderHallDetailResponse> {
  const { data } = await apiClient.get<ProviderHallDetailResponse>(`${base}/groups/${id}`, {
    params: range ? { range } : undefined,
    signal,
  })
  return data
}

export async function listVerifications(
  id: number,
  params: ProviderHallVerificationParams = {},
  signal?: AbortSignal
): Promise<ProviderHallVerificationListResponse> {
  const query: Record<string, number> = {}
  if (params.profile_id) query.profile_id = params.profile_id
  if (params.page) query.page = params.page
  if (params.page_size) query.page_size = params.page_size
  const { data } = await apiClient.get<ProviderHallVerificationListResponse>(`${base}/groups/${id}/verifications`, {
    params: query,
    signal,
  })
  return data
}

export const providerHallAPI = { list, getGroup, listVerifications, serializeSort }

export default providerHallAPI
