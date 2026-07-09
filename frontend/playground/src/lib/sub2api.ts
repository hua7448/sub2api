import { DEFAULT_STREAM_PARTIAL_IMAGES, type ApiProfile, type AppSettings, type TaskParams } from '../types'
import { DEFAULT_IMAGES_MODEL, DEFAULT_RESPONSES_MODEL } from './apiProfiles'

export interface Sub2APISettings {
  enabled: boolean
  max_upload_mb: number
  allowed_models: string[]
  default_model: string
  allowed_agent_models: string[]
  agent_model: string
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

function authHeaders(init?: HeadersInit): Headers {
  const headers = new Headers(init)
  const token = typeof window !== 'undefined' ? window.localStorage.getItem('auth_token') : null
  if (token && !headers.has('Authorization')) headers.set('Authorization', `Bearer ${token}`)
  return headers
}

async function readJSON<T>(response: Response): Promise<T> {
  const payload = await response.json().catch(() => null)
  if (!response.ok) {
    const message = payload?.message || payload?.error?.message || payload?.error || `HTTP ${response.status}`
    throw new Error(String(message))
  }
  return (payload?.data ?? payload) as T
}

export async function fetchSub2APISettings(): Promise<Sub2APISettings> {
  const response = await fetch(`${API_PREFIX}/settings`, { credentials: 'same-origin', headers: authHeaders() })
  cachedSettings = await readJSON<Sub2APISettings>(response)
  return cachedSettings
}

export async function fetchSub2APIEligibleKeys(): Promise<Sub2APIEligibleKey[]> {
  const response = await fetch(`${API_PREFIX}/eligible-keys`, { credentials: 'same-origin', headers: authHeaders() })
  cachedKeys = await readJSON<Sub2APIEligibleKey[]>(response)
  return cachedKeys
}

export function getCachedSub2APISettings() {
  return cachedSettings
}

export function getCachedSub2APIEligibleKeys() {
  return cachedKeys ?? []
}

export function getSub2APIModelsForMode(settings: Sub2APISettings | null | undefined, apiMode: AppSettings['apiMode']): string[] {
  const models = apiMode === 'responses'
    ? settings?.allowed_agent_models ?? []
    : settings?.allowed_models ?? []
  return models.filter((model) => model.trim())
}

export function getSub2APIDefaultModelForMode(settings: Sub2APISettings | null | undefined, apiMode: AppSettings['apiMode']): string {
  if (apiMode === 'responses') {
    return settings?.agent_model || settings?.allowed_agent_models?.[0] || DEFAULT_RESPONSES_MODEL
  }
  return settings?.default_model || settings?.allowed_models?.[0] || DEFAULT_IMAGES_MODEL
}

function isSub2APIModelAllowed(settings: Sub2APISettings, apiMode: AppSettings['apiMode'], model: string | undefined): model is string {
  if (!model) return false
  const allowedModels = getSub2APIModelsForMode(settings, apiMode)
  return allowedModels.length === 0 || allowedModels.includes(model)
}

export function selectSub2APIKeyId(profile: ApiProfile): number {
  const configured = Number(profile.sub2apiKeyId)
  if (Number.isFinite(configured) && configured > 0) return Math.trunc(configured)
  const first = cachedKeys?.[0]?.id
  if (first) return first
  throw new Error('请先在生图广场中选择可用 API Key')
}

export function applySub2APISettings(settings: AppSettings, remote: Sub2APISettings, keys: Sub2APIEligibleKey[]): AppSettings {
  const persistedProfile = settings.profiles?.find((profile) => profile.provider === 'sub2api')
  const persistedKeyId = Number(persistedProfile?.sub2apiKeyId)
  const activeKeyId = keys.some((key) => key.id === persistedKeyId) ? persistedKeyId : keys[0]?.id ?? null
  const persistedApiMode = persistedProfile?.apiMode === 'responses' || settings.apiMode === 'responses' ? 'responses' : 'images'
  const persistedModel = persistedProfile?.model || settings.model
  const model = isSub2APIModelAllowed(remote, persistedApiMode, persistedModel)
    ? persistedModel
    : getSub2APIDefaultModelForMode(remote, persistedApiMode)
  const activeKey = activeKeyId ? keys.find((key) => key.id === activeKeyId) ?? null : null
  const profile: ApiProfile = {
    id: 'sub2api-default',
    name: activeKey ? activeKey.name || `API Key #${activeKey.id}` : 'sub2api',
    provider: 'sub2api',
    baseUrl: '',
    apiKey: '',
    sub2apiKeyId: activeKeyId,
    model,
    timeout: 600,
    apiMode: persistedApiMode,
    codexCli: false,
    apiProxy: false,
    responseFormatB64Json: true,
    streamImages: persistedProfile?.streamImages ?? true,
    streamPartialImages: persistedProfile?.streamPartialImages ?? DEFAULT_STREAM_PARTIAL_IMAGES,
  }

  return {
    ...settings,
    baseUrl: '',
    apiKey: '',
    model,
    apiMode: persistedApiMode,
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
  const headers = authHeaders(init.headers)
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
