<template>
  <BaseDialog :show="show" :title="t('providerHall.sort.dialogTitle')" width="narrow" @close="emit('close')">
    <p class="mb-4 text-xs text-gray-500 dark:text-dark-400">{{ t('providerHall.sort.dialogHint') }}</p>
    <div class="space-y-3">
      <div v-for="(level, index) in draft" :key="index" class="grid grid-cols-[4rem_1fr_7rem] items-center gap-2">
        <span class="text-xs font-medium text-gray-600 dark:text-dark-300">{{ t('providerHall.sort.level', { n: index + 1 }) }}</span>
        <Select
          :model-value="level.field"
          :options="fieldOptions(index)"
          :placeholder="t('providerHall.sort.none')"
          :aria-label="t('providerHall.sort.level', { n: index + 1 })"
          @update:model-value="(value) => setField(index, value)"
        />
        <Select
          :model-value="level.direction"
          :options="directionOptions"
          :disabled="!level.field"
          @update:model-value="(value) => setDirection(index, value)"
        />
      </div>
    </div>
    <template #footer>
      <button type="button" class="btn btn-ghost" @click="clear">{{ t('providerHall.sort.clear') }}</button>
      <button type="button" class="btn btn-secondary" @click="emit('close')">{{ t('providerHall.sort.cancel') }}</button>
      <button type="button" class="btn btn-primary" data-testid="sort-apply" @click="apply">{{ t('providerHall.sort.apply') }}</button>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Select from '@/components/common/Select.vue'
import {
  PROVIDER_HALL_SORT_FIELDS,
  type ProviderHallSortDirection,
  type ProviderHallSortField,
  type ProviderHallSortRule,
} from '@/types/providerHall'

const props = defineProps<{ show: boolean; rules: ProviderHallSortRule[] }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'apply', rules: ProviderHallSortRule[]): void }>()
const { t } = useI18n()

interface Level {
  field: ProviderHallSortField | ''
  direction: ProviderHallSortDirection
}

const draft = ref<Level[]>(seed(props.rules))

function seed(rules: ProviderHallSortRule[]): Level[] {
  const levels: Level[] = [0, 1, 2].map((i) => ({
    field: rules[i]?.field ?? '',
    direction: rules[i]?.direction ?? 'desc',
  }))
  return levels
}

watch(
  () => props.show,
  (open) => {
    if (open) draft.value = seed(props.rules)
  }
)

const directionOptions = computed(() => [
  { value: 'desc', label: t('providerHall.sort.desc') },
  { value: 'asc', label: t('providerHall.sort.asc') },
])

function fieldOptions(index: number) {
  const used = new Set(draft.value.filter((_, i) => i !== index).map((l) => l.field).filter(Boolean))
  return [
    { value: '', label: t('providerHall.sort.none') },
    ...PROVIDER_HALL_SORT_FIELDS.filter((f) => !used.has(f)).map((f) => ({ value: f, label: t(`providerHall.sort.${f}`) })),
  ]
}

function setField(index: number, value: string | number | boolean | null) {
  draft.value[index].field = (value ? String(value) : '') as ProviderHallSortField | ''
}
function setDirection(index: number, value: string | number | boolean | null) {
  draft.value[index].direction = value === 'asc' ? 'asc' : 'desc'
}
function clear() {
  draft.value = seed([])
}
function apply() {
  const rules = draft.value
    .filter((l): l is Level & { field: ProviderHallSortField } => Boolean(l.field))
    .map((l) => ({ field: l.field, direction: l.direction }))
  emit('apply', rules)
  emit('close')
}
</script>
