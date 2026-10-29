// Lookup tables of product-db (all share: id, code, name, description, is_active, created_at, updated_at)
export const LOOKUP_GROUPS = [
  {
    key: 'catalog',
    label: 'Catalog & Shared',
    icon: 'bi-collection',
    tables: [
      { table: 'component_types', label: 'Component Types' },
      { table: 'manufacturers', label: 'Manufacturers' },
      { table: 'cpu_sockets', label: 'CPU Sockets' },
      { table: 'memory_types', label: 'Memory Types' },
      { table: 'pcie_generations', label: 'PCIe Generations' },
      { table: 'motherboard_form_factors', label: 'Motherboard Form Factors' },
      { table: 'display_connector_types', label: 'Display Connector Types' },
      { table: 'power_connector_types', label: 'Power Connector Types' },
      { table: 'electrical_connector_types', label: 'Electrical Connector Types' },
      { table: 'product_colors', label: 'Product Colors' },
      { table: 'lighting_types', label: 'Lighting Types' },
    ],
  },
  {
    key: 'cpu',
    label: 'CPU',
    icon: 'bi-cpu',
    tables: [
      { table: 'cpu_architectures', label: 'CPU Architectures' },
      { table: 'cpu_series', label: 'CPU Series' },
      { table: 'cpu_instruction_sets', label: 'CPU Instruction Sets' },
    ],
  },
  {
    key: 'gpu',
    label: 'GPU',
    icon: 'bi-gpu-card',
    tables: [
      { table: 'gpu_chipsets', label: 'GPU Chipsets' },
      { table: 'gpu_architectures', label: 'GPU Architectures' },
      { table: 'vram_types', label: 'VRAM Types' },
      { table: 'graphics_api_versions', label: 'Graphics API Versions' },
      { table: 'upscaling_technologies', label: 'Upscaling Technologies' },
    ],
  },
  {
    key: 'motherboard',
    label: 'Motherboard',
    icon: 'bi-motherboard',
    tables: [
      { table: 'motherboard_chipsets', label: 'Motherboard Chipsets' },
      { table: 'audio_codecs', label: 'Audio Codecs' },
      { table: 'usb_port_types', label: 'USB Port Types' },
      { table: 'wireless_standards', label: 'Wireless Standards' },
      { table: 'bluetooth_versions', label: 'Bluetooth Versions' },
    ],
  },
  {
    key: 'ram',
    label: 'RAM',
    icon: 'bi-memory',
    tables: [
      { table: 'memory_module_types', label: 'Memory Module Types' },
      { table: 'ram_heat_spreader_types', label: 'RAM Heat Spreader Types' },
      { table: 'ram_profile_types', label: 'RAM Profile Types' },
    ],
  },
  {
    key: 'storage',
    label: 'Storage',
    icon: 'bi-device-ssd',
    tables: [
      { table: 'storage_device_types', label: 'Storage Device Types' },
      { table: 'storage_form_factors', label: 'Storage Form Factors' },
      { table: 'storage_interfaces', label: 'Storage Interfaces' },
      { table: 'storage_protocols', label: 'Storage Protocols' },
      { table: 'nand_types', label: 'NAND Types' },
      { table: 'encryption_standards', label: 'Encryption Standards' },
    ],
  },
  {
    key: 'psu',
    label: 'Power Supply',
    icon: 'bi-lightning-charge',
    tables: [
      { table: 'psu_efficiency_certifications', label: 'PSU Efficiency Certifications' },
      { table: 'psu_standards', label: 'PSU Standards' },
      { table: 'psu_form_factors', label: 'PSU Form Factors' },
      { table: 'psu_modular_types', label: 'PSU Modular Types' },
      { table: 'psu_protection_types', label: 'PSU Protection Types' },
    ],
  },
  {
    key: 'case',
    label: 'Case',
    icon: 'bi-pc-display',
    tables: [
      { table: 'case_types', label: 'Case Types' },
      { table: 'case_panel_materials', label: 'Case Panel Materials' },
      { table: 'drive_bay_types', label: 'Drive Bay Types' },
      { table: 'io_port_types', label: 'Front I/O Port Types' },
    ],
  },
  {
    key: 'cooling',
    label: 'Cooling',
    icon: 'bi-fan',
    tables: [
      { table: 'cooling_types', label: 'Cooling Types' },
      { table: 'fan_sizes', label: 'Fan Sizes' },
      { table: 'radiator_sizes', label: 'Radiator Sizes' },
      { table: 'fan_bearing_types', label: 'Fan Bearing Types' },
    ],
  },
]

export const LOOKUP_TABLES = LOOKUP_GROUPS.flatMap((g) =>
  g.tables.map((t) => ({ ...t, group: g.key, groupLabel: g.label })),
)

export function getLookupMeta(table) {
  return LOOKUP_TABLES.find((t) => t.table === table)
}
