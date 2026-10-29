<template>
  <BaseModal :title="`Add to ${meta?.label ?? table}`" @close="$emit('close', null)">
    <form id="quickLookupForm" @submit.prevent="save">
      <div class="mb-3">
        <label class="form-label" for="qlCode">Code <span class="text-danger">*</span></label>
        <input id="qlCode" v-model.trim="form.code" class="form-control font-monospace text-uppercase" :class="{ 'is-invalid': errors.code }" maxlength="50" placeholder="e.g. DDR6" />
        <div class="invalid-feedback">{{ errors.code }}</div>
      </div>
      <div class="mb-3">
        <label class="form-label" for="qlName">Name <span class="text-danger">*</span></label>
        <input id="qlName" v-model="form.name" class="form-control" :class="{ 'is-invalid': errors.name }" maxlength="100" />
        <div class="invalid-feedback">{{ errors.name }}</div>
      </div>
      <div>
        <label class="form-label" for="qlDesc">Description</label>
        <textarea id="qlDesc" v-model="form.description" class="form-control" rows="2"></textarea>
      </div>
    </form>
    <template #footer>
      <button type="button" class="btn btn-secondary" @click="$emit('close', null)">Cancel</button>
      <button type="submit" form="quickLookupForm" class="btn btn-primary" :disabled="saving">Add</button>
    </template>
  </BaseModal>
</template>

<script setup>
import { reactive, ref } from 'vue'
import BaseModal from '@/components/common/BaseModal.vue'
import { getLookupMeta } from '@/config/lookupTables'
import { lookupApi } from '@/services/lookupApi'
import { useLookupStore } from '@/stores/lookups'
import { toast } from '@/stores/toast'

const props = defineProps({ table: { type: String, required: true } })
const emit = defineEmits(['close'])

const meta = getLookupMeta(props.table)
const lookups = useLookupStore()
const form = reactive({ code: '', name: '', description: '', is_active: true })
const errors = ref({})
const saving = ref(false)

async function save() {
  saving.value = true
  errors.value = {}
  try {
    const row = await lookupApi.create(props.table, form)
    await lookups.load(true)
    toast.success(`Added "${row.name}"`)
    emit('close', row.id)
  } catch (e) {
    errors.value = e.details ?? {}
    if (!e.details) toast.error(e.message)
  } finally {
    saving.value = false
  }
}
</script>
