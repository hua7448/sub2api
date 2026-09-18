<template>
  <div class="hall-detail" data-testid="hall-detail">
    <header class="flex flex-wrap items-center justify-between gap-2">
      <div class="flex items-center gap-2">
        <span class="font-mono text-sm font-bold text-gray-900 dark:text-white">{{ row.name }}</span>
        <span class="badge" :class="healthClass">{{ t(`providerHall.health.${row.health.status}`) }}</span>
      </div>
      <div class="flex items-center gap-2">
        <HallVerificationBadge :verification="row.verification" @view-report="(id) => emit('view-report', { groupId: row.group_id, reportId: id })" />
        <button type="button" class="btn btn-secondary btn-sm" data-testid="detail-use-group" @click="emit('use-group')">
          {{ t('providerHall.row.useGroup') }}
        </button>
      </div>
    </header>

    <div v-if="loading && !detail" class="flex items-center gap-2 py-6 text-sm text-gray-500 dark:text-dark-400" data-testid="detail-loading">
      <LoadingSpinner size="sm" />
      {{ t('providerHall.detail.loading') }}
    </div>
    <div v-else-if="error && !detail" class="flex items-center gap-3 py-4 text-sm text-red-600 dark:text-red-400" data-testid="detail-error">
      <span>{{ error }}</span>
      <button type="button" class="btn btn-ghost btn-sm" @click="emit('retry')">{{ t('providerHall.retry') }}</button>
    </div>

    <div v-else-if="detail" class="mt-3 grid gap-4 lg:grid-cols-[minmax(0,1fr)_minmax(0,1.2fr)]">
      <!-- Metric grid -->
      <dl class="grid grid-cols-2 gap-px overflow-hidden rounded-xl border border-gray-200 bg-gray-200 text-xs dark:border-dark-700 dark:bg-dark-700 sm:grid-cols-4" data-testid="detail-metrics">
        <div v-for="cell in cells" :key="cell.key" class="bg-white p-2.5 dark:bg-dark-800" :data-cell="cell.key">
          <dt class="text-[10px] text-gray-500 dark:text-dark-400">{{ cell.label }}</dt>
          <dd class="mt-0.5 font-mono text-sm font-semibold tabular-nums" :class="cell.muted ? 'text-gray-400 dark:text-dark-500' : 'text-gray-900 dark:text-white'">
            {{ cell.value }}
            <span v-if="cell.stale" class="ml-1 text-[10px] font-normal text-amber-600 dark:text-amber-400">{{ t('providerHall.state.stale') }}</span>
          </dd>
        </div>
      </dl>

      <!-- Trend chart -->
      <section class="rounded-xl border border-gray-200 bg-white p-3 dark:border-dark-700 dark:bg-dark-800">
        <div class="flex items-center justify-between gap-2">
          <h3 class="text-[11px] font-semibold text-gray-600 dark:text-dark-300">{{ t('providerHall.detail.trendTitle') }}</h3>
          <span class="text-[10px] text-gray-400">{{ detail.trend.range }}</span>
        </div>
        <div v-if="chartData" class="mt-2 h-40" data-testid="detail-trend">
          <Line :data="chartData" :options="chartOptions" :plugins="[gapPlugin]" />
        </div>
        <p v-else class="mt-4 text-center text-xs text-gray-400" data-testid="detail-trend-empty">{{ t('providerHall.detail.trendEmpty') }}</p>
        <ul class="mt-2 flex flex-wrap items-center gap-3 text-[10px] text-gray-500 dark:text-dark-400">
          <li class="inline-flex items-center gap-1"><i class="hall-key" style="background: #eb6834" aria-hidden="true"></i>{{ t('providerHall.legend.probeE2E') }}</li>
          <li class="inline-flex items-center gap-1"><i class="hall-key" style="background: #1baf7a" aria-hidden="true"></i>{{ t('providerHall.legend.probeTtft') }}</li>
          <li class="inline-flex items-center gap-1"><i class="hall-key" style="background: #2a78d6" aria-hidden="true"></i>{{ t('providerHall.legend.fast95') }}</li>
          <li class="inline-flex items-center gap-1"><i class="hall-key" style="background: rgba(208,59,59,0.25); height: 8px" aria-hidden="true"></i>{{ t('providerHall.legend.gap') }}</li>
        </ul>
      </section>

      <!-- Model health + tokens -->
      <footer class="flex flex-wrap items-center justify-between gap-3 lg:col-span-2">
        <div class="flex flex-wrap items-center gap-2">
          <span class="text-[11px] text-gray-500 dark:text-dark-400">{{ t('providerHall.detail.profiles') }}</span>
          <HallModelDots :models="detail.profiles.map((p) => p.health)" chips />
        </div>
        <div v-if="detail.detail.probe_tokens" class="text-[11px] text-gray-500 dark:text-dark-400" data-testid="detail-tokens">
          {{ t('providerHall.detail.tokens') }}
          <span class="ml-1 font-semibold text-gray-900 dark:text-white">{{ t('providerHall.detail.tokenInput') }} {{ detail.detail.probe_tokens.input }}</span>
          <span class="ml-2 font-semibold text-gray-900 dark:text-white">{{ t('providerHall.detail.tokenOutput') }} {{ detail.detail.probe_tokens.output }}</span>
        </div>
      </footer>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  CategoryScale,
  Chart as ChartJS,
  Legend,
  LineElement,
  LinearScale,
  PointElement,
  Tooltip,
  type Plugin,
} from 'chart.js'
import { Line } from 'vue-chartjs'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import type { ProviderHallDetailResponse, ProviderHallMetric, ProviderHallRow } from '@/types/providerHall'
import {
  formatDateTime,
  formatMs,
  formatMultiplier,
  formatPrice,
  formatRate,
  formatShortTime,
  formatTps,
  metricHasValue,
  stateLabelKeys,
} from '@/utils/providerHallFormat'
import HallModelDots from './HallModelDots.vue'
import HallVerificationBadge from './HallVerificationBadge.vue'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Tooltip, Legend)

