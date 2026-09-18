<template>
  <div class="relative">
    <!-- Admin: Full version badge with dropdown -->
    <template v-if="isAdmin">
      <button
        @click="toggleDropdown"
        class="flex items-center gap-1.5 rounded-lg px-2 py-1 text-xs transition-colors"
        :class="[
          hasUpdate
            ? 'bg-amber-100 text-amber-700 hover:bg-amber-200 bg-amber-900/30 text-amber-400 hover:bg-amber-900/50'
            : 'bg-cyber-elevated text-txt-secondary hover:bg-gray-200'
        ]"
        :title="hasUpdate ? t('version.updateAvailable') : t('version.upToDate')"
      >
        <span class="font-medium">
          {{ hasUpdate ? t('version.updateAvailable') : t('version.upToDate') }}
        </span>
        <!-- Update indicator -->
        <span v-if="hasUpdate" class="relative flex h-2 w-2">
          <span
            class="absolute inline-flex h-full w-full animate-ping rounded-full bg-amber-400 opacity-75"
          ></span>
          <span class="relative inline-flex h-2 w-2 rounded-full bg-amber-500"></span>
        </span>
      </button>

      <!-- Dropdown -->
      <transition name="dropdown">
        <div
          v-if="dropdownOpen"
          ref="dropdownRef"
          class="absolute left-0 z-50 mt-2 w-64 overflow-hidden rounded-xl border border-cyber-border bg-cyber-surface shadow-lg"
        >
          <!-- Header with refresh button -->
          <div class="flex items-center justify-end border-b border-gray-100 px-4 py-3">
            <button
              @click="refreshVersion(true)"
              class="rounded-lg p-1.5 text-txt-dim transition-colors hover:bg-cyber-elevated hover:text-txt-secondary"
              :disabled="loading"
              :title="t('version.refresh')"
            >
              <Icon
                name="refresh"
                size="sm"
                :stroke-width="2"
                :class="{ 'animate-spin': loading }"
              />
            </button>
          </div>

          <div class="p-4">
            <!-- Loading state -->
            <div v-if="loading" class="flex items-center justify-center py-6">
              <svg class="h-6 w-6 animate-spin text-neon-cyan" fill="none" viewBox="0 0 24 24">
                <circle
                  class="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  stroke-width="4"
                ></circle>
                <path
                  class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
            </div>

            <!-- Content -->
            <template v-else>
              <!-- Version status - centered and prominent -->
              <div class="mb-4 text-center">
                <div class="inline-flex items-center gap-2">
                  <span class="text-sm font-semibold text-txt-primary">
                    {{ hasUpdate ? t('version.updateAvailable') : t('version.upToDate') }}
                  </span>
                  <!-- Show check mark when up to date -->
                  <span
                    v-if="!hasUpdate"
                    class="flex h-5 w-5 items-center justify-center rounded-full bg-neon-green/10 bg-green-900/30"
                  >
                    <svg
                      class="h-3 w-3 text-neon-green text-green-400"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fill-rule="evenodd"
                        d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                        clip-rule="evenodd"
                      />
                    </svg>
                  </span>
                </div>
              </div>

              <!-- Priority 1: Update error (must check before hasUpdate) -->
              <div v-if="updateError" class="space-y-2">
                <div
                  class="flex items-center gap-3 rounded-lg border border-red-200 bg-neon-pink/10 p-3 border-red-800/50 bg-red-900/20"
                >
                  <div
                    class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full bg-neon-pink/10 bg-red-900/50"
                  >
                    <Icon
                      name="x"
                      size="sm"
                      :stroke-width="2"
                      class="text-neon-pink text-red-400"
                    />
                  </div>
                  <div class="min-w-0 flex-1">
                    <p class="text-sm font-medium text-red-700 text-red-300">
                      {{ t('version.updateFailed') }}
                    </p>
                    <p class="truncate text-xs text-neon-pink/70 text-red-400/70">
                      {{ updateError }}
                    </p>
                  </div>
                </div>

                <!-- Retry button -->
                <button
                  @click="handleUpdate"
                  :disabled="updating"
                  class="flex w-full items-center justify-center gap-2 rounded-lg bg-neon-pink/100 px-4 py-2 text-sm font-medium transition-colors hover:bg-red-600 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  {{ t('version.retry') }}
                </button>
              </div>

              <!-- Priority 2: Update success - need restart -->
              <div v-else-if="updateSuccess && needRestart" class="space-y-2">
                <div
                  class="flex items-center gap-3 rounded-lg border border-green-200 bg-neon-green/10 p-3 border-green-800/50 bg-green-900/20"
                >
                  <div
                    class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full bg-neon-green/10 bg-green-900/50"
                  >
                    <svg
                      class="h-4 w-4 text-neon-green text-green-400"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                      stroke-width="2"
                    >
                      <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                    </svg>
                  </div>
                  <div class="min-w-0 flex-1">
                    <p class="text-sm font-medium text-green-700 text-green-300">
                      {{ t('version.updateComplete') }}
                    </p>
                    <p class="text-xs text-neon-green/70 text-green-400/70">
                      {{ t('version.restartRequired') }}
                    </p>
                  </div>
                </div>

                <!-- Restart button with countdown -->
                <button
                  @click="handleRestart"
                  :disabled="restarting"
                  class="flex w-full items-center justify-center gap-2 rounded-lg bg-neon-green/100 px-4 py-2 text-sm font-medium transition-colors hover:bg-green-600 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  <svg
                    v-if="restarting"
                    class="h-4 w-4 animate-spin"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <circle
                      class="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      stroke-width="4"
                    ></circle>
                    <path
                      class="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    ></path>
                  </svg>
                  <svg
                    v-else
                    class="h-4 w-4"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="2"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
                    />
                  </svg>
                  <template v-if="restarting">
                    <span>{{ t('version.restarting') }}</span>
                    <span v-if="restartCountdown > 0" class="tabular-nums"
                      >({{ restartCountdown }}s)</span
                    >
                  </template>
                  <span v-else>{{ t('version.restartNow') }}</span>
                </button>
              </div>

              <!-- Priority 3: Update available for source build -->
              <div v-else-if="hasUpdate && !isReleaseBuild" class="space-y-2">
                <div
                  class="flex items-center gap-2 rounded-lg border border-blue-200 bg-blue-50 p-2 border-blue-800/50 bg-blue-900/20"
                >
                  <svg
                    class="h-3.5 w-3.5 flex-shrink-0 text-blue-500 text-blue-400"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="2"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                  <p class="text-xs text-blue-600 text-blue-400">
                    {{ t('version.sourceModeHint') }}
                  </p>
                </div>
              </div>

              <!-- Priority 4: Update available for release build - show update button -->
              <div v-else-if="hasUpdate && isReleaseBuild" class="space-y-2">
                <!-- Update info card -->
                <div
                  class="flex items-center gap-3 rounded-lg border border-amber-200 bg-amber-50 p-3 border-amber-800/50 bg-amber-900/20"
                >
                <div
                  class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full bg-amber-100 bg-amber-900/50"
                >
                  <Icon
                    name="download"
                    size="sm"
                    :stroke-width="2"
                    class="text-neon-amber text-amber-400"
                  />
                </div>
                  <div class="min-w-0 flex-1">
                    <p class="text-sm font-medium text-amber-700 text-amber-300">
                      {{ t('version.updateAvailable') }}
                    </p>
                  </div>
                </div>

                <!-- Update button -->
                <button
                  @click="handleUpdate"
                  :disabled="updating"
                  class="flex w-full items-center justify-center gap-2 rounded-lg bg-neon-cyan/100 px-4 py-2 text-sm font-medium transition-colors hover:bg-neon-cyan disabled:cursor-not-allowed disabled:opacity-50"
                >
                  <svg v-if="updating" class="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle
                      class="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      stroke-width="4"
                    ></circle>
                    <path
                      class="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    ></path>
                  </svg>
                  <Icon v-else name="download" size="sm" :stroke-width="2" />
                  {{ updating ? t('version.updating') : t('version.updateNow') }}
                </button>
              </div>

              <!-- Up to date: expose the preserved rollback workflow. -->
              <div v-else class="space-y-2 border-t border-cyber-border pt-2">
                <button
                  type="button"
                  class="flex w-full items-center justify-between rounded-lg px-2 py-1.5 text-xs text-txt-secondary transition-colors hover:bg-cyber-elevated hover:text-txt-primary"
                  @click="toggleRollbackPanel"
                >
                  <span class="flex items-center gap-1.5">
                    <Icon name="clock" size="xs" :stroke-width="2" />
                    {{ t('version.rollback') }}
                  </span>
                  <Icon
                    name="chevronDown"
                    size="xs"
                    :stroke-width="2"
                    class="transition-transform"
                    :class="{ 'rotate-180': rollbackPanelOpen }"
                  />
                </button>

                <transition name="rollback">
                  <div v-if="rollbackPanelOpen" class="space-y-2">
                    <p v-if="!isReleaseBuild" class="rounded-lg bg-blue-50 p-2 text-xs text-blue-600">
                      {{ t('version.rollbackSourceHint') }}
                    </p>
                    <template v-else>
                      <p v-if="rollbackVersionsLoading" class="py-3 text-center text-xs text-txt-secondary">
                        {{ t('common.loading') }}
                      </p>
                      <p v-else-if="rollbackVersionsError" class="rounded-lg bg-red-50 p-2 text-xs text-red-600">
                        {{ rollbackVersionsError }}
                      </p>
                      <div v-else class="max-h-36 space-y-1 overflow-y-auto">
                        <button
                          v-for="item in rollbackVersions"
                          :key="item.version"
                          type="button"
                          class="flex w-full items-center justify-between rounded-lg border px-2 py-1.5 text-left text-xs"
                          :class="selectedRollbackVersion === item.version
                            ? 'border-primary-500 bg-primary-500/10 text-primary-600'
                            : 'border-cyber-border text-txt-secondary hover:bg-cyber-elevated'"
                          @click="selectRollbackVersion(item.version)"
                        >
                          <span>v{{ item.version }}</span>
                          <span class="text-txt-dim">{{ formatPublishedAt(item.published_at) }}</span>
                        </button>
                      </div>

                      <div v-if="selectedRollbackVersion" class="space-y-2 rounded-lg border border-cyber-border p-2">
                        <div class="flex gap-1">
                          <button
                            v-for="tab in manualTabs"
                            :key="tab.key"
                            type="button"
                            class="rounded px-2 py-1 text-[11px]"
                            :class="manualTab === tab.key ? 'bg-primary-500 text-white' : 'bg-cyber-elevated text-txt-secondary'"
                            @click="manualTab = tab.key"
                          >
                            {{ tab.label }}
                          </button>
                        </div>
                        <code class="block max-h-24 overflow-auto whitespace-pre-wrap break-all rounded bg-cyber-bg p-2 text-[10px] text-txt-secondary">{{ activeManualCommand }}</code>
                        <button
                          type="button"
                          class="btn btn-secondary w-full text-xs"
                          @click="copyToClipboard(activeManualCommand)"
                        >
                          {{ copied ? t('common.copied') : t('common.copy') }}
                        </button>
                        <p v-if="rollbackError" class="text-xs text-red-500">{{ rollbackError }}</p>
                        <button
                          type="button"
                          class="btn btn-primary w-full text-xs"
                          :disabled="rollingBack"
                          @click="handleRollback"
                        >
                          {{ rollingBack ? t('version.rollingBack') : t('version.rollback') }}
                        </button>
                      </div>
                    </template>
                  </div>
                </transition>
              </div>
            </template>
          </div>
        </div>
      </transition>
    </template>

    <!-- Non-admin: no version display -->
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import {
  performUpdate,
  restartService,
  getRollbackVersions,
  rollback as rollbackAPI,
  type RollbackVersionInfo
} from '@/api/admin/system'
import { useClipboard } from '@/composables/useClipboard'
import Icon from '@/components/icons/Icon.vue'

