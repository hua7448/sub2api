import { apiClient } from '../client'
import type { PaginatedResponse } from '@/types'
import type { ImageGalleryAsset, ImageGallerySettings, ImageGalleryTemplate } from '@/api/imageGallery'

export interface ImageGalleryStorageStats {
  storage_path: string
  file_count: number
  total_bytes: number
}

export interface ImageGalleryCleanupResult {
  deleted_files: number
  deleted_bytes: number
}

export async function getSettings(): Promise<ImageGallerySettings> {
  const { data } = await apiClient.get<ImageGallerySettings>('/admin/image-gallery/settings')
  return data
}

export async function updateSettings(payload: ImageGallerySettings): Promise<ImageGallerySettings> {
  const { data } = await apiClient.put<ImageGallerySettings>('/admin/image-gallery/settings', payload)
  return data
}

export async function listItems(params?: { review_status?: string; page?: number; page_size?: number }): Promise<PaginatedResponse<ImageGalleryAsset>> {
  const { data } = await apiClient.get<PaginatedResponse<ImageGalleryAsset>>('/admin/image-gallery/items', { params })
  return data
}

export async function approveItem(id: number): Promise<ImageGalleryAsset> {
  const { data } = await apiClient.post<ImageGalleryAsset>(`/admin/image-gallery/items/${id}/approve`)
  return data
}

export async function rejectItem(id: number): Promise<ImageGalleryAsset> {
  const { data } = await apiClient.post<ImageGalleryAsset>(`/admin/image-gallery/items/${id}/reject`)
  return data
}

export async function hideItem(id: number): Promise<ImageGalleryAsset> {
  const { data } = await apiClient.post<ImageGalleryAsset>(`/admin/image-gallery/items/${id}/hide`)
  return data
}

export async function deleteItem(id: number): Promise<void> {
  await apiClient.delete(`/admin/image-gallery/items/${id}`)
}

export async function listTemplates(params?: { q?: string; category?: string; page?: number; page_size?: number }): Promise<PaginatedResponse<ImageGalleryTemplate>> {
  const { data } = await apiClient.get<PaginatedResponse<ImageGalleryTemplate>>('/admin/image-gallery/templates', { params })
  return data
}

export async function updateTemplate(id: number, payload: Partial<ImageGalleryTemplate>): Promise<ImageGalleryTemplate> {
  const { data } = await apiClient.patch<ImageGalleryTemplate>(`/admin/image-gallery/templates/${id}`, payload)
  return data
}

export async function getStorageStats(): Promise<ImageGalleryStorageStats> {
  const { data } = await apiClient.get<ImageGalleryStorageStats>('/admin/image-gallery/storage/stats')
  return data
}

export async function cleanupStorage(): Promise<ImageGalleryCleanupResult> {
  const { data } = await apiClient.post<ImageGalleryCleanupResult>('/admin/image-gallery/storage/cleanup')
  return data
}

export const adminImageGalleryAPI = {
  getSettings,
  updateSettings,
  listItems,
  approveItem,
  rejectItem,
  hideItem,
  deleteItem,
  listTemplates,
  updateTemplate,
  getStorageStats,
  cleanupStorage,
}
