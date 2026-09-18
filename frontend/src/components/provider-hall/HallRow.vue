<template>
  <!-- Desktop row -->
  <tr
    v-if="variant === 'row'"
    class="hall-row group border-b border-gray-100 align-top transition-colors last:border-b-0 hover:bg-gray-50/60 dark:border-dark-700 dark:hover:bg-dark-700/40"
    :class="expanded && 'bg-blue-50/40 dark:bg-blue-900/10'"
    :data-group-id="row.group_id"
    :aria-expanded="expanded"
    tabindex="0"
    role="row"
    :aria-label="expanded ? t('providerHall.aria.collapseRow', { name: row.name }) : t('providerHall.aria.expandRow', { name: row.name })"
    data-testid="hall-row"
    @click="emit('toggle', row.group_id)"
    @keydown.enter.prevent="emit('toggle', row.group_id)"
    @keydown.space.prevent="emit('toggle', row.group_id)"
  >
    <td class="px-3 py-3" :class="expanded && 'border-l-2 border-primary-500'">
      <div class="font-semibold text-gray-900 dark:text-white">{{ row.name }}</div>
      <div class="mt-0.5 line-clamp-2 text-[11px] text-gray-500 dark:text-dark-400">{{ row.description || t('providerHall.row.noDescription') }}</div>
    </td>
    <td class="px-3 py-3">
      <span class="text-base font-bold tabular-nums text-gray-900 dark:text-white" data-testid="row-rate">{{ formatMultiplier(row.rate) }}</span>
    </td>
    <td class="px-3 py-3">
      <div class="space-y-1">
        <HallMetricTooltip
          :title="t('providerHall.tooltip.historicalPriceTitle')"
          :paragraphs="[t('providerHall.tooltip.historicalPriceBody', { model: row.default_profile.model })]"
          :metric="row.historical_price"
        >
          <span v-if="metricHasValue(row.historical_price)" class="inline-flex items-baseline gap-1 tabular-nums" data-testid="row-price">
            <span class="text-[10px] text-gray-400">{{ currencyPrefix(row.historical_price) }}</span>
            <span class="font-bold text-gray-900 dark:text-white">{{ formatPrice(row.historical_price.value) }}</span>
            <span class="text-[10px] text-gray-400">{{ t('providerHall.row.perMillion') }}</span>
            <span v-if="row.historical_price.state === 'stale'" class="text-[10px] text-amber-600 dark:text-amber-400">{{ t('providerHall.state.stale') }}</span>
          </span>
          <span v-else class="text-xs text-gray-400 dark:text-dark-500" data-testid="row-price-placeholder">{{ placeholder(row.historical_price) }}</span>
        </HallMetricTooltip>
        <div class="border-t border-gray-100 dark:border-dark-700"></div>
        <HallMetricTooltip
          :title="t('providerHall.tooltip.predictedRateTitle')"
          :paragraphs="[t('providerHall.tooltip.predictedRateBody')]"
          :footnote="t('providerHall.tooltip.predictedRateFootnote')"
          :metric="row.predicted_rate"
        >
          <span v-if="metricHasValue(row.predicted_rate)" class="inline-flex items-baseline gap-1 tabular-nums" data-testid="row-predicted">
            <span class="font-bold text-gray-900 dark:text-white">{{ formatDecimal(row.predicted_rate.value, 4, 2) }}</span>
            <span class="text-[10px] text-gray-400">{{ t('providerHall.row.multiplierSuffix') }}</span>
            <span v-if="row.predicted_rate.state === 'stale'" class="text-[10px] text-amber-600 dark:text-amber-400">{{ t('providerHall.state.stale') }}</span>
          </span>
          <span v-else class="text-xs text-gray-400 dark:text-dark-500" data-testid="row-predicted-placeholder">{{ placeholder(row.predicted_rate) }}</span>
        </HallMetricTooltip>
      </div>
    </td>
    <td class="px-3 py-3">
      <span class="badge" :class="healthClass" data-testid="row-health">
        <span class="h-1.5 w-1.5 rounded-full" :class="healthDot" aria-hidden="true"></span>
        {{ t(`providerHall.health.${row.health.status}`) }}
      </span>
      <HallModelDots class="mt-1.5" :models="row.health.models" />
    </td>
    <td class="px-3 py-3">
      <HallMetricTooltip
        :title="t('providerHall.tooltip.ttftTitle')"
        :paragraphs="[t('providerHall.tooltip.ttftBody')]"
        :method="[t('providerHall.tooltip.ttftMethod'), t('providerHall.tooltip.ttftExample')]"
        :metric="row.metrics.ttft_fast95_ms"
      >
        <HallTtftBar :metric="row.metrics.ttft_fast95_ms" :max="ttftMax" />
      </HallMetricTooltip>
    </td>
    <td class="px-3 py-3 text-center">
      <HallMetricTooltip
        :title="t('providerHall.tooltip.cacheTitle')"
        :paragraphs="[t('providerHall.tooltip.cacheBody', { model: row.default_profile.model })]"
        :method="[t('providerHall.tooltip.cacheMethod')]"
        :metric="row.metrics.cache_rate"
      >
        <HallRing :metric="row.metrics.cache_rate" :label="t('providerHall.row.cache1h')" mode="plain" />
      </HallMetricTooltip>
    </td>
    <td class="px-3 py-3 text-center">
      <HallMetricTooltip
        :title="t('providerHall.tooltip.successTitle')"
        :paragraphs="[t('providerHall.tooltip.successBody')]"
        :method="[t('providerHall.tooltip.successMethod'), t('providerHall.tooltip.successScore')]"
        :metric="row.metrics.success_rate"
      >
        <HallRing :metric="row.metrics.success_rate" :label="t('providerHall.row.success1h')" mode="score" />
      </HallMetricTooltip>
    </td>
    <td class="px-3 py-3">
      <HallSparkline :points="row.sparkline.points" :name="row.name" />
    </td>
    <td class="px-3 py-3 text-center">
      <HallVerificationBadge :verification="row.verification" @view-report="onViewReport" />
    </td>
    <td class="px-3 py-3 text-right">
      <button
        type="button"
        class="btn btn-secondary btn-sm whitespace-nowrap"
        data-testid="use-group"
        @click.stop="emit('use-group', row)"
      >{{ t('providerHall.row.useGroup') }}</button>
    </td>
  </tr>

  <!-- Mobile card -->
  <article
    v-else
    class="hall-card rounded-2xl border border-gray-200 bg-white p-4 shadow-sm dark:border-dark-700 dark:bg-dark-800"
    :class="expanded && 'ring-2 ring-primary-300 dark:ring-primary-700'"
    :data-group-id="row.group_id"
    data-testid="hall-card"
  >
    <header class="flex items-start justify-between gap-2">
      <button
        type="button"
        class="min-w-0 flex-1 text-left"
        :aria-expanded="expanded"
        data-testid="card-toggle"
        @click="emit('toggle', row.group_id)"
      >
        <div class="truncate font-semibold text-gray-900 dark:text-white">{{ row.name }}</div>
        <div class="mt-0.5 line-clamp-2 text-[11px] text-gray-500 dark:text-dark-400">{{ row.description || t('providerHall.row.noDescription') }}</div>
      </button>
      <span class="badge shrink-0" :class="healthClass">
        <span class="h-1.5 w-1.5 rounded-full" :class="healthDot" aria-hidden="true"></span>
        {{ t(`providerHall.health.${row.health.status}`) }}
      </span>
    </header>
    <dl class="mt-3 grid grid-cols-2 gap-x-3 gap-y-2 text-xs">
      <div>
        <dt class="text-[10px] text-gray-500 dark:text-dark-400">{{ t('providerHall.columns.rate') }}</dt>
        <dd class="font-bold tabular-nums text-gray-900 dark:text-white">{{ formatMultiplier(row.rate) }}</dd>
      </div>
      <div>
        <dt class="text-[10px] text-gray-500 dark:text-dark-400">{{ t('providerHall.columns.price') }}</dt>
        <dd class="tabular-nums text-gray-900 dark:text-white">
          <template v-if="metricHasValue(row.historical_price)">{{ currencyPrefix(row.historical_price) }}{{ formatPrice(row.historical_price.value) }}{{ t('providerHall.row.perMillion') }}</template>
          <template v-else>{{ placeholder(row.historical_price) }}</template>
          <span class="text-gray-400"> / </span>
          <template v-if="metricHasValue(row.predicted_rate)">{{ formatDecimal(row.predicted_rate.value, 4, 2) }}{{ t('providerHall.row.multiplierSuffix') }}</template>
          <template v-else>{{ placeholder(row.predicted_rate) }}</template>
        </dd>
      </div>
      <div class="col-span-2">
        <dt class="text-[10px] text-gray-500 dark:text-dark-400">{{ t('providerHall.columns.ttft') }}</dt>
        <dd><HallTtftBar :metric="row.metrics.ttft_fast95_ms" :max="ttftMax" /></dd>
      </div>
      <div class="flex items-center gap-4">
        <HallRing :metric="row.metrics.cache_rate" :label="t('providerHall.row.cache1h')" mode="plain" />
        <HallRing :metric="row.metrics.success_rate" :label="t('providerHall.row.success1h')" mode="score" />
      </div>
      <div class="flex flex-col items-end justify-center gap-1">
        <HallVerificationBadge :verification="row.verification" @view-report="onViewReport" />
      </div>
      <div class="col-span-2 overflow-x-auto">
        <HallSparkline :points="row.sparkline.points" :name="row.name" :width="260" />
      </div>
      <div class="col-span-2">
        <HallModelDots :models="row.health.models" />
      </div>
    </dl>
    <footer class="mt-3 flex items-center justify-between gap-2">
      <button type="button" class="btn btn-ghost btn-sm" @click="emit('toggle', row.group_id)">
        {{ expanded ? t('providerHall.row.collapse') : t('providerHall.row.expand') }}
      </button>
      <button type="button" class="btn btn-secondary btn-sm" data-testid="use-group" @click="emit('use-group', row)">
        {{ t('providerHall.row.useGroup') }}
      </button>
    </footer>
  </article>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { ProviderHallMetric, ProviderHallRow } from '@/types/providerHall'
