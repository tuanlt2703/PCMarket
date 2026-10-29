import { getDb, commit, nextId, clone, nowIso, delay, ApiError, notFound } from './mockDb'
import { getSchema, validateSpec, normalizeSpec, emptySpec } from '@/config/componentSchemas'
import { RESERVING_ORDER_STATUSES } from '@/config/statuses'

const typeCodeOf = (db, typeId) => db.lookups.component_types.find((t) => t.id === typeId)?.code

export function reservedQuantity(db, modelId) {
  return db.orders
    .filter((o) => RESERVING_ORDER_STATUSES.includes(o.status))
    .flatMap((o) => o.items)
    .filter((it) => it.model_id === modelId)
    .reduce((s, it) => s + it.quantity, 0)
}

function withStock(db, model) {
  const inv = db.inventory.find((i) => i.model_id === model.id)
  const reserved = reservedQuantity(db, model.id)
  return {
    ...model,
    sku: inv?.sku ?? null,
    selling_price: inv?.selling_price ?? null,
    quantity_on_hand: inv?.quantity_on_hand ?? 0,
    quantity_available: (inv?.quantity_on_hand ?? 0) - reserved,
  }
}

function validateModel(db, payload, id = null) {
  const errors = {}
  const component_type_id = Number(payload.component_type_id) || null
  const manufacturer_id = Number(payload.manufacturer_id) || null
  const model_name = String(payload.model_name ?? '').trim()

  if (!component_type_id) errors.component_type_id = 'Component type is required'
  if (!model_name) errors.model_name = 'Model name is required'
  else if (model_name.length > 200) errors.model_name = 'Max 200 characters'
  if (
    model_name &&
    db.models.some(
      (m) =>
        m.id !== id &&
        m.component_type_id === component_type_id &&
        m.manufacturer_id === manufacturer_id &&
        m.model_name.toLowerCase() === model_name.toLowerCase(),
    )
  ) {
    errors.model_name = 'A model with the same type, manufacturer and name already exists'
  }

  const schema = getSchema(typeCodeOf(db, component_type_id))
  const specInput = { ...emptySpec(schema), ...(payload.spec ?? {}) }
  for (const [k, v] of Object.entries(validateSpec(schema, specInput))) errors[`spec.${k}`] = v

  if (Object.keys(errors).length) throw new ApiError(422, 'Validation failed', errors)

  return {
    component_type_id,
    manufacturer_id,
    model_name,
    description: payload.description?.trim() || null,
    release_date: payload.release_date || null,
    is_active: payload.is_active !== false,
    spec: schema ? normalizeSpec(schema, specInput) : {},
  }
}

export const productApi = {
  async list({ componentTypeId, manufacturerId, q, active } = {}) {
    await delay()
    const db = getDb()
    const term = q?.trim().toLowerCase()
    return clone(
      db.models
        .filter((m) => !componentTypeId || m.component_type_id === componentTypeId)
        .filter((m) => !manufacturerId || m.manufacturer_id === manufacturerId)
        .filter((m) => active === undefined || active === null || m.is_active === active)
        .map((m) => withStock(db, m))
        .filter((m) => !term || m.model_name.toLowerCase().includes(term) || m.sku?.toLowerCase().includes(term))
        .sort((a, b) => b.id - a.id),
    )
  },

  async get(id) {
    await delay()
    const db = getDb()
    const model = db.models.find((m) => m.id === id)
    if (!model) throw notFound('Product')
    return clone(withStock(db, model))
  },

  async create(payload) {
    await delay(200)
    const db = getDb()
    const data = validateModel(db, payload)
    const model = { id: nextId(db.models), ...data, created_at: nowIso(), updated_at: nowIso() }
    db.models.push(model)
    // inventory-service would create this on a ProductCreated event
    db.inventory.push({
      model_id: model.id,
      sku: `SKU-${String(model.id).padStart(6, '0')}`,
      selling_price: 0,
      unit_cost: 0,
      quantity_on_hand: 0,
      reorder_level: 5,
      location: null,
      updated_at: nowIso(),
    })
    commit()
    return clone(model)
  },

  async update(id, payload) {
    await delay(200)
    const db = getDb()
    const model = db.models.find((m) => m.id === id)
    if (!model) throw notFound('Product')
    if (Number(payload.component_type_id) !== model.component_type_id) {
      throw new ApiError(422, 'Component type cannot be changed', {
        component_type_id: 'Component type cannot be changed',
      })
    }
    Object.assign(model, validateModel(db, payload, id), { updated_at: nowIso() })
    commit()
    return clone(model)
  },

  async remove(id) {
    await delay()
    const db = getDb()
    const idx = db.models.findIndex((m) => m.id === id)
    if (idx < 0) throw notFound('Product')
    if (reservedQuantity(db, id) > 0) {
      throw new ApiError(409, 'Cannot delete: product is reserved by open orders. Deactivate it instead.')
    }
    db.models.splice(idx, 1)
    db.inventory = db.inventory.filter((i) => i.model_id !== id)
    commit()
  },
}
