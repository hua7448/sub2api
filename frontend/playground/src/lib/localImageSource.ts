import { ensureImageCached } from '../store'

export async function resolveLocalImageSource(imageIdOrSrc: string | undefined): Promise<string | undefined> {
  const value = imageIdOrSrc?.trim()
  if (!value) return undefined
  if (isInlineOrRemoteImageSource(value)) return value
  return await ensureImageCached(value) ?? undefined
}

export async function resolveLocalImageSourceOrFallback(imageId: string | undefined, fallbackSrc: string): Promise<string> {
  return await resolveLocalImageSource(imageId) ?? fallbackSrc
}

function isInlineOrRemoteImageSource(value: string): boolean {
  return value.startsWith('data:') || /^https?:\/\//i.test(value) || value.startsWith('blob:')
}
