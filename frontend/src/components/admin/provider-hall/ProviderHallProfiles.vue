<template>
  <TablePageLayout>
    <template #filters>
      <div class="flex flex-col justify-between gap-4 lg:flex-row lg:items-center">
        <div class="flex flex-1 flex-wrap items-center gap-3">
          <div class="relative w-full sm:w-64">
            <Search class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400 dark:text-gray-500" />
            <input v-model="search" class="input pl-9" :placeholder="t('admin.providerHall.searchProfiles')" :aria-label="t('admin.providerHall.searchProfiles')" />
          </div>
          <select v-model="groupFilter" class="input w-auto" :aria-label="t('admin.providerHall.linkedGroups')">
            <option value="">{{ t('admin.providerHall.allGroups') }}</option>
            <option v-for="g in groupOptions" :key="g.value" :value="g.value">{{ g.label }}</option>
          </select>
          <select v-model="referenceFilter" class="input w-auto" :aria-label="t('admin.providerHall.reference')">
            <option value="">{{ t('admin.providerHall.allReferences') }}</option>
            <option value="valid">{{ t('admin.providerHall.referenceValid') }}</option>
            <option value="invalid">{{ t('admin.providerHall.referenceInvalid') }}</option>
          </select>
        </div>
        <div class="flex w-full flex-wrap items-center justify-end gap-3 lg:w-auto">
          <button class="btn btn-secondary btn-icon" :disabled="loading" :title="t('admin.providerHall.reload')" :aria-label="t('admin.providerHall.reload')" @click="emit('reload')">
            <RefreshCw :size="16" :class="loading && 'animate-spin'" />
          </button>
          <button class="btn btn-secondary" :disabled="candidatesLoading" @click="openCandidates">
            <Sparkles :size="16" />{{ t('admin.providerHall.batchProfiles') }}
          </button>
          <button class="btn btn-primary" @click="open()"><Plus :size="16" />{{ t('admin.providerHall.addProfile') }}</button>
        </div>
      </div>
    </template>

    <template #table>
      <p v-if="error && !props.profiles.length" class="hall-error m-4" role="alert">{{ error }}</p>
      <DataTable :columns="columns" :data="paged" :loading="loading" row-key="id">
        <template #cell-model="{ row: p }">
          <div class="min-w-0">
            <span class="font-medium text-txt-primary">{{ p.model }}</span>
            <span v-if="isDefault(p)" class="ml-2 badge badge-success">{{ t('admin.providerHall.defaultMarker') }}</span>
            <span class="block truncate text-xs text-txt-dim">{{ p.model_aliases.length ? p.model_aliases.join(', ') : t('admin.providerHall.noAliases') }}</span>
          </div>
        </template>
        <template #cell-protocol="{ row: p }">{{ protocolLabel(p.protocol) }}</template>
        <template #cell-groups="{ row: p }">
          <span v-if="p.groups?.length" class="text-xs">{{ groupNames(p) }}</span>
          <span v-else class="text-xs text-txt-dim">{{ t('admin.providerHall.unusedProfile') }}</span>
        </template>
        <template #cell-candidates="{ row: p }">
          <span v-if="groupCount(p)" class="text-xs">{{ t('admin.providerHall.availableInGroups', { count: groupCount(p) }) }}</span>
          <span v-else class="text-xs text-txt-dim">—</span>
        </template>
        <template #cell-output_limit="{ row: p }">{{ p.output_limit }}</template>
        <template #cell-tools="{ row: p }">
          <Check v-if="p.supports_tools" :size="17" class="text-emerald-600" :aria-label="t('common.yes')" />
          <span v-else class="text-txt-dim">{{ t('common.no') }}</span>
        </template>
        <template #cell-reference="{ row: p }">
          <span class="text-xs" :class="referenceValid(p) ? 'text-emerald-600' : 'text-amber-600'">{{ t(referenceValid(p) ? 'admin.providerHall.referenceValid' : 'admin.providerHall.referenceInvalid') }}</span>
          <span class="block text-xs text-txt-dim">{{ p.reference_input_price ?? '-' }} / {{ p.reference_cache_price ?? '-' }}</span>
        </template>
        <template #cell-actions="{ row: p }">
          <button class="btn btn-ghost btn-icon" :title="t('admin.providerHall.editProfile')" :aria-label="`${t('common.edit')} ${p.model}`" @click="open(p)"><Pencil :size="16" /></button>
          <button class="btn btn-ghost btn-icon text-red-600 dark:text-red-400" :title="t('common.delete')" :aria-label="`${t('common.delete')} ${p.model}`" @click="askDelete(p)"><Trash2 :size="16" /></button>
        </template>
      </DataTable>
    </template>

    <template #pagination>
      <Pagination v-if="filtered.length" :total="filtered.length" :page="page" :page-size="pageSize" :page-size-options="[20, 50, 100]" @update:page="page = $event" @update:page-size="setPageSize" />
    </template>
  </TablePageLayout>

  <BaseDialog :show="show" :title="t(editingID ? 'admin.providerHall.editProfile' : 'admin.providerHall.addProfile')" width="wide" :close-on-escape="!saving" :show-close-button="!saving" @close="close">
    <form id="hall-profile-form" class="hall-form" @submit.prevent="save">
      <p v-if="error" role="alert" class="hall-error">{{ error }}</p>
      <p v-if="editingProfile?.groups?.length" class="text-sm">{{ t('admin.providerHall.sharedProfile') }} {{ editingProfile.groups.map(g => g.name).join(', ') }}</p>
      <p v-if="editingProfile && isDefault(editingProfile)" class="text-sm text-amber-600">{{ t('admin.providerHall.defaultIdentityLocked') }} <button type="button" class="underline" @click="show = false; emit('config')">{{ t('admin.providerHall.settingsEntry') }}</button></p>

      <div v-if="!editingID" class="space-y-3">
        <label class="hall-field">
          <span>{{ t('admin.providerHall.chooseModel') }}</span>
          <Select :model-value="typedModel ? null : pickValue" :options="pickOptions" :disabled="!candidateOptions.length" searchable clearable :aria-label="t('admin.providerHall.chooseModel')" @update:model-value="pickModel" />
          <span class="text-xs text-txt-dim">{{ candidateOptions.length ? t('admin.providerHall.chooseModelHint', { count: candidateOptions.length }) : t('admin.providerHall.chooseModelEmpty') }}</span>
        </label>
      </div>

      <fieldset :disabled="saving" class="hall-fields">
        <label class="hall-field">
          <span>{{ t('admin.providerHall.model') }}</span>
          <input v-model.trim="draft.model" class="input" maxlength="200" required :readonly="!typedModel && !!pickValue" :class="{ 'opacity-70': !typedModel && !!pickValue }" />
        </label>
        <label class="hall-field">
          <span>{{ t('admin.providerHall.protocol') }}</span>
          <select v-model="draft.protocol" class="input" :disabled="!typedModel && !!pickValue">
            <option v-for="p in protocols" :key="p.value" :value="p.value">{{ p.label }}</option>
          </select>
        </label>
        <label class="flex items-center gap-2 text-sm sm:col-span-2">
          <input v-model="typedModel" type="checkbox" class="h-4 w-4" />{{ t('admin.providerHall.manualModel') }}
          <span class="text-xs text-txt-dim">{{ t('admin.providerHall.manualModelHint') }}</span>
        </label>
        <label class="hall-field"><span>{{ t('admin.providerHall.outputLimit') }}</span><input v-model.number="draft.output_limit" class="input" type="number" min="1" max="1024" step="1" required /></label>
        <label class="flex items-center gap-2 text-sm"><input v-model="draft.supports_tools" type="checkbox" class="h-4 w-4" />{{ t('admin.providerHall.tools') }}</label>
        <label class="hall-field sm:col-span-2"><span>{{ t('admin.providerHall.aliases') }}</span><textarea v-model="aliasText" class="input font-mono" rows="3" /></label>
        <label class="flex items-center gap-2 text-sm sm:col-span-2"><input v-model="hasReference" type="checkbox" class="h-4 w-4" />{{ t('admin.providerHall.reference') }}</label>
        <template v-if="hasReference">
          <label class="hall-field"><span>{{ t('admin.providerHall.inputPrice') }}</span><input v-model="draft.reference_input_price" class="input" inputmode="decimal" pattern="[0-9]+(\.[0-9]{1,10})?" required @input="referenceEdited" /></label>
          <label class="hall-field"><span>{{ t('admin.providerHall.cachePrice') }}</span><input v-model="draft.reference_cache_price" class="input" inputmode="decimal" pattern="[0-9]+(\.[0-9]{1,10})?" required @input="referenceEdited" /></label>
          <label class="hall-field"><span>{{ t('admin.providerHall.cacheRate') }}</span><input v-model="cachePercent" class="input" type="number" min="0" max="100" step="any" required @input="referenceEdited" /></label>
          <div class="hall-field">
            <span>{{ t('admin.providerHall.confirmedAt') }}</span>
            <div class="flex flex-wrap items-center gap-2">
              <time class="text-xs">{{ draft.reference_confirmed_at ? new Date(draft.reference_confirmed_at).toLocaleString(locale) : t('admin.providerHall.none') }}</time>
              <button type="button" class="btn btn-secondary" @click="draft.reference_confirmed_at = new Date().toISOString()"><Check :size="16" />{{ t('admin.providerHall.confirmNow') }}</button>
            </div>
          </div>
        </template>
      </fieldset>
    </form>
    <template #footer>
      <div class="flex w-full flex-wrap justify-end gap-3">
        <button type="button" class="btn btn-secondary" :disabled="saving" @click="close">{{ t('common.cancel') }}</button>
        <button type="submit" form="hall-profile-form" class="btn btn-primary flex items-center gap-2" :disabled="saving"><Save :size="16" />{{ t(saving ? 'common.saving' : 'common.save') }}</button>
      </div>
    </template>
  </BaseDialog>

  <BaseDialog :show="candidatesOpen" :title="t('admin.providerHall.batchProfiles')" width="extra-wide" :close-on-escape="!creating" :show-close-button="!creating" @close="candidatesOpen = false">
    <div class="space-y-4">
      <p v-if="candidateError" class="hall-error" role="alert">{{ candidateError }}</p>
      <p class="text-sm text-txt-secondary">{{ t('admin.providerHall.batchProfilesHint') }}</p>
      <div class="flex flex-wrap items-center gap-3">
        <input v-model="candidateSearch" class="input max-w-xs" :placeholder="t('admin.providerHall.searchModel')" :aria-label="t('admin.providerHall.searchModel')" />
        <label class="flex items-center gap-2 text-sm"><input v-model="hideExisting" type="checkbox" class="h-4 w-4" />{{ t('admin.providerHall.hideExisting') }}</label>
        <span class="flex-1" />
        <button type="button" class="btn btn-secondary btn-sm" @click="toggleAll">{{ allSelected ? t('admin.providerHall.clearSelection') : t('admin.providerHall.selectAllVisible') }}</button>
      </div>
      <div class="max-h-96 overflow-auto border-y border-gray-200 py-2 dark:border-dark-700">
        <label v-for="c in visibleCandidates" :key="`${c.model}:${c.protocol}`" class="flex items-start gap-3 py-2 text-sm">
          <input v-model="selection" type="checkbox" :value="candidateKey(c)" :disabled="!c.available" class="mt-1" />
          <span class="min-w-0 flex-1 break-all">
            <span class="font-medium">{{ c.model }}</span> / {{ protocolLabel(c.protocol) }}
            <span v-if="existingProfile(c)" class="ml-2 badge badge-secondary">{{ t('admin.providerHall.existingProfile') }}</span>
            <span class="block text-xs text-txt-dim">{{ sourceLabel(c) }} · {{ t('admin.providerHall.availableInGroups', { count: c.groups.length }) }}<span v-if="!c.available"> · {{ t('admin.providerHall.invalidRoute') }}</span></span>
          </span>
        </label>
        <p v-if="!visibleCandidates.length && !candidatesLoading" class="py-6 text-sm text-txt-dim">{{ t('admin.providerHall.noModels') }}</p>
        <p v-if="candidatesLoading" class="py-6 text-sm text-txt-dim">{{ t('common.loading') }}</p>
      </div>
      <p class="text-xs text-txt-dim">{{ t('admin.providerHall.importConfirm') }}</p>
    </div>
    <template #footer>
      <div class="flex w-full flex-wrap justify-end gap-3">
        <button type="button" class="btn btn-secondary" :disabled="creating" @click="candidatesOpen = false">{{ t('common.cancel') }}</button>
        <button type="button" class="btn btn-primary" :disabled="creating || !selection.length" @click="createSelected"><Plus :size="16" />{{ t('admin.providerHall.createProfiles', { count: selection.length }) }}</button>
      </div>
    </template>
  </BaseDialog>

  <ConfirmDialog :show="pendingDelete !== null" :title="t('admin.providerHall.deleteProfile')" :message="deleteMessage" :confirm-text="t('common.delete')" danger @confirm="confirmDelete" @cancel="pendingDelete = null" />
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { Check, Pencil, Plus, RefreshCw, Save, Search, Sparkles, Trash2 } from 'lucide-vue-next'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import Select from '@/components/common/Select.vue'
import Pagination from '@/components/common/Pagination.vue'
import DataTable from '@/components/common/DataTable.vue'
import * as hall from '@/api/admin/providerHall'
import { useAppStore } from '@/stores/app'
import { hallError, lines, protocols } from './helpers'
import { extractApiErrorCode } from '@/utils/apiError'

