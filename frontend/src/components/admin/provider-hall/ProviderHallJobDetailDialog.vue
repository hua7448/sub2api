<template>
  <BaseDialog :show="jobId !== null" :title="`${t('admin.providerHall.jobDetail')} #${jobId ?? ''}`" width="wide" @close="emit('close')">
    <p v-if="error" role="alert" class="hall-error">{{ error }}</p>
    <p v-else-if="loading" role="status" class="py-6 text-sm text-gray-500">{{ t('common.loading') }}</p>
    <div v-else-if="detail" class="space-y-6 text-sm">
      <dl class="grid gap-x-6 gap-y-2 sm:grid-cols-2">
        <div><dt class="hall-dt">{{ t('admin.providerHall.jobKind') }}</dt><dd>{{ t(`admin.providerHall.kind_${detail.job.kind}`) }}</dd></div>
        <div><dt class="hall-dt">{{ t('admin.providerHall.jobStatus') }}</dt><dd><span class="hall-pill" :class="`hall-pill-${detail.job.status}`">{{ t(`admin.providerHall.status_${detail.job.status}`) }}</span></dd></div>
        <div><dt class="hall-dt">{{ t('admin.providerHall.jobGroup') }} / {{ t('admin.providerHall.jobProfile') }}</dt><dd>#{{ detail.job.group_id }} / {{ detail.job.config_snapshot?.profile?.model }} ({{ detail.job.config_snapshot?.profile?.protocol }})</dd></div>
        <div><dt class="hall-dt">{{ t('admin.providerHall.jobAttempts') }}</dt><dd>{{ detail.job.attempts }}</dd></div>
        <div><dt class="hall-dt">{{ t('admin.providerHall.jobCreated') }}</dt><dd>{{ formatDate(detail.job.created_at) }}</dd></div>
        <div><dt class="hall-dt">{{ t('admin.providerHall.jobFinished') }}</dt><dd>{{ detail.job.finished_at ? formatDate(detail.job.finished_at) : '—' }}</dd></div>
        <div v-if="detail.job.error_code" class="sm:col-span-2"><dt class="hall-dt">{{ t('admin.providerHall.jobError') }}</dt><dd class="break-all">{{ detail.job.error_code }}<span v-if="detail.job.error_message" class="text-gray-500"> — {{ detail.job.error_message }}</span></dd></div>
      </dl>
      <section v-if="detail.verification" class="space-y-2 border-t border-gray-200 pt-4 dark:border-dark-700">
        <h3 class="font-semibold">{{ t('admin.providerHall.report') }}</h3>
        <dl class="grid gap-x-6 gap-y-2 sm:grid-cols-2">
          <div><dt class="hall-dt">{{ t('admin.providerHall.verdict') }}</dt><dd><span class="hall-pill" :class="`hall-pill-${detail.verification.verdict}`">{{ t(`admin.providerHall.verdict_${detail.verification.verdict}`) }}</span><span v-if="detail.verification.stale" class="ml-2 text-xs text-amber-700 dark:text-amber-400">{{ t('admin.providerHall.stale') }}</span></dd></div>
          <div><dt class="hall-dt">{{ t('admin.providerHall.executionStatus') }}</dt><dd>{{ detail.verification.execution_status }}<span v-if="detail.verification.reason_code" class="text-gray-500"> · {{ detail.verification.reason_code }}</span></dd></div>
          <div v-for="suite in suites" :key="suite.key"><dt class="hall-dt">{{ suite.label }}</dt><dd>{{ suite.text }}</dd></div>
          <div><dt class="hall-dt">{{ t('admin.providerHall.expiresAt') }}</dt><dd>{{ formatDate(detail.verification.expires_at) }}</dd></div>
        </dl>
      </section>
      <section class="space-y-2 border-t border-gray-200 pt-4 dark:border-dark-700">
        <h3 class="font-semibold">{{ t('admin.providerHall.samples') }}</h3>
        <div class="overflow-x-auto">
          <table class="hall-table min-w-[760px]">
            <thead><tr><th class="w-[12%]">{{ t('admin.providerHall.sampleTest') }}</th><th class="w-[12%]">{{ t('admin.providerHall.sampleStatus') }}</th><th class="w-[14%]">{{ t('admin.providerHall.sampleResult') }}</th><th class="w-[16%]">{{ t('admin.providerHall.sampleLatency') }}</th><th class="w-[14%]">{{ t('admin.providerHall.sampleTokens') }}</th><th class="w-[18%]">{{ t('admin.providerHall.sampleModel') }}</th><th class="w-[14%]">{{ t('admin.providerHall.sampleBilling') }}</th></tr></thead>
            <tbody>
              <tr v-for="sample in detail.samples" :key="sample.id">
                <td>{{ sample.test_id }} #{{ sample.seq }}</td>
                <td><span class="hall-pill" :class="`hall-pill-${sample.status}`">{{ t(`admin.providerHall.sample_${sample.status}`) }}</span></td>
                <td class="break-all text-xs">{{ sample.result ?? '—' }}<span v-if="sample.error_code" class="text-red-600 dark:text-red-400"> {{ sample.error_code }}</span></td>
                <td class="text-xs">{{ ms(sample.ttft_ms) }} / {{ ms(sample.total_ms) }}</td>
                <td class="text-xs">{{ sample.input_tokens ?? '—' }} / {{ sample.output_tokens ?? '—' }}</td>
                <td class="break-all text-xs">{{ sample.response_model || '—' }}</td>
                <td class="text-xs">{{ sample.billing_status }}<span v-if="sample.actual_cost !== null"> · {{ sample.actual_cost }}</span></td>
              </tr>
              <tr v-if="!detail.samples.length"><td colspan="7" class="py-6 text-center text-gray-500">{{ t('admin.providerHall.noSamples') }}</td></tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>
    <template #footer>
      <div class="flex justify-end"><button type="button" class="btn btn-secondary" @click="emit('close')">{{ t('common.close') }}</button></div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import * as hall from '@/api/admin/providerHall'
