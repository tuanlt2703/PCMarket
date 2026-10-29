// Spec form definitions mapped to product-db *_specs tables and their child/junction tables.
// Keyed by component_types.code.

const GT0 = { gt: 0 }
const GTE0 = { min: 0 }

const num = (key, label, unit, extra = {}) => ({ key, label, type: 'number', unit, step: 'any', ...extra })
const int = (key, label, unit, extra = {}) => ({ key, label, type: 'number', unit, step: 1, integer: true, ...extra })
const bool = (key, label) => ({ key, label, type: 'boolean' })
const text = (key, label, maxLength = 100, extra = {}) => ({ key, label, type: 'text', maxLength, ...extra })
const lookup = (key, label, table, extra = {}) => ({ key, label, type: 'lookup', table, ...extra })
const multi = (key, label, table) => ({ key, label, type: 'multi', table })
// `uniqueBy` mirrors composite primary keys of junction tables
const rows = (key, label, columns, uniqueBy) => ({ key, label, type: 'rows', columns, uniqueBy })

export const COMPONENT_SCHEMAS = {
  CPU: {
    icon: 'bi-cpu',
    summary: ['socket_id', 'core_count', 'thread_count', 'max_boost_clock_mhz'],
    sections: [
      {
        title: 'General',
        fields: [
          lookup('architecture_id', 'Architecture', 'cpu_architectures'),
          lookup('cpu_series_id', 'Series', 'cpu_series'),
          lookup('socket_id', 'Socket', 'cpu_sockets'),
        ],
      },
      {
        title: 'Cores & Clocks',
        fields: [
          int('core_count', 'Cores', null, GT0),
          int('thread_count', 'Threads', null, GT0),
          num('base_clock_mhz', 'Base clock', 'MHz', GT0),
          num('max_boost_clock_mhz', 'Max boost clock', 'MHz', GT0),
        ],
      },
      {
        title: 'Cache',
        fields: [
          num('l1_cache_kb', 'L1 cache', 'KB', GTE0),
          num('l2_cache_mb', 'L2 cache', 'MB', GTE0),
          num('l3_cache_mb', 'L3 cache', 'MB', GTE0),
        ],
      },
      {
        title: 'Memory & I/O',
        fields: [
          num('max_memory_capacity_gb', 'Max memory capacity', 'GB', GT0),
          int('memory_channel_count', 'Memory channels', null, GT0),
          lookup('pcie_generation_id', 'PCIe generation', 'pcie_generations'),
          multi('supported_memory_type_ids', 'Supported memory types', 'memory_types'),
        ],
      },
      {
        title: 'Power, Graphics & Features',
        fields: [
          num('tdp_w', 'TDP', 'W', GTE0),
          num('max_operating_temperature_c', 'Max operating temperature', '°C'),
          bool('has_integrated_graphics', 'Integrated graphics'),
          lookup('integrated_gpu_chipset_id', 'Integrated GPU', 'gpu_chipsets'),
          bool('unlocked_multiplier', 'Unlocked multiplier'),
          bool('cooler_included', 'Cooler included'),
          multi('supported_instruction_set_ids', 'Instruction sets', 'cpu_instruction_sets'),
        ],
      },
    ],
  },

  GPU: {
    icon: 'bi-gpu-card',
    summary: ['chipset_id', 'vram_capacity_gb', 'vram_type_id', 'boost_clock_mhz'],
    sections: [
      {
        title: 'Chip',
        fields: [
          lookup('chipset_id', 'Chipset', 'gpu_chipsets'),
          lookup('architecture_id', 'Architecture', 'gpu_architectures'),
          int('compute_unit_count', 'Compute units / SMs', null, GT0),
          num('base_clock_mhz', 'Base clock', 'MHz', GT0),
          num('boost_clock_mhz', 'Boost clock', 'MHz', GT0),
          lookup('graphics_api_version_id', 'Graphics API', 'graphics_api_versions'),
          bool('ray_tracing_supported', 'Ray tracing'),
          multi('upscaling_technology_ids', 'Upscaling technologies', 'upscaling_technologies'),
        ],
      },
      {
        title: 'Memory',
        fields: [
          num('vram_capacity_gb', 'VRAM capacity', 'GB', GT0),
          lookup('vram_type_id', 'VRAM type', 'vram_types'),
          int('memory_bus_width_bits', 'Memory bus width', 'bit', GT0),
          num('memory_bandwidth_gbps', 'Memory bandwidth', 'GB/s', GT0),
        ],
      },
      {
        title: 'Interface & Power',
        fields: [
          lookup('pcie_generation_id', 'PCIe generation', 'pcie_generations'),
          int('pcie_lane_count', 'PCIe lanes', null, GT0),
          int('recommended_psu_w', 'Recommended PSU', 'W', GT0),
          num('total_graphics_power_w', 'Total graphics power', 'W', GTE0),
          rows('power_connectors', 'Power connectors', [
            lookup('power_connector_type_id', 'Connector', 'power_connector_types', { required: true }),
            int('connector_count', 'Count', null, { required: true, gt: 0, default: 1 }),
          ], ['power_connector_type_id']),
        ],
      },
      {
        title: 'Display',
        fields: [
          int('max_resolution_width_px', 'Max resolution width', 'px', GT0),
          int('max_resolution_height_px', 'Max resolution height', 'px', GT0),
          int('max_displays_supported', 'Max displays', null, GT0),
          rows('display_outputs', 'Display outputs', [
            lookup('display_connector_type_id', 'Connector', 'display_connector_types', { required: true }),
            int('port_count', 'Ports', null, { required: true, gt: 0, default: 1 }),
          ], ['display_connector_type_id']),
        ],
      },
      {
        title: 'Physical & Cooling',
        fields: [
          num('length_mm', 'Length', 'mm', GT0),
          num('thickness_slots', 'Thickness', 'slots', GT0),
          lookup('cooling_type_id', 'Cooling type', 'cooling_types'),
        ],
      },
    ],
  },

  MOTHERBOARD: {
    icon: 'bi-motherboard',
    summary: ['socket_id', 'chipset_id', 'form_factor_id'],
    sections: [
      {
        title: 'General',
        fields: [
          lookup('chipset_id', 'Chipset', 'motherboard_chipsets'),
          lookup('socket_id', 'CPU socket', 'cpu_sockets'),
          lookup('form_factor_id', 'Form factor', 'motherboard_form_factors'),
          multi('supported_cpu_series_ids', 'Supported CPU series', 'cpu_series'),
        ],
      },
      {
        title: 'Memory',
        fields: [
          int('memory_slot_count', 'Memory slots', null, GTE0),
          num('max_memory_capacity_gb', 'Max memory capacity', 'GB', GT0),
          int('max_memory_data_rate_mt_s', 'Max memory speed', 'MT/s', GT0),
          int('memory_channel_count', 'Memory channels', null, GT0),
          multi('supported_memory_type_ids', 'Supported memory types', 'memory_types'),
        ],
      },
      {
        title: 'Expansion & Storage',
        fields: [
          int('sata_port_count', 'SATA ports', null, GTE0),
          rows('expansion_slots', 'Expansion slots', [
            lookup('pcie_generation_id', 'PCIe gen', 'pcie_generations'),
            int('lane_count', 'Lanes', null, GT0),
            int('slot_count', 'Slots', null, { required: true, gt: 0, default: 1 }),
            text('description', 'Description', 200),
          ]),
          rows('storage_slots', 'Storage slots', [
            lookup('storage_interface_id', 'Interface', 'storage_interfaces', { required: true }),
            lookup('pcie_generation_id', 'PCIe gen', 'pcie_generations'),
            int('slot_count', 'Slots', null, { required: true, gt: 0, default: 1 }),
            lookup('form_factor_id', 'Form factor', 'storage_form_factors'),
            text('description', 'Description', 200),
          ]),
        ],
      },
      {
        title: 'Connectivity',
        fields: [
          int('ethernet_speed_mbps', 'Ethernet speed', 'Mbps', GT0),
          multi('wireless_standard_ids', 'Wireless standards', 'wireless_standards'),
          multi('bluetooth_version_ids', 'Bluetooth versions', 'bluetooth_versions'),
          rows('usb_ports', 'USB ports', [
            lookup('usb_port_type_id', 'Port type', 'usb_port_types', { required: true }),
            int('port_count', 'Ports', null, { required: true, gt: 0, default: 1 }),
          ], ['usb_port_type_id']),
          rows('display_outputs', 'Display outputs', [
            lookup('display_connector_type_id', 'Connector', 'display_connector_types', { required: true }),
            int('port_count', 'Ports', null, { required: true, gt: 0, default: 1 }),
          ], ['display_connector_type_id']),
        ],
      },
      {
        title: 'Audio & Features',
        fields: [
          lookup('audio_codec_id', 'Audio codec', 'audio_codecs'),
          text('vrm_phase_description', 'VRM phases', 100),
          bool('bios_flashback_supported', 'BIOS flashback'),
          int('rgb_header_count', 'RGB headers', null, GTE0),
        ],
      },
    ],
  },

  RAM: {
    icon: 'bi-memory',
    summary: ['memory_type_id', 'total_capacity_gb', 'data_rate_mt_s', 'cas_latency_cl'],
    sections: [
      {
        title: 'General',
        fields: [
          lookup('memory_type_id', 'Memory type', 'memory_types', { required: true }),
          lookup('module_type_id', 'Module type', 'memory_module_types'),
          num('total_capacity_gb', 'Total capacity', 'GB', GT0),
          int('module_count', 'Modules', null, GT0),
          num('capacity_per_module_gb', 'Capacity per module', 'GB', GT0),
        ],
      },
      {
        title: 'Performance',
        fields: [
          int('data_rate_mt_s', 'Data rate', 'MT/s', GT0),
          int('cas_latency_cl', 'CAS latency', 'CL', GT0),
          int('timing_t_r_c_d', 'tRCD', null, GT0),
          int('timing_t_r_p', 'tRP', null, GT0),
          int('timing_t_r_a_s', 'tRAS', null, GT0),
          num('operating_voltage_v', 'Operating voltage', 'V', { gt: 0, step: 0.001 }),
        ],
      },
      {
        title: 'Features',
        fields: [
          bool('ecc_supported', 'ECC'),
          bool('registered_buffered', 'Registered / buffered'),
          bool('on_die_ecc', 'On-die ECC'),
          bool('xmp_supported', 'XMP'),
          bool('expo_supported', 'EXPO'),
          rows('memory_profiles', 'Memory profiles', [
            lookup('profile_type_id', 'Profile type', 'ram_profile_types', { required: true }),
            text('profile_name', 'Name', 100),
            int('data_rate_mt_s', 'MT/s', null, { required: true, gt: 0 }),
            num('voltage_v', 'Voltage (V)', null, { gt: 0, step: 0.001 }),
            int('cas_latency_cl', 'CL', null, GT0),
            int('timing_t_r_c_d', 'tRCD'),
            int('timing_t_r_p', 'tRP'),
            int('timing_t_r_a_s', 'tRAS'),
          ]),
        ],
      },
      {
        title: 'Physical',
        fields: [
          lookup('heat_spreader_type_id', 'Heat spreader', 'ram_heat_spreader_types'),
          bool('rgb_supported', 'RGB'),
          lookup('color_id', 'Color', 'product_colors'),
          num('module_height_mm', 'Module height', 'mm', GT0),
        ],
      },
    ],
  },

  STORAGE: {
    icon: 'bi-device-ssd',
    summary: ['device_type_id', 'capacity_gb', 'interface_id', 'sequential_read_mbps'],
    sections: [
      {
        title: 'General',
        fields: [
          lookup('device_type_id', 'Device type', 'storage_device_types', { required: true }),
          num('capacity_gb', 'Capacity', 'GB', GT0),
          lookup('form_factor_id', 'Form factor', 'storage_form_factors'),
          lookup('interface_id', 'Interface', 'storage_interfaces'),
          lookup('protocol_id', 'Protocol', 'storage_protocols'),
          lookup('pcie_generation_id', 'PCIe generation', 'pcie_generations'),
        ],
      },
      {
        title: 'Performance',
        fields: [
          int('sequential_read_mbps', 'Sequential read', 'MB/s', GT0),
          int('sequential_write_mbps', 'Sequential write', 'MB/s', GT0),
          int('random_read_iops', 'Random read', 'IOPS', GT0),
          int('random_write_iops', 'Random write', 'IOPS', GT0),
          int('rpm', 'Spindle speed', 'RPM', GT0),
          int('cache_mb', 'Cache', 'MB', GTE0),
        ],
      },
      {
        title: 'NAND, Reliability & Security',
        fields: [
          lookup('nand_type_id', 'NAND type', 'nand_types'),
          bool('dram_cache_present', 'DRAM cache'),
          num('endurance_tbw', 'Endurance', 'TBW', GTE0),
          int('mtbf_hours', 'MTBF', 'hours', GT0),
          num('warranty_years', 'Warranty', 'years', { min: 0, step: 0.5 }),
          bool('power_loss_protection', 'Power loss protection'),
          multi('encryption_standard_ids', 'Encryption standards', 'encryption_standards'),
        ],
      },
    ],
  },

  PSU: {
    icon: 'bi-lightning-charge',
    summary: ['rated_power_w', 'efficiency_certification_id', 'modular_type_id'],
    sections: [
      {
        title: 'General',
        fields: [
          int('rated_power_w', 'Rated power', 'W', GT0),
          lookup('efficiency_certification_id', 'Efficiency', 'psu_efficiency_certifications'),
          lookup('psu_standard_id', 'Standard', 'psu_standards'),
          lookup('form_factor_id', 'Form factor', 'psu_form_factors'),
          lookup('modular_type_id', 'Modular type', 'psu_modular_types'),
        ],
      },
      {
        title: 'Electrical',
        fields: [
          num('input_voltage_min_v', 'Input voltage min', 'V'),
          num('input_voltage_max_v', 'Input voltage max', 'V'),
          num('frequency_min_hz', 'Frequency min', 'Hz'),
          num('frequency_max_hz', 'Frequency max', 'Hz'),
          num('twelve_v_output_w', '+12V output', 'W', GTE0),
          multi('protection_type_ids', 'Protection features', 'psu_protection_types'),
        ],
      },
      {
        title: 'Connectors',
        fields: [
          int('sata_power_connector_count', 'SATA power connectors', null, GTE0),
          int('molex_connector_count', 'Molex connectors', null, GTE0),
          rows('power_connectors', 'Power connectors', [
            lookup('power_connector_type_id', 'Connector', 'power_connector_types', { required: true }),
            int('connector_count', 'Count', null, { required: true, gt: 0, default: 1 }),
          ], ['power_connector_type_id']),
        ],
      },
      {
        title: 'Cooling & Physical',
        fields: [
          lookup('fan_size_id', 'Fan size', 'fan_sizes'),
          bool('zero_rpm_mode', 'Zero RPM mode'),
          num('length_mm', 'Length', 'mm'),
          num('width_mm', 'Width', 'mm'),
          num('height_mm', 'Height', 'mm'),
          num('warranty_years', 'Warranty', 'years', { min: 0, step: 0.5 }),
        ],
      },
    ],
  },

  CASE: {
    icon: 'bi-pc-display',
    summary: ['case_type_id', 'max_gpu_length_mm', 'color_id'],
    sections: [
      {
        title: 'General',
        fields: [
          lookup('case_type_id', 'Case type', 'case_types'),
          lookup('color_id', 'Color', 'product_colors'),
          multi('supported_motherboard_form_factor_ids', 'Supported motherboards', 'motherboard_form_factors'),
          lookup('side_panel_material_id', 'Side panel', 'case_panel_materials'),
          lookup('front_panel_material_id', 'Front panel', 'case_panel_materials'),
          bool('dust_filters_present', 'Dust filters'),
          bool('rgb_supported', 'RGB'),
          bool('vertical_gpu_mount_supported', 'Vertical GPU mount'),
        ],
      },
      {
        title: 'Dimensions & Clearance',
        fields: [
          num('height_mm', 'Height', 'mm'),
          num('width_mm', 'Width', 'mm'),
          num('depth_mm', 'Depth', 'mm'),
          num('volume_l', 'Volume', 'L'),
          num('max_gpu_length_mm', 'Max GPU length', 'mm'),
          num('max_cpu_cooler_height_mm', 'Max CPU cooler height', 'mm'),
          num('max_psu_length_mm', 'Max PSU length', 'mm'),
          num('cable_management_space_mm', 'Cable management space', 'mm'),
          int('expansion_slot_count', 'Expansion slots', null, GTE0),
        ],
      },
      {
        title: 'Cooling Support',
        fields: [
          int('included_fan_count', 'Included fans', null, GTE0),
          int('max_fan_count', 'Max fans', null, GTE0),
          rows('fan_support', 'Fan support', [
            lookup('fan_size_id', 'Fan size', 'fan_sizes', { required: true }),
            int('supported_count', 'Count', null, { required: true, gt: 0, default: 1 }),
          ], ['fan_size_id']),
          rows('radiator_support', 'Radiator support', [
            lookup('radiator_size_id', 'Radiator size', 'radiator_sizes', { required: true }),
            int('supported_count', 'Count', null, { required: true, gt: 0, default: 1 }),
          ], ['radiator_size_id']),
        ],
      },
      {
        title: 'Drive Bays & Front I/O',
        fields: [
          rows('drive_bays', 'Drive bays', [
            lookup('drive_bay_type_id', 'Bay type', 'drive_bay_types', { required: true }),
            int('bay_count', 'Count', null, { required: true, gt: 0, default: 1 }),
          ], ['drive_bay_type_id']),
          rows('front_io_ports', 'Front I/O ports', [
            lookup('io_port_type_id', 'Port type', 'io_port_types', { required: true }),
            int('port_count', 'Ports', null, { required: true, gt: 0, default: 1 }),
          ], ['io_port_type_id']),
        ],
      },
    ],
  },

  COOLING: {
    icon: 'bi-fan',
    summary: ['cooling_type_id', 'rated_tdp_w', 'radiator_size_id'],
    sections: [
      {
        title: 'General',
        fields: [
          lookup('cooling_type_id', 'Cooling type', 'cooling_types', { required: true }),
          int('rated_tdp_w', 'Rated TDP', 'W', GT0),
          multi('supported_cpu_socket_ids', 'Supported sockets', 'cpu_sockets'),
          lookup('radiator_size_id', 'Radiator size', 'radiator_sizes'),
          num('cooler_height_mm', 'Cooler height', 'mm'),
          lookup('lighting_type_id', 'Lighting', 'lighting_types'),
          bool('thermal_paste_included', 'Thermal paste included'),
          num('warranty_years', 'Warranty', 'years', { min: 0, step: 0.5 }),
        ],
      },
      {
        title: 'Fan (summary)',
        fields: [
          num('fan_size_mm', 'Fan size', 'mm'),
          int('fan_count', 'Fan count', null, GTE0),
          int('fan_speed_min_rpm', 'Fan speed min', 'RPM'),
          int('fan_speed_max_rpm', 'Fan speed max', 'RPM'),
          num('max_airflow_cfm', 'Max airflow', 'CFM'),
          num('static_pressure_mmh2o', 'Static pressure', 'mmH₂O', { step: 0.001 }),
          num('noise_level_dba', 'Noise level', 'dBA'),
          lookup('fan_bearing_type_id', 'Bearing', 'fan_bearing_types'),
          lookup('fan_connector_type_id', 'Fan connector', 'electrical_connector_types'),
        ],
      },
      {
        title: 'Pump',
        fields: [
          int('pump_speed_min_rpm', 'Pump speed min', 'RPM'),
          int('pump_speed_max_rpm', 'Pump speed max', 'RPM'),
          lookup('pump_connector_type_id', 'Pump connector', 'electrical_connector_types'),
        ],
      },
      {
        title: 'Fans (detail)',
        fields: [
          rows('fans', 'Fans', [
            lookup('fan_size_id', 'Size', 'fan_sizes'),
            int('fan_count', 'Count', null, { required: true, gt: 0, default: 1 }),
            int('speed_min_rpm', 'Min RPM'),
            int('speed_max_rpm', 'Max RPM'),
            num('max_airflow_cfm', 'CFM'),
            num('static_pressure_mmh2o', 'mmH₂O', null, { step: 0.001 }),
            num('noise_level_dba', 'dBA'),
            lookup('fan_bearing_type_id', 'Bearing', 'fan_bearing_types'),
            lookup('fan_connector_type_id', 'Connector', 'electrical_connector_types'),
          ]),
        ],
      },
    ],
  },
}

