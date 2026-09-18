<template>
  <BaseDialog :show="show" :title="t('providerHall.verification.dialogTitle')" width="wide" @close="emit('close')">
    <div class="space-y-3" data-testid="verification-dialog">
      <div class="flex flex-wrap items-center justify-between gap-2">
        <div class="text-sm font-semibold text-gray-900 dark:text-white">{{ groupName }}</div>
        <div v-if="profiles.length > 1" class="w-64">
          <Select
            :model-value="profileId ?? ''"
            :options="profileOptions"
            :aria-label="t('providerHall.verification.profile')"
            @update:model-value="onProfileChange"
          />
        </div>
      </div>

      <div v-if="loading && !items.length" class="flex items-center gap-2 py-6 text-sm text-gray-500 dark:text-dark-400">
        <LoadingSpinner size="sm" />
        {{ t('providerHall.loading') }}
      </div>
      <div v-else-if="error && !items.length" class="py-4 text-sm text-red-600 dark:text-red-400">
        {{ error }}
        <button type="button" class="btn btn-ghost btn-sm ml-2" @click="load()">{{ t('providerHall.retry') }}</button>
      </div>

      <!-- Single report -->
      <article v-else-if="selected" class="space-y-3" data-testid="verification-report">
        <button type="button" class="btn btn-ghost btn-sm" @click="selected = null">← {{ t('providerHall.verification.back') }}</button>
        <header class="flex flex-wrap items-center gap-2">
          <h4 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('providerHall.verification.reportTitle', { id: selected.report_id }) }}</h4>
          <span class="badge" :class="verdictClass(selected.verdict, selected.expired)">{{ verdictLabel(selected.verdict, selected.expired) }}</span>
          <span class="text-xs text-gray-500 dark:text-dark-400">{{ selected.model }} · {{ selected.protocol }}</span>
        </header>
        <dl class="grid grid-cols-2 gap-2 text-xs sm:grid-cols-4">
          <div><dt class="text-gray-500 dark:text-dark-400">{{ t('providerHall.verification.executionStatus') }}</dt><dd class="font-medium text-gray-900 dark:text-white">{{ executionLabel(selected.execution_status) }}</dd></div>
          <div><dt class="text-gray-500 dark:text-dark-400">{{ t('providerHall.verification.reasonCode') }}</dt><dd class="font-mono text-gray-900 dark:text-white">{{ selected.reason_code || '-' }}</dd></div>
          <div><dt class="text-gray-500 dark:text-dark-400">{{ t('providerHall.verification.completedAt') }}</dt><dd class="text-gray-900 dark:text-white">{{ formatDateTime(selected.completed_at, locale) }}</dd></div>
          <div><dt class="text-gray-500 dark:text-dark-400">{{ t('providerHall.verification.expiresAt') }}</dt><dd class="text-gray-900 dark:text-white">{{ formatDateTime(selected.expires_at, locale) }}</dd></div>
        </dl>
        <section>
          <h5 class="text-xs font-semibold text-gray-600 dark:text-dark-300">{{ t('providerHall.verification.suites') }}</h5>
          <table class="mt-1 w-full text-xs">
            <tbody>
              <tr v-for="suite in suites" :key="suite.key" class="border-t border-gray-100 dark:border-dark-700">
                <th scope="row" class="py-1.5 pr-3 text-left font-medium text-gray-700 dark:text-dark-200">{{ t(`providerHall.verification.suite.${suite.key}`) }}</th>
                <td class="py-1.5 text-gray-900 dark:text-white">
                  <template v-if="suite.counts">
                    <span class="mr-3 text-emerald-700 dark:text-emerald-300">{{ t('providerHall.verification.passedCount', { n: suite.counts.passed }) }}</span>
                    <span class="mr-3 text-red-700 dark:text-red-300">{{ t('providerHall.verification.failedCount', { n: suite.counts.failed }) }}</span>
                    <span class="text-amber-700 dark:text-amber-300">{{ t('providerHall.verification.errorCount', { n: suite.counts.error }) }}</span>
                  </template>
                  <template v-else-if="suite.model">
                    <span class="mr-3 text-emerald-700 dark:text-emerald-300">{{ t('providerHall.verification.matched', { n: suite.model.matched }) }}</span>
                    <span class="mr-3 text-red-700 dark:text-red-300">{{ t('providerHall.verification.mismatched', { n: suite.model.mismatched }) }}</span>
                    <span class="mr-3 text-amber-700 dark:text-amber-300">{{ t('providerHall.verification.missing', { n: suite.model.missing }) }}</span>
                    <div v-if="suite.model.seen?.length" class="mt-0.5 font-mono text-[11px] text-gray-500 dark:text-dark-400">
                      {{ t('providerHall.verification.seen', { models: suite.model.seen.join(', ') }) }}
                    </div>
                  </template>
                  <span v-else class="text-gray-400">{{ t('providerHall.verification.notRun') }}</span>
                </td>
              </tr>
            </tbody>
          </table>
        </section>
      </article>

      <!-- Report list -->
      <div v-else>
        <p v-if="!items.length" class="py-6 text-center text-sm text-gray-400" data-testid="verification-empty">{{ t('providerHall.verification.empty') }}</p>
        <ul v-else class="divide-y divide-gray-100 dark:divide-dark-700" data-testid="verification-list">
          <li v-for="report in items" :key="report.report_id">
            <button type="button" class="flex w-full items-center justify-between gap-3 py-2 text-left hover:bg-gray-50 dark:hover:bg-dark-700/40" @click="selected = report">
              <span class="flex min-w-0 items-center gap-2">
                <span class="badge shrink-0" :class="verdictClass(report.verdict, report.expired)">{{ verdictLabel(report.verdict, report.expired) }}</span>
                <span class="truncate text-xs text-gray-700 dark:text-dark-200">{{ report.model }} · {{ report.protocol }}</span>
              </span>
              <span class="shrink-0 text-xs text-gray-500 dark:text-dark-400">{{ formatDateTime(report.completed_at, locale) }}</span>
            </button>
          </li>
        </ul>
        <Pagination
          v-if="total > pageSize"
          :total="total"
          :page="page"
          :page-size="pageSize"
          :show-page-size-selector="false"
          @update:page="onPage"
        />
      </div>
    </div>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import Pagination from '@/components/common/Pagination.vue'
