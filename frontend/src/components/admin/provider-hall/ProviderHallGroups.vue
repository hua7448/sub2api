<template>
  <div class="space-y-5">
    <div class="flex max-w-4xl items-end gap-3">
      <div class="hall-field flex-1">
        <label for="hall-group">{{ t('admin.providerHall.chooseGroup') }}</label>
        <Select id="hall-group" :model-value="selectedID" :options="groupOptions" searchable :disabled="saving"
          :aria-label="t('admin.providerHall.chooseGroup')" @update:model-value="selectGroup" />
      </div>
      <button type="button" class="btn btn-secondary btn-icon" :disabled="loading || saving" :title="t('admin.providerHall.reload')" :aria-label="t('admin.providerHall.reload')" @click="reload"><RefreshCw :size="16" /></button>
    </div>
    <p v-if="error" role="alert" class="hall-error">{{ error }}</p>
    <p v-if="loading" role="status" class="py-6 text-sm text-gray-500">{{ t('common.loading') }}</p>
    <template v-if="listing && targets && !loading">
      <form class="hall-form" @submit.prevent="saveListing">
        <fieldset :disabled="saving" class="hall-fields">
          <label class="hall-field"><span>{{ t('admin.providerHall.displayName') }}</span><input v-model="listing.display_name" class="input" maxlength="100" /></label>
          <label class="hall-field"><span>{{ t('admin.providerHall.displayOrder') }}</span><input v-model.number="listing.display_order" class="input" type="number" min="-1000000" max="1000000" step="1" required /></label>
          <label class="hall-field sm:col-span-2"><span>{{ t('admin.providerHall.description') }}</span><textarea v-model="listing.description" class="input" maxlength="2000" rows="2" /></label>
          <label class="flex items-center gap-2 text-sm"><input v-model="listing.listed" type="checkbox" class="h-4 w-4" />{{ t('admin.providerHall.listed') }}</label>
        </fieldset>
        <div class="hall-actions"><span class="text-xs text-gray-500">{{ t('admin.providerHall.version', { version: listing.version }) }}</span><button type="submit" class="btn btn-primary" :disabled="saving"><Save :size="16" />{{ t('common.save') }}</button></div>
      </form>
      <form class="space-y-4 border-t border-gray-200 pt-6 dark:border-dark-700" @submit.prevent="saveTargets">
        <h2 class="text-base font-semibold">{{ t('admin.providerHall.targets') }}</h2>
        <p class="max-w-4xl text-sm text-gray-500 dark:text-gray-400">{{ t('admin.providerHall.keyNotice') }}</p>
        <p v-if="!operatorId" class="text-sm text-amber-700 dark:text-amber-400">{{ t('admin.providerHall.operatorRequired') }}</p>
        <p v-if="keyError" role="alert" class="hall-error">{{ keyError }}</p>
        <div class="flex max-w-4xl flex-wrap items-center gap-3">
          <Select v-model="newProfileID" class="min-w-0 flex-1 basis-52" :options="remainingProfiles" :disabled="saving || rows.length >= 100" :aria-label="t('admin.providerHall.profiles')" />
          <button type="button" class="btn btn-secondary" :disabled="saving || !newProfileID || rows.length >= 100" @click="addTarget"><Plus :size="16" />{{ t('admin.providerHall.addTarget') }}</button>
          <button v-if="hasMoreKeys" type="button" class="btn btn-secondary" :disabled="keysLoading" @click="loadKeys(false)">{{ t('admin.providerHall.moreKeys') }}</button>
        </div>
        <fieldset :disabled="saving" class="min-w-0">
          <div class="overflow-x-auto">
            <table class="hall-table min-w-[1040px]">
              <thead><tr><th class="w-[24%]">{{ t('admin.providerHall.profiles') }}</th><th class="w-[26%]">{{ t('admin.providerHall.probeKey') }}</th><th class="w-[17%]">{{ t('admin.providerHall.probeInterval') }}</th><th class="w-[19%]">{{ t('admin.providerHall.verificationInterval') }}</th><th class="w-[8%]">{{ t('admin.providerHall.enabled') }}</th><th class="w-[16%]">{{ t('admin.providerHall.jobActions') }}</th></tr></thead>
              <tbody>
                <tr v-for="row in rows" :key="row.profile_id">
                  <td><span class="block break-all font-medium">{{ profileName(row.profile_id) }}</span><span class="mt-1 block text-xs text-gray-500">{{ profiles.find(p => p.id === row.profile_id)?.protocol }}</span></td>
                  <td><Select v-model="row.probe_key_id" :options="keyOptions(row.probe_key_id)" :disabled="saving || keysLoading" :aria-label="`${t('admin.providerHall.probeKey')} ${profileName(row.profile_id)}`" searchable clearable /></td>
                  <td><input v-model.number="row.probe_interval_seconds" type="number" min="60" max="86400" step="1" required class="input" :aria-label="`${t('admin.providerHall.probeInterval')} ${profileName(row.profile_id)}`" /></td>
                  <td><input v-model.number="row.verification_interval_seconds" type="number" min="3600" max="604800" step="1" required class="input" :aria-label="`${t('admin.providerHall.verificationInterval')} ${profileName(row.profile_id)}`" /></td>
                  <td><input v-model="row.enabled" type="checkbox" class="mt-3 h-4 w-4" :aria-label="`${t('admin.providerHall.enabled')} ${profileName(row.profile_id)}`" /></td>
                  <td class="whitespace-nowrap">
                    <template v-if="targets.items.some(item => item.profile_id === row.profile_id)">
                      <button type="button" class="btn btn-ghost btn-sm" :disabled="enqueueing !== null || !targets.items.find(item => item.profile_id === row.profile_id)?.enabled" :aria-label="`${t('admin.providerHall.probeNow')} ${profileName(row.profile_id)}`" @click="enqueue('probe', row.profile_id)">{{ t('admin.providerHall.probeNow') }}</button>
                      <button type="button" class="btn btn-ghost btn-sm" :disabled="enqueueing !== null || !targets.items.find(item => item.profile_id === row.profile_id)?.enabled" :aria-label="`${t('admin.providerHall.verifyNow')} ${profileName(row.profile_id)}`" @click="enqueue('verification', row.profile_id)">{{ t('admin.providerHall.verifyNow') }}</button>
                    </template>
                    <button v-else type="button" class="btn btn-ghost btn-icon" :title="t('common.delete')" :aria-label="`${t('common.delete')} ${profileName(row.profile_id)}`" @click="rows = rows.filter(r => r !== row)"><Trash2 :size="16" /></button>
                  </td>
                </tr>
                <tr v-if="!rows.length"><td colspan="6" class="py-10 text-center text-gray-500">{{ t('admin.providerHall.noTargets') }}</td></tr>
              </tbody>
            </table>
          </div>
        </fieldset>
        <div class="hall-actions"><button type="submit" class="btn btn-primary" :disabled="saving"><Save :size="16" />{{ t('admin.providerHall.saveTargets') }}</button></div>
      </form>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { Plus, RefreshCw, Save, Trash2 } from 'lucide-vue-next'
