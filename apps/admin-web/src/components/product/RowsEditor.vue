<template>
  <div class="mb-3">
    <div class="d-flex justify-content-between align-items-center mb-1">
      <label class="form-label mb-0">{{ field.label }}</label>
      <button type="button" class="btn btn-sm btn-outline-primary" @click="addRow">
        <i class="bi bi-plus-lg me-1"></i>Add row
      </button>
    </div>
    <div class="table-responsive">
      <table class="table table-sm table-bordered rows-editor mb-0">
        <thead class="table-light">
          <tr>
            <th v-for="c in field.columns" :key="c.key" class="small">
              {{ c.label }}<span v-if="c.unit"> ({{ c.unit }})</span><span v-if="c.required" class="text-danger">*</span>
            </th>
            <th style="width: 44px"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(row, i) in modelValue" :key="i">
            <td v-for="c in field.columns" :key="c.key">
              <SpecFieldInput
                :field="c"
                :model-value="row[c.key]"
                :error="errors[`${field.key}.${i}.${c.key}`]"
                compact
                @update:model-value="updateCell(i, c.key, $event)"
              />
            </td>
            <td class="text-center">
              <button type="button" class="btn btn-sm btn-outline-danger" title="Remove row" @click="removeRow(i)">
                <i class="bi bi-x-lg"></i>
              </button>
            </td>
          </tr>
          <tr v-if="!modelValue?.length">
            <td :colspan="field.columns.length + 1" class="text-center text-secondary small">No rows</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import SpecFieldInput from './SpecFieldInput.vue'
import { emptyRow } from '@/config/componentSchemas'

const props = defineProps({
  field: { type: Object, required: true },
  modelValue: { type: Array, default: () => [] },
  errors: { type: Object, default: () => ({}) },
})
const emit = defineEmits(['update:modelValue'])

function addRow() {
  emit('update:modelValue', [...(props.modelValue ?? []), emptyRow(props.field)])
}

function removeRow(i) {
  emit('update:modelValue', props.modelValue.filter((_, idx) => idx !== i))
}

function updateCell(i, key, value) {
  emit(
    'update:modelValue',
    props.modelValue.map((row, idx) => (idx === i ? { ...row, [key]: value } : row)),
  )
}
</script>
