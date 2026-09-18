<template>
  <AppLayout class="provider-hall-layout">
    <div class="provider-hall w-full min-w-0 pb-8">
      <header class="border-b border-gray-200 pb-4 dark:border-dark-700">
        <h1 class="flex items-center gap-2 text-xl font-semibold"><Building2 :size="22" class="text-emerald-600" />{{ t('admin.providerHall.title') }}</h1>
        <div class="mt-5 flex gap-1 overflow-x-auto" role="tablist" :aria-label="t('admin.providerHall.title')">
          <button v-for="item in tabs" :id="`hall-tab-${item}`" :key="item" type="button" role="tab"
            :aria-controls="`hall-panel-${item}`" :aria-selected="tab === item" :tabindex="tab === item ? 0 : -1"
            class="hall-tab" :class="tab === item && 'hall-tab-active'" @click="tab = item" @keydown="tabKey($event, item)">
            {{ t(`admin.providerHall.${item}`) }}
          </button>
        </div>
      </header>
      <p v-if="error" role="alert" class="hall-error mt-4">{{ error }}</p>
      <p v-if="loading && !config" role="status" class="py-10 text-sm text-gray-500">{{ t('common.loading') }}</p>
      <button v-else-if="!config" class="btn btn-secondary mt-4" @click="load"><RefreshCw :size="16" />{{ t('admin.providerHall.reload') }}</button>
      <template v-if="config">
        <section id="hall-panel-config" v-show="tab === 'config'" role="tabpanel" aria-labelledby="hall-tab-config" class="py-6">
          <ProviderHallConfigForm :config="config" :profiles="profiles" :loading="loading" @saved="configSaved" @reload="load" @dirty="dirty.config = $event" />
        </section>
        <section id="hall-panel-groups" v-show="tab === 'groups'" role="tabpanel" aria-labelledby="hall-tab-groups" class="py-6">
          <ProviderHallGroups :profiles="profiles" :operator-id="config.operator_user_id" :focus-group="focusGroup" @open="openJob = $event" @jobs="showJobs" @profile="profileSaved" @dirty="dirty.groups = $event" />
        </section>
        <section id="hall-panel-profiles" v-show="tab === 'profiles'" role="tabpanel" aria-labelledby="hall-tab-profiles" class="py-6">
          <ProviderHallProfiles :profiles="profiles" :config="config" @config="tab = 'config'" @saved="profileSaved" @reload="reloadProfiles" @dirty="dirty.profiles = $event" />
        </section>
        <section id="hall-panel-jobs" v-show="tab === 'jobs'" role="tabpanel" aria-labelledby="hall-tab-jobs" class="py-6">
          <ProviderHallJobs v-if="visited.jobs" :focus-group="jobGroup" :focus-profile="jobProfile" :active="tab === 'jobs'" @open="openJob = $event" @target="showTarget" />
        </section>
        <section id="hall-panel-health" v-show="tab === 'health'" role="tabpanel" aria-labelledby="hall-tab-health" class="py-6">
          <ProviderHallHealth v-if="visited.health" @config="tab = 'config'" />
        </section>
        <ProviderHallJobDetailDialog :job-id="openJob" @close="openJob = null" @target="openJob = null; showTarget($event)" />
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { nextTick, onMounted, onUnmounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router'
import { Building2, RefreshCw } from 'lucide-vue-next'
import AppLayout from '@/components/layout/AppLayout.vue'
import ProviderHallConfigForm from '@/components/admin/provider-hall/ProviderHallConfigForm.vue'
import ProviderHallGroups from '@/components/admin/provider-hall/ProviderHallGroups.vue'
import ProviderHallProfiles from '@/components/admin/provider-hall/ProviderHallProfiles.vue'
import ProviderHallJobs from '@/components/admin/provider-hall/ProviderHallJobs.vue'
import ProviderHallHealth from '@/components/admin/provider-hall/ProviderHallHealth.vue'
import ProviderHallJobDetailDialog from '@/components/admin/provider-hall/ProviderHallJobDetailDialog.vue'
import { hallError } from '@/components/admin/provider-hall/helpers'
import * as hall from '@/api/admin/providerHall'
import { useAppStore } from '@/stores/app'

const { t } = useI18n()
const app = useAppStore()
const route = useRoute(), router = useRouter()
const tabs = ['groups', 'profiles', 'config', 'jobs', 'health'] as const
type Tab = typeof tabs[number]
const tab = ref<Tab>(tabs.includes(route.query.tab as Tab) ? route.query.tab as Tab : 'groups')
// Task and health panels load on first visit so the config page stays cheap.
const visited = reactive({ jobs: false, health: false })
const openJob = ref<number | null>(Number(route.query.job) || null)
const focusGroup = ref<number | null>(tab.value === 'groups' ? Number(route.query.group) || null : null), jobGroup = ref<number | null>(Number(route.query.group) || null), jobProfile = ref<number | null>(Number(route.query.profile) || null)
function showJobs(group: number, profile: number) { jobGroup.value = group; jobProfile.value = profile; tab.value = 'jobs' }
async function showTarget(group: number) { focusGroup.value = null; await nextTick(); focusGroup.value = group; tab.value = 'groups' }
// Keep one route writer so tab changes cannot overwrite target/job navigation.
watch([tab, jobGroup, jobProfile, focusGroup, openJob], ([active, group, profile, focus, job]) => {
  void router.replace({ query: { ...route.query, tab: active,
    group: (active === 'jobs' ? group : active === 'groups' ? focus : null) || undefined,
    profile: active === 'jobs' ? profile || undefined : undefined,
    job: job || undefined } })
})
watch(tab, value => { if (value === 'jobs') visited.jobs = true; if (value === 'health') visited.health = true }, { immediate: true })
const config = ref<hall.ProviderHallConfig | null>(null)
const profiles = ref<hall.ProviderHallProfile[]>([])
const loading = ref(false)
const error = ref('')
const dirty = reactive({ config: false, groups: false, profiles: false })
let request: AbortController | undefined
let profileRequest: AbortController | undefined

async function load() {
  request?.abort()
  const current = new AbortController()
  request = current
  loading.value = true
  error.value = ''
  try {
    const cfg = await hall.getConfig(current.signal)
    if (!current.signal.aborted) config.value = cfg
  } catch (err) { if (!current.signal.aborted) error.value = hallError(err, t, 'loadFailed') }
  finally { if (!current.signal.aborted) loading.value = false }
}
async function reloadProfiles() {
  profileRequest?.abort()
  const current = new AbortController()
  profileRequest = current
  try {
    const rows = await hall.listProfiles(current.signal)
    if (!current.signal.aborted) profiles.value = rows
  } catch (err) { if (!current.signal.aborted) error.value = hallError(err, t, 'loadFailed') }
}
function configSaved(value: hall.ProviderHallConfig) {
  request?.abort()
  loading.value = false
  config.value = value
  app.showSuccess(t('admin.providerHall.saved'))
}
function profileSaved(value: hall.ProviderHallProfile) {
  profileRequest?.abort()
  profiles.value = [...profiles.value.filter(p => p.id !== value.id), value].sort((a, b) => a.id - b.id)
  app.showSuccess(t('admin.providerHall.saved'))
}
async function tabKey(event: KeyboardEvent, item: Tab) {
  const index = tabs.indexOf(item)
  if (event.key !== 'ArrowRight' && event.key !== 'ArrowLeft' && event.key !== 'Home' && event.key !== 'End') return
  event.preventDefault()
  const next = event.key === 'Home' ? 0 : event.key === 'End' ? tabs.length - 1 : (index + (event.key === 'ArrowRight' ? 1 : -1) + tabs.length) % tabs.length
  tab.value = tabs[next]
  await nextTick()
  document.getElementById(`hall-tab-${tab.value}`)?.focus()
}
function beforeUnload(event: BeforeUnloadEvent) {
  if (Object.values(dirty).some(Boolean)) { event.preventDefault(); event.returnValue = '' }
}
onBeforeRouteLeave(() => !Object.values(dirty).some(Boolean) || window.confirm(t('admin.providerHall.discard')))
onMounted(() => { void load(); void reloadProfiles(); window.addEventListener('beforeunload', beforeUnload) })
onUnmounted(() => { request?.abort(); profileRequest?.abort(); window.removeEventListener('beforeunload', beforeUnload) })
</script>

<style>
.provider-hall-layout main { @apply bg-white dark:bg-zinc-900; min-height: calc(100vh - 64px); }
.provider-hall { @apply text-gray-900 dark:text-gray-100; letter-spacing: 0; }
.provider-hall .hall-tab { @apply shrink-0 border-b-2 border-transparent px-4 py-2 text-sm text-gray-500 hover:text-gray-900 dark:hover:text-white; }
.provider-hall .hall-tab-active { @apply border-emerald-600 font-medium text-emerald-700 dark:text-emerald-400; }
.hall-form { @apply max-w-4xl space-y-5; }
.hall-fields { @apply grid min-w-0 gap-5 sm:grid-cols-2; }
.hall-field { @apply flex min-w-0 flex-col gap-2 text-sm; }
.hall-field .input { @apply w-full min-w-0; }
.hall-actions { @apply flex flex-wrap items-center justify-end gap-3 border-t border-gray-200 pt-4 dark:border-dark-700; }
.hall-actions .btn, .provider-hall .btn { @apply inline-flex items-center justify-center gap-2; }
.hall-error { @apply break-words border-l-2 border-red-500 bg-red-50 p-3 text-sm text-red-700 dark:bg-red-950/20 dark:text-red-300; }
.hall-table { @apply w-full table-fixed text-left text-sm; }
.hall-table th { @apply relative border-b border-gray-200 px-3 py-3 font-medium text-gray-500 dark:border-dark-700; }
.hall-table td { @apply border-b border-gray-100 px-3 py-4 align-top dark:border-dark-700; }
.hall-table .input { @apply w-full min-w-0; }
</style>
