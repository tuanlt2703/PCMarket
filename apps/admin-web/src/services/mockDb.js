import { createSeed } from '@/mock/seed'

// In-browser mock database persisted to localStorage until the real backend services exist.
const STORAGE_KEY = 'pcmarket-admin-mock-db'
const DB_VERSION = 1

let db = load()

function load() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) {
      const parsed = JSON.parse(raw)
      if (parsed?.version === DB_VERSION) return parsed
    }
  } catch {
    // corrupted storage -> fall through to a fresh seed
  }
  const seeded = createSeed()
  save(seeded)
  return seeded
}

function save(data) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(data))
  } catch {
    // storage full / unavailable: keep working in memory
  }
}

export class ApiError extends Error {
  constructor(status, message, details) {
    super(message)
    this.status = status
    this.details = details
  }
}

export const getDb = () => db
export const commit = () => save(db)

export function resetDb() {
  db = createSeed()
  save(db)
}

export const nextId = (rows) => rows.reduce((max, r) => Math.max(max, r.id), 0) + 1
export const clone = (v) => structuredClone(v)
export const nowIso = () => new Date().toISOString()

// Simulates network latency so loading states are exercised
export const delay = (ms = 120) => new Promise((resolve) => setTimeout(resolve, ms))

export function notFound(entity) {
  return new ApiError(404, `${entity} not found`)
}