export function getSchema(componentTypeCode) {
  return COMPONENT_SCHEMAS[componentTypeCode] ?? null
}

export function schemaFields(schema) {
  return schema ? schema.sections.flatMap((s) => s.fields) : []
}

export function emptySpec(schema) {
  const spec = {}
  for (const f of schemaFields(schema)) {
    spec[f.key] = f.type === 'multi' || f.type === 'rows' ? [] : null
  }
  return spec
}

export function emptyRow(rowsField) {
  const row = {}
  for (const c of rowsField.columns) row[c.key] = c.default ?? null
  return row
}

function checkNumber(field, value) {
  if (value === null || value === undefined || value === '') return null
  const n = Number(value)
  if (Number.isNaN(n)) return 'Must be a number'
  if (field.integer && !Number.isInteger(n)) return 'Must be an integer'
  if (field.gt !== undefined && !(n > field.gt)) return `Must be greater than ${field.gt}`
  if (field.min !== undefined && n < field.min) return `Must be at least ${field.min}`
  return null
}

function checkField(field, value) {
  const empty = value === null || value === undefined || value === ''
  if (field.required && empty) return 'Required'
  if (field.type === 'number') return checkNumber(field, value)
  if (field.type === 'text' && !empty && field.maxLength && String(value).length > field.maxLength) {
    return `Max ${field.maxLength} characters`
  }
  return null
}

