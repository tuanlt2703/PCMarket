<template>
  <div v-if="!order" class="text-center py-5"><div class="spinner-border text-primary"></div></div>
  <template v-else>
    <div class="row">
      <div class="col-lg-8">
        <div class="invoice p-3 mb-4 card">
          <div class="d-flex justify-content-between align-items-start mb-3">
            <div>
              <h4 class="mb-1"><i class="bi bi-receipt me-2"></i>{{ order.order_number }}</h4>
              <small class="text-secondary">Placed {{ formatDateTime(order.created_at) }}</small>
            </div>
            <div class="text-end">
              <StatusBadge :value="order.status" :map="ORDER_STATUSES" />
              <div class="mt-1"><StatusBadge :value="order.payment_status" :map="PAYMENT_STATUSES" /></div>
            </div>
          </div>

          <div class="row mb-3">
            <div class="col-sm-6">
              <strong>Customer</strong>
              <address class="mb-0">
                <RouterLink :to="{ name: 'customer-detail', params: { id: order.customer_id } }">{{ order.customer_name }}</RouterLink><br />
                {{ order.customer_email }}<br />
                {{ order.customer_phone }}
              </address>
            </div>
            <div class="col-sm-6">
              <strong>Shipping address</strong>
              <address class="mb-0">{{ order.shipping_address }}</address>
              <div class="mt-1"><strong>Payment:</strong> {{ PAYMENT_METHODS[order.payment_method] }}</div>
            </div>
          </div>

          <div class="table-responsive">
            <table class="table table-striped">
              <thead>
                <tr><th>Product</th><th>SKU</th><th class="text-center">Qty</th><th class="text-end">Unit price</th><th class="text-end">Subtotal</th></tr>
              </thead>
              <tbody>
                <tr v-for="it in order.items" :key="it.model_id">
                  <td><RouterLink :to="{ name: 'product-detail', params: { id: it.model_id } }">{{ it.product_name }}</RouterLink></td>
                  <td class="font-monospace small">{{ it.sku }}</td>
                  <td class="text-center">{{ it.quantity }}</td>
                  <td class="text-end">{{ formatCurrency(it.unit_price) }}</td>
                  <td class="text-end">{{ formatCurrency(it.unit_price * it.quantity) }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <div class="row justify-content-end">
            <div class="col-md-6">
              <table class="table mb-0">
                <tbody>
                  <tr><th>Subtotal</th><td class="text-end">{{ formatCurrency(order.subtotal) }}</td></tr>
                  <tr><th>Discount</th><td class="text-end">-{{ formatCurrency(order.discount) }}</td></tr>
                  <tr><th>Shipping</th><td class="text-end">{{ formatCurrency(order.shipping_fee) }}</td></tr>
                  <tr><th>Total</th><td class="text-end fw-bold fs-5">{{ formatCurrency(order.total) }}</td></tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <div class="col-lg-4">
        <div class="card card-primary card-outline mb-4">
          <div class="card-header"><h3 class="card-title">Update status</h3></div>
          <div class="card-body">
            <template v-if="nextStatuses.length">
              <div class="mb-2">
                <textarea v-model="note" class="form-control" rows="2" placeholder="Note (optional)"></textarea>
              </div>
              <div class="d-flex flex-wrap gap-2">
                <button
                  v-for="s in nextStatuses"
                  :key="s"
                  type="button"
                  class="btn btn-sm"
                  :class="`btn-${ORDER_STATUSES[s].color}`"
                  :disabled="busy"
                  @click="changeStatus(s)"
                >
                  Mark as {{ ORDER_STATUSES[s].label }}
                </button>
              </div>
              <div v-if="nextStatuses.includes('SHIPPED')" class="form-text">
                Shipping deducts the items from stock on hand.
              </div>
            </template>
            <p v-else class="text-secondary mb-0">This order is closed.</p>
            <hr />
            <label class="form-label" for="odPay">Payment status</label>
            <select id="odPay" class="form-select" :value="order.payment_status" :disabled="busy" @change="changePayment($event.target.value)">
              <option v-for="(m, key) in PAYMENT_STATUSES" :key="key" :value="key">{{ m.label }}</option>
            </select>
          </div>
        </div>

        <div class="card mb-4">
          <div class="card-header"><h3 class="card-title">History</h3></div>
          <div class="card-body">
            <div class="timeline">
              <div v-for="(h, i) in [...order.status_history].reverse()" :key="i">
                <i class="timeline-icon bi bi-circle-fill" :class="`text-bg-${ORDER_STATUSES[h.status]?.color ?? 'secondary'}`"></i>
                <div class="timeline-item">
                  <span class="time"><i class="bi bi-clock-fill"></i> {{ formatDateTime(h.at) }}</span>
                  <h3 class="timeline-header">{{ ORDER_STATUSES[h.status]?.label ?? h.status }}</h3>
                  <div v-if="h.note" class="timeline-body">{{ h.note }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </template>
</template>

<script setup>
import { computed, inject, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { orderApi } from '@/services/orderApi'
import { ORDER_STATUSES, ORDER_TRANSITIONS, PAYMENT_STATUSES, PAYMENT_METHODS } from '@/config/statuses'
import { formatCurrency, formatDateTime } from '@/utils/format'
import { toast } from '@/stores/toast'
import StatusBadge from '@/components/common/StatusBadge.vue'

const route = useRoute()
const router = useRouter()
const setPageTitle = inject('setPageTitle')

const order = ref(null)
const note = ref('')
const busy = ref(false)

const nextStatuses = computed(() => ORDER_TRANSITIONS[order.value?.status] ?? [])

async function run(action, msg) {
  busy.value = true
  try {
    order.value = await action()
    toast.success(msg)
  } catch (e) {
    toast.error(e.message)
  } finally {
    busy.value = false
  }
}

async function changeStatus(s) {
  await run(() => orderApi.updateStatus(order.value.id, s, note.value), `Order marked as ${ORDER_STATUSES[s].label}`)
  note.value = ''
}

const changePayment = (s) => run(() => orderApi.updatePaymentStatus(order.value.id, s), 'Payment status updated')

onMounted(async () => {
  try {
    order.value = await orderApi.get(Number(route.params.id))
    setPageTitle(`Order ${order.value.order_number}`)
  } catch (e) {
    toast.error(e.message)
    router.replace({ name: 'orders' })
  }
})
</script>
