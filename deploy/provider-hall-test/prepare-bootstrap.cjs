const fs = require('node:fs');
const crypto = require('node:crypto');
const path = require('node:path');
const base = '/opt/sub2api-provider-hall-test';
const dir = path.join(base, 'credentials');
const target = path.join(dir, 'login.json');
if (fs.existsSync(target)) throw new Error('Login record already exists; refusing to replace credentials');
const email = `hall-${crypto.randomBytes(8).toString('hex')}@test.example.com`;
const password = crypto.randomBytes(24).toString('base64url') + '!aA7';
const login = {url: 'http://127.0.0.1:18082/login', email, password, created_at: new Date().toISOString()};
fs.writeFileSync(target, JSON.stringify(login, null, 2) + '\n', {mode: 0o600, flag: 'wx'});
const secrets = {
  DATABASE_PASSWORD: fs.readFileSync(path.join(dir, 'db-password'), 'utf8').trim(),
  ADMIN_EMAIL: email,
  ADMIN_PASSWORD: password,
  JWT_SECRET: crypto.randomBytes(32).toString('hex'),
  TOTP_ENCRYPTION_KEY: crypto.randomBytes(32).toString('hex'),
};
fs.writeFileSync(path.join(dir, 'runtime.env'), Object.entries(secrets).map(([k,v]) => `${k}=${v}`).join('\n') + '\n', {mode: 0o600});
console.log('Bootstrap credentials stored in restricted files; secret values suppressed.');
