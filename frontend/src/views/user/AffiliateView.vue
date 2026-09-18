<template>
  <AppLayout>
    <div class="space-y-5 pb-6">
      <div class="flex items-start justify-between gap-3">
        <div class="min-w-0">
          <h1 class="text-xl font-semibold text-gray-900 dark:text-white">{{ t('affiliate.title') }}</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">{{ t('affiliate.description') }}</p>
        </div>
        <button
          class="btn btn-secondary h-9 w-9 p-0"
          :disabled="refreshing"
          :title="t('common.refresh')"
          @click="refreshAll"
        >
          <Icon name="refresh" size="sm" :class="refreshing ? 'animate-spin' : ''" />
        </button>
      </div>

      <div v-if="loading" class="flex justify-center py-16">
        <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
      </div>

      <template v-else-if="detail">
        <div class="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
          <div class="card p-4">
            <div class="flex items-center justify-between">
              <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('affiliate.stats.rebateRate') }}</p>
              <Icon name="dollar" size="sm" class="text-primary-500" />
            </div>
            <p class="mt-2 text-2xl font-semibold text-primary-600 dark:text-primary-400">{{ formattedRebateRate }}%</p>
            <p class="mt-1 text-xs text-gray-400 dark:text-dark-500">{{ t('affiliate.stats.rebateRateHint') }}</p>
          </div>
          <div class="card p-4">
            <div class="flex items-center justify-between">
              <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('affiliate.stats.customers') }}</p>
              <Icon name="users" size="sm" class="text-sky-500" />
            </div>
            <p class="mt-2 text-2xl font-semibold text-gray-900 dark:text-white">
              {{ formatCount(detail.qualified_invitee_count) }}
              <span class="text-sm font-normal text-gray-400">/ {{ formatCount(detail.aff_count) }}</span>
            </p>
            <p class="mt-1 text-xs text-gray-400 dark:text-dark-500">{{ t('affiliate.stats.customersHint') }}</p>
          </div>
          <div class="card p-4">
            <div class="flex items-center justify-between">
              <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('affiliate.stats.availableQuota') }}</p>
              <Icon name="gift" size="sm" class="text-emerald-500" />
            </div>
            <p class="mt-2 text-2xl font-semibold text-emerald-600 dark:text-emerald-400">{{ formatCurrency(detail.aff_quota) }}</p>
            <p class="mt-1 text-xs text-gray-400 dark:text-dark-500">
              {{ t('affiliate.stats.frozenQuota') }} {{ formatCurrency(detail.aff_frozen_quota) }}
            </p>
          </div>
          <div class="card p-4">
            <div class="flex items-center justify-between">
              <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('affiliate.stats.totalQuota') }}</p>
              <Icon name="chart" size="sm" class="text-amber-500" />
            </div>
            <p class="mt-2 text-2xl font-semibold text-gray-900 dark:text-white">{{ formatCurrency(detail.aff_history_quota) }}</p>
            <p class="mt-1 text-xs text-gray-400 dark:text-dark-500">{{ t('affiliate.stats.totalQuotaHint') }}</p>
          </div>
        </div>

        <div
          class="flex flex-col gap-4 border px-4 py-4 sm:flex-row sm:items-center sm:justify-between"
          :class="eligibilityMet
            ? 'border-emerald-200 bg-emerald-50 dark:border-emerald-900/40 dark:bg-emerald-900/20'
            : 'border-amber-200 bg-amber-50 dark:border-amber-900/40 dark:bg-amber-900/20'"
        >
          <div class="flex min-w-0 items-start gap-3">
            <Icon
              :name="eligibilityMet ? 'checkCircle' : 'clock'"
              size="lg"
              class="mt-0.5 shrink-0"
              :class="eligibilityMet ? 'text-emerald-600' : 'text-amber-600'"
            />
            <div>
              <p class="text-sm font-semibold" :class="eligibilityMet ? 'text-emerald-800 dark:text-emerald-200' : 'text-amber-800 dark:text-amber-200'">
                {{ eligibilityMet ? t('affiliate.threshold.met') : t('affiliate.threshold.notMet', { remaining: eligibilityRemaining }) }}
              </p>
              <p class="mt-1 text-xs" :class="eligibilityMet ? 'text-emerald-700 dark:text-emerald-300' : 'text-amber-700 dark:text-amber-300'">
                {{ thresholdSummary }}
              </p>
            </div>
          </div>
          <div v-if="detail.min_qualified_invitees > 0" class="w-full shrink-0 sm:w-64">
            <div class="mb-1 flex justify-between text-xs text-gray-500 dark:text-dark-400">
              <span>{{ t('affiliate.threshold.progress') }}</span>
              <span>{{ eligibilityProgressCurrent }}/{{ detail.min_qualified_invitees }}</span>
            </div>
            <div class="h-2 overflow-hidden rounded-full bg-white/80 dark:bg-dark-800">
              <div class="h-full rounded-full bg-emerald-500 transition-all" :style="{ width: `${eligibilityProgress}%` }"></div>
            </div>
          </div>
        </div>

        <div class="grid gap-4 lg:grid-cols-2">
          <section class="card min-w-0 p-5">
            <div class="flex items-start gap-3">
              <div class="rounded-lg bg-primary-50 p-2 text-primary-600 dark:bg-primary-900/30 dark:text-primary-300">
                <Icon name="userPlus" size="md" />
              </div>
              <div>
                <h2 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('affiliate.share.title') }}</h2>
                <p class="mt-0.5 text-sm text-gray-500 dark:text-dark-400">{{ t('affiliate.share.description') }}</p>
              </div>
            </div>

            <div class="mt-5 grid gap-4 md:grid-cols-2">
              <div class="min-w-0">
                <p class="mb-1.5 text-xs font-medium text-gray-500 dark:text-dark-400">{{ t('affiliate.yourCode') }}</p>
                <div class="flex flex-col items-stretch gap-2 rounded-lg border border-gray-200 bg-gray-50 px-3 py-2 dark:border-dark-700 dark:bg-dark-900 sm:h-11 sm:flex-row sm:items-center sm:py-0">
                  <code class="min-w-0 break-all text-sm font-semibold text-gray-900 dark:text-white sm:flex-1 sm:truncate">{{ detail.aff_code }}</code>
                  <button class="btn btn-secondary btn-sm w-full sm:w-auto sm:shrink-0" @click="copyCode">
                    <Icon name="copy" size="sm" />
                    <span>{{ t('affiliate.copyCode') }}</span>
                  </button>
                </div>
              </div>
              <div class="min-w-0">
                <p class="mb-1.5 text-xs font-medium text-gray-500 dark:text-dark-400">{{ t('affiliate.inviteLink') }}</p>
                <div class="flex flex-col items-stretch gap-2 rounded-lg border border-gray-200 bg-gray-50 px-3 py-2 dark:border-dark-700 dark:bg-dark-900 sm:h-11 sm:flex-row sm:items-center sm:py-0">
                  <code class="min-w-0 break-all text-sm text-gray-700 dark:text-gray-300 sm:flex-1 sm:truncate">{{ inviteLink }}</code>
                  <button class="btn btn-secondary btn-sm w-full sm:w-auto sm:shrink-0" @click="copyInviteLink">
                    <Icon name="link" size="sm" />
                    <span>{{ t('affiliate.copyLink') }}</span>
                  </button>
                </div>
              </div>
            </div>

            <div class="mt-5 border-t border-gray-100 pt-4 dark:border-dark-700">
              <p class="text-sm font-medium text-gray-900 dark:text-white">{{ t('affiliate.tips.title') }}</p>
              <ol class="mt-2 grid gap-x-6 gap-y-2 text-sm text-gray-600 dark:text-dark-300 md:grid-cols-2">
                <li>1. {{ t('affiliate.tips.line1') }}</li>
                <li>2. {{ t('affiliate.tips.line2', { rate: `${formattedRebateRate}%` }) }}</li>
                <li>3. {{ t('affiliate.tips.line3') }}</li>
                <li>
                  4.
                  {{ detail.min_qualified_invitees > 0
                    ? t('affiliate.tips.line4', { count: detail.min_qualified_invitees, start: detail.min_qualified_invitees + 1 })
                    : t('affiliate.tips.line4NoThreshold') }}
                </li>
                <li v-if="detail.aff_frozen_quota > 0">5. {{ t('affiliate.tips.line5') }}</li>
              </ol>
            </div>
          </section>

          <section class="card min-w-0 p-5">
            <div class="flex items-start justify-between gap-3">
              <div>
                <h2 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('affiliate.transfer.title') }}</h2>
                <p class="mt-0.5 text-sm text-gray-500 dark:text-dark-400">{{ t('affiliate.transfer.description') }}</p>
              </div>
              <Icon name="swap" size="md" class="shrink-0 text-emerald-500" />
            </div>
            <div class="mt-5 flex items-end justify-between gap-4 border-b border-gray-100 pb-4 dark:border-dark-700">
              <div>
                <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('affiliate.stats.availableQuota') }}</p>
                <p class="mt-1 text-2xl font-semibold text-emerald-600 dark:text-emerald-400">{{ formatCurrency(detail.aff_quota) }}</p>
              </div>
              <button class="btn btn-primary" :disabled="transferring || detail.aff_quota <= 0 || !eligibilityMet" @click="transferQuota">
                <Icon v-if="transferring" name="refresh" size="sm" class="animate-spin" />
                <Icon v-else name="arrowRight" size="sm" />
                <span>{{ transferring ? t('affiliate.transfer.transferring') : t('affiliate.transfer.button') }}</span>
              </button>
            </div>

            <div class="mt-4">
              <div class="flex items-center justify-between">
                <p class="text-sm font-medium text-gray-900 dark:text-white">{{ t('affiliate.subscriptionDays.title') }}</p>
                <span class="text-xs text-gray-400">{{ t('affiliate.subscriptionDays.allTypes') }}</span>
              </div>
              <div v-if="detail.subscription_days_quota.length" class="mt-2 divide-y divide-gray-100 dark:divide-dark-700">
                <div v-for="quota in detail.subscription_days_quota" :key="quota.group_id" class="flex items-center justify-between gap-3 py-3">
                  <div class="min-w-0">
                    <p class="truncate text-sm font-medium text-gray-900 dark:text-white">{{ quota.group_name || `Group #${quota.group_id}` }}</p>
                    <div class="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-gray-500 dark:text-dark-400">
                      <span>
                        {{ t('affiliate.subscriptionDays.pendingDays') }}
                        <span class="font-semibold text-emerald-600 dark:text-emerald-400">{{ quota.pending_days }}</span>
                      </span>
                      <span v-if="quota.frozen_days > 0" class="text-amber-600 dark:text-amber-400">
                        {{ t('affiliate.subscriptionDays.frozenDays') }} {{ quota.frozen_days }}
                      </span>
                    </div>
                  </div>
                  <button class="btn btn-secondary btn-sm" :disabled="quota.pending_days <= 0 || transferringDays[quota.group_id] || !eligibilityMet" @click="transferDays(quota)">
                    <Icon v-if="transferringDays[quota.group_id]" name="refresh" size="sm" class="animate-spin" />
                    <span>{{ transferringDays[quota.group_id] ? t('affiliate.subscriptionDays.transferring') : t('affiliate.subscriptionDays.transferButton') }}</span>
                  </button>
                </div>
              </div>
              <p v-else class="mt-3 text-sm text-gray-400 dark:text-dark-500">{{ t('affiliate.subscriptionDays.empty') }}</p>
            </div>
          </section>
        </div>

        <section class="card overflow-hidden">
          <div class="border-b border-gray-200 px-4 pt-4 dark:border-dark-700 sm:px-5">
            <div class="flex gap-1 overflow-x-auto" role="tablist">
              <button v-for="tab in activityTabs" :key="tab.key" type="button" role="tab" :aria-selected="activeTab === tab.key" class="whitespace-nowrap border-b-2 px-3 py-2 text-sm font-medium transition-colors" :class="activeTab === tab.key ? 'border-primary-500 text-primary-600 dark:text-primary-400' : 'border-transparent text-gray-500 hover:text-gray-900 dark:text-dark-400 dark:hover:text-dark-200'" @click="activeTab = tab.key">
                {{ tab.label }} <span class="ml-1 text-xs text-gray-400">{{ tab.count }}</span>
              </button>
            </div>
          </div>

          <div class="p-5">
            <template v-if="activeTab === 'invitees'">
              <div class="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <h2 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('affiliate.invitees.title') }}</h2>
                  <p class="mt-0.5 text-sm text-gray-500 dark:text-dark-400">{{ t('affiliate.invitees.description') }}</p>
                </div>
                <div class="relative w-full sm:w-64">
                  <Icon name="search" size="sm" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                  <input v-model="inviteeSearch" class="input h-9 pl-9" :placeholder="t('affiliate.invitees.search')" />
                </div>
              </div>
              <div v-if="filteredInvitees.length === 0" class="border border-dashed border-gray-300 py-10 text-center text-sm text-gray-500 dark:border-dark-700 dark:text-dark-400">{{ detail.invitees.length ? t('affiliate.invitees.noMatch') : t('affiliate.invitees.empty') }}</div>
              <div v-else>
                <div class="divide-y divide-gray-100 dark:divide-dark-800 sm:hidden">
                  <div v-for="item in filteredInvitees" :key="item.user_id" class="py-4">
                    <div class="flex items-start justify-between gap-3">
                      <div class="min-w-0">
                        <p class="truncate text-sm font-medium text-gray-900 dark:text-white">{{ item.email || '-' }}</p>
                        <p class="mt-0.5 truncate text-xs text-gray-400">{{ item.username || '-' }}</p>
                      </div>
                      <span class="inline-flex shrink-0 items-center gap-1 text-xs font-medium" :class="item.qualified ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400'">
                        <Icon :name="item.qualified ? 'checkCircle' : 'clock'" size="xs" />
                        {{ item.qualified ? t('affiliate.invitees.qualified') : t('affiliate.invitees.pending') }}
                      </span>
                    </div>
                    <div class="mt-3 grid grid-cols-2 gap-3 text-xs">
                      <div>
                        <p class="text-gray-400 dark:text-dark-500">{{ t('affiliate.invitees.columns.joinedAt') }}</p>
                        <p class="mt-1 text-gray-600 dark:text-dark-300">{{ formatDateTime(item.created_at) || '-' }}</p>
                      </div>
                      <div class="text-right">
                        <p class="text-gray-400 dark:text-dark-500">{{ t('affiliate.invitees.columns.rebate') }}</p>
                        <p class="mt-1 font-semibold text-emerald-600 dark:text-emerald-400">{{ formatCurrency(item.total_rebate) }}</p>
                      </div>
                    </div>
                  </div>
                </div>
                <div class="hidden overflow-x-auto sm:block">
                  <table class="w-full min-w-[720px] text-left text-sm">
                  <thead><tr class="border-b border-gray-200 text-gray-500 dark:border-dark-700 dark:text-dark-400">
                    <th class="px-3 py-2 font-medium">{{ t('affiliate.invitees.columns.customer') }}</th>
                    <th class="px-3 py-2 font-medium">{{ t('affiliate.invitees.columns.status') }}</th>
                    <th class="px-3 py-2 text-right font-medium">{{ t('affiliate.invitees.columns.rebate') }}</th>
                    <th class="px-3 py-2 font-medium">{{ t('affiliate.invitees.columns.joinedAt') }}</th>
                  </tr></thead>
                  <tbody><tr v-for="item in filteredInvitees" :key="item.user_id" class="border-b border-gray-100 last:border-0 dark:border-dark-800">
                    <td class="px-3 py-3"><p class="text-gray-900 dark:text-white">{{ item.email || '-' }}</p><p class="mt-0.5 text-xs text-gray-400">{{ item.username || '-' }}</p></td>
                    <td class="px-3 py-3"><span class="inline-flex items-center gap-1 text-xs font-medium" :class="item.qualified ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400'"><Icon :name="item.qualified ? 'checkCircle' : 'clock'" size="xs" />{{ item.qualified ? t('affiliate.invitees.qualified') : t('affiliate.invitees.pending') }}</span></td>
                    <td class="px-3 py-3 text-right font-medium text-emerald-600 dark:text-emerald-400">{{ formatCurrency(item.total_rebate) }}</td>
                    <td class="px-3 py-3 text-gray-600 dark:text-dark-300">{{ formatDateTime(item.created_at) || '-' }}</td>
                  </tr></tbody>
                  </table>
                </div>
              </div>
            </template>

            <template v-else-if="activeTab === 'rewards'">
              <div class="mb-4"><h2 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('affiliate.records.rebateTitle') }}</h2><p class="mt-0.5 text-sm text-gray-500 dark:text-dark-400">{{ t('affiliate.records.rebateDescription') }}</p></div>
              <div v-if="rebateLoading" class="flex justify-center py-10"><div class="h-6 w-6 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div></div>
              <div v-else-if="rebateRecords.length === 0" class="border border-dashed border-gray-300 py-10 text-center text-sm text-gray-500 dark:border-dark-700 dark:text-dark-400">{{ t('affiliate.records.rebateEmpty') }}</div>
              <div v-else>
                <div class="divide-y divide-gray-100 dark:divide-dark-800 sm:hidden">
                  <div v-for="row in rebateRecords" :key="row.ledger_id" class="py-4">
                    <div class="flex items-start justify-between gap-3">
                      <div class="min-w-0">
                        <p class="flex items-center gap-1.5 text-sm font-medium text-gray-900 dark:text-white">
                          <Icon :name="row.reward_type === 'subscription_days' ? 'calendar' : 'gift'" size="sm" :class="row.reward_type === 'subscription_days' ? 'text-sky-500' : 'text-emerald-500'" />
                          <span class="truncate">{{ rewardSourceLabel(row) }}</span>
                        </p>
                        <p v-if="row.redeem_code_id" class="mt-0.5 text-xs text-gray-400">#{{ row.redeem_code_id }}</p>
                      </div>
                      <p class="shrink-0 text-sm font-semibold text-emerald-600 dark:text-emerald-400">
                        {{ row.reward_type === 'subscription_days' ? t('affiliate.records.daysReward', { days: row.rebate_days }) : formatCurrency(row.rebate_amount) }}
                      </p>
                    </div>
                    <div class="mt-3 grid grid-cols-2 gap-3 text-xs">
                      <div>
                        <p class="text-gray-400 dark:text-dark-500">{{ t('affiliate.records.columns.invitee') }}</p>
                        <p class="mt-1 truncate text-gray-600 dark:text-dark-300">{{ row.invitee_email || row.invitee_username || '-' }}</p>
                      </div>
                      <div class="text-right">
                        <p class="text-gray-400 dark:text-dark-500">{{ t('affiliate.records.columns.group') }}</p>
                        <p class="mt-1 truncate text-gray-600 dark:text-dark-300">{{ row.group_name || '-' }}</p>
                      </div>
                      <div class="col-span-2">
                        <p class="text-gray-400 dark:text-dark-500">{{ t('affiliate.records.columns.rebatedAt') }}</p>
                        <p class="mt-1 text-gray-600 dark:text-dark-300">{{ formatDateTime(row.created_at) || '-' }}</p>
                      </div>
                    </div>
                  </div>
                </div>
                <div class="hidden overflow-x-auto sm:block">
                  <table class="w-full min-w-[720px] text-left text-sm">
                  <thead><tr class="border-b border-gray-200 text-gray-500 dark:border-dark-700 dark:text-dark-400">
                    <th class="px-3 py-2 font-medium">{{ t('affiliate.records.columns.source') }}</th>
                    <th class="px-3 py-2 font-medium">{{ t('affiliate.records.columns.invitee') }}</th>
                    <th class="px-3 py-2 font-medium">{{ t('affiliate.records.columns.group') }}</th>
                    <th class="px-3 py-2 text-right font-medium">{{ t('affiliate.records.columns.reward') }}</th>
                    <th class="px-3 py-2 font-medium">{{ t('affiliate.records.columns.rebatedAt') }}</th>
                  </tr></thead>
                  <tbody><tr v-for="row in rebateRecords" :key="row.ledger_id" class="border-b border-gray-100 last:border-0 dark:border-dark-800">
                    <td class="px-3 py-3"><span class="inline-flex items-center gap-1.5 font-medium text-gray-900 dark:text-white"><Icon :name="row.reward_type === 'subscription_days' ? 'calendar' : 'gift'" size="sm" :class="row.reward_type === 'subscription_days' ? 'text-sky-500' : 'text-emerald-500'" />{{ rewardSourceLabel(row) }}</span><p v-if="row.redeem_code_id" class="mt-0.5 text-xs text-gray-400">#{{ row.redeem_code_id }}</p></td>
                    <td class="px-3 py-3 text-gray-600 dark:text-dark-300">{{ row.invitee_email || row.invitee_username || '-' }}</td>
                    <td class="px-3 py-3 text-gray-600 dark:text-dark-300">{{ row.group_name || '-' }}</td>
                    <td class="px-3 py-3 text-right font-semibold text-emerald-600 dark:text-emerald-400">{{ row.reward_type === 'subscription_days' ? t('affiliate.records.daysReward', { days: row.rebate_days }) : formatCurrency(row.rebate_amount) }}</td>
                    <td class="px-3 py-3 text-gray-600 dark:text-dark-300">{{ formatDateTime(row.created_at) || '-' }}</td>
                  </tr></tbody>
                  </table>
                </div>
              </div>
              <Pagination v-if="rebatePagination.total > rebatePagination.page_size" class="mt-4" :page="rebatePagination.page" :total="rebatePagination.total" :page-size="rebatePagination.page_size" @update:page="handleRebatePageChange" />
            </template>

            <template v-else>
              <div class="mb-4"><h2 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('affiliate.records.transferTitle') }}</h2><p class="mt-0.5 text-sm text-gray-500 dark:text-dark-400">{{ t('affiliate.records.transferDescription') }}</p></div>
              <div v-if="transferLoading" class="flex justify-center py-10"><div class="h-6 w-6 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div></div>
              <div v-else-if="transferRecords.length === 0" class="border border-dashed border-gray-300 py-10 text-center text-sm text-gray-500 dark:border-dark-700 dark:text-dark-400">{{ t('affiliate.records.transferEmpty') }}</div>
              <div v-else>
                <div class="divide-y divide-gray-100 dark:divide-dark-800 sm:hidden">
                  <div v-for="row in transferRecords" :key="row.ledger_id" class="py-4">
                    <div class="flex items-start justify-between gap-3">
                      <div>
                        <p class="text-xs text-gray-400 dark:text-dark-500">{{ t('affiliate.records.columns.transferAmount') }}</p>
                        <p class="mt-1 text-base font-semibold text-emerald-600 dark:text-emerald-400">{{ formatCurrency(row.amount) }}</p>
                      </div>
                      <div class="text-right">
                        <p class="text-xs text-gray-400 dark:text-dark-500">{{ t('affiliate.records.columns.balanceAfter') }}</p>
                        <p class="mt-1 text-sm text-gray-600 dark:text-dark-300">{{ row.balance_after == null ? '-' : formatCurrency(row.balance_after) }}</p>
                      </div>
                    </div>
                    <p class="mt-3 text-xs text-gray-500 dark:text-dark-400">{{ formatDateTime(row.created_at) || '-' }}</p>
                  </div>
                </div>
                <div class="hidden overflow-x-auto sm:block">
                  <table class="w-full min-w-[560px] text-left text-sm">
                  <thead><tr class="border-b border-gray-200 text-gray-500 dark:border-dark-700 dark:text-dark-400"><th class="px-3 py-2 text-right font-medium">{{ t('affiliate.records.columns.transferAmount') }}</th><th class="px-3 py-2 text-right font-medium">{{ t('affiliate.records.columns.balanceAfter') }}</th><th class="px-3 py-2 font-medium">{{ t('affiliate.records.columns.transferredAt') }}</th></tr></thead>
                  <tbody><tr v-for="row in transferRecords" :key="row.ledger_id" class="border-b border-gray-100 last:border-0 dark:border-dark-800"><td class="px-3 py-3 text-right font-semibold text-emerald-600 dark:text-emerald-400">{{ formatCurrency(row.amount) }}</td><td class="px-3 py-3 text-right text-gray-600 dark:text-dark-300">{{ row.balance_after == null ? '-' : formatCurrency(row.balance_after) }}</td><td class="px-3 py-3 text-gray-600 dark:text-dark-300">{{ formatDateTime(row.created_at) || '-' }}</td></tr></tbody>
                  </table>
                </div>
              </div>
              <Pagination v-if="transferPagination.total > transferPagination.page_size" class="mt-4" :page="transferPagination.page" :total="transferPagination.total" :page-size="transferPagination.page_size" @update:page="handleTransferPageChange" />
            </template>
          </div>
        </section>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import Pagination from '@/components/common/Pagination.vue'
