export interface ParsedOpenAIRefreshTokenBatchItem {
  raw: string
  account: string
  password: string
  emailPassword: string
  refreshToken: string
}

const BATCH_LINE_SEPARATOR = '----'

export function parseOpenAIRefreshTokenBatchLine(
  line: string
): ParsedOpenAIRefreshTokenBatchItem | null {
  const raw = line.trim()
  if (!raw) {
    return null
  }

  const segments = raw.split(BATCH_LINE_SEPARATOR).map((item) => item.trim())
  if (segments.length < 4) {
    return {
      raw,
      account: '',
      password: '',
      emailPassword: '',
      refreshToken: raw
    }
  }

  const refreshToken = segments[segments.length - 1]
  const account = segments[0]
  const password = segments[1]
  const emailPassword = segments.slice(2, -1).join(BATCH_LINE_SEPARATOR).trim()

  if (!refreshToken) {
    return null
  }

  return {
    raw,
    account,
    password,
    emailPassword,
    refreshToken
  }
}

export function parseOpenAIRefreshTokenBatchInput(
  input: string
): ParsedOpenAIRefreshTokenBatchItem[] {
  return input
    .split('\n')
    .map((line) => parseOpenAIRefreshTokenBatchLine(line))
    .filter((item): item is ParsedOpenAIRefreshTokenBatchItem => item !== null)
}
