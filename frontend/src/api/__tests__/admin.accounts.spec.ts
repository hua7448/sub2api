import { beforeEach, describe, expect, it, vi } from 'vitest'

const { get } = vi.hoisted(() => ({
  get: vi.fn(),
}))

vi.mock('@/api/client', () => ({
  apiClient: {
    get,
  },
}))

import { list, listWithEtag } from '@/api/admin/accounts'

describe('admin accounts API', () => {
  beforeEach(() => {
    get.mockReset()
    get.mockResolvedValue({ data: { items: [], total: 0, pages: 0 } })
  })

  it('omits empty account list filters', async () => {
    await list(2, 50, {
      platform: '',
      type: '',
      status: '',
      search: '   ',
      created_date: '',
      sort_by: 'name',
      sort_order: 'asc',
    })

    expect(get).toHaveBeenCalledWith('/admin/accounts', expect.objectContaining({
      params: {
        page: 2,
        page_size: 50,
        sort_by: 'name',
        sort_order: 'asc',
      },
    }))
  })

  it('keeps selected account type and created date filters', async () => {
    await listWithEtag(1, 20, {
      type: 'oauth',
      created_date: '2026-06-09',
    })

    expect(get).toHaveBeenCalledWith('/admin/accounts', expect.objectContaining({
      params: {
        page: 1,
        page_size: 20,
        type: 'oauth',
        created_date: '2026-06-09',
      },
    }))
  })

  it('keeps selected plan type filters', async () => {
    await list(1, 20, {
      plan_type: 'pro',
    })

    expect(get).toHaveBeenCalledWith('/admin/accounts', expect.objectContaining({
      params: {
        page: 1,
        page_size: 20,
        plan_type: 'pro',
      },
    }))
  })
})
