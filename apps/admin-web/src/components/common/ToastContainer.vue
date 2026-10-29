<template>
  <div class="toast-stack">
    <div
      v-for="t in toasts"
      :key="t.id"
      class="toast show align-items-center border-0 mb-2"
      :class="`text-bg-${t.type}`"
      role="alert"
      aria-live="assertive"
    >
      <div class="d-flex">
        <div class="toast-body">{{ t.message }}</div>
        <button type="button" class="btn-close btn-close-white me-2 m-auto" aria-label="Close" @click="dismiss(t.id)"></button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { toasts, dismiss } from '@/stores/toast'
</script>
