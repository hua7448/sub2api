<template>
  <form class="hall-form" :aria-busy="busy" @submit.prevent="save">
    <p v-if="error" role="alert" class="hall-error">{{ error }}</p>
    <fieldset :disabled="busy" class="space-y-5">
      <h2 class="font-medium">{{ t('admin.providerHall.switchesTitle') }}</h2>
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
        <label v-for="flag in switches" :key="flag.key" class="flex items-start gap-3 text-sm">
          <input v-model="draft[flag.field]" type="checkbox" :disabled="!flag.ready" class="mt-0.5 h-4 w-4" />
          <span class="flex flex-col gap-0.5">
            <span>{{ t(`admin.providerHall.${flag.key}`) }}
              <span v-if="!flag.ready" class="text-gray-500">({{ t('admin.providerHall.inactive') }})</span></span>
            <span class="text-xs text-gray-500">{{ t(flag.ready ? `admin.providerHall.switchHint.${flag.key}` : 'admin.providerHall.switchNotReady') }}</span>
          </span>
        </label>
      </div>
      <label class="flex items-center gap-3 text-sm"><input v-model="draft.auto_schedule_enabled" type="checkbox" />{{ t('admin.providerHall.autoSchedule') }}</label>
      <p v-if="(draft.display_enabled || draft.tasks_enabled) && !draft.collection_enabled" class="text-sm text-amber-600">{{ t('admin.providerHall.fieldHint.tasks_enabled') }}</p>
      <div class="hall-fields">
        <h2 class="font-medium sm:col-span-2">{{ t('admin.providerHall.connectionSettings') }}</h2>
        <label class="hall-field sm:col-span-2">
          <span>{{ t('admin.providerHall.gatewayOrigin') }}</span>
          <input v-model.trim="draft.gateway_origin" class="input" type="url" maxlength="512" />
          <button type="button" class="btn btn-secondary" :disabled="busy || checking" @click="checkConnection">{{ t('admin.providerHall.checkGateway') }}</button>
          <span v-if="connection" role="status" class="text-xs">{{ connection }}</span>
        </label>
        <h2 class="font-medium sm:col-span-2">{{ t('admin.providerHall.taskSettings') }}</h2>
        <div class="hall-field">
          <label for="hall-operator">{{ t('admin.providerHall.operator') }}</label>
          <Select id="hall-operator" v-model="draft.operator_user_id" :options="operatorOptions" remote clearable
            :loading="usersLoading" :disabled="busy" :aria-label="t('admin.providerHall.operator')" @search="searchUsers" />
          <span v-if="selectedUser?.id === draft.operator_user_id" class="text-xs text-gray-500">{{ selectedUser.email }} · {{ t(selectedUser.status === 'active' ? 'admin.providerHall.enabled' : 'admin.providerHall.inactive') }} · {{ selectedUser.balance }} USD</span>
        </div>
        <label class="hall-field">
          <span>{{ t('admin.providerHall.budget') }}</span>
          <input v-model="draft.daily_budget" class="input" inputmode="decimal" pattern="[0-9]+(\.[0-9]{1,8})?" required />
          <span class="text-xs text-gray-500">{{ t('admin.providerHall.budgetDay') }}</span>
        </label>
        <h2 class="font-medium sm:col-span-2">{{ t('admin.providerHall.displaySettings') }}</h2>
        <label class="hall-field">
          <span>{{ t('admin.providerHall.defaultRange') }}</span>
          <select v-model="draft.default_range" class="input">
            <option v-for="range in ['6h', '24h', '7d', '30d']" :key="range" :value="range">{{ range }}</option>
          </select>
        </label>
        <div class="hall-field sm:col-span-2">
          <label for="hall-default-profile">{{ t('admin.providerHall.defaultProfile') }}</label>
          <Select id="hall-default-profile" :model-value="defaultProfileID" :options="profileOptions" :disabled="busy"
            :aria-label="t('admin.providerHall.defaultProfile')" @update:model-value="selectProfile" />
          <span class="text-xs text-gray-500">{{ profiles.find(p => p.id === defaultProfileID)?.groups?.map(g => g.name).join(', ') || t('admin.providerHall.defaultFallback') }}</span>
        </div>
        <h2 class="font-medium sm:col-span-2">{{ t('admin.providerHall.collectionSettings') }}</h2>
        <label class="hall-field sm:col-span-2">
          <span>{{ t('admin.providerHall.expectedNodes') }}</span>
          <textarea v-model="nodeText" rows="4" class="input font-mono" spellcheck="false" />
          <span class="text-xs text-gray-500">{{ t('admin.providerHall.discoveredNodes') }}</span>
          <div class="flex flex-wrap gap-2"><button v-for="node in nodes" :key="node" type="button" class="btn btn-secondary btn-sm" @click="nodeText = [...new Set([...lines(nodeText), node])].join('\n')">{{ node }}</button></div>
        </label>
      </div>
    </fieldset>
    <div class="hall-actions">
      <span class="text-xs text-gray-500">{{ t('admin.providerHall.version', { version: config.version }) }}</span>
      <button type="button" class="btn btn-secondary" :disabled="busy" @click="reload"><RefreshCw :size="16" :class="loading && 'animate-spin'" />{{ t('admin.providerHall.reload') }}</button>
      <button type="submit" class="btn btn-primary" :disabled="busy"><Save :size="16" />{{ t(saving ? 'common.saving' : 'common.save') }}</button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { RefreshCw, Save } from 'lucide-vue-next'
import Select from '@/components/common/Select.vue'
import * as hall from '@/api/admin/providerHall'
import * as users from '@/api/admin/users'
import type { AdminUser } from '@/types'
import { hallError, lines } from './helpers'

