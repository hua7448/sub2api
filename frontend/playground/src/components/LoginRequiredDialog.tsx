import { getHostSiteLogo, getHostSiteName } from '../lib/sub2apiHost'
import { text } from '../lib/playgroundI18n'

interface LoginRequiredDialogProps {
  open: boolean
}

function getLoginHref() {
  const redirect = `${window.location.pathname}${window.location.search}${window.location.hash}`
  return `/login?redirect=${encodeURIComponent(redirect)}`
}

export default function LoginRequiredDialog({ open }: LoginRequiredDialogProps) {
  if (!open) return null

  const siteName = getHostSiteName()
  const siteLogo = getHostSiteLogo()

  return (
    <div
      data-no-drag-select
      className="fixed inset-0 z-[130] flex items-center justify-center p-4"
    >
      <div className="absolute inset-0 bg-stone-950/35 backdrop-blur-md animate-overlay-in" />
      <div className="relative z-10 w-full max-w-sm rounded-2xl border border-amber-200/80 bg-stone-50/95 p-6 shadow-[0_18px_60px_rgba(68,45,18,0.22)] ring-1 ring-amber-900/5 animate-confirm-in dark:border-amber-400/20 dark:bg-stone-950/95 dark:ring-white/10">
        <div className="mb-4 flex items-center gap-3">
          <img src={siteLogo} alt="" className="h-9 w-9 rounded-lg object-contain" />
          <div className="min-w-0">
            <div className="truncate text-sm font-semibold text-stone-900 dark:text-stone-100">{siteName}</div>
            <div className="text-xs text-stone-500 dark:text-stone-400">
              {text('生图广场', 'Image Playground')}
            </div>
          </div>
        </div>
        <h2 className="text-lg font-semibold text-stone-900 dark:text-stone-100">
          {text('请先登录主站', 'Sign in required')}
        </h2>
        <p className="mt-2 text-sm leading-6 text-stone-600 dark:text-stone-300">
          {text('当前没有检测到登录状态。生图广场会使用主站账号和 API Key 权限，请先回到主站登录后再继续。', 'No active session was found. Image Playground uses your main-site account and API key permissions, so sign in first to continue.')}
        </p>
        <div className="mt-6 flex gap-2">
          <a
            href={getLoginHref()}
            className="flex-1 rounded-lg bg-amber-600 px-4 py-2.5 text-center text-sm font-semibold text-white transition hover:bg-amber-700 focus:outline-none focus:ring-2 focus:ring-amber-500/35"
          >
            {text('去登录', 'Sign in')}
          </a>
          <a
            href="/home"
            className="rounded-lg border border-stone-200 bg-white px-4 py-2.5 text-sm font-medium text-stone-600 transition hover:bg-stone-50 hover:text-stone-900 dark:border-white/10 dark:bg-white/[0.04] dark:text-stone-300 dark:hover:bg-white/[0.08] dark:hover:text-white"
          >
            {text('主站首页', 'Home')}
          </a>
        </div>
      </div>
    </div>
  )
}
