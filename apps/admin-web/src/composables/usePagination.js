import { computed, ref, watch } from 'vue'

// Client-side pagination over a reactive list; resets to page 1 when `resetOn` changes
export function usePagination(items, { pageSize = 10, resetOn } = {}) {
  const page = ref(1)
  const paged = computed(() => items.value.slice((page.value - 1) * pageSize, page.value * pageSize))
  if (resetOn) watch(resetOn, () => (page.value = 1), { deep: true })
  watch(
    () => items.value.length,
    (len) => {
      const max = Math.max(1, Math.ceil(len / pageSize))
      if (page.value > max) page.value = max
    },
  )
  return { page, paged, pageSize }
}
