import { afterEach, describe, expect, it, vi } from 'vitest'
import { copyImageSourceToClipboard } from './clipboard'

describe('copyImageSourceToClipboard', () => {
  const originalClipboard = navigator.clipboard
  const originalClipboardItem = globalThis.ClipboardItem

  afterEach(() => {
    vi.restoreAllMocks()
    Object.defineProperty(navigator, 'clipboard', {
      configurable: true,
      value: originalClipboard,
    })
    Object.defineProperty(globalThis, 'ClipboardItem', {
      configurable: true,
      value: originalClipboardItem,
    })
  })

  it('calls clipboard.write with the clipboard object as receiver', async () => {
    const clipboard = {
      write: vi.fn(function (this: unknown) {
        if (this !== clipboard) throw new TypeError('Illegal invocation')
        return Promise.resolve()
      }),
    }

    class TestClipboardItem {
      static supports(type: string) {
        return type === 'image/png'
      }

      constructor(public readonly items: Record<string, Blob | Promise<Blob>>) {}
    }

    Object.defineProperty(navigator, 'clipboard', {
      configurable: true,
      value: clipboard,
    })
    Object.defineProperty(globalThis, 'ClipboardItem', {
      configurable: true,
      value: TestClipboardItem,
    })

    await copyImageSourceToClipboard('data:image/png;base64,aGVsbG8=')

    expect(clipboard.write).toHaveBeenCalledOnce()
  })
})
