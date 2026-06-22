import { afterEach, describe, expect, it, vi } from 'vitest'
import { dataUrlToBlob } from './canvasImage'

describe('dataUrlToBlob', () => {
  const originalFetch = globalThis.fetch

  afterEach(() => {
    vi.restoreAllMocks()
    globalThis.fetch = originalFetch
  })

  it('decodes data URLs without using fetch', async () => {
    globalThis.fetch = vi.fn(async () => {
      throw new Error('fetch should not be called for data URLs')
    }) as unknown as typeof fetch

    const blob = await dataUrlToBlob('data:image/png;base64,aGVsbG8=')

    expect(globalThis.fetch).not.toHaveBeenCalled()
    expect(blob.type).toBe('image/png')
    expect(await blob.text()).toBe('hello')
  })
})
