<template>
  <div class="hall-verification inline-flex flex-col items-center gap-1">
    <span
      class="badge whitespace-nowrap"
      :class="pillClass"
      :data-verdict="kind"
      :title="title"
      data-testid="verification-pill"
    >{{ label }}</span>
    <button
      v-if="verification?.report_id"
      type="button"
      class="text-[11px] leading-none text-primary-600 underline-offset-2 hover:underline dark:text-primary-400"
      data-testid="verification-report-link"
      @click.stop="emit('view-report', verification.report_id as number)"
    >{{ t('providerHall.verification.viewReport') }}</button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { ProviderHallVerificationBadge } from '@/types/providerHall'
import { formatDateTime } from '@/utils/providerHallFormat'

const props = defineProps<{ verification: ProviderHallVerificationBadge | null | undefined }>()
const emit = defineEmits<{ (e: 'view-report', reportId: number): void }>()
const { t, locale } = useI18n()

type Kind = 'passed' | 'failed' | 'suspected' | 'insufficient' | 'expired' | 'none'

const kind = computed<Kind>(() => {
  const v = props.verification
  if (!v || !v.verdict) return 'none'
  if (v.expired) return 'expired'
  return v.verdict
})

const label = computed(() => t(`providerHall.verification.${kind.value}`))

const pillClass = computed(() => {
  switch (kind.value) {
    case 'passed':
      return 'badge-success'
    case 'failed':
      return 'badge-danger'
    case 'suspected':
      return 'badge-warning'
    case 'insufficient':
      return 'badge-warning'
    default:
      return 'badge-gray'
  }
})

const title = computed(() => {
  const v = props.verification
  if (!v || !v.verdict) return label.value
  const parts = [
    t('providerHall.tooltip.verificationMethod', { method: t('providerHall.tooltip.verificationMethodValue') }),
    t('providerHall.tooltip.verificationCompleted', { time: formatDateTime(v.completed_at, locale.value) }),
  ]
  if (v.reason_code) parts.push(t('providerHall.tooltip.verificationReason', { reason: v.reason_code }))
  return parts.join('\n')
})
</script>
