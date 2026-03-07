import { spawnSync } from 'node:child_process'

const backendBaseArgIndex = process.argv.indexOf('--backend-base')

const backendBase =
  backendBaseArgIndex > -1 && process.argv[backendBaseArgIndex + 1]
    ? process.argv[backendBaseArgIndex + 1].replace(/\/$/, '')
    : process.env.ACCEPTANCE_BACKEND_BASE || 'http://127.0.0.1:18080'

const steps = [
  {
    name: 'frontend-smoke',
    args: ['scripts/frontend-smoke.mjs'],
  },
  {
    name: 'bookings-export-real-smoke',
    args: ['scripts/bookings-export-real-smoke.mjs', '--backend-base', backendBase],
  },
]

function runStep(step) {
  const startedAt = new Date().toISOString()
  const result = spawnSync(process.execPath, step.args, {
    cwd: process.cwd(),
    env: process.env,
    stdio: 'inherit',
  })

  return {
    name: step.name,
    startedAt,
    finishedAt: new Date().toISOString(),
    exitCode: result.status ?? 1,
  }
}

const results = []

for (const step of steps) {
  results.push(runStep(step))
  if (results.at(-1)?.exitCode !== 0) {
    console.error(JSON.stringify({ ok: false, backendBase, results }, null, 2))
    process.exit(results.at(-1)?.exitCode || 1)
  }
}

console.log(JSON.stringify({ ok: true, backendBase, results }, null, 2))