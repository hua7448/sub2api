<template>
  <AppLayout>
    <div class="mx-auto max-w-4xl space-y-6">
      <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">生图广场设置</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-dark-300">配置内嵌 Image Playground 的可用模型、参数和 Agent 能力。</p>
        </div>
        <button type="button" class="btn btn-primary" :disabled="saving || !settings" @click="saveSettings">
          {{ saving ? '保存中...' : '保存设置' }}
        </button>
      </div>

      <div v-if="loading" class="flex justify-center py-16">
        <div class="h-8 w-8 animate-spin rounded-full border-b-2 border-primary-600"></div>
      </div>

      <template v-else-if="settings">
        <section class="rounded-lg border border-gray-200 bg-white p-5 dark:border-dark-700 dark:bg-dark-800">
          <h2 class="mb-4 text-base font-semibold text-gray-900 dark:text-white">功能开关</h2>
          <div class="space-y-4">
            <label class="flex items-center justify-between gap-3">
              <span class="text-sm text-gray-700 dark:text-dark-100">启用用户入口</span>
              <input v-model="settings.enabled" type="checkbox" class="h-5 w-5" />
            </label>
            <label class="flex items-center justify-between gap-3">
              <span class="text-sm text-gray-700 dark:text-dark-100">启用 Agent 模式</span>
              <input v-model="settings.agent_enabled" type="checkbox" class="h-5 w-5" />
            </label>
            <label class="flex items-center justify-between gap-3">
              <span class="text-sm text-gray-700 dark:text-dark-100">允许 Agent Web Search</span>
              <input v-model="settings.agent_web_search_enabled" :disabled="!settings.agent_enabled" type="checkbox" class="h-5 w-5 disabled:opacity-50" />
            </label>
          </div>
        </section>

        <section class="rounded-lg border border-gray-200 bg-white p-5 dark:border-dark-700 dark:bg-dark-800">
          <h2 class="mb-4 text-base font-semibold text-gray-900 dark:text-white">模型与参数限制</h2>
          <div class="grid gap-4 md:grid-cols-2">
            <label class="block">
              <span class="input-label">默认模型</span>
              <input v-model="settings.default_model" class="input" />
            </label>
            <label class="block">
              <span class="input-label">最大生成数量</span>
              <input v-model.number="settings.max_n" type="number" min="1" class="input" />
            </label>
            <label class="block md:col-span-2">
              <span class="input-label">允许模型</span>
              <input v-model="allowedModelsText" class="input" placeholder="gpt-image-2, gpt-image-1" />
            </label>
            <label class="block md:col-span-2">
              <span class="input-label">允许尺寸</span>
              <input v-model="allowedSizesText" class="input" placeholder="1024x1024, 1024x1536, 1536x1024, auto" />
            </label>
            <label class="block">
              <span class="input-label">允许质量</span>
              <input v-model="allowedQualityText" class="input" placeholder="auto, low, medium, high" />
            </label>
            <label class="block">
              <span class="input-label">允许输出格式</span>
              <input v-model="allowedFormatsText" class="input" placeholder="png, jpeg, webp" />
            </label>
            <label class="block">
              <span class="input-label">最大上传 MB</span>
              <input v-model.number="settings.max_upload_mb" type="number" min="1" class="input" />
            </label>
          </div>
        </section>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import { useAppStore } from '@/stores/app'
import { adminImageGalleryAPI } from '@/api/admin/imageGallery'
import type { ImageGallerySettings } from '@/api/imageGallery'

const appStore = useAppStore()
const loading = ref(true)
const saving = ref(false)
const settings = ref<ImageGallerySettings | null>(null)

const allowedModelsText = listSetting('allowed_models')
const allowedSizesText = listSetting('allowed_sizes')
const allowedQualityText = listSetting('allowed_quality')
const allowedFormatsText = listSetting('allowed_output_formats')

function listSetting(key: 'allowed_models' | 'allowed_sizes' | 'allowed_quality' | 'allowed_output_formats') {
  return computed({
    get: () => settings.value?.[key].join(', ') ?? '',
    set: (value: string) => {
      if (settings.value) settings.value[key] = splitList(value)
    },
  })
}

async function loadSettings() {
  loading.value = true
  try {
    settings.value = await adminImageGalleryAPI.getSettings()
  } catch (error: any) {
    appStore.showError(error?.message || '加载生图广场设置失败')
  } finally {
    loading.value = false
  }
}

async function saveSettings() {
  if (!settings.value) return
  saving.value = true
  try {
    settings.value = await adminImageGalleryAPI.updateSettings(settings.value)
    appStore.showSuccess('已保存')
  } catch (error: any) {
    appStore.showError(error?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

function splitList(value: string): string[] {
  return value.split(',').map((item) => item.trim()).filter(Boolean)
}

onMounted(loadSettings)
</script>
