<template>
  <AppLayout>
    <div class="mx-auto max-w-7xl space-y-5">
      <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">生图广场</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-dark-300">使用自己的 API Key 创作、管理历史并发布到公开广场。</p>
        </div>
        <div class="flex rounded-lg border border-gray-200 bg-white p-1 dark:border-dark-700 dark:bg-dark-800">
          <button v-for="tab in tabs" :key="tab.key" type="button" class="rounded-md px-3 py-2 text-sm font-medium transition" :class="activeTab === tab.key ? 'bg-primary-600 text-white' : 'text-gray-600 hover:bg-gray-100 dark:text-dark-200 dark:hover:bg-dark-700'" @click="activeTab = tab.key">
            {{ tab.label }}
          </button>
        </div>
      </div>

      <div v-if="loading" class="flex justify-center py-16">
        <div class="h-8 w-8 animate-spin rounded-full border-b-2 border-primary-600"></div>
      </div>

      <div v-else-if="settings && !settings.enabled" class="rounded-lg border border-amber-200 bg-amber-50 p-5 text-sm text-amber-800 dark:border-amber-800 dark:bg-amber-900/20 dark:text-amber-200">
        生图广场当前未启用。
      </div>

      <template v-else-if="settings">
        <section v-show="activeTab === 'create'" class="grid gap-5 lg:grid-cols-[360px,1fr]">
          <form class="space-y-4 rounded-lg border border-gray-200 bg-white p-5 dark:border-dark-700 dark:bg-dark-800" @submit.prevent="submitGenerate">
            <div>
              <label class="input-label">API Key</label>
              <select v-model.number="form.api_key_id" class="input">
                <option :value="0">请选择 KEY</option>
                <option v-for="key in keys" :key="key.id" :value="key.id">{{ key.name }} · {{ key.group_name || '默认分组' }}</option>
              </select>
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="input-label">模型</label>
                <select v-model="form.model" class="input">
                  <option v-for="model in settings.allowed_models" :key="model" :value="model">{{ model }}</option>
                </select>
              </div>
              <div>
                <label class="input-label">数量</label>
                <input v-model.number="form.n" type="number" min="1" :max="settings.max_n" class="input" />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="input-label">尺寸</label>
                <select v-model="form.size" class="input">
                  <option value="">默认</option>
                  <option v-for="size in settings.allowed_sizes" :key="size" :value="size">{{ size }}</option>
                </select>
              </div>
              <div>
                <label class="input-label">质量</label>
                <select v-model="form.quality" class="input">
                  <option value="">默认</option>
                  <option v-for="quality in settings.allowed_quality" :key="quality" :value="quality">{{ quality }}</option>
                </select>
              </div>
            </div>
            <div>
              <label class="input-label">输出格式</label>
              <select v-model="form.output_format" class="input">
                <option value="">默认</option>
                <option v-for="format in settings.allowed_output_formats" :key="format" :value="format">{{ format }}</option>
              </select>
            </div>
            <div>
              <label class="input-label">提示词</label>
              <textarea v-model="form.prompt" class="input min-h-[180px] resize-y" placeholder="描述你想生成的图片"></textarea>
            </div>
            <div>
              <label class="input-label">参考图</label>
              <input type="file" accept="image/*" multiple class="input" @change="onReferenceFilesChange" />
              <p v-if="referenceFiles.length" class="mt-1 text-xs text-gray-500 dark:text-dark-300">已选择 {{ referenceFiles.length }} 张，提交时会转为图片编辑请求。</p>
            </div>
            <div>
              <label class="input-label">蒙版</label>
              <input type="file" accept="image/*" class="input" @change="onMaskFileChange" />
              <p v-if="maskFile" class="mt-1 text-xs text-gray-500 dark:text-dark-300">{{ maskFile.name }}</p>
            </div>
            <p class="text-xs text-gray-500 dark:text-dark-300">生成前费用按模型和尺寸预估可能不准确，最终以实际用量记录为准。</p>
            <button type="submit" class="btn btn-primary w-full" :disabled="generating || !canSubmit">
              {{ generating ? '生成中...' : '开始生成' }}
            </button>
          </form>

          <div class="rounded-lg border border-gray-200 bg-white p-5 dark:border-dark-700 dark:bg-dark-800">
            <div class="mb-4 flex items-center justify-between">
              <h2 class="text-base font-semibold text-gray-900 dark:text-white">生成结果</h2>
              <button type="button" class="btn btn-secondary btn-sm" @click="loadHistory">刷新历史</button>
            </div>
            <div v-if="results.length === 0" class="flex min-h-[360px] items-center justify-center rounded-lg border border-dashed border-gray-300 text-sm text-gray-400 dark:border-dark-600">
              生成后的图片会显示在这里。
            </div>
            <div v-else class="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
              <ImageAssetCard v-for="item in results" :key="item.id" :item="item" @publish="publish" @delete="remove" />
            </div>
          </div>
        </section>

        <section v-show="activeTab === 'templates'" class="space-y-4">
          <div class="flex max-w-xl gap-2">
            <input v-model="templateQuery" class="input" placeholder="搜索模板" @keyup.enter="loadTemplates" />
            <button type="button" class="btn btn-secondary" @click="loadTemplates">搜索</button>
          </div>
          <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
            <div v-for="tpl in templates" :key="tpl.id" class="rounded-lg border border-gray-200 bg-white p-4 dark:border-dark-700 dark:bg-dark-800">
              <div class="flex items-start justify-between gap-3">
                <div>
                  <h3 class="font-semibold text-gray-900 dark:text-white">{{ tpl.title }}</h3>
                  <p class="mt-1 text-xs text-gray-500">{{ tpl.category || '未分类' }}</p>
                </div>
                <button type="button" class="btn btn-primary btn-sm" @click="applyTemplate(tpl.prompt)">套用</button>
              </div>
              <p class="mt-3 line-clamp-4 whitespace-pre-wrap text-sm text-gray-600 dark:text-dark-200">{{ tpl.prompt }}</p>
              <div class="mt-3 flex flex-wrap gap-1">
                <span v-for="tag in tpl.tags" :key="tag" class="rounded bg-gray-100 px-2 py-0.5 text-xs text-gray-600 dark:bg-dark-700 dark:text-dark-200">{{ tag }}</span>
              </div>
            </div>
          </div>
        </section>

        <section v-show="activeTab === 'history'" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          <ImageAssetCard v-for="item in history" :key="item.id" :item="item" @publish="publish" @delete="remove" />
          <div v-if="history.length === 0" class="col-span-full rounded-lg border border-dashed border-gray-300 p-10 text-center text-sm text-gray-400 dark:border-dark-600">暂无历史记录</div>
        </section>

        <section v-show="activeTab === 'public'" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          <ImageAssetCard v-for="item in publicItems" :key="item.id" :item="item" public-card @favorite="favorite" @reuse="reusePrompt" />
          <div v-if="publicItems.length === 0" class="col-span-full rounded-lg border border-dashed border-gray-300 p-10 text-center text-sm text-gray-400 dark:border-dark-600">暂无公开图片</div>
        </section>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, defineComponent, h, onMounted, reactive, ref, watch } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import { useAppStore } from '@/stores/app'
