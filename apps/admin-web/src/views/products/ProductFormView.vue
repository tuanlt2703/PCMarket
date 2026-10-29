<template>
  <div v-if="loading" class="text-center py-5"><div class="spinner-border text-primary"></div></div>
  <form v-else novalidate @submit.prevent="save">
    <div class="card card-primary card-outline mb-4">
      <div class="card-header"><h3 class="card-title">Model information</h3></div>
      <div class="card-body">
        <div class="row">
          <div class="col-md-4 mb-3">
            <label class="form-label" for="pfType">Component type <span class="text-danger">*</span></label>
            <select
              id="pfType"
              v-model="form.component_type_id"
              class="form-select"
              :class="{ 'is-invalid': errors.component_type_id }"
              :disabled="isEdit"
              @change="onTypeChange"
            >
              <option :value="null">-- select --</option>
              <option v-for="t in lookups.options('component_types', form.component_type_id)" :key="t.id" :value="t.id">
                {{ t.name }}
              </option>
            </select>
            <div class="invalid-feedback">{{ errors.component_type_id }}</div>
          </div>
          <div class="col-md-4">
            <SpecFieldInput v-model="form.manufacturer_id" :field="manufacturerField" :error="errors.manufacturer_id" />
          </div>
          <div class="col-md-4 mb-3">
            <label class="form-label" for="pfRelease">Release date</label>
            <input id="pfRelease" v-model="form.release_date" type="date" class="form-control" />
          </div>
          <div class="col-md-8 mb-3">
            <label class="form-label" for="pfName">Model name <span class="text-danger">*</span></label>
            <input id="pfName" v-model="form.model_name" class="form-control" :class="{ 'is-invalid': errors.model_name }" maxlength="200" />
            <div class="invalid-feedback">{{ errors.model_name }}</div>
          </div>
          <div class="col-md-4 mb-3 d-flex align-items-end">
            <div class="form-check form-switch mb-2">
              <input id="pfActive" v-model="form.is_active" class="form-check-input" type="checkbox" />
              <label class="form-check-label" for="pfActive">Active (visible for sale)</label>
            </div>
          </div>
          <div class="col-12">
            <label class="form-label" for="pfDesc">Description</label>
            <textarea id="pfDesc" v-model="form.description" class="form-control" rows="3"></textarea>
          </div>
        </div>
      </div>
    </div>

    <div v-if="form.component_type_id && !schema" class="alert alert-info">
      <i class="bi bi-info-circle me-1"></i>
      There is no spec definition for this component type yet — only general information will be stored.
    </div>

    <template v-if="schema">
      <div v-for="section in schema.sections" :key="section.title" class="card mb-4">
        <div class="card-header">
          <h3 class="card-title">{{ section.title }}</h3>
        </div>
        <div class="card-body">
          <div class="row">
            <div
              v-for="f in section.fields"
              :key="f.key"
              :class="f.type === 'rows' || f.type === 'multi' ? 'col-12' : 'col-lg-3 col-md-4 col-sm-6'"
            >
              <RowsEditor
                v-if="f.type === 'rows'"
                v-model="form.spec[f.key]"
                :field="f"
                :errors="specErrors"
              />
              <SpecFieldInput v-else v-model="form.spec[f.key]" :field="f" :error="specErrors[f.key]" />
            </div>
          </div>
        </div>
      </div>
    </template>

    <div class="d-flex justify-content-end gap-2 mb-4">
      <button type="button" class="btn btn-secondary" @click="router.back()">Cancel</button>
      <button type="submit" class="btn btn-primary" :disabled="saving">
        <span v-if="saving" class="spinner-border spinner-border-sm me-1"></span>
        <i v-else class="bi bi-save me-1"></i>Save
      </button>
    </div>
  </form>

  <QuickAddLookupModal v-if="quickAdd" :table="quickAdd.table" @close="onQuickAddClose" />
</template>

<script setup>
import { computed, inject, onMounted, provide, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { productApi } from '@/services/productApi'
import { useLookupStore } from '@/stores/lookups'
import { getSchema, emptySpec } from '@/config/componentSchemas'
import { toast } from '@/stores/toast'
import SpecFieldInput from '@/components/product/SpecFieldInput.vue'
import RowsEditor from '@/components/product/RowsEditor.vue'
import QuickAddLookupModal from '@/components/product/QuickAddLookupModal.vue'

const route = useRoute()
const router = useRouter()
const lookups = useLookupStore()
const setPageTitle = inject('setPageTitle')

const isEdit = computed(() => route.name === 'product-edit')
const loading = ref(true)
const saving = ref(false)
const errors = ref({})
const form = ref(null)

const manufacturerField = { key: 'manufacturer_id', label: 'Manufacturer', type: 'lookup', table: 'manufacturers' }

const typeCode = computed(() => lookups.find('component_types', form.value?.component_type_id)?.code)
const schema = computed(() => getSchema(typeCode.value))

const specErrors = computed(() =>
  Object.fromEntries(
    Object.entries(errors.value)
      .filter(([k]) => k.startsWith('spec.'))
      .map(([k, v]) => [k.slice(5), v]),
  ),
)

function onTypeChange() {
  form.value.spec = emptySpec(schema.value)
  errors.value = {}
}

// Inline creation of lookup values from any select in the form
const quickAdd = ref(null)
provide('quickAddLookup', (table) => new Promise((resolve) => (quickAdd.value = { table, resolve })))
function onQuickAddClose(id) {
  quickAdd.value.resolve(id)
  quickAdd.value = null
}

async function save() {
  saving.value = true
  errors.value = {}
  try {
    const saved = isEdit.value
      ? await productApi.update(form.value.id, form.value)
      : await productApi.create(form.value)
    toast.success(isEdit.value ? 'Product updated' : 'Product created')
    router.push({ name: 'product-detail', params: { id: saved.id } })
  } catch (e) {
    errors.value = e.details ?? {}
    toast.error(e.details ? 'Please fix the highlighted fields' : e.message)
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  if (isEdit.value) {
    try {
      const p = await productApi.get(Number(route.params.id))
      form.value = {
        id: p.id,
        component_type_id: p.component_type_id,
        manufacturer_id: p.manufacturer_id,
        model_name: p.model_name,
        description: p.description,
        release_date: p.release_date,
        is_active: p.is_active,
        spec: { ...emptySpec(getSchema(lookups.find('component_types', p.component_type_id)?.code)), ...p.spec },
      }
      setPageTitle(`Edit: ${p.model_name}`)
    } catch (e) {
      toast.error(e.message)
      router.replace({ name: 'products' })
      return
    }
  } else {
    const type = lookups.typeByCode(route.query.type)
    form.value = {
      component_type_id: type?.id ?? null,
      manufacturer_id: null,
      model_name: '',
      description: '',
      release_date: null,
      is_active: true,
      spec: emptySpec(getSchema(type?.code)),
    }
    if (type) setPageTitle(`New ${type.name}`)
  }
  loading.value = false
})
</script>
