import { defineStore } from 'pinia'
import { lookupApi } from '@/services/lookupApi'

// Cache of all lookup tables, used to render dropdowns and resolve ids to names
export const useLookupStore = defineStore('lookups', {
  state: () => ({ tables: {}, loaded: false }),
  getters: {
    componentTypes: (s) => s.tables.component_types ?? [],
    manufacturers: (s) => s.tables.manufacturers ?? [],
  },
  actions: {
    async load(force = false) {
      if (this.loaded && !force) return
      this.tables = await lookupApi.getAll()
      this.loaded = true
    },
    // Active rows plus the currently selected ones (so inactive values still render on edit)
    options(table, keepIds = []) {
      const keep = new Set([keepIds].flat().filter((v) => v != null).map(Number))
      return (this.tables[table] ?? [])
        .filter((r) => r.is_active || keep.has(r.id))
        .sort((a, b) => a.name.localeCompare(b.name))
    },
    find(table, id) {
      return (this.tables[table] ?? []).find((r) => r.id === Number(id)) ?? null
    },
    name(table, id) {
      return this.find(table, id)?.name ?? (id == null ? '' : `#${id}`)
    },
    typeByCode(code) {
      return this.componentTypes.find((t) => t.code === code) ?? null
    },
  },
})