// Returns { [path]: message }, path is `key` or `key.rowIndex.column`
export function validateSpec(schema, spec) {
  const errors = {}
  for (const field of schemaFields(schema)) {
    const value = spec[field.key]
    if (field.type === 'rows') {
      const seen = new Set()
      ;(value ?? []).forEach((row, i) => {
        for (const col of field.columns) {
          const msg = checkField(col, row[col.key])
          if (msg) errors[`${field.key}.${i}.${col.key}`] = msg
        }
        if (field.uniqueBy) {
          const k = field.uniqueBy.map((c) => row[c]).join('|')
          if (seen.has(k)) errors[`${field.key}.${i}.${field.uniqueBy[0]}`] = 'Duplicated'
          seen.add(k)
        }
      })
    } else if (field.type !== 'multi') {
      const msg = checkField(field, value)
      if (msg) errors[field.key] = msg
    }
  }
  return errors
}

// Converts form strings to numbers / nulls before saving
export function normalizeSpec(schema, spec) {
  const toValue = (field, v) => {
    if (v === '' || v === undefined) return null
    if (field.type === 'number' || field.type === 'lookup') return v === null ? null : Number(v)
    return v
  }
  const out = {}
  for (const field of schemaFields(schema)) {
    const v = spec[field.key]
    if (field.type === 'multi') out[field.key] = [...new Set((v ?? []).map(Number))]
    else if (field.type === 'rows') {
      out[field.key] = (v ?? []).map((row) => {
        const r = {}
        for (const c of field.columns) r[c.key] = toValue(c, row[c.key])
        return r
      })
    } else out[field.key] = toValue(field, v)
  }
  return out
}

// Lists every (table, id) referenced by a spec — used for lookup delete protection
export function collectLookupRefs(schema, spec) {
  const refs = []
  if (!schema || !spec) return refs
  for (const field of schemaFields(schema)) {
    const v = spec[field.key]
    if (field.type === 'lookup' && v != null) refs.push([field.table, v])
    else if (field.type === 'multi') (v ?? []).forEach((id) => refs.push([field.table, id]))
    else if (field.type === 'rows') {
      for (const row of v ?? []) {
        for (const c of field.columns) {
          if (c.type === 'lookup' && row[c.key] != null) refs.push([c.table, row[c.key]])
        }
      }
    }
  }
  return refs
}
