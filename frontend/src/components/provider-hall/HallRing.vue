<template>
  <div class="hall-ring inline-flex flex-col items-center" :aria-label="ariaLabel" role="img">
    <svg :width="size" :height="size" viewBox="0 0 44 44" class="block" aria-hidden="true">
      <circle
        cx="22"
        cy="22"
        :r="RING_RADIUS"
        fill="none"
        class="hall-ring-track"
        stroke-width="4"
      />
      <circle
        v-if="hasValue"
        cx="22"
        cy="22"
        :r="RING_RADIUS"
        fill="none"
        class="hall-ring-arc"
        :data-tone="tone"
        stroke-width="4"
        stroke-linecap="round"
        :stroke-dasharray="`${arc} ${RING_CIRCUMFERENCE}`"
        transform="rotate(-90 22 22)"
      />
      <text
        x="22"
        y="22"
        text-anchor="middle"
        dominant-baseline="central"
        class="hall-ring-text"
        :class="hasValue ? 'text-[10px] font-bold' : 'text-[9px] font-semibold'"
      >{{ display }}</text>
    </svg>
    <span v-if="label" class="mt-1 text-[10px] leading-none text-gray-500 dark:text-dark-400">{{ label }}</span>
    <span
      v-if="!hasValue && placeholder"
      class="mt-0.5 text-[10px] leading-none text-gray-400 dark:text-dark-500"
      data-testid="ring-placeholder"
    >{{ placeholder }}</span>
    <span
      v-else-if="metric?.state === 'stale'"
      class="mt-0.5 text-[10px] leading-none text-amber-600 dark:text-amber-400"
      data-testid="ring-stale"
    >{{ staleText }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { ProviderHallMetric } from '@/types/providerHall'
import { formatRate, metricHasValue, stateLabelKeys } from '@/utils/providerHallFormat'
import { RING_CIRCUMFERENCE, RING_RADIUS, ringArc, successTone, type RingTone } from './hallGeometry'

const props = withDefaults(
  defineProps<{
    metric: ProviderHallMetric<number> | null | undefined
    label?: string
    /** 'score' colours by success-rate score; 'plain' always uses the series hue. */
    mode?: 'score' | 'plain'
    size?: number
  }>(),
  { label: '', mode: 'plain', size: 44 }
)

const { t, te } = useI18n()

const hasValue = computed(() => metricHasValue(props.metric))
const arc = computed(() => (hasValue.value ? ringArc(props.metric!.value) : 0))
const tone = computed<RingTone>(() => {
  if (!hasValue.value) return 'muted'
  return props.mode === 'score' ? successTone(props.metric!.value) : 'series'
})
const display = computed(() => (hasValue.value ? formatRate(props.metric!.value, 1) : '-'))
const placeholder = computed(() => {
  if (hasValue.value || !props.metric) return ''
  for (const key of stateLabelKeys(props.metric.state, props.metric.reason_code)) {
    if (te(key)) return t(key)
  }
  return t('providerHall.state.insufficient')
})
const staleText = computed(() => t('providerHall.state.stale'))
const ariaLabel = computed(() => t('providerHall.aria.ring', { label: props.label, value: hasValue.value ? display.value : placeholder.value }))
</script>

<style scoped>
.hall-ring-track {
  stroke: #e5e7eb;
}
.dark .hall-ring-track,
.dark-theme .hall-ring-track {
  stroke: #2c2c2a;
}
.hall-ring-text {
  fill: #111827;
}
.dark .hall-ring-text,
.dark-theme .hall-ring-text {
  fill: #ffffff;
}
/* Status palette (dataviz skill): good / warning / critical; series slot 1 blue. */
.hall-ring-arc[data-tone='good'] {
  stroke: #0ca30c;
}
.hall-ring-arc[data-tone='series'] {
  stroke: #2a78d6;
}
.dark .hall-ring-arc[data-tone='series'],
.dark-theme .hall-ring-arc[data-tone='series'] {
  stroke: #3987e5;
}
.hall-ring-arc[data-tone='warning'] {
  stroke: #fab219;
}
.hall-ring-arc[data-tone='critical'] {
  stroke: #d03b3b;
}
.hall-ring-arc[data-tone='muted'] {
  stroke: #9ca3af;
}
</style>
