<template>
  <div class="space-y-4">
    <div class="flex flex-wrap items-center gap-3">
      <input v-model="filter.search" class="input w-full sm:w-64" :placeholder="t('admin.providerHall.searchGroups')" :aria-label="t('admin.providerHall.searchGroups')" />
      <select v-model="filter.platform" class="input w-auto" :aria-label="t('admin.providerHall.platform')"><option value="">{{ t('admin.providerHall.allPlatforms') }}</option><option v-for="p in platforms" :key="p">{{ p }}</option></select>
      <select v-model="filter.listed" class="input w-auto" :aria-label="t('admin.providerHall.listed')"><option value="">{{ t('admin.providerHall.allStatuses') }}</option><option value="true">{{ t('admin.providerHall.listed') }}</option><option value="false">{{ t('admin.providerHall.unlisted') }}</option></select>
      <select v-model="filter.sort" class="input w-auto" :aria-label="t('admin.providerHall.displayOrder')"><option value="name">{{ t('admin.providerHall.nameAsc') }}</option><option value="name_desc">{{ t('admin.providerHall.nameDesc') }}</option></select>
      <button class="btn btn-secondary btn-icon" :disabled="loading" :title="t('admin.providerHall.reload')" :aria-label="t('admin.providerHall.reload')" @click="load"><RefreshCw :size="16" /></button>
      <button class="btn btn-secondary" :disabled="!selected.length" @click="openBatch"><Layers :size="16" />{{ t('admin.providerHall.batch') }} ({{ selected.length }})</button>
    </div>
    <p v-if="error && !listing" class="hall-error" role="alert">{{ error }}</p>
    <DataTable :columns="columns" :data="groups" :loading="loading" row-key="group_id" selectable :selected-keys="selected" server-side-sort @update:selected-keys="selectKeys" @sort="(_key, order) => filter.sort = order === 'desc' ? 'name_desc' : 'name'">
      <template #cell-name="{ row: g }"><span class="break-words font-medium">{{ g.name }}</span><span v-if="g.display_name" class="block text-xs text-gray-500">{{ g.display_name }}</span><span v-if="!g.supported" class="block text-xs text-amber-600">{{ t('admin.providerHall.unsupported') }}</span></template>
      <template #cell-platform="{ row: g }">{{ g.platform }}<span v-if="g.status !== 'active'" class="block text-xs text-amber-600">{{ t('admin.providerHall.inactive') }}</span></template>
      <template #cell-listed="{ row: g }">{{ t(g.listed ? 'admin.providerHall.listed' : 'admin.providerHall.unlisted') }}</template>
      <template #cell-targets="{ row: g }"><span class="break-all">{{ g.targets.length }}</span><span class="block max-w-xs break-all text-xs text-gray-500">{{ g.models.join(', ') || t('admin.providerHall.noTargets') }}</span><span v-if="g.effective_profile_id" class="block text-xs">{{ t('admin.providerHall.defaultProfile') }}: {{ profileName(g.effective_profile_id) }}</span></template>
      <template #cell-keys="{ row: g }"><span v-for="(count, status) in g.key_statuses" :key="status" class="block text-xs">{{ t(`admin.providerHall.key_${status}`) }}: {{ count }}</span></template>
      <template #cell-results="{ row: g }"><button v-for="result in [g.latest_probe, g.latest_verification].filter(Boolean)" :key="result.job_id" class="block text-left text-xs text-emerald-600" @click="emit('open', result.job_id)">{{ t(`admin.providerHall.status_${result.status}`) }} {{ result.verdict ? t(`admin.providerHall.verdict_${result.verdict}`) : '' }}<span class="block text-gray-500">{{ new Date(result.at).toLocaleString() }}</span></button></template>
      <template #cell-actions="{ row: g }"><button class="btn btn-ghost btn-sm" :disabled="!g.supported" @click="open(g.group_id)"><Settings2 :size="15" />{{ t('admin.providerHall.manage') }}</button></template>
    </DataTable>
    <Pagination v-if="total" :total="total" :page="page" :page-size="pageSize" @update:page="page = $event; selected = []; load()" @update:page-size="pageSize = $event; page = 1; selected = []; load()" />
    <BaseDialog :show="listing !== null" :title="groups.find(g => g.group_id === listing?.group_id)?.name || t('admin.providerHall.manage')" width="extra-wide" :close-on-escape="!busy" :show-close-button="!busy" @close="close">
      <form v-if="listing" id="hall-settings" class="space-y-5" @submit.prevent="save">
        <p v-if="error" class="hall-error" role="alert">{{ error }}</p>
        <fieldset :disabled="busy" class="hall-fields">
          <label class="hall-field"><span>{{ t('admin.providerHall.displayName') }}</span><input v-model="listing.display_name" class="input" maxlength="100" /></label>
          <label class="hall-field"><span>{{ t('admin.providerHall.displayOrder') }}</span><input v-model.number="listing.display_order" class="input" type="number" step="1" /></label>
          <label class="hall-field sm:col-span-2"><span>{{ t('admin.providerHall.description') }}</span><textarea v-model="listing.description" class="input" maxlength="2000" rows="2" /></label>
          <label class="flex items-center gap-2 text-sm"><input v-model="listing.listed" type="checkbox" />{{ t('admin.providerHall.listed') }}</label>
        </fieldset>
        <div class="border-t border-gray-200 pt-4 dark:border-dark-700">
          <h3 class="mb-3 font-medium">{{ t('admin.providerHall.targets') }}</h3>
          <p v-if="!operatorId" class="text-sm text-amber-600">{{ t('admin.providerHall.operatorRequired') }}</p>
          <div class="mb-3 flex flex-wrap gap-2">
            <Select v-model="newProfile" class="min-w-0 flex-1 basis-52" :options="remainingProfiles" :disabled="busy" searchable :aria-label="t('admin.providerHall.profiles')" />
            <button type="button" class="btn btn-secondary" :disabled="busy || !newProfile" @click="addTarget"><Plus :size="16" />{{ t('admin.providerHall.addTarget') }}</button>
            <button type="button" class="btn btn-secondary" :disabled="busy" @click="candidateOpen = !candidateOpen">{{ t('admin.providerHall.chooseModels') }}</button>
            <button type="button" class="btn btn-secondary" :disabled="busy" @click="refreshModels"><RefreshCw :size="16" />{{ t('admin.providerHall.refreshUpstream') }}</button>
          </div>
          <div v-if="candidateOpen" class="mb-4 max-h-64 overflow-auto border-y border-gray-200 py-3 dark:border-dark-700">
            <label v-for="(m, index) in candidates" :key="`${m.model}:${m.protocol}:${m.source}:${index}`" class="flex items-start gap-3 py-2 text-sm"><input v-model="candidateSelection" type="checkbox" :value="index" :disabled="!m.available || busy" class="mt-1" /><span class="min-w-0 break-all">{{ m.model }} / {{ m.protocol }}<span class="block text-xs text-gray-500">{{ t(`admin.providerHall.source_${m.source}`) }} · {{ m.upstream_model }}<span v-if="!m.available"> · {{ t('admin.providerHall.invalidRoute') }}</span></span></span></label>
            <p v-if="!candidates.length" class="text-sm text-gray-500">{{ t('admin.providerHall.noModels') }}</p>
            <button type="button" class="btn btn-secondary mt-3" :disabled="busy || !candidateSelection.length" @click="importModels">{{ t('admin.providerHall.importSelected', { count: candidateSelection.length }) }}</button>
          </div>
          <p v-if="refreshResult" class="mb-3 text-sm" role="status">{{ refreshResult }}</p>
          <p class="mb-3 text-xs text-gray-500">{{ t('admin.providerHall.keyNotice') }}</p>
          <div class="overflow-x-auto"><table class="hall-table min-w-[1040px]">
            <thead><tr><th class="w-[20%]">{{ t('admin.providerHall.model') }}</th><th class="w-[25%]">{{ t('admin.providerHall.probeKey') }}</th><th class="w-[12%]">{{ t('admin.providerHall.probeMinutes') }}</th><th class="w-[12%]">{{ t('admin.providerHall.verifyHours') }}</th><th class="w-[8%]">{{ t('admin.providerHall.enabled') }}</th><th class="w-[8%]">{{ t('admin.providerHall.autoSchedule') }}</th><th class="w-[15%]">{{ t('admin.providerHall.jobActions') }}</th></tr></thead>
            <tbody><tr v-for="row in rows" :key="row.profile_id">
              <td class="break-all">{{ profileName(row.profile_id) }}<span class="block text-xs text-gray-500">{{ allProfiles.find(p => p.id === row.profile_id)?.protocol }}</span></td>
              <td><Select v-model="row.probe_key_id" :options="keyOptions(row.probe_key_id)" :disabled="busy" clearable searchable :aria-label="`${t('admin.providerHall.probeKey')} ${profileName(row.profile_id)}`" /><button type="button" class="btn btn-ghost btn-sm mt-1" :disabled="busy || !operatorId" @click="ensureKey(row)"><KeyRound :size="14" />{{ t('admin.providerHall.ensureKey') }}</button></td>
              <td><input :value="(row.probe_interval_seconds || 300) / 60" class="input" type="number" min="1" max="1440" step="any" :disabled="busy" :aria-label="t('admin.providerHall.probeMinutes')" @input="row.probe_interval_seconds = Math.round(Number(($event.target as HTMLInputElement).value) * 60)" /></td>
              <td><input :value="(row.verification_interval_seconds || 86400) / 3600" class="input" type="number" min="1" max="168" step="any" :disabled="busy" :aria-label="t('admin.providerHall.verifyHours')" @input="row.verification_interval_seconds = Math.round(Number(($event.target as HTMLInputElement).value) * 3600)" /></td>
              <td><input v-model="row.enabled" type="checkbox" :disabled="busy" :aria-label="`${t('admin.providerHall.enabled')} ${profileName(row.profile_id)}`" /></td>
              <td><input v-model="row.auto_schedule_enabled" type="checkbox" :disabled="busy" :aria-label="`${t('admin.providerHall.autoSchedule')} ${profileName(row.profile_id)}`" /></td>
              <td><button type="button" class="btn btn-ghost btn-sm whitespace-normal" :disabled="busy" @click="test(row, 'probe')">{{ t(!row.enabled ? 'admin.providerHall.enableAndTest' : dirty ? 'admin.providerHall.saveAndTest' : 'admin.providerHall.probeNow') }}</button><button type="button" class="btn btn-ghost btn-sm" :disabled="busy" @click="test(row, 'verification')">{{ t('admin.providerHall.verifyNow') }}</button><button type="button" class="btn btn-ghost btn-sm" :disabled="busy" @click="viewJobs(row.profile_id)">{{ t('admin.providerHall.results') }}</button></td>
            </tr></tbody>
          </table></div>
          <p v-if="!rows.length" class="py-6 text-sm text-gray-500">{{ t('admin.providerHall.noTargets') }}</p>
          <button v-if="lastJob" type="button" class="btn btn-secondary mt-3" @click="emit('open', lastJob)">{{ t('admin.providerHall.jobDetail') }} #{{ lastJob }}</button>
        </div>
      </form>
      <template #footer><div class="flex w-full flex-wrap justify-end gap-3"><button class="btn btn-secondary" :disabled="busy" @click="preflight">{{ t('admin.providerHall.preflight') }}</button><button class="btn btn-primary" type="submit" form="hall-settings" :disabled="busy"><Save :size="16" />{{ t('common.save') }}</button></div></template>
    </BaseDialog>
    <ProviderHallBatchDialog :show="batchOpen" :groups="batchSelection" :profiles="allProfiles" @close="batchOpen = false" @saved="selected = $event; load()" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { KeyRound, Layers, Plus, RefreshCw, Save, Settings2 } from 'lucide-vue-next'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Pagination from '@/components/common/Pagination.vue'