const GITHUB_REPO = 'Wei-Shaw/sub2api'
// Docker Hub image published by CI (tags carry no "v" prefix, e.g. weishaw/sub2api:0.1.146)
const DOCKER_IMAGE = 'weishaw/sub2api'

const { t } = useI18n()

defineProps<{
  version?: string
}>()

const authStore = useAuthStore()
const appStore = useAppStore()

const isAdmin = computed(() => authStore.isAdmin)

const dropdownOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)

// Use store's cached version state
const loading = computed(() => appStore.versionLoading)
const hasUpdate = computed(() => appStore.hasUpdate)
const buildType = computed(() => appStore.buildType)

// Update process states (local to this component)
const updating = ref(false)
const restarting = ref(false)
const needRestart = ref(false)
const updateError = ref('')
const updateSuccess = ref(false)
const restartCountdown = ref(0)
// Distinguishes the success + restart panel between update and rollback flows
const successKind = ref<'update' | 'rollback'>('update')

// Rollback states
const rollbackPanelOpen = ref(false)
const rollbackVersions = ref<RollbackVersionInfo[]>([])
const rollbackVersionsLoading = ref(false)
const rollbackVersionsError = ref('')
const selectedRollbackVersion = ref('')
const rollingBack = ref(false)
const rollbackError = ref('')

