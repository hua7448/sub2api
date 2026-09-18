const fs = require('node:fs');
const assert = require('node:assert/strict');
const base = 'http://127.0.0.1:18082';
async function main() {
  const login = JSON.parse(fs.readFileSync('/opt/sub2api-provider-hall-test/credentials/login.json', 'utf8'));
  const response = await fetch(base + '/api/v1/auth/login', {
    method: 'POST', headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({email: login.email, password: login.password}),
  });
  assert.equal(response.status, 200, 'Login HTTP status');
  const result = await response.json();
  const token = result.data?.access_token;
  assert.ok(token, 'Login must return access token');
  const me = await fetch(base + '/api/v1/auth/me', {headers: {Authorization: `Bearer ${token}`}});
  assert.equal(me.status, 200, 'Authenticated profile HTTP status');
  const profile = (await me.json()).data;
  assert.equal(profile.email, login.email);
  assert.equal(profile.role, 'admin');
  const page = await fetch(base + '/login');
  assert.equal(page.status, 200, 'Login page HTTP status');
  const html = await page.text();
  assert.ok(html.includes('<html'), 'Embedded frontend HTML');
  const assets = [...html.matchAll(/(?:src|href)="([^" ]+\.(?:js|css))"/g)].map(m => m[1]);
  assert.ok(assets.length > 0, 'Frontend asset references');
  for (const asset of assets) {
    const res = await fetch(new URL(asset, base));
    assert.equal(res.status, 200, `Frontend asset ${asset}`);
  }
  console.log(JSON.stringify({login: 'passed', role: profile.role, frontend: 'passed', assets: assets.length, checked_at: new Date().toISOString()}));
}
main().catch(err => {console.error(err.message); process.exitCode = 1;});
