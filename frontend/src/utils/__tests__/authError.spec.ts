import { describe, expect, it } from 'vitest'
import { buildAuthErrorMessage } from '@/utils/authError'

describe('buildAuthErrorMessage', () => {
  it('prefers response detail message when available', () => {
    const message = buildAuthErrorMessage(
      {
        response: {
          data: {
            detail: 'detailed message',
            message: 'plain message'
          }
        },
      },
      { fallback: 'fallback' }
    )
    expect(message).toBe('detailed message')
  })

  it('falls back to response message when detail is unavailable', () => {
    const message = buildAuthErrorMessage(
      {
        response: {
          data: {
            message: 'plain message'
          }
        },
      },
      { fallback: 'fallback' }
    )
    expect(message).toBe('plain message')
  })

  it('falls back to error.message when response payload is unavailable', () => {
    const message = buildAuthErrorMessage(
      {
        message: 'error message'
      },
      { fallback: 'fallback' }
    )
    expect(message).toBe('error message')
  })

  it('uses fallback when no message can be extracted', () => {
    expect(buildAuthErrorMessage({}, { fallback: 'fallback' })).toBe('fallback')
  })

  it('maps a structured reason to a localized message', () => {
    const message = buildAuthErrorMessage(
      {
        response: {
          data: {
            reason: 'EMAIL_ALIAS_NOT_ALLOWED',
            message: 'provider email aliases are not allowed for registration'
          }
        }
      },
      {
        fallback: 'fallback',
        reasonMessages: {
          EMAIL_ALIAS_NOT_ALLOWED: 'localized alias error'
        }
      }
    )

    expect(message).toBe('localized alias error')
  })
})
