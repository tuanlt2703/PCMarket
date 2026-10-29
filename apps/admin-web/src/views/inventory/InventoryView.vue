<template>
  <div class="row">
    <div v-for="b in infoBoxes" :key="b.label" class="col-lg-3 col-sm-6">
      <div class="info-box mb-3">
        <span class="info-box-icon shadow-sm" :class="`text-bg-${b.color}`"><i :class="['bi', b.icon]"></i></span>
        <div class="info-box-content">
          <span class="info-box-text">{{ b.label }}</span>
          <span class="info-box-number">{{ b.value }}</span>
        </div>
      </div>
    </div>
  </div>

  <div class="card mb-4">
    <div class="card-header">
      <h3 class="card-title">Stock levels</h3>
      <div class="card-tools">
        <RouterLink :to="{ name: 'stock-movements' }" class="btn btn-sm btn-outline-secondary">
          <i class="bi bi-arrow-left-right me-1"></i>Movements
        </RouterLink>
      </div>
    </div>
    <div class="card-body">
      <div class="row g-2 mb-3">
        <div class="col-md-5">
          <input v-model.trim="search" type="search" class="form-control" placeholder="Search SKU or product…" />
        </div>
        <div class="col-md-4">
          <select v-model="typeId" class="form-select">
            <option :value="null">All component types</option>
            <option v-for="t in lookups.options('component_types')" :key="t.id" :value="t.id">{{ t.name }}</option>
          </select>
        </div>
        <div class="col-md-3">
          <select v-model="stockFilter" class="form-select">
            <option value="">All stock levels</option>
            <option value="low">Low / out of stock</option>
            <option value="out">Out of stock</option>
          </select>
        </div>
      </div>

      <div v-if="loading" class="text-center py-4"><div class="spinner-border text-primary"></div></div>
      <div v-else class="table-responsive">
        <table class="table table-bordered table-hover align-middle">
          <thead>
            <tr>
              <th>SKU</th>
              <th>Product</th>
              <th>Location</th>
              <th class="text-end">Unit cost</th>
              <th class="text-end">Selling price</th>
              <th class="text-center">On hand</th>
              <th class="text-center">Reserved</th>
              <th class="text-center">Available</th>
              <th class="text-center">Reorder at</th>
              <th class="text-end" style="width: 110px">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="r in paged" :key="r.model_id" :class="{ 'table-danger': r.quantity_available <= 0, 'table-warning': r.quantity_available > 0 && r.quantity_available <= r.reorder_level }">
              <td class="font-monospace small">{{ r.sku }}</td>
              <td>
                <RouterLink :to="{ name: 'product-detail', params: { id: r.model_id } }">{{ r.model_name }}</RouterLink>
                <div class="small text-secondary">{{ lookups.name('component_types', r.component_type_id) }}</div>
              </td>
              <td>{{ r.location ?? '—' }}</td>
              <td class="text-end text-nowrap">{{ formatCurrency(r.unit_cost) }}</td>
              <td class="text-end text-nowrap">{{ formatCurrency(r.selling_price) }}</td>
              <td class="text-center">{{ r.quantity_on_hand }}</td>
              <td class="text-center">{{ r.quantity_reserved }}</td>
              <td class="text-center fw-bold">{{ r.quantity_available }}</td>
              <td class="text-center">{{ r.reorder_level }}</td>
              <td class="text-end text-nowrap">
                <button type="button" class="btn btn-sm btn-outline-success me-1" title="Adjust stock" @click="openAdjust(r)">
                  <i class="bi bi-plus-slash-minus"></i>
                </button>
                <button type="button" class="btn btn-sm btn-outline-primary" title="Edit SKU / price" @click="openEdit(r)">
                  <i class="bi bi-pencil"></i>
                </button>
              </td>
            </tr>
            <tr v-if="!filtered.length"><td colspan="10" class="text-center text-secondary">No data</td></tr>
          </tbody>
        </table>
      </div>
    </div>
    <div class="card-footer">
      <PaginationBar v-model="page" :total="filtered.length" :page-size="pageSize" />
    </div>
  </div>

  <BaseModal v-if="adjust" :title="`Adjust stock — ${adjust.model_name}`" @close="adjust = null">
    <form id="adjustForm" @submit.prevent="saveAdjust">
      <p class="mb-3">Current on hand: <strong>{{ adjust.quantity_on_hand }}</strong></p>
      <div class="mb-3">
        <label class="form-label">Movement type</label>
        <div class="btn-group w-100" role="group">
          <template v-for="(m, key) in STOCK_MOVEMENT_TYPES" :key="key">
            <input :id="`mt-${key}`" v-model="adjust.type" type="radio" class="btn-check" :value="key" />
            <label class="btn" :class="`btn-outline-${m.color}`" :for="`mt-${key}`">{{ m.label }}</label>
          </template>
        </div>
      </div>
      <div class="mb-3">
        <label class="form-label" for="adjQty">
          {{ adjust.type === 'ADJUST' ? 'Counted quantity (new on-hand value)' : 'Quantity' }}
        </label>
        <input id="adjQty" v-model.number="adjust.quantity" type="number" min="0" step="1" class="form-control" :class="{ 'is-invalid': errors.quantity }" />
        <div class="invalid-feedback">{{ errors.quantity }}</div>
      </div>
      <div class="mb-3">
        <label class="form-label" for="adjReason">Reason <span class="text-danger">*</span></label>
        <input id="adjReason" v-model="adjust.reason" class="form-control" :class="{ 'is-invalid': errors.reason }" placeholder="e.g. Purchase order received" />
        <div class="invalid-feedback">{{ errors.reason }}</div>
      </div>
      <div>
        <label class="form-label" for="adjRef">Reference</label>
        <input id="adjRef" v-model="adjust.reference" class="form-control" placeholder="e.g. PO-2026-001" />
      </div>
    </form>
    <template #footer>
      <button type="button" class="btn btn-secondary" @click="adjust = null">Cancel</button>
      <button type="submit" form="adjustForm" class="btn btn-primary" :disabled="saving">Save</button>
    </template>
  </BaseModal>

  <BaseModal v-if="edit" :title="`Edit inventory — ${edit.model_name}`" @close="edit = null">
    <form id="invEditForm" @submit.prevent="saveEdit">
      <div class="row">
        <div class="col-md-6 mb-3">
          <label class="form-label" for="ieSku">SKU <span class="text-danger">*</span></label>
          <input id="ieSku" v-model.trim="edit.sku" class="form-control font-monospace text-uppercase" :class="{ 'is-invalid': errors.sku }" />
          <div class="invalid-feedback">{{ errors.sku }}</div>
        </div>
        <div class="col-md-6 mb-3">
          <label class="form-label" for="ieLoc">Location</label>
          <input id="ieLoc" v-model="edit.location" class="form-control" />
        </div>
        <div class="col-md-6 mb-3">
          <label class="form-label" for="ieCost">Unit cost (VND)</label>
          <input id="ieCost" v-model.number="edit.unit_cost" type="number" min="0" step="1000" class="form-control" :class="{ 'is-invalid': errors.unit_cost }" />
          <div class="invalid-feedback">{{ errors.unit_cost }}</div>
        </div>
        <div class="col-md-6 mb-3">
          <label class="form-label" for="iePrice">Selling price (VND)</label>
          <input id="iePrice" v-model.number="edit.selling_price" type="number" min="0" step="1000" class="form-control" :class="{ 'is-invalid': errors.selling_price }" />
          <div class="invalid-feedback">{{ errors.selling_price }}</div>
        </div>
        <div class="col-md-6">
          <label class="form-label" for="ieReorder">Reorder level</label>
          <input id="ieReorder" v-model.number="edit.reorder_level" type="number" min="0" step="1" class="form-control" :class="{ 'is-invalid': errors.reorder_level }" />
          <div class="invalid-feedback">{{ errors.reorder_level }}</div>
        </div>
      </div>
    </form>
    <template #footer>
      <button type="button" class="btn btn-secondary" @click="edit = null">Cancel</button>
      <button type="submit" form="invEditForm" class="btn btn-primary" :disabled="saving">Save</button>
    </template>
  </BaseModal>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { inventoryApi } from '@/services/inventoryApi'
