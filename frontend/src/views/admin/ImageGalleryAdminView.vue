<template>
  <AppLayout>
    <div class="mx-auto max-w-7xl space-y-6">
      <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">生图广场管理</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-dark-300">配置入口、审核公开图片、管理模板和本地存储。</p>
        </div>
        <button type="button" class="btn btn-primary" :disabled="saving || !settings" @click="saveSettings">
          {{ saving ? '保存中...' : '保存设置' }}
        </button>
      </div>

      <div v-if="loading" class="flex justify-center py-16">
        <div class="h-8 w-8 animate-spin rounded-full border-b-2 border-primary-600"></div>
      </div>

      <template v-else-if="settings">
        <section class="grid gap-5 lg:grid-cols-3">
          <div class="rounded-lg border border-gray-200 bg-white p-5 dark:border-dark-700 dark:bg-dark-800">
            <h2 class="mb-4 text-base font-semibold text-gray-900 dark:text-white">功能开关</h2>
            <div class="space-y-4">
              <label class="flex items-center justify-between gap-3">
                <span class="text-sm text-gray-700 dark:text-dark-100">启用用户菜单</span>
                <input v-model="settings.enabled" type="checkbox" class="h-5 w-5" />
              </label>
              <label class="flex items-center justify-between gap-3">
                <span class="text-sm text-gray-700 dark:text-dark-100">启用公开广场</span>
                <input v-model="settings.public_enabled" type="checkbox" class="h-5 w-5" />
              </label>
              <label class="flex items-center justify-between gap-3">
                <span class="text-sm text-gray-700 dark:text-dark-100">公开前需要审核</span>
                <input v-model="settings.publish_requires_review" type="checkbox" class="h-5 w-5" />
              </label>
              <label class="flex items-center justify-between gap-3">
                <span class="text-sm text-gray-700 dark:text-dark-100">启用模板</span>
                <input v-model="settings.templates_enabled" type="checkbox" class="h-5 w-5" />
              </label>
              <label class="flex items-center justify-between gap-3">
                <span class="text-sm text-gray-700 dark:text-dark-100">允许模板导入</span>
                <input v-model="settings.template_import_enabled" type="checkbox" class="h-5 w-5" />
              </label>
            </div>
          </div>

          <div class="rounded-lg border border-gray-200 bg-white p-5 dark:border-dark-700 dark:bg-dark-800">
            <h2 class="mb-4 text-base font-semibold text-gray-900 dark:text-white">生成限制</h2>
            <div class="space-y-3">
              <label class="block">
                <span class="input-label">默认模型</span>
                <input v-model="settings.default_model" class="input" />
              </label>
              <label class="block">
                <span class="input-label">允许模型</span>
                <input v-model="allowedModelsText" class="input" placeholder="gpt-image-2, gpt-image-1" />
              </label>
              <label class="block">
                <span class="input-label">允许尺寸</span>
                <input v-model="allowedSizesText" class="input" />
              </label>
              <div class="grid grid-cols-2 gap-3">
                <label class="block">
                  <span class="input-label">最大数量</span>
                  <input v-model.number="settings.max_n" type="number" min="1" class="input" />
                </label>
                <label class="block">
                  <span class="input-label">上传上限 MB</span>
                  <input v-model.number="settings.max_upload_mb" type="number" min="1" class="input" />
                </label>
              </div>
            </div>
          </div>

          <div class="rounded-lg border border-gray-200 bg-white p-5 dark:border-dark-700 dark:bg-dark-800">
            <h2 class="mb-4 text-base font-semibold text-gray-900 dark:text-white">本地存储</h2>
            <div class="space-y-3">
              <label class="block">
                <span class="input-label">路径</span>
                <input v-model="settings.storage_path" class="input" />
              </label>
              <div class="grid grid-cols-2 gap-3">
                <label class="block">
                  <span class="input-label">用户配额 MB</span>
                  <input v-model.number="settings.user_quota_mb" type="number" min="0" class="input" />
                </label>
                <label class="block">
                  <span class="input-label">保留天数</span>
                  <input v-model.number="settings.retention_days" type="number" min="0" class="input" />
                </label>
              </div>
              <div class="rounded-lg bg-gray-50 p-3 text-sm text-gray-600 dark:bg-dark-700 dark:text-dark-200">
                <div>文件数：{{ storageStats?.file_count ?? 0 }}</div>
                <div>占用：{{ formatBytes(storageStats?.total_bytes ?? 0) }}</div>
              </div>
              <button type="button" class="btn btn-secondary btn-sm" @click="cleanup">清理已删除/过期文件</button>
            </div>
          </div>
        </section>

        <section class="rounded-lg border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-800">
          <div class="flex items-center justify-between border-b border-gray-100 p-5 dark:border-dark-700">
            <h2 class="text-base font-semibold text-gray-900 dark:text-white">公开审核</h2>
            <button type="button" class="btn btn-secondary btn-sm" @click="loadItems">刷新</button>
          </div>
          <div class="grid gap-4 p-5 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            <article v-for="item in items" :key="item.id" class="overflow-hidden rounded-lg border border-gray-200 dark:border-dark-700">
              <img :src="item.url" class="aspect-square w-full bg-gray-100 object-cover dark:bg-dark-700" loading="lazy" />
              <div class="space-y-3 p-3">
                <p class="line-clamp-3 text-sm text-gray-700 dark:text-dark-200">{{ item.prompt }}</p>
                <div class="text-xs text-gray-500">{{ item.review_status }} · {{ item.visibility }}</div>
                <div class="flex flex-wrap gap-2">
                  <button type="button" class="btn btn-primary btn-sm" @click="approve(item.id)">通过</button>
                  <button type="button" class="btn btn-secondary btn-sm" @click="reject(item.id)">拒绝</button>
                  <button type="button" class="btn btn-secondary btn-sm" @click="hide(item.id)">隐藏</button>
                  <button type="button" class="btn btn-secondary btn-sm text-red-600 dark:text-red-400" @click="remove(item.id)">删除</button>
                </div>
              </div>
            </article>
            <div v-if="items.length === 0" class="col-span-full p-10 text-center text-sm text-gray-400">暂无待审核图片</div>
          </div>
        </section>

        <section class="rounded-lg border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-800">
          <div class="flex items-center justify-between border-b border-gray-100 p-5 dark:border-dark-700">
            <h2 class="text-base font-semibold text-gray-900 dark:text-white">灵感模板</h2>
            <button type="button" class="btn btn-secondary btn-sm" @click="loadTemplates">刷新</button>
          </div>
          <div class="divide-y divide-gray-100 dark:divide-dark-700">
            <div v-for="tpl in templates" :key="tpl.id" class="grid gap-3 p-4 md:grid-cols-[1fr,140px,120px] md:items-center">
              <div>
                <div class="font-medium text-gray-900 dark:text-white">{{ tpl.title }}</div>
                <p class="mt-1 line-clamp-2 text-sm text-gray-500">{{ tpl.prompt }}</p>
              </div>
              <div class="text-sm text-gray-500">{{ tpl.category || '未分类' }}</div>
              <label class="flex items-center gap-2 text-sm">
                <input type="checkbox" :checked="tpl.enabled" @change="toggleTemplate(tpl)" />
                启用
              </label>
            </div>
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
import {
  adminImageGalleryAPI,
  type ImageGalleryStorageStats,
} from '@/api/admin/imageGallery'
import type { ImageGalleryAsset, ImageGallerySettings, ImageGalleryTemplate } from '@/api/imageGallery'

