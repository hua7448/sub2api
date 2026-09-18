import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { defineComponent, h, nextTick } from 'vue'
import { flushPromises, mount } from '@vue/test-utils'
import { createMemoryHistory, createRouter, type Router } from 'vue-router'
import * as api from '@/api/providerHall'
import { useProviderHall, sortStorageKey, type ProviderHallStore } from '@/composables/useProviderHall'
import { makeList, makeRow } from '@/components/provider-hall/__tests__/testUtils'

vi.mock('@/stores/auth', () => ({ useAuthStore: () => ({ user: { id: 42 } }) }))
vi.mock('@/api/providerHall', async () => {
  const actual = await vi.importActual<typeof import('@/api/providerHall')>('@/api/providerHall')
  return { ...actual, list: vi.fn(), getGroup: vi.fn(), listVerifications: vi.fn() }
})

function deferred<T>() {
  let resolve!: (value: T) => void
  let reject!: (reason: unknown) => void
  const promise = new Promise<T>((yes, no) => {
    resolve = yes
    reject = no
  })
  return { promise, resolve, reject }
}

async function setup(initialPath = '/providers') {
  const router: Router = createRouter({
    history: createMemoryHistory(),
    routes: [{ path: '/providers', component: { template: '<div />' } }],
  })
  await router.push(initialPath)
  await router.isReady()
  let store!: ProviderHallStore
  const Host = defineComponent({
    setup() {
      store = useProviderHall()
      return () => h('div')
    },
  })
  const wrapper = mount(Host, { global: { plugins: [router] } })
  return { wrapper, router, store }
}

beforeEach(() => {
  vi.clearAllMocks()
  localStorage.clear()
  vi.mocked(api.list).mockResolvedValue(makeList([makeRow({ group_id: 1 })]))
})

afterEach(() => {
  vi.useRealTimers()
  Object.defineProperty(document, 'hidden', { value: false, configurable: true })
})

