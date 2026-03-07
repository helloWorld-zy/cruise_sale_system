import { execFileSync } from 'node:child_process'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

const backendBaseArgIndex = process.argv.indexOf('--backend-base')
const usernameArgIndex = process.argv.indexOf('--username')
const passwordArgIndex = process.argv.indexOf('--password')

const scriptDir = dirname(fileURLToPath(import.meta.url))
const dockerDir = resolve(scriptDir, '../docker')
const postgresUser = process.env.SMOKE_PGUSER || 'cruise'
const postgresDatabase = process.env.SMOKE_PGDATABASE || 'cruise_booking'

const backendBase =
  backendBaseArgIndex > -1 && process.argv[backendBaseArgIndex + 1]
    ? process.argv[backendBaseArgIndex + 1].replace(/\/$/, '')
    : 'http://127.0.0.1:8080'

const username =
  usernameArgIndex > -1 && process.argv[usernameArgIndex + 1]
    ? process.argv[usernameArgIndex + 1]
    : 'admin'

const password =
  passwordArgIndex > -1 && process.argv[passwordArgIndex + 1]
    ? process.argv[passwordArgIndex + 1]
    : 'admin123'

function sqlLiteral(value) {
  if (value === null || value === undefined) {
    return 'NULL'
  }
  return `'${String(value).replace(/'/g, "''")}'`
}

function runPsql(sql) {
  return execFileSync(
    'docker',
    [
      'compose',
      'exec',
      '-T',
      'postgres',
      'psql',
      '-U',
      postgresUser,
      '-d',
      postgresDatabase,
      '-At',
      '-F',
      '\t',
      '-c',
      sql,
    ],
    {
      cwd: dockerDir,
      encoding: 'utf8',
      stdio: ['ignore', 'pipe', 'pipe'],
    },
  ).trim()
}

function parseFirstRow(output) {
  const firstLine = output
    .split(/\r?\n/)
    .map((line) => line.trim())
    .find(Boolean)

  if (!firstLine) {
    return []
  }

  return firstLine.split('\t').map((part) => part.trim())
}

function parseRequiredId(value, fieldName) {
  const parsed = Number(value)
  if (!Number.isFinite(parsed) || parsed <= 0) {
    throw new Error(`fixture query returned invalid ${fieldName}: ${value}`)
  }
  return parsed
}

function prepareFixture() {
  const basis = parseFirstRow(
    runPsql(`
      SELECT u.id, v.id, ct.id
      FROM users u
      JOIN voyages v ON v.status = 1
      JOIN cabin_types ct
        ON ct.cruise_id = v.cruise_id
       AND ct.status = 1
       AND ct.deleted_at IS NULL
      WHERE u.status = 1
      ORDER BY u.id, v.id, ct.id
      LIMIT 1;
    `),
  )

  if (basis.length < 3) {
    throw new Error('cannot prepare fixture: missing active user/voyage/cabin type base data')
  }

  const [userIdValue, voyageIdValue, cabinTypeIdValue] = basis
  const userId = parseRequiredId(userIdValue, 'user_id')
  const voyageId = parseRequiredId(voyageIdValue, 'voyage_id')
  const cabinTypeId = parseRequiredId(cabinTypeIdValue, 'cabin_type_id')
  const fixtureCode = `SMOKE-REAL-${Date.now()}`
  const occupancy = 2
  const priceCents = 199900

  const skuId = parseRequiredId(
    parseFirstRow(runPsql(`
      INSERT INTO cabin_skus (
        voyage_id,
        cabin_type_id,
        code,
        max_guests,
        status,
        created_at,
        updated_at
      )
      VALUES (
        ${voyageId},
        ${cabinTypeId},
        ${sqlLiteral(fixtureCode)},
        ${occupancy},
        1,
        NOW(),
        NOW()
      )
      RETURNING id;
    `))[0],
    'sku_id',
  )

  runPsql(`
    INSERT INTO cabin_inventories (
      cabin_sku_id,
      total,
      locked,
      sold,
      alert_threshold,
      updated_at
    )
    VALUES (
      ${skuId},
      1,
      0,
      0,
      0,
      NOW()
    );

    INSERT INTO cabin_prices (
      cabin_sku_id,
      date,
      occupancy,
      price_cents,
      child_price_cents,
      single_supplement_cents,
      price_type,
      created_at,
      updated_at
    )
    VALUES (
      ${skuId},
      NOW(),
      ${occupancy},
      ${priceCents},
      0,
      0,
      'base',
      NOW(),
      NOW()
    );
  `)

  return { userId, voyageId, cabinTypeId, skuId, fixtureCode, occupancy, priceCents }
}

