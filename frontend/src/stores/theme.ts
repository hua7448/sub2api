import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export type Theme = 'dark' | 'light'

export const useThemeStore = defineStore('theme', () => {
  // 强制浅色：本项目大量模板使用硬编码的近黑文字/白底（约 760 处 text-gray-900/800/700、
  // 200 处 bg-white）且无 dark: 变体，深色主题目前不可读。在完成全站深色适配之前，
  // 暂时锁定为浅色，并主动清理历史残留的 localStorage='dark'，避免界面不可读。
  const theme = ref<Theme>('light')

  // 清理历史 dark 残留（旧版本默认 dark，会导致老用户一进来就是坏掉的深色界面）
  if (localStorage.getItem('theme') === 'dark') {
    localStorage.setItem('theme', 'light')
  }

  // Apply theme to document.
  const applyTheme = (newTheme: Theme) => {
    document.documentElement.classList.toggle('dark-theme', newTheme === 'dark')
  }

  // Watch theme changes
  watch(theme, (newTheme) => {
    localStorage.setItem('theme', newTheme)
    applyTheme(newTheme)
  }, { immediate: true })

  // 深色主题尚未完成全站适配，切换暂时保持浅色（避免切到不可读的深色）。
  const toggleTheme = () => {
    theme.value = 'light'
  }

  return {
    theme,
    toggleTheme
  }
})
