<template>
  <header class="glass sticky top-0 z-30 border-b border-blue-100">
    <div class="flex h-16 items-center justify-between gap-2 px-2 sm:px-4 md:px-6">
      <!-- Left: Mobile Menu Toggle + Page Title -->
      <div class="flex shrink-0 items-center gap-2 sm:gap-4">
        <button
          @click="toggleMobileSidebar"
          class="btn-ghost btn-icon lg:hidden"
          :aria-label="t('common.toggleMenu')"
        >
          <Menu class="h-5 w-5" :stroke-width="1.5" />
        </button>

        <div class="hidden lg:block">
          <h1 class="text-lg font-semibold text-gray-900">
            {{ pageTitle }}
          </h1>
          <p v-if="pageDescription" class="text-xs text-gray-600">
            {{ pageDescription }}
          </p>
        </div>
      </div>

      <!-- Right: Announcements + Docs + Language + Subscriptions + Balance + User Dropdown -->
      <div class="flex min-w-0 items-center gap-1 sm:gap-3">
        <!-- Announcement Bell -->
        <AnnouncementBell v-if="user" />

        <!-- Docs Link -->
        <a
          v-if="docUrl"
          :href="docUrl"
          target="_blank"
          rel="noopener noreferrer"
          class="hidden items-center gap-1.5 rounded-lg px-2.5 py-1.5 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-100 hover:text-blue-600 sm:flex"
        >
          <BookOpen class="h-4 w-4" :stroke-width="1.5" />
          <span class="hidden sm:inline">{{ t('nav.docs') }}</span>
        </a>

        <!-- Theme Switcher -->
        <button
          @click="toggleTheme"
          class="hidden items-center gap-1.5 rounded-lg px-2.5 py-1.5 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-100 hover:text-blue-600 sm:flex"
          :title="themeStore.theme === 'dark' ? t('nav.lightMode') : t('nav.darkMode')"
        >
          <Sun v-if="themeStore.theme === 'dark'" class="h-4 w-4" :stroke-width="1.5" />
          <Moon v-else class="h-4 w-4" :stroke-width="1.5" />
        </button>

        <!-- Model Plaza Entry -->
        <router-link
          v-if="user && modelPlazaEnabled"
          :to="{ path: '/model-plaza', query: { embedded: '1' } }"
          class="hidden items-center gap-1.5 rounded-lg px-2.5 py-1.5 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white sm:flex"
        >
          <Icon name="grid" size="sm" />
          <span class="hidden sm:inline">{{ t('nav.modelPlaza') }}</span>
        </router-link>

        <!-- Language Switcher -->
        <LocaleSwitcher />

        <!-- Subscription Progress -->
        <SubscriptionProgressMini v-if="user" />

        <!-- Balance Display -->
        <div
          v-if="user"
          class="hidden items-center gap-2 rounded-xl border border-blue-200 bg-blue-50 px-3 py-1.5 sm:flex"
        >
          <Wallet class="h-4 w-4 text-blue-600" :stroke-width="1.5" />
          <span class="text-sm font-semibold font-mono text-blue-700">
            ${{ user.balance?.toFixed(2) || '0.00' }}
          </span>
          <span
            v-if="frozenBalance > 0"
            class="rounded-full bg-amber-100 px-1.5 py-0.5 text-xs font-medium text-amber-700 dark:bg-amber-900/40 dark:text-amber-200"
          >
            {{ balanceFrozenLabel }}
          </span>
          <div
            class="pointer-events-none absolute right-0 top-full mt-2 hidden w-56 rounded-lg border border-gray-200 bg-white p-3 text-xs shadow-lg group-hover:block dark:border-dark-700 dark:bg-dark-800"
          >
            <div class="flex items-center justify-between">
              <span class="text-gray-500 dark:text-dark-400">{{ balanceAvailableText }}</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ formatHeaderMoney(availableBalance) }}</span>
            </div>
            <div class="mt-2 flex items-center justify-between">
              <span class="text-gray-500 dark:text-dark-400">{{ balanceFrozenText }}</span>
              <span class="font-medium text-amber-700 dark:text-amber-200">{{ formatHeaderMoney(frozenBalance) }}</span>
            </div>
            <div class="mt-2 border-t border-gray-100 pt-2 dark:border-dark-700">
              <div class="flex items-center justify-between">
                <span class="text-gray-500 dark:text-dark-400">{{ balanceTotalText }}</span>
                <span class="font-semibold text-gray-900 dark:text-white">{{ formatHeaderMoney(totalBalance) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- User Dropdown -->
        <div v-if="user" class="relative" ref="dropdownRef">
          <button
            @click="toggleDropdown"
            class="flex items-center gap-2 rounded-xl p-1.5 transition-colors hover:bg-gray-100 dark:hover:bg-dark-800"
            :aria-label="t('common.userMenu')"
          >
            <div
              class="flex h-8 w-8 items-center justify-center overflow-hidden rounded-xl bg-gradient-to-br from-primary-500 to-primary-600 text-sm font-medium text-white shadow-lg shadow-primary-500/30"
            >
              <img
                v-if="avatarUrl"
                :src="avatarUrl"
                :alt="displayName"
                class="h-full w-full object-cover"
              >
              <span v-else>{{ userInitials }}</span>
            </div>
            <div class="hidden text-left md:block">
              <div class="text-sm font-medium text-gray-900">
                {{ displayName }}
              </div>
              <div class="text-xs text-txt-secondary">
                {{ t('admin.users.roles.' + user.role) }}
              </div>
            </div>
            <ChevronDown class="hidden h-4 w-4 text-gray-500 md:block" :stroke-width="1.5" />
          </button>

          <!-- Dropdown Menu -->
          <transition name="dropdown">
            <div v-if="dropdownOpen" class="dropdown right-0 mt-2 w-56">
              <!-- User Info -->
              <div class="border-b border-gray-200 px-4 py-3">
                <div class="text-sm font-medium text-gray-900">
                  {{ displayName }}
                </div>
                <div class="text-xs text-gray-600">{{ user.email }}</div>
              </div>

              <!-- Balance (mobile only) -->
              <div class="border-b border-gray-200 px-4 py-2 sm:hidden">
                <div class="text-xs text-gray-600">
                  {{ t('common.balance') }}
                </div>
                <div class="text-sm font-semibold font-mono text-blue-700">
                  ${{ user.balance?.toFixed(2) || '0.00' }}
                </div>
              </div>

              <div class="py-1">
                <router-link to="/profile" @click="closeDropdown" class="dropdown-item">
                  <UserIcon class="h-4 w-4" :stroke-width="1.5" />
                  {{ t('nav.profile') }}
                </router-link>

                <router-link to="/keys" @click="closeDropdown" class="dropdown-item">
                  <KeyIcon class="h-4 w-4" :stroke-width="1.5" />
                  {{ t('nav.apiKeys') }}
                </router-link>

              </div>

              <!-- Contact Support -->
              <div
                v-if="contactInfo"
                class="border-t border-cyber-border px-4 py-2.5"
              >
                <div class="flex items-center gap-2 text-xs text-txt-secondary">
                  <MessageSquare class="h-3.5 w-3.5 flex-shrink-0" :stroke-width="1.5" />
                  <span>{{ t('common.contactSupport') }}:</span>
                  <span class="font-medium text-txt-primary">{{ contactInfo }}</span>
                </div>
              </div>

              <div v-if="showOnboardingButton" class="border-t border-cyber-border py-1">
                <button @click="handleReplayGuide" class="dropdown-item w-full">
                  <HelpCircle class="h-4 w-4" :stroke-width="1.5" />
                  {{ $t('onboarding.restartTour') }}
                </button>
              </div>

              <div class="border-t border-cyber-border py-1">
                <button
                  @click="handleLogout"
                  class="dropdown-item w-full text-neon-pink hover:bg-neon-pink/5"
                >
                  <LogOut class="h-4 w-4" :stroke-width="1.5" />
                  {{ t('nav.logout') }}
                </button>
              </div>
            </div>
          </transition>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAppStore, useAuthStore, useOnboardingStore } from '@/stores'
import { useAdminSettingsStore } from '@/stores/adminSettings'
import { useThemeStore } from '@/stores/theme'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import SubscriptionProgressMini from '@/components/common/SubscriptionProgressMini.vue'
import AnnouncementBell from '@/components/common/AnnouncementBell.vue'
import Icon from '@/components/icons/Icon.vue'
import { sanitizeUrl } from '@/utils/url'
import { FeatureFlags, isFeatureFlagEnabled } from '@/utils/featureFlags'
import {
  Menu,
  BookOpen,
  Wallet,
  ChevronDown,
  User as UserIcon,
  Key as KeyIcon,
  MessageSquare,
  HelpCircle,
  LogOut,
  Sun,
  Moon
} from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()
const adminSettingsStore = useAdminSettingsStore()
const onboardingStore = useOnboardingStore()
const themeStore = useThemeStore()

const user = computed(() => authStore.user)
const dropdownOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)
const contactInfo = computed(() => appStore.contactInfo)
const docUrl = computed(() => sanitizeUrl(appStore.docUrl))
const modelPlazaEnabled = computed(() => isFeatureFlagEnabled(FeatureFlags.modelPlaza))
const avatarUrl = computed(() => user.value?.avatar_url?.trim() || '')
const availableBalance = computed(() => Number(user.value?.balance || 0))
const frozenBalance = computed(() => Number(user.value?.frozen_balance || 0))
const totalBalance = computed(() => availableBalance.value + frozenBalance.value)
const balanceAvailableText = computed(() => t('common.availableBalance') === 'common.availableBalance' ? '可用余额' : t('common.availableBalance'))
const balanceFrozenText = computed(() => t('common.frozenBalance') === 'common.frozenBalance' ? '冻结金额' : t('common.frozenBalance'))
const balanceTotalText = computed(() => t('common.totalBalance') === 'common.totalBalance' ? '总余额' : t('common.totalBalance'))
const balanceFrozenLabel = computed(() => `${balanceFrozenText.value} ${formatHeaderMoney(frozenBalance.value)}`)