const props = defineProps<{ profiles: hall.ProviderHallProfile[]; config?: hall.ProviderHallConfig; loading?: boolean }>()
const emit = defineEmits<{ saved: [hall.ProviderHallProfile]; reload: []; dirty: [boolean]; config: [] }>()
const { t, locale } = useI18n()
const app = useAppStore()

const error = ref(''), saving = ref(false), creating = ref(false), candidateError = ref('')
const show = ref(false), editingID = ref<number | null>(null)
const aliasText = ref(''), hasReference = ref(false), typedModel = ref(false)
const draft = ref<hall.ProviderHallProfileInput>(emptyProfile())
const search = ref(''), groupFilter = ref(''), referenceFilter = ref(''), page = ref(1), pageSize = ref(20)
const candidateOptions = ref<hall.ProviderHallProfileCandidate[]>([])
const pickValue = ref<string | null>(null)

const candidatesOpen = ref(false), candidatesLoading = ref(false), candidateSearch = ref('')
const hideExisting = ref(true), selection = ref<string[]>([]), pendingDelete = ref<hall.ProviderHallProfile | null>(null)

const columns = computed(() => [
  { key: 'model', label: t('admin.providerHall.model'), sortable: true },
  { key: 'protocol', label: t('admin.providerHall.protocol') },
  { key: 'groups', label: t('admin.providerHall.linkedGroups') },
  { key: 'candidates', label: t('admin.providerHall.upstreamAvailability') },
  { key: 'output_limit', label: t('admin.providerHall.outputLimit'), sortable: true },
  { key: 'tools', label: t('admin.providerHall.tools') },
  { key: 'reference', label: t('admin.providerHall.reference') },
  { key: 'actions', label: t('admin.providerHall.jobActions') }
])

