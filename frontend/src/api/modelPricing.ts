/**
 * User-facing model pricing board API.
 */

import { apiClient } from './client'

export interface ModelPricingBoardItem {
  model_id: string
  platform: string
  group_id: number
  group_name: string
  channel_id: number
  channel_name: string
  rate_multiplier: number
  user_rate_overridden: boolean
  site_input_price: number | null
  site_output_price: number | null
  site_cache_read_price: number | null
  official_input_price: number | null
  official_output_price: number | null
  official_cache_read_price: number | null
  input_savings: number | null
  output_savings: number | null
  cache_read_savings: number | null
}

export interface ModelPricingBoardResponse {
  items: ModelPricingBoardItem[]
}

export async function getBoard(options?: { signal?: AbortSignal }): Promise<ModelPricingBoardResponse> {
  const { data } = await apiClient.get<ModelPricingBoardResponse>('/model-pricing/board', {
    signal: options?.signal,
  })
  return data
}

export const modelPricingAPI = { getBoard }

export default modelPricingAPI
