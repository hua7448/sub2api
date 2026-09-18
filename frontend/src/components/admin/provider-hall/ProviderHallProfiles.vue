<template>
  <div>
    <div class="mb-4 flex items-center justify-between gap-3">
      <h2 class="text-base font-semibold">{{ t('admin.providerHall.profiles') }}</h2>
      <div class="flex gap-2">
        <button class="btn btn-secondary btn-icon" :title="t('admin.providerHall.reload')" :aria-label="t('admin.providerHall.reload')" @click="emit('reload')"><RefreshCw :size="16" /></button>
        <button class="btn btn-primary" @click="open()"><Plus :size="16" />{{ t('admin.providerHall.addProfile') }}</button>
      </div>
    </div>
    <div class="overflow-x-auto">
      <table class="hall-table min-w-[650px]">
        <thead><tr><th class="w-[32%]">{{ t('admin.providerHall.model') }}</th><th class="w-[22%]">{{ t('admin.providerHall.protocol') }}</th><th class="w-[16%]">{{ t('admin.providerHall.outputLimit') }}</th><th class="w-[20%]">{{ t('admin.providerHall.tools') }}</th><th class="w-[10%]"><span class="sr-only">{{ t('common.edit') }}</span></th></tr></thead>
        <tbody>
          <tr v-for="profile in profiles" :key="profile.id">
            <td class="break-all font-medium">{{ profile.model }}</td>
            <td>{{ protocols.find(p => p.value === profile.protocol)?.label }}</td>
            <td>{{ profile.output_limit }}</td>
            <td><Check v-if="profile.supports_tools" :size="17" class="text-emerald-600" :aria-label="t('common.yes')" /><span v-else class="text-gray-500">{{ t('common.no') }}</span></td>
            <td><button class="btn btn-ghost btn-icon" :title="t('admin.providerHall.editProfile')" :aria-label="`${t('common.edit')} ${profile.model}`" @click="open(profile)"><Pencil :size="16" /></button></td>
          </tr>
          <tr v-if="!profiles.length"><td colspan="5" class="py-12 text-center text-gray-500">{{ t('admin.providerHall.noProfiles') }}</td></tr>
        </tbody>
      </table>
    </div>
    <BaseDialog :show="show" :title="t(editingID ? 'admin.providerHall.editProfile' : 'admin.providerHall.addProfile')" width="wide" :close-on-escape="!saving" :show-close-button="!saving" @close="close">
      <form id="hall-profile-form" class="hall-form" @submit.prevent="save">
        <p v-if="error" role="alert" class="hall-error">{{ error }}</p>
        <fieldset :disabled="saving" class="hall-fields">
          <label class="hall-field"><span>{{ t('admin.providerHall.model') }}</span><input v-model.trim="draft.model" class="input" maxlength="200" required /></label>
          <label class="hall-field"><span>{{ t('admin.providerHall.protocol') }}</span><select v-model="draft.protocol" class="input"><option v-for="p in protocols" :key="p.value" :value="p.value">{{ p.label }}</option></select></label>
          <label class="hall-field"><span>{{ t('admin.providerHall.outputLimit') }}</span><input v-model.number="draft.output_limit" class="input" type="number" min="1" max="1024" step="1" required /></label>
          <label class="flex items-center gap-2 text-sm"><input v-model="draft.supports_tools" type="checkbox" class="h-4 w-4" />{{ t('admin.providerHall.tools') }}</label>
          <label class="hall-field sm:col-span-2"><span>{{ t('admin.providerHall.aliases') }}</span><textarea v-model="aliasText" class="input font-mono" rows="3" /></label>
          <label class="flex items-center gap-2 text-sm sm:col-span-2"><input v-model="hasReference" type="checkbox" class="h-4 w-4" />{{ t('admin.providerHall.reference') }}</label>
          <template v-if="hasReference">
            <label class="hall-field"><span>{{ t('admin.providerHall.inputPrice') }}</span><input v-model="draft.reference_input_price" class="input" inputmode="decimal" pattern="[0-9]+(\.[0-9]{1,10})?" required @input="referenceEdited" /></label>
            <label class="hall-field"><span>{{ t('admin.providerHall.cachePrice') }}</span><input v-model="draft.reference_cache_price" class="input" inputmode="decimal" pattern="[0-9]+(\.[0-9]{1,10})?" required @input="referenceEdited" /></label>
            <label class="hall-field"><span>{{ t('admin.providerHall.cacheRate') }}</span><input v-model="draft.reference_cache_rate" class="input" inputmode="decimal" pattern="(0(\.[0-9]{1,10})?|1(\.0{1,10})?)" required @input="referenceEdited" /></label>
            <div class="hall-field">
              <span>{{ t('admin.providerHall.confirmedAt') }}</span>
              <div class="flex flex-wrap items-center gap-2"><time class="text-xs">{{ draft.reference_confirmed_at ? new Date(draft.reference_confirmed_at).toLocaleString(locale) : t('admin.providerHall.none') }}</time><button type="button" class="btn btn-secondary" @click="draft.reference_confirmed_at = new Date().toISOString()"><Check :size="16" />{{ t('admin.providerHall.confirmNow') }}</button></div>
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
  </div>
</template>

<script setup lang="ts">
import { computed, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { Check, Pencil, Plus, RefreshCw, Save } from 'lucide-vue-next'
import BaseDialog from '@/components/common/BaseDialog.vue'
import * as hall from '@/api/admin/providerHall'
import { hallError, lines, protocols } from './helpers'

defineProps<{ profiles: hall.ProviderHallProfile[] }>()
const emit = defineEmits<{ saved: [hall.ProviderHallProfile]; reload: []; dirty: [boolean] }>()
const { t, locale } = useI18n()
const show = ref(false)
const saving = ref(false)
const editingID = ref<number | null>(null)
const error = ref('')
const aliasText = ref('')
const hasReference = ref(false)
const draft = ref<hall.ProviderHallProfileInput>(emptyProfile())
let disposed = false
const snapshot = computed(() => JSON.stringify({ ...draft.value, model_aliases: lines(aliasText.value), hasReference: hasReference.value }))
const initial = ref('')
const dirty = computed(() => show.value && snapshot.value !== initial.value)
watch(dirty, value => emit('dirty', value))
function emptyProfile(): hall.ProviderHallProfileInput {
  return { version: 0, model: '', protocol: 'responses', supports_tools: false, output_limit: 256, model_aliases: [],
    reference_input_price: null, reference_cache_price: null, reference_cache_rate: null, reference_confirmed_at: null }
}
function open(profile?: hall.ProviderHallProfile) {
  draft.value = profile ? { ...profile, model_aliases: [...profile.model_aliases] } : emptyProfile()
  editingID.value = profile?.id ?? null
  aliasText.value = draft.value.model_aliases.join('\n')
  hasReference.value = draft.value.reference_input_price !== null
  initial.value = snapshot.value
  error.value = ''
  show.value = true
}
function referenceEdited() { draft.value.reference_confirmed_at = null }
function close() {
  if (!saving.value && (!dirty.value || window.confirm(t('admin.providerHall.discard')))) show.value = false
}
async function save() {
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
    if (!disposed) { show.value = false; emit('saved', result) }
  } catch (err) { if (!disposed) error.value = hallError(err, t) }
  finally { saving.value = false }
}
onUnmounted(() => { disposed = true })
</script>

<style>
.modal-content:has(#hall-profile-form) { @apply rounded-lg bg-white dark:bg-zinc-900; backdrop-filter: none; }
</style>
