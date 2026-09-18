import { extractApiErrorCode, extractApiErrorMessage } from '@/utils/apiError'

export const protocols = [
  { value: 'responses', label: 'Responses' },
  { value: 'chat_completions', label: 'Chat Completions' },
  { value: 'messages', label: 'Messages' },
]

export function hallError(error: unknown, t: (key: string) => string, fallback = 'saveFailed'): string {
  const keys: Record<string, string> = {
    PROVIDER_HALL_VERSION_CONFLICT: 'conflict',
    PROVIDER_HALL_INVALID_CONFIG: 'invalid',
    PROVIDER_HALL_PROBE_KEY_INVALID: 'invalidKey',
    PROVIDER_HALL_OPERATOR_NOT_CONFIGURED: 'invalidOperator',
    PROVIDER_HALL_TARGET_ROUTE_UNSUPPORTED: 'invalidRoute',
    PROVIDER_HALL_NOT_READY: 'notReady',
    PROVIDER_HALL_PROFILE_EXISTS: 'duplicateProfile',
    PROVIDER_HALL_NOT_FOUND: 'missing',
    PROVIDER_HALL_TASKS_DISABLED: 'tasksDisabled',
    PROVIDER_HALL_BUDGET_EXHAUSTED: 'budgetExhausted',
    PROVIDER_HALL_TARGET_DISABLED: 'targetDisabled',
    PROVIDER_HALL_JOB_IN_FLIGHT: 'jobInFlight',
    PROVIDER_HALL_JOB_NOT_CANCELLABLE: 'jobNotCancellable',
  }
  const code = extractApiErrorCode(error) ?? ''
  const key = keys[code]
  if (!key) return extractApiErrorMessage(error, t(`admin.providerHall.${fallback}`))
  const field = code === 'PROVIDER_HALL_INVALID_CONFIG' ? invalidField(error) : ''
  if (!field) return t(`admin.providerHall.${key}`)
  const hintKey = `admin.providerHall.fieldHint.${field}`
  const hint = t(hintKey)
  return `${t(`admin.providerHall.${key}`)}：${hint === hintKey ? field : hint}`
}

function invalidField(error: unknown): string {
  const metadata = (error as { metadata?: Record<string, unknown> } | null)?.metadata
  const field = metadata?.field
  return typeof field === 'string' ? field : ''
}

export function lines(value: string): string[] {
  return value.split(/\r?\n/).map(line => line.trim()).filter(Boolean)
}

/** Idempotency key for one admin click; kept until the response returns. */
export function idempotencyKey(prefix: string): string {
  const random = typeof crypto !== 'undefined' && 'randomUUID' in crypto ? crypto.randomUUID() : `${Date.now()}-${Math.random().toString(16).slice(2)}`
  return `${prefix}-${random}`.slice(0, 128)
}
