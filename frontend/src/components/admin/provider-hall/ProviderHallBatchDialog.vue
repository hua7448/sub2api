<template>
  <BaseDialog :show="show" :title="t('admin.providerHall.batch')" width="wide" :close-on-escape="!busy" :show-close-button="!busy" @close="emit('close')">
    <div class="space-y-4">
      <select v-model="action" class="input" :disabled="busy" :aria-label="t('admin.providerHall.jobActions')"><option v-for="a in actions" :key="a" :value="a">{{ t(`admin.providerHall.batch_${a}`) }}</option></select>
      <label v-if="action === 'intervals'" class="hall-field">{{ t('admin.providerHall.probeMinutes') }}<input v-model.number="minutes" class="input" type="number" min="1" max="1440" /></label>
      <label v-if="action === 'intervals'" class="hall-field">{{ t('admin.providerHall.verifyHours') }}<input v-model.number="hours" class="input" type="number" min="1" max="168" /></label>
      <div v-for="g in selection" :key="g.group_id" class="border-b border-gray-200 py-2 dark:border-dark-700"><p class="font-medium">{{ g.name }}</p><label v-for="target in g.targets" :key="target.id" class="mt-2 flex items-center gap-2 text-sm"><input v-model="targets" type="checkbox" :value="target.id" :disabled="busy" />{{ profiles.find(p => p.id === target.profile_id)?.model }}</label></div>
      <button v-if="executed && selection.length" class="btn btn-secondary" :disabled="busy" @click="reloadFailed">{{ t('admin.providerHall.reloadFailed') }}</button>
      <p v-if="error" class="hall-error">{{ error }}</p>
      <ul v-if="results.length" class="space-y-2 text-sm"><li v-for="r in results" :key="r.id" :class="!r.success && 'text-red-600'">{{ groups.find(g => g.group_id === r.id)?.name }}: {{ r.success ? t(executed ? 'admin.providerHall.saved' : 'admin.providerHall.batchReady') : hallError({ reason: r.reason, message: r.error, metadata: r.metadata }, t) }}</li></ul>
    </div>
    <template #footer><div class="flex flex-wrap gap-3"><button class="btn btn-secondary" :disabled="busy" @click="run(true)">{{ t('admin.providerHall.preview') }}</button><button class="btn btn-primary" :disabled="busy || !previewed" @click="run(false)">{{ t('admin.providerHall.execute') }}</button></div></template>
  </BaseDialog>
</template>
<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import * as hall from '@/api/admin/providerHall'
import { hallError } from './helpers'
const props = defineProps<{ show: boolean; groups: hall.ProviderHallGroupSummary[]; profiles: hall.ProviderHallProfile[] }>()
const emit = defineEmits<{ close: []; saved: [number[]] }>()
const { t } = useI18n()
const action = ref('list'), targets = ref<number[]>([]), minutes = ref(5), hours = ref(24), busy = ref(false), previewed = ref(false), executed = ref(false), error = ref('')
const actions = ['list', 'unlist', 'enable', 'disable', 'auto_on', 'auto_off', 'intervals']
const selection = ref<hall.ProviderHallGroupSummary[]>([]), results = ref<hall.ProviderHallBatchResult[]>([])
watch(() => props.show, show => { if (show) { selection.value = props.groups; targets.value = []; results.value = []; previewed.value = false; executed.value = false; error.value = '' } })
watch([action, targets, minutes, hours], () => { previewed.value = false }, { deep: true })
async function reloadFailed() {
  if (busy.value) return
  busy.value = true; error.value = ''; previewed.value = false
  try {
    const updated = await Promise.all(selection.value.map(async g => {
      const settings = await hall.getTargets(g.group_id)
      return { ...g, ...settings.listing, version: settings.version, targets: settings.items }
    }))
    selection.value = updated
    targets.value = targets.value.filter(id => updated.some(g => g.targets.some(target => target.id === id)))
    results.value = []
  } catch (err) { error.value = hallError(err, t) } finally { busy.value = false }
}
async function run(preview: boolean) {
  if (busy.value || !selection.value.length) return
  if (!['list', 'unlist'].includes(action.value) && !targets.value.length) { error.value = t('admin.providerHall.selectTargets'); return }
  busy.value = true; error.value = ''
  try {
    const input = selection.value.map(g => ({ id: g.group_id, version: g.version, listed: action.value === 'list' ? true : action.value === 'unlist' ? false : g.listed, display_name: g.display_name, description: g.description, display_order: g.display_order, items: g.targets.map(target => {
      const row: hall.ProviderHallTargetInput = { profile_id: target.profile_id, probe_key_id: target.probe_key_id, enabled: target.enabled, auto_schedule_enabled: target.auto_schedule_enabled, probe_interval_seconds: target.probe_interval_seconds, verification_interval_seconds: target.verification_interval_seconds }
      if (targets.value.includes(target.id)) { if (action.value === 'enable') row.enabled = true; if (action.value === 'disable') row.enabled = false; if (action.value === 'auto_on') row.auto_schedule_enabled = true; if (action.value === 'auto_off') row.auto_schedule_enabled = false; if (action.value === 'intervals') { row.probe_interval_seconds = Math.round(minutes.value * 60); row.verification_interval_seconds = Math.round(hours.value * 3600) } } return row
    }) }))
    results.value = await hall.batchGroups(input, preview); previewed.value = preview && results.value.some(r => r.success); executed.value = !preview
    if (!preview) { const failed = results.value.filter(r => !r.success).map(r => r.id); selection.value = selection.value.filter(g => failed.includes(g.group_id)); emit('saved', failed) }
  } catch (err) { error.value = hallError(err, t) } finally { busy.value = false }
}
</script>
