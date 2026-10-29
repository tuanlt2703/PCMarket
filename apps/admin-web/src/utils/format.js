const vnd = new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND', maximumFractionDigits: 0 })
const number = new Intl.NumberFormat('en-US', { maximumFractionDigits: 3 })

export const formatCurrency = (v) => (v == null ? '—' : vnd.format(v))
export const formatNumber = (v) => (v == null || v === '' ? '—' : number.format(Number(v)))

export function formatDate(v) {
  if (!v) return '—'
  return new Date(v).toLocaleDateString('vi-VN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}

export function formatDateTime(v) {
  if (!v) return '—'
  return new Date(v).toLocaleString('vi-VN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function formatBool(v) {
  if (v === true) return 'Yes'
  if (v === false) return 'No'
  return '—'
}

// Extracts a user-facing message (and field errors) from an ApiError
export function errorMessage(err) {
  return err?.message || 'Unexpected error'
}