const groupOptions = computed(() => {
  const seen = new Map<number, string>()
  for (const p of props.profiles) for (const g of p.groups ?? []) seen.set(g.id, g.name)
  return [...seen].map(([value, label]) => ({ value, label }))
})
const filtered = computed(() => props.profiles.filter(p => {
  const haystack = `${p.model} ${p.protocol} ${p.model_aliases.join(' ')} ${(p.groups ?? []).map(g => g.name).join(' ')}`.toLowerCase()
  if (!haystack.includes(search.value.toLowerCase())) return false
  if (groupFilter.value && !(p.groups ?? []).some(g => String(g.id) === String(groupFilter.value))) return false
  if (referenceFilter.value === 'valid' && !referenceValid(p)) return false
  if (referenceFilter.value === 'invalid' && referenceValid(p)) return false
  return true
}))
const paged = computed(() => filtered.value.slice((page.value - 1) * pageSize.value, page.value * pageSize.value))
const editingProfile = computed(() => props.profiles.find(p => p.id === editingID.value))
const cachePercent = computed({ get: () => draft.value.reference_cache_rate === null ? '' : String(Number(draft.value.reference_cache_rate) * 100), set: value => { draft.value.reference_cache_rate = value === '' ? null : String(Number(value) / 100) } })

const candidateIndex = computed(() => new Map(candidateOptions.value.map(c => [candidateKey(c), c])))
const visibleCandidates = computed(() => candidateOptions.value.filter(c => {
  if (candidateSearch.value && !c.model.toLowerCase().includes(candidateSearch.value.toLowerCase())) return false
  if (hideExisting.value && existingProfile(c)) return false
  return true
}))
const allSelected = computed(() => visibleCandidates.value.length > 0 && visibleCandidates.value.every(c => selection.value.includes(candidateKey(c))))