import Select from '@/components/common/Select.vue'
import * as hall from '@/api/admin/providerHall'
import * as groupsAPI from '@/api/admin/groups'
import * as usersAPI from '@/api/admin/users'
import type { AdminGroup } from '@/types'
import { useAppStore } from '@/stores/app'
import { extractApiErrorCode } from '@/utils/apiError'
import { hallError, idempotencyKey } from './helpers'

const props = defineProps<{ profiles: hall.ProviderHallProfile[]; operatorId: number | null }>()
const emit = defineEmits<{ dirty: [boolean] }>()
const { t } = useI18n()
const app = useAppStore()
const groups = ref<AdminGroup[]>([])
const selectedID = ref<number | null>(null)
const listing = ref<hall.ProviderHallGroup | null>(null)
const targets = ref<hall.ProviderHallTargetSet | null>(null)
const rows = ref<hall.ProviderHallTargetInput[]>([])
const newProfileID = ref<number | null>(null)
const error = ref('')
const keyError = ref('')
const loading = ref(false)
const saving = ref(false)
const keysLoading = ref(false)
const enqueueing = ref<string | null>(null)
// One idempotency key per click, reused while that click's request is in
// flight so a double click can never create two jobs.
const pendingKeys = new Map<string, string>()
const hasMoreKeys = ref(false)
const keys = ref<{ id: number; name: string; group_id: number | null; status: string }[]>([])
const listingBaseline = ref('')
const targetsBaseline = ref('[]')
const dirty = computed(() => (listing.value && listingSnapshot() !== listingBaseline.value) || JSON.stringify(rows.value) !== targetsBaseline.value)
watch(dirty, value => emit('dirty', Boolean(value)))
let request: AbortController | undefined
let keyRequest: AbortController | undefined
let keyPage = 0
let disposed = false