import userAPI from '@/api/user'
import type { UserAffiliateDetail } from '@/types'
import type { UserAffiliateRebateRecord, UserAffiliateTransferRecord } from '@/api/user'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { useSubscriptionStore } from '@/stores/subscriptions'
import { useClipboard } from '@/composables/useClipboard'
import { formatCurrency, formatDateTime } from '@/utils/format'
import { extractApiErrorMessage } from '@/utils/apiError'

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()
const subscriptionStore = useSubscriptionStore()
const { copyToClipboard } = useClipboard()

const loading = ref(true)
const refreshing = ref(false)
const transferring = ref(false)
const transferringDays = ref<Record<number, boolean>>({})
const detail = ref<UserAffiliateDetail | null>(null)
const activeTab = ref<ActivityTab>('invitees')
const inviteeSearch = ref('')

type ActivityTab = 'invitees' | 'rewards' | 'transfers'

// 返利记录 / 提取记录
const rebateRecords = ref<UserAffiliateRebateRecord[]>([])
const rebateLoading = ref(false)
const rebatePagination = ref({ page: 1, page_size: 10, total: 0 })
const transferRecords = ref<UserAffiliateTransferRecord[]>([])
const transferLoading = ref(false)
const transferPagination = ref({ page: 1, page_size: 10, total: 0 })

