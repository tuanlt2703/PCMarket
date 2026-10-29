<template>
  <Teleport to="body">
    <div class="modal d-block" tabindex="-1" role="dialog" aria-modal="true" @keydown.esc="$emit('close')">
      <div class="modal-dialog modal-dialog-scrollable" :class="sizeClass">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ title }}</h5>
            <button type="button" class="btn-close" aria-label="Close" @click="$emit('close')"></button>
          </div>
          <div class="modal-body">
            <slot />
          </div>
          <div v-if="$slots.footer" class="modal-footer">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </div>
    <div class="modal-backdrop-custom" @click="$emit('close')"></div>
  </Teleport>
</template>

<script setup>
import { computed, onMounted, onBeforeUnmount } from 'vue'

const props = defineProps({
  title: { type: String, required: true },
  size: { type: String, default: '' }, // sm | lg | xl
})
defineEmits(['close'])

const sizeClass = computed(() => (props.size ? `modal-${props.size}` : ''))

onMounted(() => document.body.classList.add('modal-open'))
onBeforeUnmount(() => document.body.classList.remove('modal-open'))
</script>
