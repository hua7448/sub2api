import type { ApiProfile, TaskParams } from '../types'
import { dataUrlToBlob, imageDataUrlToPngBlob, maskDataUrlToPngBlob } from './canvasImage'
import { fetchImageUrlAsDataUrl, getDataUrlEncodedByteSize, getDataUrlDecodedByteSize, MIME_MAP } from './imageApiShared'
import { assertSub2APIParams, getCachedSub2APISettings, selectSub2APIKeyId } from './sub2api'

const API_PREFIX = '/api/v1/image-playground'
const JOB_POLL_MS = 5_000
const REQUEST_TIMEOUT_MS = 30_000

export type Sub2APIJobStatus = 'pending' | 'running' | 'succeeded' | 'failed' | 'save_failed' | 'cancelled'

export interface Sub2APIJob {
  id: number
  status: Sub2APIJobStatus
  error?: string
  model?: string
  prompt?: string
}

export interface Sub2APIAsset {
  id: number
  url: string
  download_url?: string
  mime_type?: string
  width?: number
  height?: number
  params?: Partial<TaskParams>
}

export interface Sub2APIJobResult {
  job: Sub2APIJob
  assets?: Sub2APIAsset[]
}

export interface Sub2APIJobSubmitOptions {
  profile: ApiProfile
  clientTaskId: string
  prompt: string
  params: TaskParams
  inputImageDataUrls: string[]
  maskDataUrl?: string
  onJobEnqueued?: (job: { jobId: number; status?: string }) => void
}

export class Sub2APIJobRecoverableError extends Error {
  jobId?: number
  clientTaskId?: string

  constructor(message: string, details: { jobId?: number; clientTaskId?: string } = {}) {
    super(message)
    this.name = 'Sub2APIJobRecoverableError'
    this.jobId = details.jobId
    this.clientTaskId = details.clientTaskId
  }
}

export function isSub2APIJobRecoverableError(err: unknown): err is Sub2APIJobRecoverableError {
  return err instanceof Sub2APIJobRecoverableError
}

export function isRecoverableSub2APINetworkError(err: unknown): boolean {
  if (typeof DOMException !== 'undefined' && err instanceof DOMException && err.name === 'AbortError') return true
  if (err instanceof TypeError) {
    return /failed to fetch|fetch failed|load failed|networkerror|network request failed/i.test(err.message)
  }
  const message = err instanceof Error ? err.message : String(err)
  return /abort|network|failed to fetch|fetch failed|load failed|timeout|连接|断开|中断/i.test(message)
}

function authHeaders(init?: HeadersInit): Headers {
  const headers = new Headers(init)
  const token = typeof window !== 'undefined' ? window.localStorage.getItem('auth_token') : null
  if (token && !headers.has('Authorization')) headers.set('Authorization', `Bearer ${token}`)
  return headers
}

async function fetchWithTimeout(path: string, init: RequestInit = {}, timeoutMs = REQUEST_TIMEOUT_MS): Promise<Response> {
  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), timeoutMs)
  try {
    return await fetch(`${API_PREFIX}${path}`, {
      ...init,
      credentials: 'same-origin',
      headers: authHeaders(init.headers),
      signal: controller.signal,
      cache: 'no-store',
    })
  } finally {
    clearTimeout(timeoutId)
  }
}

