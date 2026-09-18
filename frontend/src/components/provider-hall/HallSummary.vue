<template>
  <section class="hall-summary flex flex-wrap items-center justify-between gap-4" data-testid="hall-summary">
    <div class="flex flex-wrap items-center gap-5">
      <div class="flex items-center gap-3">
        <div class="relative inline-flex h-12 w-12 items-center justify-center" role="img" :aria-label="`${t('providerHall.summary.available')} ${available}`">
          <svg width="48" height="48" viewBox="0 0 48 48" aria-hidden="true" class="absolute inset-0">
            <circle cx="24" cy="24" r="21" fill="none" stroke-width="3" class="hall-summary-track" />
            <circle
              cx="24"
              cy="24"
              r="21"
              fill="none"
              stroke-width="3"
              stroke-linecap="round"
              class="hall-summary-arc"
              :stroke-dasharray="`${arc} ${circumference}`"
              transform="rotate(-90 24 24)"
            />
          </svg>
          <div class="relative flex flex-col items-center leading-none">
            <span class="text-base font-black text-gray-900 dark:text-white" data-testid="summary-available">{{ available }}</span>
            <span class="text-[9px] text-gray-500 dark:text-dark-400">{{ t('providerHall.summary.available') }}</span>
          </div>
        </div>
      </div>
      <dl class="flex flex-wrap items-center gap-x-6 gap-y-2">
        <div class="border-l border-gray-200 pl-4 dark:border-dark-700">
          <dt class="text-[11px] text-gray-500 dark:text-dark-400">{{ t('providerHall.summary.listed') }}</dt>
          <dd class="text-base font-bold text-gray-900 dark:text-white" data-testid="summary-listed">{{ listed }}</dd>
        </div>
        <div class="border-l border-gray-200 pl-4 dark:border-dark-700">
          <dt class="text-[11px] text-gray-500 dark:text-dark-400">{{ t('providerHall.summary.abnormal') }}</dt>
          <dd class="text-base font-bold text-gray-900 dark:text-white" data-testid="summary-abnormal">{{ abnormal }}</dd>
        </div>
        <div class="border-l border-gray-200 pl-4 dark:border-dark-700">
          <dt class="text-[11px] text-gray-500 dark:text-dark-400">{{ t('providerHall.summary.verified') }}</dt>
          <dd class="text-base font-bold text-gray-900 dark:text-white" data-testid="summary-verified">{{ verified }}</dd>
        </div>
      </dl>
    </div>
    <div class="flex flex-wrap items-center gap-3">
      <ul class="m-0 flex items-center gap-3 p-0 text-[11px] text-gray-500 dark:text-dark-400" aria-label="legend">
        <li class="inline-flex items-center gap-1"><span class="h-2 w-2 rounded-full bg-emerald-500" aria-hidden="true"></span>{{ t('providerHall.legend.up') }}</li>
        <li class="inline-flex items-center gap-1"><span class="h-2 w-2 rounded-full bg-red-500" aria-hidden="true"></span>{{ t('providerHall.legend.down') }}</li>
        <li class="inline-flex items-center gap-1"><span class="h-2 w-2 rounded-full bg-gray-400" aria-hidden="true"></span>{{ t('providerHall.legend.unknown') }}</li>
      </ul>
      <slot name="actions" />
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { ProviderHallSummary } from '@/types/providerHall'

const props = defineProps<{ summary: ProviderHallSummary | null | undefined }>()
const { t } = useI18n()

const available = computed(() => props.summary?.available ?? 0)
const listed = computed(() => props.summary?.listed ?? 0)
const abnormal = computed(() => props.summary?.abnormal ?? 0)
const verified = computed(() => props.summary?.verified ?? 0)
const circumference = 2 * Math.PI * 21
const arc = computed(() => {
  if (!listed.value) return 0
  return Math.min(1, available.value / listed.value) * circumference
})
</script>

<style scoped>
.hall-summary-track {
  stroke: #e5e7eb;
}
.dark .hall-summary-track,
.dark-theme .hall-summary-track {
  stroke: #2c2c2a;
}
.hall-summary-arc {
  stroke: #0ca30c;
}
</style>
