import { afterEach, describe, expect, it, vi } from 'vitest'
import { downloadImageIds, getImageZipEntries } from './downloadImages'

describe('downloadImageIds', () => {
  const originalFetch = globalThis.fetch
  const originalDocument = (globalThis as typeof globalThis & { document?: Document }).document
  const originalWindow = (globalThis as typeof globalThis & { window?: Window }).window
  const originalCreateObjectURL = URL.createObjectURL
  const originalRevokeObjectURL = URL.revokeObjectURL

  afterEach(() => {
    vi.restoreAllMocks()
    globalThis.fetch = originalFetch
    if (originalDocument) {
      Object.defineProperty(globalThis, 'document', { configurable: true, value: originalDocument })
    } else {
      delete (globalThis as typeof globalThis & { document?: Document }).document
    }
    if (originalWindow) {
      Object.defineProperty(globalThis, 'window', { configurable: true, value: originalWindow })
    } else {
      delete (globalThis as typeof globalThis & { window?: Window }).window
    }
    URL.createObjectURL = originalCreateObjectURL
    URL.revokeObjectURL = originalRevokeObjectURL
  })

  it('downloads data URL images without using fetch', async () => {
    const clickSpy = vi.fn()
    const appendSpy = vi.fn()
    const removeSpy = vi.fn()
    const anchor = {
      href: '',
      download: '',
      click: clickSpy,
    }

    Object.defineProperty(globalThis, 'document', {
      configurable: true,
      value: {
        createElement: vi.fn(() => anchor),
        body: {
          appendChild: appendSpy,
          removeChild: removeSpy,
        },
      },
    })
    Object.defineProperty(globalThis, 'window', {
      configurable: true,
      value: {
        setTimeout: vi.fn((callback: () => void) => {
          callback()
          return 0
        }),
      },
    })
    URL.createObjectURL = vi.fn(() => 'blob:download-url')
    URL.revokeObjectURL = vi.fn()
    globalThis.fetch = vi.fn(async () => {
      throw new Error('fetch should not be called for data URLs')
    }) as unknown as typeof fetch

    const result = await downloadImageIds(['data:image/png;base64,aGVsbG8='], 'image')

    expect(result).toEqual({ successCount: 1, failCount: 0 })
    expect(globalThis.fetch).not.toHaveBeenCalled()
    expect(URL.createObjectURL).toHaveBeenCalledWith(expect.any(Blob))
    expect(clickSpy).toHaveBeenCalledTimes(1)
    expect(appendSpy).toHaveBeenCalledTimes(1)
    expect(removeSpy).toHaveBeenCalledTimes(1)
  })

  it('uses stable numbered bases for ZIP entries', () => {
    expect(getImageZipEntries(['img-1'], 'task-task-1')).toEqual([
      { imageId: 'img-1', fileNameBase: 'task-task-1' },
    ])
    expect(getImageZipEntries(['img-1', 'img-2'], 'task-task-1')).toEqual([
      { imageId: 'img-1', fileNameBase: 'task-task-1-01' },
      { imageId: 'img-2', fileNameBase: 'task-task-1-02' },
    ])
  })
})
