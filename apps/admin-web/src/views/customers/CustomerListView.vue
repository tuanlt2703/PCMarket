<template>
  <div class="card mb-4">
    <div class="card-header">
      <h3 class="card-title">Customers</h3>
      <div class="card-tools">
        <button type="button" class="btn btn-primary btn-sm" @click="editing = {}">
          <i class="bi bi-person-plus me-1"></i>New customer
        </button>
      </div>
    </div>
    <div class="card-body">
      <div class="row g-2 mb-3">
        <div class="col-md-6">
          <input v-model.trim="search" type="search" class="form-control" placeholder="Search name, email, phone, code…" />
        </div>
        <div class="col-md-3">
          <select v-model="tier" class="form-select">
            <option value="">All tiers</option>
            <option v-for="(m, key) in CUSTOMER_TIERS" :key="key" :value="key">{{ m.label }}</option>
          </select>
        </div>
      </div>
      <div v-if="loading" class="text-center py-4"><div class="spinner-border text-primary"></div></div>
      <div v-else class="table-responsive">
        <table class="table table-hover table-bordered align-middle">
          <thead>
            <tr>
              <th>Code</th>
              <th>Name</th>
              <th>Contact</th>
              <th>Tier</th>
              <th class="text-end">Points</th>
              <th class="text-center">Orders</th>
              <th class="text-end">Total spent</th>
              <th>Status</th>
              <th class="text-end" style="width: 130px">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="c in paged" :key="c.id">
              <td class="font-monospace small">{{ c.code }}</td>
              <td><RouterLink :to="{ name: 'customer-detail', params: { id: c.id } }">{{ c.full_name }}</RouterLink></td>
              <td class="small">{{ c.email }}<br />{{ c.phone }}</td>
              <td><StatusBadge :value="c.tier" :map="CUSTOMER_TIERS" /></td>
              <td class="text-end">{{ formatNumber(c.loyalty_points) }}</td>
              <td class="text-center">{{ c.order_count }}</td>
              <td class="text-end text-nowrap">{{ formatCurrency(c.total_spent) }}</td>
              <td>
                <span class="badge" :class="c.is_active ? 'text-bg-success' : 'text-bg-secondary'">{{ c.is_active ? 'Active' : 'Inactive' }}</span>
              </td>
              <td class="text-end text-nowrap">
                <RouterLink :to="{ name: 'customer-detail', params: { id: c.id } }" class="btn btn-sm btn-outline-info me-1" title="View">
                  <i class="bi bi-eye"></i>
                </RouterLink>
                <button type="button" class="btn btn-sm btn-outline-primary me-1" title="Edit" @click="editing = c">
                  <i class="bi bi-pencil"></i>
                </button>
                <button type="button" class="btn btn-sm btn-outline-danger" title="Delete" @click="deleting = c">
                  <i class="bi bi-trash"></i>
                </button>
              </td>
            </tr>
            <tr v-if="!filtered.length"><td colspan="9" class="text-center text-secondary">No customers</td></tr>
          </tbody>
        </table>
      </div>
    </div>
    <div class="card-footer">
      <PaginationBar v-model="page" :total="filtered.length" :page-size="pageSize" />
    </div>
  </div>

  <CustomerFormModal v-if="editing" :customer="editing" @close="editing = null" @saved="onSaved" />
  <ConfirmDialog
    v-if="deleting"
    title="Delete customer"
    :message="`Delete customer '${deleting.full_name}'?`"
    :busy="busy"
    @confirm="remove"
    @cancel="deleting = null"
  />
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { customerApi } from '@/services/customerApi'
import { CUSTOMER_TIERS } from '@/config/statuses'
import { formatCurrency, formatNumber } from '@/utils/format'
import { toast } from '@/stores/toast'
import { usePagination } from '@/composables/usePagination'
import StatusBadge from '@/components/common/StatusBadge.vue'
import PaginationBar from '@/components/common/PaginationBar.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import CustomerFormModal from '@/components/customer/CustomerFormModal.vue'

const rows = ref([])
const loading = ref(true)
const search = ref('')
const tier = ref('')
const editing = ref(null)
const deleting = ref(null)
const busy = ref(false)

const filtered = computed(() => {
  const term = search.value.toLowerCase()
  return rows.value
    .filter((c) => !tier.value || c.tier === tier.value)
    .filter((c) => !term || [c.full_name, c.email, c.phone, c.code].some((v) => v?.toLowerCase().includes(term)))
})
const { page, paged, pageSize } = usePagination(filtered, { pageSize: 15, resetOn: [search, tier] })

async function load() {
  loading.value = true
  rows.value = await customerApi.list()
  loading.value = false
}

function onSaved() {
  editing.value = null
  load()
}

async function remove() {
  busy.value = true
  try {
    await customerApi.remove(deleting.value.id)
    toast.success('Customer deleted')
    await load()
  } catch (e) {
    toast.error(e.message)
  } finally {
    busy.value = false
    deleting.value = null
  }
}

onMounted(load)
</script>
