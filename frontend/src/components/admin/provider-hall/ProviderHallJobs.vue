<template>
  <div class="space-y-5">
    <div class="flex flex-wrap items-end gap-3">
      <div class="hall-field w-44">
        <label for="hall-job-kind">{{ t('admin.providerHall.jobKind') }}</label>
        <Select id="hall-job-kind" v-model="filter.kind" :options="kindOptions" :disabled="loading" :aria-label="t('admin.providerHall.jobKind')" />
      </div>
      <div class="hall-field w-44">
        <label for="hall-job-status">{{ t('admin.providerHall.jobStatus') }}</label>
        <Select id="hall-job-status" v-model="filter.status" :options="statusOptions" :disabled="loading" :aria-label="t('admin.providerHall.jobStatus')" />
      </div>
      <label class="hall-field w-44"><span>{{ t('admin.providerHall.jobGroup') }}</span><input v-model="filter.group_name" class="input" :placeholder="t('admin.providerHall.searchGroups')" @input="filter.group_id = null; filter.profile_id = null" /></label>
      <label class="hall-field w-44"><span>{{ t('admin.providerHall.model') }}</span><input v-model="filter.model" class="input" :placeholder="t('admin.providerHall.searchModel')" /></label>
      <label class="hall-field w-44"><span>{{ t('admin.providerHall.jobSource') }}</span><select v-model="filter.source" class="input"><option value="">{{ t('admin.providerHall.allSources') }}</option><option value="manual">{{ t('admin.providerHall.source_manual') }}</option><option value="scheduled">{{ t('admin.providerHall.source_scheduled') }}</option></select></label>
      <button type="button" class="btn btn-secondary btn-icon" :disabled="loading" :title="t('admin.providerHall.reload')" :aria-label="t('admin.providerHall.reload')" @click="load"><RefreshCw :size="16" /></button>
    </div>
    <p v-if="error" role="alert" class="hall-error">{{ error }}</p>
    <p v-if="loading && !jobs.length" role="status" class="py-6 text-sm text-gray-500">{{ t('common.loading') }}</p>
    <div class="overflow-x-auto">
      <table class="hall-table min-w-[900px]">
        <thead><tr><th class="w-[8%]">ID</th><th class="w-[10%]">{{ t('admin.providerHall.jobKind') }}</th><th class="w-[12%]">{{ t('admin.providerHall.jobStatus') }}</th><th class="w-[12%]">{{ t('admin.providerHall.jobGroup') }} / {{ t('admin.providerHall.jobProfile') }}</th><th class="w-[16%]">{{ t('admin.providerHall.jobCreated') }}</th><th class="w-[16%]">{{ t('admin.providerHall.jobFinished') }}</th><th class="w-[14%]">{{ t('admin.providerHall.jobError') }}</th><th class="w-[12%]">{{ t('admin.providerHall.jobActions') }}</th></tr></thead>
        <tbody>
          <tr v-for="job in jobs" :key="job.id" :data-job-id="job.id">
            <td class="font-mono text-xs">#{{ job.id }}</td>
            <td>{{ t(`admin.providerHall.kind_${job.kind}`) }}<span class="block text-xs text-gray-500">{{ t(job.slot_at ? 'admin.providerHall.source_scheduled' : 'admin.providerHall.source_manual') }}</span></td>
            <td><span class="hall-pill" :class="`hall-pill-${job.status}`">{{ t(`admin.providerHall.status_${job.status}`) }}</span></td>
            <td class="break-all text-xs"><button class="text-emerald-600 underline" @click="emit('target', job.group_id)">{{ job.group_name || `#${job.group_id}` }}</button> / {{ job.config_snapshot?.profile?.model || `#${job.profile_id}` }}</td>
            <td class="text-xs">{{ formatDate(job.created_at) }}</td>
            <td class="text-xs">{{ job.finished_at ? formatDate(job.finished_at) : '—' }}</td>
            <td class="break-all text-xs">{{ job.error_code ? jobReason(job.error_code, t) : '—' }}</td>
            <td class="whitespace-nowrap">
              <button type="button" class="btn btn-ghost btn-sm" :aria-label="`${t('admin.providerHall.jobDetail')} #${job.id}`" @click="emit('open', job.id)">{{ t('common.view') }}</button>
              <button v-if="job.status === 'queued' || job.status === 'running' || job.status === 'unknown'" type="button" class="btn btn-ghost btn-sm text-red-600 dark:text-red-400" :disabled="cancelling === job.id" :aria-label="`${t('admin.providerHall.cancelJob')} #${job.id}`" @click="askCancel(job.id)">{{ t('common.cancel') }}</button>
            </td>
          </tr>
          <tr v-if="!jobs.length && !loading"><td colspan="8" class="py-10 text-center text-gray-500">{{ t('admin.providerHall.noJobs') }}</td></tr>
        </tbody>
      </table>
    </div>
    <Pagination v-if="total > 0" :total="total" :page="page" :page-size="pageSize" :page-size-options="[20, 50, 100]" @update:page="setPage" @update:page-size="setPageSize" />
    <ConfirmDialog :show="confirmID !== null" :title="t('admin.providerHall.cancelJob')" :message="t('admin.providerHall.cancelJobConfirm', { id: confirmID ?? 0 })" danger @confirm="confirmCancel" @cancel="confirmID = null" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { RefreshCw } from 'lucide-vue-next'
