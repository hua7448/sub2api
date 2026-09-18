<template>
  <TablePageLayout>
    <template #filters>
      <div class="flex flex-col justify-between gap-4 lg:flex-row lg:items-center">
        <div class="flex flex-1 flex-wrap items-center gap-3">
          <div class="relative w-full sm:w-64">
            <Search class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400 dark:text-gray-500" />
            <input v-model="filter.search" class="input pl-9" :placeholder="t('admin.providerHall.searchGroups')" :aria-label="t('admin.providerHall.searchGroups')" />
          </div>
          <select v-model="filter.platform" class="input w-auto" :aria-label="t('admin.providerHall.platform')">
            <option value="">{{ t('admin.providerHall.allPlatforms') }}</option>
            <option v-for="p in platforms" :key="p" :value="p">{{ p }}</option>
          </select>
          <select v-model="filter.listed" class="input w-auto" :aria-label="t('admin.providerHall.listed')">
            <option value="">{{ t('admin.providerHall.allStatuses') }}</option>
            <option value="true">{{ t('admin.providerHall.listed') }}</option>
            <option value="false">{{ t('admin.providerHall.unlisted') }}</option>
          </select>
          <select v-model="sortMode" class="input w-auto" :aria-label="t('admin.providerHall.displayOrder')">
            <option value="name">{{ t('admin.providerHall.nameAsc') }}</option>
            <option value="name_desc">{{ t('admin.providerHall.nameDesc') }}</option>
            <option value="display_order">{{ t('admin.providerHall.sortByDisplayOrder') }}</option>
          </select>
        </div>
        <div class="flex w-full flex-wrap items-center justify-end gap-3 lg:w-auto">
          <button class="btn btn-secondary btn-icon" :disabled="loading" :title="t('admin.providerHall.reload')" :aria-label="t('admin.providerHall.reload')" @click="reloadAll">
            <RefreshCw :size="16" :class="loading && 'animate-spin'" />
          </button>
          <button v-if="sortMode === 'display_order'" class="btn btn-secondary" @click="openSort">
            <ArrowUpDown :size="16" />{{ t('admin.providerHall.sortOrder') }}
          </button>
          <button class="btn btn-secondary" :disabled="!selected.length" @click="openBatch">
            <Layers :size="16" />{{ t('admin.providerHall.batch') }} ({{ selected.length }})
          </button>
        </div>
      </div>
    </template>

    <template #table>
      <p v-if="error && !groups.length" class="hall-error m-4" role="alert">{{ error }}</p>
      <DataTable
        :columns="columns"
        :data="groups"
        :loading="loading"
        row-key="group_id"
        selectable
        expandable
        :expanded-keys="expanded"
        :selected-keys="selected"
        :expand-label="row => t('admin.providerHall.expandTargets', { name: row.display_name || row.name })"
        @update:selected-keys="selectKeys"
        @update:expanded-keys="setExpanded"
        @expand="onExpand"
      >
        <template #cell-name="{ row: g }">
          <div class="min-w-0">
            <span class="font-medium text-txt-primary">{{ g.name }}</span>
            <span v-if="g.display_name" class="ml-2 text-xs text-txt-dim">{{ g.display_name }}</span>
            <span v-if="!g.supported" class="ml-2 text-xs text-amber-600">{{ t('admin.providerHall.unsupported') }}</span>
            <span v-if="dirtyGroups.has(g.group_id)" class="ml-2 text-xs text-amber-600">{{ t('admin.providerHall.unsaved') }}</span>
            <span v-if="g.description" class="block truncate text-xs text-txt-dim">{{ g.description }}</span>
          </div>
        </template>
        <template #cell-listed="{ row: g }">
          <span :class="g.listed ? 'badge badge-success' : 'badge badge-secondary'">{{ t(g.listed ? 'admin.providerHall.listed' : 'admin.providerHall.unlisted') }}</span>
          <span v-if="g.status !== 'active'" class="block text-xs text-amber-600">{{ t('admin.providerHall.inactive') }}</span>
        </template>
        <template #cell-platform="{ row: g }">{{ g.platform }}</template>
        <template #cell-targets="{ row: g }">
          <span class="text-sm">{{ t('admin.providerHall.targetCount', { count: g.targets.length }) }}</span>
          <span class="block max-w-xs truncate text-xs text-txt-dim">{{ g.models.join(', ') || t('admin.providerHall.noTargets') }}</span>
        </template>
        <template #cell-keys="{ row: g }">
          <span v-for="entry in keySummary(g)" :key="entry.label" class="mr-2 whitespace-nowrap text-xs" :class="{ 'text-amber-600': entry.warn }">{{ entry.label }}: {{ entry.count }}</span>
          <span v-if="!keySummary(g).length" class="text-xs text-txt-dim">—</span>
        </template>
        <template #cell-results="{ row: g }">
          <button v-for="result in [g.latest_probe, g.latest_verification].filter(Boolean)" :key="result.job_id" type="button" class="block text-left text-xs text-emerald-600 hover:underline" @click.stop="emit('open', result.job_id)">
            {{ t(`admin.providerHall.status_${result.status}`) }} {{ result.verdict ? t(`admin.providerHall.verdict_${result.verdict}`) : '' }}
            <span class="block text-txt-dim">{{ formatDate(result.at) }}</span>
          </button>
          <span v-if="!g.latest_probe && !g.latest_verification" class="text-xs text-txt-dim">—</span>
        </template>
        <template #cell-actions="{ row: g }">
          <button class="btn btn-ghost btn-sm" :disabled="!g.supported" @click.stop="editListing(g)">
            <Settings2 :size="15" />{{ t('admin.providerHall.manage') }}
          </button>
        </template>

        <template #expanded="{ row: g }">
          <div v-if="!g.supported" class="p-3 text-sm text-amber-600">{{ t('admin.providerHall.unsupported') }}</div>
          <div v-else class="space-y-4 p-3">
            <p v-if="errorFor(g.group_id)" class="hall-error" role="alert">{{ errorFor(g.group_id) }}</p>
            <div class="flex flex-wrap items-center gap-2">
              <Select v-model="newProfile[g.group_id]" class="min-w-0 flex-1 basis-56" :options="addableProfiles(g)" :disabled="busy(g.group_id)" searchable :aria-label="t('admin.providerHall.profiles')" />
              <button type="button" class="btn btn-secondary" :disabled="busy(g.group_id) || !newProfile[g.group_id]" @click="addTarget(g)">
                <Plus :size="16" />{{ t('admin.providerHall.addTarget') }}
              </button>
              <button type="button" class="btn btn-secondary" :disabled="busy(g.group_id)" @click="loadCandidates(g)">
                <RefreshCw :size="16" />{{ t('admin.providerHall.chooseModels') }}
              </button>
              <button type="button" class="btn btn-secondary" :disabled="busy(g.group_id)" @click="refreshModels(g)">
                <RefreshCw :size="16" />{{ t('admin.providerHall.refreshUpstream') }}
              </button>
              <span class="flex-1" />
              <button type="button" class="btn btn-secondary" :disabled="busy(g.group_id)" @click="preflight(g)">{{ t('admin.providerHall.preflight') }}</button>
              <button type="button" class="btn btn-primary" :disabled="busy(g.group_id)" @click="saveGroup(g)">
                <Save :size="16" />{{ t('admin.providerHall.saveGroup') }}
              </button>
            </div>
            <div v-if="candidateOpen[g.group_id]" class="max-h-64 overflow-auto border-y border-gray-200 py-3 dark:border-dark-700">
              <label v-for="(m, index) in candidates[g.group_id] || []" :key="`${m.model}:${m.protocol}:${m.source}:${index}`" class="flex items-start gap-3 py-2 text-sm">
                <input v-model="candidateSelection[g.group_id]" type="checkbox" :value="index" :disabled="!m.available || busy(g.group_id)" class="mt-1" />
                <span class="min-w-0 break-all">{{ m.model }} / {{ protocolLabel(m.protocol) }}
                  <span class="block text-xs text-txt-dim">{{ t(`admin.providerHall.source_${m.source}`) }} · {{ m.upstream_model }}<span v-if="!m.available"> · {{ t('admin.providerHall.invalidRoute') }}</span></span>
                </span>
              </label>
              <p v-if="!(candidates[g.group_id] || []).length" class="text-sm text-txt-dim">{{ t('admin.providerHall.noModels') }}</p>
              <button type="button" class="btn btn-secondary mt-3" :disabled="busy(g.group_id) || !(candidateSelection[g.group_id] || []).length" @click="importModels(g)">{{ t('admin.providerHall.importSelected', { count: (candidateSelection[g.group_id] || []).length }) }}</button>
            </div>
            <p v-if="refreshResult[g.group_id]" class="text-sm" role="status">{{ refreshResult[g.group_id] }}</p>
            <p v-if="!operatorId" class="text-sm text-amber-600">{{ t('admin.providerHall.operatorRequired') }}</p>
            <p class="text-xs text-txt-dim">{{ t('admin.providerHall.keyNotice') }}</p>

            <div v-if="(draft[g.group_id] || []).length" class="overflow-x-auto">
              <table class="hall-table min-w-[960px]">
                <thead>
                  <tr>
                    <th class="w-[22%]">{{ t('admin.providerHall.model') }}</th>
                    <th class="w-[24%]">{{ t('admin.providerHall.probeKey') }}</th>
                    <th class="w-[11%]">{{ t('admin.providerHall.probeMinutes') }}</th>
                    <th class="w-[11%]">{{ t('admin.providerHall.verifyHours') }}</th>
                    <th class="w-[7%]">{{ t('admin.providerHall.enabled') }}</th>
                    <th class="w-[8%]">{{ t('admin.providerHall.autoSchedule') }}</th>
                    <th class="w-[17%]">{{ t('admin.providerHall.jobActions') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="row in draft[g.group_id]" :key="row.profile_id">
                    <td class="break-all">
                      {{ profileName(row.profile_id) }}
                      <span class="block text-xs text-txt-dim">{{ protocolLabel(profileProtocol(row.profile_id)) }}</span>
                    </td>
                    <td>
                      <Select v-model="row.probe_key_id" :options="keyOptions(g, row.probe_key_id)" :disabled="busy(g.group_id)" clearable searchable :aria-label="`${t('admin.providerHall.probeKey')} ${profileName(row.profile_id)}`" />
                      <button type="button" class="btn btn-ghost btn-sm mt-1" :disabled="busy(g.group_id) || !operatorId" @click="ensureKey(g, row)">
                        <KeyRound :size="14" />{{ t('admin.providerHall.ensureKey') }}
                      </button>
                    </td>
                    <td><input :value="(row.probe_interval_seconds || 300) / 60" class="input" type="number" min="1" max="1440" step="any" :disabled="busy(g.group_id)" :aria-label="t('admin.providerHall.probeMinutes')" @input="row.probe_interval_seconds = Math.round(Number(($event.target as HTMLInputElement).value) * 60)" /></td>
                    <td><input :value="(row.verification_interval_seconds || 86400) / 3600" class="input" type="number" min="1" max="168" step="any" :disabled="busy(g.group_id)" :aria-label="t('admin.providerHall.verifyHours')" @input="row.verification_interval_seconds = Math.round(Number(($event.target as HTMLInputElement).value) * 3600)" /></td>
                    <td class="text-center"><Toggle :model-value="row.enabled" @update:model-value="row.enabled = $event" /></td>
                    <td class="text-center"><Toggle :model-value="row.auto_schedule_enabled === true" @update:model-value="row.auto_schedule_enabled = $event" /></td>
                    <td class="whitespace-nowrap">
                      <button type="button" class="btn btn-ghost btn-sm whitespace-normal" :disabled="busy(g.group_id)" @click="test(g, row, 'probe')">{{ t(!row.enabled ? 'admin.providerHall.enableAndTest' : isDirty(g.group_id) ? 'admin.providerHall.saveAndTest' : 'admin.providerHall.probeNow') }}</button>
                      <button type="button" class="btn btn-ghost btn-sm" :disabled="busy(g.group_id)" @click="test(g, row, 'verification')">{{ t('admin.providerHall.verifyNow') }}</button>
                      <button type="button" class="btn btn-ghost btn-sm" :disabled="busy(g.group_id)" @click="viewJobs(g, row)">{{ t('admin.providerHall.results') }}</button>
                      <button type="button" class="btn btn-ghost btn-sm text-red-600 dark:text-red-400" :disabled="busy(g.group_id)" :aria-label="`${t('common.delete')} ${profileName(row.profile_id)}`" @click="removeTarget(g, row)"><Trash2 :size="14" /></button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <p v-else class="py-4 text-sm text-txt-dim">{{ t('admin.providerHall.noTargets') }}</p>
            <button v-if="lastJob[g.group_id]" type="button" class="btn btn-secondary" @click="emit('open', lastJob[g.group_id]!)">{{ t('admin.providerHall.jobDetail') }} #{{ lastJob[g.group_id] }}</button>
          </div>
        </template>
      </DataTable>
    </template>

    <template #pagination>
      <Pagination v-if="total" :total="total" :page="page" :page-size="pageSize" @update:page="setPage" @update:page-size="setPageSize" />
    </template>
  </TablePageLayout>

  <BaseDialog :show="sortOpen" :title="t('admin.providerHall.sortOrder')" width="normal" @close="sortOpen = false">
    <p class="mb-4 text-sm text-txt-secondary">{{ t('admin.providerHall.sortOrderHint') }}</p>
    <VueDraggable v-model="sortableGroups" :animation="200" class="space-y-2">
      <div v-for="g in sortableGroups" :key="g.group_id" class="flex cursor-grab items-center gap-3 rounded-lg border border-gray-200 bg-white p-3 transition-shadow hover:shadow-md active:cursor-grabbing dark:border-dark-600 dark:bg-dark-700">
        <GripVertical :size="16" class="text-gray-400" />
        <div class="min-w-0 flex-1">
          <div class="truncate font-medium text-txt-primary">{{ g.name }}</div>
          <div class="text-xs text-txt-dim">{{ g.platform }} · #{{ g.group_id }}</div>
        </div>
        <span :class="g.listed ? 'badge badge-success' : 'badge badge-secondary'">{{ t(g.listed ? 'admin.providerHall.listed' : 'admin.providerHall.unlisted') }}</span>
      </div>
    </VueDraggable>
    <template #footer>
      <div class="flex justify-end gap-3">
        <button type="button" class="btn btn-secondary" :disabled="sortSaving" @click="sortOpen = false">{{ t('common.cancel') }}</button>
        <button type="button" class="btn btn-primary" :disabled="sortSaving" @click="saveSort"><Save :size="16" />{{ t('common.save') }}</button>
      </div>
    </template>
  </BaseDialog>

  <BaseDialog :show="editing !== null" :title="editing?.name || t('admin.providerHall.manage')" width="wide" :close-on-escape="!editingBusy" :show-close-button="!editingBusy" @close="closeEditing">
    <form v-if="editing" id="hall-listing" class="space-y-5" @submit.prevent="saveListing">
      <p v-if="editingError" class="hall-error" role="alert">{{ editingError }}</p>
      <fieldset :disabled="editingBusy" class="hall-fields">
        <label class="hall-field"><span>{{ t('admin.providerHall.displayName') }}</span><input v-model="editing.display_name" class="input" maxlength="100" /></label>
        <label class="hall-field"><span>{{ t('admin.providerHall.displayOrder') }}</span><input v-model.number="editing.display_order" class="input" type="number" step="1" /></label>
        <label class="hall-field sm:col-span-2"><span>{{ t('admin.providerHall.description') }}</span><textarea v-model="editing.description" class="input" maxlength="2000" rows="2" /></label>
        <label class="flex items-center gap-2 text-sm sm:col-span-2"><Toggle :model-value="editing.listed" @update:model-value="editing!.listed = $event" />{{ t('admin.providerHall.listed') }}</label>
      </fieldset>
      <p class="text-xs text-txt-dim">{{ t('admin.providerHall.listingHint') }}</p>
    </form>
    <template #footer>
      <div class="flex w-full flex-wrap justify-end gap-3">
        <button type="button" class="btn btn-secondary" :disabled="editingBusy" @click="closeEditing">{{ t('common.cancel') }}</button>
        <button type="submit" form="hall-listing" class="btn btn-primary" :disabled="editingBusy"><Save :size="16" />{{ t('common.save') }}</button>
      </div>
    </template>
  </BaseDialog>

  <ProviderHallBatchDialog :show="batchOpen" :groups="batchSelection" :profiles="allProfiles" @close="batchOpen = false" @saved="onBatchSaved" />
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ArrowUpDown, GripVertical, KeyRound, Layers, Plus, RefreshCw, Save, Search, Settings2, Trash2 } from 'lucide-vue-next'
import { VueDraggable } from 'vue-draggable-plus'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Pagination from '@/components/common/Pagination.vue'
import DataTable from '@/components/common/DataTable.vue'
import Select from '@/components/common/Select.vue'
import Toggle from '@/components/common/Toggle.vue'
import ProviderHallBatchDialog from './ProviderHallBatchDialog.vue'
import * as hall from '@/api/admin/providerHall'
import { useAppStore } from '@/stores/app'
import { formatDate } from '@/utils/format'
import { hallError, idempotencyKey, protocols } from './helpers'
import { extractApiErrorCode } from '@/utils/apiError'

const props = defineProps<{ profiles: hall.ProviderHallProfile[]; operatorId: number | null; focusGroup?: number | null }>()
const emit = defineEmits<{ dirty: [boolean]; open: [number]; jobs: [number, number]; profile: [hall.ProviderHallProfile] }>()
const { t } = useI18n()
const app = useAppStore()
const platforms = ['openai', 'composite', 'anthropic', 'gemini', 'antigravity', 'grok', 'deepseek', 'kimi', 'zhipu']

const groups = ref<hall.ProviderHallGroupSummary[]>([])
const total = ref(0), page = ref(1), pageSize = ref(20), selected = ref<number[]>([])
const filter = reactive({ search: '', platform: '', listed: '' })
const sortMode = ref<'name' | 'name_desc' | 'display_order'>('display_order')
const loading = ref(false), error = ref('')

// Per-group draft state. Each expanded group keeps its own editable row set,
// key/ candidate caches and error, so opening a second group never discards
// or overwrites what the admin typed in the first.
const draft = reactive<Record<number, hall.ProviderHallTargetInput[]>>({})
const baseline = reactive<Record<number, string>>({})
// The target-set version for the expanded group. It is what makes the
// optimistic concurrency check work, so it must come from the last read or
// write of that group rather than from the list row.
const targetVersion = reactive<Record<number, number>>({})
const keys = reactive<Record<number, hall.ProviderHallProbeKeyOption[]>>({})
const candidates = reactive<Record<number, hall.ProviderHallModelCandidate[]>>({})
const candidateSelection = reactive<Record<number, number[]>>({})
const candidateOpen = reactive<Record<number, boolean>>({})
const refreshResult = reactive<Record<number, string>>({})
const lastJob = reactive<Record<number, number | null>>({})
const errors = reactive<Record<number, string>>({})
const busyGroups = ref<number[]>([])
const newProfile = reactive<Record<number, number | null>>({})
const localProfiles = ref<hall.ProviderHallProfile[]>([])
const expanded = ref<number[]>([])
const pendingKeys = new Map<string, string>()
const dirtyGroups = computed(() => new Set(Object.keys(baseline).map(Number).filter(id => isDirty(id))))
const allProfiles = computed(() => [...props.profiles, ...localProfiles.value.filter(p => !props.profiles.some(x => x.id === p.id))])

const columns = computed(() => [
  { key: 'name', label: t('admin.providerHall.jobGroup'), sortable: true },
  { key: 'platform', label: t('admin.providerHall.platform') },
  { key: 'listed', label: t('admin.providerHall.listed') },
  { key: 'targets', label: t('admin.providerHall.targets') },
  { key: 'keys', label: t('admin.providerHall.probeKey') },
  { key: 'results', label: t('admin.providerHall.results') },
  { key: 'actions', label: t('admin.providerHall.jobActions') }
])

let request: AbortController | undefined, timer: ReturnType<typeof setTimeout> | undefined
let disposed = false

function busy(id: number) { return busyGroups.value.includes(id) }
function errorFor(id: number) { return errors[id] ?? '' }
function profileName(id: number) { return allProfiles.value.find(p => p.id === id)?.model ?? `#${id}` }
function profileProtocol(id: number) { return allProfiles.value.find(p => p.id === id)?.protocol ?? 'responses' }
function protocolLabel(value: string) { return protocols.find(p => p.value === value)?.label ?? value }
/** Probe-key status counts, labelled and flagged for the ones that need action. */
function keySummary(g: hall.ProviderHallGroupSummary) {
  const statuses = g.key_statuses ?? {}
  return Object.keys(statuses)
    .filter(status => statuses[status] > 0)
    .map(status => ({ key: status, label: t(`admin.providerHall.key_${status}`), count: statuses[status], warn: status === 'missing' || status === 'inactive' }))
}
function payloadFor(id: number): hall.ProviderHallSettingsInput {
  const g = groups.value.find(x => x.group_id === id)!
  return { version: targetVersion[id] ?? 0, listed: g.listed, display_name: g.display_name, description: g.description, display_order: g.display_order, items: draft[id] ?? [] }
}
function isDirty(id: number) { return (id in baseline) && JSON.stringify(draft[id] ?? []) !== baseline[id] }

async function load() {
  request?.abort()
  const current = new AbortController()
  request = current
  loading.value = true
  error.value = ''
  try {
    const result = await hall.listGroups({ ...filter, sort: sortMode.value, page: page.value, page_size: pageSize.value }, current.signal)
    if (current.signal.aborted) return
    groups.value = result.items
    total.value = result.total
    // Drop drafts whose group left the current page, but keep dirty ones so a
    // filter change cannot silently discard typed work.
    for (const id of Object.keys(baseline).map(Number)) {
      if (!groups.value.some(g => g.group_id === id) && !isDirty(id)) clearDraft(id)
    }
  } catch (err) { if (!current.signal.aborted) error.value = hallError(err, t, 'loadFailed') }
  finally { if (!current.signal.aborted) loading.value = false }
}

function clearDraft(id: number) {
  delete draft[id]; delete baseline[id]; delete errors[id]; delete targetVersion[id]
  delete candidates[id]; delete candidateSelection[id]; delete candidateOpen[id]; delete refreshResult[id]
}

async function openGroup(id: number) {
  try {
    const [set, ks, ms] = await Promise.all([hall.getTargets(id), hall.listProbeKeys(id), hall.listModels(id)])
    if (disposed) return
    draft[id] = set.items.map(r => ({ profile_id: r.profile_id, probe_key_id: r.probe_key_id, enabled: r.enabled, auto_schedule_enabled: r.auto_schedule_enabled ?? false, probe_interval_seconds: r.probe_interval_seconds, verification_interval_seconds: r.verification_interval_seconds }))
    baseline[id] = JSON.stringify(draft[id])
    targetVersion[id] = set.version
    keys[id] = ks
    candidates[id] = ms
    candidateSelection[id] = []
    candidateOpen[id] = false
    refreshResult[id] = ''
    errors[id] = ''
  } catch (err) { if (!disposed) errors[id] = hallError(err, t, 'loadFailed') }
}

function setExpanded(keysList: (number | string)[]) {
  const next = keysList.map(Number)
  // Collapsing a group with unsaved edits asks first. Declining writes the
  // previous keys back, which the table treats as "keep it open".
  const closing = expanded.value.filter(id => !next.includes(id))
  if (closing.some(id => isDirty(id)) && !window.confirm(t('admin.providerHall.discard'))) {
    expanded.value = []
    expanded.value = next.filter(id => !closing.includes(id))
    return
  }
  expanded.value = next
  for (const id of closing) if (!isDirty(id)) clearDraft(id)
}
// Expansion is the load trigger: a group's targets, keys and candidates are
// fetched only when the admin actually opens it.
function onExpand(key: string | number, isExpanded: boolean) {
  if (isExpanded) void openGroup(Number(key))
}
function selectKeys(keysList: (number | string)[]) {
  selected.value = keysList.map(Number).filter(id => groups.value.some(g => g.group_id === id && g.supported))
}
function setPage(value: number) { page.value = value; selected.value = []; void load() }
function setPageSize(value: number) { pageSize.value = value; page.value = 1; selected.value = []; void load() }
async function reloadAll() {
  if (dirtyGroups.value.size && !window.confirm(t('admin.providerHall.discard'))) return
  for (const id of Object.keys(baseline).map(Number)) clearDraft(id)
  expanded.value = []
  await load()
}

function addableProfiles(g: hall.ProviderHallGroupSummary) {
  const used = new Set((draft[g.group_id] ?? []).map(r => r.profile_id))
  return allProfiles.value.filter(p => !used.has(p.id)).map(p => ({ value: p.id, label: `${p.model} / ${protocolLabel(p.protocol)}` }))
}
function addTarget(g: hall.ProviderHallGroupSummary) {
  const id = g.group_id
  const profileID = newProfile[id]
  if (!profileID || (draft[id] ?? []).length >= 100) return
  draft[id] = [...(draft[id] ?? []), { profile_id: profileID, probe_key_id: null, enabled: false, auto_schedule_enabled: false, probe_interval_seconds: 300, verification_interval_seconds: 86400 }]
  newProfile[id] = null
}
function removeTarget(g: hall.ProviderHallGroupSummary, row: hall.ProviderHallTargetInput) {
  const id = g.group_id
  draft[id] = (draft[id] ?? []).filter(r => r.profile_id !== row.profile_id)
}
function keyOptions(g: hall.ProviderHallGroupSummary, current: number | null) {
  const options = (keys[g.group_id] ?? []).map(k => ({ value: k.id as number | null, label: `${k.name} #${k.id} · ${t(`admin.providerHall.key_${k.status}`)}`, disabled: !['unused', 'registered'].includes(k.status) }))
  if (current && !options.some(k => k.value === current)) options.push({ value: current, label: t('admin.providerHall.unavailableKey', { id: current }), disabled: true })
  return [{ value: null, label: t('admin.providerHall.none'), disabled: false }, ...options]
}
async function ensureKey(g: hall.ProviderHallGroupSummary, row: hall.ProviderHallTargetInput) {
  const id = g.group_id
  if (busy(id)) return
  busyGroups.value = [...busyGroups.value, id]
  errors[id] = ''
  try {
    const key = await hall.ensureProbeKey(id, row.profile_id)
    if (disposed) return
    keys[id] = [...(keys[id] ?? []).filter(k => k.id !== key.id), key]
    row.probe_key_id = key.id
  } catch (err) { errors[id] = hallError(err, t) }
  finally { busyGroups.value = busyGroups.value.filter(x => x !== id) }
}

async function persist(g: hall.ProviderHallGroupSummary) {
  const id = g.group_id
  const result = await hall.saveSettings(id, payloadFor(id))
  if (disposed) return null
  draft[id] = result.items.map(r => ({ profile_id: r.profile_id, probe_key_id: r.probe_key_id, enabled: r.enabled, auto_schedule_enabled: r.auto_schedule_enabled ?? false, probe_interval_seconds: r.probe_interval_seconds, verification_interval_seconds: r.verification_interval_seconds }))
  baseline[id] = JSON.stringify(draft[id])
  targetVersion[id] = result.version
  if (result.listing) Object.assign(g, result.listing)
  app.showSuccess(t('admin.providerHall.groupSaved'))
  void load()
  return result
}
async function saveGroup(g: hall.ProviderHallGroupSummary) {
  const id = g.group_id
  if (busy(id)) return
  busyGroups.value = [...busyGroups.value, id]
  errors[id] = ''
  try { await persist(g) } catch (err) { errors[id] = hallError(err, t) }
  finally { busyGroups.value = busyGroups.value.filter(x => x !== id) }
}
async function preflight(g: hall.ProviderHallGroupSummary) {
  const id = g.group_id
  if (busy(id)) return
  busyGroups.value = [...busyGroups.value, id]
  errors[id] = ''
  try { await hall.preflightGroup(id, payloadFor(id)); app.showSuccess(t('admin.providerHall.preflightPassed')) }
  catch (err) { errors[id] = hallError(err, t) }
  finally { busyGroups.value = busyGroups.value.filter(x => x !== id) }
}

async function test(g: hall.ProviderHallGroupSummary, row: hall.ProviderHallTargetInput, kind: hall.ProviderHallJobKind) {
  const id = g.group_id
  if (busy(id)) return
  if (!row.enabled && !window.confirm(t('admin.providerHall.enableAndTestConfirm'))) return
  row.enabled = true
  busyGroups.value = [...busyGroups.value, id]
  errors[id] = ''
  try {
    // Enqueue must use the versions the save just returned. Re-reading the
    // targets here would race the reload and send stale versions.
    let items: hall.ProviderHallTarget[] | null = null
    if (isDirty(id)) {
      const saved = await persist(g)
      if (disposed || !saved) return
      items = saved.items
    } else {
      items = (await hall.getTargets(id)).items
      if (disposed) return
    }
    const target = items.find(r => r.profile_id === row.profile_id)
    const profile = allProfiles.value.find(p => p.id === row.profile_id)
    if (!target || !profile) return
    const slot = `${kind}:${target.id}:${target.version}:${profile.version}`
    const key = pendingKeys.get(slot) ?? idempotencyKey(slot)
    pendingKeys.set(slot, key)
    const result = await (kind === 'probe' ? hall.enqueueProbe : hall.enqueueVerification)(id, row.profile_id, key, { target_version: target.version, profile_version: profile.version })
    lastJob[id] = result.job_id
    pendingKeys.delete(slot)
    app.showSuccess(t(result.reused ? 'admin.providerHall.enqueuedReused' : 'admin.providerHall.enqueued', { id: result.job_id }))
  } catch (err) { errors[id] = hallError(err, t) }
  finally { busyGroups.value = busyGroups.value.filter(x => x !== id) }
}

async function loadCandidates(g: hall.ProviderHallGroupSummary) {
  const id = g.group_id
  if (busy(id)) return
  busyGroups.value = [...busyGroups.value, id]
  try {
    candidates[id] = await hall.listModels(id)
    candidateSelection[id] = []
    candidateOpen[id] = true
  } catch (err) { errors[id] = hallError(err, t) }
  finally { busyGroups.value = busyGroups.value.filter(x => x !== id) }
}
async function refreshModels(g: hall.ProviderHallGroupSummary) {
  const id = g.group_id
  if (busy(id)) return
  busyGroups.value = [...busyGroups.value, id]
  errors[id] = ''
  try {
    const result = await hall.refreshModels(id)
    if (disposed) return
    candidates[id] = [...(candidates[id] ?? []).filter(m => m.source !== 'upstream'), ...result.flatMap(r => r.models)]
    candidateSelection[id] = []
    candidateOpen[id] = true
    refreshResult[id] = result.map(r => `#${r.account_id}: ${t(r.success ? 'admin.providerHall.refreshOK' : 'admin.providerHall.refreshFailed')}`).join(' · ')
  } catch (err) { errors[id] = hallError(err, t) }
  finally { busyGroups.value = busyGroups.value.filter(x => x !== id) }
}
async function importModels(g: hall.ProviderHallGroupSummary) {
  const id = g.group_id
  if (busy(id)) return
  const chosen = (candidateSelection[id] ?? []).map(i => (candidates[id] ?? [])[i]).filter(Boolean)
  if (!window.confirm(`${t('admin.providerHall.importConfirm')}\n${chosen.map(m => `${m.model} / ${protocolLabel(m.protocol)}`).join('\n')}`)) return
  busyGroups.value = [...busyGroups.value, id]
  errors[id] = ''
  try {
    for (const m of chosen) {
      let p = allProfiles.value.find(x => x.model === m.model && x.protocol === m.protocol)
      if (!p) {
        try { p = await hall.createProfile({ version: 0, model: m.model, protocol: m.protocol, supports_tools: false, output_limit: 256, model_aliases: [], reference_input_price: null, reference_cache_price: null, reference_cache_rate: null, reference_confirmed_at: null }) }
        catch (err) {
          if (extractApiErrorCode(err) !== 'PROVIDER_HALL_PROFILE_EXISTS') throw err
          p = (await hall.listProfiles()).find(x => x.model === m.model && x.protocol === m.protocol)
          if (!p) throw err
        }
        localProfiles.value.push(p)
        emit('profile', p)
      }
      if (!(draft[id] ?? []).some(r => r.profile_id === p!.id)) {
        draft[id] = [...(draft[id] ?? []), { profile_id: p.id, probe_key_id: null, enabled: false, auto_schedule_enabled: false, probe_interval_seconds: 300, verification_interval_seconds: 86400 }]
      }
    }
    candidateSelection[id] = []
  } catch (err) { errors[id] = hallError(err, t) }
  finally { busyGroups.value = busyGroups.value.filter(x => x !== id) }
}
function viewJobs(g: hall.ProviderHallGroupSummary, row: hall.ProviderHallTargetInput) { emit('jobs', g.group_id, row.profile_id) }

// ---- Listing metadata dialog ----
const editing = ref<hall.ProviderHallGroupSummary | null>(null)
const editingBusy = ref(false), editingError = ref('')
function editListing(g: hall.ProviderHallGroupSummary) {
  editing.value = JSON.parse(JSON.stringify(g))
  editingError.value = ''
}
function closeEditing() { if (!editingBusy.value) editing.value = null }
async function saveListing() {
  const g = editing.value
  if (!g || editingBusy.value) return
  // The listing write carries the group's targets (the API saves both under one
  // version), so an open target draft would be overwritten by this save.
  if (isDirty(g.group_id) && !window.confirm(t('admin.providerHall.listingDiscardsTargetDraft'))) return
  editingBusy.value = true
  editingError.value = ''
  try {
    const set = await hall.getTargets(g.group_id)
    if (disposed) return
    const items = set.items.map(r => ({ profile_id: r.profile_id, probe_key_id: r.probe_key_id, enabled: r.enabled, auto_schedule_enabled: r.auto_schedule_enabled ?? false, probe_interval_seconds: r.probe_interval_seconds, verification_interval_seconds: r.verification_interval_seconds }))
    await hall.saveSettings(g.group_id, { version: set.version, listed: g.listed, display_name: g.display_name, description: g.description, display_order: g.display_order, items })
    if (disposed) return
    targetVersion[g.group_id] = set.version + 1
    clearDraft(g.group_id)
    editing.value = null
    app.showSuccess(t('admin.providerHall.groupSaved'))
    void load()
  } catch (err) { if (!disposed) editingError.value = hallError(err, t) }
  finally { editingBusy.value = false }
}

// ---- Ordering ----
const sortOpen = ref(false), sortSaving = ref(false), sortableGroups = ref<hall.ProviderHallGroupSummary[]>([])
async function openSort() {
  try {
    // The sort dialog spans the whole hall, not just the current page.
    const result = await hall.listGroups({ sort: 'display_order', page: 1, page_size: 100 })
    if (disposed) return
    sortableGroups.value = result.items
    sortOpen.value = true
  } catch (err) { error.value = hallError(err, t, 'loadFailed') }
}
async function saveSort() {
  if (sortSaving.value) return
  sortSaving.value = true
  error.value = ''
  try {
    // One batch call writes the whole new order; each group keeps its targets.
    const input = await Promise.all(sortableGroups.value.map(async (g, index) => {
      const set = await hall.getTargets(g.group_id)
      return {
        id: g.group_id, version: set.version, listed: g.listed, display_name: g.display_name,
        description: g.description, display_order: index,
        items: set.items.map(r => ({ profile_id: r.profile_id, probe_key_id: r.probe_key_id, enabled: r.enabled, auto_schedule_enabled: r.auto_schedule_enabled ?? false, probe_interval_seconds: r.probe_interval_seconds, verification_interval_seconds: r.verification_interval_seconds }))
      }
    }))
    const results = await hall.batchGroups(input, false)
    if (disposed) return
    const failed = results.filter(r => !r.success)
    if (failed.length) error.value = `${t('admin.providerHall.sortPartial')}: ${failed.map(r => `#${r.id}`).join(', ')}`
    else { sortOpen.value = false; app.showSuccess(t('admin.providerHall.saved')) }
    void load()
  } catch (err) { if (!disposed) error.value = hallError(err, t) }
  finally { sortSaving.value = false }
}

// ---- Batch ----
const batchOpen = ref(false), batchSelection = ref<hall.ProviderHallGroupSummary[]>([])
function openBatch() {
  batchSelection.value = JSON.parse(JSON.stringify(groups.value.filter(g => selected.value.includes(g.group_id))))
  batchOpen.value = true
}
function onBatchSaved(failed: number[]) { selected.value = failed; void load() }

watch([filter, sortMode], () => { clearTimeout(timer); selected.value = []; page.value = 1; timer = setTimeout(() => void load(), 250) })
watch(() => props.focusGroup, id => { if (id) { expanded.value = [...new Set([...expanded.value, id])] } })
onMounted(() => { void load() })
onUnmounted(() => { disposed = true; request?.abort(); clearTimeout(timer) })
</script>

<style>
.hall-listing-row { @apply align-top; }
</style>
