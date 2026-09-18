<template>
  <ul class="hall-model-dots m-0 flex flex-wrap gap-x-2 gap-y-1 p-0" :class="chips && 'gap-1.5'">
    <li
      v-for="item in models"
      :key="item.profile_id"
      class="inline-flex items-center gap-1 text-[11px] leading-none"
      :class="chips ? chipClass(item.status) : 'text-gray-600 dark:text-dark-300'"
      :title="tooltip(item)"
      :data-status="item.status"
      data-testid="model-dot"
    >
      <span class="max-w-[9rem] truncate">{{ item.model }}</span>
      <span
        class="inline-block h-1.5 w-1.5 shrink-0 rounded-full"
        :class="dotClass(item.status)"
        aria-hidden="true"
      ></span>
      <span class="sr-only">{{ t(`providerHall.health.${item.status}`) }}</span>
    </li>
  </ul>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { ProviderHallHealthStatus, ProviderHallModelHealth } from '@/types/providerHall'
import { formatDateTime } from '@/utils/providerHallFormat'

withDefaults(defineProps<{ models: ProviderHallModelHealth[]; chips?: boolean }>(), { chips: false })

const { t, locale } = useI18n()

function dotClass(status: ProviderHallHealthStatus): string {
  if (status === 'up') return 'bg-emerald-500'
  if (status === 'down') return 'bg-red-500'
  return 'bg-gray-400 dark:bg-dark-500'
}

function chipClass(status: ProviderHallHealthStatus): string {
  const base = 'rounded border px-1.5 py-1 '
  if (status === 'up') return base + 'border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-800 dark:bg-emerald-900/20 dark:text-emerald-300'
  if (status === 'down') return base + 'border-red-200 bg-red-50 text-red-700 dark:border-red-800 dark:bg-red-900/20 dark:text-red-300'
  return base + 'border-gray-200 bg-gray-50 text-gray-600 dark:border-dark-600 dark:bg-dark-800 dark:text-dark-300'
}

function tooltip(item: ProviderHallModelHealth): string {
  const status = t(`providerHall.health.${item.status}`)
  const when = item.checked_at
    ? t('providerHall.health.checkedAt', { time: formatDateTime(item.checked_at, locale.value) })
    : t('providerHall.health.neverChecked')
  return `${item.model} · ${item.protocol} · ${status} · ${when}`
}
</script>
