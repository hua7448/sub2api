<template>
  <div class="hall-health space-y-6">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <p class="text-xs text-gray-500">{{ health ? t('admin.providerHall.healthGeneratedAt', { at: formatDate(health.generated_at) }) : '' }}</p>
      <div class="flex items-center gap-2">
        <AutoRefreshButton :enabled="autoRefresh.enabled.value" :interval-seconds="autoRefresh.intervalSeconds.value" :countdown="autoRefresh.countdown.value" :intervals="autoRefresh.intervals" @update:enabled="autoRefresh.setEnabled" @update:interval="autoRefresh.setInterval" />
        <button type="button" class="btn btn-secondary btn-icon" :disabled="loading" :title="t('admin.providerHall.healthRefresh')" :aria-label="t('admin.providerHall.healthRefresh')" @click="load"><RefreshCw :size="16" /></button>
      </div>
    </div>
    <p v-if="error" role="alert" class="hall-error">{{ error }}</p>
    <p v-if="loading && !health" role="status" class="py-6 text-sm text-gray-500">{{ t('common.loading') }}</p>
    <template v-if="health">
      <div class="flex flex-wrap items-center gap-4 border-b border-gray-200 pb-4 text-sm dark:border-dark-700">
        <span>{{ t('admin.providerHall.tasks') }}: {{ t(health.tasks_enabled ? 'admin.providerHall.enabled' : 'admin.providerHall.inactive') }}</span>
        <span>{{ t('admin.providerHall.autoSchedule') }}: {{ t(health.auto_schedule_enabled ? 'admin.providerHall.enabled' : 'admin.providerHall.inactive') }}</span>
        <button class="btn btn-secondary" @click="emit('config')">{{ t('admin.providerHall.settingsEntry') }}</button>
      </div>
      <p v-if="health.budget.paused_reason || health.aggregator.last_error" class="hall-error">{{ health.budget.paused_reason ? t(`admin.providerHall.paused_${health.budget.paused_reason}`) : health.aggregator.last_error }}</p>
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <StatCard :title="t('admin.providerHall.healthPending')" :value="health.reconciliation.pending" :icon="Hourglass" icon-variant="primary" />
        <StatCard :title="t('admin.providerHall.healthUncertain')" :value="health.reconciliation.uncertain" :icon="CircleHelp" :icon-variant="health.reconciliation.uncertain > 0 ? 'warning' : 'primary'" />
        <StatCard :title="t('admin.providerHall.healthConfirmedSpend')" :value="`${health.budget.confirmed_spend} / ${health.budget.budget}`" :icon="Wallet" :icon-variant="health.budget.paused_reason === 'budget_exhausted' ? 'danger' : 'success'" />
        <StatCard :title="t('admin.providerHall.healthJobs')" :value="`${health.jobs.queued} / ${health.jobs.running} / ${health.jobs.unknown}`" :icon="Activity" :icon-variant="health.jobs.unknown > 0 ? 'warning' : 'primary'" />
      </div>
      <section class="space-y-2">
        <h2 class="text-base font-semibold">{{ t('admin.providerHall.healthNodes') }}</h2>
        <p v-if="health.collection.missing_expected.length" role="alert" class="hall-error">{{ t('admin.providerHall.healthMissing') }}: {{ health.collection.missing_expected.join(', ') }}</p>
        <div class="overflow-x-auto">
          <table class="hall-table min-w-[820px]">
            <thead><tr><th class="w-[20%]">{{ t('admin.providerHall.healthNode') }}</th><th class="w-[8%]">{{ t('admin.providerHall.healthEpoch') }}</th><th class="w-[12%]">{{ t('admin.providerHall.healthVersion') }}</th><th class="w-[18%]">{{ t('admin.providerHall.healthHeartbeat') }}</th><th class="w-[18%]">{{ t('admin.providerHall.healthConfirmed') }}</th><th class="w-[10%]">{{ t('admin.providerHall.healthSeq') }}</th><th class="w-[14%]">{{ t('admin.providerHall.jobStatus') }}</th></tr></thead>
            <tbody>
              <tr v-for="node in health.collection.nodes" :key="node.node_id">
                <td class="break-all font-mono text-xs">{{ node.node_id }}</td>
                <td class="text-xs">{{ node.epoch_id }}</td>
                <td class="break-all text-xs">{{ node.version || '—' }}</td>
                <td class="text-xs">{{ formatDate(node.heartbeat_at) }}</td>
                <td class="text-xs">{{ node.confirmed_at ? formatDate(node.confirmed_at) : '—' }}</td>
                <td class="text-xs">{{ node.persisted_seq }}</td>
                <td><span class="hall-pill" :class="node.lost ? 'hall-pill-failed' : 'hall-pill-succeeded'">{{ t(node.lost ? 'admin.providerHall.healthLost' : 'admin.providerHall.healthAlive') }}</span><span v-if="node.overflowed" class="hall-pill hall-pill-unknown ml-1">{{ t('admin.providerHall.healthOverflowed') }}</span></td>
              </tr>
              <tr v-if="!health.collection.nodes.length"><td colspan="7" class="py-6 text-center text-gray-500">{{ t('admin.providerHall.healthNoNodes') }}</td></tr>
            </tbody>
          </table>
        </div>
      </section>
      <section class="space-y-2">
        <h2 class="text-base font-semibold">{{ t('admin.providerHall.healthGaps') }}</h2>
        <ul v-if="health.collection.open_gaps.length" class="space-y-1 text-sm">
          <li v-for="gap in health.collection.open_gaps" :key="gap.id" class="flex flex-wrap gap-2"><span class="font-mono text-xs">{{ gap.node_id }}</span><span>{{ gap.scope }} · {{ gap.reason }}</span><span class="text-gray-500">{{ formatDate(gap.started_at) }}</span></li>
        </ul>
        <p v-else class="text-sm text-gray-500">{{ t('admin.providerHall.healthNoGaps') }}</p>
      </section>
      <details><summary class="cursor-pointer py-3 text-sm font-medium">{{ t('admin.providerHall.diagnostics') }}</summary><div class="mb-6 grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <StatCard :title="t('admin.providerHall.healthCollectionEnabled')" :value="health.collection.enabled ? t('admin.providerHall.enabled') : t('admin.providerHall.inactive')" :icon="Radio" :icon-variant="health.collection.enabled ? 'success' : 'warning'" />
        <StatCard :title="t('admin.providerHall.healthQueue')" :value="t('admin.providerHall.healthQueueValue', health.collection.local_queue)" :icon="Layers" :icon-variant="health.collection.local_queue.dropped > 0 ? 'danger' : 'primary'" />
        <StatCard :title="t('admin.providerHall.healthLag')" :value="health.aggregator.lag_seconds === null ? t('admin.providerHall.never') : t('admin.providerHall.seconds', { n: health.aggregator.lag_seconds })" :icon="Timer" :icon-variant="health.aggregator.lag_seconds !== null && health.aggregator.lag_seconds > 300 ? 'warning' : 'primary'" />
        <StatCard :title="t('admin.providerHall.healthDirty')" :value="health.aggregator.dirty_count" :icon="ListChecks" icon-variant="primary" />
      </div><div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <section class="space-y-2">
          <h2 class="text-base font-semibold">{{ t('admin.providerHall.healthAggregator') }}</h2>
          <dl class="hall-dl">
            <dt>{{ t('admin.providerHall.healthWatermark') }}</dt><dd>{{ health.aggregator.watermark ? formatDate(health.aggregator.watermark) : t('admin.providerHall.never') }}</dd>
            <dt>{{ t('admin.providerHall.healthLastRun') }}</dt><dd>{{ health.aggregator.last_run_at ? formatDate(health.aggregator.last_run_at) : t('admin.providerHall.never') }}</dd>
            <dt>{{ t('admin.providerHall.healthLastError') }}</dt><dd class="break-all">{{ health.aggregator.last_error || '—' }}</dd>
          </dl>
        </section>
        <section class="space-y-2">
          <h2 class="text-base font-semibold">{{ t('admin.providerHall.healthBudget') }}</h2>
          <dl class="hall-dl">
            <dt>{{ t('admin.providerHall.healthBudgetDay') }}</dt><dd>{{ health.budget.day }}</dd>
            <dt>{{ t('admin.providerHall.healthUncertainSpend') }}</dt><dd>{{ health.budget.uncertain_spend }}</dd>
            <dt>{{ t('admin.providerHall.healthInFlight') }}</dt><dd>{{ health.budget.in_flight }}</dd>
            <dt>{{ t('admin.providerHall.healthPaused') }}</dt><dd>{{ health.budget.paused_reason ? t(`admin.providerHall.paused_${health.budget.paused_reason}`) : '—' }}</dd>
          </dl>
        </section>
        <section class="space-y-2">
          <h2 class="text-base font-semibold">{{ t('admin.providerHall.healthReconciliation') }}</h2>
          <dl class="hall-dl">
            <dt>{{ t('admin.providerHall.healthFailed24h') }}</dt><dd>{{ health.reconciliation.failed_24h }}</dd>
            <dt>{{ t('admin.providerHall.healthJobsFailed24h') }}</dt><dd>{{ health.jobs.failed_24h }}</dd>
          </dl>
        </section>
      </div></details>
    </template>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Activity, CircleHelp, Hourglass, Layers, ListChecks, Radio, RefreshCw, Timer, Wallet } from 'lucide-vue-next'
