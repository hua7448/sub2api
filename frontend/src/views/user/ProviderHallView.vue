<template>
  <AppLayout>
    <div class="space-y-4 pb-12">
      <section class="card !rounded-3xl !border-0 p-0 shadow-sm ring-1 ring-gray-900/5 dark:!bg-dark-800 dark:ring-dark-700">
        <header class="flex flex-wrap items-start justify-between gap-3 border-b border-gray-100 px-5 py-4 dark:border-dark-700 sm:px-6">
          <div class="min-w-0">
            <span class="badge badge-success inline-flex items-center gap-1 !text-[11px]" data-testid="hall-monitoring">
              <span class="h-1.5 w-1.5 rounded-full bg-emerald-500" aria-hidden="true"></span>
              {{ t('providerHall.monitoring') }}
            </span>
            <p class="mt-1.5 text-xs text-gray-500 dark:text-dark-400">{{ t('providerHall.intro') }}</p>
          </div>
          <div class="flex flex-wrap items-center gap-3 text-xs text-gray-500 dark:text-dark-400">
            <span v-if="hall.refreshing.value" class="inline-flex items-center gap-1 text-primary-600 dark:text-primary-300">
              <LoadingSpinner size="sm" />
              {{ t('providerHall.loading') }}
            </span>
            <span v-else-if="hall.dataThrough.value" class="font-mono" data-testid="hall-updated">
              {{ t('providerHall.updatedAt', { time: formatDateTime(hall.dataThrough.value, locale) }) }}
            </span>
            <button
              type="button"
              class="btn btn-ghost btn-sm !px-2"
              :title="t('providerHall.retry')"
              :aria-label="t('providerHall.retry')"
              data-testid="hall-refresh"
              @click="hall.fetch(true)"
            >
              <Icon name="refresh" size="sm" />
            </button>
            <AutoRefreshButton
              :enabled="hall.autoRefresh.enabled.value"
              :interval-seconds="hall.autoRefresh.intervalSeconds.value"
              :countdown="hall.autoRefresh.countdown.value"
              :intervals="hall.autoRefresh.intervals"
              @update:enabled="hall.autoRefresh.setEnabled"
              @update:interval="hall.autoRefresh.setInterval"
            />
          </div>
        </header>

        <div class="space-y-4 px-5 py-4 sm:px-6">
          <HallSummary :summary="hall.summary.value" />

          <HallToolbar
            :range="hall.range.value"
            :primary-sort="hall.primarySort.value"
            :is-custom-sort="hall.isCustomSort.value"
            :model="hall.model.value"
            :models="hall.catalog.value?.models ?? []"
            :search="hall.search.value"
            @update:range="(v) => (hall.range.value = v)"
            @cycle-sort="hall.cycleSort"
            @clear-sort="hall.clearSort"
            @open-custom="sortDialogOpen = true"
            @update:model="(v) => (hall.model.value = v)"
            @update:search="(v) => (hall.search.value = v)"
          />

          <div
            v-if="hall.stale.value"
            class="flex flex-wrap items-center justify-between gap-2 rounded-xl border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-800 dark:border-amber-800 dark:bg-amber-900/20 dark:text-amber-200"
            role="status"
            data-testid="hall-stale"
          >
            <span>{{ t('providerHall.staleBanner') }}<span v-if="hall.error.value" class="ml-1 opacity-80">({{ hall.error.value }})</span></span>
            <button type="button" class="btn btn-ghost btn-sm" @click="hall.fetch(false)">{{ t('providerHall.retry') }}</button>
          </div>

          <div v-if="hall.loading.value && !hall.data.value" class="flex items-center justify-center gap-2 py-16 text-sm text-gray-500 dark:text-dark-400" data-testid="hall-loading">
            <LoadingSpinner size="sm" />
            {{ t('providerHall.loading') }}
          </div>
          <div v-else-if="hall.error.value && !hall.data.value" class="py-12 text-center" data-testid="hall-error">
            <p class="text-sm text-red-600 dark:text-red-400">{{ hall.error.value }}</p>
            <button type="button" class="btn btn-secondary btn-sm mt-3" @click="hall.fetch(false)">{{ t('providerHall.retry') }}</button>
          </div>
          <EmptyState
            v-else-if="!hall.items.value.length"
            :title="t('providerHall.empty.title')"
            :description="t('providerHall.empty.description')"
          />
          <template v-else>
            <HallTable
              :items="hall.items.value"
              :expanded="hall.expanded.value"
              :details="hall.details.value"
              :detail-loading="hall.detailLoading.value"
              :detail-errors="hall.detailErrors.value"
              :range="hall.range.value"
              @toggle="hall.toggleExpand"
              @use-group="openUseGroup"
              @view-report="openReport"
              @retry-detail="(id) => hall.loadDetail(id, true)"
            />
            <Pagination
              v-if="hall.pagination.value.total > hall.pageSize"
              :total="hall.pagination.value.total"
              :page="hall.page.value"
              :page-size="hall.pageSize"
              :show-page-size-selector="false"
              @update:page="(p) => (hall.page.value = p)"
            />
          </template>
        </div>
      </section>
    </div>

    <HallSortDialog :show="sortDialogOpen" :rules="hall.sort.value" @close="sortDialogOpen = false" @apply="hall.setSort" />
    <UseGroupDialog :show="useGroupOpen" :group="useGroupRow" @close="useGroupOpen = false" />
    <HallVerificationDialog
      :show="reportOpen"
      :group-id="reportGroup?.group_id ?? null"
      :group-name="reportGroup?.name ?? ''"
      :profiles="reportProfiles"
      :initial-report-id="reportId"
      @close="reportOpen = false"
    />
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import AutoRefreshButton from '@/components/common/AutoRefreshButton.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import Pagination from '@/components/common/Pagination.vue'
import HallSortDialog from '@/components/provider-hall/HallSortDialog.vue'
import HallSummary from '@/components/provider-hall/HallSummary.vue'
import HallTable from '@/components/provider-hall/HallTable.vue'
import HallToolbar from '@/components/provider-hall/HallToolbar.vue'
import HallVerificationDialog from '@/components/provider-hall/HallVerificationDialog.vue'
import UseGroupDialog from '@/components/provider-hall/UseGroupDialog.vue'
import { useProviderHall } from '@/composables/useProviderHall'
import { formatDateTime } from '@/utils/providerHallFormat'
import type { ProviderHallProfileRef, ProviderHallRow } from '@/types/providerHall'

const { t, locale } = useI18n()

const hall = useProviderHall({
  loadFailedText: () => t('providerHall.loadFailed'),
  detailFailedText: () => t('providerHall.detailLoadFailed'),
})

const sortDialogOpen = ref(false)
const useGroupOpen = ref(false)
const useGroupRow = ref<ProviderHallRow | null>(null)
const reportOpen = ref(false)
const reportGroup = ref<ProviderHallRow | null>(null)
const reportId = ref<number | null>(null)

const reportProfiles = computed<ProviderHallProfileRef[]>(() => {
  if (!reportGroup.value) return []
  const detail = hall.getDetail(reportGroup.value.group_id)
  if (detail) return detail.profiles.map(({ profile_id, model, protocol }) => ({ profile_id, model, protocol }))
  return [reportGroup.value.default_profile]
})

function openUseGroup(row: ProviderHallRow) {
  useGroupRow.value = row
  useGroupOpen.value = true
}

function openReport(payload: { groupId: number; reportId: number }) {
  reportGroup.value = hall.items.value.find((row) => row.group_id === payload.groupId) ?? null
  reportId.value = payload.reportId
  reportOpen.value = true
}

onMounted(() => {
  void hall.fetch(false)
  hall.autoRefresh.start()
})
</script>
