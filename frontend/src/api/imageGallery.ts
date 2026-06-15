import { apiClient } from './client'
import type { PaginatedResponse } from '@/types'

export interface ImageGallerySettings {
  enabled: boolean
  public_enabled: boolean
  publish_requires_review: boolean
  storage_backend: string
  storage_path?: string
  max_upload_mb: number
  user_quota_mb: number
  retention_days: number
  allowed_models: string[]
  default_model: string
  allowed_sizes: string[]
  allowed_quality: string[]
  allowed_output_formats: string[]
  max_n: number
  templates_enabled: boolean
  template_import_enabled: boolean
}

export interface ImageGalleryEligibleKey {
  id: number
  name: string
  group_name: string
  status: string
  quota_remaining: number
}

export interface ImageGalleryTemplate {
  id: number
  title: string
  prompt: string
  category: string
  tags: string[]
  source: string
  source_url: string
  author: string
  license: string
  enabled: boolean
  featured: boolean
}

export interface ImageGalleryAsset {
  id: number
  user_id: number
  job_id?: number
  usage_log_id?: number
  api_key_id?: number
  url: string
  download_url: string
  sha256: string
  mime_type: string
  width: number
  height: number
  size_bytes: number
  visibility: 'private' | 'public'
  review_status: 'none' | 'pending' | 'approved' | 'rejected'
  hidden: boolean
  prompt: string
  model: string
  params: Record<string, unknown>
  actual_cost?: number
  favorited: boolean
  created_at: string
}

export interface ImageGalleryJob {
  id: number
  status: string
  model: string
  prompt: string
  error: string
  created_at: string
  completed_at?: string
}

export interface GenerateImageRequest {
  api_key_id: number
  prompt: string
  model: string
  size?: string
  quality?: string
  n: number
  output_format?: string
  reference_images?: string[]
  mask_image?: string
}

export interface GenerateImageResult {
  job: ImageGalleryJob
  assets: ImageGalleryAsset[]
}

export async function getSettings(): Promise<ImageGallerySettings> {
  const { data } = await apiClient.get<ImageGallerySettings>('/image-gallery/settings')
  return data
}

export async function getEligibleKeys(): Promise<ImageGalleryEligibleKey[]> {
  const { data } = await apiClient.get<ImageGalleryEligibleKey[]>('/image-gallery/eligible-keys')
  return data
}

export async function listTemplates(params?: { q?: string; category?: string; page?: number; page_size?: number }): Promise<PaginatedResponse<ImageGalleryTemplate>> {
  const { data } = await apiClient.get<PaginatedResponse<ImageGalleryTemplate>>('/image-gallery/templates', { params })
  return data
}

export async function generateImage(payload: GenerateImageRequest): Promise<GenerateImageResult> {
  const { data } = await apiClient.post<GenerateImageResult>('/image-gallery/generations', payload)
  return data
}

export async function listHistory(page = 1, pageSize = 20): Promise<PaginatedResponse<ImageGalleryAsset>> {
  const { data } = await apiClient.get<PaginatedResponse<ImageGalleryAsset>>('/image-gallery/history', { params: { page, page_size: pageSize } })
  return data
}

export async function listPublic(page = 1, pageSize = 20): Promise<PaginatedResponse<ImageGalleryAsset>> {
  const { data } = await apiClient.get<PaginatedResponse<ImageGalleryAsset>>('/image-gallery/public', { params: { page, page_size: pageSize } })
  return data
}

export async function publishItem(id: number): Promise<ImageGalleryAsset> {
  const { data } = await apiClient.post<ImageGalleryAsset>(`/image-gallery/items/${id}/publish`)
  return data
}

export async function deleteItem(id: number): Promise<void> {
  await apiClient.delete(`/image-gallery/items/${id}`)
}

export async function favoritePublicItem(id: number): Promise<void> {
  await apiClient.post(`/image-gallery/public/${id}/favorite`)
}

export const imageGalleryAPI = {
  getSettings,
  getEligibleKeys,
  listTemplates,
  generateImage,
  listHistory,
  listPublic,
  publishItem,
  deleteItem,
  favoritePublicItem,
}
