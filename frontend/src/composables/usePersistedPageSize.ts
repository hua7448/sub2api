import { getConfiguredTableDefaultPageSize, normalizeTablePageSize } from '@/utils/tablePreferences'

const STORAGE_KEY = 'table-page-size'

function hasInjectedTablePageSizeConfig(): boolean {
  if (typeof window === 'undefined') return false
  const config = window.__APP_CONFIG__
  return Boolean(
    config &&
      (config.table_default_page_size !== undefined ||
        config.table_page_size_options !== undefined)
  )
}

export function getPersistedPageSize(fallback = getConfiguredTableDefaultPageSize()): number {
  if (hasInjectedTablePageSizeConfig()) {
    return normalizeTablePageSize(getConfiguredTableDefaultPageSize() || fallback)
  }

  if (typeof window !== 'undefined') {
    try {
      const stored = window.localStorage.getItem(STORAGE_KEY)
      if (stored !== null) {
        const parsed = Number(stored)
        if (Number.isFinite(parsed)) {
          return normalizeTablePageSize(parsed)
        }
      }
    } catch (error) {
      console.warn('Failed to read persisted page size:', error)
    }
  }
  return normalizeTablePageSize(getConfiguredTableDefaultPageSize() || fallback)
}

export function setPersistedPageSize(size: number): void {
  if (typeof window === 'undefined') return
  try {
    window.localStorage.setItem(STORAGE_KEY, String(size))
  } catch (error) {
    console.warn('Failed to persist page size:', error)
  }
}