async function readJSON<T>(response: Response): Promise<T> {
  const payload = await response.json().catch(() => null)
  if (!response.ok) {
    const message = payload?.message || payload?.error?.message || payload?.error || `HTTP ${response.status}`
    throw new Error(String(message))
  }
  return (payload?.data ?? payload) as T
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

function estimateUploadBytes(inputImageDataUrls: string[], maskDataUrl?: string): number {
  return inputImageDataUrls.reduce((sum, dataUrl) => sum + getDataUrlEncodedByteSize(dataUrl), 0) +
    (maskDataUrl ? getDataUrlEncodedByteSize(maskDataUrl) : 0)
}

function assertSub2APIUploadLimits(inputImageDataUrls: string[], maskDataUrl?: string) {
  const settings = getCachedSub2APISettings()
  const maxUploadBytes = Math.max(1, settings?.max_upload_mb ?? 20) * 1024 * 1024
  const files = [...inputImageDataUrls, ...(maskDataUrl ? [maskDataUrl] : [])]
  for (const dataUrl of files) {
    if (getDataUrlDecodedByteSize(dataUrl) > maxUploadBytes) {
      throw new Error(`文件太大：单文件上限为 ${(maxUploadBytes / 1024 / 1024).toFixed(1)} MiB`)
    }
  }
  const estimatedBytes = estimateUploadBytes(inputImageDataUrls, maskDataUrl)
  const maxBodyBytes = Math.max(maxUploadBytes * Math.max(1, files.length), 64 * 1024 * 1024)
  if (estimatedBytes > maxBodyBytes) {
    throw new Error(`请求体太大：当前约 ${(estimatedBytes / 1024 / 1024).toFixed(1)} MiB`)
  }
}

async function createJobBody(opts: Sub2APIJobSubmitOptions): Promise<{ body: BodyInit; headers: HeadersInit; path: string }> {
  const { profile, clientTaskId, prompt, params, inputImageDataUrls, maskDataUrl } = opts
  const apiKeyId = selectSub2APIKeyId(profile)
  const common = {
    api_key_id: apiKeyId,
    client_task_id: clientTaskId,
    model: profile.model,
    prompt,
    size: params.size,
    quality: params.quality,
    n: params.n,
    output_format: params.output_format,
  }

  if (inputImageDataUrls.length === 0) {
    return {
      path: '/jobs/images/generations',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(common),
    }
  }

  const formData = new FormData()
  for (const [key, value] of Object.entries(common)) formData.append(key, String(value))
  for (let i = 0; i < inputImageDataUrls.length; i++) {
    const dataUrl = inputImageDataUrls[i]
    const blob = maskDataUrl && i === 0 ? await imageDataUrlToPngBlob(dataUrl) : await dataUrlToBlob(dataUrl)
    const ext = blob.type.split('/')[1] || 'png'
    formData.append('image[]', blob, `input-${i + 1}.${ext}`)
  }
  if (maskDataUrl) {
    formData.append('mask', await maskDataUrlToPngBlob(maskDataUrl), 'mask.png')
  }
  return { path: '/jobs/images/edits', headers: {}, body: formData }
}

export async function createSub2APIImageJob(opts: Sub2APIJobSubmitOptions): Promise<Sub2APIJobResult> {
  assertSub2APIParams(opts.params, getCachedSub2APISettings())
  assertSub2APIUploadLimits(opts.inputImageDataUrls, opts.maskDataUrl)
  const { body, headers, path } = await createJobBody(opts)
  const response = await fetchWithTimeout(path, { method: 'POST', headers, body })
  return readJSON<Sub2APIJobResult>(response)
}

export async function getSub2APIImageJob(jobId: number): Promise<Sub2APIJobResult> {
  const response = await fetchWithTimeout(`/jobs/${encodeURIComponent(String(jobId))}`)
  return readJSON<Sub2APIJobResult>(response)
}

export async function getSub2APIImageJobByClientTask(clientTaskId: string): Promise<Sub2APIJobResult> {
  const response = await fetchWithTimeout(`/jobs/by-client-task/${encodeURIComponent(clientTaskId)}`)
  return readJSON<Sub2APIJobResult>(response)
}

export async function cancelSub2APIImageJob(jobId: number): Promise<Sub2APIJobResult> {
  const response = await fetchWithTimeout(`/jobs/${encodeURIComponent(String(jobId))}/cancel`, { method: 'POST' })
  return readJSON<Sub2APIJobResult>(response)
}

export async function waitForSub2APIImageJob(jobId: number, params: TaskParams): Promise<{ images: string[]; assets: Sub2APIAsset[]; job: Sub2APIJob }> {
  const mime = MIME_MAP[params.output_format] || 'image/png'
  while (true) {
    const result = await getSub2APIImageJob(jobId)
    if (result.job.status === 'succeeded') {
      const assets = result.assets ?? []
      const images: string[] = []
      for (const asset of assets) {
        const url = asset.download_url || asset.url
        images.push(await fetchImageUrlAsDataUrl(url, asset.mime_type || mime))
      }
      if (!images.length) throw new Error('任务已完成，但没有返回可下载的图片。')
      return { images, assets, job: result.job }
    }
    if (result.job.status === 'failed' || result.job.status === 'save_failed' || result.job.status === 'cancelled') {
      throw new Error(result.job.error || `任务${result.job.status}`)
    }
    await sleep(JOB_POLL_MS)
  }
}

export async function callSub2APIImageJob(opts: Sub2APIJobSubmitOptions): Promise<{ images: string[]; assets: Sub2APIAsset[]; job: Sub2APIJob }> {
  let jobId: number | undefined
  try {
    const created = await createSub2APIImageJob(opts)
    jobId = created.job.id
    opts.onJobEnqueued?.({ jobId, status: created.job.status })
    return await waitForSub2APIImageJob(jobId, opts.params)
  } catch (err) {
    if (jobId && isRecoverableSub2APINetworkError(err)) {
      throw new Sub2APIJobRecoverableError('与服务器的连接已断开，任务仍会在服务端继续运行。', {
        jobId,
        clientTaskId: opts.clientTaskId,
      })
    }
    if (!jobId && isRecoverableSub2APINetworkError(err)) {
      try {
        const recovered = await getSub2APIImageJobByClientTask(opts.clientTaskId)
        if (recovered.job.id) {
          throw new Sub2APIJobRecoverableError('提交响应丢失，已通过任务 ID 找回服务端任务。', {
            jobId: recovered.job.id,
            clientTaskId: opts.clientTaskId,
          })
        }
      } catch (recoverErr) {
        if (isSub2APIJobRecoverableError(recoverErr)) throw recoverErr
      }
      throw new Sub2APIJobRecoverableError('上传过程中断，未能确认服务端是否创建任务。请重新上传后再试。', {
        clientTaskId: opts.clientTaskId,
      })
    }
    throw err
  }
}