import { formatDecimal, formatMultiplier, formatPrice, metricHasValue, stateLabelKeys } from '@/utils/providerHallFormat'
import HallMetricTooltip from './HallMetricTooltip.vue'
import HallModelDots from './HallModelDots.vue'
import HallRing from './HallRing.vue'
import HallSparkline from './HallSparkline.vue'
import HallTtftBar from './HallTtftBar.vue'
import HallVerificationBadge from './HallVerificationBadge.vue'

const props = withDefaults(
  defineProps<{
    row: ProviderHallRow
    ttftMax: number
    expanded: boolean
    variant?: 'row' | 'card'
  }>(),
  { variant: 'row' }
)

const emit = defineEmits<{
  (e: 'toggle', groupId: number): void
  (e: 'use-group', row: ProviderHallRow): void
  (e: 'view-report', payload: { groupId: number; reportId: number }): void
}>()

const { t, te } = useI18n()

const healthClass = computed(() => {
  switch (props.row.health.status) {
    case 'up':
      return 'badge-success'
    case 'down':
      return 'badge-danger'
    default:
      return 'badge-gray'
  }
})
const healthDot = computed(() => {
  switch (props.row.health.status) {
    case 'up':
      return 'bg-emerald-500'
    case 'down':
      return 'bg-red-500'
    default:
      return 'bg-gray-400'
  }
})

function placeholder(metric: ProviderHallMetric<unknown>): string {
  for (const key of stateLabelKeys(metric.state, metric.reason_code)) {
    if (te(key)) return t(key)
  }
  return t('providerHall.state.insufficient')
}

function currencyPrefix(metric: ProviderHallMetric<string>): string {
  // Unit comes from the quote (same default profile); quota-denominated
  // history shows the local quota glyph, USD shows a dollar sign.
  const unit = props.row.quote.unit
  if (unit === 'quota_per_million') return '¥'
  if (metric.reason_code === 'quota_per_million') return '¥'
  return '$'
}

function onViewReport(reportId: number) {
  emit('view-report', { groupId: props.row.group_id, reportId })
}
</script>
