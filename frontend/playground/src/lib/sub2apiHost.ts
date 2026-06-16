type Locale = 'en' | 'zh'

interface HostConfig {
  site_name?: string
  site_logo?: string
}

interface AuthUser {
  role?: string
}

declare global {
  interface Window {
    __APP_CONFIG__?: HostConfig
  }
}

const LOCALE_KEY = 'sub2api_locale'

export function getHostLocale(): Locale {
  const saved = typeof window !== 'undefined' ? window.localStorage.getItem(LOCALE_KEY) : null
  return saved === 'zh' ? 'zh' : 'en'
}

export function hostText(zh: string, en: string): string {
  return getHostLocale() === 'zh' ? zh : en
}

export function getHostSiteName(): string {
  return window.__APP_CONFIG__?.site_name?.trim() || 'Sub2API'
}

export function getHostSiteLogo(): string {
  return window.__APP_CONFIG__?.site_logo?.trim() || '/logo.png'
}

export function getConsolePath(): string {
  try {
    const raw = window.localStorage.getItem('auth_user')
    const user = raw ? JSON.parse(raw) as AuthUser : null
    return user?.role === 'admin' ? '/admin/dashboard' : '/dashboard'
  } catch {
    return '/dashboard'
  }
}

export function applyHostDocumentChrome() {
  const locale = getHostLocale()
  const siteName = getHostSiteName()
  document.documentElement.lang = locale === 'zh' ? 'zh-CN' : 'en'
  document.title = `${siteName} - ${hostText('生图广场', 'Image Playground')}`

  const logo = getHostSiteLogo()
  const iconLinks = document.querySelectorAll<HTMLLinkElement>('link[rel="icon"], link[rel="apple-touch-icon"]')
  iconLinks.forEach((link) => {
    link.href = logo
    link.removeAttribute('type')
  })
}

export function subscribeHostLocaleChange(callback: () => void): () => void {
  const handler = (event: StorageEvent) => {
    if (event.key === LOCALE_KEY) callback()
  }
  window.addEventListener('storage', handler)
  return () => window.removeEventListener('storage', handler)
}
