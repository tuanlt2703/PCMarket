<template>
  <div :class="compact ? '' : 'mb-3'">
    <label v-if="!compact" class="form-label" :for="inputId">
      {{ field.label }}<span v-if="field.unit" class="text-secondary"> ({{ field.unit }})</span>
      <span v-if="field.required" class="text-danger">*</span>
    </label>

    <!-- lookup -->
    <div v-if="field.type === 'lookup'" class="input-group" :class="{ 'input-group-sm': compact, 'has-validation': error }">
      <select
        :id="inputId"
        class="form-select"
        :class="{ 'is-invalid': error, 'form-select-sm': compact }"
        :value="modelValue ?? ''"
        @change="emit('update:modelValue', $event.target.value === '' ? null : Number($event.target.value))"
      >
        <option value="">{{ field.required ? '-- select --' : '—' }}</option>
        <option v-for="o in lookups.options(field.table, modelValue)" :key="o.id" :value="o.id">
          {{ o.name }}{{ o.is_active ? '' : ' (inactive)' }}
        </option>
      </select>
      <button
        type="button"
        class="btn btn-outline-secondary"
        :title="`Add new value to ${field.table}`"
        @click="addLookup"
      >
        <i class="bi bi-plus-lg"></i>
      </button>
      <div v-if="error" class="invalid-feedback">{{ error }}</div>
    </div>

    <!-- multi lookup (junction table) -->
    <div v-else-if="field.type === 'multi'" class="border rounded p-2">
      <div class="d-flex flex-wrap gap-3">
        <div v-for="o in lookups.options(field.table, modelValue)" :key="o.id" class="form-check mb-0">
          <input
            :id="`${inputId}-${o.id}`"
            class="form-check-input"
            type="checkbox"
            :checked="(modelValue ?? []).includes(o.id)"
            @change="toggleMulti(o.id, $event.target.checked)"
          />
          <label class="form-check-label" :for="`${inputId}-${o.id}`">{{ o.name }}</label>
        </div>
        <button type="button" class="btn btn-sm btn-link p-0" @click="addLookup">
          <i class="bi bi-plus-circle me-1"></i>Add value
        </button>
      </div>
    </div>

    <!-- nullable boolean -->
    <select
      v-else-if="field.type === 'boolean'"
      :id="inputId"
      class="form-select"
      :class="{ 'form-select-sm': compact }"
      :value="modelValue === true ? 'true' : modelValue === false ? 'false' : ''"
      @change="emit('update:modelValue', $event.target.value === '' ? null : $event.target.value === 'true')"
    >
      <option value="">—</option>
      <option value="true">Yes</option>
      <option value="false">No</option>
    </select>

    <!-- number -->
    <template v-else-if="field.type === 'number'">
      <input
        :id="inputId"
        type="number"
        class="form-control"
        :class="{ 'is-invalid': error, 'form-control-sm': compact }"
        :step="field.step ?? 'any'"
        :min="field.min ?? undefined"
        :value="modelValue ?? ''"
        @input="emit('update:modelValue', $event.target.value === '' ? null : Number($event.target.value))"
      />
      <div v-if="error" class="invalid-feedback">{{ error }}</div>
    </template>

    <!-- text -->
    <template v-else>
      <input
        :id="inputId"
        type="text"
        class="form-control"
        :class="{ 'is-invalid': error, 'form-control-sm': compact }"
        :maxlength="field.maxLength"
        :value="modelValue ?? ''"
        @input="emit('update:modelValue', $event.target.value === '' ? null : $event.target.value)"
      />
      <div v-if="error" class="invalid-feedback">{{ error }}</div>
    </template>
  </div>
</template>

<script setup>
import { inject, useId } from 'vue'
import { useLookupStore } from '@/stores/lookups'

const props = defineProps({
  field: { type: Object, required: true },
  modelValue: { type: [String, Number, Boolean, Array], default: null },
  error: { type: String, default: null },
  compact: { type: Boolean, default: false },
})
const emit = defineEmits(['update:modelValue'])

const lookups = useLookupStore()
const quickAddLookup = inject('quickAddLookup', null)
const inputId = `spec-${props.field.key}-${useId()}`

function toggleMulti(id, checked) {
  const set = new Set(props.modelValue ?? [])
  checked ? set.add(id) : set.delete(id)
  emit('update:modelValue', [...set])
}

async function addLookup() {
  if (!quickAddLookup) return
  const id = await quickAddLookup(props.field.table)
  if (!id) return
  if (props.field.type === 'multi') toggleMulti(id, true)
  else emit('update:modelValue', id)
}
</script>
