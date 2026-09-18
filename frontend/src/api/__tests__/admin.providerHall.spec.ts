import { beforeEach, describe, expect, it, vi } from 'vitest'

const { get, put, post } = vi.hoisted(() => ({ get: vi.fn(), put: vi.fn(), post: vi.fn() }))

vi.mock('@/api/client', () => ({ apiClient: { get, put, post } }))

import { getConfig, updateConfig, getTargets, updateTargets, enqueueProbe, enqueueVerification, listJobs, getJob, cancelJob, getHealth, listGroupVerifications, type ProviderHallConfigInput, type ProviderHallTargetSet } from '@/api/admin/providerHall'

describe('provider hall configuration contract', () => {
  const input: ProviderHallConfigInput = {
    version: 2,
    collection_enabled: false,
    display_enabled: false,
    tasks_enabled: false,
    default_model: '',
    default_protocol: 'responses',
    default_range: '6h',
    gateway_origin: '',
    operator_user_id: null,
    daily_budget: '999999999999.00000001',
    expected_nodes: []
  }

  beforeEach(() => vi.resetAllMocks())

  it('preserves monetary precision and optimistic version on write', async () => {
    put.mockResolvedValue({ data: { ...input, version: 3 } })
    const result = await updateConfig(input)
    expect(put).toHaveBeenCalledWith('/admin/provider-hall/config', input)
    expect(JSON.parse(JSON.stringify(put.mock.calls[0][1])).daily_budget).toBe('999999999999.00000001')
    expect(result.version).toBe(3)
  })

  it('propagates conflicts without changing the unsaved input', async () => {
    const conflict = { response: { status: 409, data: { reason: 'PROVIDER_HALL_VERSION_CONFLICT' } } }
    put.mockRejectedValue(conflict)
    await expect(updateConfig(input)).rejects.toBe(conflict)
    expect(input.version).toBe(2)
    expect(input.daily_budget).toBe('999999999999.00000001')
  })

  it('passes cancellation to the shared client', async () => {
    const controller = new AbortController()
    get.mockResolvedValue({ data: input })
    await getConfig(controller.signal)
    expect(get).toHaveBeenCalledWith('/admin/provider-hall/config', { signal: controller.signal })
  })

  it('writes target edits without server-owned IDs and audit metadata', async () => {
    const set: ProviderHallTargetSet = {
      group_id: 4, version: 9, updated_at: '2026-09-10T00:00:00Z', updated_by: 3,
      items: [{ id: 15, group_id: 4, version: 2, updated_at: '2026-09-10T00:00:00Z', updated_by: 3,
        profile_id: 8, probe_key_id: 11, enabled: false, probe_interval_seconds: 300, verification_interval_seconds: 86400 }]
    }
    put.mockResolvedValue({ data: set })
    await updateTargets(4, set)
    expect(put).toHaveBeenCalledWith('/admin/provider-hall/groups/4/targets', {
      version: 9, items: [{ profile_id: 8, probe_key_id: 11, enabled: false, probe_interval_seconds: 300, verification_interval_seconds: 86400 }]
    })
    expect(set.items[0].id).toBe(15)
  })

  it('sends an explicit empty array to disable all targets', async () => {
    put.mockResolvedValue({ data: { items: [] } })
    await updateTargets(4, { version: 9, items: [] })
    expect(put).toHaveBeenCalledWith('/admin/provider-hall/groups/4/targets', { version: 9, items: [] })
  })

  it('cancels target reads through the shared HTTP client', async () => {
    const controller = new AbortController()
    get.mockResolvedValue({ data: { items: [] } })
    await getTargets(4, controller.signal)
    expect(get).toHaveBeenCalledWith('/admin/provider-hall/groups/4/targets', { signal: controller.signal })
  })
})

describe('provider hall task contract', () => {
  beforeEach(() => vi.resetAllMocks())

  it('enqueues probes and verifications with the click idempotency key', async () => {
    post.mockResolvedValue({ data: { job_id: 42, status: 'queued', reused: false } })
    const result = await enqueueProbe(4, 8, 'click-1')
    expect(post).toHaveBeenCalledWith('/admin/provider-hall/groups/4/probes', { profile_id: 8, idempotency_key: 'click-1' })
    expect(result.job_id).toBe(42)
    await enqueueVerification(4, 8, 'click-2')
    expect(post).toHaveBeenLastCalledWith('/admin/provider-hall/groups/4/verifications', { profile_id: 8, idempotency_key: 'click-2' })
  })

  it('sends only set filters when listing jobs', async () => {
    get.mockResolvedValue({ data: { items: [], total: 0, page: 1, page_size: 20, pages: 0 } })
    const controller = new AbortController()
    await listJobs({ status: '', kind: 'probe', group_id: null, page: 2, page_size: 50 }, controller.signal)
    expect(get).toHaveBeenCalledWith('/admin/provider-hall/jobs', { params: { kind: 'probe', page: 2, page_size: 50 }, signal: controller.signal })
  })

  it('cancels, reads detail, health and group reports', async () => {
    post.mockResolvedValue({ data: { job_id: 1, status: 'cancelled' } })
    await cancelJob(1)
    expect(post).toHaveBeenCalledWith('/admin/provider-hall/jobs/1/cancel')
    get.mockResolvedValue({ data: {} })
    await getJob(1)
    expect(get).toHaveBeenCalledWith('/admin/provider-hall/jobs/1', { signal: undefined })
    await getHealth()
    expect(get).toHaveBeenCalledWith('/admin/provider-hall/health', { signal: undefined })
    await listGroupVerifications(4, { profile_id: 8, page: 1, page_size: 10 })
    expect(get).toHaveBeenLastCalledWith('/admin/provider-hall/groups/4/verifications', { params: { profile_id: 8, page: 1, page_size: 10 }, signal: undefined })
  })
})
