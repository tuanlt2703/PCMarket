import { getDb, commit, nextId, clone, nowIso, delay, ApiError, notFound } from './mockDb'
import { ORDER_TRANSITIONS, PAYMENT_STATUSES } from '@/config/statuses'

function withCustomer(db, o) {
  const c = db.customers.find((x) => x.id === o.customer_id)
  return { ...o, customer_name: c?.full_name ?? '(deleted)', customer_email: c?.email ?? null, customer_phone: c?.phone ?? null }
}

export const orderApi = {
  async list({ status, paymentStatus, q, from, to } = {}) {
    await delay()
    const db = getDb()
    const term = q?.trim().toLowerCase()
    return clone(
      db.orders
        .map((o) => withCustomer(db, o))
        .filter((o) => !status || o.status === status)
        .filter((o) => !paymentStatus || o.payment_status === paymentStatus)
        .filter((o) => !from || o.created_at.slice(0, 10) >= from)
        .filter((o) => !to || o.created_at.slice(0, 10) <= to)
        .filter(
          (o) => !term || o.order_number.toLowerCase().includes(term) || o.customer_name.toLowerCase().includes(term),
        )
        .sort((a, b) => b.created_at.localeCompare(a.created_at)),
    )
  },

  async get(id) {
    await delay()
    const db = getDb()
    const o = db.orders.find((x) => x.id === id)
    if (!o) throw notFound('Order')
    return clone(withCustomer(db, o))
  },

  async updateStatus(id, status, note) {
    await delay(200)
    const db = getDb()
    const o = db.orders.find((x) => x.id === id)
    if (!o) throw notFound('Order')
    if (!ORDER_TRANSITIONS[o.status]?.includes(status)) {
      throw new ApiError(422, `Cannot change status from ${o.status} to ${status}`)
    }

    // Shipping consumes physical stock (reservation is released at the same time)
    if (status === 'SHIPPED') {
      const shortages = o.items.filter((it) => {
        const inv = db.inventory.find((i) => i.model_id === it.model_id)
        return !inv || inv.quantity_on_hand < it.quantity
      })
      if (shortages.length) {
        throw new ApiError(409, `Not enough stock to ship: ${shortages.map((s) => s.product_name).join(', ')}`)
      }
      for (const it of o.items) {
        const inv = db.inventory.find((i) => i.model_id === it.model_id)
        inv.quantity_on_hand -= it.quantity
        inv.updated_at = nowIso()
        db.stockMovements.push({
          id: nextId(db.stockMovements),
          model_id: it.model_id,
          type: 'OUT',
          quantity: -it.quantity,
          reason: 'Order shipped',
          reference: o.order_number,
          created_at: nowIso(),
        })
      }
    }

    o.status = status
    if (status === 'DELIVERED' && o.payment_method === 'COD') o.payment_status = 'PAID'
    o.status_history.push({ status, note: note?.trim() || null, at: nowIso() })
    o.updated_at = nowIso()
    commit()
    return clone(withCustomer(db, o))
  },

  async updatePaymentStatus(id, paymentStatus) {
    await delay()
    const db = getDb()
    const o = db.orders.find((x) => x.id === id)
    if (!o) throw notFound('Order')
    if (!PAYMENT_STATUSES[paymentStatus]) throw new ApiError(422, 'Invalid payment status')
    o.payment_status = paymentStatus
    o.updated_at = nowIso()
    commit()
    return clone(withCustomer(db, o))
  },
}
