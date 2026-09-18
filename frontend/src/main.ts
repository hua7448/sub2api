import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import i18n, { initI18n } from './i18n'
import { useAppStore } from '@/stores/app'
import { updateFavicon } from '@/utils/branding'
import '@fontsource/inter/400.css'
import '@fontsource/inter/500.css'
import '@fontsource/inter/600.css'
import '@fontsource/inter/700.css'
import '@fontsource-variable/jetbrains-mono'
import { isIOSDevice } from '@/utils/device'
import './style.css'
import './styles/cyber-effects.css'
import './styles/enhanced-badges.css'

function initIOSViewportZoomFix() {
  // iOS Safari 在输入框字号小于 16px 时聚焦会自动放大页面，且失焦后不会恢复。
  // 限制 maximum-scale 可阻止该行为；iOS 10+ 用户仍可双指手动缩放，不影响可访问性。
  // 仅在 iOS 设备上注入，避免影响 Android Chrome 的手动缩放能力。
  if (!isIOSDevice()) return

  const viewport = document.querySelector('meta[name="viewport"]')
  if (!viewport) return

  const content = viewport.getAttribute('content') || ''
  if (/maximum-scale/i.test(content)) return
  viewport.setAttribute('content', `${content}, maximum-scale=1.0`)
}

function applyInitialTheme() {
  if (localStorage.getItem('theme') === 'dark') {
    localStorage.setItem('theme', 'light')
  }
  const savedTheme = localStorage.getItem('theme')
  const isDark = savedTheme === 'dark'
  document.documentElement.classList.toggle('dark-theme', isDark)
  document.documentElement.classList.remove('dark')

  const themeColor = document.querySelector<HTMLMetaElement>('meta[name="theme-color"]')
  if (themeColor) {
    themeColor.content = isDark ? '#0f172a' : '#eef4fb'
  }
}

async function bootstrap() {
  applyInitialTheme()
  initIOSViewportZoomFix()

  const app = createApp(App)
  const pinia = createPinia()
  app.use(pinia)

  // Initialize settings from injected config BEFORE mounting (prevents flash)
  // This must happen after pinia is installed but before router and i18n
  const appStore = useAppStore()
  appStore.initFromInjectedConfig()

  // Set document title immediately after config is loaded
  if (appStore.siteName && appStore.siteName !== 'Sub2API') {
    document.title = `${appStore.siteName} - AI API Gateway`
  }
  updateFavicon(appStore.siteLogo)

  await initI18n()

  app.use(router)
  app.use(i18n)

  app.mount('#app')

  router.isReady().catch((error) => {
    console.error('Router initial navigation failed:', error)
  })
}

bootstrap().catch((error) => {
  console.error('Failed to bootstrap frontend:', error)

  const appRoot = document.getElementById('app')
  if (!appRoot) {
    return
  }

  const shell = document.createElement('div')
  shell.style.cssText =
    'min-height:100vh;display:flex;align-items:center;justify-content:center;background:#eef4fb;color:#0b1324;font-family:system-ui,-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif;padding:24px;text-align:center'

  const panel = document.createElement('div')
  const title = document.createElement('div')
  title.textContent = '前端启动失败'
  title.style.cssText = 'font-size:18px;font-weight:700;margin-bottom:8px'

  const message = document.createElement('div')
  message.textContent = '请重新加载页面。'
  message.style.cssText = 'font-size:14px;color:#475569;margin-bottom:16px'

  const button = document.createElement('button')
  button.type = 'button'
  button.textContent = '重新加载'
  button.style.cssText =
    'border:0;border-radius:8px;background:#2563eb;color:white;font-weight:600;padding:10px 16px;cursor:pointer'
  button.addEventListener('click', () => window.location.reload())

  panel.append(title, message, button)
  shell.appendChild(panel)
  appRoot.replaceChildren(shell)
})