import { formatDate } from '@/utils/format'
import { hallError } from './helpers'

const props = defineProps<{ jobId: number | null }>()
const emit = defineEmits<{ close: [] }>()
const { t } = useI18n()
const detail = ref<hall.ProviderHallJobDetail | null>(null)
const loading = ref(false)
const error = ref('')
let request: AbortController | undefined

function ms(value: number | null) { return value === null ? '—' : `${value} ms` }
const suites = computed(() => {
  const summary = detail.value?.verification?.summary
  if (!summary) return []
  const out: { key: string; label: string; text: string }[] = []
  for (const [key, label] of [['arithmetic', 'suiteArithmetic'], ['json', 'suiteJson'], ['tool', 'suiteTool']] as const) {
    const suite = summary[key]
    if (suite) out.push({ key, label: t(`admin.providerHall.${label}`), text: t('admin.providerHall.suiteResult', suite) })
  }
  if (summary.model) out.push({ key: 'model', label: t('admin.providerHall.suiteModel'), text: t('admin.providerHall.modelResult', summary.model) })
  return out
})
async function load(id: number) {
  request?.abort()
  const current = new AbortController()
  request = current
  loading.value = true
  error.value = ''
  detail.value = null
  try {
    const result = await hall.getJob(id, current.signal)
    if (!current.signal.aborted) detail.value = result
  } catch (err) { if (!current.signal.aborted) error.value = hallError(err, t, 'loadFailed') }
  finally { if (!current.signal.aborted) loading.value = false }
}
watch(() => props.jobId, id => { if (id !== null) void load(id); else { request?.abort(); detail.value = null } }, { immediate: true })
onUnmounted(() => request?.abort())
</script>

<style>
.hall-dt { @apply text-xs text-gray-500 dark:text-gray-400; }
</style>