const { copied, copyToClipboard } = useClipboard()

// Manual rollback methods differ by deployment: script installs use install.sh,
// docker deployments pin the image tag instead
const manualTab = ref<'script' | 'docker'>('script')

const manualTabs = computed(() => [
  { key: 'script' as const, label: t('version.deployScript') },
  { key: 'docker' as const, label: t('version.deployDocker') }
])

const scriptRollbackCommand = computed(() => {
  if (!selectedRollbackVersion.value) return ''
  const tag = `v${selectedRollbackVersion.value}`
  return `curl -sSL https://raw.githubusercontent.com/${GITHUB_REPO}/${tag}/deploy/install.sh | sudo bash -s -- rollback ${tag}`
})

const dockerRollbackCommand = computed(() => {
  if (!selectedRollbackVersion.value) return ''
  return [
    `# ${t('version.dockerEditCompose')}`,
    `image: ${DOCKER_IMAGE}:${selectedRollbackVersion.value}`,
    '',
    `# ${t('version.dockerRecreate')}`,
    'docker compose up -d'
  ].join('\n')
})

const activeManualCommand = computed(() =>
  manualTab.value === 'docker' ? dockerRollbackCommand.value : scriptRollbackCommand.value
)

// Only show update check for release builds (binary/docker deployment)
const isReleaseBuild = computed(() => buildType.value === 'release')