import { useLookupStore } from '@/stores/lookups'
import { STOCK_MOVEMENT_TYPES } from '@/config/statuses'
import { formatCurrency } from '@/utils/format'
import { toast } from '@/stores/toast'
import { usePagination } from '@/composables/usePagination'
import BaseModal from '@/components/common/BaseModal.vue'
import PaginationBar from '@/components/common/PaginationBar.vue'

const route = useRoute()
const lookups = useLookupStore()

const rows = ref([])
const loading = ref(true)
const search = ref(route.query.q ?? '')
const typeId = ref(null)
const stockFilter = ref(route.query.filter ?? '')
const adjust = ref(null)
const edit = ref(null)
const errors = ref({})
const saving = ref(false)

const filtered = computed(() => {
  const term = search.value.toLowerCase()
  return rows.value
    .filter((r) => !typeId.value || r.component_type_id === typeId.value)
    .filter((r) => {
      if (stockFilter.value === 'out') return r.quantity_available <= 0
      if (stockFilter.value === 'low') return r.quantity_available <= r.reorder_level
      return true
    })
    .filter((r) => !term || r.sku.toLowerCase().includes(term) || r.model_name.toLowerCase().includes(term))
})
const { page, paged, pageSize } = usePagination(filtered, { pageSize: 15, resetOn: [search, typeId, stockFilter] })

const infoBoxes = computed(() => [
  { label: 'SKUs', value: rows.value.length, icon: 'bi-upc-scan', color: 'info' },
  { label: 'Units on hand', value: rows.value.reduce((s, r) => s + r.quantity_on_hand, 0), icon: 'bi-boxes', color: 'success' },
  { label: 'Stock value (cost)', value: formatCurrency(rows.value.reduce((s, r) => s + r.quantity_on_hand * r.unit_cost, 0)), icon: 'bi-cash-coin', color: 'warning' },
  { label: 'Low / out of stock', value: rows.value.filter((r) => r.quantity_available <= r.reorder_level).length, icon: 'bi-exclamation-triangle', color: 'danger' },
])

async function load() {
  loading.value = true
  rows.value = await inventoryApi.list()
  loading.value = false
}

function openAdjust(r) {
  errors.value = {}
  adjust.value = { model_id: r.model_id, model_name: r.model_name, quantity_on_hand: r.quantity_on_hand, type: 'IN', quantity: 1, reason: '', reference: '' }
}

function openEdit(r) {
  errors.value = {}
  edit.value = { ...r }
}

async function submit(action, successMsg, close) {
  saving.value = true
  errors.value = {}
  try {
    await action()
    toast.success(successMsg)
    close()
    await load()
  } catch (e) {
    errors.value = e.details ?? {}
    if (!e.details) toast.error(e.message)
  } finally {
    saving.value = false
  }
}

const saveAdjust = () =>
  submit(() => inventoryApi.adjust(adjust.value.model_id, adjust.value), 'Stock updated', () => (adjust.value = null))
const saveEdit = () =>
  submit(() => inventoryApi.update(edit.value.model_id, edit.value), 'Inventory updated', () => (edit.value = null))

onMounted(load)
</script>
