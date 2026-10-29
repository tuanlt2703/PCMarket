<template>
  <div class="card mb-4">
    <div class="card-header">
      <h3 class="card-title">Orders</h3>
    </div>
    <div class="card-body">
      <ul class="nav nav-pills mb-3 flex-wrap">
        <li class="nav-item">
          <a href="#" class="nav-link" :class="{ active: !status }" @click.prevent="status = ''">
            All <span class="badge text-bg-light ms-1">{{ allOrders.length }}</span>
          </a>
        </li>
        <li v-for="(m, key) in ORDER_STATUSES" :key="key" class="nav-item">
          <a href="#" class="nav-link" :class="{ active: status === key }" @click.prevent="status = key">
            {{ m.label }} <span class="badge text-bg-light ms-1">{{ countBy(key) }}</span>
          </a>
        </li>
      </ul>
      <div class="row g-2 mb-3">
        <div class="col-md-4">
          <input v-model.trim="search" type="search" class="form-control" placeholder="Search order # or customer…" />
        </div>
        <div class="col-md-2">
          <select v-model="paymentStatus" class="form-select">
            <option value="">All payments</option>
            <option v-for="(m, key) in PAYMENT_STATUSES" :key="key" :value="key">{{ m.label }}</option>
          </select>
        </div>
        <div class="col-md-3">
          <div class="input-group">
            <span class="input-group-text">From</span>
            <input v-model="from" type="date" class="form-control" />
          </div>
        </div>
        <div class="col-md-3">
          <div class="input-group">
            <span class="input-group-text">To</span>
            <input v-model="to" type="date" class="form-control" />
          </div>
        </div>
      </div>

      <div v-if="loading" class="text-center py-4"><div class="spinner-border text-primary"></div></div>
      <div v-else class="table-responsive">
        <table class="table table-hover table-bordered align-middle">
          <thead>
            <tr>
              <th>Order #</th>
              <th>Customer</th>
              <th class="text-center">Items</th>
              <th class="text-end">Total</th>
              <th>Payment</th>
              <th>Status</th>
              <th>Created</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="o in paged" :key="o.id">
              <td><RouterLink :to="{ name: 'order-detail', params: { id: o.id } }" class="font-monospace">{{ o.order_number }}</RouterLink></td>
              <td>
                <RouterLink :to="{ name: 'customer-detail', params: { id: o.customer_id } }">{{ o.customer_name }}</RouterLink>
              </td>
              <td class="text-center">{{ o.items.reduce((s, i) => s + i.quantity, 0) }}</td>
              <td class="text-end text-nowrap">{{ formatCurrency(o.total) }}</td>
              <td>
                <StatusBadge :value="o.payment_status" :map="PAYMENT_STATUSES" />
                <div class="small text-secondary">{{ PAYMENT_METHODS[o.payment_method] }}</div>
              </td>
              <td><StatusBadge :value="o.status" :map="ORDER_STATUSES" /></td>
              <td class="text-nowrap">{{ formatDateTime(o.created_at) }}</td>
              <td class="text-end">
                <RouterLink :to="{ name: 'order-detail', params: { id: o.id } }" class="btn btn-sm btn-outline-info" title="View">
                  <i class="bi bi-eye"></i>
                </RouterLink>
              </td>
            </tr>
            <tr v-if="!filtered.length"><td colspan="8" class="text-center text-secondary">No orders</td></tr>
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
import { orderApi } from '@/services/orderApi'
import { ORDER_STATUSES, PAYMENT_STATUSES, PAYMENT_METHODS } from '@/config/statuses'
import { formatCurrency, formatDateTime } from '@/utils/format'
import { usePagination } from '@/composables/usePagination'
import StatusBadge from '@/components/common/StatusBadge.vue'
import PaginationBar from '@/components/common/PaginationBar.vue'

const allOrders = ref([])
const loading = ref(true)
const status = ref('')
const paymentStatus = ref('')
const search = ref('')
const from = ref('')
const to = ref('')

const baseFiltered = computed(() => {
  const term = search.value.toLowerCase()
  return allOrders.value
    .filter((o) => !paymentStatus.value || o.payment_status === paymentStatus.value)
    .filter((o) => !from.value || o.created_at.slice(0, 10) >= from.value)
    .filter((o) => !to.value || o.created_at.slice(0, 10) <= to.value)
    .filter((o) => !term || o.order_number.toLowerCase().includes(term) || o.customer_name.toLowerCase().includes(term))
})
const filtered = computed(() => baseFiltered.value.filter((o) => !status.value || o.status === status.value))
const countBy = (s) => baseFiltered.value.filter((o) => o.status === s).length

const { page, paged, pageSize } = usePagination(filtered, { pageSize: 15, resetOn: [status, paymentStatus, search, from, to] })

onMounted(async () => {
  allOrders.value = await orderApi.list()
  loading.value = false
})
</script>