const appStore = useAppStore()
const loading = ref(true)
const saving = ref(false)
const settings = ref<ImageGallerySettings | null>(null)
const storageStats = ref<ImageGalleryStorageStats | null>(null)
const items = ref<ImageGalleryAsset[]>([])
const templates = ref<ImageGalleryTemplate[]>([])

const allowedModelsText = computed({
  get: () => settings.value?.allowed_models.join(', ') ?? '',
  set: (value: string) => {
    if (settings.value) settings.value.allowed_models = splitList(value)
  },
})

const allowedSizesText = computed({
  get: () => settings.value?.allowed_sizes.join(', ') ?? '',
  set: (value: string) => {
    if (settings.value) settings.value.allowed_sizes = splitList(value)
  },
})

async function loadAll() {
  loading.value = true
  try {
    await Promise.all([loadSettings(), loadItems(), loadTemplates(), loadStorageStats()])
  } catch (error: any) {
    appStore.showError(error?.message || '加载生图广场管理失败')
  } finally {
    loading.value = false
  }
}

async function loadSettings() {
  settings.value = await adminImageGalleryAPI.getSettings()
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

async function loadItems() {
  const res = await adminImageGalleryAPI.listItems({ review_status: 'pending', page_size: 24 })
  items.value = res.items
}

async function loadTemplates() {
  const res = await adminImageGalleryAPI.listTemplates({ page_size: 100 })
  templates.value = res.items
}

async function loadStorageStats() {
  storageStats.value = await adminImageGalleryAPI.getStorageStats()
}

async function approve(id: number) {
  await adminImageGalleryAPI.approveItem(id)
  await loadItems()
}

async function reject(id: number) {
  await adminImageGalleryAPI.rejectItem(id)
  await loadItems()
}

async function hide(id: number) {
  await adminImageGalleryAPI.hideItem(id)
  await loadItems()
}

async function remove(id: number) {
  await adminImageGalleryAPI.deleteItem(id)
  await loadItems()
}

async function cleanup() {
  try {
    const result = await adminImageGalleryAPI.cleanupStorage()
    appStore.showSuccess(`已清理 ${result.deleted_files} 个文件`)
    await loadStorageStats()
  } catch (error: any) {
    appStore.showError(error?.message || '清理失败')
  }
}

async function toggleTemplate(tpl: ImageGalleryTemplate) {
  const updated = await adminImageGalleryAPI.updateTemplate(tpl.id, { enabled: !tpl.enabled })
  const idx = templates.value.findIndex((item) => item.id === tpl.id)
  if (idx >= 0) templates.value[idx] = updated
}

function splitList(value: string): string[] {
  return value.split(',').map((item) => item.trim()).filter(Boolean)
}

function formatBytes(value: number): string {
  if (value < 1024) return `${value} B`
  if (value < 1024 * 1024) return `${(value / 1024).toFixed(1)} KB`
  if (value < 1024 * 1024 * 1024) return `${(value / 1024 / 1024).toFixed(1)} MB`
  return `${(value / 1024 / 1024 / 1024).toFixed(1)} GB`
}

onMounted(loadAll)
</script>
