<template>
  <div v-if="!meta" class="alert alert-danger">Unknown lookup table "{{ table }}".</div>
  <div v-else class="card mb-4">
    <div class="card-header">
      <h3 class="card-title">
        {{ meta.label }}
        <small class="text-secondary font-monospace ms-2">{{ meta.table }}</small>
      </h3>
      <div class="card-tools">
        <button type="button" class="btn btn-primary btn-sm" @click="openCreate">
          <i class="bi bi-plus-lg me-1"></i>Add new
        </button>
      </div>
    </div>
    <div class="card-body">
      <div class="row g-2 mb-3">
        <div class="col-md-6">
          <input v-model.trim="search" type="search" class="form-control" placeholder="Search code, name, description…" />
        </div>
        <div class="col-md-3">
          <select v-model="statusFilter" class="form-select">
            <option value="">All statuses</option>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
        </div>
      </div>

      <div v-if="loading" class="text-center py-4"><div class="spinner-border text-primary"></div></div>
      <div v-else class="table-responsive">
        <table class="table table-bordered table-hover align-middle">
          <thead>
            <tr>
              <th style="width: 70px">ID</th>
              <th>Code</th>
              <th>Name</th>
              <th>Description</th>
              <th style="width: 100px">Status</th>
              <th style="width: 160px">Updated</th>
              <th style="width: 130px" class="text-end">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="r in paged" :key="r.id">
              <td>{{ r.id }}</td>
              <td class="font-monospace">{{ r.code }}</td>
              <td>{{ r.name }}</td>
              <td class="text-secondary">{{ r.description }}</td>
              <td>
                <span class="badge" :class="r.is_active ? 'text-bg-success' : 'text-bg-secondary'">
                  {{ r.is_active ? 'Active' : 'Inactive' }}
                </span>
              </td>
              <td>{{ formatDateTime(r.updated_at) }}</td>
              <td class="text-end text-nowrap">
                <button type="button" class="btn btn-sm btn-outline-primary me-1" title="Edit" @click="openEdit(r)">
                  <i class="bi bi-pencil"></i>
                </button>
                <button type="button" class="btn btn-sm btn-outline-secondary me-1" :title="r.is_active ? 'Deactivate' : 'Activate'" @click="toggleActive(r)">
                  <i :class="['bi', r.is_active ? 'bi-toggle-on' : 'bi-toggle-off']"></i>
                </button>
                <button type="button" class="btn btn-sm btn-outline-danger" title="Delete" @click="deleting = r">
                  <i class="bi bi-trash"></i>
                </button>
              </td>
            </tr>
            <tr v-if="!filtered.length">
              <td colspan="7" class="text-center text-secondary">No data</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    <div class="card-footer">
      <PaginationBar v-model="page" :total="filtered.length" :page-size="pageSize" />
    </div>
  </div>

  <BaseModal v-if="form" :title="form.id ? `Edit ${meta.label}` : `New ${meta.label}`" @close="form = null">
    <form id="lookupForm" @submit.prevent="save">
      <div class="mb-3">
        <label class="form-label" for="lkCode">Code <span class="text-danger">*</span></label>
        <input
          id="lkCode"
          v-model.trim="form.code"
          class="form-control font-monospace text-uppercase"
          :class="{ 'is-invalid': errors.code }"
          maxlength="50"
          :disabled="codeLocked"
        />
        <div class="invalid-feedback">{{ errors.code }}</div>
        <div v-if="codeLocked" class="form-text">
          This code is bound to a built-in spec form and cannot be changed.
        </div>
      </div>
      <div class="mb-3">
        <label class="form-label" for="lkName">Name <span class="text-danger">*</span></label>
        <input id="lkName" v-model="form.name" class="form-control" :class="{ 'is-invalid': errors.name }" maxlength="100" />
        <div class="invalid-feedback">{{ errors.name }}</div>
      </div>
      <div class="mb-3">
        <label class="form-label" for="lkDesc">Description</label>
        <textarea id="lkDesc" v-model="form.description" class="form-control" rows="3"></textarea>
      </div>
      <div class="form-check form-switch">
        <input id="lkActive" v-model="form.is_active" class="form-check-input" type="checkbox" />
        <label class="form-check-label" for="lkActive">Active</label>
      </div>
    </form>
    <template #footer>
      <button type="button" class="btn btn-secondary" @click="form = null">Cancel</button>
      <button type="submit" form="lookupForm" class="btn btn-primary" :disabled="saving">
        <span v-if="saving" class="spinner-border spinner-border-sm me-1"></span>Save
      </button>
    </template>
  </BaseModal>

  <ConfirmDialog
    v-if="deleting"
    title="Delete item"
    :message="`Delete '${deleting.name}' (${deleting.code})? This cannot be undone.`"
    :busy="saving"
    @confirm="remove"
    @cancel="deleting = null"
  />