function toggleDropdown() {
  dropdownOpen.value = !dropdownOpen.value
}

function closeDropdown() {
  dropdownOpen.value = false
}

async function refreshVersion(force = true) {
  if (!isAdmin.value) return

  // Reset update states when refreshing
  updateError.value = ''
  updateSuccess.value = false
  needRestart.value = false
  resetRollbackState()

  await appStore.fetchVersion(force)
}

async function handleUpdate() {
  if (updating.value) return

  updating.value = true
  updateError.value = ''
  updateSuccess.value = false

  try {
    const result = await performUpdate()
    successKind.value = 'update'
    updateSuccess.value = true
    needRestart.value = result.need_restart
    // Clear version cache to reflect update completed
    appStore.clearVersionCache()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } }; message?: string }
    updateError.value = err.response?.data?.message || err.message || t('version.updateFailed')
  } finally {
    updating.value = false
  }
}

function resetRollbackState() {
  rollbackPanelOpen.value = false
  rollbackVersions.value = []
  rollbackVersionsError.value = ''
  selectedRollbackVersion.value = ''
  rollbackError.value = ''
  manualTab.value = 'script'
}

async function toggleRollbackPanel() {
  if (!isAdmin.value) return
  rollbackPanelOpen.value = !rollbackPanelOpen.value
  // Source builds only show a hint, no version list to fetch
  if (
    rollbackPanelOpen.value &&
    isReleaseBuild.value &&
    rollbackVersions.value.length === 0 &&
    !rollbackVersionsLoading.value
  ) {
    await loadRollbackVersions()
  }
}

