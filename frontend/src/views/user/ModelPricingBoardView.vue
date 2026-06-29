<template>
  <AppLayout>
    <section class="space-y-6">
      <div class="card overflow-hidden border-2 border-gray-300 bg-[#fffaf2] shadow-none dark:border-dark-700 dark:bg-dark-900">
        <div class="border-b border-gray-200 px-5 py-4 dark:border-dark-700">
          <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
            <div class="flex flex-wrap items-center gap-2 rounded-xl bg-gray-100 p-1 dark:bg-dark-800">
              <button
                v-for="tab in categoryTabs"
                :key="tab.value"
                type="button"
                class="inline-flex items-center rounded-lg px-4 py-2 text-sm font-semibold transition"
                :class="activeCategory === tab.value ? tab.activeClass : tab.inactiveClass"
                :data-testid="`model-pricing-tab-${tab.value}`"
                @click="activeCategory = tab.value"
              >
                {{ tab.label }}
              </button>
            </div>

            <div class="flex w-full flex-col gap-3 sm:flex-row lg:w-auto">
              <div class="relative min-w-0 flex-1 sm:w-80">
                <Icon
                  name="search"
                  size="md"
                  class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 dark:text-gray-500"
                />
                <input
                  v-model="searchQuery"
                  type="text"
                  :placeholder="t('modelPricing.searchPlaceholder')"
                  class="input pl-10"
                />
              </div>

              <button
                type="button"
                class="btn btn-secondary"
                :disabled="loading"
                :title="t('common.refresh', 'Refresh')"
                @click="loadBoard"
              >
                <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
              </button>
            </div>
          </div>
        </div>

        <div class="p-5 sm:p-6">
          <div
            v-if="loading"
            class="rounded-2xl border border-dashed border-gray-300 bg-white/70 px-6 py-16 text-center text-sm text-gray-500 dark:border-dark-700 dark:bg-dark-800/60 dark:text-gray-400"
          >
            {{ t('common.loading') }}
          </div>

          <div
            v-else-if="filteredItems.length === 0"
            class="rounded-2xl border border-dashed border-gray-300 bg-white/70 px-6 py-16 text-center dark:border-dark-700 dark:bg-dark-800/60"
          >
            <p class="text-base font-semibold text-gray-900 dark:text-white">
              {{ emptyStateTitle }}
            </p>
            <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">
              {{ emptyStateDescription }}
            </p>
          </div>

          <div v-else class="grid grid-cols-1 gap-5 xl:grid-cols-2">
            <article
              v-for="item in filteredItems"
              :key="`${item.platform}:${item.model_id}:${item.channel_id}:${item.group_id}`"
              class="rounded-2xl border border-gray-200 bg-white/90 p-5 shadow-sm transition hover:border-primary-300 dark:border-dark-700 dark:bg-dark-800/90 dark:hover:border-primary-700"
              :data-testid="`model-pricing-card-${item.model_id}`"
            >
              <div class="flex items-start justify-between gap-4">
                <div class="min-w-0">
                  <div class="flex items-center gap-3">
                    <div class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-2xl border border-gray-200 bg-[#fffaf2] dark:border-dark-700 dark:bg-dark-900">
                      <ModelIcon :model="item.model_id" size="24px" />
                    </div>
                    <div class="min-w-0">
                      <h2 class="truncate font-mono text-sm font-semibold text-gray-900 dark:text-white">
                        {{ item.model_id }}
                      </h2>
                      <div class="mt-1 flex flex-wrap items-center gap-2">
                        <span :class="platformBadgeClass(item)" :data-testid="`model-pricing-provider-${item.model_id}`">
                          {{ modelVendorLabel(item) }}
                        </span>
                        <span
                          v-if="item.user_rate_overridden"
                          class="inline-flex rounded-full border border-emerald-200 bg-emerald-50 px-2.5 py-1 text-xs font-medium text-emerald-700 dark:border-emerald-900/60 dark:bg-emerald-950/30 dark:text-emerald-300"
                        >
                          {{ t('modelPricing.userRate') }}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>

                <div
                  class="inline-flex flex-shrink-0 rounded-full border px-3 py-1 text-xs font-semibold"
                  :class="categoryBadgeClass(item)"
                >
                  {{ categoryLabel(item) }}
                </div>
              </div>

              <div class="mt-5 space-y-3">
                <div
                  v-for="row in pricingRows(item)"
                  :key="row.key"
                  class="rounded-xl border border-gray-200 bg-[#fffaf2]/80 p-4 dark:border-dark-700 dark:bg-dark-900/80"
                >
                  <div class="flex flex-wrap items-start justify-between gap-3">
                    <div>
                      <div class="text-xs font-medium uppercase tracking-[0.14em] text-gray-500 dark:text-gray-400">
                        {{ row.label }}
                      </div>
                      <div class="mt-2 text-xs font-semibold text-primary-700 dark:text-primary-300">
                        {{ t('modelPricing.card.sitePrice') }}
                      </div>
                      <div class="mt-2 text-2xl font-black text-primary-700 dark:text-primary-300">
                        {{ formatTokenPrice(row.sitePrice, item) }}
                      </div>
                      <p class="mt-2 max-w-md text-xs leading-relaxed text-gray-500 dark:text-gray-400">
                        {{ priceFormulaText(item) }}
                      </p>
                      <div
                        v-if="row.officialPrice != null"
                        class="mt-2 flex flex-wrap items-center gap-2 text-sm text-gray-400 dark:text-gray-500"
                      >
                        <span>{{ t('modelPricing.card.officialPrice') }}</span>
                        <span class="line-through">
                          {{ formatTokenPrice(row.officialPrice, item) }}
                        </span>
                      </div>
                    </div>

                    <div
                      v-if="row.savings != null && row.officialPrice != null"
                      class="inline-flex rounded-full border border-amber-200 bg-amber-50 px-3 py-1 text-xs font-semibold text-amber-800 dark:border-amber-900/60 dark:bg-amber-950/30 dark:text-amber-200"
                    >
                      {{ t('modelPricing.card.savings', { percent: formatSavings(row.savings) }) }}
                    </div>
                  </div>
                </div>
              </div>

              <div class="mt-5 grid grid-cols-1 gap-3 rounded-2xl border border-dashed border-gray-300 bg-gray-50/80 p-4 text-sm dark:border-dark-700 dark:bg-dark-900/60 sm:grid-cols-2">
                <div>
                  <div class="text-xs font-medium uppercase tracking-[0.14em] text-gray-500 dark:text-gray-400">
                    {{ t('modelPricing.card.group') }}
                  </div>
                  <div class="mt-1 font-medium text-gray-900 dark:text-white">
                    {{ item.group_name }}
                  </div>
                </div>
                <div>
                  <div class="text-xs font-medium uppercase tracking-[0.14em] text-gray-500 dark:text-gray-400">
                    {{ t('modelPricing.card.channel') }}
                  </div>
                  <div class="mt-1 break-words font-medium text-gray-900 dark:text-white">
                    {{ item.channel_name }}
                  </div>
                </div>
                <div>
                  <div class="text-xs font-medium uppercase tracking-[0.14em] text-gray-500 dark:text-gray-400">
                    {{ t('modelPricing.card.rateLabel') }}
                  </div>
                  <div class="mt-1 text-gray-700 dark:text-gray-300">
                    {{ t('modelPricing.card.rate', { rate: formatRate(item.rate_multiplier) }) }}
                  </div>
                </div>
                <div>
                  <div class="text-xs font-medium uppercase tracking-[0.14em] text-gray-500 dark:text-gray-400">
                    {{ t('modelPricing.card.source') }}
                  </div>
                  <div class="mt-1 text-gray-700 dark:text-gray-300">
                    {{ platformLabel(item.platform) }}
                  </div>
                </div>
              </div>
            </article>
          </div>
        </div>
      </div>
    </section>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import ModelIcon from '@/components/common/ModelIcon.vue'
import modelPricingAPI, { type ModelPricingBoardItem } from '@/api/modelPricing'
import { useAppStore } from '@/stores/app'
import { extractApiErrorMessage } from '@/utils/apiError'

type PricingCategory = 'claude' | 'codex' | 'domestic'

interface PricingRow {
  key: 'input' | 'output' | 'cacheRead'
  label: string
  sitePrice: number | null
  officialPrice: number | null
  savings: number | null
}

const { t } = useI18n()
const appStore = useAppStore()

const items = ref<ModelPricingBoardItem[]>([])
const loading = ref(false)
const searchQuery = ref('')
const activeCategory = ref<PricingCategory>('codex')

const categoryTabs = computed(() => [
  {
    value: 'codex' as const,
    label: t('modelPricing.tabs.codex'),
    activeClass: 'bg-primary-600 text-[#fffaf2] shadow-sm dark:bg-primary-400 dark:text-primary-950',
    inactiveClass: 'text-primary-800/70 hover:text-primary-950 dark:text-primary-200/70 dark:hover:text-primary-100',
  },
  {
    value: 'domestic' as const,
    label: t('modelPricing.tabs.domestic'),
    activeClass: 'bg-emerald-600 text-white shadow-sm dark:bg-emerald-400 dark:text-emerald-950',
    inactiveClass: 'text-emerald-800/70 hover:text-emerald-950 dark:text-emerald-200/70 dark:hover:text-emerald-100',
  },
  {
    value: 'claude' as const,
    label: t('modelPricing.tabs.claude'),
    activeClass: 'bg-amber-500 text-amber-950 shadow-sm',
    inactiveClass: 'text-amber-900/70 hover:text-amber-950 dark:text-amber-200/70 dark:hover:text-amber-100',
  },
])

const categorizedItems = computed(() =>
  items.value.filter((item) => classifyItem(item) === activeCategory.value),
)

const filteredItems = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return categorizedItems.value
  return categorizedItems.value.filter((item) =>
    item.model_id.toLowerCase().includes(q) ||
    item.platform.toLowerCase().includes(q) ||
    item.group_name.toLowerCase().includes(q) ||
    item.channel_name.toLowerCase().includes(q),
  )
})

