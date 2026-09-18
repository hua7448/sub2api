/**
 * User Groups API endpoints (non-admin)
 * Handles group-related operations for regular users
 */

import { apiClient } from './client'
import type { Group, GroupPlatform } from '@/types'

export interface GroupModelsResponse {
  group_id: number
  group_name: string
  platform: GroupPlatform
  models: string[]
  is_default: boolean
}

/**
 * Get available groups that the current user can bind to API keys
 * This returns groups based on user's permissions:
 * - Standard groups: public (non-exclusive) or explicitly allowed
 * - Subscription groups: user has active subscription
 * @returns List of available groups
 */
export async function getAvailable(): Promise<Group[]> {
  const { data } = await apiClient.get<Group[]>('/groups/available')
  return data
}

/**
 * Get current user's custom group rate multipliers
 * @returns Map of group_id to custom rate_multiplier
 */
export async function getUserGroupRates(): Promise<Record<number, number>> {
  const { data } = await apiClient.get<Record<number, number> | null>('/groups/rates')
  return data || {}
}

/**
 * Get available model IDs for a specific group the current user can access.
 */
export async function getModels(groupId: number): Promise<GroupModelsResponse> {
  const { data } = await apiClient.get<GroupModelsResponse>(`/groups/${groupId}/models`)
  return data
}

export const userGroupsAPI = {
  getAvailable,
  getUserGroupRates,
  getModels
}

export default userGroupsAPI