const props = defineProps<{
  row: ProviderHallRow
  detail: ProviderHallDetailResponse | null
  loading: boolean
  error: string | null
}>()

const emit = defineEmits<{
  (e: 'use-group'): void
  (e: 'view-report', payload: { groupId: number; reportId: number }): void
  (e: 'retry'): void
}>()

const { t, te, locale } = useI18n()

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

function placeholder(metric: ProviderHallMetric<unknown> | null | undefined): string {
  if (!metric) return '-'
  for (const key of stateLabelKeys(metric.state, metric.reason_code)) {
    if (te(key)) return t(key)
  }
  return t('providerHall.state.insufficient')
}

interface Cell {
  key: string
  label: string
  value: string
  muted?: boolean
  stale?: boolean
}

function metricCell<T>(key: string, label: string, metric: ProviderHallMetric<T> | null | undefined, format: (v: T) => string): Cell {
  if (metricHasValue(metric)) {
    return { key, label, value: format(metric.value), stale: metric.state === 'stale' }
  }
  return { key, label, value: placeholder(metric), muted: true }
}

const cells = computed<Cell[]>(() => {
  const d = props.detail
  if (!d) return []
  const quote = d.quote
  const priceCell: Cell = quote.applicable
    ? { key: 'input_price', label: t('providerHall.detail.inputPrice'), value: `${quote.unit === 'quota_per_million' ? '¥' : '$'}${formatPrice(quote.input_price)}/M` }
    : { key: 'input_price', label: t('providerHall.detail.inputPrice'), value: te(`providerHall.reason.${quote.reason_code}`) ? t(`providerHall.reason.${quote.reason_code}`) : t('providerHall.state.not_applicable'), muted: true }
  return [
    { key: 'rate', label: t('providerHall.detail.rate'), value: formatMultiplier(d.rate) },
    priceCell,
    metricCell('success_1h', t('providerHall.detail.success1h'), d.metrics.success_rate, (v) => formatRate(v, 0)),
    metricCell('availability_6h', t('providerHall.detail.availability6h'), d.detail.e2e_availability_6h, (v) => formatRate(v, 0)),
    metricCell('cache_1h', t('providerHall.detail.cache1h'), d.metrics.cache_rate, (v) => formatRate(v, 2)),
    metricCell('probe_ttft', t('providerHall.detail.probeTtft'), d.detail.probe_ttft_ms, (v) => formatMs(v, locale.value)),
    metricCell('probe_total', t('providerHall.detail.probeTotal'), d.detail.probe_total_ms, (v) => formatMs(v, locale.value)),
    metricCell('fast95', t('providerHall.detail.fast95'), d.metrics.ttft_fast95_ms, (v) => formatMs(v, locale.value)),
    metricCell('p90', t('providerHall.detail.p90', { model: d.default_profile.model }), d.detail.p90_ms, (v) => formatMs(v, locale.value)),
    metricCell('tps', t('providerHall.detail.tps'), d.detail.tps, (v) => formatTps(v)),
    {
      key: 'last_probe',
      label: t('providerHall.detail.lastProbe'),
      value: d.detail.last_probe_at ? formatDateTime(d.detail.last_probe_at, locale.value) : placeholder({ state: 'insufficient', reason_code: 'no_probe', value: null, window_start: null, window_end: null, computed_at: null }),
      muted: !d.detail.last_probe_at,
    },
  ]
})

