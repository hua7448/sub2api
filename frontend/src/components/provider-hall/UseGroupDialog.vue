<template>
  <BaseDialog :show="show" :title="t('providerHall.useGroup.title', { name: group?.name ?? '' })" width="normal" @close="onClose">
    <div v-if="group" class="space-y-4" data-testid="use-group-dialog">
      <nav class="tabs w-full" role="tablist">
        <button
          type="button"
          role="tab"
          class="tab flex-1"
          :class="tab === 'switch' ? 'tab-active' : ''"
          :aria-selected="tab === 'switch'"
          data-testid="tab-switch"
          @click="tab = 'switch'"
        >{{ t('providerHall.useGroup.switchTab') }}</button>
        <button
          type="button"
          role="tab"
          class="tab flex-1"
          :class="tab === 'create' ? 'tab-active' : ''"
          :aria-selected="tab === 'create'"
          data-testid="tab-create"
          @click="tab = 'create'"
        >{{ t('providerHall.useGroup.createTab') }}</button>
      </nav>

      <!-- Switch an existing key -->
      <form v-if="tab === 'switch'" class="space-y-3" @submit.prevent="submitSwitch">
        <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('providerHall.useGroup.selectKey') }}</p>
        <div v-if="keysLoading && !keys.length" class="flex items-center gap-2 py-3 text-sm text-gray-500 dark:text-dark-400">
          <LoadingSpinner size="sm" />
          {{ t('providerHall.loading') }}
        </div>
        <p v-else-if="!keys.length" class="py-3 text-sm text-gray-400" data-testid="no-keys">{{ t('providerHall.useGroup.noKeys') }}</p>
        <ul v-else class="max-h-64 space-y-1 overflow-y-auto" data-testid="key-list">
          <li v-for="key in keys" :key="key.id">
            <label
              class="flex cursor-pointer items-center justify-between gap-3 rounded-lg border px-3 py-2 text-sm"
              :class="selectedKeyId === key.id ? 'border-primary-400 bg-primary-50 dark:border-primary-600 dark:bg-primary-900/20' : 'border-gray-200 dark:border-dark-700'"
            >
              <span class="flex min-w-0 items-center gap-2">
                <input v-model="selectedKeyId" type="radio" name="hall-key" :value="key.id" :disabled="key.group_id === group.group_id" />
                <span class="truncate font-medium text-gray-900 dark:text-white">{{ key.name }}</span>
              </span>
              <span class="shrink-0 text-xs text-gray-500 dark:text-dark-400">
                <template v-if="key.group_id === group.group_id">{{ t('providerHall.useGroup.alreadyInGroup') }}</template>
                <template v-else>{{ key.group?.name || t('providerHall.useGroup.noGroup') }} → {{ group.name }}</template>
              </span>
            </label>
          </li>
        </ul>
        <button v-if="keysHasMore" type="button" class="btn btn-ghost btn-sm" :disabled="keysLoading" @click="loadKeys(keysPage + 1)">
          {{ t('providerHall.useGroup.loadMore') }}
        </button>
        <div class="flex justify-end gap-2 pt-2">
          <button type="button" class="btn btn-secondary" @click="onClose">{{ t('providerHall.useGroup.cancel') }}</button>
          <button type="submit" class="btn btn-primary" :disabled="!selectedKeyId || submitting" data-testid="switch-submit">
            {{ t('providerHall.useGroup.switch') }}
          </button>
        </div>
      </form>

      <!-- Create a new key -->
      <form v-else class="space-y-3" @submit.prevent="submitCreate">
        <label class="block text-xs text-gray-500 dark:text-dark-400" for="hall-key-name">{{ t('providerHall.useGroup.keyName') }}</label>
        <input
          id="hall-key-name"
          v-model="newKeyName"
          type="text"
          class="input"
          maxlength="100"
          :placeholder="t('providerHall.useGroup.keyNamePlaceholder', { name: group.name })"
          data-testid="create-name"
        />
        <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('providerHall.useGroup.newGroup') }}: <span class="font-medium text-gray-900 dark:text-white">{{ group.name }}</span></p>
        <div class="flex justify-end gap-2 pt-2">
          <button type="button" class="btn btn-secondary" @click="onClose">{{ t('providerHall.useGroup.cancel') }}</button>
          <button type="submit" class="btn btn-primary" :disabled="!newKeyName.trim() || submitting" data-testid="create-submit">
            {{ t('providerHall.useGroup.create') }}
          </button>
        </div>
      </form>
    </div>
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import { keysAPI } from '@/api/keys'
import { useAppStore } from '@/stores/app'
import { extractApiErrorMessage } from '@/utils/apiError'
import type { ApiKey } from '@/types'
import type { ProviderHallRow } from '@/types/providerHall'

const props = defineProps<{ show: boolean; group: ProviderHallRow | null }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'done', payload: { action: 'switch' | 'create'; key: ApiKey }): void }>()
const { t } = useI18n()
const appStore = useAppStore()

const tab = ref<'switch' | 'create'>('switch')
const keys = ref<ApiKey[]>([])
const keysPage = ref(0)
const keysHasMore = ref(false)
const keysLoading = ref(false)
const selectedKeyId = ref<number | null>(null)
const newKeyName = ref('')
const submitting = ref(false)
const KEYS_PAGE_SIZE = 20

async function loadKeys(page: number) {
  keysLoading.value = true
  try {
    const response = await keysAPI.list(page, KEYS_PAGE_SIZE)
    keys.value = page === 1 ? response.items : [...keys.value, ...response.items]
    keysPage.value = page
    keysHasMore.value = page < (response.pages ?? 1)
  } catch (err) {
    appStore.showError(extractApiErrorMessage(err, t('providerHall.useGroup.loadKeysFailed')))
  } finally {
    keysLoading.value = false
  }
}

watch(
  () => [props.show, props.group?.group_id] as const,
  ([open]) => {
    if (open && props.group) {
      tab.value = 'switch'
      selectedKeyId.value = null
      newKeyName.value = ''
      keys.value = []
      void loadKeys(1)
    }
  },
  { immediate: true }
)

async function submitSwitch() {
  if (!props.group || !selectedKeyId.value || submitting.value) return
  submitting.value = true
  try {
    const key = await keysAPI.update(selectedKeyId.value, { group_id: props.group.group_id })
    appStore.showSuccess(t('providerHall.useGroup.switchSuccess', { name: props.group.name }))
    emit('done', { action: 'switch', key })
    emit('close')
  } catch (err) {
    // Keep the selection so the user can retry.
    appStore.showError(extractApiErrorMessage(err, t('providerHall.useGroup.switchFailed')))
  } finally {
    submitting.value = false
  }
}

async function submitCreate() {
  if (!props.group || !newKeyName.value.trim() || submitting.value) return
  submitting.value = true
  try {
    const key = await keysAPI.create(newKeyName.value.trim(), props.group.group_id)
    appStore.showSuccess(t('providerHall.useGroup.createSuccess', { name: props.group.name }))
    emit('done', { action: 'create', key })
    emit('close')
  } catch (err) {
    // Keep the typed name so the user can retry.
    appStore.showError(extractApiErrorMessage(err, t('providerHall.useGroup.createFailed')))
  } finally {
    submitting.value = false
  }
}

function onClose() {
  if (submitting.value) return
  emit('close')
}
</script>
