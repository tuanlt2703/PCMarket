<template>
  <div v-if="total > 0" class="d-flex flex-wrap justify-content-between align-items-center gap-2">
    <small class="text-secondary">
      Showing {{ from }}–{{ to }} of {{ total }}
    </small>
    <ul class="pagination pagination-sm m-0">
      <li class="page-item" :class="{ disabled: modelValue <= 1 }">
        <a class="page-link" href="#" @click.prevent="go(modelValue - 1)">&laquo;</a>
      </li>
      <li v-for="p in pages" :key="p" class="page-item" :class="{ active: p === modelValue }">
        <a class="page-link" href="#" @click.prevent="go(p)">{{ p }}</a>
      </li>
      <li class="page-item" :class="{ disabled: modelValue >= pageCount }">
        <a class="page-link" href="#" @click.prevent="go(modelValue + 1)">&raquo;</a>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: { type: Number, required: true },
  total: { type: Number, required: true },
  pageSize: { type: Number, default: 10 },
})
const emit = defineEmits(['update:modelValue'])

const pageCount = computed(() => Math.max(1, Math.ceil(props.total / props.pageSize)))
const from = computed(() => (props.modelValue - 1) * props.pageSize + 1)
const to = computed(() => Math.min(props.total, props.modelValue * props.pageSize))
const pages = computed(() => {
  const start = Math.max(1, Math.min(props.modelValue - 2, pageCount.value - 4))
  const end = Math.min(pageCount.value, start + 4)
  return Array.from({ length: end - start + 1 }, (_, i) => start + i)
})

function go(p) {
  if (p >= 1 && p <= pageCount.value) emit('update:modelValue', p)
}
</script>
