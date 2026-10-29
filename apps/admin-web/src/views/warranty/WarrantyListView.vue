<template>
  <div class="card mb-4">
    <div class="card-header">
      <h3 class="card-title">Warranty claims</h3>
      <div class="card-tools">
        <button type="button" class="btn btn-primary btn-sm" @click="openCreate">
          <i class="bi bi-plus-lg me-1"></i>New claim
        </button>
      </div>
    </div>
    <div class="card-body">
      <div class="row g-2 mb-3">
        <div class="col-md-6">
          <input v-model.trim="search" type="search" class="form-control" placeholder="Search claim #, customer, product, serial…" />
        </div>
        <div class="col-md-3">
          <select v-model="status" class="form-select">
            <option value="">All statuses</option>
            <option v-for="(m, key) in WARRANTY_STATUSES" :key="key" :value="key">{{ m.label }}</option>
          </select>
        </div>
      </div>
      <div v-if="loading" class="text-center py-4"><div class="spinner-border text-primary"></div></div>
      <div v-else class="table-responsive">
        <table class="table table-hover table-bordered align-middle">
          <thead>
            <tr>
              <th>Claim #</th>
              <th>Customer</th>
              <th>Product / Serial</th>
              <th>Order</th>
              <th>Issue</th>
              <th>Status</th>
              <th>Expiry</th>
              <th>Created</th>
              <th class="text-end" style="width: 100px">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="w in paged" :key="w.id">
              <td class="font-monospace">{{ w.claim_number }}</td>
              <td><RouterLink :to="{ name: 'customer-detail', params: { id: w.customer_id } }">{{ w.customer_name }}</RouterLink></td>
              <td>
                <RouterLink :to="{ name: 'product-detail', params: { id: w.model_id } }">{{ w.model_name }}</RouterLink>
                <div class="small text-secondary font-monospace">{{ w.serial_number }}</div>
              </td>
              <td><RouterLink :to="{ name: 'order-detail', params: { id: w.order_id } }" class="font-monospace small">{{ w.order_number }}</RouterLink></td>
              <td class="small">{{ w.issue_description }}</td>
              <td><StatusBadge :value="w.status" :map="WARRANTY_STATUSES" /></td>
              <td class="text-nowrap">{{ formatDate(w.warranty_expiry) }}</td>
              <td class="text-nowrap">{{ formatDate(w.created_at) }}</td>
              <td class="text-end text-nowrap">
                <button type="button" class="btn btn-sm btn-outline-primary me-1" title="Edit" @click="openEdit(w)">
                  <i class="bi bi-pencil"></i>
                </button>
                <button type="button" class="btn btn-sm btn-outline-danger" title="Delete" @click="deleting = w">
                  <i class="bi bi-trash"></i>
                </button>
              </td>
            </tr>
            <tr v-if="!filtered.length"><td colspan="9" class="text-center text-secondary">No claims</td></tr>
          </tbody>
        </table>
      </div>
    </div>
    <div class="card-footer">
      <PaginationBar v-model="page" :total="filtered.length" :page-size="pageSize" />
    </div>
  </div>

  <BaseModal v-if="form" :title="form.id ? `Edit claim ${form.claim_number}` : 'New warranty claim'" size="lg" @close="form = null">
    <form id="warrantyForm" @submit.prevent="save">
      <div class="row">
        <div class="col-md-6 mb-3">
          <label class="form-label" for="wfOrder">Order (delivered) <span class="text-danger">*</span></label>
          <select id="wfOrder" v-model="form.order_id" class="form-select" :class="{ 'is-invalid': errors.order_id }" @change="form.model_id = null">
            <option :value="null">-- select --</option>
            <option v-for="o in deliveredOrders" :key="o.id" :value="o.id">{{ o.order_number }} — {{ o.customer_name }}</option>
          </select>
          <div class="invalid-feedback">{{ errors.order_id }}</div>
        </div>
        <div class="col-md-6 mb-3">
          <label class="form-label" for="wfProduct">Product <span class="text-danger">*</span></label>
          <select id="wfProduct" v-model="form.model_id" class="form-select" :class="{ 'is-invalid': errors.model_id }" :disabled="!selectedOrder">
            <option :value="null">-- select --</option>
            <option v-for="it in selectedOrder?.items ?? []" :key="it.model_id" :value="it.model_id">{{ it.product_name }}</option>
          </select>
          <div class="invalid-feedback">{{ errors.model_id }}</div>
        </div>
        <div class="col-md-6 mb-3">
          <label class="form-label" for="wfSerial">Serial number <span class="text-danger">*</span></label>
          <input id="wfSerial" v-model.trim="form.serial_number" class="form-control font-monospace" :class="{ 'is-invalid': errors.serial_number }" />
          <div class="invalid-feedback">{{ errors.serial_number }}</div>
        </div>
        <div class="col-md-3 mb-3">
          <label class="form-label" for="wfExpiry">Warranty expiry</label>
          <input id="wfExpiry" v-model="form.warranty_expiry" type="date" class="form-control" />
        </div>
        <div class="col-md-3 mb-3">
          <label class="form-label" for="wfStatus">Status</label>
          <select id="wfStatus" v-model="form.status" class="form-select">
            <option v-for="(m, key) in WARRANTY_STATUSES" :key="key" :value="key">{{ m.label }}</option>
          </select>
        </div>
        <div class="col-12 mb-3">
          <label class="form-label" for="wfIssue">Issue description <span class="text-danger">*</span></label>
          <textarea id="wfIssue" v-model="form.issue_description" class="form-control" :class="{ 'is-invalid': errors.issue_description }" rows="3"></textarea>
          <div class="invalid-feedback">{{ errors.issue_description }}</div>
        </div>
        <div class="col-12">
          <label class="form-label" for="wfRes">Resolution</label>
          <textarea id="wfRes" v-model="form.resolution" class="form-control" rows="2"></textarea>
        </div>
      </div>
    </form>
    <template #footer>
      <button type="button" class="btn btn-secondary" @click="form = null">Cancel</button>
      <button type="submit" form="warrantyForm" class="btn btn-primary" :disabled="saving">Save</button>
    </template>
  </BaseModal>

  <ConfirmDialog
    v-if="deleting"
    title="Delete claim"
    :message="`Delete warranty claim ${deleting.claim_number}?`"
    :busy="saving"
    @confirm="remove"
    @cancel="deleting = null"
  />
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { warrantyApi } from '@/services/warrantyApi'
import { orderApi } from '@/services/orderApi'
import { WARRANTY_STATUSES } from '@/config/statuses'
import { formatDate } from '@/utils/format'
import { toast } from '@/stores/toast'
import { usePagination } from '@/composables/usePagination'
import BaseModal from '@/components/common/BaseModal.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import PaginationBar from '@/components/common/PaginationBar.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'

