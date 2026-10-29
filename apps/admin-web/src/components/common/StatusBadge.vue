<template>
  <span class="badge" :class="`text-bg-${meta?.color ?? 'secondary'}`">{{ meta?.label ?? value }}</span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  value: { type: String, default: null },
  map: { type: Object, required: true },
})
const meta = computed(() => props.map[props.value])
</script>
