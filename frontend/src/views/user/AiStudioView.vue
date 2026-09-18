<template>
  <AppLayout>
    <div
      data-testid="ai-studio-shell"
      class="-m-2 overflow-hidden rounded-[28px] border border-cyber-border bg-cyber-surface/95 shadow-glow backdrop-blur md:-m-4 lg:-m-6 xl:-m-8"
    >
      <section class="border-b border-cyber-border bg-gradient-to-br from-cyber-elevated/80 via-cyber-surface/95 to-cyber-surface px-4 py-2.5 sm:px-5">
        <div class="flex flex-col gap-2 xl:flex-row xl:items-center xl:justify-between">
          <div class="min-w-0 space-y-0.5">
            <h1 class="text-lg font-semibold tracking-tight text-txt-primary sm:text-xl">{{ t('imageStudio.title') }}</h1>
            <p class="max-w-3xl text-sm font-medium leading-5 text-txt-secondary sm:text-base">
              {{ t('imageStudio.description') }}
            </p>
          </div>

          <div
            data-testid="ai-studio-toolbar"
            class="grid shrink-0 gap-2 rounded-2xl border border-cyber-border bg-cyber-elevated/70 p-2.5 shadow-sm xl:w-[460px] xl:grid-cols-[minmax(0,1fr)_auto] xl:items-center"
          >
            <div v-if="loadingKeys" class="flex items-center gap-3 rounded-xl border border-cyber-border bg-cyber-surface/70 p-3 text-sm text-txt-secondary xl:col-span-2">
              <LoadingSpinner />
              <span>{{ t('imageStudio.loadingKeys') }}</span>
            </div>

            <div v-else-if="availableKeys.length === 0" class="flex flex-col gap-3 rounded-xl border border-dashed border-cyber-border bg-cyber-surface/70 p-3 sm:flex-row sm:items-center sm:justify-between xl:col-span-2">
              <div>
                <h2 class="text-sm font-semibold text-txt-primary">{{ t('imageStudio.noKeysTitle') }}</h2>
                <p class="mt-1 text-xs leading-5 text-txt-secondary">
                  {{ t('imageStudio.noKeysDescription') }}
                </p>
              </div>
              <button class="btn btn-primary justify-center" @click="goToKeys">
                <Icon name="plus" size="md" class="mr-2" />
                {{ t('imageStudio.manageKeys') }}
              </button>
            </div>

            <template v-else>
              <Select
                :model-value="selectedKeyId"
                :options="keyOptions"
                class="w-full"
                @update:model-value="handleKeyChange"
              />
              <div class="flex gap-2">
                <button class="btn btn-secondary btn-sm" :disabled="loadingKeys" @click="loadKeys">
                  <Icon name="refresh" size="sm" :class="loadingKeys ? 'animate-spin' : ''" />
                </button>
                <button class="btn btn-primary btn-sm" @click="goToKeys">
                  <Icon name="key" size="sm" />
                </button>
              </div>
              <p class="text-xs leading-5 text-txt-dim xl:col-span-2">
                {{ selectedKeyHint }}
              </p>
            </template>
          </div>
        </div>
      </section>

      <section data-testid="ai-studio-frame" class="bg-white">
        <div v-if="loadingKeys" class="flex min-h-[calc(100vh-132px)] items-center justify-center bg-cyber-surface/80 p-6 text-center">
          <div class="flex items-center gap-3 rounded-2xl border border-cyber-border bg-cyber-surface px-4 py-3 text-sm text-txt-secondary shadow-glow">
            <LoadingSpinner />
            <span>{{ t('imageStudio.loadingKeys') }}</span>
          </div>
        </div>

        <div v-else-if="availableKeys.length === 0" class="flex min-h-[calc(100vh-132px)] items-center justify-center bg-cyber-surface/80 p-6 text-center">
          <div class="max-w-md space-y-3">
            <h2 class="text-xl font-semibold text-txt-primary">{{ t('imageStudio.noKeysTitle') }}</h2>
            <p class="text-sm leading-6 text-txt-secondary">{{ t('imageStudio.noKeysDescription') }}</p>
          </div>
        </div>

        <div v-else class="relative min-h-[calc(100vh-132px)]">
          <div
            v-if="!iframeLoaded"
            class="absolute inset-0 z-10 flex items-center justify-center bg-cyber-surface/90 backdrop-blur-sm"
          >
            <div class="flex items-center gap-3 rounded-2xl border border-cyber-border bg-cyber-surface px-4 py-3 text-sm text-txt-secondary shadow-glow">
              <LoadingSpinner />
              <span>{{ t('imageStudio.frameLoading') }}</span>
            </div>
          </div>

          <iframe
            ref="studioFrame"
            :src="studioFrameSrc"
            class="h-[calc(100vh-132px)] min-h-[820px] w-full border-0 bg-[#f5f5f3]"
            :title="t('imageStudio.title')"
            @load="handleFrameLoad"
          />
        </div>
      </section>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import Select from '@/components/common/Select.vue'
