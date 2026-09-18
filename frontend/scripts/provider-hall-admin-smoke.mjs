import assert from 'node:assert/strict'
import { mkdir } from 'node:fs/promises'
import { resolve } from 'node:path'
import { pathToFileURL } from 'node:url'
import { chromium } from 'playwright'

/**
 * Walks the admin provider-hall configuration the way an operator does now:
 * probe targets live in an inline disclosure row on the group table, while
 * listing metadata (display name, description, order, listed) sits in a
 * separate dialog. Both halves are exercised so the round trip is covered.
 */
export async function verifyAdmin({ origin, token, user, artifacts, groupName }) {
  await mkdir(artifacts, { recursive: true })
  const browser = await chromium.launch({ headless: true, ...(process.env.PROVIDER_HALL_BROWSER_CHANNEL ? { channel: process.env.PROVIDER_HALL_BROWSER_CHANNEL } : {}) })
  const errors = []
  try {
    for (const width of [390, 1440, 1920]) {
      for (const dark of [false, true]) {
        const context = await browser.newContext({ viewport: { width, height: 1000 }, reducedMotion: 'reduce' })
        await context.addInitScript(({ token, user }) => {
          localStorage.setItem('auth_token', token)
          localStorage.setItem('auth_user', JSON.stringify(user))
          localStorage.setItem('sub2api_locale', 'zh')
          localStorage.setItem(`admin_guide_${user.id}_${user.role}_v4_interactive`, 'true')
        }, { token, user })
        const page = await context.newPage()
        page.on('pageerror', error => errors.push(error.message))
        page.on('dialog', dialog => dialog.accept())
        await page.goto(`${origin}/admin/provider-hall`)
        await page.getByRole('tab', { name: '分组与目标', exact: true }).waitFor()
        await page.evaluate(dark => document.documentElement.classList.toggle('dark-theme', dark), dark)
        const search = page.getByRole('textbox', { name: '搜索分组名称', exact: true }).first()
        const filtered = page.waitForResponse(response => response.url().includes('/admin/provider-hall/groups?') && new URL(response.url()).searchParams.get('search') === groupName && response.ok())
        await search.fill(groupName)
        await filtered
        await page.screenshot({ path: resolve(artifacts, `groups-${width}-${dark ? 'dark' : 'light'}.png`), fullPage: true })

        // The targets panel is a disclosure row, not a dialog. The desktop
        // toggle and the mobile card button share a test id but not a label,
        // so select on the id.
        const expand = page.locator('[data-test="row-expand"]').first()
        await expand.waitFor()
        await expand.click()
        const panel = page.locator('[data-expanded-for]').first()
        await panel.waitFor()
        await panel.getByRole('button', { name: '预检查', exact: true }).click()
        await page.waitForTimeout(150)
        await page.screenshot({ path: resolve(artifacts, `editor-${width}-${dark ? 'dark' : 'light'}.png`), fullPage: true })
        const overflow = await page.evaluate(() => ({ page: document.documentElement.scrollWidth > innerWidth + 2,
          dialog: [...document.querySelectorAll('.modal-content')].some(e => e.getBoundingClientRect().right > innerWidth + 2 || e.getBoundingClientRect().left < -2) }))
        assert.deepEqual(overflow, { page: false, dialog: false }, `overflow at ${width}, dark=${dark}`)

        if (width === 1440 && !dark) {
          // Issue (or reuse) the dedicated key, then save the group and probe.
          const keyResponse = page.waitForResponse(r => r.url().endsWith('/probe-keys') && r.request().method() === 'POST')
          await panel.getByRole('button', { name: '创建或复用', exact: true }).click()
          assert.equal((await keyResponse).status(), 200)
          const settingsResponse = page.waitForResponse(r => r.url().endsWith('/settings') && r.request().method() === 'PUT')
          await panel.getByRole('button', { name: '保存本组', exact: true }).click()
          assert.equal((await settingsResponse).status(), 200)

          // The target is enabled and clean now, so the button is a plain probe.
          const jobResponse = page.waitForResponse(r => r.url().endsWith('/probes') && r.request().method() === 'POST')
          await panel.getByRole('button', { name: '立即探测', exact: true }).click()
          assert.ok((await jobResponse).ok())

          // Listing metadata is saved through its own dialog.
          await panel.getByRole('button', { name: '查看结果', exact: true }).click()
          await page.waitForURL(url => url.searchParams.get('tab') === 'jobs' && url.searchParams.has('profile'))
          await page.locator('#hall-panel-jobs').getByRole('button', { name: groupName, exact: true }).first().click()
          await page.waitForURL(url => url.searchParams.get('tab') === 'groups' && url.searchParams.has('group'))
          await page.getByRole('button', { name: '管理', exact: true }).first().click()
          const dialog = page.getByRole('dialog').filter({ has: page.locator('#hall-listing') })
          await dialog.waitFor()
          await dialog.locator('textarea').fill('Saved and tested from the browser')
          await dialog.getByRole('button', { name: '保存', exact: true }).click()
          await page.waitForTimeout(150)
          await page.getByRole('button', { name: '管理', exact: true }).first().click()
          await dialog.waitFor()
          assert.equal(await dialog.locator('textarea').inputValue(), 'Saved and tested from the browser')
          await dialog.getByRole('button', { name: '取消', exact: true }).click()
        }

        for (const label of ['模型档案', '全局配置', '任务', '运行状态']) {
          await page.getByRole('tab', { name: label, exact: true }).click()
          await page.waitForTimeout(150)
          await page.screenshot({ path: resolve(artifacts, `${label}-${width}-${dark ? 'dark' : 'light'}.png`), fullPage: true })
          assert.equal(await page.evaluate(() => document.documentElement.scrollWidth > innerWidth + 2), false, `${label} overflows at ${width}`)
        }
        await context.close()
      }
    }
    assert.deepEqual(errors, [], 'browser exceptions')
  } finally { await browser.close() }
}

if (process.argv[1] && import.meta.url === pathToFileURL(resolve(process.argv[1])).href) {
  assert.ok(process.env.PROVIDER_HALL_UI_URL && process.env.PROVIDER_HALL_TEST_TOKEN && process.env.PROVIDER_HALL_TEST_USER, 'Use the isolated test:e2e:provider-hall runner, or supply UI URL, test token and test user JSON')
  await verifyAdmin({ origin: process.env.PROVIDER_HALL_UI_URL, token: process.env.PROVIDER_HALL_TEST_TOKEN, user: JSON.parse(process.env.PROVIDER_HALL_TEST_USER), artifacts: process.env.PROVIDER_HALL_ARTIFACTS || '/tmp/provider-hall-admin-screenshots', groupName: process.env.PROVIDER_HALL_TEST_GROUP || 'Hall E2E' })
}
