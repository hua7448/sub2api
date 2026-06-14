<template>
  <div
    v-if="variant === 'pixel'"
    class="auth-pixel-shell"
  >
    <div class="auth-pixel-grid" aria-hidden="true"></div>
    <div class="auth-pixel-dither auth-pixel-dither-top" aria-hidden="true"></div>
    <div class="auth-pixel-dither auth-pixel-dither-bottom" aria-hidden="true"></div>

    <main class="auth-pixel-stage">
      <section class="auth-pixel-art" aria-hidden="true">
        <div class="pixel-window">
          <div class="pixel-window-bar">
            <span></span>
            <span></span>
            <span></span>
          </div>
          <div class="pixel-window-body">
            <div class="pixel-code">
              <div><span class="pixel-code-line">01</span> const hello = 'world'</div>
              <div><span class="pixel-code-line">02</span> auth.node.boot()</div>
              <div><span class="pixel-code-line">03</span> session.ready = true</div>
              <div><span class="pixel-code-line">04</span> route.guard.pass()</div>
            </div>
            <div class="pixel-ridge pixel-ridge-a"></div>
            <div class="pixel-ridge pixel-ridge-b"></div>
            <div class="pixel-ridge pixel-ridge-c"></div>
            <div class="pixel-cursor"></div>
          </div>
        </div>
        <div class="pixel-stack pixel-stack-a"></div>
        <div class="pixel-stack pixel-stack-b"></div>
      </section>

      <section class="auth-pixel-panel">
        <div v-if="settingsLoaded" class="auth-pixel-brand">
          <div class="auth-pixel-logo">
            <img :src="siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
          </div>
        </div>

        <div class="auth-pixel-card">
          <slot />
        </div>

        <div class="auth-pixel-footer">
          <slot name="footer" />
        </div>
      </section>
    </main>
  </div>

  <div v-else class="relative flex min-h-screen items-center justify-center overflow-hidden p-4">
    <!-- Background -->
    <div
      class="absolute inset-0 bg-gradient-to-br from-gray-50 via-primary-50/30 to-gray-100 dark:from-dark-950 dark:via-dark-900 dark:to-dark-950"
    ></div>

    <!-- Decorative Elements -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <!-- Gradient Orbs -->
      <div
        class="absolute -right-40 -top-40 h-80 w-80 rounded-full bg-primary-400/20 blur-3xl"
      ></div>
      <div
        class="absolute -bottom-40 -left-40 h-80 w-80 rounded-full bg-primary-500/15 blur-3xl"
      ></div>
      <div
        class="absolute left-1/2 top-1/2 h-96 w-96 -translate-x-1/2 -translate-y-1/2 rounded-full bg-primary-300/10 blur-3xl"
      ></div>

      <!-- Grid Pattern -->
      <div
        class="absolute inset-0 bg-[linear-gradient(rgba(184,95,63,0.04)_1px,transparent_1px),linear-gradient(90deg,rgba(184,95,63,0.04)_1px,transparent_1px)] bg-[size:64px_64px]"
      ></div>
    </div>

    <!-- Content Container -->
    <div class="relative z-10 w-full max-w-md">
      <!-- Logo/Brand -->
      <div class="mb-8 text-center">
        <!-- Custom Logo or Default Logo -->
        <template v-if="settingsLoaded">
          <div
            class="mb-4 inline-flex h-16 w-16 items-center justify-center overflow-hidden rounded-2xl shadow-lg shadow-primary-500/30"
          >
            <img :src="siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <h1 class="text-gradient mb-2 text-3xl font-bold">
            {{ siteName }}
          </h1>
          <p class="text-sm text-gray-500 dark:text-dark-400">
            {{ siteSubtitle }}
          </p>
        </template>
      </div>

      <!-- Card Container -->
      <div class="card-glass rounded-2xl p-8 shadow-glass">
        <slot />
      </div>

      <!-- Footer Links -->
      <div class="mt-6 text-center text-sm">
        <slot name="footer" />
      </div>

      <!-- Copyright -->
      <div class="mt-8 text-center text-xs text-gray-400 dark:text-dark-500">
        &copy; {{ currentYear }} {{ siteName }}. All rights reserved.
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useAppStore } from '@/stores'
import { sanitizeUrl } from '@/utils/url'

withDefaults(defineProps<{
  variant?: 'default' | 'pixel'
}>(), {
  variant: 'default'
})

const appStore = useAppStore()