</template>

<script setup>
import { computed, inject, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { getLookupMeta } from '@/config/lookupTables'
import { COMPONENT_SCHEMAS } from '@/config/componentSchemas'
import { lookupApi } from '@/services/lookupApi'
import { useLookupStore } from '@/stores/lookups'
import { toast } from '@/stores/toast'
import { formatDateTime } from '@/utils/format'
import { usePagination } from '@/composables/usePagination'
import BaseModal from '@/components/common/BaseModal.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import PaginationBar from '@/components/common/PaginationBar.vue'

const route = useRoute()
const lookups = useLookupStore()
const setPageTitle = inject('setPageTitle')

const table = route.params.table
const meta = getLookupMeta(table)

const rows = ref([])
const loading = ref(true)
const search = ref('')
const statusFilter = ref('')
const form = ref(null)
const errors = ref({})
const saving = ref(false)
const deleting = ref(null)

const filtered = computed(() => {
  const term = search.value.toLowerCase()
  return rows.value
    .filter((r) => !statusFilter.value || (statusFilter.value === 'active') === r.is_active)
    .filter((r) => !term || [r.code, r.name, r.description].some((v) => v?.toLowerCase().includes(term)))
})
const { page, paged, pageSize } = usePagination(filtered, { pageSize: 15, resetOn: [search, statusFilter] })

const codeLocked = computed(() => {
  if (table !== 'component_types' || !form.value?.id) return false
  const original = rows.value.find((r) => r.id === form.value.id)
  return !!COMPONENT_SCHEMAS[original?.code]
})

async function load() {
  loading.value = true
  rows.value = await lookupApi.list(table)
  loading.value = false
}

function openCreate() {
  errors.value = {}
  form.value = { id: null, code: '', name: '', description: '', is_active: true }
}

function openEdit(r) {
  errors.value = {}
  form.value = { ...r }
}

async function refreshAll() {
  await Promise.all([load(), lookups.load(true)])
}

async function save() {
  saving.value = true
  errors.value = {}
  try {
    if (form.value.id) await lookupApi.update(table, form.value.id, form.value)
    else await lookupApi.create(table, form.value)
    toast.success('Saved successfully')
    form.value = null
    await refreshAll()
  } catch (e) {
    errors.value = e.details ?? {}
    if (!e.details) toast.error(e.message)
  } finally {
    saving.value = false
  }
}

async function toggleActive(r) {
  try {
    await lookupApi.update(table, r.id, { ...r, is_active: !r.is_active })
    toast.success(r.is_active ? 'Deactivated' : 'Activated')
    await refreshAll()
  } catch (e) {
    toast.error(e.message)
  }
}

async function remove() {
  saving.value = true
  try {
    await lookupApi.remove(table, deleting.value.id)
    toast.success('Deleted')
    await refreshAll()
  } catch (e) {
    toast.error(e.message)
  } finally {
    saving.value = false
    deleting.value = null
  }
}

onMounted(() => {
  if (!meta) return
  setPageTitle(meta.label)
  load()
})
</script>
