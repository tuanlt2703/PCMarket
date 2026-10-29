import { getDb, commit, nextId, clone, nowIso, delay, ApiError, notFound } from './mockDb'
import { reservedQuantity } from './productApi'

function row(db, inv) {
  const model = db.models.find((m) => m.id === inv.model_id)
  const reserved = reservedQuantity(db, inv.model_id)
  return {
    ...inv,
    model_name: model?.model_name ?? '(deleted product)',
    component_type_id: model?.component_type_id ?? null,
    manufacturer_id: model?.manufacturer_id ?? null,
    is_active: model?.is_active ?? false,
    quantity_reserved: reserved,
    quantity_available: inv.quantity_on_hand - reserved,
  }
}

export const inventoryApi = {
  async list() {
    await delay()
    const db = getDb()
    return clone(db.inventory.map((inv) => row(db, inv)))
  },

  async get(modelId) {
    await delay()
    const db = getDb()
    const inv = db.inventory.find((i) => i.model_id === modelId)
    if (!inv) throw notFound('Inventory record')
    return clone(row(db, inv))
  },

  async update(modelId, payload) {
    await delay()
    const db = getDb()
    const inv = db.inventory.find((i) => i.model_id === modelId)
    if (!inv) throw notFound('Inventory record')
    const errors = {}
    const sku = String(payload.sku ?? '').trim().toUpperCase()
    if (!sku) errors.sku = 'SKU is required'
    else if (db.inventory.some((i) => i.sku === sku && i.model_id !== modelId)) errors.sku = 'SKU already exists'
    for (const k of ['selling_price', 'unit_cost', 'reorder_level']) {
      const n = Number(payload[k])
      if (payload[k] === '' || payload[k] === null || Number.isNaN(n) || n < 0) errors[k] = 'Must be ≥ 0'
    }
    if (Object.keys(errors).length) throw new ApiError(422, 'Validation failed', errors)
    Object.assign(inv, {
      sku,
      selling_price: Number(payload.selling_price),
      unit_cost: Number(payload.unit_cost),
      reorder_level: Math.floor(Number(payload.reorder_level)),
      location: payload.location?.trim() || null,
      updated_at: nowIso(),
    })
    commit()
    return clone(row(db, inv))
  },

  // type: IN (receive), OUT (issue), ADJUST (set absolute counted quantity)
  async adjust(modelId, { type, quantity, reason, reference }) {
    await delay()
    const db = getDb()
    const inv = db.inventory.find((i) => i.model_id === modelId)
    if (!inv) throw notFound('Inventory record')
    const qty = Number(quantity)
    const errors = {}
    if (!['IN', 'OUT', 'ADJUST'].includes(type)) errors.type = 'Invalid type'
    if (!Number.isInteger(qty) || qty < 0 || (type !== 'ADJUST' && qty === 0)) errors.quantity = 'Invalid quantity'
    if (type === 'OUT' && qty > inv.quantity_on_hand) errors.quantity = 'Not enough stock on hand'
    if (!reason?.trim()) errors.reason = 'Reason is required'
    if (Object.keys(errors).length) throw new ApiError(422, 'Validation failed', errors)

    const delta = type === 'IN' ? qty : type === 'OUT' ? -qty : qty - inv.quantity_on_hand
    inv.quantity_on_hand += delta
    inv.updated_at = nowIso()
    db.stockMovements.push({
      id: nextId(db.stockMovements),
      model_id: modelId,
      type,
      quantity: delta,
      reason: reason.trim(),
      reference: reference?.trim() || null,
      created_at: nowIso(),
    })
    commit()
    return clone(row(db, inv))
  },

  async movements({ modelId } = {}) {
    await delay()
    const db = getDb()
    return clone(
      db.stockMovements
        .filter((m) => !modelId || m.model_id === modelId)
        .map((m) => ({ ...m, model_name: db.models.find((x) => x.id === m.model_id)?.model_name ?? '(deleted product)' }))
        .sort((a, b) => b.created_at.localeCompare(a.created_at)),
    )
  },
}
