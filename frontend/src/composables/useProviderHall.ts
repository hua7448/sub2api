import { computed, onBeforeUnmount, ref, watch, type Ref } from 'vue'
import { useRoute, useRouter, type LocationQuery, type LocationQueryRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useAutoRefresh } from '@/composables/useAutoRefresh'
import { extractApiErrorMessage } from '@/utils/apiError'
import * as api from '@/api/providerHall'
import {
  PROVIDER_HALL_RANGES,
  PROVIDER_HALL_SORT_FIELDS,
  type ProviderHallDetailResponse,
  type ProviderHallListResponse,
  type ProviderHallProfileRef,
  type ProviderHallRange,
  type ProviderHallSortField,
  type ProviderHallSortRule,
} from '@/types/providerHall'

export const PROVIDER_HALL_AUTO_REFRESH_INTERVALS = [30, 60, 120, 300] as const
export const PROVIDER_HALL_DEFAULT_PAGE_SIZE = 50
export const PROVIDER_HALL_MAX_SORT_LEVELS = 3

/** Model filter value: `model::protocol` (a profile ref), or '' for all. */
export function modelFilterValue(ref: ProviderHallProfileRef | null | undefined): string {
  if (!ref || !ref.model) return ''
  return `${ref.model}::${ref.protocol}`
}

export function parseModelFilter(value: string | null | undefined): { model: string; protocol: string } | null {
  if (!value) return null
  const idx = value.lastIndexOf('::')
  if (idx <= 0) return { model: value, protocol: '' }
  return { model: value.slice(0, idx), protocol: value.slice(idx + 2) }
}

export function parseSortRules(raw: string | null | undefined): ProviderHallSortRule[] {
  if (!raw) return []
  const rules: ProviderHallSortRule[] = []
  const seen = new Set<string>()
  for (const part of String(raw).split(',')) {
    const [field, dir] = part.trim().split(':')
    if (!PROVIDER_HALL_SORT_FIELDS.includes(field as ProviderHallSortField)) continue
    if (seen.has(field)) continue
    seen.add(field)
    rules.push({ field: field as ProviderHallSortField, direction: dir === 'asc' ? 'asc' : 'desc' })
    if (rules.length >= PROVIDER_HALL_MAX_SORT_LEVELS) break
  }
  return rules
}

export function sortStorageKey(userId: number | string | null | undefined): string {
  return `providerHall.sort.${userId ?? 'anon'}`
}

function first(value: LocationQuery[string]): string {
  if (Array.isArray(value)) return String(value[0] ?? '')
  return value == null ? '' : String(value)
}

function parseRange(value: string, fallback: ProviderHallRange): ProviderHallRange {
  return (PROVIDER_HALL_RANGES as readonly string[]).includes(value) ? (value as ProviderHallRange) : fallback
}

function isCanceled(error: unknown): boolean {
  const e = error as { name?: string; code?: string }
  return e?.name === 'AbortError' || e?.name === 'CanceledError' || e?.code === 'ERR_CANCELED'
}

export interface UseProviderHallOptions {
  defaultRange?: ProviderHallRange
  pageSize?: number
  /** i18n'd fallback error text. */
  loadFailedText?: () => string
  detailFailedText?: () => string
}