const inviteLink = computed(() => {
  if (!detail.value) return ''
  if (typeof window === 'undefined') return `/register?aff=${encodeURIComponent(detail.value.aff_code)}`
  return `${window.location.origin}/register?aff=${encodeURIComponent(detail.value.aff_code)}`
})

// Rebate rate is a percentage in the range [0, 100]; backend already clamps it.
// We trim trailing zeros (e.g. 20.00 → "20", 12.50 → "12.5") for a cleaner UI.
const formattedRebateRate = computed(() => {
  const v = detail.value?.effective_rebate_rate_percent ?? 0
  const rounded = Math.round(v * 100) / 100
  return Number.isInteger(rounded) ? String(rounded) : rounded.toString()
})

const eligibilityMet = computed(() => {
  if (!detail.value) return false
  return detail.value.min_qualified_invitees <= 0
    || detail.value.qualified_invitee_count >= detail.value.min_qualified_invitees
})

const eligibilityRemaining = computed(() => {
  if (!detail.value) return 0
  return Math.max(0, detail.value.min_qualified_invitees - detail.value.qualified_invitee_count)
})

const eligibilityProgress = computed(() => {
  if (!detail.value || detail.value.min_qualified_invitees <= 0) return 100
  return Math.min(100, Math.round(
    (detail.value.qualified_invitee_count / detail.value.min_qualified_invitees) * 100,
  ))
})