function cleanupFixture(fixture, bookingId) {
  if (!fixture?.skuId) {
    return
  }

  const statements = []

  if (bookingId > 0) {
    statements.push(`DELETE FROM booking_passengers WHERE booking_id = ${bookingId};`)
    statements.push(`DELETE FROM order_status_logs WHERE order_id = ${bookingId};`)
    statements.push(`DELETE FROM refunds WHERE payment_id IN (SELECT id FROM payments WHERE order_id = ${bookingId});`)
    statements.push(`DELETE FROM payments WHERE order_id = ${bookingId};`)
    statements.push(`DELETE FROM bookings WHERE id = ${bookingId};`)
  }

  statements.push(`DELETE FROM cabin_holds WHERE cabin_sku_id = ${fixture.skuId};`)
  statements.push(`DELETE FROM inventory_logs WHERE cabin_sku_id = ${fixture.skuId};`)
  statements.push(`DELETE FROM cabin_prices WHERE cabin_sku_id = ${fixture.skuId};`)
  statements.push(`DELETE FROM cabin_inventories WHERE cabin_sku_id = ${fixture.skuId};`)
  statements.push(`DELETE FROM cabin_skus WHERE id = ${fixture.skuId};`)

  runPsql(statements.join('\n'))
}

async function requestJson(url, options = {}) {
  const response = await fetch(url, options)
  const text = await response.text()
  let payload = null
  try {
    payload = text ? JSON.parse(text) : null
  } catch {
    payload = text
  }
  return { response, payload, text }
}

async function login() {
  const { response, payload, text } = await requestJson(`${backendBase}/api/v1/admin/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  })

  if (!response.ok) {
    throw new Error(`login failed: HTTP ${response.status} ${text}`)
  }

  const token = payload?.data?.token || payload?.token || ''
  if (!token) {
    throw new Error('login response missing token')
  }
  return token
}

async function createBooking(token, fixture) {
  const { response, payload, text } = await requestJson(`${backendBase}/api/v1/admin/bookings`, {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      user_id: fixture.userId,
      voyage_id: fixture.voyageId,
      cabin_sku_id: fixture.skuId,
      guests: fixture.occupancy,
    }),
  })

  if (!response.ok) {
    throw new Error(`create booking failed: HTTP ${response.status} ${text}`)
  }

  const bookingId = Number(payload?.data?.id || payload?.id || 0)
  if (!Number.isFinite(bookingId) || bookingId <= 0) {
    throw new Error(`create booking returned invalid id: ${text}`)
  }
  return bookingId
}

async function exportBooking(token, bookingId) {
  const exportUrl = new URL(`${backendBase}/api/v1/admin/bookings/export`)
  exportUrl.searchParams.set('booking_no', String(bookingId))
  exportUrl.searchParams.set('keyword', String(bookingId))

  const response = await fetch(exportUrl, {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })

  const contentType = response.headers.get('content-type') || ''
  const text = await response.text()

  if (!response.ok) {
    throw new Error(`export failed: HTTP ${response.status} ${text}`)
  }
  if (!contentType.includes('text/csv')) {
    throw new Error(`export content-type mismatch: ${contentType}`)
  }
  if (!text.includes('booking_no,phone,voyage_code,cruise_name')) {
    throw new Error(`export header mismatch: ${text}`)
  }
  if (!text.includes(String(bookingId))) {
    throw new Error(`export content missing booking id ${bookingId}: ${text}`)
  }

  return { exportUrl: exportUrl.toString(), contentType, text }
}

async function deleteBooking(token, bookingId) {
  const response = await fetch(`${backendBase}/api/v1/admin/bookings/${bookingId}`, {
    method: 'DELETE',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })

  if (response.status !== 204) {
    const text = await response.text()
    throw new Error(`delete booking failed: HTTP ${response.status} ${text}`)
  }
}

async function main() {
  let createdBookingId = 0
  let token = ''
  let fixture = null
  const startedAt = new Date().toISOString()

  try {
    fixture = prepareFixture()
    token = await login()
    createdBookingId = await createBooking(token, fixture)
    const exported = await exportBooking(token, createdBookingId)
    await deleteBooking(token, createdBookingId)
    cleanupFixture(fixture, createdBookingId)

    console.log(JSON.stringify({
      ok: true,
      backendBase,
      startedAt,
      finishedAt: new Date().toISOString(),
      bookingId: createdBookingId,
      fixture,
      exportUrl: exported.exportUrl,
      contentType: exported.contentType,
      preview: exported.text.split('\n').slice(0, 2),
    }, null, 2))
  } catch (error) {
    const message = String(error?.message || error)
    if (token && createdBookingId > 0) {
      try {
        await deleteBooking(token, createdBookingId)
      } catch {}
    }
    if (fixture?.skuId) {
      try {
        cleanupFixture(fixture, createdBookingId)
      } catch {}
    }
    console.error(JSON.stringify({
      ok: false,
      backendBase,
      startedAt,
      finishedAt: new Date().toISOString(),
      bookingId: createdBookingId,
      fixture,
      error: message,
    }, null, 2))
    process.exit(1)
  }
}

main()