const pickOptions = computed(() => candidateOptions.value.filter(c => !existingProfile(c)).map(c => ({ value: candidateKey(c), label: `${c.model} / ${protocolLabel(c.protocol)}` })))
const snapshot = computed(() => JSON.stringify({ ...draft.value, model_aliases: lines(aliasText.value), hasReference: hasReference.value, typedModel: typedModel.value }))
const initial = ref('')
const dirty = computed(() => show.value && snapshot.value !== initial.value)

function candidateKey(c: { model: string; protocol: string }) { return `${c.model}:${c.protocol}` }
function protocolLabel(value: string) { return protocols.find(p => p.value === value)?.label ?? value }
function sourceLabel(c: hall.ProviderHallProfileCandidate) { return (c.sources ?? [c.source]).map(s => t(`admin.providerHall.source_${s}`)).join(' / ') }
function existingProfile(c: { model: string; protocol: string }) { return props.profiles.find(p => p.model === c.model && p.protocol === c.protocol) }
function groupNames(p: hall.ProviderHallProfile) { return (p.groups ?? []).map(g => g.name).join(', ') }
function groupCount(p: hall.ProviderHallProfile) {
  // How many groups could probe this profile, from the merged candidate list.
  return (candidateIndex.value.get(candidateKey(p))?.groups ?? p.groups?.map(g => g.id) ?? []).length
}
function referenceValid(p: hall.ProviderHallProfile) {
  return Boolean(p.reference_confirmed_at) && Date.now() - Date.parse(p.reference_confirmed_at!) <= 7 * 86400000 &&
    p.reference_input_price !== null && p.reference_cache_price !== null && p.reference_cache_rate !== null
}
function isDefault(p: hall.ProviderHallProfile) { return props.config ? p.model === props.config.default_model && p.protocol === props.config.default_protocol : p.is_default === true }
function setPageSize(value: number) { pageSize.value = value; page.value = 1 }

