import { apiClient } from '../client'
import type { ImageGallerySettings } from '@/api/imageGallery'

export async function getSettings(): Promise<ImageGallerySettings> {
  const { data } = await apiClient.get<ImageGallerySettings>('/admin/image-gallery/settings')
  return data
}

export async function updateSettings(payload: ImageGallerySettings): Promise<ImageGallerySettings> {
  const { data } = await apiClient.put<ImageGallerySettings>('/admin/image-gallery/settings', payload)
  return data
}

export const adminImageGalleryAPI = {
  getSettings,
  updateSettings,
}
