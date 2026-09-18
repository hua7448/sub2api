<template>
  <div class="card">
    <div class="border-b border-gray-100 px-6 py-4">
      <h2 class="text-lg font-semibold text-txt-primary">{{ t('dashboard.quickActions') }}</h2>
    </div>
    <div class="space-y-3 p-4">
      <button @click="router.push('/keys')" class="group flex w-full items-center gap-4 rounded-xl bg-cyber-bg p-4 text-left transition-all duration-200 hover:bg-cyber-elevated">
        <div class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-xl bg-primary-100 transition-transform group-hover:scale-105 bg-primary-900/30">
          <Icon name="key" size="lg" class="text-neon-cyan text-primary-400" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-sm font-medium text-txt-primary">{{ t('dashboard.createApiKey') }}</p>
          <p class="text-xs text-txt-secondary">{{ t('dashboard.generateNewKey') }}</p>
        </div>
        <Icon
          name="chevronRight"
          size="md"
          class="text-txt-dim transition-colors group-hover:text-neon-cyan"
        />
      </button>

      <button @click="router.push('/usage')" class="group flex w-full items-center gap-4 rounded-xl bg-cyber-bg p-4 text-left transition-all duration-200 hover:bg-cyber-elevated">
        <div class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-xl bg-emerald-100 transition-transform group-hover:scale-105 bg-emerald-900/30">
          <Icon name="chart" size="lg" class="text-emerald-600 text-emerald-400" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-sm font-medium text-txt-primary">{{ t('dashboard.viewUsage') }}</p>
          <p class="text-xs text-txt-secondary">{{ t('dashboard.checkDetailedLogs') }}</p>
        </div>
        <Icon
          name="chevronRight"
          size="md"
          class="text-txt-dim transition-colors group-hover:text-emerald-500"
        />
      </button>

      <button @click="goToRecharge" class="group flex w-full items-center gap-4 rounded-xl bg-cyber-bg p-4 text-left transition-all duration-200 hover:bg-cyber-elevated">
        <div class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-xl bg-sky-100 transition-transform group-hover:scale-105 bg-sky-900/30">
          <Icon name="creditCard" size="lg" class="text-sky-600 text-sky-400" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-sm font-medium text-txt-primary">{{ t('dashboard.goToRecharge') }}</p>
          <p class="text-xs text-txt-secondary">{{ t('dashboard.goToRechargeHint') }}</p>
        </div>
        <Icon
          name="externalLink"
          size="md"
          class="text-txt-dim transition-colors group-hover:text-sky-500"
        />
      </button>

      <button @click="router.push('/redeem')" class="group flex w-full items-center gap-4 rounded-xl bg-cyber-bg p-4 text-left transition-all duration-200 hover:bg-cyber-elevated">
        <div class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-xl bg-amber-100 transition-transform group-hover:scale-105 bg-amber-900/30">
          <Icon name="gift" size="lg" class="text-amber-600 text-amber-400" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-sm font-medium text-txt-primary">{{ t('dashboard.redeemCode') }}</p>
          <p class="text-xs text-txt-secondary">{{ t('dashboard.addBalanceWithCode') }}</p>
        </div>
        <Icon
          name="chevronRight"
          size="md"
          class="text-txt-dim transition-colors group-hover:text-neon-amber"
        />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'

const RECHARGE_URL = 'https://fk.sharesai.xyz/cat/4'

const router = useRouter()
const { t } = useI18n()

function goToRecharge() {
  window.location.assign(RECHARGE_URL)
}
</script>
