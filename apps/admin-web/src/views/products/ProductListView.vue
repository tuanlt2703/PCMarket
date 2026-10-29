<template>
  <div class="card mb-4">
    <div class="card-header">
      <h3 class="card-title">{{ currentType ? currentType.name : 'All products' }}</h3>
      <div class="card-tools">
        <RouterLink
          :to="{ name: 'product-create', query: currentType ? { type: currentType.code } : {} }"
          class="btn btn-primary btn-sm"
        >
          <i class="bi bi-plus-lg me-1"></i>New product
        </RouterLink>
      </div>
    </div>
    <div class="card-body">
      <div class="row g-2 mb-3">
        <div class="col-lg-4">
          <input v-model.trim="search" type="search" class="form-control" placeholder="Search model name or SKU…" />
        </div>
        <div class="col-lg-3 col-md-4">
          <select v-model="typeCode" class="form-select">
            <option value="">All component types</option>
            <option v-for="t in lookups.options('component_types')" :key="t.id" :value="t.code">{{ t.name }}</option>
          </select>
        </div>
        <div class="col-lg-3 col-md-4">
          <select v-model="manufacturerId" class="form-select">
            <option :value="null">All manufacturers</option>
            <option v-for="m in lookups.options('manufacturers')" :key="m.id" :value="m.id">{{ m.name }}</option>
          </select>
        </div>
        <div class="col-lg-2 col-md-4">
          <select v-model="activeFilter" class="form-select">
            <option :value="null">All statuses</option>
            <option :value="true">Active</option>
            <option :value="false">Inactive</option>
          </select>
        </div>
      </div>

      <div v-if="loading" class="text-center py-4"><div class="spinner-border text-primary"></div></div>
      <div v-else class="table-responsive">
        <table class="table table-hover table-bordered align-middle">
          <thead>
            <tr>
              <th style="width: 60px">ID</th>
              <th>Model</th>
              <th>Type</th>
              <th>Manufacturer</th>
              <th>Key specs</th>
              <th class="text-end">Price</th>
              <th class="text-center">Available</th>
              <th>Status</th>
              <th class="text-end" style="width: 130px">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="p in paged" :key="p.id">
              <td>{{ p.id }}</td>
              <td>
                <RouterLink :to="{ name: 'product-detail', params: { id: p.id } }" class="fw-semibold">{{ p.model_name }}</RouterLink>
                <div class="small text-secondary font-monospace">{{ p.sku }}</div>
              </td>
              <td>{{ lookups.name('component_types', p.component_type_id) }}</td>
              <td>{{ lookups.name('manufacturers', p.manufacturer_id) || '—' }}</td>
              <td class="small">{{ summary(p) }}</td>
              <td class="text-end text-nowrap">{{ formatCurrency(p.selling_price) }}</td>
              <td class="text-center">
                <span class="badge" :class="stockClass(p)">{{ p.quantity_available }}</span>
              </td>
              <td>
                <span class="badge" :class="p.is_active ? 'text-bg-success' : 'text-bg-secondary'">
                  {{ p.is_active ? 'Active' : 'Inactive' }}
                </span>
              </td>
              <td class="text-end text-nowrap">
                <RouterLink :to="{ name: 'product-detail', params: { id: p.id } }" class="btn btn-sm btn-outline-info me-1" title="View">
                  <i class="bi bi-eye"></i>
                </RouterLink>
                <RouterLink :to="{ name: 'product-edit', params: { id: p.id } }" class="btn btn-sm btn-outline-primary me-1" title="Edit">
                  <i class="bi bi-pencil"></i>
                </RouterLink>
                <button type="button" class="btn btn-sm btn-outline-danger" title="Delete" @click="deleting = p">
                  <i class="bi bi-trash"></i>
                </button>
              </td>
            </tr>
            <tr v-if="!products.length">
              <td colspan="9" class="text-center text-secondary">No products found</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    <div class="card-footer">
      <PaginationBar v-model="page" :total="products.length" :page-size="pageSize" />
    </div>
  </div>

  <ConfirmDialog
    v-if="deleting"
    title="Delete product"
    :message="`Delete '${deleting.model_name}' and all of its specs? Its inventory record will be removed too.`"
    :busy="busy"
    @confirm="remove"
    @cancel="deleting = null"
  />
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { productApi } from '@/services/productApi'
import { useLookupStore } from '@/stores/lookups'
import { getSchema } from '@/config/componentSchemas'
import { specSummary } from '@/utils/spec'
import { formatCurrency } from '@/utils/format'
import { toast } from '@/stores/toast'
import { usePagination } from '@/composables/usePagination'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import PaginationBar from '@/components/common/PaginationBar.vue'

const route = useRoute()
const router = useRouter()
const lookups = useLookupStore()

const products = ref([])
const loading = ref(true)
const search = ref('')
const typeCode = ref(route.query.type ?? '')
const manufacturerId = ref(null)
const activeFilter = ref(null)
const deleting = ref(null)
const busy = ref(false)

const currentType = computed(() => (typeCode.value ? lookups.typeByCode(typeCode.value) : null))
const { page, paged, pageSize } = usePagination(products, { pageSize: 15, resetOn: [search, typeCode, manufacturerId, activeFilter] })

const typeCodeById = (id) => lookups.find('component_types', id)?.code
const summary = (p) => specSummary(getSchema(typeCodeById(p.component_type_id)), p.spec, lookups)
const stockClass = (p) =>
  p.quantity_available <= 0 ? 'text-bg-danger' : p.quantity_available <= 5 ? 'text-bg-warning' : 'text-bg-success'

let requestSeq = 0
async function load() {
  const seq = ++requestSeq
  loading.value = true
  const result = await productApi.list({
    componentTypeId: currentType.value?.id,
    manufacturerId: manufacturerId.value,
    active: activeFilter.value,
    q: search.value,
  })
  if (seq !== requestSeq) return
  products.value = result
  loading.value = false
}

async function remove() {
  busy.value = true
  try {
    await productApi.remove(deleting.value.id)
    toast.success('Product deleted')
    await load()
  } catch (e) {
    toast.error(e.message)
  } finally {
    busy.value = false
    deleting.value = null
  }
}

watch(
  () => route.query.type,
  (t) => (typeCode.value = t ?? ''),
)
watch(typeCode, (t) => {
  if ((route.query.type ?? '') !== t) router.replace({ query: t ? { type: t } : {} })
})
watch([search, typeCode, manufacturerId, activeFilter], load, { immediate: true })
</script>
