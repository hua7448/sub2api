import { afterEach, describe, expect, it, vi } from 'vitest'
import { resolveLocalImageSource, resolveLocalImageSourceOrFallback } from './localImageSource'

vi.mock('../store', () => ({
  ensureImageCached: vi.fn(async (id: string) => (id === 'stored-image' ? 'data:image/png;base64,b3JpZw==' : undefined)),
}))

describe('local image source resolution', () => {
  afterEach(() => {
    vi.clearAllMocks()
  })

  it('resolves stored image ids to original data URLs', async () => {
    await expect(resolveLocalImageSource('stored-image')).resolves.toBe('data:image/png;base64,b3JpZw==')
  })

  it('keeps inline and remote image sources unchanged', async () => {
    await expect(resolveLocalImageSource('data:image/png;base64,abc')).resolves.toBe('data:image/png;base64,abc')
    await expect(resolveLocalImageSource('https://example.com/image.png')).resolves.toBe('https://example.com/image.png')
  })

  it('falls back when a stored image id cannot be resolved', async () => {
    await expect(resolveLocalImageSourceOrFallback('missing-image', 'blob:preview')).resolves.toBe('blob:preview')
  })
})
