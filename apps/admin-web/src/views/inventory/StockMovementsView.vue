<template>
  <div class="card mb-4">
    <div class="card-header">
      <h3 class="card-title">Stock movements</h3>
    </div>
    <div class="card-body">
      <div class="row g-2 mb-3">
        <div class="col-md-6">
          <input v-model.trim="search" type="search" class="form-control" placeholder="Search product, reason or reference…" />
        </div>
        <div class="col-md-3">
          <select v-model="type" class="form-select">
            <option value="">All types</option>
            <option v-for="(m, key) in STOCK_MOVEMENT_TYPES" :key="key" :value="key">{{ m.label }}</option>
          </select>
        </div>
      </div>
      <div v-if="loading" class="text-center py-4"><div class="spinner-border text-primary"></div></div>
      <div v-else class="table-responsive">
        <table class="table table-striped table-bordered align-middle">
          <thead>
            <tr><th>Date</th><th>Product</th><th>Type</th><th class="text-center">Quantity</th><th>Reason</th><th>Reference</th></tr>
          </thead>
          <tbody>
            <tr v-for="m in paged" :key="m.id">
              <td class="text-nowrap">{{ formatDateTime(m.created_at) }}</td>
              <td><RouterLink :to="{ name: 'product-detail', params: { id: m.model_id } }">{{ m.model_name }}</RouterLink></td>
              <td><StatusBadge :value="m.type" :map="STOCK_MOVEMENT_TYPES" /></td>
              <td class="text-center fw-bold" :class="m.quantity >= 0 ? 'text-success' : 'text-danger'">
                {{ m.quantity > 0 ? '+' : '' }}{{ m.quantity }}
              </td>
              <td>{{ m.reason }}</td>
              <td class="font-monospace small">{{ m.reference ?? '—' }}</td>
            </tr>
            <tr v-if="!filtered.length"><td colspan="6" class="text-center text-secondary">No data</td></tr>
          </tbody>
        </table>
      </div>
    </div>
    <div class="card-footer">
      <PaginationBar v-model="page" :total="filtered.length" :page-size="pageSize" />
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { inventoryApi } from '@/services/inventoryApi'
import { STOCK_MOVEMENT_TYPES } from '@/config/statuses'
import { formatDateTime } from '@/utils/format'
import { usePagination } from '@/composables/usePagination'
import StatusBadge from '@/components/common/StatusBadge.vue'
import PaginationBar from '@/components/common/PaginationBar.vue'

const rows = ref([])
const loading = ref(true)
const search = ref('')
const type = ref('')

const filtered = computed(() => {
  const term = search.value.toLowerCase()
  return rows.value
    .filter((m) => !type.value || m.type === type.value)
    .filter((m) => !term || [m.model_name, m.reason, m.reference].some((v) => v?.toLowerCase().includes(term)))
})
const { page, paged, pageSize } = usePagination(filtered, { pageSize: 20, resetOn: [search, type] })

onMounted(async () => {
  rows.value = await inventoryApi.movements()
  loading.value = false
})
</script>