import Select from '@/components/common/Select.vue'
import * as api from '@/api/providerHall'
import { extractApiErrorMessage } from '@/utils/apiError'
import { formatDateTime } from '@/utils/providerHallFormat'
import type { ProviderHallProfileRef, ProviderHallVerificationReport, ProviderHallVerificationSummary } from '@/types/providerHall'

const props = defineProps<{
  show: boolean
  groupId: number | null
  groupName: string
  profiles: ProviderHallProfileRef[]
  initialProfileId?: number | null
  initialReportId?: number | null
}>()
const emit = defineEmits<{ (e: 'close'): void }>()
const { t, locale } = useI18n()

const pageSize = 20
const page = ref(1)
const total = ref(0)
const items = ref<ProviderHallVerificationReport[]>([])
const selected = ref<ProviderHallVerificationReport | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)
const profileId = ref<number | null>(props.initialProfileId ?? null)
let controller: AbortController | null = null

const profileOptions = computed(() => [
  { value: '', label: t('providerHall.verification.profile') },
  ...props.profiles.map((p) => ({ value: p.profile_id, label: `${p.model} · ${p.protocol}` })),
])

async function load() {
  if (!props.show || !props.groupId) return
  controller?.abort()
  const request = new AbortController()
  controller = request
  loading.value = true
  error.value = null
  try {
    const response = await api.listVerifications(
      props.groupId,
      { profile_id: profileId.value ?? undefined, page: page.value, page_size: pageSize },
      request.signal
    )
    if (controller !== request) return
    items.value = response.items
    total.value = response.pagination.total
    if (props.initialReportId && !selected.value) {
      selected.value = response.items.find((r) => r.report_id === props.initialReportId) ?? null
    }
  } catch (err) {
    const e = err as { name?: string; code?: string }
    if (e?.name === 'AbortError' || e?.name === 'CanceledError' || e?.code === 'ERR_CANCELED') return
    error.value = extractApiErrorMessage(err, t('providerHall.reportLoadFailed'))
  } finally {
    if (controller === request) loading.value = false
  }
}

watch(
  () => [props.show, props.groupId] as const,
  ([open]) => {
    if (open) {
      page.value = 1
      selected.value = null
      profileId.value = props.initialProfileId ?? null
      void load()
    } else {
      controller?.abort()
    }
  },
  { immediate: true }
)

function onProfileChange(value: string | number | boolean | null) {
  profileId.value = value ? Number(value) : null
  page.value = 1
  selected.value = null
  void load()
}

function onPage(next: number) {
  page.value = next
  void load()
}

type SuiteRow = {
  key: 'arithmetic' | 'json' | 'tool' | 'model'
  counts?: { passed: number; failed: number; error: number } | null
  model?: { matched: number; mismatched: number; missing: number; seen?: string[] } | null
}

const suites = computed<SuiteRow[]>(() => {
  const summary = (selected.value?.summary ?? {}) as ProviderHallVerificationSummary
  return [
    { key: 'arithmetic', counts: summary.arithmetic ?? null },
    { key: 'json', counts: summary.json ?? null },
    { key: 'tool', counts: summary.tool ?? null },
    { key: 'model', model: summary.model ?? null },
  ]
})

function verdictLabel(verdict: string, expired: boolean): string {
  if (expired) return t('providerHall.verification.expired')
  return t(`providerHall.verification.${verdict}`)
}
function verdictClass(verdict: string, expired: boolean): string {
  if (expired) return 'badge-gray'
  if (verdict === 'passed') return 'badge-success'
  if (verdict === 'failed') return 'badge-danger'
  return 'badge-warning'
}
function executionLabel(status: string): string {
  const key = `providerHall.verification.execution.${status}`
  return ['completed', 'partial', 'error'].includes(status) ? t(key) : status
}

onBeforeUnmount(() => controller?.abort())
</script>