const props = defineProps<{ config: hall.ProviderHallConfig; profiles: hall.ProviderHallProfile[]; loading?: boolean }>()
const emit = defineEmits<{ saved: [hall.ProviderHallConfig]; reload: []; dirty: [boolean] }>()
const { t } = useI18n()
const draft = ref({ ...props.config })
const nodeText = ref(props.config.expected_nodes.join('\n'))
const saving = ref(false)
const busy = computed(() => saving.value || props.loading === true)
const error = ref('')
const checking = ref(false), connection = ref(''), nodes = ref<string[]>([])
async function checkConnection() {
  if (checking.value) return
  checking.value = true; connection.value = ''
  try { const result = await hall.checkGateway(draft.value.gateway_origin); connection.value = t(result.healthy ? 'admin.providerHall.gatewayHealthy' : 'admin.providerHall.gatewayUnhealthy') + ` (${result.status})` }
  catch (err) { connection.value = hallError(err, t) } finally { checking.value = false }
}
const userRows = ref<AdminUser[]>([])
const selectedUser = ref<AdminUser | null>(null)
const usersLoading = ref(false)
let userRequest: AbortController | undefined
let disposed = false
const baseline = computed(() => JSON.stringify({ ...props.config, expected_nodes: lines(props.config.expected_nodes.join('\n')) }))
const dirty = computed(() => JSON.stringify({ ...draft.value, expected_nodes: lines(nodeText.value) }) !== baseline.value)
watch(dirty, value => emit('dirty', value))
watch(() => props.config, config => {
  draft.value = { ...config }
  nodeText.value = config.expected_nodes.join('\n')
  error.value = ''
  void loadSelectedUser()
})
watch(() => draft.value.operator_user_id, () => { void loadSelectedUser() })

const switches = computed(() => {
  const ready = props.config.readiness
  return ([
    { key: 'collection', field: 'collection_enabled', ready: ready?.collection === true },
    { key: 'tasks', field: 'tasks_enabled', ready: ready?.tasks === true },
    { key: 'display', field: 'display_enabled', ready: ready?.display === true },
  ] as const)
})
const profileOptions = computed(() => [{ value: 0, label: t('admin.providerHall.none') },
  ...props.profiles.map(p => ({ value: p.id, label: `${p.model} / ${p.protocol}` }))])
const defaultProfileID = computed(() => props.profiles.find(p => p.model === draft.value.default_model && p.protocol === draft.value.default_protocol)?.id ?? 0)
function selectProfile(value: unknown) {
  const profile = props.profiles.find(p => p.id === value)
  draft.value.default_model = profile?.model ?? ''
  if (profile) draft.value.default_protocol = profile.protocol
}
const operatorOptions = computed(() => {
  const options = userRows.value.map(u => ({ value: u.id, label: `${u.email} (#${u.id})` }))
  const id = draft.value.operator_user_id
  if (id && !options.some(o => o.value === id)) options.unshift({ value: id, label: selectedUser.value?.id === id ? `${selectedUser.value.email} (#${id})` : `#${id}` })
  return [{ value: null, label: t('admin.providerHall.none') }, ...options]
})
async function loadSelectedUser() {
  const id = draft.value.operator_user_id
  if (!id) { selectedUser.value = null; return }
  try {
    const user = await users.getById(id)
    if (!disposed && id === draft.value.operator_user_id) selectedUser.value = user
  } catch { /* The saved ID remains visible even if the user was deleted. */ }
}
async function searchUsers(search = '') {
  userRequest?.abort()
  const request = new AbortController()
  userRequest = request
  usersLoading.value = true
  try {
    const result = await users.list(1, 50, { search, status: 'active' }, { signal: request.signal })
    if (!request.signal.aborted) userRows.value = result.items
  } catch (err) {
    if (!request.signal.aborted) error.value = hallError(err, t, 'loadFailed')
  } finally {
    if (!request.signal.aborted) usersLoading.value = false
  }
}
function reload() {
  if (busy.value) return
  if (!dirty.value || window.confirm(t('admin.providerHall.discard'))) emit('reload')
}
async function save() {
  if (busy.value) return
  saving.value = true
  error.value = ''
  try {
    const d = draft.value
    const input = { version: props.config.version, collection_enabled: d.collection_enabled,
      auto_schedule_enabled: d.auto_schedule_enabled ?? false,
      display_enabled: d.display_enabled, tasks_enabled: d.tasks_enabled, default_model: d.default_model, default_protocol: d.default_protocol,
      default_range: d.default_range, gateway_origin: d.gateway_origin, operator_user_id: d.operator_user_id,
      daily_budget: d.daily_budget, expected_nodes: lines(nodeText.value) }
    const preview = await hall.preflightConfig(input)
    if (preview.disabled_targets.length && !window.confirm(`${t('admin.providerHall.operatorImpact')}\n${preview.disabled_targets.map(target => `#${target.group_id} / ${props.profiles.find(p => p.id === target.profile_id)?.model || target.profile_id}`).join('\n')}`)) return
    const result = await hall.updateConfig(input)
    if (!disposed) emit('saved', result)
  } catch (err) { if (!disposed) error.value = hallError(err, t) }
  finally { saving.value = false }
}
onMounted(() => { void searchUsers(); void loadSelectedUser(); void hall.getHealth().then(h => { if (!disposed) nodes.value = h.collection.nodes.map(n => n.node_id) }).catch(() => {}) })
onUnmounted(() => { disposed = true; userRequest?.abort() })
</script>