const emptyStateTitle = computed(() => t(`modelPricing.emptyState.${activeCategory.value}`))
const emptyStateDescription = computed(() => t(`modelPricing.emptyStateDescription.${activeCategory.value}`))

function classifyItem(item: Pick<ModelPricingBoardItem, 'platform' | 'model_id'>): PricingCategory | null {
  if (isDomesticModel(item)) return 'domestic'
  if (item.platform === 'anthropic') return 'claude'
  if (item.platform === 'openai') return 'codex'
  return null
}

function isDomesticModel(item: Pick<ModelPricingBoardItem, 'platform' | 'model_id'>): boolean {
  const platform = item.platform.toLowerCase()
  const model = item.model_id.toLowerCase()
  return ['deepseek', 'zhipu', 'glm', 'kimi', 'moonshot', 'minimax', 'doubao'].some((needle) =>
    platform.includes(needle) || model.includes(needle),
  )
}

function pricingRows(item: ModelPricingBoardItem): PricingRow[] {
  const rows: PricingRow[] = [
    {
      key: 'input',
      label: t('modelPricing.card.input'),
      sitePrice: item.site_input_price,
      officialPrice: item.official_input_price,
      savings: item.input_savings,
    },
    {
      key: 'output',
      label: t('modelPricing.card.output'),
      sitePrice: item.site_output_price,
      officialPrice: item.official_output_price,
      savings: item.output_savings,
    },
    {
      key: 'cacheRead',
      label: t('modelPricing.card.cacheRead'),
      sitePrice: item.site_cache_read_price,
      officialPrice: item.official_cache_read_price,
      savings: item.cache_read_savings,
    },
  ]

  return rows.filter((row) => row.sitePrice != null)
}

