<template>
  <div v-if="!data" class="text-center py-5"><div class="spinner-border text-primary"></div></div>
  <template v-else>
    <div class="row">
      <div v-for="box in boxes" :key="box.label" class="col-lg-3 col-6">
        <div class="small-box" :class="`text-bg-${box.color}`">
          <div class="inner">
            <h3>{{ box.value }}</h3>
            <p>{{ box.label }}</p>
          </div>
          <i :class="['small-box-icon', 'bi', box.icon]" aria-hidden="true"></i>
          <RouterLink
            :to="box.to"
            class="small-box-footer link-underline-opacity-0 link-underline-opacity-50-hover"
            :class="box.color === 'warning' ? 'link-dark' : 'link-light'"
          >
            More info <i class="bi bi-link-45deg"></i>
          </RouterLink>
        </div>
      </div>
    </div>

    <div class="row">
      <div class="col-lg-8">
        <div class="card mb-4">
          <div class="card-header"><h3 class="card-title">Sales — last 30 days</h3></div>
          <div class="card-body">
            <apexchart type="area" height="300" :options="salesOptions" :series="salesSeries" />
          </div>
        </div>

        <div class="card mb-4">
          <div class="card-header">
            <h3 class="card-title">Latest Orders</h3>
            <div class="card-tools">
              <RouterLink :to="{ name: 'orders' }" class="btn btn-sm btn-secondary">View all</RouterLink>
            </div>
          </div>
          <div class="card-body p-0">
            <div class="table-responsive">
              <table class="table table-striped m-0">
                <thead>
                  <tr><th>Order</th><th>Customer</th><th>Status</th><th class="text-end">Total</th><th>Date</th></tr>
                </thead>
                <tbody>
                  <tr v-for="o in data.recentOrders" :key="o.id">
                    <td><RouterLink :to="{ name: 'order-detail', params: { id: o.id } }">{{ o.order_number }}</RouterLink></td>
                    <td>{{ o.customer_name }}</td>
                    <td><StatusBadge :value="o.status" :map="ORDER_STATUSES" /></td>
                    <td class="text-end">{{ formatCurrency(o.total) }}</td>
                    <td>{{ formatDateTime(o.created_at) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <div class="col-lg-4">
        <div class="card mb-4">
          <div class="card-header"><h3 class="card-title">Revenue by component type</h3></div>
          <div class="card-body">
            <apexchart type="donut" height="300" :options="typeOptions" :series="typeSeries" />
          </div>
        </div>

        <div class="card mb-4">
          <div class="card-header">
            <h3 class="card-title">Low stock</h3>
            <div class="card-tools">
              <span class="badge text-bg-danger">{{ data.lowStockCount }}</span>
            </div>
          </div>
          <div class="card-body p-0">
            <ul class="list-group list-group-flush">
              <li v-for="s in data.lowStock" :key="s.model_id" class="list-group-item d-flex justify-content-between align-items-center">
                <RouterLink :to="{ name: 'product-detail', params: { id: s.model_id } }" class="text-truncate me-2 min-w-0">
                  {{ s.model_name }}
                </RouterLink>
                <span class="badge" :class="s.available <= 0 ? 'text-bg-danger' : 'text-bg-warning'">{{ s.available }} left</span>
              </li>
              <li v-if="!data.lowStock.length" class="list-group-item text-secondary">All products are sufficiently stocked.</li>
            </ul>
          </div>
        </div>

        <div class="card mb-4">
          <div class="card-header"><h3 class="card-title">Warranty claims</h3></div>
          <div class="card-body">
            <div v-for="(meta, key) in WARRANTY_STATUSES" :key="key" class="d-flex justify-content-between mb-1">
              <span><StatusBadge :value="key" :map="WARRANTY_STATUSES" /></span>
              <strong>{{ data.warrantyByStatus[key] ?? 0 }}</strong>
            </div>
          </div>
        </div>
      </div>
    </div>
  </template>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { dashboardApi } from '@/services/dashboardApi'
import { ORDER_STATUSES, WARRANTY_STATUSES } from '@/config/statuses'
import { formatCurrency, formatDateTime } from '@/utils/format'
import StatusBadge from '@/components/common/StatusBadge.vue'

const data = ref(null)

const compactVnd = (v) =>
  v >= 1e9 ? `${(v / 1e9).toFixed(1)}B` : v >= 1e6 ? `${(v / 1e6).toFixed(1)}M` : new Intl.NumberFormat('vi-VN').format(v)

const boxes = computed(() => [
  { label: `Orders (30 days) · ${data.value.pendingOrders} pending`, value: data.value.newOrders, color: 'primary', icon: 'bi-cart-fill', to: { name: 'orders' } },
  { label: 'Revenue (30 days, VND)', value: compactVnd(data.value.revenue), color: 'success', icon: 'bi-cash-stack', to: { name: 'orders' } },
  { label: `Customers · ${data.value.newCustomers} new`, value: data.value.customers, color: 'warning', icon: 'bi-person-plus-fill', to: { name: 'customers' } },
  { label: 'Low-stock products', value: data.value.lowStockCount, color: 'danger', icon: 'bi-exclamation-triangle-fill', to: { name: 'inventory', query: { filter: 'low' } } },
])

const salesSeries = computed(() => [
  { name: 'Revenue', data: data.value.series.map((d) => d.revenue) },
  { name: 'Orders', data: data.value.series.map((d) => d.orders) },
])

const salesOptions = computed(() => ({
  chart: { toolbar: { show: false }, background: 'transparent' },
  theme: { mode: document.documentElement.getAttribute('data-bs-theme') === 'dark' ? 'dark' : 'light' },
  colors: ['#0d6efd', '#20c997'],
  dataLabels: { enabled: false },
  stroke: { curve: 'smooth', width: 2 },
  xaxis: { categories: data.value.series.map((d) => d.date.slice(5)), tickAmount: 10 },
  yaxis: [
    { title: { text: 'Revenue (VND)' }, labels: { formatter: compactVnd } },
    { opposite: true, title: { text: 'Orders' }, labels: { formatter: (v) => Math.round(v) } },
  ],
  tooltip: { y: { formatter: (v, { seriesIndex }) => (seriesIndex === 0 ? formatCurrency(v) : v) } },
  legend: { position: 'top' },
}))

const typeSeries = computed(() => Object.values(data.value.salesByType))
const typeOptions = computed(() => ({
  labels: Object.keys(data.value.salesByType),
  colors: ['#0d6efd', '#198754', '#ffc107', '#dc3545', '#6f42c1', '#0dcaf0', '#fd7e14', '#20c997', '#6c757d'],
  chart: { background: 'transparent' },
  theme: { mode: document.documentElement.getAttribute('data-bs-theme') === 'dark' ? 'dark' : 'light' },
  legend: { position: 'bottom' },
  tooltip: { y: { formatter: formatCurrency } },
  dataLabels: { formatter: (v) => `${v.toFixed(0)}%` },
}))

onMounted(async () => {
  data.value = await dashboardApi.summary({ days: 30 })
})
</script>