const isDark = computed(() => typeof document !== 'undefined' && (document.documentElement.classList.contains('dark') || document.documentElement.classList.contains('dark-theme')))

const chartData = computed(() => {
  const points = props.detail?.trend.points ?? []
  if (!points.length) return null
  const anyValue = points.some((p) => p.real_fast95_ms != null || p.probe_total_ms != null || p.probe_ttft_ms != null)
  if (!anyValue) return null
  const labels = points.map((p) => formatShortTime(p.t, locale.value))
  const dark = isDark.value
  const common = { tension: 0.3, pointRadius: 0, pointHoverRadius: 4, pointHitRadius: 8, borderWidth: 2, spanGaps: false }
  return {
    labels,
    datasets: [
      { label: t('providerHall.legend.probeE2E'), data: points.map((p) => p.probe_total_ms), borderColor: dark ? '#d95926' : '#eb6834', backgroundColor: 'transparent', ...common },
      { label: t('providerHall.legend.probeTtft'), data: points.map((p) => p.probe_ttft_ms), borderColor: dark ? '#199e70' : '#1baf7a', backgroundColor: 'transparent', ...common, borderWidth: 1.5 },
      { label: t('providerHall.legend.fast95'), data: points.map((p) => p.real_fast95_ms), borderColor: dark ? '#3987e5' : '#2a78d6', backgroundColor: 'transparent', ...common },
    ],
  }
})

/** Draws a translucent wash over buckets flagged as data gaps. */
const gapPlugin = computed<Plugin<'line'>>(() => {
  const points = props.detail?.trend.points ?? []
  return {
    id: 'providerHallGaps',
    beforeDatasetsDraw(chart) {
      const x = chart.scales.x
      const area = chart.chartArea
      if (!x || !area || !points.length) return
      const ctx = chart.ctx
      ctx.save()
      ctx.fillStyle = 'rgba(208, 59, 59, 0.16)'
      let start = -1
      const flush = (end: number) => {
        if (start < 0) return
        const half = points.length > 1 ? (x.getPixelForValue(1) - x.getPixelForValue(0)) / 2 : 4
        const x0 = Math.max(area.left, x.getPixelForValue(start) - half)
        const x1 = Math.min(area.right, x.getPixelForValue(end) + half)
        ctx.fillRect(x0, area.top, Math.max(1, x1 - x0), area.bottom - area.top)
        start = -1
      }
      points.forEach((p, i) => {
        if (p.gap) {
          if (start < 0) start = i
        } else {
          flush(i - 1)
        }
      })
      flush(points.length - 1)
      ctx.restore()
    },
  }
})

const chartOptions = computed(() => {
  const dark = isDark.value
  const ink = dark ? '#c3c2b7' : '#52514e'
  const grid = dark ? '#2c2c2a' : '#e1e0d9'
  return {
    responsive: true,
    maintainAspectRatio: false,
    animation: false as const,
    interaction: { mode: 'index' as const, intersect: false },
    plugins: {
      legend: { display: false },
      tooltip: {
        callbacks: {
          label: (item: { dataset: { label?: string }; parsed: { y: number | null } }) =>
            `${item.dataset.label ?? ''}: ${item.parsed.y == null ? '-' : formatMs(item.parsed.y, locale.value)}`,
        },
      },
    },
    scales: {
      x: { ticks: { color: ink, maxTicksLimit: 6, font: { size: 10 } }, grid: { display: false }, border: { color: grid } },
      y: {
        ticks: { color: ink, font: { size: 10 }, callback: (v: string | number) => `${Number(v).toLocaleString(locale.value)}` },
        grid: { color: grid, lineWidth: 1 },
        border: { display: false },
        beginAtZero: true,
      },
    },
  }
})
</script>

<style scoped>
.hall-key {
  display: inline-block;
  width: 12px;
  height: 2px;
  border-radius: 1px;
}
</style>
