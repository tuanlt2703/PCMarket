import { getDb, commit, nextId, clone, nowIso, delay, ApiError, notFound } from './mockDb'
import { CUSTOMER_TIERS } from '@/config/statuses'

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const PHONE_RE = /^0\d{9,10}$/

function validate(db, payload, id = null) {
  const errors = {}
  const full_name = String(payload.full_name ?? '').trim()
  const email = String(payload.email ?? '').trim().toLowerCase()
  const phone = String(payload.phone ?? '').replace(/\s/g, '')
  if (!full_name) errors.full_name = 'Full name is required'
  else if (full_name.length > 150) errors.full_name = 'Max 150 characters'
  if (!EMAIL_RE.test(email)) errors.email = 'Invalid email'
  else if (db.customers.some((c) => c.email === email && c.id !== id)) errors.email = 'Email already exists'
  if (phone && !PHONE_RE.test(phone)) errors.phone = 'Invalid phone number'
  if (!CUSTOMER_TIERS[payload.tier]) errors.tier = 'Invalid tier'
  const points = Number(payload.loyalty_points ?? 0)
  if (!Number.isInteger(points) || points < 0) errors.loyalty_points = 'Must be a non-negative integer'
  if (Object.keys(errors).length) throw new ApiError(422, 'Validation failed', errors)
  return {
    full_name,
    email,
    phone: phone || null,
    address: payload.address?.trim() || null,
    tier: payload.tier,
    loyalty_points: points,
    is_active: payload.is_active !== false,
  }
}

function withStats(db, c) {
  const orders = db.orders.filter((o) => o.customer_id === c.id)
  return {
    ...c,
    order_count: orders.length,
    total_spent: orders.filter((o) => o.status === 'DELIVERED').reduce((s, o) => s + o.total, 0),
  }
}

export const customerApi = {
  async list({ q, tier } = {}) {
    await delay()
    const db = getDb()
    const term = q?.trim().toLowerCase()
    return clone(
      db.customers
        .filter((c) => !tier || c.tier === tier)
        .filter(
          (c) =>
            !term ||
            [c.full_name, c.email, c.phone, c.code].some((v) => v?.toLowerCase().includes(term)),
        )
        .map((c) => withStats(db, c))
        .sort((a, b) => b.id - a.id),
    )
  },

  async get(id) {
    await delay()
    const db = getDb()
    const c = db.customers.find((x) => x.id === id)
    if (!c) throw notFound('Customer')
    return clone({
      ...withStats(db, c),
      orders: db.orders.filter((o) => o.customer_id === id).sort((a, b) => b.id - a.id),
      warranties: db.warranties.filter((w) => w.customer_id === id),
    })
  },

  async create(payload) {
    await delay()
    const db = getDb()
    const id = nextId(db.customers)
    const customer = {
      id,
      code: `CUS${String(id).padStart(5, '0')}`,
      ...validate(db, payload),
      created_at: nowIso(),
      updated_at: nowIso(),
    }
    db.customers.push(customer)
    commit()
    return clone(customer)
  },

  async update(id, payload) {
    await delay()
    const db = getDb()
    const c = db.customers.find((x) => x.id === id)
    if (!c) throw notFound('Customer')
    Object.assign(c, validate(db, payload, id), { updated_at: nowIso() })
    commit()
    return clone(c)
  },

  async remove(id) {
    await delay()
    const db = getDb()
    const idx = db.customers.findIndex((x) => x.id === id)
    if (idx < 0) throw notFound('Customer')
    if (db.orders.some((o) => o.customer_id === id)) {
      throw new ApiError(409, 'Cannot delete a customer with orders. Deactivate the account instead.')
    }
    db.customers.splice(idx, 1)
    commit()
  },
}
