<template>
  <div v-if="!customer" class="text-center py-5"><div class="spinner-border text-primary"></div></div>
  <div v-else class="row">
    <div class="col-lg-4">
      <div class="card card-primary card-outline mb-4">
        <div class="card-body box-profile">
          <div class="text-center mb-2">
            <i class="bi bi-person-circle display-4 text-secondary"></i>
          </div>
          <h3 class="profile-username text-center">{{ customer.full_name }}</h3>
          <p class="text-secondary text-center mb-2">{{ customer.code }}</p>
          <p class="text-center"><StatusBadge :value="customer.tier" :map="CUSTOMER_TIERS" /></p>
          <ul class="list-group list-group-flush mb-3">
            <li class="list-group-item"><b>Email</b> <span class="float-end">{{ customer.email }}</span></li>
            <li class="list-group-item"><b>Phone</b> <span class="float-end">{{ customer.phone ?? '—' }}</span></li>
            <li class="list-group-item"><b>Loyalty points</b> <span class="float-end">{{ formatNumber(customer.loyalty_points) }}</span></li>
            <li class="list-group-item"><b>Orders</b> <span class="float-end">{{ customer.order_count }}</span></li>
            <li class="list-group-item"><b>Total spent</b> <span class="float-end">{{ formatCurrency(customer.total_spent) }}</span></li>
            <li class="list-group-item"><b>Joined</b> <span class="float-end">{{ formatDate(customer.created_at) }}</span></li>
            <li class="list-group-item">
              <b>Status</b>
              <span class="float-end badge" :class="customer.is_active ? 'text-bg-success' : 'text-bg-secondary'">
                {{ customer.is_active ? 'Active' : 'Inactive' }}
              </span>
            </li>
          </ul>
          <p class="small"><b>Address:</b> {{ customer.address ?? '—' }}</p>
          <button type="button" class="btn btn-primary w-100" @click="editing = true">
            <i class="bi bi-pencil me-1"></i>Edit
          </button>
        </div>
      </div>
    </div>

    <div class="col-lg-8">
      <div class="card mb-4">
        <div class="card-header"><h3 class="card-title">Orders</h3></div>
        <div class="card-body p-0">
          <table class="table table-striped mb-0">
            <thead><tr><th>Order #</th><th>Date</th><th>Status</th><th class="text-end">Total</th></tr></thead>
            <tbody>
              <tr v-for="o in customer.orders" :key="o.id">
                <td><RouterLink :to="{ name: 'order-detail', params: { id: o.id } }" class="font-monospace">{{ o.order_number }}</RouterLink></td>
                <td>{{ formatDateTime(o.created_at) }}</td>
                <td><StatusBadge :value="o.status" :map="ORDER_STATUSES" /></td>
                <td class="text-end">{{ formatCurrency(o.total) }}</td>
              </tr>
              <tr v-if="!customer.orders.length"><td colspan="4" class="text-center text-secondary">No orders</td></tr>
            </tbody>
          </table>
        </div>
      </div>
      <div class="card mb-4">
        <div class="card-header"><h3 class="card-title">Warranty claims</h3></div>
        <div class="card-body p-0">
          <table class="table table-striped mb-0">
            <thead><tr><th>Claim #</th><th>Serial</th><th>Status</th><th>Created</th></tr></thead>
            <tbody>
              <tr v-for="w in customer.warranties" :key="w.id">
                <td class="font-monospace">{{ w.claim_number }}</td>
                <td class="font-monospace small">{{ w.serial_number }}</td>
                <td><StatusBadge :value="w.status" :map="WARRANTY_STATUSES" /></td>
                <td>{{ formatDate(w.created_at) }}</td>
              </tr>
              <tr v-if="!customer.warranties.length"><td colspan="4" class="text-center text-secondary">No claims</td></tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>

  <CustomerFormModal v-if="editing" :customer="customer" @close="editing = false" @saved="onSaved" />
</template>

<script setup>
import { inject, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { customerApi } from '@/services/customerApi'
import { CUSTOMER_TIERS, ORDER_STATUSES, WARRANTY_STATUSES } from '@/config/statuses'
import { formatCurrency, formatDate, formatDateTime, formatNumber } from '@/utils/format'
import { toast } from '@/stores/toast'
import StatusBadge from '@/components/common/StatusBadge.vue'
import CustomerFormModal from '@/components/customer/CustomerFormModal.vue'

const route = useRoute()
const router = useRouter()
const setPageTitle = inject('setPageTitle')

const customer = ref(null)
const editing = ref(false)

async function load() {
  try {
    customer.value = await customerApi.get(Number(route.params.id))
    setPageTitle(customer.value.full_name)
  } catch (e) {
    toast.error(e.message)
    router.replace({ name: 'customers' })
  }
}

function onSaved() {
  editing.value = false
  load()
}

onMounted(load)
</script>
