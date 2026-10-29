<template>
  <BaseModal :title="customer?.id ? 'Edit customer' : 'New customer'" size="lg" @close="$emit('close')">
    <form id="customerForm" @submit.prevent="save">
      <div class="row">
        <div class="col-md-6 mb-3">
          <label class="form-label" for="cfName">Full name <span class="text-danger">*</span></label>
          <input id="cfName" v-model="form.full_name" class="form-control" :class="{ 'is-invalid': errors.full_name }" maxlength="150" />
          <div class="invalid-feedback">{{ errors.full_name }}</div>
        </div>
        <div class="col-md-6 mb-3">
          <label class="form-label" for="cfEmail">Email <span class="text-danger">*</span></label>
          <input id="cfEmail" v-model.trim="form.email" type="email" class="form-control" :class="{ 'is-invalid': errors.email }" />
          <div class="invalid-feedback">{{ errors.email }}</div>
        </div>
        <div class="col-md-6 mb-3">
          <label class="form-label" for="cfPhone">Phone</label>
          <input id="cfPhone" v-model.trim="form.phone" class="form-control" :class="{ 'is-invalid': errors.phone }" placeholder="0901234567" />
          <div class="invalid-feedback">{{ errors.phone }}</div>
        </div>
        <div class="col-md-3 mb-3">
          <label class="form-label" for="cfTier">Loyalty tier</label>
          <select id="cfTier" v-model="form.tier" class="form-select" :class="{ 'is-invalid': errors.tier }">
            <option v-for="(m, key) in CUSTOMER_TIERS" :key="key" :value="key">{{ m.label }}</option>
          </select>
        </div>
        <div class="col-md-3 mb-3">
          <label class="form-label" for="cfPoints">Points</label>
          <input id="cfPoints" v-model.number="form.loyalty_points" type="number" min="0" step="1" class="form-control" :class="{ 'is-invalid': errors.loyalty_points }" />
          <div class="invalid-feedback">{{ errors.loyalty_points }}</div>
        </div>
        <div class="col-12 mb-3">
          <label class="form-label" for="cfAddr">Address</label>
          <textarea id="cfAddr" v-model="form.address" class="form-control" rows="2"></textarea>
        </div>
        <div class="col-12">
          <div class="form-check form-switch">
            <input id="cfActive" v-model="form.is_active" class="form-check-input" type="checkbox" />
            <label class="form-check-label" for="cfActive">Active account</label>
          </div>
        </div>
      </div>
    </form>
    <template #footer>
      <button type="button" class="btn btn-secondary" @click="$emit('close')">Cancel</button>
      <button type="submit" form="customerForm" class="btn btn-primary" :disabled="saving">Save</button>
    </template>
  </BaseModal>
</template>

<script setup>
import { reactive, ref } from 'vue'
import BaseModal from '@/components/common/BaseModal.vue'
import { customerApi } from '@/services/customerApi'
import { CUSTOMER_TIERS } from '@/config/statuses'
import { toast } from '@/stores/toast'

const props = defineProps({ customer: { type: Object, default: null } })
const emit = defineEmits(['close', 'saved'])

const form = reactive({
  full_name: '',
  email: '',
  phone: '',
  address: '',
  tier: 'MEMBER',
  loyalty_points: 0,
  is_active: true,
  ...(props.customer ?? {}),
})
const errors = ref({})
const saving = ref(false)

async function save() {
  saving.value = true
  errors.value = {}
  try {
    const saved = props.customer?.id
      ? await customerApi.update(props.customer.id, form)
      : await customerApi.create(form)
    toast.success('Customer saved')
    emit('saved', saved)
  } catch (e) {
    errors.value = e.details ?? {}
    if (!e.details) toast.error(e.message)
  } finally {
    saving.value = false
  }
}
</script>
