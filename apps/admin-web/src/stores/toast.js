import { reactive } from 'vue'

let seq = 0
export const toasts = reactive([])

function push(type, message, timeout = 3500) {
  const id = ++seq
  toasts.push({ id, type, message })
  setTimeout(() => dismiss(id), timeout)
}

export function dismiss(id) {
  const idx = toasts.findIndex((t) => t.id === id)
  if (idx >= 0) toasts.splice(idx, 1)
}

export const toast = {
  success: (msg) => push('success', msg),
  error: (msg) => push('danger', msg, 6000),
  info: (msg) => push('info', msg),
}