import {
  deleteItem,
  favoritePublicItem,
  generateImage,
  getEligibleKeys,
  getSettings,
  listHistory,
  listPublic,
  listTemplates,
  publishItem,
  type ImageGalleryAsset,
  type ImageGalleryEligibleKey,
  type ImageGallerySettings,
  type ImageGalleryTemplate,
} from '@/api/imageGallery'

const appStore = useAppStore()
const loading = ref(true)
const generating = ref(false)
const settings = ref<ImageGallerySettings | null>(null)
const keys = ref<ImageGalleryEligibleKey[]>([])
const results = ref<ImageGalleryAsset[]>([])
const history = ref<ImageGalleryAsset[]>([])
const publicItems = ref<ImageGalleryAsset[]>([])
const templates = ref<ImageGalleryTemplate[]>([])
const referenceFiles = ref<File[]>([])
const maskFile = ref<File | null>(null)
const templateQuery = ref('')
const activeTab = ref<'create' | 'templates' | 'history' | 'public'>('create')

const tabs = [
  { key: 'create', label: '创作' },
  { key: 'templates', label: '灵感模板' },
  { key: 'history', label: '我的历史' },
  { key: 'public', label: '公开广场' },
] as const

const form = reactive({
  api_key_id: 0,
  prompt: '',
  model: '',
  size: '',
  quality: '',
  n: 1,
  output_format: '',
})

const canSubmit = computed(() => form.api_key_id > 0 && form.prompt.trim().length > 0 && form.n > 0)

async function loadInitial() {
  loading.value = true
  try {
    settings.value = await getSettings()
    if (!settings.value.enabled) return
    form.model = settings.value.default_model
    const [eligible, tpl, hist, pub] = await Promise.all([
      getEligibleKeys(),
      listTemplates({ page_size: 24 }),
      listHistory(1, 24),
      settings.value.public_enabled ? listPublic(1, 24) : Promise.resolve({ items: [], total: 0, page: 1, page_size: 24, pages: 1 }),
    ])
    keys.value = eligible
    if (eligible[0]) form.api_key_id = eligible[0].id
    templates.value = tpl.items
    history.value = hist.items
    publicItems.value = pub.items
  } catch (error: any) {
    appStore.showError(error?.message || '加载生图广场失败')
  } finally {
    loading.value = false
  }
}