const eligibilityProgressCurrent = computed(() => {
  if (!detail.value) return 0
  return Math.min(detail.value.qualified_invitee_count, detail.value.min_qualified_invitees)
})

const thresholdSummary = computed(() => {
  if (!detail.value || detail.value.min_qualified_invitees <= 0) {
    return t('affiliate.threshold.noThreshold')
  }
  return eligibilityMet.value
    ? t('affiliate.threshold.readySummary', { start: detail.value.min_qualified_invitees + 1 })
    : t('affiliate.threshold.pendingSummary')
})

const filteredInvitees = computed(() => {
  if (!detail.value) return []
  const query = inviteeSearch.value.trim().toLocaleLowerCase()
  if (!query) return detail.value.invitees
  return detail.value.invitees.filter((item) => (
    (item.email || '').toLocaleLowerCase().includes(query)
    || (item.username || '').toLocaleLowerCase().includes(query)
  ))
})

const activityTabs = computed(() => [
  { key: 'invitees' as const, label: t('affiliate.invitees.title'), count: detail.value?.invitees.length || 0 },
  { key: 'rewards' as const, label: t('affiliate.records.rebateTitle'), count: rebatePagination.value.total },
  { key: 'transfers' as const, label: t('affiliate.records.transferTitle'), count: transferPagination.value.total },
])

