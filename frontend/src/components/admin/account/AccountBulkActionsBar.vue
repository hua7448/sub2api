<template>
  <div class="mb-4 flex items-center justify-between rounded-lg bg-primary-50 p-3 dark:bg-primary-900/20">
    <div class="flex flex-wrap items-center gap-2">
      <span v-if="allResultsSelected" class="text-sm font-medium text-primary-900 dark:text-primary-100">
        {{ t('admin.accounts.bulkActions.selectedAll', { count: selectedIds.length }) }}
      </span>
      <span v-else-if="selectedIds.length > 0" class="text-sm font-medium text-primary-900 dark:text-primary-100">
        {{ t('admin.accounts.bulkActions.selected', { count: selectedIds.length }) }}
      </span>
      <span v-else class="text-sm font-medium text-primary-900 dark:text-primary-100">
        {{ t('admin.accounts.bulkEdit.title') }}
      </span>
      <template v-if="selectedIds.length > 0">
        <button
          @click="$emit('select-page')"
          class="text-xs font-medium text-primary-700 hover:text-primary-800 dark:text-primary-300 dark:hover:text-primary-200"
        >
          {{ t('admin.accounts.bulkActions.selectCurrentPage') }}
        </button>
      </template>
      <template v-if="!allResultsSelected && totalResults > selectedIds.length">
        <span v-if="selectedIds.length > 0" class="text-gray-300 dark:text-primary-800">•</span>
        <button
          :disabled="selectingAll"
          @click="$emit('select-all-results')"
          class="text-xs font-medium text-primary-700 hover:text-primary-800 disabled:cursor-not-allowed disabled:opacity-60 dark:text-primary-300 dark:hover:text-primary-200"
        >
          {{
            selectingAll
              ? t('admin.accounts.bulkActions.selectingAll')
              : t('admin.accounts.bulkActions.selectAllResults', { count: totalResults })
          }}
        </button>
      </template>
      <template v-if="selectedIds.length > 0">
        <span class="text-gray-300 dark:text-primary-800">•</span>
        <button
          @click="$emit('clear')"
          class="text-xs font-medium text-primary-700 hover:text-primary-800 dark:text-primary-300 dark:hover:text-primary-200"
        >
          {{ t('admin.accounts.bulkActions.clear') }}
        </button>
      </template>
    </div>
    <div class="flex flex-wrap justify-end gap-2">
      <div class="inline-flex h-8 items-stretch">
        <select
          :value="testScope"
          :disabled="testing || testingInvalid"
          class="min-w-[108px] rounded-l-md border border-r-0 border-gray-300 bg-white px-2 text-xs text-gray-700 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 disabled:cursor-not-allowed disabled:opacity-60 dark:border-gray-600 dark:bg-dark-700 dark:text-gray-200"
          @change="handleTestScopeChange"
        >
          <option value="all">{{ t('admin.accounts.bulkActions.testScopeAll') }}</option>
          <option value="error">{{ t('admin.accounts.bulkActions.testScopeError') }}</option>
        </select>
        <button
          @click="$emit('test-accounts')"
          :disabled="testing || testingInvalid"
          class="inline-flex items-center gap-1.5 rounded-r-md border border-primary-600 bg-primary-600 px-3 text-xs font-medium text-white transition-colors hover:bg-primary-700 disabled:cursor-not-allowed disabled:opacity-60"
        >
          <Icon name="refresh" size="sm" :class="testingInvalid ? 'animate-spin' : ''" />
          {{
            testingInvalid
              ? t('admin.accounts.bulkActions.testingInvalid', { completed: testProgress.completed, total: testProgress.total })
              : t('admin.accounts.bulkActions.startTest')
          }}
        </button>
      </div>
      <template v-if="selectedIds.length > 0">
        <button
          @click="$emit('test-selected')"
          :disabled="testing || testingInvalid"
          class="btn btn-secondary btn-sm disabled:cursor-not-allowed disabled:opacity-60"
        >
          {{ testing ? t('admin.accounts.bulkActions.testingSelected') : t('admin.accounts.bulkActions.testSelected') }}
        </button>
        <button @click="$emit('delete')" class="btn btn-danger btn-sm">{{ t('admin.accounts.bulkActions.delete') }}</button>
        <button @click="$emit('reset-status')" class="btn btn-secondary btn-sm">{{ t('admin.accounts.bulkActions.resetStatus') }}</button>
        <button @click="$emit('refresh-token')" class="btn btn-secondary btn-sm">{{ t('admin.accounts.bulkActions.refreshToken') }}</button>
        <button @click="$emit('probe-upstream-billing')" class="btn btn-secondary btn-sm">{{ t('admin.accounts.bulkActions.probeUpstreamBilling') }}</button>
        <button @click="$emit('toggle-keep-status-active', true)" class="btn btn-success btn-sm">{{ t('admin.accounts.bulkActions.enableKeepStatusActive') }}</button>
        <button @click="$emit('toggle-keep-status-active', false)" class="btn btn-secondary btn-sm">{{ t('admin.accounts.bulkActions.disableKeepStatusActive') }}</button>
        <button @click="$emit('toggle-schedulable', true)" class="btn btn-success btn-sm">{{ t('admin.accounts.bulkActions.enableScheduling') }}</button>
        <button @click="$emit('toggle-schedulable', false)" class="btn btn-warning btn-sm">{{ t('admin.accounts.bulkActions.disableScheduling') }}</button>
        <button @click="$emit('edit-selected')" class="btn btn-primary btn-sm">{{ t('admin.accounts.bulkActions.edit') }}</button>
      </template>
      <button @click="$emit('edit-filtered')" class="btn btn-primary btn-sm">
        {{ t('admin.accounts.bulkEdit.submit') }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'

type AccountTestScope = 'all' | 'error'

withDefaults(defineProps<{
  selectedIds: number[]
  totalResults: number
  selectingAll: boolean
  allResultsSelected: boolean
  testing?: boolean
  testingInvalid?: boolean
  testScope?: AccountTestScope
  testProgress?: {
    completed: number
    total: number
  }
}>(), {
  testing: false,
  testingInvalid: false,
  testScope: 'all',
  testProgress: () => ({ completed: 0, total: 0 })
})

const emit = defineEmits<{
  (e: 'delete'): void
  (e: 'edit-selected'): void
  (e: 'edit-filtered'): void
  (e: 'clear'): void
  (e: 'select-page'): void
  (e: 'select-all-results'): void
  (e: 'toggle-schedulable', enabled: boolean): void
  (e: 'toggle-keep-status-active', enabled: boolean): void
  (e: 'reset-status'): void
  (e: 'refresh-token'): void
  (e: 'test-selected'): void
  (e: 'test-accounts'): void
  (e: 'update:test-scope', scope: AccountTestScope): void
  (e: 'probe-upstream-billing'): void
}>()

const { t } = useI18n()

const handleTestScopeChange = (event: Event) => {
  emit('update:test-scope', (event.target as HTMLSelectElement).value as AccountTestScope)
}
</script>