function formatTokenPrice(value: number | null, item: ModelPricingBoardItem): string {
  if (value == null) return `- ${t('modelPricing.unitPerMillion')}`
  const symbol = isDomesticModel(item) ? '¥' : '$'
  return `${symbol}${formatPriceAmount(value * 1_000_000)} ${t('modelPricing.unitPerMillion')}`
}

function formatPriceAmount(value: number): string {
  if (!Number.isFinite(value)) return '-'
  const abs = Math.abs(value)
  const maximumFractionDigits = abs > 0 && abs < 0.01 ? 6 : 4
  return value.toLocaleString(undefined, {
    minimumFractionDigits: 0,
    maximumFractionDigits,
  })
}

function priceFormulaText(item: ModelPricingBoardItem): string {
  const key = isDomesticModel(item) ? 'modelPricing.card.priceFormulaDomestic' : 'modelPricing.card.priceFormula'
  return t(key, { rate: formatRate(item.rate_multiplier) })
}

function formatSavings(value: number | null): string {
  if (value == null) return '-'
  return `${(value * 100).toFixed(1).replace(/\.0$/, '')}%`
}

function formatRate(value: number): string {
  return Number.isFinite(value) ? value.toFixed(4).replace(/\.?0+$/, '') : '-'
}

function modelVendorLabel(item: Pick<ModelPricingBoardItem, 'platform' | 'model_id'>): string {
  const platform = item.platform.toLowerCase()
  const model = item.model_id.toLowerCase()
  const combined = `${platform} ${model}`
  if (combined.includes('deepseek')) return 'DeepSeek'
  if (combined.includes('zhipu') || combined.includes('chatglm') || combined.includes('glm')) return 'GLM'
  if (combined.includes('kimi') || combined.includes('moonshot')) return 'Kimi'
  if (combined.includes('minimax')) return 'MiniMax'
  if (combined.includes('doubao')) return 'Doubao'
  return platformLabel(item.platform)
}