import DataTable from '@/components/common/DataTable.vue'
import Select from '@/components/common/Select.vue'
import ProviderHallBatchDialog from './ProviderHallBatchDialog.vue'
import * as hall from '@/api/admin/providerHall'
import { useAppStore } from '@/stores/app'
import { hallError, idempotencyKey } from './helpers'
import { extractApiErrorCode } from '@/utils/apiError'

const props = defineProps<{ profiles: hall.ProviderHallProfile[]; operatorId: number | null; focusGroup?: number | null }>()
const emit = defineEmits<{ dirty: [boolean]; open: [number]; jobs: [number, number]; profile: [hall.ProviderHallProfile] }>()
const { t } = useI18n()
const app = useAppStore()
const groups = ref<hall.ProviderHallGroupSummary[]>([])
const total = ref(0), page = ref(1), pageSize = ref(20), selected = ref<number[]>([])
const platforms = ['openai', 'composite', 'anthropic', 'gemini', 'antigravity', 'grok', 'deepseek', 'kimi', 'zhipu']
const filter = reactive({ search: '', platform: '', listed: '', sort: 'name' })
const listing = ref<hall.ProviderHallGroup | null>(null), saved = ref<hall.ProviderHallTargetSet | null>(null)
const rows = ref<hall.ProviderHallTargetInput[]>([]), localProfiles = ref<hall.ProviderHallProfile[]>([])
const allProfiles = computed(() => [...props.profiles, ...localProfiles.value.filter(p => !props.profiles.some(x => x.id === p.id))])
const keys = ref<hall.ProviderHallProbeKeyOption[]>([]), candidates = ref<hall.ProviderHallModelCandidate[]>([])
const candidateSelection = ref<number[]>([]), candidateOpen = ref(false), newProfile = ref<number | null>(null)
const loading = ref(false), busy = ref(false), error = ref(''), refreshResult = ref('')
const lastJob = ref<number | null>(null), baseline = ref(''), batchOpen = ref(false), batchSelection = ref<hall.ProviderHallGroupSummary[]>([])
const pendingKeys = new Map<string, string>()
const dirty = computed(() => listing.value !== null && snapshot() !== baseline.value)
const columns = computed(() => [{ key: 'name', label: t('admin.providerHall.jobGroup'), sortable: true }, { key: 'platform', label: t('admin.providerHall.platform') }, { key: 'listed', label: t('admin.providerHall.listed') }, { key: 'targets', label: t('admin.providerHall.targets') }, { key: 'keys', label: t('admin.providerHall.probeKey') }, { key: 'results', label: t('admin.providerHall.results') }, { key: 'actions', label: t('admin.providerHall.jobActions') }])
const remainingProfiles = computed(() => allProfiles.value.filter(p => !rows.value.some(r => r.profile_id === p.id)).map(p => ({ value: p.id, label: `${p.model} / ${p.protocol}` })))
let request: AbortController | undefined, groupRequest: AbortController | undefined, timer: ReturnType<typeof setTimeout> | undefined
let disposed = false
function profileName(id: number) { return allProfiles.value.find(p => p.id === id)?.model ?? `#${id}` }
function payload(): hall.ProviderHallSettingsInput { const g = listing.value!; return { version: g.version, listed: g.listed, display_name: g.display_name, description: g.description, display_order: g.display_order, items: rows.value } }
function snapshot() { return listing.value ? JSON.stringify(payload()) : '' }
function writable(set: hall.ProviderHallTargetSet) { return set.items.map(r => ({ profile_id: r.profile_id, probe_key_id: r.probe_key_id, enabled: r.enabled, auto_schedule_enabled: r.auto_schedule_enabled ?? false, probe_interval_seconds: r.probe_interval_seconds, verification_interval_seconds: r.verification_interval_seconds })) }
async function load() {
  request?.abort(); const current = new AbortController(); request = current; loading.value = true
  try { const result = await hall.listGroups({ ...filter, page: page.value, page_size: pageSize.value }, current.signal); if (!current.signal.aborted) { groups.value = result.items; total.value = result.total } }
  catch (err) { if (!current.signal.aborted) error.value = hallError(err, t, 'loadFailed') }
  finally { if (!current.signal.aborted) loading.value = false }
}
async function open(id: number) {
  if (busy.value || (dirty.value && !window.confirm(t('admin.providerHall.discard')))) return
  groupRequest?.abort(); const current = new AbortController(); groupRequest = current; error.value = ''; lastJob.value = null
  try {
    const [g, set, ks, ms] = await Promise.all([hall.getGroup(id, current.signal), hall.getTargets(id, current.signal), hall.listProbeKeys(id, current.signal), hall.listModels(id, current.signal)])
    if (current.signal.aborted) return
    if (g.version !== set.version) { error.value = t('admin.providerHall.conflict'); return }
    listing.value = g; saved.value = set; rows.value = writable(set); keys.value = ks; candidates.value = ms; candidateSelection.value = []; candidateOpen.value = false; baseline.value = snapshot(); refreshResult.value = ''
  } catch (err) { if (!current.signal.aborted) error.value = hallError(err, t, 'loadFailed') }
}
function viewJobs(profileID: number) { const groupID = listing.value?.group_id; if (!groupID) return; close(); if (!listing.value) emit('jobs', groupID, profileID) }
function close() { if (!busy.value && (!dirty.value || window.confirm(t('admin.providerHall.discard')))) { listing.value = null; error.value = '' } }
function selectKeys(keys: (number | string)[]) { selected.value = keys.map(Number).filter(id => groups.value.some(g => g.group_id === id && g.supported)) }
function openBatch() { batchSelection.value = JSON.parse(JSON.stringify(groups.value.filter(g => selected.value.includes(g.group_id)))); batchOpen.value = true }
function addTarget() { if (!newProfile.value || rows.value.length >= 100) return; rows.value.push({ profile_id: newProfile.value, probe_key_id: null, enabled: false, auto_schedule_enabled: false, probe_interval_seconds: 300, verification_interval_seconds: 86400 }); newProfile.value = null }
function keyOptions(current: number | null) {
  const opts = keys.value.map(k => ({ value: k.id as number | null, label: `${k.name} #${k.id} · ${t(`admin.providerHall.key_${k.status}`)}`, disabled: !['unused', 'registered'].includes(k.status) }))
  if (current && !opts.some(k => k.value === current)) opts.push({ value: current, label: t('admin.providerHall.unavailableKey', { id: current }), disabled: true })
  return [{ value: null, label: t('admin.providerHall.none'), disabled: false }, ...opts]
}
async function ensureKey(row: hall.ProviderHallTargetInput) {
  if (busy.value || !listing.value) return; busy.value = true; error.value = ''
  try { const key = await hall.ensureProbeKey(listing.value.group_id, row.profile_id); if (disposed) return; keys.value = [...keys.value.filter(k => k.id !== key.id), key]; row.probe_key_id = key.id }
  catch (err) { error.value = hallError(err, t) } finally { busy.value = false }
}
async function persist() {
  const result = await hall.saveSettings(listing.value!.group_id, payload())
  if (disposed) return false
  saved.value = result; if (result.listing) listing.value = result.listing; else listing.value!.version = result.version
  rows.value = writable(result); baseline.value = snapshot(); app.showSuccess(t('admin.providerHall.saved')); void load(); return true
}
async function save() { if (busy.value || !listing.value) return; busy.value = true; error.value = ''; try { await persist() } catch (err) { error.value = hallError(err, t) } finally { busy.value = false } }
async function preflight() { if (busy.value || !listing.value) return; busy.value = true; error.value = ''; try { await hall.preflightGroup(listing.value.group_id, payload()); app.showSuccess(t('admin.providerHall.preflightPassed')) } catch (err) { error.value = hallError(err, t) } finally { busy.value = false } }
async function test(row: hall.ProviderHallTargetInput, kind: hall.ProviderHallJobKind) {
  if (busy.value || !listing.value) return
  if (!row.enabled && !window.confirm(t('admin.providerHall.enableAndTestConfirm'))) return
  row.enabled = true; busy.value = true; error.value = ''
  try {
    if (dirty.value && !(await persist())) return
    const target = saved.value?.items.find(r => r.profile_id === row.profile_id), profile = allProfiles.value.find(p => p.id === row.profile_id)
    if (!target || !profile) return
    const slot = `${kind}:${target.id}:${target.version}:${profile.version}`, key = pendingKeys.get(slot) ?? idempotencyKey(slot); pendingKeys.set(slot, key)
    const result = await (kind === 'probe' ? hall.enqueueProbe : hall.enqueueVerification)(listing.value.group_id, row.profile_id, key, { target_version: target.version, profile_version: profile.version })
    lastJob.value = result.job_id; pendingKeys.delete(slot); app.showSuccess(t(result.reused ? 'admin.providerHall.enqueuedReused' : 'admin.providerHall.enqueued', { id: result.job_id }))
  } catch (err) { error.value = hallError(err, t) } finally { busy.value = false }
}
async function refreshModels() {
  if (busy.value || !listing.value) return; busy.value = true; error.value = ''
  try { const result = await hall.refreshModels(listing.value.group_id); if (disposed) return; candidates.value = [...candidates.value.filter(m => m.source !== 'upstream'), ...result.flatMap(r => r.models)]; candidateSelection.value = []; candidateOpen.value = true; refreshResult.value = result.map(r => `#${r.account_id}: ${t(r.success ? 'admin.providerHall.refreshOK' : 'admin.providerHall.refreshFailed')}`).join(' · ') }
  catch (err) { error.value = hallError(err, t) } finally { busy.value = false }
}
async function importModels() {
  if (busy.value) return
  const chosen = candidateSelection.value.map(i => candidates.value[i]).filter(Boolean)
  if (!window.confirm(`${t('admin.providerHall.importConfirm')}\n${chosen.map(m => `${m.model} / ${m.protocol}`).join('\n')}`)) return
  busy.value = true; error.value = ''
  try { for (const m of chosen) {
    let p = allProfiles.value.find(p => p.model === m.model && p.protocol === m.protocol)
    if (!p) {
      try { p = await hall.createProfile({ version: 0, model: m.model, protocol: m.protocol, supports_tools: false, output_limit: 256, model_aliases: [], reference_input_price: null, reference_cache_price: null, reference_cache_rate: null, reference_confirmed_at: null }) }
      catch (err) { if (extractApiErrorCode(err) !== 'PROVIDER_HALL_PROFILE_EXISTS') throw err; p = (await hall.listProfiles()).find(p => p.model === m.model && p.protocol === m.protocol); if (!p) throw err }
      localProfiles.value.push(p); emit('profile', p)
    }
    if (!rows.value.some(r => r.profile_id === p!.id)) { newProfile.value = p.id; addTarget() }
  } candidateSelection.value = [] }
  catch (err) { error.value = hallError(err, t) } finally { busy.value = false }
}
watch(dirty, v => emit('dirty', v))
watch(filter, () => { clearTimeout(timer); selected.value = []; page.value = 1; timer = setTimeout(() => void load(), 250) })
watch(() => props.focusGroup, id => { if (id) void open(id) })
watch(() => props.operatorId, () => { if (listing.value && !dirty.value && !busy.value) void open(listing.value.group_id) })
onMounted(() => { void load(); if (props.focusGroup) void open(props.focusGroup) })
onUnmounted(() => { disposed = true; request?.abort(); groupRequest?.abort(); clearTimeout(timer) })
</script>

<style>
.modal-content:has(#hall-settings) { @apply rounded-lg bg-white dark:bg-zinc-900; backdrop-filter: none; }
</style>
