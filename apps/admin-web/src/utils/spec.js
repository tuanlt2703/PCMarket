import { formatBool, formatNumber } from '@/utils/format'

// Renders a spec value for display using the field definition and the lookup store
export function formatSpecValue(field, value, lookups) {
  switch (field.type) {
    case 'lookup':
      return value == null ? '—' : lookups.name(field.table, value)
    case 'multi':
      return value?.length ? value.map((id) => lookups.name(field.table, id)).join(', ') : '—'
    case 'boolean':
      return formatBool(value)
    case 'number':
      return value == null ? '—' : `${formatNumber(value)}${field.unit ? ` ${field.unit}` : ''}`
    default:
      return value || '—'
  }
}

export function specSummary(schema, spec, lookups) {
  if (!schema || !spec) return ''
  const fields = schema.sections.flatMap((s) => s.fields)
  return schema.summary
    .map((key) => fields.find((f) => f.key === key))
    .filter((f) => f && spec[f.key] != null && spec[f.key] !== '')
    .map((f) => {
      const v = formatSpecValue(f, spec[f.key], lookups)
      return f.type === 'number' && !f.unit ? `${v} ${f.label.toLowerCase()}` : v
    })
    .join(' · ')
}