import StatCard from '@/components/common/StatCard.vue'
import AutoRefreshButton from '@/components/common/AutoRefreshButton.vue'
import { useAutoRefresh } from '@/composables/useAutoRefresh'
import * as hall from '@/api/admin/providerHall'
import { formatDate } from '@/utils/format'
import { hallError } from './helpers'

const { t } = useI18n()
const emit = defineEmits<{ config: [] }>()
const health = ref<hall.ProviderHallHealth | null>(null)
const loading = ref(false)
const error = ref('')
let request: AbortController | undefined

async function load() {
  request?.abort()
  const current = new AbortController()
  request = current
  loading.value = true
  error.value = ''
  try {
    const result = await hall.getHealth(current.signal)
    if (!current.signal.aborted) health.value = result
  } catch (err) { if (!current.signal.aborted) error.value = hallError(err, t, 'loadFailed') }
  finally { if (!current.signal.aborted) loading.value = false }
}
const autoRefresh = useAutoRefresh({ storageKey: 'provider_hall_health_auto_refresh', intervals: [15, 30, 60] as const, defaultInterval: 30, onRefresh: load, shouldPause: () => document.hidden })
defineExpose({ load })
onMounted(() => { void load() })
onUnmounted(() => request?.abort())
</script>

<style>
.hall-health .stat-value { white-space: normal; overflow-wrap: anywhere; font-size: 1.125rem; }
.hall-dl { @apply grid grid-cols-[auto_1fr] gap-x-4 gap-y-1 text-sm; }
.hall-dl dt { @apply text-gray-500 dark:text-gray-400; }
</style>
