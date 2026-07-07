import { getActiveApiProfile, getCustomProviderDefinition } from './apiProfiles'
import { callFalAiImageApi } from './falAiImageApi'
import { callOpenAICompatibleImageApi } from './openaiCompatibleImageApi'
import type { CallApiOptions, CallApiResult } from './imageApiShared'
import { callSub2APIImageJob } from './sub2apiJobs'

export type { CallApiOptions, CallApiResult } from './imageApiShared'
export { normalizeBaseUrl } from './devProxy'

export async function callImageApi(opts: CallApiOptions): Promise<CallApiResult> {
  const profile = getActiveApiProfile(opts.settings)
  if (profile.provider === 'sub2api' && profile.apiMode === 'images') {
    const result = await callSub2APIImageJob({
      profile,
      clientTaskId: opts.clientTaskId || `${Date.now().toString(36)}-${Math.random().toString(36).slice(2)}`,
      prompt: opts.prompt,
      params: opts.params,
      inputImageDataUrls: opts.inputImageDataUrls,
      maskDataUrl: opts.maskDataUrl,
      onJobEnqueued: opts.onServerJobEnqueued,
    })
    return {
      images: result.images,
      actualParams: { n: result.images.length },
      actualParamsList: result.assets.map((asset) => asset.params),
      rawImageUrls: result.assets.map((asset) => asset.download_url || asset.url).filter(Boolean),
      serverJobId: result.job.id,
      serverJobStatus: result.job.status,
      serverAssetIds: result.assets.map((asset) => asset.id),
    }
  }
  if (profile.provider === 'fal') return callFalAiImageApi(opts, profile)

  return callOpenAICompatibleImageApi(opts, profile, getCustomProviderDefinition(opts.settings, profile.provider))
}