import Icon from '@/components/icons/Icon.vue'
import keysAPI from '@/api/keys'
import { useAppStore, useAuthStore } from '@/stores'
import type { ApiKey } from '@/types'
import { extractApiErrorMessage } from '@/utils/apiError'

type StudioBridgePayload = {
  baseUrl: string
  apiKey: string
  apiKeyName?: string
  userScope?: string
  siteName?: string
}

const STUDIO_USAGE_UPDATED_MESSAGE = 'sub2api-image-studio-usage-updated'
const studioFrameSrc = '/image-studio/?v=scope5'

const { t } = useI18n()
const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()

const loadingKeys = ref(false)
const availableKeys = ref<ApiKey[]>([])
const selectedKeyId = ref<number | null>(null)
const iframeLoaded = ref(false)
const studioFrame = ref<HTMLIFrameElement | null>(null)

const selectedKey = computed(() => availableKeys.value.find((item) => item.id === selectedKeyId.value) ?? null)

const keyOptions = computed(() =>
  availableKeys.value.map((item) => ({
    value: item.id,
    label: item.group?.name ? `${item.name} · ${item.group.name}` : item.name,
  })),
)

const effectiveBaseUrl = computed(() => {
  if (typeof window === 'undefined') return ''
  return window.location.origin
})

const studioFrameOrigin = computed(() => {
  if (typeof window === 'undefined') return ''
  try {
    const rawSrc = studioFrame.value?.getAttribute('src') || studioFrameSrc
    return new URL(rawSrc, window.location.href).origin
  } catch {
    return window.location.origin
  }
})

const selectedKeyHint = computed(() => {
  if (!selectedKey.value) return t('imageStudio.loadingKeys')
  return t('imageStudio.activeKeyHint', { name: selectedKey.value.name })
})

function buildUserScope(): string {
  const userId = authStore.user?.id
  if (typeof userId === 'number' && userId > 0) {
    return `user:${userId}`
  }

  const keyOwnerId = selectedKey.value?.user_id
  if (typeof keyOwnerId === 'number' && keyOwnerId > 0) {
    return `user:${keyOwnerId}`
  }

  const keyId = selectedKey.value?.id
  if (typeof keyId === 'number' && keyId > 0) {
    return `key:${keyId}`
  }

  return 'anonymous'
}

function buildBridgePayload(): StudioBridgePayload | null {
  if (!selectedKey.value || !effectiveBaseUrl.value) {
    return null
  }
  return {
    baseUrl: effectiveBaseUrl.value,
    apiKey: selectedKey.value.key,
    apiKeyName: selectedKey.value.name,
    userScope: buildUserScope(),
    siteName: appStore.siteName,
  }
}

function postBridgeConfig() {
  const payload = buildBridgePayload()
  const frameWindow = studioFrame.value?.contentWindow
  if (!payload || !frameWindow || typeof window === 'undefined') {
    return
  }
  frameWindow.postMessage(
    {
      type: 'sub2api-image-studio-config',
      payload,
    },
    studioFrameOrigin.value || window.location.origin,
  )
}

async function loadKeys() {
  loadingKeys.value = true
  try {
    const response = await keysAPI.list(1, 100, { status: 'active' })
    availableKeys.value = (response.items || []).filter((item) => item.status === 'active')
    if (!availableKeys.value.some((item) => item.id === selectedKeyId.value)) {
      selectedKeyId.value = availableKeys.value[0]?.id ?? null
    }
  } catch (error: unknown) {
    availableKeys.value = []
    selectedKeyId.value = null
    appStore.showError(extractApiErrorMessage(error, t('common.error')))
  } finally {
    loadingKeys.value = false
  }
}

function handleFrameLoad() {
  iframeLoaded.value = true
  postBridgeConfig()
}

function handleKeyChange(value: string | number | boolean | null) {
  selectedKeyId.value = typeof value === 'number' ? value : Number(value)
  postBridgeConfig()
}

function goToKeys() {
  void router.push('/keys')
}

function handleStudioMessage(event: MessageEvent) {
  if (typeof window === 'undefined' || event.origin !== (studioFrameOrigin.value || window.location.origin)) {
    return
  }

  if (event.source !== studioFrame.value?.contentWindow) {
    return
  }

  const data = event.data as { type?: string } | null
  if (data?.type !== STUDIO_USAGE_UPDATED_MESSAGE) {
    return
  }

  void authStore.refreshUser()
}

watch(
  () => [selectedKey.value?.id, selectedKey.value?.user_id, effectiveBaseUrl.value, studioFrameOrigin.value, appStore.siteName, authStore.user?.id],
  async () => {
    if (!iframeLoaded.value) return
    await nextTick()
    postBridgeConfig()
  },
)

onMounted(() => {
  window.addEventListener('message', handleStudioMessage)
  void loadKeys()
})

onBeforeUnmount(() => {
  window.removeEventListener('message', handleStudioMessage)
})
</script>