function formatCount(value: number): string {
  return value.toLocaleString()
}

function rewardSourceLabel(row: UserAffiliateRebateRecord): string {
  const labels: Record<UserAffiliateRebateRecord['source_type'], string> = {
    balance_redeem: 'affiliate.records.source.balanceRedeem',
    subscription_redeem: 'affiliate.records.source.subscriptionRedeem',
    balance_payment: 'affiliate.records.source.balancePayment',
    subscription_payment: 'affiliate.records.source.subscriptionPayment',
  }
  return t(labels[row.source_type] || 'affiliate.records.source.unknown')
}

async function loadAffiliateDetail(silent = false): Promise<void> {
  if (!silent) {
    loading.value = true
  }
  try {
    const loaded = await userAPI.getAffiliateDetail()
    detail.value = {
      ...loaded,
      invitees: loaded.invitees || [],
      subscription_days_quota: loaded.subscription_days_quota || [],
    }
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('affiliate.loadFailed')))
  } finally {
    if (!silent) {
      loading.value = false
    }
  }
}

function userTimezone(): string {
  try {
    return Intl.DateTimeFormat().resolvedOptions().timeZone
  } catch {
    return 'UTC'
  }
}

async function loadRebateRecords(): Promise<void> {
  rebateLoading.value = true
  try {
    const res = await userAPI.listAffiliateRebateRecords({
      page: rebatePagination.value.page,
      page_size: rebatePagination.value.page_size,
      timezone: userTimezone(),
    })
    rebateRecords.value = res.items || []
    rebatePagination.value.total = res.total || 0
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('affiliate.records.loadFailed')))
  } finally {
    rebateLoading.value = false
  }
}