import Select from '@/components/common/Select.vue'
import Pagination from '@/components/common/Pagination.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import * as hall from '@/api/admin/providerHall'
import { formatDate } from '@/utils/format'
import { useAppStore } from '@/stores/app'
import { hallError, jobReason } from './helpers'

const props = defineProps<{ focusGroup?: number | null; focusProfile?: number | null; active?: boolean }>()
const emit = defineEmits<{ open: [number]; target: [number] }>()
const { t } = useI18n()
const app = useAppStore()
const jobs = ref<hall.ProviderHallJob[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)
const error = ref('')
const cancelling = ref<number | null>(null)
const confirmID = ref<number | null>(null)
const filter = reactive<hall.ProviderHallJobFilter>({ kind: '', status: '', group_id: props.focusGroup, profile_id: props.focusProfile, group_name: '', model: '', source: '' })
let refreshTimer: ReturnType<typeof setInterval> | undefined, filterTimer: ReturnType<typeof setTimeout> | undefined
let request: AbortController | undefined
let disposed = false

const kindOptions = computed(() => [{ value: '', label: t('admin.providerHall.allKinds') }, ...(['probe', 'verification'] as const).map(k => ({ value: k, label: t(`admin.providerHall.kind_${k}`) }))])
const statusOptions = computed(() => [{ value: '', label: t('admin.providerHall.allStatuses') }, ...(['queued', 'running', 'succeeded', 'failed', 'cancelled', 'unknown'] as const).map(s => ({ value: s, label: t(`admin.providerHall.status_${s}`) }))])

async function load() {
  request?.abort()
  const current = new AbortController()
  request = current
  loading.value = true
  error.value = ''
  try {
    const result = await hall.listJobs({ ...filter, group_id: filter.group_id || null, page: page.value, page_size: pageSize.value }, current.signal)
    if (current.signal.aborted) return
    jobs.value = result.items
    total.value = result.total
  } catch (err) { if (!current.signal.aborted) error.value = hallError(err, t, 'loadFailed') }
  finally { if (!current.signal.aborted) loading.value = false }
}
function setPage(value: number) { page.value = value; void load() }
function setPageSize(value: number) { pageSize.value = value; page.value = 1; void load() }
function askCancel(id: number) { confirmID.value = id }
async function confirmCancel() {
  const id = confirmID.value
  confirmID.value = null
  if (id === null || cancelling.value !== null) return
  cancelling.value = id
  error.value = ''
  try {
    await hall.cancelJob(id)
    if (disposed) return
    app.showSuccess(t('admin.providerHall.cancelled'))
    await load()
  } catch (err) { if (!disposed) error.value = hallError(err, t) }
  finally { cancelling.value = null }
}
watch(filter, () => { page.value = 1; clearTimeout(filterTimer); filterTimer = setTimeout(() => void load(), 250) })
watch(() => [props.focusGroup, props.focusProfile], () => { filter.group_id = props.focusGroup; filter.profile_id = props.focusProfile; filter.group_name = ''; filter.model = '' })
defineExpose({ load })
onMounted(() => { void load(); refreshTimer = setInterval(() => { if (props.active !== false && !loading.value && jobs.value.some(j => ['queued', 'running', 'unknown'].includes(j.status))) void load() }, 5000) })
onUnmounted(() => { disposed = true; request?.abort(); clearInterval(refreshTimer); clearTimeout(filterTimer) })
</script>

<style>
.provider-hall .hall-pill { @apply inline-block rounded-full px-2 py-0.5 text-xs font-medium; }
.provider-hall .hall-pill-queued { @apply bg-gray-100 text-gray-700 dark:bg-dark-700 dark:text-gray-300; }
.provider-hall .hall-pill-running { @apply bg-blue-50 text-blue-700 dark:bg-blue-950/40 dark:text-blue-300; }
.provider-hall .hall-pill-succeeded, .provider-hall .hall-pill-passed, .provider-hall .hall-pill-received { @apply bg-emerald-50 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-300; }
.provider-hall .hall-pill-failed, .provider-hall .hall-pill-suspected, .provider-hall .hall-pill-model_mismatch { @apply bg-red-50 text-red-700 dark:bg-red-950/40 dark:text-red-300; }
.provider-hall .hall-pill-cancelled, .provider-hall .hall-pill-prepared { @apply bg-gray-100 text-gray-500 dark:bg-dark-700 dark:text-gray-400; }
.provider-hall .hall-pill-unknown, .provider-hall .hall-pill-uncertain, .provider-hall .hall-pill-insufficient, .provider-hall .hall-pill-error, .provider-hall .hall-pill-dispatched { @apply bg-amber-50 text-amber-700 dark:bg-amber-950/40 dark:text-amber-300; }
</style>
