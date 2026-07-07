import { describe, expect, it, vi, afterEach } from 'vitest'
import { DEFAULT_PARAMS, type ApiProfile, type AppSettings } from '../types'
import { DEFAULT_SETTINGS } from './apiProfiles'
import { callImageApi } from './api'

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
})
