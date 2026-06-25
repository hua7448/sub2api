<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex flex-col justify-between gap-4 lg:flex-row lg:items-start">
          <div class="flex flex-1 flex-wrap items-center gap-3">
            <div class="relative w-full sm:w-80">
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
          </div>

          <div class="flex w-full flex-shrink-0 flex-wrap items-center justify-end gap-3 lg:w-auto">
            <button
              @click="loadBoard"
              :disabled="loading"
              class="btn btn-secondary"
              :title="t('common.refresh', 'Refresh')"
            >
              <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
            </button>
          </div>
        </div>
      </template>

      <template #table>
        <div class="table-wrapper">
          <table>
            <thead>
              <tr>
                <th>{{ t('modelPricing.columns.model') }}</th>
                <th>{{ t('modelPricing.columns.platform') }}</th>
                <th>{{ t('modelPricing.columns.source') }}</th>
                <th>{{ t('modelPricing.columns.input') }}</th>
                <th>{{ t('modelPricing.columns.output') }}</th>
                <th>{{ t('modelPricing.columns.cacheRead') }}</th>
                <th>{{ t('modelPricing.columns.official') }}</th>
                <th>{{ t('modelPricing.columns.savings') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="loading">
                <td colspan="8" class="py-10 text-center text-sm text-gray-500 dark:text-gray-400">
                  {{ t('common.loading') }}
                </td>
              </tr>
              <tr v-else-if="filteredItems.length === 0">
                <td colspan="8" class="py-10 text-center text-sm text-gray-500 dark:text-gray-400">
                  {{ t('modelPricing.empty') }}
                </td>
              </tr>
              <template v-else>
                <tr v-for="item in filteredItems" :key="`${item.platform}:${item.model_id}`">
                  <td>
                    <div class="flex min-w-72 items-center gap-3">
                      <ModelIcon :model="item.model_id" size="20px" />
                      <span class="font-mono text-sm font-medium text-gray-900 dark:text-white">
                        {{ item.model_id }}
                      </span>
                    </div>
                  </td>
                  <td>
                    <span :class="platformBadgeClass(item.platform)">
                      {{ platformLabel(item.platform) }}
                    </span>
                  </td>
                  <td>
                    <div class="min-w-52 space-y-1">
                      <div class="font-medium text-gray-900 dark:text-white">{{ item.group_name }}</div>
                      <div class="text-xs text-gray-500 dark:text-gray-400">{{ item.channel_name }}</div>
                      <div class="flex items-center gap-2 text-xs">
                        <span class="rounded border border-gray-200 px-1.5 py-0.5 text-gray-500 dark:border-dark-700 dark:text-gray-400">
                          x{{ formatRate(item.rate_multiplier) }}
                        </span>
                        <span
                          v-if="item.user_rate_overridden"
                          class="rounded border border-emerald-200 bg-emerald-50 px-1.5 py-0.5 text-emerald-700 dark:border-emerald-900/60 dark:bg-emerald-950/30 dark:text-emerald-300"
                        >
                          {{ t('modelPricing.userRate') }}
                        </span>
                      </div>
                    </div>
                  </td>
                  <td>{{ formatTokenPrice(item.site_input_price) }}</td>
                  <td>{{ formatTokenPrice(item.site_output_price) }}</td>
                  <td>{{ formatTokenPrice(item.site_cache_read_price) }}</td>
                  <td>
                    <div class="min-w-44 space-y-1 text-xs">
                      <div>{{ t('modelPricing.official.input') }} {{ formatTokenPrice(item.official_input_price) }}</div>
                      <div>{{ t('modelPricing.official.output') }} {{ formatTokenPrice(item.official_output_price) }}</div>
                      <div>{{ t('modelPricing.official.cacheRead') }} {{ formatTokenPrice(item.official_cache_read_price) }}</div>
                    </div>
                  </td>
                  <td>
                    <div class="min-w-32 space-y-1 text-xs">
                      <div>{{ t('modelPricing.official.input') }} {{ formatSavings(item.input_savings) }}</div>
                      <div>{{ t('modelPricing.official.output') }} {{ formatSavings(item.output_savings) }}</div>
                      <div>{{ t('modelPricing.official.cacheRead') }} {{ formatSavings(item.cache_read_savings) }}</div>
                    </div>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>
      </template>
    </TablePageLayout>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import ModelIcon from '@/components/common/ModelIcon.vue'
import modelPricingAPI, { type ModelPricingBoardItem } from '@/api/modelPricing'
import { useAppStore } from '@/stores/app'
import { extractApiErrorMessage } from '@/utils/apiError'
import { formatScaled } from '@/utils/pricing'

const { t } = useI18n()
const appStore = useAppStore()

const items = ref<ModelPricingBoardItem[]>([])
const loading = ref(false)
const searchQuery = ref('')

const filteredItems = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return items.value
  return items.value.filter((item) =>
    item.model_id.toLowerCase().includes(q) ||
    item.platform.toLowerCase().includes(q) ||
    item.group_name.toLowerCase().includes(q) ||
    item.channel_name.toLowerCase().includes(q),
  )
})

function formatTokenPrice(value: number | null): string {
  return `${formatScaled(value, 1_000_000)} ${t('modelPricing.unitPerMillion')}`
}

function formatSavings(value: number | null): string {
  if (value == null) return '-'
  return `${(value * 100).toFixed(1).replace(/\.0$/, '')}%`
}

function formatRate(value: number): string {
  return Number.isFinite(value) ? value.toFixed(4).replace(/\.?0+$/, '') : '-'
}

function platformLabel(platform: string): string {
  return t(`monitorCommon.providers.${platform}`, platform)
}

function platformBadgeClass(platform: string): string {
  const base = 'inline-flex rounded px-2 py-1 text-xs font-medium'
  if (platform === 'anthropic') return `${base} bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-200`
  if (platform === 'openai') return `${base} bg-gray-100 text-gray-800 dark:bg-dark-700 dark:text-gray-200`
  return `${base} bg-primary-100 text-primary-800 dark:bg-primary-900/30 dark:text-primary-200`
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
</script>