function platformLabel(platform: string): string {
  const normalized = platform.toLowerCase()
  if (normalized === 'deepseek') return 'DeepSeek'
  if (normalized === 'zhipu' || normalized === 'glm') return 'GLM'
  if (normalized === 'kimi' || normalized === 'moonshot') return 'Kimi'
  if (normalized === 'minimax') return 'MiniMax'
  if (normalized === 'doubao') return 'Doubao'
  return t(`monitorCommon.providers.${platform}`, platform)
}

function platformBadgeClass(item: Pick<ModelPricingBoardItem, 'platform' | 'model_id'>): string {
  const base = 'inline-flex rounded-full border px-2.5 py-1 text-xs font-medium'
  if (isDomesticModel(item)) return `${base} border-emerald-200 bg-emerald-50 text-emerald-800 dark:border-emerald-900/60 dark:bg-emerald-950/30 dark:text-emerald-200`
  if (item.platform === 'anthropic') return `${base} border-amber-200 bg-amber-50 text-amber-800 dark:border-amber-900/60 dark:bg-amber-950/30 dark:text-amber-200`
  if (item.platform === 'openai') return `${base} border-gray-300 bg-gray-100 text-gray-800 dark:border-dark-600 dark:bg-dark-700 dark:text-gray-200`
  return `${base} border-primary-200 bg-primary-50 text-primary-800 dark:border-primary-900/60 dark:bg-primary-950/30 dark:text-primary-200`
}

function categoryLabel(item: ModelPricingBoardItem): string {
  const category = classifyItem(item)
  if (category === 'claude') return t('modelPricing.tabs.claude')
  if (category === 'codex') return t('modelPricing.tabs.codex')
  if (category === 'domestic') return t('modelPricing.tabs.domestic')
  return platformLabel(item.platform)
}

function categoryBadgeClass(item: ModelPricingBoardItem): string {
  const category = classifyItem(item)
  if (category === 'claude') {
    return 'border-amber-200 bg-amber-50 text-amber-800 dark:border-amber-900/60 dark:bg-amber-950/30 dark:text-amber-200'
  }
  if (category === 'codex') {
    return 'border-primary-200 bg-primary-100 text-primary-800 dark:border-primary-800 dark:bg-primary-950/40 dark:text-primary-200'
  }
  if (category === 'domestic') {
    return 'border-emerald-200 bg-emerald-50 text-emerald-800 dark:border-emerald-900/60 dark:bg-emerald-950/30 dark:text-emerald-200'
  }
  return 'border-primary-200 bg-primary-50 text-primary-800 dark:border-primary-900/60 dark:bg-primary-950/30 dark:text-primary-200'
}

async function loadBoard() {
  loading.value = true
  try {
    const res = await modelPricingAPI.getBoard()
    items.value = res.items || []
  } catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('modelPricing.loadError')))
  } finally {
    loading.value = false
  }
}

onMounted(loadBoard)

defineExpose({
  filteredItems,
  activeCategory,
  classifyItem,
  pricingRows,
})
</script>