async function submitGenerate() {
  if (!canSubmit.value || !settings.value) return
  generating.value = true
  try {
    const referenceImages = await Promise.all(referenceFiles.value.map(readFileAsDataURL))
    const payload = {
      ...form,
      reference_images: referenceImages,
      mask_image: maskFile.value ? await readFileAsDataURL(maskFile.value) : undefined,
    }
    const result = await generateImage(payload)
    results.value = result.assets
    await loadHistory()
    if (result.assets.length === 0 && result.job?.error) {
      appStore.showWarning(result.job.error)
    } else {
      appStore.showSuccess('生成完成')
    }
  } catch (error: any) {
    appStore.showError(error?.message || '生成失败')
  } finally {
    generating.value = false
  }
}

function onReferenceFilesChange(event: Event) {
  const input = event.target as HTMLInputElement
  referenceFiles.value = Array.from(input.files || []).slice(0, 4)
}

function onMaskFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  maskFile.value = input.files?.[0] || null
}

function readFileAsDataURL(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(String(reader.result || ''))
    reader.onerror = () => reject(reader.error || new Error('读取图片失败'))
    reader.readAsDataURL(file)
  })
}

async function loadHistory() {
  const res = await listHistory(1, 24)
  history.value = res.items
}

async function loadTemplates() {
  const res = await listTemplates({ q: templateQuery.value, page_size: 24 })
  templates.value = res.items
}

async function loadPublicItems() {
  if (!settings.value?.public_enabled) return
  const res = await listPublic(1, 24)
  publicItems.value = res.items
}

function applyTemplate(prompt: string) {
  form.prompt = prompt
  activeTab.value = 'create'
}

function reusePrompt(prompt: string) {
  form.prompt = prompt
  activeTab.value = 'create'
}

async function publish(item: ImageGalleryAsset) {
  try {
    const updated = await publishItem(item.id)
    replaceLocalItem(updated)
    appStore.showSuccess(updated.review_status === 'pending' ? '已提交审核' : '已发布')
  } catch (error: any) {
    appStore.showError(error?.message || '发布失败')
  }
}

async function remove(item: ImageGalleryAsset) {
  try {
    await deleteItem(item.id)
    results.value = results.value.filter((x) => x.id !== item.id)
    history.value = history.value.filter((x) => x.id !== item.id)
    appStore.showSuccess('已删除')
  } catch (error: any) {
    appStore.showError(error?.message || '删除失败')
  }
}

async function favorite(item: ImageGalleryAsset) {
  try {
    await favoritePublicItem(item.id)
    item.favorited = !item.favorited
  } catch (error: any) {
    appStore.showError(error?.message || '操作失败')
  }
}

function replaceLocalItem(item: ImageGalleryAsset) {
  for (const list of [results.value, history.value]) {
    const idx = list.findIndex((x) => x.id === item.id)
    if (idx >= 0) list[idx] = item
  }
}

watch(activeTab, (tab) => {
  if (tab === 'history') loadHistory()
  if (tab === 'public') loadPublicItems()
})

onMounted(loadInitial)

const ImageAssetCard = defineComponent({
  props: {
    item: { type: Object as () => ImageGalleryAsset, required: true },
    publicCard: { type: Boolean, default: false },
  },
  emits: ['publish', 'delete', 'favorite', 'reuse'],
  setup(props, { emit }) {
    return () => h('article', { class: 'overflow-hidden rounded-lg border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-800' }, [
      h('img', { src: props.item.url, class: 'aspect-square w-full bg-gray-100 object-cover dark:bg-dark-700', loading: 'lazy' }),
      h('div', { class: 'space-y-3 p-3' }, [
        h('p', { class: 'line-clamp-3 min-h-[3.75rem] text-sm text-gray-700 dark:text-dark-200' }, props.item.prompt),
        h('div', { class: 'flex flex-wrap items-center gap-2 text-xs text-gray-500 dark:text-dark-300' }, [
          h('span', props.item.model),
          props.item.actual_cost != null ? h('span', `$${props.item.actual_cost.toFixed(6)}`) : null,
          props.item.review_status !== 'none' ? h('span', props.item.review_status) : null,
        ]),
        h('div', { class: 'flex flex-wrap gap-2' }, props.publicCard ? [
          h('button', { type: 'button', class: 'btn btn-secondary btn-sm', onClick: () => emit('reuse', props.item.prompt) }, '再创作'),
          h('button', { type: 'button', class: 'btn btn-secondary btn-sm', onClick: () => emit('favorite', props.item) }, props.item.favorited ? '已收藏' : '收藏'),
          h('a', { class: 'btn btn-secondary btn-sm', href: props.item.download_url }, '下载'),
        ] : [
          h('button', { type: 'button', class: 'btn btn-secondary btn-sm', onClick: () => emit('publish', props.item) }, props.item.review_status === 'pending' ? '审核中' : '发布'),
          h('a', { class: 'btn btn-secondary btn-sm', href: props.item.download_url }, '下载'),
          h('button', { type: 'button', class: 'btn btn-secondary btn-sm text-red-600 dark:text-red-400', onClick: () => emit('delete', props.item) }, '删除'),
        ]),
      ]),
    ])
  },
})
</script>