const groupOptions = computed(() => groups.value.map(g => ({ value: g.id, label: `${g.name} / ${g.platform} (#${g.id})` })))
const remainingProfiles = computed(() => props.profiles.filter(p => !rows.value.some(r => r.profile_id === p.id)).map(p => ({ value: p.id, label: `${p.model} / ${p.protocol}` })))
function profileName(id: number) { return props.profiles.find(p => p.id === id)?.model ?? `#${id}` }
function listingSnapshot() {
  const g = listing.value
  return JSON.stringify(g && { listed: g.listed, display_name: g.display_name, description: g.description, display_order: g.display_order })
}
function writableRows(set: hall.ProviderHallTargetSet) {
  return set.items.map(({ profile_id, probe_key_id, enabled, probe_interval_seconds, verification_interval_seconds }) => ({ profile_id, probe_key_id, enabled, probe_interval_seconds, verification_interval_seconds }))
}
function keyOptions(current: number | null) {
  const options = keys.value.filter(k => k.group_id === selectedID.value).map(k => ({ value: k.id, label: `${k.name} (#${k.id})`, disabled: k.status !== 'active' }))
  if (current && !options.some(o => o.value === current)) options.unshift({ value: current, label: t('admin.providerHall.unavailableKey', { id: current }), disabled: true })
  return [{ value: null, label: t('admin.providerHall.none'), disabled: false }, ...options]
}
function selectGroup(value: unknown) {
  if (typeof value !== 'number' || value === selectedID.value || saving.value) return
  if (dirty.value && !window.confirm(t('admin.providerHall.discard'))) return
  selectedID.value = value
  void loadGroup()
}
async function loadGroup() {
  request?.abort()
  const id = selectedID.value
  if (!id) return
  const current = new AbortController()
  request = current
  loading.value = true
  error.value = ''
  listing.value = null
  targets.value = null
  rows.value = []
  targetsBaseline.value = '[]'
  try {
    const [group, set] = await Promise.all([hall.getGroup(id, current.signal), hall.getTargets(id, current.signal)])
    if (current.signal.aborted) return
    if (group.version !== set.version) { error.value = t('admin.providerHall.conflict'); return }
    listing.value = group
    targets.value = set
    rows.value = writableRows(set)
    listingBaseline.value = listingSnapshot()
    targetsBaseline.value = JSON.stringify(rows.value)
  } catch (err) { if (!current.signal.aborted) error.value = hallError(err, t, 'loadFailed') }
  finally { if (!current.signal.aborted) loading.value = false }
}
async function loadGroups() {
  try {
    const result = await groupsAPI.getAllIncludingInactive()
    if (disposed) return
    groups.value = result.filter(g => g.platform === 'openai' || g.platform === 'composite')
    if (!selectedID.value && groups.value.length) selectedID.value = groups.value[0].id
    await loadGroup()
  } catch (err) { if (!disposed) error.value = hallError(err, t, 'loadFailed') }
}
function reload() {
  if (!dirty.value || window.confirm(t('admin.providerHall.discard'))) { void loadGroups(); void loadKeys(true) }
}
async function loadKeys(reset: boolean) {
  if (reset) { keyRequest?.abort(); keys.value = []; keyPage = 0; hasMoreKeys.value = false; keysLoading.value = false }
  const id = props.operatorId
  if (!id || keysLoading.value) return
  const current = new AbortController()
  keyRequest = current
  keysLoading.value = true
  keyError.value = ''
  try {
    const page = keyPage + 1
    const result = await usersAPI.getUserApiKeys(id, { page, page_size: 100, signal: current.signal })
    if (current.signal.aborted) return
    keys.value.push(...result.items.map(k => ({ id: k.id, name: k.name, group_id: k.group_id, status: k.status })))
    keyPage = page
    hasMoreKeys.value = page * 100 < result.total
  } catch (err) { if (!current.signal.aborted) keyError.value = hallError(err, t, 'loadFailed') }
  finally { if (!current.signal.aborted) keysLoading.value = false }
}
function addTarget() {
  if (!newProfileID.value || rows.value.some(r => r.profile_id === newProfileID.value)) return
  rows.value.push({ profile_id: newProfileID.value, probe_key_id: null, enabled: false, probe_interval_seconds: 300, verification_interval_seconds: 86400 })
  newProfileID.value = null
}
async function enqueue(kind: hall.ProviderHallJobKind, profileID: number) {
  const groupID = listing.value?.group_id
  if (!groupID) return
  const slot = `${kind}:${groupID}:${profileID}`
  let key = pendingKeys.get(slot)
  if (!key) { key = idempotencyKey(`admin-${kind}-${groupID}-${profileID}`); pendingKeys.set(slot, key) }
  if (enqueueing.value !== null) return
  enqueueing.value = slot
  error.value = ''
  try {
    const result = await (kind === 'probe' ? hall.enqueueProbe : hall.enqueueVerification)(groupID, profileID, key)
    if (disposed) return
    app.showSuccess(t(result.reused ? 'admin.providerHall.enqueuedReused' : 'admin.providerHall.enqueued', { id: result.job_id }))
    pendingKeys.delete(slot)
  } catch (err) {
    if (disposed) return
    error.value = hallError(err, t)
    // Business rejections are final for this click; transport failures keep the key for a retry.
    if (['PROVIDER_HALL_TASKS_DISABLED', 'PROVIDER_HALL_BUDGET_EXHAUSTED', 'PROVIDER_HALL_TARGET_DISABLED', 'PROVIDER_HALL_NOT_READY'].includes(extractApiErrorCode(err) ?? '')) pendingKeys.delete(slot)
  } finally { enqueueing.value = null }
}
async function saveListing() {
  if (!listing.value || !targets.value || saving.value) return
  saving.value = true
  error.value = ''
  const g = listing.value
  try {
    const result = await hall.updateGroup(g.group_id, { version: g.version, listed: g.listed, display_name: g.display_name, description: g.description, display_order: g.display_order })
    if (disposed) return
    listing.value = result
    targets.value.version = result.version
    listingBaseline.value = listingSnapshot()
    app.showSuccess(t('admin.providerHall.groupSaved'))
  } catch (err) { if (!disposed) error.value = hallError(err, t) }
  finally { saving.value = false }
}
async function saveTargets() {
  if (!listing.value || !targets.value || saving.value) return
  saving.value = true
  error.value = ''
  try {
    const result = await hall.updateTargets(listing.value.group_id, { version: targets.value.version, items: rows.value })
    if (disposed) return
    targets.value = result
    listing.value.version = result.version
    rows.value = writableRows(result)
    targetsBaseline.value = JSON.stringify(rows.value)
    app.showSuccess(t('admin.providerHall.saved'))
  } catch (err) { if (!disposed) error.value = hallError(err, t) }
  finally { saving.value = false }
}
watch(() => props.operatorId, () => {
  void loadKeys(true)
  // Keep unsaved target edits, but let their old version fail CAS after an operator change.
  if (!dirty.value && !saving.value) void loadGroup()
})
onMounted(() => { void loadGroups(); void loadKeys(true) })
onUnmounted(() => { disposed = true; request?.abort(); keyRequest?.abort() })
</script>
