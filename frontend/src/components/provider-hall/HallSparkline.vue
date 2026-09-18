<template>
  <figure class="hall-sparkline m-0 inline-flex flex-col" :style="{ width: `${width}px` }">
    <svg
      :width="width"
      :height="height"
      :viewBox="`0 0 ${width} ${height}`"
      role="img"
      :aria-label="ariaLabel"
      class="block overflow-visible"
    >
      <rect
        v-for="(gap, index) in geometry.gaps"
        :key="`gap-${index}`"
        :x="gap.x"
        y="0"
        :width="gap.width"
        :height="height"
        class="hall-sparkline-gap"
        data-testid="sparkline-gap"
      />
      <path
        v-if="geometry.hasProbe"
        :d="geometry.probe"
        fill="none"
        class="hall-sparkline-probe"
        stroke-width="1.5"
        stroke-linejoin="round"
        stroke-linecap="round"
        data-testid="sparkline-probe"
      />
      <path
        v-if="geometry.hasReal"
        :d="geometry.real"
        fill="none"
        class="hall-sparkline-real"
        stroke-width="2"
        stroke-linejoin="round"
        stroke-linecap="round"
        data-testid="sparkline-real"
      />
      <text
        v-if="!geometry.hasReal && !geometry.hasProbe"
        :x="width / 2"
        :y="height / 2"
        text-anchor="middle"
        dominant-baseline="central"
        class="hall-sparkline-empty text-[10px]"
      >-</text>
    </svg>
    <figcaption v-if="showLegend" class="mt-0.5 flex items-center gap-2 text-[10px] leading-none text-gray-400 dark:text-dark-500">
      <span class="inline-flex items-center gap-1"><i class="hall-legend-probe" aria-hidden="true"></i>{{ t('providerHall.legend.probe') }}</span>
      <span class="inline-flex items-center gap-1"><i class="hall-legend-real" aria-hidden="true"></i>{{ t('providerHall.legend.fast95') }}</span>
    </figcaption>
  </figure>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { ProviderHallSparkPoint } from '@/types/providerHall'
import { buildSparkline } from './hallGeometry'

const props = withDefaults(
  defineProps<{
    points: ProviderHallSparkPoint[]
    name?: string
    width?: number
    height?: number
    showLegend?: boolean
  }>(),
  { name: '', width: 200, height: 40, showLegend: true }
)

const { t } = useI18n()
const geometry = computed(() => buildSparkline(props.points || [], props.width, props.height))
const ariaLabel = computed(() =>
  t('providerHall.aria.sparkline', {
    name: props.name,
    points: props.points?.length ?? 0,
    gaps: geometry.value.gaps.length,
  })
)
</script>

<style scoped>
/* Series slot 1 (blue) = real fastest-95% mean; slot 2 (orange) = probe. */
.hall-sparkline-real {
  stroke: #2a78d6;
}
.hall-sparkline-probe {
  stroke: #eb6834;
}
.dark .hall-sparkline-real,
.dark-theme .hall-sparkline-real {
  stroke: #3987e5;
}
.dark .hall-sparkline-probe,
.dark-theme .hall-sparkline-probe {
  stroke: #d95926;
}
/* Gap wash: critical status hue at low opacity, never a saturated block. */
.hall-sparkline-gap {
  fill: rgba(208, 59, 59, 0.16);
}
.hall-sparkline-empty {
  fill: #9ca3af;
}
.hall-legend-probe,
.hall-legend-real {
  display: inline-block;
  width: 10px;
  height: 2px;
  border-radius: 1px;
}
.hall-legend-probe {
  background: #eb6834;
}
.hall-legend-real {
  background: #2a78d6;
}
</style>