function emptyProfile(): hall.ProviderHallProfileInput {
  return { version: 0, model: '', protocol: 'responses', supports_tools: false, output_limit: 256, model_aliases: [],
    reference_input_price: null, reference_cache_price: null, reference_cache_rate: null, reference_confirmed_at: null }
}
function pickModel(value: unknown) {
  const candidate = candidateOptions.value.find(c => candidateKey(c) === value)
  pickValue.value = candidate ? String(value) : null
  if (candidate) { draft.value.model = candidate.model; draft.value.protocol = candidate.protocol }
}
function open(profile?: hall.ProviderHallProfile) {
  draft.value = profile ? { ...profile, model_aliases: [...profile.model_aliases] } : emptyProfile()
  editingID.value = profile?.id ?? null
  aliasText.value = draft.value.model_aliases.join('\n')
  hasReference.value = draft.value.reference_input_price !== null
  // Editing keeps whatever the profile already is; creating starts from the picker.
  typedModel.value = Boolean(profile)
  pickValue.value = null
  error.value = ''
  initial.value = snapshot.value
  show.value = true
}
function referenceEdited() { draft.value.reference_confirmed_at = null }
function close() { if (!saving.value && (!dirty.value || window.confirm(t('admin.providerHall.discard')))) show.value = false }

