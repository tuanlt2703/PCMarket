import { getDb, commit, nextId, clone, nowIso, delay, ApiError, notFound } from './mockDb'
import { WARRANTY_STATUSES } from '@/config/statuses'

function enrich(db, w) {
  return {
    ...w,
    customer_name: db.customers.find((c) => c.id === w.customer_id)?.full_name ?? '(deleted)',
    order_number: db.orders.find((o) => o.id === w.order_id)?.order_number ?? null,
    model_name: db.models.find((m) => m.id === w.model_id)?.model_name ?? '(deleted product)',
  }
}

function validate(db, payload) {
  const errors = {}
  const order = db.orders.find((o) => o.id === Number(payload.order_id))
  if (!order) errors.order_id = 'Order is required'
  else if (order.status !== 'DELIVERED') errors.order_id = 'Only delivered orders can have warranty claims'
  const modelId = Number(payload.model_id)
  if (!modelId) errors.model_id = 'Product is required'
  else if (order && !order.items.some((it) => it.model_id === modelId)) errors.model_id = 'Product is not in this order'
  if (!payload.serial_number?.trim()) errors.serial_number = 'Serial number is required'
  if (!payload.issue_description?.trim()) errors.issue_description = 'Issue description is required'
  if (!WARRANTY_STATUSES[payload.status]) errors.status = 'Invalid status'
  if (Object.keys(errors).length) throw new ApiError(422, 'Validation failed', errors)
  return { order, modelId }
}

export const warrantyApi = {
  async list({ status, q } = {}) {
    await delay()
    const db = getDb()
    const term = q?.trim().toLowerCase()
    return clone(
      db.warranties
        .map((w) => enrich(db, w))
        .filter((w) => !status || w.status === status)
        .filter(
          (w) =>
            !term ||
            [w.claim_number, w.customer_name, w.model_name, w.serial_number].some((v) => v?.toLowerCase().includes(term)),
        )
        .sort((a, b) => b.created_at.localeCompare(a.created_at)),
    )
  },

  async create(payload) {
    await delay()
    const db = getDb()
    const { order, modelId } = validate(db, payload)
    const id = nextId(db.warranties)
    const claim = {
      id,
      claim_number: `WAR-${String(id).padStart(5, '0')}`,
      order_id: order.id,
      customer_id: order.customer_id,
      model_id: modelId,
      serial_number: payload.serial_number.trim(),
      purchase_date: order.created_at.slice(0, 10),
      warranty_expiry: payload.warranty_expiry || null,
      issue_description: payload.issue_description.trim(),
      status: payload.status,
      resolution: payload.resolution?.trim() || null,
      created_at: nowIso(),
      updated_at: nowIso(),
    }
    db.warranties.push(claim)
    commit()
    return clone(enrich(db, claim))
  },

  async update(id, payload) {
    await delay()
    const db = getDb()
    const w = db.warranties.find((x) => x.id === id)
    if (!w) throw notFound('Warranty claim')
    const { order, modelId } = validate(db, payload)
    Object.assign(w, {
      order_id: order.id,
      customer_id: order.customer_id,
      model_id: modelId,
      serial_number: payload.serial_number.trim(),
      warranty_expiry: payload.warranty_expiry || null,
      issue_description: payload.issue_description.trim(),
      status: payload.status,
      resolution: payload.resolution?.trim() || null,
      updated_at: nowIso(),
    })
    commit()
    return clone(enrich(db, w))
  },

  async remove(id) {
    await delay()
    const db = getDb()
    const idx = db.warranties.findIndex((x) => x.id === id)
    if (idx < 0) throw notFound('Warranty claim')
    db.warranties.splice(idx, 1)
    commit()
  },
}
