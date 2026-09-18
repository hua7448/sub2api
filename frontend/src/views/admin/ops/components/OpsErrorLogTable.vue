<template>
  <div
    class="flex h-full min-h-0 flex-col overflow-hidden"
    :class="flat ? 'bg-cyber-surface' : 'card'"
  >
    <IpGeoBatchToolbar :ips="rows.map((row) => row.client_ip)" @failed="emit('ipGeoBatchFailed')" />

    <DataTable
      class="min-h-0 flex-1"
      :columns="columns"
      :data="rows"
      :loading="loading"
      clickable-rows
      server-side-sort
      default-sort-key="created_at"
      default-sort-order="desc"
      @sort="onSort"
      @rowClick="(row) => emit('openErrorDetail', row.id)"
    >
      <template #cell-created_at="{ row }">
        <span class="text-sm text-txt-secondary" :title="row.request_id || row.client_request_id">
          {{ formatDateTime(row.created_at) }}
        </span>
      </template>

      <template #cell-type="{ row }">
        <span class="inline-flex items-center rounded px-2 py-0.5 text-xs font-medium" :class="getTypeBadge(row).className">
          {{ getTypeBadge(row).label }}
        </span>
      </template>

      <template #cell-endpoint="{ row }">
        <div class="max-w-[320px] space-y-1 text-xs">
          <div class="break-all text-txt-primary">
            <span class="font-medium text-txt-secondary">{{ t('usage.inbound') }}:</span>
            <span class="ml-1">{{ row.inbound_endpoint?.trim() || '-' }}</span>
          </div>
          <div v-if="row.upstream_endpoint" class="break-all text-txt-primary">
            <span class="font-medium text-txt-secondary">{{ t('usage.upstream') }}:</span>
            <span class="ml-1">{{ row.upstream_endpoint.trim() || '-' }}</span>
          </div>
        </div>
      </template>

      <template #cell-platform="{ row }">
        <span class="text-sm text-txt-primary">{{ row.platform || '-' }}</span>
      </template>

      <template #cell-model="{ row }">
        <div v-if="hasModelMapping(row)" class="space-y-0.5 text-xs">
          <div class="break-all font-medium text-txt-primary">{{ row.requested_model }}</div>
          <div class="break-all text-txt-secondary"><span class="mr-0.5">↳</span>{{ row.upstream_model }}</div>
        </div>
        <span v-else-if="displayModel(row)" class="text-sm font-medium text-txt-primary">{{ displayModel(row) }}</span>
        <span v-else class="text-sm text-txt-dim">-</span>
      </template>

      <template #cell-group="{ row }">
        <span
          v-if="row.group_id"
          class="inline-flex items-center rounded bg-indigo-900/30 px-2 py-0.5 text-xs font-medium text-indigo-300"
          :title="t('admin.ops.errorLog.id') + ' ' + row.group_id"
        >
          {{ row.group_name || '#' + row.group_id }}
        </span>
        <span v-else class="text-sm text-txt-dim">-</span>
      </template>

      <template #cell-user="{ row }">
        <div v-if="row.user_id" class="text-sm">
          <button
            v-if="userClickable && row.user_email"
            class="font-medium text-neon-cyan underline decoration-dashed underline-offset-2 transition-colors hover:text-primary-300"
            :title="t('admin.usage.clickToViewBalance')"
            @click.stop="emit('userClick', row.user_id, row.user_email)"
          >
            {{ row.user_email }}
          </button>
          <span v-else class="font-medium text-txt-primary">{{ row.user_email || '-' }}</span>
          <span class="ml-1 text-txt-secondary">#{{ row.user_id }}</span>
        </div>
        <div v-else-if="row.deleted_key_owner_user_id" class="text-sm">
          <button
            v-if="userClickable && row.deleted_key_owner_email"
            class="font-medium text-neon-cyan underline decoration-dashed underline-offset-2 transition-colors hover:text-primary-300"
            :title="t('admin.usage.clickToViewBalance')"
            @click.stop="emit('userClick', row.deleted_key_owner_user_id, row.deleted_key_owner_email ?? undefined)"
          >
            {{ row.deleted_key_owner_email }}
          </button>
          <span v-else class="font-medium text-txt-primary">{{ row.deleted_key_owner_email || '-' }}</span>
          <span class="ml-1 text-txt-secondary">#{{ row.deleted_key_owner_user_id }}</span>
        </div>
        <span v-else class="text-sm text-txt-dim">-</span>
      </template>

      <template #cell-api_key="{ row }">
        <div v-if="row.api_key_id || row.api_key_name" class="text-sm">
          <span class="text-txt-primary">{{ row.api_key_name || '#' + row.api_key_id }}</span>
          <span
            v-if="row.api_key_deleted"
            class="ml-1 inline-flex items-center rounded bg-neon-pink/10 px-1 py-px text-[10px] font-medium leading-tight text-neon-pink ring-1 ring-inset ring-neon-pink/30"
          >
            {{ t('admin.ops.errorLog.keyDeletedBadge') }}
          </span>
        </div>
        <span v-else class="text-sm text-txt-dim">-</span>
      </template>

      <template #cell-account="{ row }">
        <span
          v-if="row.account_id"
          class="text-sm text-txt-primary"
          :title="t('admin.ops.errorLog.accountId') + ' ' + row.account_id"
        >
          {{ row.account_name || '#' + row.account_id }}
        </span>
        <span v-else class="text-sm text-txt-dim">-</span>
      </template>

      <template #cell-category="{ row }">
        <span class="text-sm text-txt-primary">
          {{ t('usage.errors.categories.' + mapErrorCategory(row.phase, row.type)) }}
        </span>
      </template>

      <template #cell-status="{ row }">
        <div class="flex items-center gap-1.5">
          <span class="inline-flex items-center rounded px-2 py-0.5 text-xs font-medium" :class="getStatusClass(row.status_code)">
            {{ row.status_code }}
          </span>
          <span v-if="row.severity" :class="['rounded px-1.5 py-0.5 text-[10px] font-medium', getSeverityClass(row.severity)]">
            {{ row.severity }}
          </span>
          <span v-if="row.request_type != null && row.request_type > 0" class="inline-flex items-center rounded bg-cyber-elevated px-2 py-0.5 text-xs font-medium text-txt-primary">
            {{ formatRequestType(row.request_type) }}
          </span>
        </div>
      </template>

      <template #cell-message="{ row }">
        <span v-if="row.message" class="block max-w-[280px] truncate text-sm text-txt-secondary" :title="row.message">
          {{ formatSmartMessage(row.message) || '-' }}
        </span>
        <span v-else class="text-sm text-txt-dim">-</span>
      </template>

      <template #cell-user_agent="{ row }">
        <span v-if="row.user_agent" class="block max-w-[320px] truncate text-sm text-txt-secondary" :title="row.user_agent">
          {{ row.user_agent }}
        </span>
        <span v-else class="text-sm text-txt-dim">-</span>
      </template>

      <template #cell-client_ip="{ row }">
        <div @click.stop>
          <div v-if="row.client_ip">
            <span class="font-mono text-sm text-txt-secondary">{{ row.client_ip }}</span>
            <IpGeoCell :ip="row.client_ip" />
          </div>
          <span v-else class="text-sm text-txt-dim">-</span>
        </div>
      </template>

      <template #cell-actions="{ row }">
        <button
          type="button"
          class="rounded p-1 text-txt-dim transition-colors hover:bg-cyber-elevated hover:text-neon-cyan"
          :title="t('admin.ops.errorLog.details')"
          @click.stop="emit('openErrorDetail', row.id)"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>
        </button>
      </template>

      <template #empty><EmptyState :message="t('admin.ops.errorLog.noErrors')" /></template>
    </DataTable>

    <Pagination
      v-if="total > 0"
      class="flex-shrink-0"
      :total="total"
      :page="page"
      :page-size="pageSize"
      @update:page="emit('update:page', $event)"
      @update:pageSize="emit('update:pageSize', $event)"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import DataTable from '@/components/common/DataTable.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import Pagination from '@/components/common/Pagination.vue'