export function useProviderHall(options: UseProviderHallOptions = {}) {
  const route = useRoute()
  const router = useRouter()
  const authStore = useAuthStore()
  const userId = computed(() => authStore.user?.id ?? null)
  const defaultRange = options.defaultRange ?? '6h'
  const pageSize = options.pageSize ?? PROVIDER_HALL_DEFAULT_PAGE_SIZE

  // ---- filter state (mirrored into the route query) ----
  const initialQuery = route.query
  const range = ref<ProviderHallRange>(parseRange(first(initialQuery.range), defaultRange))
  const model = ref<string>(
    first(initialQuery.model)
      ? `${first(initialQuery.model)}::${first(initialQuery.protocol)}`
      : ''
  )
  const search = ref<string>(first(initialQuery.search))
  const page = ref<number>(Math.max(1, Number(first(initialQuery.page)) || 1))
  const sort = ref<ProviderHallSortRule[]>(loadInitialSort(first(initialQuery.sort)))

  function loadInitialSort(fromQuery: string): ProviderHallSortRule[] {
    const fromRoute = parseSortRules(fromQuery)
    if (fromRoute.length) return fromRoute
    try {
      return parseSortRules(localStorage.getItem(sortStorageKey(userId.value)))
    } catch {
      return []
    }
  }

  // ---- data state ----
  const data: Ref<ProviderHallListResponse | null> = ref(null)
  const loading = ref(false)
  const refreshing = ref(false)
  const stale = ref(false)
  const error = ref<string | null>(null)
  const lastLoadedAt = ref<Date | null>(null)

  let controller: AbortController | null = null
  let sequence = 0

  async function fetch(silent = false): Promise<void> {
    controller?.abort()
    const request = new AbortController()
    controller = request
    const id = ++sequence
    refreshing.value = true
    if (!silent || !data.value) loading.value = true
    try {
      const parsed = parseModelFilter(model.value)
      const response = await api.list(
        {
          range: range.value,
          model: parsed?.model || undefined,
          protocol: parsed?.protocol || undefined,
          search: search.value.trim() || undefined,
          sort: sort.value,
          page: page.value,
          page_size: pageSize,
        },
        request.signal
      )
      if (id !== sequence) return
      data.value = response
      stale.value = false
      error.value = null
      lastLoadedAt.value = new Date()
    } catch (err) {
      if (isCanceled(err) || id !== sequence) return
      // Keep whatever was shown; flag it as stale so the user knows.
      stale.value = data.value != null
      error.value = extractApiErrorMessage(err, options.loadFailedText?.() ?? 'Failed to load')
    } finally {
      if (id === sequence) {
        loading.value = false
        refreshing.value = false
      }
    }
  }

  // ---- expansion + detail cache ----
  const expanded = ref(new Set<number>())
  const details = ref(new Map<string, ProviderHallDetailResponse>())
  const detailLoading = ref(new Set<number>())
  const detailErrors = ref(new Map<number, string>())
  const detailControllers = new Map<number, AbortController>()

  function detailKey(groupId: number, r: ProviderHallRange = range.value): string {
    return `${groupId}:${r}`
  }

  function getDetail(groupId: number): ProviderHallDetailResponse | null {
    return details.value.get(detailKey(groupId)) ?? null
  }

  async function loadDetail(groupId: number, force = false): Promise<void> {
    const key = detailKey(groupId)
    if (!force && details.value.has(key)) return
    detailControllers.get(groupId)?.abort()
    const request = new AbortController()
    detailControllers.set(groupId, request)
    const r = range.value
    detailLoading.value = new Set(detailLoading.value).add(groupId)
    try {
      const response = await api.getGroup(groupId, r, request.signal)
      if (detailControllers.get(groupId) !== request) return
      const next = new Map(details.value)
      next.set(detailKey(groupId, r), response)
      details.value = next
      const errs = new Map(detailErrors.value)
      errs.delete(groupId)
      detailErrors.value = errs
    } catch (err) {
      if (isCanceled(err) || detailControllers.get(groupId) !== request) return
      const errs = new Map(detailErrors.value)
      errs.set(groupId, extractApiErrorMessage(err, options.detailFailedText?.() ?? 'Failed to load'))
      detailErrors.value = errs
    } finally {
      if (detailControllers.get(groupId) === request) {
        const next = new Set(detailLoading.value)
        next.delete(groupId)
        detailLoading.value = next
        detailControllers.delete(groupId)
      }
    }
  }

  function toggleExpand(groupId: number): void {
    const next = new Set(expanded.value)
    if (next.has(groupId)) {
      next.delete(groupId)
    } else {
      next.add(groupId)
      void loadDetail(groupId)
    }
    expanded.value = next
  }

  function isExpanded(groupId: number): boolean {
    return expanded.value.has(groupId)
  }

  // ---- sort helpers ----
  function persistSort(): void {
    try {
      const serialized = api.serializeSort(sort.value)
      if (serialized) localStorage.setItem(sortStorageKey(userId.value), serialized)
      else localStorage.removeItem(sortStorageKey(userId.value))
    } catch {
      /* ignore quota / privacy mode */
    }
  }

  function setSort(rules: ProviderHallSortRule[]): void {
    const seen = new Set<string>()
    sort.value = rules
      .filter((rule) => {
        if (seen.has(rule.field)) return false
        seen.add(rule.field)
        return true
      })
      .slice(0, PROVIDER_HALL_MAX_SORT_LEVELS)
  }

  /** Single-chip cycle: off → desc → asc → off. A custom (multi-level) sort restarts at desc. */
  function cycleSort(field: ProviderHallSortField): void {
    const current = sort.value.length === 1 ? sort.value[0] : null
    if (!current || current.field !== field) {
      setSort([{ field, direction: 'desc' }])
    } else if (current.direction === 'desc') {
      setSort([{ field, direction: 'asc' }])
    } else {
      setSort([])
    }
  }

  function clearSort(): void {
    setSort([])
  }

  const primarySort = computed(() => sort.value[0] ?? null)
  const isCustomSort = computed(() => sort.value.length > 1)

  // ---- route sync ----
  function buildQuery(): LocationQueryRaw {
    const query: LocationQueryRaw = { ...route.query }
    const parsed = parseModelFilter(model.value)
    const assign = (key: string, value: string | undefined) => {
      if (value) query[key] = value
      else delete query[key]
    }
    assign('range', range.value === defaultRange ? undefined : range.value)
    assign('model', parsed?.model)
    assign('protocol', parsed?.protocol)
    assign('search', search.value.trim() || undefined)
    assign('sort', api.serializeSort(sort.value))
    assign('page', page.value > 1 ? String(page.value) : undefined)
    return query
  }

  let syncingFromRoute = false
  function syncRoute(): void {
    if (syncingFromRoute) return
    const next = buildQuery()
    const current = route.query
    const keys = new Set([...Object.keys(next), ...Object.keys(current)])
    let changed = false
    for (const key of keys) {
      if (String(next[key] ?? '') !== String(first(current[key]))) {
        changed = true
        break
      }
    }
    if (changed) void router.replace({ query: next })
  }

  watch([range, model, search, sort], ([, , , newSort], [oldRange, oldModel, oldSearch]) => {
    // Any filter/sort change resets the page.
    if (range.value !== oldRange || model.value !== oldModel || search.value !== oldSearch) page.value = 1
    if (newSort !== undefined) persistSort()
    syncRoute()
    void fetch(false)
  })

  watch(page, () => {
    syncRoute()
    void fetch(false)
  })

  watch(
    () => route.query,
    (query) => {
      // Browser back/forward: pull state from the route without re-writing it.
      syncingFromRoute = true
      try {
        range.value = parseRange(first(query.range), defaultRange)
        model.value = first(query.model) ? `${first(query.model)}::${first(query.protocol)}` : ''
        search.value = first(query.search)
        page.value = Math.max(1, Number(first(query.page)) || 1)
        const rules = parseSortRules(first(query.sort))
        if (api.serializeSort(rules) !== api.serializeSort(sort.value)) sort.value = rules
      } finally {
        syncingFromRoute = false
      }
    }
  )

  // ---- auto refresh ----
  const autoRefresh = useAutoRefresh({
    storageKey: 'providerHall.autoRefresh',
    intervals: PROVIDER_HALL_AUTO_REFRESH_INTERVALS,
    defaultInterval: 60,
    onRefresh: async () => {
      await fetch(true)
      // Keep open detail panels in step with the list.
      await Promise.all(Array.from(expanded.value).map((id) => loadDetail(id, true)))
    },
    shouldPause: () => (typeof document !== 'undefined' && document.hidden) || loading.value,
  })

  onBeforeUnmount(() => {
    controller?.abort()
    for (const c of detailControllers.values()) c.abort()
    detailControllers.clear()
  })

  const items = computed(() => data.value?.items ?? [])
  const summary = computed(() => data.value?.summary ?? null)
  const catalog = computed(() => data.value?.catalog ?? null)
  const pagination = computed(() => data.value?.pagination ?? { page: page.value, page_size: pageSize, total: 0 })
  const dataThrough = computed(() => data.value?.data_through ?? null)
  const pricingAt = computed(() => data.value?.pricing_at ?? null)

  return {
    range,
    model,
    search,
    page,
    sort,
    pageSize,
    data,
    items,
    summary,
    catalog,
    pagination,
    dataThrough,
    pricingAt,
    loading,
    refreshing,
    stale,
    error,
    lastLoadedAt,
    fetch,
    expanded,
    details,
    detailLoading,
    detailErrors,
    getDetail,
    loadDetail,
    toggleExpand,
    isExpanded,
    setSort,
    cycleSort,
    clearSort,
    primarySort,
    isCustomSort,
    autoRefresh,
  }
}

export type ProviderHallStore = ReturnType<typeof useProviderHall>
