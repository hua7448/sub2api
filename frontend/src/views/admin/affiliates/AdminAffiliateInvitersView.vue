<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex flex-wrap items-center gap-3">
          <div class="relative w-full md:w-80">
            <Icon name="search" size="md" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
            <input v-model="search" type="text" class="input pl-10" :placeholder="t('common.search')" @input="debounceLoad" />
          </div>
          <button class="btn btn-secondary px-2 md:px-3" :disabled="loading" :title="t('common.refresh')" @click="loadData">
            <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
          </button>
        </div>
      </template>

      <template #table>
        <div class="overflow-x-auto">
          <table class="w-full min-w-[800px] text-left text-sm">
            <thead>
              <tr class="border-b border-gray-200 text-gray-500 dark:border-dark-700 dark:text-dark-400">
                <th class="px-3 py-2 font-medium">ID</th>
                <th class="px-3 py-2 font-medium">{{ t('common.email') }}</th>
                <th class="px-3 py-2 font-medium">{{ t('common.username') }}</th>
                <th class="px-3 py-2 font-medium">{{ t('affiliate.stats.invitedUsers') }}</th>
                <th class="px-3 py-2 font-medium">{{ t('affiliate.threshold.qualifiedCount', { count: '', required: '' }).split(':')[0] }}</th>
                <th class="px-3 py-2 font-medium text-right">{{ t('affiliate.stats.availableQuota') }}</th>
                <th class="px-3 py-2 font-medium text-right">{{ t('affiliate.stats.totalQuota') }}</th>
                <th class="px-3 py-2 font-medium">{{ t('common.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="loading" class="border-b border-gray-100 dark:border-dark-800">
                <td colspan="8" class="px-3 py-8 text-center text-gray-400">{{ t('common.loading') }}</td>
              </tr>
              <tr v-else-if="items.length === 0" class="border-b border-gray-100 dark:border-dark-800">
                <td colspan="8" class="px-3 py-8 text-center text-gray-400">{{ t('common.noData') }}</td>
              </tr>
              <tr v-for="item in items" :key="item.user_id" class="border-b border-gray-100 last:border-b-0 dark:border-dark-800">
                <td class="px-3 py-3 font-mono text-xs text-gray-500">{{ item.user_id }}</td>
                <td class="px-3 py-3 text-gray-900 dark:text-white">{{ item.email || '-' }}</td>
                <td class="px-3 py-3 text-gray-700 dark:text-gray-300">{{ item.username || '-' }}</td>
                <td class="px-3 py-3">{{ item.aff_count }}</td>
                <td class="px-3 py-3">
                  <span class="font-medium" :class="item.qualified_count > 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-500'">{{ item.qualified_count }}</span>
                </td>
                <td class="px-3 py-3 text-right font-medium text-emerald-600 dark:text-emerald-400">{{ formatCurrency(item.aff_quota) }}</td>
                <td class="px-3 py-3 text-right text-gray-700 dark:text-gray-300">{{ formatCurrency(item.aff_history_quota) }}</td>
                <td class="px-3 py-3">
                  <button class="btn btn-secondary btn-sm" @click="openAdjust(item)">
                    {{ t('common.edit') }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </template>

      <template #pagination>
        <Pagination
          :page="page"
          :page-size="pageSize"
          :total="total"
          @update:page="handlePageChange"
          @update:pageSize="handlePageSizeChange"
        />
      </template>
    </TablePageLayout>

    <BaseDialog :show="adjustDialog" :title="t('common.edit')" @close="adjustDialog = false">
      <div v-if="adjustTarget" class="space-y-4">
        <p class="text-sm text-gray-600 dark:text-gray-400">
          {{ adjustTarget.email }} (ID: {{ adjustTarget.user_id }})
        </p>
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('affiliate.stats.availableQuota') }}</label>
          <input v-model.number="adjustAmount" type="number" step="0.01" min="0" class="input" />
        </div>
        <div class="flex justify-end gap-2">
          <button class="btn btn-secondary" @click="adjustDialog = false">{{ t('common.cancel') }}</button>
          <button class="btn btn-primary" :disabled="adjusting" @click="doAdjust">{{ t('common.save') }}</button>
        </div>
      </div>
    </BaseDialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Pagination from '@/components/common/Pagination.vue'
import { affiliatesAPI, type AffiliateInviterEntry } from '@/api/admin/affiliates'
import { formatCurrency } from '@/utils/format'
import { useAppStore } from '@/stores/app'
import { useDebounceFn } from '@vueuse/core'

const { t } = useI18n()
const appStore = useAppStore()

const items = ref<AffiliateInviterEntry[]>([])
const loading = ref(false)
const search = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

const adjustDialog = ref(false)
const adjustTarget = ref<AffiliateInviterEntry | null>(null)
const adjustAmount = ref(0)
const adjusting = ref(false)

async function loadData() {
  loading.value = true
  try {
    const resp = await affiliatesAPI.listInviters({ page: page.value, page_size: pageSize.value, search: search.value })
    items.value = resp.items ?? []
    total.value = resp.total ?? 0
  } catch (e: any) {
    appStore.showError(e?.message || 'Failed to load')
  } finally {
    loading.value = false
  }
}

const debounceLoad = useDebounceFn(() => { page.value = 1; loadData() }, 300)

function handlePageChange(p: number) {
  page.value = p
  loadData()
}

function handlePageSizeChange(size: number) {
  pageSize.value = size
  page.value = 1
  loadData()
}

function openAdjust(item: AffiliateInviterEntry) {
  adjustTarget.value = item
  adjustAmount.value = item.aff_quota
  adjustDialog.value = true
}

async function doAdjust() {
  if (!adjustTarget.value) return
  adjusting.value = true
  try {
    await affiliatesAPI.adjustUserQuota(adjustTarget.value.user_id, 'set', adjustAmount.value)
    appStore.showSuccess(t('common.saved'))
    adjustDialog.value = false
    loadData()
  } catch (e: any) {
    appStore.showError(e?.message || 'Failed')
  } finally {
    adjusting.value = false
  }
}

onMounted(loadData)
</script>