async function loadTransferRecords(): Promise<void> {
  transferLoading.value = true
  try {
    const res = await userAPI.listAffiliateTransferRecords({
      page: transferPagination.value.page,
      page_size: transferPagination.value.page_size,
      timezone: userTimezone(),
    })
    transferRecords.value = res.items || []
    transferPagination.value.total = res.total || 0
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('affiliate.records.loadFailed')))
  } finally {
    transferLoading.value = false
  }
}

async function refreshAll(): Promise<void> {
  if (refreshing.value) return
  refreshing.value = true
  try {
    await Promise.all([
      loadAffiliateDetail(true),
      loadRebateRecords(),
      loadTransferRecords(),
    ])
  } finally {
    refreshing.value = false
  }
}

function handleRebatePageChange(page: number): void {
  rebatePagination.value.page = page
  void loadRebateRecords()
}

function handleTransferPageChange(page: number): void {
  transferPagination.value.page = page
  void loadTransferRecords()
}

async function copyCode(): Promise<void> {
  if (!detail.value?.aff_code) return
  await copyToClipboard(detail.value.aff_code, t('affiliate.codeCopied'))
}

async function copyInviteLink(): Promise<void> {
  if (!inviteLink.value) return
  await copyToClipboard(inviteLink.value, t('affiliate.linkCopied'))
}

