<template>
  <div class="hall-ttft min-w-[120px]">
    <div class="flex items-baseline justify-between gap-2">
      <span
        class="text-base font-bold tabular-nums text-gray-900 dark:text-white"
        :class="!hasValue && 'text-sm font-medium text-gray-400 dark:text-dark-500'"
        data-testid="ttft-value"
      >{{ display }}</span>
      <span
        v-if="speed"
        class="rounded px-1.5 py-0.5 text-[10px] font-semibold leading-none"
        :class="speedClass"
        :data-speed="speed"
        data-testid="ttft-speed"
      >{{ t(`providerHall.speed.${speed}`) }}</span>
      <span
        v-else-if="metric?.state === 'stale'"
        class="text-[10px] text-amber-600 dark:text-amber-400"
      >{{ t('providerHall.state.stale') }}</span>
    </div>
    <div class="hall-ttft-track mt-1.5 h-1.5 w-full overflow-hidden rounded-full" aria-hidden="true">
      <div
        v-if="hasValue"
        class="hall-ttft-bar h-full rounded-full"
        :data-speed="speed"
        :style="{ width: `${barPercent}%` }"
        data-testid="ttft-bar"
      ></div>
    </div>
    <span
      v-if="metric?.state === 'stale' && speed"
      class="mt-0.5 block text-[10px] leading-none text-amber-600 dark:text-amber-400"
    >{{ t('providerHall.state.stale') }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { ProviderHallMetric } from '@/types/providerHall'
import { formatMs, metricHasValue, stateLabelKeys, ttftSpeed } from '@/utils/providerHallFormat'

const props = defineProps<{
  metric: ProviderHallMetric<number> | null | undefined
  /** Largest fast95 value on the current page, for the relative bar. */
  max: number
}>()

const { t, te, locale } = useI18n()
const hasValue = computed(() => metricHasValue(props.metric))
const speed = computed(() => (hasValue.value ? ttftSpeed(props.metric!.value) : null))
const display = computed(() => {
  if (hasValue.value) return formatMs(props.metric!.value, locale.value)
  if (!props.metric) return '-'
  for (const key of stateLabelKeys(props.metric.state, props.metric.reason_code)) {
    if (te(key)) return t(key)
  }
  return t('providerHall.state.insufficient')
})
const barPercent = computed(() => {
  const value = props.metric?.value
  if (!hasValue.value || value == null || !props.max || props.max <= 0) return 0
  return Math.max(4, Math.min(100, (value / props.max) * 100))
})
const speedClass = computed(() => {
  switch (speed.value) {
    case 'fast':
      return 'bg-green-50 text-green-700 dark:bg-green-900/30 dark:text-green-300'
    case 'medium':
      return 'bg-amber-50 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300'
    default:
      return 'bg-red-50 text-red-700 dark:bg-red-900/30 dark:text-red-300'
  }
})
</script>

<style scoped>
.hall-ttft-track {
  background: #e5e7eb;
}
.dark .hall-ttft-track,
.dark-theme .hall-ttft-track {
  background: #2c2c2a;
}
/* Meter: fill carries the speed class (good / warning / critical). */
.hall-ttft-bar[data-speed='fast'] {
  background: #0ca30c;
}
.hall-ttft-bar[data-speed='medium'] {
  background: #fab219;
}
.hall-ttft-bar[data-speed='slow'] {
  background: #d03b3b;
}
</style>