async function loadRollbackVersions() {
  if (!isAdmin.value) return
  rollbackVersionsLoading.value = true
  rollbackVersionsError.value = ''
  try {
    const data = await getRollbackVersions()
    rollbackVersions.value = data.versions || []
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } }; message?: string }
    rollbackVersionsError.value =
      err.response?.data?.message || err.message || t('version.loadVersionsFailed')
  } finally {
    rollbackVersionsLoading.value = false
  }
}

function selectRollbackVersion(version: string) {
  if (rollingBack.value) return
  rollbackError.value = ''
  selectedRollbackVersion.value = selectedRollbackVersion.value === version ? '' : version
}

function formatPublishedAt(publishedAt: string): string {
  if (!publishedAt) return ''
  const date = new Date(publishedAt)
  if (Number.isNaN(date.getTime())) return ''
  return date.toLocaleDateString()
}

async function handleRollback() {
  if (!isAdmin.value) return
  if (rollingBack.value || !selectedRollbackVersion.value) return

  rollingBack.value = true
  rollbackError.value = ''

  try {
    const result = await rollbackAPI(selectedRollbackVersion.value)
    successKind.value = 'rollback'
    updateSuccess.value = true
    needRestart.value = result.need_restart
    rollbackPanelOpen.value = false
    // Clear version cache so the next check reflects the rolled-back version
    appStore.clearVersionCache()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } }; message?: string }
    rollbackError.value = err.response?.data?.message || err.message || t('version.rollbackFailed')
  } finally {
    rollingBack.value = false
  }
}

async function handleRestart() {
  if (restarting.value) return

  restarting.value = true
  restartCountdown.value = 8

  try {
    await restartService()
    // Service will restart, page will reload automatically or show disconnected
  } catch (error) {
    // Expected - connection will be lost during restart
    console.log('Service restarting...')
  }

  // Start countdown
  const countdownInterval = setInterval(() => {
    restartCountdown.value--
    if (restartCountdown.value <= 0) {
      clearInterval(countdownInterval)
      // Try to check if service is back before reload
      checkServiceAndReload()
    }
  }, 1000)
}

async function checkServiceAndReload() {
  const maxRetries = 5
  const retryDelay = 1000

  for (let i = 0; i < maxRetries; i++) {
    try {
      const response = await fetch('/health', {
        method: 'GET',
        cache: 'no-cache'
      })
      if (response.ok) {
        // Service is back, reload page
        window.location.reload()
        return
      }
    } catch {
      // Service not ready yet
    }

    if (i < maxRetries - 1) {
      await new Promise((resolve) => setTimeout(resolve, retryDelay))
    }
  }

  // After retries, reload anyway
  window.location.reload()
}

function handleClickOutside(event: MouseEvent) {
  const target = event.target as Node
  const button = (event.target as Element).closest('button')
  if (dropdownRef.value && !dropdownRef.value.contains(target) && !button?.contains(target)) {
    closeDropdown()
  }
}

onMounted(() => {
  if (isAdmin.value) {
    // Use cached version if available, otherwise fetch
    appStore.fetchVersion(false)
  }
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

.rollback-enter-active,
.rollback-leave-active {
  transition: all 0.2s ease;
}

.rollback-enter-from,
.rollback-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
