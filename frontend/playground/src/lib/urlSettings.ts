import type { AppSettings } from '../types'

const URL_SETTING_KEYS = ['settings', 'apiUrl', 'apiKey', 'codexCli', 'apiMode', 'model', 'profileName', 'streamImages', 'streamPartialImages']

export function activateFirstImportedProfile(settings: AppSettings, _importedSettings?: unknown): AppSettings {
  return settings
}

export function hasUrlSettingParams(searchParams: URLSearchParams) {
  return URL_SETTING_KEYS.some((key) => searchParams.has(key))
}

export function clearUrlSettingParams(searchParams: URLSearchParams) {
  for (const key of URL_SETTING_KEYS) searchParams.delete(key)
}

export function buildSettingsFromUrlParams(currentSettings: Partial<AppSettings> | unknown, _searchParams?: URLSearchParams): Partial<AppSettings> {
  return currentSettings && typeof currentSettings === 'object' ? currentSettings as Partial<AppSettings> : {}
}
