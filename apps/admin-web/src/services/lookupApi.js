import { getDb, commit, nextId, clone, nowIso, delay, ApiError, notFound } from './mockDb'
import { getLookupMeta } from '@/config/lookupTables'
import { getSchema, collectLookupRefs } from '@/config/componentSchemas'

function tableRows(table) {
  if (!getLookupMeta(table)) throw new ApiError(404, `Unknown lookup table: ${table}`)
  const db = getDb()
  db.lookups[table] ??= []
  return db.lookups[table]
}

function validate(table, payload, id = null) {
  const errors = {}
  const code = String(payload.code ?? '').trim().toUpperCase()
  const name = String(payload.name ?? '').trim()
  if (!code) errors.code = 'Code is required'
  else if (code.length > 50) errors.code = 'Max 50 characters'
  else if (!/^[A-Z0-9_.\-]+$/.test(code)) errors.code = 'Only A-Z, 0-9, _ . - allowed'
  else if (tableRows(table).some((r) => r.code === code && r.id !== id)) errors.code = 'Code already exists'
  if (!name) errors.name = 'Name is required'
  else if (name.length > 100) errors.name = 'Max 100 characters'
  if (Object.keys(errors).length) throw new ApiError(422, 'Validation failed', errors)
  return {
    code,
    name,
    description: payload.description?.trim() || null,
    is_active: payload.is_active !== false,
  }
}

// Mirrors the FK constraints in product-db: a lookup row cannot be removed while referenced
export function findLookupUsages(table, id) {
  const db = getDb()
  const typeCode = (typeId) => db.lookups.component_types.find((t) => t.id === typeId)?.code
  const usedBy = []
  for (const m of db.models) {
    const refs = collectLookupRefs(getSchema(typeCode(m.component_type_id)), m.spec)
    refs.push(['component_types', m.component_type_id], ['manufacturers', m.manufacturer_id])
    if (refs.some(([t, v]) => t === table && v === id)) usedBy.push({ id: m.id, model_name: m.model_name })
  }
  return usedBy
}

export const lookupApi = {
  async getAll() {
    await delay(60)
    return clone(getDb().lookups)
  },

  async list(table) {
    await delay()
    return clone(tableRows(table))
  },

  async create(table, payload) {
    await delay()
    const rows = tableRows(table)
    const data = validate(table, payload)
    const row = { id: nextId(rows), ...data, created_at: nowIso(), updated_at: nowIso() }
    rows.push(row)
    commit()
    return clone(row)
  },

  async update(table, id, payload) {
    await delay()
    const row = tableRows(table).find((r) => r.id === id)
    if (!row) throw notFound('Lookup item')
    Object.assign(row, validate(table, payload, id), { updated_at: nowIso() })
    commit()
    return clone(row)
  },

  async remove(table, id) {
    await delay()
    const rows = tableRows(table)
    const idx = rows.findIndex((r) => r.id === id)
    if (idx < 0) throw notFound('Lookup item')
    const usages = findLookupUsages(table, id)
    if (usages.length) {
      throw new ApiError(
        409,
        `Cannot delete: used by ${usages.length} product(s), e.g. "${usages[0].model_name}". Deactivate it instead.`,
        usages,
      )
    }
    rows.splice(idx, 1)
    commit()
  },

  async usageCount(table, id) {
    await delay(30)
    return findLookupUsages(table, id).length
  },
}
