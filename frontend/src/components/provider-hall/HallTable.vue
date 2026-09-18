<template>
  <div class="hall-table">
    <!-- Desktop table (md and up) -->
    <div class="hidden overflow-x-auto md:block">
      <table
        class="w-full min-w-[1080px] table-fixed border-collapse text-sm"
        :aria-label="t('providerHall.aria.table')"
        data-testid="hall-table"
        @keydown="onKeydown"
      >
        <colgroup>
          <!-- At the 1080px minimum the last two columns must still hold the
               verification link and the "use this group" button. -->
          <col style="width: 16%" />
          <col style="width: 6%" />
          <col style="width: 11%" />
          <col style="width: 12%" />
          <col style="width: 11%" />
          <col style="width: 7%" />
          <col style="width: 7%" />
          <col style="width: 13%" />
          <col style="width: 8%" />
          <col style="width: 9%" />
        </colgroup>
        <thead class="sticky top-0 z-10 bg-white text-[11px] uppercase tracking-wide text-gray-500 dark:bg-dark-800 dark:text-dark-400">
          <tr class="border-b border-gray-200 dark:border-dark-700">
            <th scope="col" class="px-3 py-2 text-left font-semibold">{{ t('providerHall.columns.group') }}</th>
            <th scope="col" class="px-3 py-2 text-left font-semibold">{{ t('providerHall.columns.rate') }}</th>
            <th scope="col" class="px-3 py-2 text-left font-semibold">{{ t('providerHall.columns.price') }}</th>
            <th scope="col" class="px-3 py-2 text-left font-semibold">{{ t('providerHall.columns.status') }}</th>
            <th scope="col" class="px-3 py-2 text-left font-semibold">{{ t('providerHall.columns.ttft') }}</th>
            <th scope="col" class="px-3 py-2 text-center font-semibold">{{ t('providerHall.columns.cache') }}</th>
            <th scope="col" class="px-3 py-2 text-center font-semibold">{{ t('providerHall.columns.success') }}</th>
            <th scope="col" class="px-3 py-2 text-left font-semibold">{{ t('providerHall.columns.trend') }}</th>
            <th scope="col" class="px-3 py-2 text-center font-semibold">{{ t('providerHall.columns.verification') }}</th>
            <th scope="col" class="px-3 py-2 text-right font-semibold">{{ t('providerHall.columns.action') }}</th>
          </tr>
        </thead>
        <tbody>
          <template v-for="row in items" :key="row.group_id">
            <HallRow
              :row="row"
              :ttft-max="ttftMax"
              :expanded="expanded.has(row.group_id)"
              variant="row"
              @toggle="emit('toggle', $event)"
              @use-group="emit('use-group', $event)"
              @view-report="emit('view-report', $event)"
            />
            <tr v-if="expanded.has(row.group_id)" class="border-b border-gray-100 dark:border-dark-700" data-testid="hall-detail-row">
              <td :colspan="10" class="bg-gray-50/70 px-3 py-3 dark:bg-dark-900/40">
                <HallDetailPanel
                  :row="row"
                  :detail="details.get(`${row.group_id}:${range}`) ?? null"
                  :loading="detailLoading.has(row.group_id)"
                  :error="detailErrors.get(row.group_id) ?? null"
                  @use-group="emit('use-group', row)"
                  @view-report="emit('view-report', $event)"
                  @retry="emit('retry-detail', row.group_id)"
                />
              </td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>

    <!-- Mobile cards -->
    <div class="space-y-3 md:hidden" data-testid="hall-cards">
      <template v-for="row in items" :key="`card-${row.group_id}`">
        <HallRow
          :row="row"
          :ttft-max="ttftMax"
          :expanded="expanded.has(row.group_id)"
          variant="card"
          @toggle="emit('toggle', $event)"
          @use-group="emit('use-group', $event)"
          @view-report="emit('view-report', $event)"
        />
        <div v-if="expanded.has(row.group_id)" class="rounded-2xl border border-gray-200 bg-gray-50/70 p-3 dark:border-dark-700 dark:bg-dark-900/40">
          <HallDetailPanel
            :row="row"
            :detail="details.get(`${row.group_id}:${range}`) ?? null"
            :loading="detailLoading.has(row.group_id)"
            :error="detailErrors.get(row.group_id) ?? null"
            @use-group="emit('use-group', row)"
            @view-report="emit('view-report', $event)"
            @retry="emit('retry-detail', row.group_id)"
          />
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { ProviderHallDetailResponse, ProviderHallRange, ProviderHallRow as HallRowType } from '@/types/providerHall'
import { metricHasValue } from '@/utils/providerHallFormat'
import HallDetailPanel from './HallDetailPanel.vue'
import HallRow from './HallRow.vue'

const props = defineProps<{
  items: HallRowType[]
  expanded: Set<number>
  details: Map<string, ProviderHallDetailResponse>
  detailLoading: Set<number>
  detailErrors: Map<number, string>
  range: ProviderHallRange
}>()

const emit = defineEmits<{
  (e: 'toggle', groupId: number): void
  (e: 'use-group', row: HallRowType): void
  (e: 'view-report', payload: { groupId: number; reportId: number }): void
  (e: 'retry-detail', groupId: number): void
}>()

const { t } = useI18n()

/** Page-wide max for the relative TTFT bar. */
const ttftMax = computed(() => {
  let max = 0
  for (const row of props.items) {
    const metric = row.metrics.ttft_fast95_ms
    if (metricHasValue(metric) && metric.value > max) max = metric.value
  }
  return max
})

/** Arrow keys move focus between rows; Enter/Space are handled on the row. */
function onKeydown(event: KeyboardEvent) {
  if (event.key !== 'ArrowDown' && event.key !== 'ArrowUp') return
  const target = event.target as HTMLElement | null
  const current = target?.closest('tr[data-group-id]') as HTMLElement | null
  if (!current) return
  const rows = Array.from((event.currentTarget as HTMLElement).querySelectorAll<HTMLElement>('tr[data-group-id]'))
  const index = rows.indexOf(current)
  if (index < 0) return
  const next = rows[event.key === 'ArrowDown' ? index + 1 : index - 1]
  if (next) {
    event.preventDefault()
    next.focus()
  }
}
</script>