import IpGeoCell from '@/components/common/IpGeoCell.vue'
import IpGeoBatchToolbar from '@/components/common/IpGeoBatchToolbar.vue'
import type { OpsErrorLog } from '@/api/admin/ops'
import type { Column } from '@/components/common/types'
import { getSeverityClass, formatDateTime } from '../utils/opsFormatters'
import { mapErrorCategory } from '@/utils/errorCategory'
import { mapErrorSortKey, statusCodeBadgeClass } from '@/utils/errorBadges'

interface Props {
  rows: OpsErrorLog[]
  total: number
  loading: boolean
  page: number
  pageSize: number
  userClickable?: boolean
  visibleColumnKeys?: string[]
  flat?: boolean
}

interface Emits {
  (e: 'openErrorDetail', id: number): void
  (e: 'update:page', value: number): void
  (e: 'update:pageSize', value: number): void
  (e: 'ipGeoBatchFailed'): void
  (e: 'sort', sortBy: string, sortOrder: 'asc' | 'desc'): void
  (e: 'userClick', userId: number, email?: string): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()
const { t } = useI18n()

const allColumns = computed<Column[]>(() => [
  { key: 'user', label: t('admin.ops.errorLog.user') },
  { key: 'api_key', label: t('admin.ops.errorLog.apiKey') },
  { key: 'account', label: t('admin.ops.errorLog.account') },
  { key: 'platform', label: t('admin.ops.errorLog.platform') },
  { key: 'model', label: t('admin.ops.errorLog.model'), sortable: true },
  { key: 'endpoint', label: t('admin.ops.errorLog.endpoint') },
  { key: 'group', label: t('admin.ops.errorLog.group') },
  { key: 'type', label: t('admin.ops.errorLog.type') },
  { key: 'category', label: t('usage.errors.category') },
  { key: 'status', label: t('admin.ops.errorLog.status'), sortable: true },
  { key: 'message', label: t('admin.ops.errorLog.message') },
  { key: 'created_at', label: t('admin.ops.errorLog.time'), sortable: true },
  { key: 'user_agent', label: t('usage.userAgent') },
  { key: 'client_ip', label: t('admin.ops.errorLog.ip') },
  { key: 'actions', label: t('admin.ops.errorLog.action') },
])

const columns = computed<Column[]>(() =>
  props.visibleColumnKeys
    ? allColumns.value.filter((column) => props.visibleColumnKeys!.includes(column.key))
    : allColumns.value
)

function isUpstreamRow(log: OpsErrorLog): boolean {
  return String(log.phase || '').toLowerCase() === 'upstream'
    && String(log.error_owner || '').toLowerCase() === 'provider'
}

function hasModelMapping(log: OpsErrorLog): boolean {
  const requested = String(log.requested_model || '').trim()
  const upstream = String(log.upstream_model || '').trim()
  return Boolean(requested && upstream && requested !== upstream)
}

function displayModel(log: OpsErrorLog): string {
  return String(log.upstream_model || log.requested_model || log.model || '').trim()
}

function formatRequestType(type: number | null | undefined): string {
  switch (type) {
    case 1: return t('admin.ops.errorLog.requestTypeSync')
    case 2: return t('admin.ops.errorLog.requestTypeStream')
    case 3: return t('admin.ops.errorLog.requestTypeWs')
    default: return ''
  }
}

function getTypeBadge(log: OpsErrorLog): { label: string; className: string } {
  const phase = String(log.phase || '').toLowerCase()
  const owner = String(log.error_owner || '').toLowerCase()

  if (isUpstreamRow(log)) return { label: t('admin.ops.errorLog.typeUpstream'), className: 'bg-neon-pink/10 text-neon-pink' }
  if (phase === 'request' && owner === 'client') return { label: t('admin.ops.errorLog.typeRequest'), className: 'bg-amber-900/30 text-amber-400' }
  if (phase === 'auth' && owner === 'client') return { label: t('admin.ops.errorLog.typeAuth'), className: 'bg-blue-900/30 text-blue-400' }
  if (phase === 'routing' && owner === 'platform') return { label: t('admin.ops.errorLog.typeRouting'), className: 'bg-purple-900/30 text-purple-400' }
  if (phase === 'internal' && owner === 'platform') return { label: t('admin.ops.errorLog.typeInternal'), className: 'bg-cyber-elevated text-txt-primary' }
  return { label: phase || owner || t('common.unknown'), className: 'bg-cyber-elevated text-txt-secondary' }
}

function onSort(key: string, order: 'asc' | 'desc') {
  emit('sort', mapErrorSortKey(key), order)
}

const getStatusClass = statusCodeBadgeClass

function formatSmartMessage(message: string): string {
  if (!message) return ''
  if (message.startsWith('{') || message.startsWith('[')) {
    try {
      const parsed = JSON.parse(message)
      if (parsed?.error?.message) return String(parsed.error.message)
      if (parsed?.message) return String(parsed.message)
      if (parsed?.detail) return String(parsed.detail)
      if (typeof parsed === 'object') return JSON.stringify(parsed).substring(0, 150)
    } catch {
      // Keep the original non-JSON message below.
    }
  }
  if (message.includes('context deadline exceeded')) return t('admin.ops.errorLog.commonErrors.contextDeadlineExceeded')
  if (message.includes('connection refused')) return t('admin.ops.errorLog.commonErrors.connectionRefused')
  if (message.toLowerCase().includes('rate limit')) return t('admin.ops.errorLog.commonErrors.rateLimit')
  return message.length > 200 ? message.substring(0, 200) + '...' : message
}
</script>