async function transferQuota(): Promise<void> {
  if (!detail.value || detail.value.aff_quota <= 0 || transferring.value) return
  transferring.value = true
  try {
    const resp = await userAPI.transferAffiliateQuota()
    appStore.showSuccess(t('affiliate.transfer.success', { amount: formatCurrency(resp.transferred_quota) }))
    await Promise.all([
      loadAffiliateDetail(true),
      authStore.refreshUser().catch(() => undefined),
    ])
    // 提取后刷新提取记录（回到第一页展示最新记录）
    transferPagination.value.page = 1
    void loadTransferRecords()
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('affiliate.transferFailed')))
  } finally {
    transferring.value = false
  }
}

async function transferDays(quota: { group_id: number; group_name: string; pending_days: number }): Promise<void> {
  if (quota.pending_days <= 0 || transferringDays.value[quota.group_id]) return
  transferringDays.value[quota.group_id] = true
  try {
    const resp = await userAPI.transferAffiliateSubscriptionDays(quota.group_id)
    appStore.showSuccess(t('affiliate.subscriptionDays.transferSuccess', { days: resp.transferred_days, group: quota.group_name || `Group #${quota.group_id}` }))
    await Promise.all([
      loadAffiliateDetail(true),
      subscriptionStore.fetchActiveSubscriptions(true).catch(() => undefined),
      authStore.refreshUser().catch(() => undefined),
    ])
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('affiliate.transferFailed')))
  } finally {
    transferringDays.value[quota.group_id] = false
  }
}

onMounted(() => {
  void loadAffiliateDetail()
  void loadRebateRecords()
  void loadTransferRecords()
})
</script>
