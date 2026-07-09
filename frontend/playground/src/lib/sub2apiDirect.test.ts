import { describe, expect, it, vi, afterEach } from 'vitest'
import { DEFAULT_PARAMS, type ApiProfile, type AppSettings } from '../types'
import { DEFAULT_SETTINGS } from './apiProfiles'
import { callImageApi } from './api'
import { applySub2APISettings, type Sub2APISettings } from './sub2api'

const remoteSettings: Sub2APISettings = {
  enabled: true,
  max_upload_mb: 20,
  allowed_models: ['gpt-image-2'],
  default_model: 'gpt-image-2',
  allowed_agent_models: ['gpt-5.5'],
  agent_model: 'gpt-5.5',
  allowed_sizes: ['1024x1024', 'auto'],
  allowed_quality: ['auto', 'high'],
  allowed_output_formats: ['png', 'jpeg'],
  max_n: 4,
  agent_enabled: true,
  agent_web_search_enabled: false,
}

describe('sub2api direct image calls', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('uses the image proxy path instead of the server job path', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify({
      data: [{ b64_json: 'aW1hZ2U=' }],
    }), {
      status: 200,
      headers: { 'Content-Type': 'application/json' },
    }))

    const profile: ApiProfile = {
      id: 'sub2api-default',
      name: 'sub2api',
      provider: 'sub2api',
      baseUrl: '',
      apiKey: '',
      sub2apiKeyId: 7,
      model: 'gpt-image-1',
      timeout: 600,
      apiMode: 'images',
      codexCli: false,
      apiProxy: false,
      responseFormatB64Json: true,
      streamImages: false,
      streamPartialImages: 1,
    }
    const settings: AppSettings = {
      ...DEFAULT_SETTINGS,
      profiles: [profile],
      activeProfileId: profile.id,
      providerOrder: ['sub2api'],
    }

    const result = await callImageApi({
      settings,
      prompt: 'prompt',
      params: { ...DEFAULT_PARAMS },
      inputImageDataUrls: [],
    })

    expect(result.images).toEqual(['data:image/png;base64,aW1hZ2U='])
    expect(fetchMock).toHaveBeenCalledTimes(1)
    expect(fetchMock.mock.calls[0][0]).toBe('/api/v1/image-playground/proxy/images/generations')
    expect(String(fetchMock.mock.calls[0][0])).not.toContain('/jobs/')
    const headers = new Headers((fetchMock.mock.calls[0][1] as RequestInit).headers)
    expect(headers.get('X-Sub2API-Key-ID')).toBe('7')
  })

  it('preserves Responses API mode and uses the configured agent model', () => {
    const profile: ApiProfile = {
      id: 'sub2api-default',
      name: 'sub2api',
      provider: 'sub2api',
      baseUrl: '',
      apiKey: '',
      sub2apiKeyId: 7,
      model: 'gpt-image-2',
      timeout: 600,
      apiMode: 'responses',
      codexCli: false,
      apiProxy: false,
      responseFormatB64Json: true,
      streamImages: true,
      streamPartialImages: 1,
    }

    const applied = applySub2APISettings({
      ...DEFAULT_SETTINGS,
      profiles: [profile],
      activeProfileId: profile.id,
      apiMode: 'responses',
      model: 'gpt-image-2',
    }, remoteSettings, [{ id: 7, name: 'Key', group_name: 'group', status: 'active', quota_remaining: 100 }])

    expect(applied.apiMode).toBe('responses')
    expect(applied.model).toBe('gpt-5.5')
    expect(applied.profiles[0]).toMatchObject({
      provider: 'sub2api',
      apiMode: 'responses',
      model: 'gpt-5.5',
      apiKey: '',
      baseUrl: '',
    })
  })

  it('uses the responses proxy path for sub2api Responses API image generation', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify({
      output: [{
        type: 'image_generation_call',
        result: 'aW1hZ2U=',
      }],
    }), {
      status: 200,
      headers: { 'Content-Type': 'application/json' },
    }))

    const profile: ApiProfile = {
      id: 'sub2api-default',
      name: 'sub2api',
      provider: 'sub2api',
      baseUrl: '',
      apiKey: '',
      sub2apiKeyId: 7,
      model: 'gpt-5.5',
      timeout: 600,
      apiMode: 'responses',
      codexCli: false,
      apiProxy: false,
      responseFormatB64Json: true,
      streamImages: false,
      streamPartialImages: 1,
    }
    const settings: AppSettings = {
      ...DEFAULT_SETTINGS,
      model: profile.model,
      apiMode: 'responses',
      profiles: [profile],
      activeProfileId: profile.id,
      providerOrder: ['sub2api'],
    }

    const result = await callImageApi({
      settings,
      prompt: 'prompt',
      params: { ...DEFAULT_PARAMS },
      inputImageDataUrls: [],
    })

    expect(result.images).toEqual(['data:image/png;base64,aW1hZ2U='])
    expect(fetchMock).toHaveBeenCalledTimes(1)
    expect(fetchMock.mock.calls[0][0]).toBe('/api/v1/image-playground/proxy/responses')
    const body = JSON.parse(String((fetchMock.mock.calls[0][1] as RequestInit).body))
    expect(body.model).toBe('gpt-5.5')
    expect(body.tools[0].type).toBe('image_generation')
  })
})