const showOnboardingButton = computed(() => {
  return !authStore.isSimpleMode && user.value?.role === 'admin'
})

const userInitials = computed(() => {
  if (!user.value) return ''
  if (user.value.username) {
    return user.value.username.substring(0, 2).toUpperCase()
  }
  if (user.value.email) {
    const localPart = user.value.email.split('@')[0]
    return localPart.substring(0, 2).toUpperCase()
  }
  return ''
})

const displayName = computed(() => {
  if (!user.value) return ''
  return user.value.username || user.value.email?.split('@')[0] || ''
})

const pageTitle = computed(() => {
  if (route.name === 'CustomPage') {
    const id = route.params.id as string
    const publicItems = appStore.cachedPublicSettings?.custom_menu_items ?? []
    const menuItem = publicItems.find((item) => item.id === id)
      ?? (authStore.isAdmin ? adminSettingsStore.customMenuItems.find((item) => item.id === id) : undefined)
    if (menuItem?.label) return menuItem.label
  }
  const titleKey = route.meta.titleKey as string
  if (titleKey) {
    return t(titleKey)
  }
  return (route.meta.title as string) || ''
})

const pageDescription = computed(() => {
  const descKey = route.meta.descriptionKey as string
  if (descKey) {
    return t(descKey)
  }
  return (route.meta.description as string) || ''
})

function toggleMobileSidebar() {
  appStore.toggleMobileSidebar()
}

function toggleTheme() {
  themeStore.toggleTheme()
}

function toggleDropdown() {
  dropdownOpen.value = !dropdownOpen.value
}

function closeDropdown() {
  dropdownOpen.value = false
}

async function handleLogout() {
  closeDropdown()
  try {
    await authStore.logout()
  } catch (error) {
    console.error('Logout error:', error)
  }
  await router.push('/login')
}

function handleReplayGuide() {
  closeDropdown()
  onboardingStore.replay()
}

function formatHeaderMoney(value: number) {
  if (!Number.isFinite(value)) return '$0.00'
  return `$${value.toFixed(2)}`
}

function handleClickOutside(event: MouseEvent) {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    closeDropdown()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: scale(0.95) translateY(-4px);
}
</style>
