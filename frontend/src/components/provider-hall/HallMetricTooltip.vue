<template>
  <HelpTooltip width-class="w-80 max-w-[90vw] text-left" :trigger="trigger">
    <template #trigger>
      <slot name="trigger">
        <span class="inline-flex cursor-help items-center border-b border-dotted border-gray-300 dark:border-dark-500">
          <slot />
        </span>
      </slot>
    </template>
    <div class="space-y-1.5">
      <p class="font-semibold text-white">{{ title }}</p>
      <p v-for="(line, index) in paragraphs" :key="`p-${index}`" class="text-gray-200">{{ line }}</p>
      <template v-if="method && method.length">
        <p class="font-semibold text-sky-300">{{ t('providerHall.tooltip.calculation') }}</p>
        <p v-for="(line, index) in method" :key="`m-${index}`" class="font-mono text-[11px] text-gray-200">{{ line }}</p>
      </template>
      <p v-if="footnote" class="text-gray-300">{{ footnote }}</p>
      <p v-if="metric && (metric.window_start || metric.window_end)" class="text-gray-400" data-testid="metric-window">
        {{ t('providerHall.tooltip.window', { start: formatShortTime(metric.window_start, locale), end: formatShortTime(metric.window_end, locale) }) }}
      </p>
      <p v-if="metric?.computed_at" class="text-gray-400">
        {{ t('providerHall.tooltip.computedAt', { time: formatDateTime(metric.computed_at, locale) }) }}
      </p>
      <p v-if="metric && metric.state !== 'ok'" class="text-amber-300" data-testid="metric-state">
        {{ t('providerHall.tooltip.stateLabel', { state: stateText }) }}
      </p>
    </div>
  </HelpTooltip>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import HelpTooltip from '@/components/common/HelpTooltip.vue'
import type { ProviderHallMetric } from '@/types/providerHall'
import { formatDateTime, formatShortTime, stateLabelKeys } from '@/utils/providerHallFormat'

const props = withDefaults(
  defineProps<{
    title: string
    paragraphs?: string[]
    method?: string[]
    footnote?: string
    metric?: ProviderHallMetric<unknown> | null
    trigger?: 'hover' | 'click'
  }>(),
  { paragraphs: () => [], method: () => [], footnote: '', metric: null, trigger: 'hover' }
)

const { t, te, locale } = useI18n()

const stateText = computed(() => {
  if (!props.metric) return ''
  for (const key of stateLabelKeys(props.metric.state, props.metric.reason_code)) {
    if (te(key)) return t(key)
  }
  return props.metric.reason_code || props.metric.state
})
</script>