async function save() {
  if (saving.value) return
  const d = draft.value
  if (hasReference.value && (!d.reference_confirmed_at || !d.reference_input_price || !d.reference_cache_price || !d.reference_cache_rate)) {
    error.value = t('admin.providerHall.referenceRequired'); return
  }
  saving.value = true
  error.value = ''
  const input: hall.ProviderHallProfileInput = { version: d.version, model: d.model, protocol: d.protocol, supports_tools: d.supports_tools,
    output_limit: d.output_limit, model_aliases: lines(aliasText.value),
    reference_input_price: hasReference.value ? d.reference_input_price : null,
    reference_cache_price: hasReference.value ? d.reference_cache_price : null,
    reference_cache_rate: hasReference.value ? d.reference_cache_rate : null,
    reference_confirmed_at: hasReference.value ? d.reference_confirmed_at : null }
  try {
    const result = editingID.value ? await hall.updateProfile(editingID.value, input) : await hall.createProfile(input)
    if (disposed) return
    show.value = false
    emit('saved', result)
  } catch (err) { if (!disposed) error.value = hallError(err, t) }
  finally { saving.value = false }
}

async function loadCandidates() {
  candidatesLoading.value = true
  candidateError.value = ''
  try {
    const result = await hall.listProfileCandidates()
    if (!disposed) candidateOptions.value = result
  } catch (err) { if (!disposed) candidateError.value = hallError(err, t) }
  finally { if (!disposed) candidatesLoading.value = false }
}
function openCandidates() {
  candidatesOpen.value = true
  selection.value = []
  void loadCandidates()
}
function toggleAll() {
  selection.value = allSelected.value ? [] : visibleCandidates.value.filter(c => c.available).map(candidateKey)
}
async function createSelected() {
  if (creating.value || !selection.value.length) return
  creating.value = true
  candidateError.value = ''
  let created = 0, skipped = 0
  try {
    for (const key of selection.value) {
      const candidate = candidateIndex.value.get(key)
      if (!candidate || existingProfile(candidate)) { skipped++; continue }
      try {
        const result = await hall.createProfile({ ...emptyProfile(), model: candidate.model, protocol: candidate.protocol })
        if (disposed) return
        emit('saved', result)
        created++
      } catch (err) { if (extractApiErrorCode(err) !== 'PROVIDER_HALL_PROFILE_EXISTS') throw err }
    }
    if (created) app.showSuccess(t('admin.providerHall.profilesCreated', { count: created }))
    if (skipped) candidateError.value = t('admin.providerHall.profilesSkipped', { count: skipped })
    if (created) candidatesOpen.value = false
  } catch (err) { if (!disposed) candidateError.value = hallError(err, t) }
  finally { creating.value = false }
}

function askDelete(p: hall.ProviderHallProfile) { pendingDelete.value = p }
const deleteMessage = computed(() => t('admin.providerHall.deleteProfileConfirm', { name: pendingDelete.value?.model ?? '' }))
async function confirmDelete() {
  const p = pendingDelete.value
  pendingDelete.value = null
  if (!p) return
  error.value = ''
  try {
    await hall.deleteProfile(p.id)
    if (disposed) return
    app.showSuccess(t('admin.providerHall.profileDeleted'))
    emit('reload')
  } catch (err) { if (!disposed) error.value = hallError(err, t) }
}

watch([search, groupFilter, referenceFilter], () => { page.value = 1 })
watch(dirty, value => emit('dirty', value))
let disposed = false
onMounted(() => { void loadCandidates() })
onUnmounted(() => { disposed = true })
</script>

<style>
.modal-content:has(#hall-profile-form) { @apply rounded-lg bg-white dark:bg-zinc-900; backdrop-filter: none; }
</style>
