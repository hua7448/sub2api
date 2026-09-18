import { readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

import { describe, expect, it } from 'vitest'

const testDir = dirname(fileURLToPath(import.meta.url))
const viewSource = readFileSync(resolve(testDir, '../AiStudioView.vue'), 'utf8')
const zhLocaleSource = [
  '../../../i18n/locales/zh/common.ts',
  '../../../i18n/locales/zh/dashboard.ts',
].map((path) => readFileSync(resolve(testDir, path), 'utf8')).join('\n')
const enLocaleSource = [
  '../../../i18n/locales/en/common.ts',
  '../../../i18n/locales/en/dashboard.ts',
].map((path) => readFileSync(resolve(testDir, path), 'utf8')).join('\n')
const studioBundleSource = readFileSync(resolve(testDir, '../../../../public/image-studio/assets/index-Dr6rC7Ut.js'), 'utf8')

describe('AiStudioView shell', () => {
  it('uses product-neutral copy and removes gateway helper text', () => {
    expect(viewSource).not.toContain('ChatGpt-Image-Studio')
    expect(viewSource).not.toContain('gatewayAutoHint')
    expect(zhLocaleSource).toContain("aiStudio: 'AI画图'")
    expect(zhLocaleSource).toContain("title: 'AI画图'")
    expect(zhLocaleSource).not.toContain('ChatGpt-Image-Studio')
    expect(zhLocaleSource).not.toContain('网关地址会自动使用当前站点')
    expect(zhLocaleSource).not.toContain('AI 创作')
    expect(enLocaleSource).not.toContain('ChatGpt-Image-Studio')
    expect(enLocaleSource).not.toContain('current site as its gateway')
  })

  it('groups key controls in a compact workspace toolbar', () => {
    expect(viewSource).toContain('data-testid="ai-studio-shell"')
    expect(viewSource).toContain('data-testid="ai-studio-toolbar"')
    expect(viewSource).toContain('data-testid="ai-studio-frame"')
    expect(viewSource).toContain('xl:flex-row xl:items-center xl:justify-between')
    expect(viewSource).toContain('xl:w-[460px] xl:grid-cols-[minmax(0,1fr)_auto]')
    expect(viewSource).toContain(':src="studioFrameSrc"')
    expect(viewSource).toContain("const studioFrameSrc = '/image-studio/?v=scope5'")
  })

  it('uses the same-origin gateway for embedded image requests', () => {
    expect(viewSource).toContain('return window.location.origin')
    expect(viewSource).not.toContain('cachedPublicSettings?.api_base_url')
  })

  it('passes a scoped user identity to isolate local image history', () => {
    expect(viewSource).toContain('userScope?: string')
    expect(viewSource).toContain('function buildUserScope()')
    expect(viewSource).toContain('userScope: buildUserScope()')
    expect(viewSource).toContain("`user:${userId}`")
    expect(viewSource).toContain('const keyOwnerId = selectedKey.value?.user_id')
    expect(viewSource).toContain("`user:${keyOwnerId}`")
    expect(viewSource).toContain('selectedKey.value?.user_id')
    expect(studioBundleSource).toContain('userScope:d||void 0')
    expect(studioBundleSource).toContain('items:"+String(n).replace')
    expect(studioBundleSource).toContain('sub2api:image-studio-config:v2')
    expect(studioBundleSource).toContain('r!==o&&(Wt=null,ur=null,Ci=null,pl=Promise.resolve())')
  })

  it('targets the actual image studio iframe origin for bridge messages', () => {
    expect(viewSource).toContain('const studioFrameOrigin = computed(() => {')
    expect(viewSource).toContain("studioFrame.value?.getAttribute('src') || studioFrameSrc")
    expect(viewSource).toContain('new URL(rawSrc, window.location.href).origin')
    expect(viewSource).toContain('studioFrameOrigin.value || window.location.origin')
    expect(viewSource).toContain("event.origin !== (studioFrameOrigin.value || window.location.origin)")
  })

  it('refreshes account data after image usage changes inside the iframe', () => {
    expect(viewSource).toContain('sub2api-image-studio-usage-updated')
    expect(viewSource).toContain('window.addEventListener(\'message\', handleStudioMessage)')
    expect(viewSource).toContain('window.removeEventListener(\'message\', handleStudioMessage)')
    expect(viewSource).toContain('void authStore.refreshUser()')
  })

  it('generates multiple images with bounded parallel requests instead of the unsupported n parameter', () => {
    expect(studioBundleSource).not.toContain('n:Math.max(1,i)')
    expect(studioBundleSource).toContain('const xh0=10')
    expect(studioBundleSource).toContain('async function xh1')
    expect(studioBundleSource).toContain('await xh1(m,async')
    expect(studioBundleSource).toContain('await xh1(p,async')
    expect(studioBundleSource).toContain('count:m=1')
    expect(studioBundleSource).not.toContain('for(let p=0;p<m;p++){const E=await fetch')
    expect(studioBundleSource).not.toContain('for(let K=0;K<p;K++){const J=new FormData')
    expect(studioBundleSource).toContain('`${f}:${h+1}`')
    expect(studioBundleSource).toContain('`${f}:${L+1}`')
    expect(studioBundleSource).toContain('data:h')
    expect(studioBundleSource).toContain('data:L')
  })
})
