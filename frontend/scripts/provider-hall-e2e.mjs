import assert from 'node:assert/strict'
import { spawn } from 'node:child_process'
import { once } from 'node:events'
import { createWriteStream } from 'node:fs'
import { mkdtemp, mkdir, writeFile, rm } from 'node:fs/promises'
import { createServer } from 'node:http'
import { createServer as createTCPServer } from 'node:net'
import { tmpdir } from 'node:os'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { randomBytes } from 'node:crypto'
import { verifyAdmin } from './provider-hall-admin-smoke.mjs'

const root = resolve(dirname(fileURLToPath(import.meta.url)), '../..')
const work = await mkdtemp(resolve(tmpdir(), 'hall-e2e-'))
const artifacts = process.env.PROVIDER_HALL_ARTIFACTS || resolve(work, 'artifacts')
await mkdir(artifacts, { recursive: true })
const pgBin = process.env.PROVIDER_HALL_PG_BIN || '/opt/homebrew/opt/postgresql@18/bin'
const children = []
let pgStarted = false
const password = randomBytes(24).toString('hex')
let token = '', paidRequests = 0, modelLists = 0, accountProbes = 0
const log = createWriteStream(resolve(artifacts, 'processes.log'))
function start(command, args, options = {}) {
  const child = spawn(command, args, { cwd: root, env: process.env, ...options, stdio: ['ignore', 'pipe', 'pipe'] })
  child.stdout.pipe(log, { end: false }); child.stderr.pipe(log, { end: false })
  children.push(child); return child
}
async function run(command, args, options) { const child = start(command, args, options); const [code] = await once(child, 'exit'); assert.equal(code, 0, `${command} failed; see ${artifacts}/processes.log`) }
async function freePort() { const server = createTCPServer(); server.listen(0, '127.0.0.1'); await once(server, 'listening'); const port = server.address().port; await new Promise(resolve => server.close(resolve)); return port }
const [pgPort, redisPort, apiPort, uiPort] = await Promise.all([freePort(), freePort(), freePort(), freePort()])
const apiOrigin = `http://127.0.0.1:${apiPort}`, uiOrigin = `http://127.0.0.1:${uiPort}`
async function waitFor(fn, message, timeout = 60000) { const start = Date.now(); while (Date.now() - start < timeout) { try { const value = await fn(); if (value) return value } catch { /* startup and task polling */ } await new Promise(resolve => setTimeout(resolve, 250)) } throw Error(message) }
async function api(path, method = 'GET', body, expected = 200) {
  const response = await fetch(`${apiOrigin}/api/v1${path}`, { method, headers: { 'Content-Type': 'application/json', ...(token ? { Authorization: `Bearer ${token}` } : {}) }, body: body === undefined ? undefined : JSON.stringify(body) })
  const json = await response.json()
  assert.ok(response.status === expected || expected === 200 && [201, 202].includes(response.status), `${method} ${path}: ${response.status} ${JSON.stringify(json)}`)
  return json.data ?? json
}
const upstream = createServer(async (req, res) => {
  if (req.url.endsWith('/models')) { modelLists++; res.writeHead(200, { 'Content-Type': 'application/json' }); res.end(JSON.stringify({ data: [{ id: 'gpt-4.1-mini', object: 'model' }] })); return }
  let raw = ''; for await (const chunk of req) raw += chunk
  const body = JSON.parse(raw)
  if (body.stream === false) { accountProbes++; res.writeHead(200, { 'Content-Type': 'application/json' }); res.end(JSON.stringify({ id: 'resp_account_probe', object: 'response', model: body.model, status: 'completed', output: [{ type: 'function_call', name: body.tools?.[0]?.name || 'test', call_id: 'call_test', arguments: '{}' }], usage: { input_tokens: 1, output_tokens: 1 } })); return }
  paidRequests++
  res.writeHead(200, { 'Content-Type': 'text/event-stream' })
  const response = { id: `resp_hall_${paidRequests}`, object: 'response', model: body.model, status: 'completed', output: [{ type: 'message', role: 'assistant', content: [{ type: 'output_text', text: '2' }] }], usage: { input_tokens: 12, output_tokens: 1, total_tokens: 13 } }
  res.write(`event: response.output_text.delta\ndata: ${JSON.stringify({ type: 'response.output_text.delta', delta: '2' })}\n\n`)
  res.end(`event: response.completed\ndata: ${JSON.stringify({ type: 'response.completed', response })}\n\n`)
})
try {
  upstream.listen(0, '127.0.0.1'); await once(upstream, 'listening')
  await run(resolve(pgBin, 'initdb'), ['-D', resolve(work, 'pg'), '-U', 'hall_test', '--auth=trust', '--no-locale', '-E', 'UTF8'])
  await run(resolve(pgBin, 'pg_ctl'), ['-D', resolve(work, 'pg'), '-l', resolve(artifacts, 'postgres.log'), '-o', `-h 127.0.0.1 -k ${work} -p ${pgPort}`, '-w', 'start']); pgStarted = true
  await run(resolve(pgBin, 'createdb'), ['-h', '127.0.0.1', '-p', String(pgPort), '-U', 'hall_test', 'hall_test'])
  start(process.env.PROVIDER_HALL_REDIS_BIN || '/opt/homebrew/bin/redis-server', ['--bind', '127.0.0.1', '--port', String(redisPort), '--save', '', '--appendonly', 'no', '--dir', work])
  await run('go', ['build', '-o', resolve(work, 'server'), './cmd/server'], { cwd: resolve(root, 'backend') })
  const data = resolve(work, 'data'); await mkdir(data)
  start(resolve(work, 'server'), [], { cwd: resolve(root, 'backend'), env: { ...process.env, DATA_DIR: data, AUTO_SETUP: 'true', RUN_MODE: 'standard', INSTANCE_ID: 'hall-e2e', INSTANCE_ROLE: 'primary', DATABASE_HOST: '127.0.0.1', DATABASE_PORT: String(pgPort), DATABASE_USER: 'hall_test', DATABASE_PASSWORD: '', DATABASE_DBNAME: 'hall_test', DATABASE_SSLMODE: 'disable', REDIS_HOST: '127.0.0.1', REDIS_PORT: String(redisPort), REDIS_PASSWORD: '', SERVER_HOST: '127.0.0.1', SERVER_PORT: String(apiPort), ADMIN_EMAIL: 'hall-admin@example.test', ADMIN_PASSWORD: password, PROVIDER_HALL_ALLOW_LOOPBACK: '1' } })
  await waitFor(async () => (await fetch(`${apiOrigin}/health`)).ok, 'API startup failed', 120000)
  const auth = await api('/auth/login', 'POST', { email: 'hall-admin@example.test', password }); token = auth.access_token
  const operator = await api('/admin/users', 'POST', { email: 'hall-operator@example.test', password, balance: 100, concurrency: 5 })
  const group = await api('/admin/groups', 'POST', { name: 'Hall E2E', platform: 'openai', rate_multiplier: 1, subscription_type: 'standard' })
  await api('/admin/accounts', 'POST', { name: 'Hall mock upstream', platform: 'openai', type: 'apikey', credentials: { api_key: 'test-only', base_url: `http://127.0.0.1:${upstream.address().port}`, model_mapping: { 'gpt-4.1-mini': 'gpt-4.1-mini' } }, concurrency: 5, priority: 1, group_ids: [group.id] })
  await waitFor(() => accountProbes === 1, 'Account capability probe did not finish')
  const hall = '/admin/provider-hall'
  let config = await api(`${hall}/config`)
  const configBody = { version: config.version, collection_enabled: true, display_enabled: false, tasks_enabled: true, auto_schedule_enabled: false, default_model: '', default_protocol: 'responses', default_range: '6h', gateway_origin: apiOrigin, operator_user_id: operator.id, daily_budget: '1', expected_nodes: ['hall-e2e'] }
  await api(`${hall}/config/preflight`, 'POST', configBody)
  config = await api(`${hall}/config`, 'PUT', configBody)
  const candidates = await api(`${hall}/groups/${group.id}/models`); assert.ok(candidates.some(m => m.model === 'gpt-4.1-mini' && m.available))
  assert.equal(modelLists, 0); assert.equal(paidRequests, 0)
  await api(`${hall}/groups/${group.id}/models/refresh`, 'POST'); assert.equal(modelLists, 1); assert.equal(paidRequests, 0)
  await api(`${hall}/config/check-gateway`, 'POST', { origin: apiOrigin }); assert.equal(paidRequests, 0)
  const profile = await api(`${hall}/profiles`, 'POST', { version: 0, model: 'gpt-4.1-mini', protocol: 'responses', supports_tools: false, output_limit: 32, model_aliases: [], reference_input_price: null, reference_cache_price: null, reference_cache_rate: null, reference_confirmed_at: null })
  const keys = await Promise.all(Array.from({ length: 5 }, () => api(`${hall}/groups/${group.id}/probe-keys`, 'POST', { profile_id: profile.id })))
  assert.equal(new Set(keys.map(k => k.id)).size, 1); assert.ok(keys.every(k => !('key' in k)))
  let settings = { version: 0, listed: false, display_name: 'Hall E2E display', description: 'Local test', display_order: 0, items: [{ profile_id: profile.id, probe_key_id: keys[0].id, enabled: true, auto_schedule_enabled: false, probe_interval_seconds: 300, verification_interval_seconds: 86400 }] }
  await api(`${hall}/groups/${group.id}/preflight`, 'POST', settings); assert.equal(paidRequests, 0)
  assert.equal((await api(`${hall}/groups/${group.id}/targets`)).items.length, 0)
  const saved = await api(`${hall}/groups/${group.id}/settings`, 'PUT', settings)
  await api(`${hall}/groups/${group.id}/settings`, 'PUT', settings, 409)
  settings.version = saved.version
  const before = paidRequests
  const enqueue = { profile_id: profile.id, target_version: saved.items[0].version, profile_version: profile.version, idempotency_key: 'hall-e2e-manual' }
  await api(`${hall}/groups/${group.id}/probes`, 'POST', { ...enqueue, target_version: 9999 }, 409)
  const queued = await api(`${hall}/groups/${group.id}/probes`, 'POST', enqueue)
  const retry = await api(`${hall}/groups/${group.id}/probes`, 'POST', enqueue); assert.equal(queued.job_id, retry.job_id)
  const detail = await waitFor(async () => { const d = await api(`${hall}/jobs/${queued.job_id}`); return ['succeeded', 'failed', 'cancelled'].includes(d.job.status) && d }, 'manual test did not finish', 60000)
  assert.equal(detail.job.status, 'succeeded', JSON.stringify(detail)); assert.equal(paidRequests, before + 1)
  const filtered = await api(`${hall}/jobs?group_name=Hall%20E2E&model=gpt-4.1&source=manual`); assert.equal(filtered.total, 1)
  const batch = await api(`${hall}/groups/batch`, 'POST', { preview: true, groups: [{ ...settings, id: group.id, listed: true }] }); assert.equal(batch[0].success, true)
  assert.equal((await api(`${hall}/groups/${group.id}`)).listed, false)
  const executed = await api(`${hall}/groups/batch`, 'POST', { preview: false, groups: [{ ...settings, id: group.id, listed: true }] }); assert.equal(executed[0].success, true)
  assert.equal(paidRequests, before + 1)
  const secondGroup = await api('/admin/groups', 'POST', { name: 'Hall second batch', platform: 'openai', rate_multiplier: 1, subscription_type: 'standard' })
  const partial = await api(`${hall}/groups/batch`, 'POST', { preview: false, groups: [
    { ...settings, id: group.id, listed: false },
    { id: secondGroup.id, version: 0, listed: true, display_name: '', description: '', display_order: 0, items: [] },
  ] })
  assert.deepEqual(partial.map(r => r.success), [false, true])
  assert.equal(partial[0].reason, 'PROVIDER_HALL_VERSION_CONFLICT')
  assert.equal((await api(`${hall}/groups/${group.id}`)).listed, true)
  assert.equal((await api(`${hall}/groups/${secondGroup.id}`)).listed, true)
  start(resolve(root, 'frontend/node_modules/.bin/vite'), ['--host', '127.0.0.1', '--port', String(uiPort), '--strictPort'], { cwd: resolve(root, 'frontend'), env: { ...process.env, VITE_DEV_PROXY_TARGET: apiOrigin } })
  await waitFor(async () => (await fetch(uiOrigin)).ok, 'Vite startup failed')
  await verifyAdmin({ origin: uiOrigin, token, user: auth.user, artifacts, groupName: 'Hall E2E' })
  // Real probe facts must retain their provenance and stay out of user metrics.
  await run(resolve(pgBin, 'psql'), ['-h', '127.0.0.1', '-p', String(pgPort), '-U', 'hall_test', '-d', 'hall_test', '-v', 'ON_ERROR_STOP=1', '-c', `DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM provider_hall_requests WHERE group_id=${group.id} AND source='probe' AND sample_id IS NOT NULL) THEN RAISE EXCEPTION 'missing probe fact'; END IF;
    IF EXISTS (SELECT 1 FROM provider_hall_requests WHERE group_id=${group.id} AND source='user') THEN RAISE EXCEPTION 'probe became user traffic'; END IF;
    IF EXISTS (SELECT 1 FROM provider_hall_metrics_1m WHERE group_id=${group.id} AND (success_count>0 OR failed_count>0 OR ttft_sample_count>0)) THEN RAISE EXCEPTION 'probe polluted user metrics'; END IF;
  END $$;`])
  await writeFile(resolve(artifacts, 'result.json'), JSON.stringify({ passed: true, paidRequests, modelLists, accountProbes, widths: [390, 1440, 1920], themes: ['light', 'dark'], runMode: 'standard' }, null, 2))
  console.log(`Provider Hall E2E passed. Evidence: ${artifacts}`)
} catch (error) {
  console.error(error)
  console.error(`Evidence: ${artifacts}`)
  process.exitCode = 1
} finally {
  for (const child of children.reverse()) if (child.exitCode === null && !child.killed) child.kill('SIGTERM')
  await Promise.all(children.filter(c => c.exitCode === null).map(c => Promise.race([once(c, 'exit'), new Promise(resolve => setTimeout(resolve, 5000))])))
  for (const child of children) if (child.exitCode === null) child.kill('SIGKILL')
  if (pgStarted) await run(resolve(pgBin, 'pg_ctl'), ['-D', resolve(work, 'pg'), '-m', 'immediate', '-w', 'stop'])
  upstream.closeAllConnections(); upstream.close()
  log.end()
  // Keep logs and screenshots; ephemeral data contains only generated test credentials.
  for (const name of ['pg', 'data', 'server']) await rm(resolve(work, name), { recursive: true, force: true })
}
