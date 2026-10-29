<template>
  <div class="row mb-3">
    <div class="col-md-6">
      <div class="input-group">
        <span class="input-group-text"><i class="bi bi-search"></i></span>
        <input v-model.trim="search" type="search" class="form-control" placeholder="Search lookup tables…" />
      </div>
    </div>
    <div class="col-md-6 mt-2 mt-md-0">
      <select v-model="group" class="form-select">
        <option value="">All groups</option>
        <option v-for="g in LOOKUP_GROUPS" :key="g.key" :value="g.key">{{ g.label }}</option>
      </select>
    </div>
  </div>

  <div class="row">
    <div v-for="g in visibleGroups" :key="g.key" class="col-xl-4 col-md-6">
      <div class="card card-outline card-primary mb-4">
        <div class="card-header">
          <h3 class="card-title"><i :class="['bi', g.icon, 'me-2']"></i>{{ g.label }}</h3>
        </div>
        <div class="card-body p-0">
          <ul class="list-group list-group-flush">
            <li v-for="t in g.tables" :key="t.table" class="list-group-item d-flex justify-content-between align-items-center">
              <RouterLink :to="{ name: 'lookup-table', params: { table: t.table } }">
                {{ t.label }}
                <div class="small text-secondary font-monospace">{{ t.table }}</div>
              </RouterLink>
              <span class="badge text-bg-secondary rounded-pill" :title="`${activeCount(t.table)} active`">
                {{ rowCount(t.table) }}
              </span>
            </li>
          </ul>
        </div>
      </div>
    </div>
    <div v-if="!visibleGroups.length" class="col-12 text-secondary">No lookup tables match your search.</div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { LOOKUP_GROUPS } from '@/config/lookupTables'
import { useLookupStore } from '@/stores/lookups'

const route = useRoute()
const router = useRouter()
const lookups = useLookupStore()

const search = ref('')
const group = ref(route.query.group ?? '')
watch(group, (g) => router.replace({ query: g ? { group: g } : {} }))
watch(
  () => route.query.group,
  (g) => (group.value = g ?? ''),
)

const visibleGroups = computed(() => {
  const term = search.value.toLowerCase()
  return LOOKUP_GROUPS.filter((g) => !group.value || g.key === group.value)
    .map((g) => ({
      ...g,
      tables: g.tables.filter((t) => !term || t.label.toLowerCase().includes(term) || t.table.includes(term)),
    }))
    .filter((g) => g.tables.length)
})

const rowCount = (table) => lookups.tables[table]?.length ?? 0
const activeCount = (table) => lookups.tables[table]?.filter((r) => r.is_active).length ?? 0
</script>
