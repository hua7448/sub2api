<template>
  <div class="hall-toolbar flex flex-wrap items-center gap-x-4 gap-y-2" data-testid="hall-toolbar">
    <div class="flex items-center gap-2">
      <span class="text-xs text-gray-500 dark:text-dark-400">{{ t('providerHall.toolbar.range') }}</span>
      <nav class="tabs inline-flex !p-0.5" role="group" :aria-label="t('providerHall.aria.ranges')">
        <button
          v-for="value in ranges"
          :key="value"
          type="button"
          class="tab !px-2.5 !py-1 !text-xs"
          :class="range === value ? 'tab-active' : ''"
          :aria-pressed="range === value"
          :data-range="value"
          @click="emit('update:range', value)"
        >{{ t(`providerHall.ranges.${value}`) }}</button>
      </nav>
    </div>

    <div class="flex flex-wrap items-center gap-2">
      <span class="text-xs text-gray-500 dark:text-dark-400">{{ t('providerHall.toolbar.sort') }}</span>
      <nav class="tabs inline-flex flex-wrap !p-0.5" role="group" :aria-label="t('providerHall.aria.sorts')">
        <button
          type="button"
          class="tab !px-2.5 !py-1 !text-xs"
          :class="!primarySort && !isCustomSort ? 'tab-active' : ''"
          data-sort="default"
          @click="emit('clear-sort')"
        >{{ t('providerHall.sort.default') }}</button>
        <button
          v-for="field in sortChips"
          :key="field"
          type="button"
          class="tab inline-flex items-center gap-1 !px-2.5 !py-1 !text-xs"
          :class="!isCustomSort && primarySort?.field === field ? 'tab-active' : ''"
          :data-sort="field"
          :aria-pressed="!isCustomSort && primarySort?.field === field"
          @click="emit('cycle-sort', field)"
        >
          {{ t(`providerHall.sort.${field}`) }}
          <span v-if="!isCustomSort && primarySort?.field === field" aria-hidden="true" data-testid="sort-arrow">
            {{ primarySort?.direction === 'asc' ? '↑' : '↓' }}
          </span>
        </button>
        <button
          type="button"
          class="tab !px-2.5 !py-1 !text-xs"
          :class="isCustomSort ? 'tab-active' : ''"
          data-sort="custom"
          @click="emit('open-custom')"
        >{{ t('providerHall.sort.custom') }}</button>
      </nav>
    </div>

    <div class="flex items-center gap-2">
      <span class="text-xs text-gray-500 dark:text-dark-400">{{ t('providerHall.toolbar.model') }}</span>
      <div class="w-56">
        <Select
          :model-value="model"
          :options="modelOptions"
          :placeholder="t('providerHall.toolbar.allModels')"
          :aria-label="t('providerHall.toolbar.model')"
          clearable
          searchable="auto"
          @update:model-value="onModelChange"
        />
      </div>
    </div>

    <div class="flex min-w-[10rem] flex-1 items-center gap-2 sm:max-w-xs">
      <label class="sr-only" for="hall-search">{{ t('providerHall.toolbar.search') }}</label>
      <input
        id="hall-search"
        :value="search"
        type="search"
        class="input !py-1.5 text-sm"
        :placeholder="t('providerHall.toolbar.searchPlaceholder')"
        maxlength="100"
        data-testid="hall-search"
        @input="onSearchInput"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import Select from '@/components/common/Select.vue'
import {
  PROVIDER_HALL_RANGES,
  type ProviderHallProfileRef,
  type ProviderHallRange,
  type ProviderHallSortField,
  type ProviderHallSortRule,
} from '@/types/providerHall'
import { modelFilterValue } from '@/composables/useProviderHall'

const props = defineProps<{
  range: ProviderHallRange
  primarySort: ProviderHallSortRule | null
  isCustomSort: boolean
  model: string
  models: ProviderHallProfileRef[]
  search: string
}>()

const emit = defineEmits<{
  (e: 'update:range', value: ProviderHallRange): void
  (e: 'cycle-sort', field: ProviderHallSortField): void
  (e: 'clear-sort'): void
  (e: 'open-custom'): void
  (e: 'update:model', value: string): void
  (e: 'update:search', value: string): void
}>()

const { t } = useI18n()
const ranges = PROVIDER_HALL_RANGES
const sortChips: ProviderHallSortField[] = ['rate', 'historical_price', 'ttft_fast95', 'cache_rate', 'success_rate']

const modelOptions = computed(() =>
  (props.models || []).map((ref) => ({
    value: modelFilterValue(ref),
    label: `${ref.model} · ${ref.protocol}`,
  }))
)

function onModelChange(value: string | number | boolean | null) {
  emit('update:model', value == null ? '' : String(value))
}

let searchTimer: number | undefined
function onSearchInput(event: Event) {
  const value = (event.target as HTMLInputElement).value
  if (searchTimer !== undefined) window.clearTimeout(searchTimer)
  searchTimer = window.setTimeout(() => emit('update:search', value), 300)
}
onBeforeUnmount(() => {
  if (searchTimer !== undefined) window.clearTimeout(searchTimer)
})
</script>