const rows = ref([])
const deliveredOrders = ref([])
const loading = ref(true)
const search = ref('')
const status = ref('')
const form = ref(null)
const errors = ref({})
const saving = ref(false)
const deleting = ref(null)

const filtered = computed(() => {
  const term = search.value.toLowerCase()
  return rows.value
    .filter((w) => !status.value || w.status === status.value)
    .filter((w) => !term || [w.claim_number, w.customer_name, w.model_name, w.serial_number].some((v) => v?.toLowerCase().includes(term)))
})
const { page, paged, pageSize } = usePagination(filtered, { pageSize: 15, resetOn: [search, status] })

const selectedOrder = computed(() => deliveredOrders.value.find((o) => o.id === form.value?.order_id))

async function load() {
  loading.value = true
  rows.value = await warrantyApi.list()
  loading.value = false
}

async function loadOrders() {
  deliveredOrders.value = await orderApi.list({ status: 'DELIVERED' })
}

function openCreate() {
  errors.value = {}
  form.value = { order_id: null, model_id: null, serial_number: '', warranty_expiry: null, status: 'RECEIVED', issue_description: '', resolution: '' }
  loadOrders()
}

function openEdit(w) {
  errors.value = {}
  form.value = { ...w }
  loadOrders()
}

async function save() {
  saving.value = true
  errors.value = {}
  try {
    if (form.value.id) await warrantyApi.update(form.value.id, form.value)
    else await warrantyApi.create(form.value)
    toast.success('Warranty claim saved')
    form.value = null
    await load()
  } catch (e) {
    errors.value = e.details ?? {}
    if (!e.details) toast.error(e.message)
  } finally {
    saving.value = false
  }
}

async function remove() {
  saving.value = true
  try {
    await warrantyApi.remove(deleting.value.id)
    toast.success('Claim deleted')
    await load()
  } catch (e) {
    toast.error(e.message)
  } finally {
    saving.value = false
    deleting.value = null
  }
}

onMounted(load)
</script>