const siteName = computed(() => appStore.siteName || 'Sub2API')
const siteLogo = computed(() => sanitizeUrl(appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || 'Subscription to API Conversion Platform')
const settingsLoaded = computed(() => appStore.publicSettingsLoaded)

const currentYear = computed(() => new Date().getFullYear())

onMounted(() => {
  appStore.fetchPublicSettings()
})
</script>

<style scoped>
.text-gradient {
  @apply bg-gradient-to-r from-primary-600 to-primary-500 bg-clip-text text-transparent;
}

.auth-pixel-shell {
  position: relative;
  min-height: 100vh;
  overflow: hidden;
  background:
    linear-gradient(135deg, rgba(255, 250, 242, 0.96), rgba(244, 228, 214, 0.9)),
    #f7f3ea;
  color: #1f1d1a;
}

.dark .auth-pixel-shell {
  background:
    linear-gradient(135deg, rgba(31, 29, 26, 0.98), rgba(48, 44, 39, 0.94)),
    #141210;
  color: #f1ede6;
}

.auth-pixel-grid {
  position: absolute;
  inset: 0;
  opacity: 0.58;
  background-image:
    linear-gradient(rgba(94, 85, 75, 0.08) 1px, transparent 1px),
    linear-gradient(90deg, rgba(94, 85, 75, 0.08) 1px, transparent 1px);
  background-size: 18px 18px;
}

.dark .auth-pixel-grid {
  opacity: 0.38;
  background-image:
    linear-gradient(rgba(208, 198, 184, 0.08) 1px, transparent 1px),
    linear-gradient(90deg, rgba(208, 198, 184, 0.08) 1px, transparent 1px);
}

.auth-pixel-dither {
  position: absolute;
  width: 312px;
  height: 312px;
  opacity: 0.35;
  background-image:
    linear-gradient(45deg, rgba(184, 95, 63, 0.22) 25%, transparent 25%),
    linear-gradient(-45deg, rgba(184, 95, 63, 0.18) 25%, transparent 25%);
  background-position: 0 0, 9px 0;
  background-size: 18px 18px;
}

.auth-pixel-dither-top {
  right: -96px;
  top: -88px;
}

.auth-pixel-dither-bottom {
  bottom: -120px;
  left: -92px;
  transform: rotate(180deg);
}

.auth-pixel-stage {
  position: relative;
  z-index: 1;
  display: grid;
  min-height: 100vh;
  grid-template-columns: minmax(0, 1fr);
  align-items: center;
  gap: 2rem;
  padding: 1.25rem;
}

@media (min-width: 1024px) {
  .auth-pixel-stage {
    grid-template-columns: minmax(360px, 0.95fr) minmax(420px, 520px);
    padding: 3rem clamp(2rem, 5vw, 5rem);
  }
}

.auth-pixel-art {
  display: none;
}

@media (min-width: 1024px) {
  .auth-pixel-art {
    position: relative;
    display: flex;
    min-height: 520px;
    align-items: center;
    justify-content: center;
  }
}

.pixel-window {
  position: relative;
  width: min(76vw, 520px);
  border: 4px solid #302c27;
  background: #fffaf2;
  box-shadow:
    12px 12px 0 #302c27,
    28px 28px 0 rgba(184, 95, 63, 0.18);
}

.dark .pixel-window {
  border-color: #f1ede6;
  background: #1f1d1a;
  box-shadow:
    12px 12px 0 #f1ede6,
    28px 28px 0 rgba(201, 123, 87, 0.28);
}

.pixel-window-bar {
  display: flex;
  gap: 8px;
  border-bottom: 4px solid #302c27;
  background: #e9c5ad;
  padding: 12px;
}

.dark .pixel-window-bar {
  border-color: #f1ede6;
  background: #663229;
}

.pixel-window-bar span {
  width: 14px;
  height: 14px;
  border: 3px solid #302c27;
  background: #fffaf2;
}

.dark .pixel-window-bar span {
  border-color: #f1ede6;
  background: #302c27;
}

.pixel-window-body {
  position: relative;
  height: 328px;
  overflow: hidden;
  background:
    linear-gradient(90deg, transparent 0 22px, rgba(184, 95, 63, 0.08) 22px 24px, transparent 24px 48px),
    linear-gradient(#fffaf2, #f4e4d6);
  background-size: 48px 48px, auto;
}

.dark .pixel-window-body {
  background:
    linear-gradient(90deg, transparent 0 22px, rgba(201, 123, 87, 0.16) 22px 24px, transparent 24px 48px),
    linear-gradient(#1f1d1a, #302c27);
}

.pixel-ridge {
  position: absolute;
  left: 52px;
  height: 26px;
  border: 4px solid #302c27;
  background: #b85f3f;
  opacity: 0.18;
}

.dark .pixel-ridge {
  border-color: #f1ede6;
  background: #c97b57;
}

.pixel-ridge-a {
  top: 70px;
  width: 58%;
}

.pixel-ridge-b {
  top: 132px;
  width: 42%;
  background: #474038;
}

.pixel-ridge-c {
  top: 194px;
  width: 70%;
  background: #daa17f;
}

.pixel-code {
  position: absolute;
  left: 38px;
  right: 38px;
  top: 52px;
  z-index: 1;
  display: grid;
  gap: 14px;
  color: #302c27;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: clamp(0.85rem, 1.3vw, 1rem);
  font-weight: 800;
  line-height: 1.4;
}

.dark .pixel-code {
  color: #f1ede6;
}

.pixel-code > div {
  width: fit-content;
  max-width: 100%;
  border: 3px solid #302c27;
  background: rgba(255, 250, 242, 0.92);
  padding: 7px 10px;
  box-shadow: 5px 5px 0 rgba(48, 44, 39, 0.22);
  white-space: nowrap;
}

.dark .pixel-code > div {
  border-color: #f1ede6;
  background: rgba(31, 29, 26, 0.9);
  box-shadow: 5px 5px 0 rgba(241, 237, 230, 0.2);
}

.pixel-code > div:nth-child(2) {
  margin-left: 22px;
}

.pixel-code > div:nth-child(3) {
  margin-left: 44px;
}

.pixel-code > div:nth-child(4) {
  margin-left: 16px;
}

.pixel-code-line {
  margin-right: 10px;
  color: #9d4b33;
}

.dark .pixel-code-line {
  color: #daa17f;
}

.pixel-cursor {
  position: absolute;
  right: 78px;
  bottom: 72px;
  width: 56px;
  height: 8px;
  background: #302c27;
  box-shadow: 5px 5px 0 rgba(48, 44, 39, 0.22);
}

.dark .pixel-cursor {
  background: #f1ede6;
  box-shadow: 5px 5px 0 rgba(241, 237, 230, 0.2);
}

.pixel-stack {
  position: absolute;
  border: 4px solid #302c27;
  background: #f4e4d6;
}

.dark .pixel-stack {
  border-color: #f1ede6;
  background: #302c27;
}

.pixel-stack-a {
  left: 8%;
  top: 18%;
  width: 72px;
  height: 72px;
  box-shadow: 10px 10px 0 rgba(48, 44, 39, 0.22);
}

.pixel-stack-b {
  bottom: 14%;
  right: 10%;
  width: 96px;
  height: 48px;
  background: #e9c5ad;
  box-shadow: 10px 10px 0 rgba(48, 44, 39, 0.22);
}

.auth-pixel-panel {
  width: min(100%, 460px);
  justify-self: center;
}

@media (min-width: 1024px) {
  .auth-pixel-panel {
    justify-self: end;
  }
}

.auth-pixel-brand {
  display: flex;
  justify-content: center;
  margin-bottom: 1rem;
}

.auth-pixel-logo {
  display: flex;
  width: 48px;
  height: 48px;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border: 3px solid #302c27;
  background: #fffaf2;
  box-shadow: 6px 6px 0 #302c27;
}

.dark .auth-pixel-logo {
  border-color: #f1ede6;
  background: #302c27;
  box-shadow: 6px 6px 0 #f1ede6;
}

.auth-pixel-card {
  border: 4px solid #302c27;
  background: rgba(255, 250, 242, 0.96);
  padding: clamp(1.25rem, 4vw, 2rem);
  box-shadow:
    10px 10px 0 #302c27,
    20px 20px 0 rgba(184, 95, 63, 0.16);
}

.dark .auth-pixel-card {
  border-color: #f1ede6;
  background: rgba(31, 29, 26, 0.96);
  box-shadow:
    10px 10px 0 #f1ede6,
    20px 20px 0 rgba(201, 123, 87, 0.24);
}

.auth-pixel-footer {
  margin-top: 1.5rem;
  text-align: center;
  font-size: 0.875rem;
}

</style>
