<template>
  <div v-if="!product" class="text-center py-5"><div class="spinner-border text-primary"></div></div>
  <template v-else>
    <div class="d-flex flex-wrap justify-content-end gap-2 mb-3">
      <RouterLink :to="{ name: 'products', query: { type: typeCode } }" class="btn btn-secondary btn-sm">
        <i class="bi bi-arrow-left me-1"></i>Back to list
      </RouterLink>
      <RouterLink :to="{ name: 'product-edit', params: { id: product.id } }" class="btn btn-primary btn-sm">
        <i class="bi bi-pencil me-1"></i>Edit
      </RouterLink>
    </div>

    <div class="row">
      <div class="col-lg-8">
        <div class="card card-primary card-outline mb-4">
          <div class="card-header"><h3 class="card-title">Model information</h3></div>
          <div class="card-body p-0">
            <table class="table spec-table mb-0">
              <tbody>
                <tr><th>Component type</th><td>{{ lookups.name('component_types', product.component_type_id) }}</td></tr>
                <tr><th>Manufacturer</th><td>{{ lookups.name('manufacturers', product.manufacturer_id) || '—' }}</td></tr>
                <tr><th>Model name</th><td>{{ product.model_name }}</td></tr>
                <tr><th>Release date</th><td>{{ formatDate(product.release_date) }}</td></tr>
                <tr><th>Description</th><td>{{ product.description || '—' }}</td></tr>
                <tr>
                  <th>Status</th>
                  <td>
                    <span class="badge" :class="product.is_active ? 'text-bg-success' : 'text-bg-secondary'">
                      {{ product.is_active ? 'Active' : 'Inactive' }}
                    </span>
                  </td>
                </tr>
                <tr><th>Last updated</th><td>{{ formatDateTime(product.updated_at) }}</td></tr>
              </tbody>
            </table>
          </div>
        </div>

        <div v-for="section in schema?.sections ?? []" :key="section.title" class="card mb-4">
          <div class="card-header"><h3 class="card-title">{{ section.title }}</h3></div>
          <div class="card-body p-0">
            <table class="table spec-table mb-0">
              <tbody>
                <tr v-for="f in section.fields.filter((x) => x.type !== 'rows')" :key="f.key">
                  <th>{{ f.label }}</th>
                  <td>{{ formatSpecValue(f, product.spec[f.key], lookups) }}</td>
                </tr>
              </tbody>
            </table>
            <div v-for="f in section.fields.filter((x) => x.type === 'rows')" :key="f.key" class="p-3 border-top">
              <div class="fw-semibold mb-2">{{ f.label }}</div>
              <table v-if="product.spec[f.key]?.length" class="table table-sm table-bordered mb-0">
                <thead>
                  <tr><th v-for="c in f.columns" :key="c.key" class="small">{{ c.label }}</th></tr>
                </thead>
                <tbody>
                  <tr v-for="(row, i) in product.spec[f.key]" :key="i">
                    <td v-for="c in f.columns" :key="c.key">{{ formatSpecValue(c, row[c.key], lookups) }}</td>
                  </tr>
                </tbody>
              </table>
              <div v-else class="text-secondary small">—</div>
            </div>
          </div>
        </div>
      </div>

      <div class="col-lg-4">
        <div class="card mb-4">
          <div class="card-header">
            <h3 class="card-title">Inventory</h3>
            <div class="card-tools">
              <RouterLink :to="{ name: 'inventory', query: { q: product.sku } }" class="btn btn-tool" title="Manage stock">
                <i class="bi bi-box-arrow-up-right"></i>
              </RouterLink>
            </div>
          </div>
          <div class="card-body">
            <dl class="row mb-0">
              <dt class="col-6">SKU</dt><dd class="col-6 font-monospace">{{ product.sku ?? '—' }}</dd>
              <dt class="col-6">Selling price</dt><dd class="col-6">{{ formatCurrency(product.selling_price) }}</dd>
              <dt class="col-6">On hand</dt><dd class="col-6">{{ product.quantity_on_hand }}</dd>
              <dt class="col-6">Available</dt>
              <dd class="col-6">
                <span class="badge" :class="product.quantity_available > 0 ? 'text-bg-success' : 'text-bg-danger'">
                  {{ product.quantity_available }}
                </span>
              </dd>
            </dl>
          </div>
        </div>

        <div class="card mb-4">
          <div class="card-header"><h3 class="card-title">Recent stock movements</h3></div>
          <div class="card-body p-0">
            <ul class="list-group list-group-flush">
              <li v-for="m in movements.slice(0, 8)" :key="m.id" class="list-group-item small">
                <div class="d-flex justify-content-between">
                  <StatusBadge :value="m.type" :map="STOCK_MOVEMENT_TYPES" />
                  <strong :class="m.quantity >= 0 ? 'text-success' : 'text-danger'">
                    {{ m.quantity > 0 ? '+' : '' }}{{ m.quantity }}
                  </strong>
                </div>
                <div>{{ m.reason }} <span v-if="m.reference" class="text-secondary">({{ m.reference }})</span></div>
                <div class="text-secondary">{{ formatDateTime(m.created_at) }}</div>
              </li>
              <li v-if="!movements.length" class="list-group-item text-secondary small">No movements</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </template>
</template>

<script setup>
import { computed, inject, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { productApi } from '@/services/productApi'
import { inventoryApi } from '@/services/inventoryApi'
import { useLookupStore } from '@/stores/lookups'
import { getSchema } from '@/config/componentSchemas'
import { STOCK_MOVEMENT_TYPES } from '@/config/statuses'
import { formatSpecValue } from '@/utils/spec'
import { formatCurrency, formatDate, formatDateTime } from '@/utils/format'
import { toast } from '@/stores/toast'
import StatusBadge from '@/components/common/StatusBadge.vue'

const route = useRoute()
const router = useRouter()
const lookups = useLookupStore()
const setPageTitle = inject('setPageTitle')

const product = ref(null)
const movements = ref([])

const typeCode = computed(() => lookups.find('component_types', product.value?.component_type_id)?.code)
const schema = computed(() => getSchema(typeCode.value))

onMounted(async () => {
  const id = Number(route.params.id)
  try {
    const [p, m] = await Promise.all([productApi.get(id), inventoryApi.movements({ modelId: id })])
    product.value = p
    movements.value = m
    setPageTitle(p.model_name)
  } catch (e) {
    toast.error(e.message)
    router.replace({ name: 'products' })
  }
})
</script>
