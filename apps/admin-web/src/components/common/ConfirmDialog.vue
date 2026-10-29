<template>
  <BaseModal :title="title" size="sm" @close="$emit('cancel')">
    <p class="mb-0">{{ message }}</p>
    <template #footer>
      <button type="button" class="btn btn-secondary" :disabled="busy" @click="$emit('cancel')">Cancel</button>
      <button type="button" class="btn" :class="`btn-${variant}`" :disabled="busy" @click="$emit('confirm')">
        <span v-if="busy" class="spinner-border spinner-border-sm me-1"></span>{{ confirmText }}
      </button>
    </template>
  </BaseModal>
</template>

<script setup>
import BaseModal from './BaseModal.vue'

defineProps({
  title: { type: String, default: 'Confirm' },
  message: { type: String, required: true },
  confirmText: { type: String, default: 'Delete' },
  variant: { type: String, default: 'danger' },
  busy: { type: Boolean, default: false },
})
defineEmits(['confirm', 'cancel'])
</script>
