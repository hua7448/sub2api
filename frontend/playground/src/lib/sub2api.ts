import type { ApiProfile, AppSettings, TaskParams } from '../types'

export interface Sub2APISettings {
  enabled: boolean
  max_upload_mb: number
  allowed_models: string[]
  default_model: string
  allowed_sizes: string[]
  allowed_quality: string[]
  allowed_output_formats: string[]
  max_n: number
  agent_enabled: boolean
  agent_web_search_enabled: boolean
}

export interface Sub2APIEligibleKey {
  id: number
  name: string
  group_name: string
  status: string
  quota_remaining: number
}

const API_PREFIX = '/api/v1/image-playground'

let cachedSettings: Sub2APISettings | null = null
let cachedKeys: Sub2APIEligibleKey[] | null = null

async function readJSON<T>(response: Response): Promise<T> {
  const payload = await response.json().catch(() => null)
  if (!response.ok) {
    const message = payload?.message || payload?.error?.message || payload?.error || `HTTP ${response.status}`
    throw new Error(String(message))
  }
  return (payload?.data ?? payload) as T
}

export async function fetchSub2APISettings(): Promise<Sub2APISettings> {
  const response = await fetch(`${API_PREFIX}/settings`, { credentials: 'same-origin' })
  cachedSettings = await readJSON<Sub2APISettings>(response)
  return cachedSettings
}

export async function fetchSub2APIEligibleKeys(): Promise<Sub2APIEligibleKey[]> {
  const response = await fetch(`${API_PREFIX}/eligible-keys`, { credentials: 'same-origin' })
  cachedKeys = await readJSON<Sub2APIEligibleKey[]>(response)
  return cachedKeys
}

export function getCachedSub2APISettings() {
  return cachedSettings
}

export function getCachedSub2APIEligibleKeys() {
  return cachedKeys ?? []
}

export function selectSub2APIKeyId(profile: ApiProfile): number {
  const configured = Number(profile.sub2apiKeyId)
  if (Number.isFinite(configured) && configured > 0) return Math.trunc(configured)
  const first = cachedKeys?.[0]?.id
  if (first) return first
  throw new Error('请先在生图广场中选择可用 API Key')
}

export function applySub2APISettings(settings: AppSettings, remote: Sub2APISettings, keys: Sub2APIEligibleKey[]): AppSettings {
  const activeKeyId = keys[0]?.id ?? null
  const model = remote.default_model || remote.allowed_models[0] || 'gpt-image-2'
  const profile: ApiProfile = {
    id: 'sub2api-default',
    name: activeKeyId ? `sub2api #${activeKeyId}` : 'sub2api',
    provider: 'sub2api',
    baseUrl: '',
    apiKey: '',
    sub2apiKeyId: activeKeyId,
    model,
    timeout: 600,
    apiMode: 'images',
    codexCli: false,
    apiProxy: false,
    responseFormatB64Json: true,
    streamImages: true,
    streamPartialImages: 1,
  }

  return {
    ...settings,
    baseUrl: '',
    apiKey: '',
    model,
    apiMode: profile.apiMode,
    codexCli: false,
    apiProxy: false,
    agentWebSearch: remote.agent_web_search_enabled ? settings.agentWebSearch : false,
    customProviders: [],
    providerOrder: ['sub2api'],
    profiles: [profile],
    activeProfileId: profile.id,
  }
}

export function assertSub2APIParams(params: TaskParams, settings: Sub2APISettings | null) {
  if (!settings) return
  if (!settings.enabled) throw new Error('生图广场当前未启用')
  if (settings.allowed_sizes.length > 0 && params.size && !settings.allowed_sizes.includes(params.size)) {
    throw new Error('当前尺寸不在管理员允许范围内')
  }
  if (settings.allowed_quality.length > 0 && params.quality && !settings.allowed_quality.includes(params.quality)) {
    throw new Error('当前质量不在管理员允许范围内')
  }
  if (settings.allowed_output_formats.length > 0 && params.output_format && !settings.allowed_output_formats.includes(params.output_format)) {
    throw new Error('当前输出格式不在管理员允许范围内')
  }
  if (settings.max_n > 0 && params.n > settings.max_n) {
    throw new Error(`单次最多生成 ${settings.max_n} 张图片`)
  }
}

export async function proxySub2API(path: 'images/generations' | 'images/edits' | 'responses', profile: ApiProfile, init: RequestInit): Promise<Response> {
  const keyId = selectSub2APIKeyId(profile)
  const headers = new Headers(init.headers)
  headers.delete('Authorization')
  headers.delete('Proxy-Authorization')
  headers.delete('x-api-key')
  headers.delete('x-goog-api-key')
  headers.set('X-Sub2API-Key-ID', String(keyId))

  return fetch(`${API_PREFIX}/proxy/${path}`, {
    ...init,
    method: 'POST',
    credentials: 'same-origin',
    headers,
  })
}
