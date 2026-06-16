import { apiClient } from './client'

export interface ImageGallerySettings {
  enabled: boolean
  max_upload_mb: number
  allowed_models: string[]
  default_model: string
  allowed_sizes: string[]
  allowed_quality: string[]
  allowed_output_formats: string[]
  max_n: number
  agent_enabled: boolean
  agent_web_search_enabled: boolean
}

export interface ImageGalleryEligibleKey {
  id: number
  name: string
  group_name: string
  status: string
  quota_remaining: number
}

export async function getSettings(): Promise<ImageGallerySettings> {
  const { data } = await apiClient.get<ImageGallerySettings>('/image-playground/settings')
  return data
}

export async function getEligibleKeys(): Promise<ImageGalleryEligibleKey[]> {
  const { data } = await apiClient.get<ImageGalleryEligibleKey[]>('/image-playground/eligible-keys')
  return data
}

export const imageGalleryAPI = {
  getSettings,
  getEligibleKeys,
}