describe('useProviderHall', () => {
  it('writes filters and sort into the route query and resets the page', async () => {
    const { store, router, wrapper } = await setup('/providers?page=3')
    expect(store.page.value).toBe(3)
    store.range.value = '24h'
    store.model.value = 'astra::responses'
    store.search.value = 'pro'
    await flushPromises()
    expect(router.currentRoute.value.query).toMatchObject({ range: '24h', model: 'astra', protocol: 'responses', search: 'pro' })
    expect(router.currentRoute.value.query.page).toBeUndefined()
    expect(store.page.value).toBe(1)
    expect(api.list).toHaveBeenLastCalledWith(
      expect.objectContaining({ range: '24h', model: 'astra', protocol: 'responses', search: 'pro', page: 1 }),
      expect.any(AbortSignal)
    )
    store.range.value = '6h'
    await flushPromises()
    expect(router.currentRoute.value.query.range).toBeUndefined()
    wrapper.unmount()
  })

  it('persists up to three sort levels per user and restores them without a query', async () => {
    const first = await setup()
    first.store.setSort([
      { field: 'rate', direction: 'desc' },
      { field: 'cache_rate', direction: 'asc' },
      { field: 'success_rate', direction: 'desc' },
      { field: 'ttft_fast95', direction: 'asc' },
      { field: 'rate', direction: 'asc' },
    ])
    await flushPromises()
    expect(first.store.sort.value).toHaveLength(3)
    expect(localStorage.getItem(sortStorageKey(42))).toBe('rate:desc,cache_rate:asc,success_rate:desc')
    expect(first.router.currentRoute.value.query.sort).toBe('rate:desc,cache_rate:asc,success_rate:desc')
    expect(api.list).toHaveBeenLastCalledWith(expect.objectContaining({ sort: first.store.sort.value }), expect.any(AbortSignal))
    first.wrapper.unmount()

    const second = await setup()
    expect(second.store.sort.value).toEqual([
      { field: 'rate', direction: 'desc' },
      { field: 'cache_rate', direction: 'asc' },
      { field: 'success_rate', direction: 'desc' },
    ])
    expect(second.store.isCustomSort.value).toBe(true)
    second.store.cycleSort('rate')
    expect(second.store.sort.value).toEqual([{ field: 'rate', direction: 'desc' }])
    second.store.cycleSort('rate')
    expect(second.store.sort.value).toEqual([{ field: 'rate', direction: 'asc' }])
    second.store.cycleSort('rate')
    expect(second.store.sort.value).toEqual([])
    await flushPromises()
    expect(localStorage.getItem(sortStorageKey(42))).toBeNull()
    second.wrapper.unmount()
  })

  it('prefers the route sort over the persisted one', async () => {
    localStorage.setItem(sortStorageKey(42), 'rate:desc')
    const { store, wrapper } = await setup('/providers?sort=success_rate:asc,bogus:desc')
    expect(store.sort.value).toEqual([{ field: 'success_rate', direction: 'asc' }])
    wrapper.unmount()
  })

  it('discards a late response from a superseded request', async () => {
    const slow = deferred<ReturnType<typeof makeList>>()
    const fast = deferred<ReturnType<typeof makeList>>()
    vi.mocked(api.list).mockReturnValueOnce(slow.promise).mockReturnValueOnce(fast.promise)
    const { store, wrapper } = await setup()
    void store.fetch(false)
    store.range.value = '7d'
    await nextTick()
    await flushPromises()
    fast.resolve(makeList([makeRow({ group_id: 7, name: 'Fast' })]))
    await flushPromises()
    slow.resolve(makeList([makeRow({ group_id: 1, name: 'Slow' })]))
    await flushPromises()
    expect(store.items.value.map((r) => r.name)).toEqual(['Fast'])
    expect(store.loading.value).toBe(false)
    wrapper.unmount()
  })

  it('keeps the previous data and flags it stale when a refresh fails', async () => {
    const { store, wrapper } = await setup()
    await store.fetch(false)
    expect(store.items.value).toHaveLength(1)
    vi.mocked(api.list).mockRejectedValueOnce({ message: 'boom' })
    await store.fetch(true)
    expect(store.items.value).toHaveLength(1)
    expect(store.stale.value).toBe(true)
    expect(store.error.value).toBe('boom')
    await store.fetch(true)
    expect(store.stale.value).toBe(false)
    expect(store.error.value).toBeNull()
    wrapper.unmount()
  })

  it('pauses auto refresh while the document is hidden', async () => {
    vi.useFakeTimers()
    const { store, wrapper } = await setup()
    await store.fetch(false)
    expect(api.list).toHaveBeenCalledTimes(1)
    Object.defineProperty(document, 'hidden', { value: true, configurable: true })
    store.autoRefresh.setInterval(30)
    store.autoRefresh.setEnabled(true)
    store.autoRefresh.start()
    await vi.advanceTimersByTimeAsync(40_000)
    expect(api.list).toHaveBeenCalledTimes(1)
    Object.defineProperty(document, 'hidden', { value: false, configurable: true })
    await vi.advanceTimersByTimeAsync(31_000)
    expect(api.list).toHaveBeenCalledTimes(2)
    wrapper.unmount()
  })

  it('caches details per group and range and reloads open panels on refresh', async () => {
    const row = makeRow({ group_id: 1 })
    vi.mocked(api.getGroup).mockResolvedValue({ ...row, detail: {} as never, trend: { range: '6h', bucket_seconds: 300, points: [] }, profiles: [] })
    const { store, wrapper } = await setup()
    store.toggleExpand(1)
    await flushPromises()
    expect(api.getGroup).toHaveBeenCalledWith(1, '6h', expect.any(AbortSignal))
    expect(store.getDetail(1)).not.toBeNull()
    store.toggleExpand(1)
    store.toggleExpand(1)
    await flushPromises()
    expect(api.getGroup).toHaveBeenCalledTimes(1)
    store.range.value = '24h'
    await flushPromises()
    expect(store.getDetail(1)).toBeNull()
    await store.loadDetail(1)
    expect(api.getGroup).toHaveBeenLastCalledWith(1, '24h', expect.any(AbortSignal))
    wrapper.unmount()
  })
})
