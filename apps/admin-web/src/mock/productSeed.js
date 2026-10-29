// L(table, code) resolves a lookup code to its id.
// Prices are VND; `stock` seeds the inventory record.
export function buildProductSeed(L) {
  const ids = (table, ...codes) => codes.map((c) => L(table, c))

  return [
    // ---------------- CPU ----------------
    {
      type: 'CPU', manufacturer: 'INTEL', model_name: 'Core i5-14600K', release_date: '2023-10-17',
      description: '14 cores (6P + 8E), unlocked desktop processor.',
      sku: 'CPU-INT-14600K', price: 7_490_000, cost: 6_600_000, stock: 24,
      spec: {
        architecture_id: L('cpu_architectures', 'RAPTOR_LAKE_R'), cpu_series_id: L('cpu_series', 'CORE_I5_14'),
        socket_id: L('cpu_sockets', 'LGA1700'), core_count: 14, thread_count: 20, base_clock_mhz: 3500,
        max_boost_clock_mhz: 5300, l1_cache_kb: 1200, l2_cache_mb: 20, l3_cache_mb: 24, tdp_w: 125,
        has_integrated_graphics: true, integrated_gpu_chipset_id: L('gpu_chipsets', 'UHD_770'),
        max_memory_capacity_gb: 192, memory_channel_count: 2, pcie_generation_id: L('pcie_generations', 'PCIE5'),
        unlocked_multiplier: true, cooler_included: false, max_operating_temperature_c: 100,
        supported_memory_type_ids: ids('memory_types', 'DDR4', 'DDR5'),
        supported_instruction_set_ids: ids('cpu_instruction_sets', 'SSE4_2', 'AVX2', 'AES_NI', 'FMA3'),
      },
    },
    {
      type: 'CPU', manufacturer: 'INTEL', model_name: 'Core Ultra 7 265K', release_date: '2024-10-24',
      description: '20 cores (8P + 12E) Arrow Lake desktop processor.',
      sku: 'CPU-INT-U7-265K', price: 9_990_000, cost: 8_900_000, stock: 12,
      spec: {
        architecture_id: L('cpu_architectures', 'ARROW_LAKE'), cpu_series_id: L('cpu_series', 'CORE_ULTRA_200S'),
        socket_id: L('cpu_sockets', 'LGA1851'), core_count: 20, thread_count: 20, base_clock_mhz: 3900,
        max_boost_clock_mhz: 5500, l2_cache_mb: 36, l3_cache_mb: 30, tdp_w: 125,
        has_integrated_graphics: true, integrated_gpu_chipset_id: L('gpu_chipsets', 'INTEL_XE_LPG'),
        max_memory_capacity_gb: 192, memory_channel_count: 2, pcie_generation_id: L('pcie_generations', 'PCIE5'),
        unlocked_multiplier: true, cooler_included: false, max_operating_temperature_c: 105,
        supported_memory_type_ids: ids('memory_types', 'DDR5'),
        supported_instruction_set_ids: ids('cpu_instruction_sets', 'SSE4_2', 'AVX2', 'AES_NI', 'FMA3'),
      },
    },
    {
      type: 'CPU', manufacturer: 'AMD', model_name: 'Ryzen 7 9800X3D', release_date: '2024-11-07',
      description: '8-core gaming processor with 2nd gen 3D V-Cache.',
      sku: 'CPU-AMD-9800X3D', price: 12_990_000, cost: 11_500_000, stock: 4,
      spec: {
        architecture_id: L('cpu_architectures', 'ZEN5'), cpu_series_id: L('cpu_series', 'RYZEN_9000'),
        socket_id: L('cpu_sockets', 'AM5'), core_count: 8, thread_count: 16, base_clock_mhz: 4700,
        max_boost_clock_mhz: 5200, l1_cache_kb: 640, l2_cache_mb: 8, l3_cache_mb: 96, tdp_w: 120,
        has_integrated_graphics: true, integrated_gpu_chipset_id: L('gpu_chipsets', 'RADEON_IGPU'),
        max_memory_capacity_gb: 192, memory_channel_count: 2, pcie_generation_id: L('pcie_generations', 'PCIE5'),
        unlocked_multiplier: true, cooler_included: false, max_operating_temperature_c: 95,
        supported_memory_type_ids: ids('memory_types', 'DDR5'),
        supported_instruction_set_ids: ids('cpu_instruction_sets', 'SSE4_2', 'AVX2', 'AVX512', 'AES_NI', 'FMA3'),
      },
    },

    // ---------------- GPU ----------------
    {
      type: 'GPU', manufacturer: 'ASUS', model_name: 'TUF Gaming GeForce RTX 4070 SUPER OC 12GB', release_date: '2024-01-17',
      description: 'Triple-fan RTX 4070 SUPER with military-grade components.',
      sku: 'GPU-ASUS-TUF-4070S', price: 17_990_000, cost: 16_200_000, stock: 8,
      spec: {
        chipset_id: L('gpu_chipsets', 'RTX_4070_SUPER'), architecture_id: L('gpu_architectures', 'ADA'),
        vram_capacity_gb: 12, vram_type_id: L('vram_types', 'GDDR6X'), memory_bus_width_bits: 192,
        memory_bandwidth_gbps: 504, compute_unit_count: 56, base_clock_mhz: 1980, boost_clock_mhz: 2610,
        pcie_generation_id: L('pcie_generations', 'PCIE4'), pcie_lane_count: 16, recommended_psu_w: 750,
        total_graphics_power_w: 220, max_resolution_width_px: 7680, max_resolution_height_px: 4320,
        max_displays_supported: 4, length_mm: 301, thickness_slots: 3.25, cooling_type_id: L('cooling_types', 'AIR'),
        graphics_api_version_id: L('graphics_api_versions', 'DX12U'), ray_tracing_supported: true,
        power_connectors: [{ power_connector_type_id: L('power_connector_types', '12VHPWR'), connector_count: 1 }],
        display_outputs: [
          { display_connector_type_id: L('display_connector_types', 'HDMI_2_1'), port_count: 2 },
          { display_connector_type_id: L('display_connector_types', 'DP_1_4A'), port_count: 3 },
        ],
        upscaling_technology_ids: ids('upscaling_technologies', 'DLSS3'),
      },
    },
    {
      type: 'GPU', manufacturer: 'MSI', model_name: 'GeForce RTX 5080 16G GAMING TRIO OC', release_date: '2025-01-30',
      description: 'Blackwell RTX 5080 with TRI FROZR 3 cooling.',
      sku: 'GPU-MSI-5080-TRIO', price: 34_990_000, cost: 32_000_000, stock: 3,
      spec: {
        chipset_id: L('gpu_chipsets', 'RTX_5080'), architecture_id: L('gpu_architectures', 'BLACKWELL'),
        vram_capacity_gb: 16, vram_type_id: L('vram_types', 'GDDR7'), memory_bus_width_bits: 256,
        memory_bandwidth_gbps: 960, compute_unit_count: 84, base_clock_mhz: 2295, boost_clock_mhz: 2730,
        pcie_generation_id: L('pcie_generations', 'PCIE5'), pcie_lane_count: 16, recommended_psu_w: 850,
        total_graphics_power_w: 360, max_resolution_width_px: 7680, max_resolution_height_px: 4320,
        max_displays_supported: 4, length_mm: 338, thickness_slots: 3, cooling_type_id: L('cooling_types', 'AIR'),
        graphics_api_version_id: L('graphics_api_versions', 'DX12U'), ray_tracing_supported: true,
        power_connectors: [{ power_connector_type_id: L('power_connector_types', '12V_2X6'), connector_count: 1 }],
        display_outputs: [
          { display_connector_type_id: L('display_connector_types', 'HDMI_2_1'), port_count: 1 },
          { display_connector_type_id: L('display_connector_types', 'DP_2_1'), port_count: 3 },
        ],
        upscaling_technology_ids: ids('upscaling_technologies', 'DLSS3', 'DLSS4'),
      },
    },
    {
      type: 'GPU', manufacturer: 'GIGABYTE', model_name: 'Radeon RX 9070 XT GAMING OC 16G', release_date: '2025-03-06',
      description: 'RDNA 4 graphics card with WINDFORCE cooling.',
      sku: 'GPU-GB-9070XT-GOC', price: 19_490_000, cost: 17_600_000, stock: 0,
      spec: {
        chipset_id: L('gpu_chipsets', 'RX_9070_XT'), architecture_id: L('gpu_architectures', 'RDNA4'),
        vram_capacity_gb: 16, vram_type_id: L('vram_types', 'GDDR6'), memory_bus_width_bits: 256,
        memory_bandwidth_gbps: 640, compute_unit_count: 64, base_clock_mhz: 2400, boost_clock_mhz: 3060,
        pcie_generation_id: L('pcie_generations', 'PCIE5'), pcie_lane_count: 16, recommended_psu_w: 850,
        total_graphics_power_w: 340, max_displays_supported: 4, length_mm: 288, thickness_slots: 3,
        cooling_type_id: L('cooling_types', 'AIR'), graphics_api_version_id: L('graphics_api_versions', 'DX12U'),
        ray_tracing_supported: true,
        power_connectors: [{ power_connector_type_id: L('power_connector_types', 'PCIE_8PIN'), connector_count: 3 }],
        display_outputs: [
          { display_connector_type_id: L('display_connector_types', 'HDMI_2_1'), port_count: 2 },
          { display_connector_type_id: L('display_connector_types', 'DP_2_1'), port_count: 2 },
        ],
        upscaling_technology_ids: ids('upscaling_technologies', 'FSR3', 'FSR4'),
      },
    },

    // ---------------- MOTHERBOARD ----------------
    {
      type: 'MOTHERBOARD', manufacturer: 'ASUS', model_name: 'ROG STRIX Z790-E GAMING WIFI II', release_date: '2023-10-17',
      description: 'High-end Z790 ATX board with Wi-Fi 7.',
      sku: 'MB-ASUS-Z790E-II', price: 13_490_000, cost: 12_000_000, stock: 6,
      spec: {
        chipset_id: L('motherboard_chipsets', 'Z790'), socket_id: L('cpu_sockets', 'LGA1700'),
        form_factor_id: L('motherboard_form_factors', 'ATX'), memory_slot_count: 4, max_memory_capacity_gb: 192,
        max_memory_data_rate_mt_s: 8000, memory_channel_count: 2, sata_port_count: 4, ethernet_speed_mbps: 2500,
        audio_codec_id: L('audio_codecs', 'ALC4080'), vrm_phase_description: '20+1+2 (90A)',
        bios_flashback_supported: true, rgb_header_count: 4,
        supported_cpu_series_ids: ids('cpu_series', 'CORE_I5_14', 'CORE_I7_14'),
        supported_memory_type_ids: ids('memory_types', 'DDR5'),
        expansion_slots: [
          { pcie_generation_id: L('pcie_generations', 'PCIE5'), lane_count: 16, slot_count: 1, description: 'CPU' },
          { pcie_generation_id: L('pcie_generations', 'PCIE4'), lane_count: 4, slot_count: 1, description: 'Chipset' },
        ],
        storage_slots: [
          { storage_interface_id: L('storage_interfaces', 'M2_PCIE'), pcie_generation_id: L('pcie_generations', 'PCIE5'), slot_count: 1, form_factor_id: L('storage_form_factors', 'M2_2280'), description: 'M.2_1 (CPU)' },
          { storage_interface_id: L('storage_interfaces', 'M2_PCIE'), pcie_generation_id: L('pcie_generations', 'PCIE4'), slot_count: 4, form_factor_id: L('storage_form_factors', 'M2_2280'), description: null },
        ],
        usb_ports: [
          { usb_port_type_id: L('usb_port_types', 'USB32G2_A'), port_count: 8 },
          { usb_port_type_id: L('usb_port_types', 'USB4_C'), port_count: 2 },
        ],
        display_outputs: [{ display_connector_type_id: L('display_connector_types', 'HDMI_2_1'), port_count: 1 }],
        wireless_standard_ids: ids('wireless_standards', 'WIFI7'),
        bluetooth_version_ids: ids('bluetooth_versions', 'BT_5_4'),
      },
    },
    {
      type: 'MOTHERBOARD', manufacturer: 'MSI', model_name: 'MAG B650 TOMAHAWK WIFI', release_date: '2022-10-10',
      description: 'Mainstream AM5 ATX board.',
      sku: 'MB-MSI-B650-TMHK', price: 5_490_000, cost: 4_800_000, stock: 15,
      spec: {
        chipset_id: L('motherboard_chipsets', 'B650'), socket_id: L('cpu_sockets', 'AM5'),
        form_factor_id: L('motherboard_form_factors', 'ATX'), memory_slot_count: 4, max_memory_capacity_gb: 192,
        max_memory_data_rate_mt_s: 6600, memory_channel_count: 2, sata_port_count: 6, ethernet_speed_mbps: 2500,
        audio_codec_id: L('audio_codecs', 'ALC4080'), vrm_phase_description: '14+2+1 (80A)',
        bios_flashback_supported: true, rgb_header_count: 4,
        supported_cpu_series_ids: ids('cpu_series', 'RYZEN_7000', 'RYZEN_9000'),
        supported_memory_type_ids: ids('memory_types', 'DDR5'),
        expansion_slots: [
          { pcie_generation_id: L('pcie_generations', 'PCIE4'), lane_count: 16, slot_count: 1, description: 'CPU' },
          { pcie_generation_id: L('pcie_generations', 'PCIE3'), lane_count: 1, slot_count: 2, description: 'Chipset' },
        ],
        storage_slots: [
          { storage_interface_id: L('storage_interfaces', 'M2_PCIE'), pcie_generation_id: L('pcie_generations', 'PCIE4'), slot_count: 3, form_factor_id: L('storage_form_factors', 'M2_2280'), description: null },
        ],
        usb_ports: [
          { usb_port_type_id: L('usb_port_types', 'USB2_A'), port_count: 4 },
          { usb_port_type_id: L('usb_port_types', 'USB32G2_A'), port_count: 3 },
          { usb_port_type_id: L('usb_port_types', 'USB32G2_C'), port_count: 1 },
        ],
        display_outputs: [
          { display_connector_type_id: L('display_connector_types', 'HDMI_2_1'), port_count: 1 },
          { display_connector_type_id: L('display_connector_types', 'DP_1_4A'), port_count: 1 },
        ],
        wireless_standard_ids: ids('wireless_standards', 'WIFI6E'),
        bluetooth_version_ids: ids('bluetooth_versions', 'BT_5_3'),
      },
    },
    {
      type: 'MOTHERBOARD', manufacturer: 'GIGABYTE', model_name: 'Z890 AORUS ELITE WIFI7', release_date: '2024-10-24',
      description: 'Z890 ATX board for Intel Core Ultra 200S.',
      sku: 'MB-GB-Z890-ELITE', price: 8_290_000, cost: 7_400_000, stock: 2,
      spec: {
        chipset_id: L('motherboard_chipsets', 'Z890'), socket_id: L('cpu_sockets', 'LGA1851'),
        form_factor_id: L('motherboard_form_factors', 'ATX'), memory_slot_count: 4, max_memory_capacity_gb: 256,
        max_memory_data_rate_mt_s: 9000, memory_channel_count: 2, sata_port_count: 4, ethernet_speed_mbps: 2500,
        audio_codec_id: L('audio_codecs', 'ALC1220'), vrm_phase_description: '16+1+2 (80A)',
        bios_flashback_supported: true, rgb_header_count: 3,
        supported_cpu_series_ids: ids('cpu_series', 'CORE_ULTRA_200S'),
        supported_memory_type_ids: ids('memory_types', 'DDR5'),
        storage_slots: [
          { storage_interface_id: L('storage_interfaces', 'M2_PCIE'), pcie_generation_id: L('pcie_generations', 'PCIE5'), slot_count: 1, form_factor_id: L('storage_form_factors', 'M2_2280'), description: null },
          { storage_interface_id: L('storage_interfaces', 'M2_PCIE'), pcie_generation_id: L('pcie_generations', 'PCIE4'), slot_count: 3, form_factor_id: L('storage_form_factors', 'M2_2280'), description: null },
        ],
        wireless_standard_ids: ids('wireless_standards', 'WIFI7'),
        bluetooth_version_ids: ids('bluetooth_versions', 'BT_5_4'),
      },
    },

    // ---------------- RAM ----------------
    {
      type: 'RAM', manufacturer: 'CORSAIR', model_name: 'Vengeance RGB DDR5 32GB (2x16GB) 6000MHz CL30', release_date: '2023-03-01',
      description: 'DDR5 kit tuned for AMD EXPO and Intel XMP 3.0.',
      sku: 'RAM-COR-VEN-32-6000', price: 3_290_000, cost: 2_800_000, stock: 40,
      spec: {
        memory_type_id: L('memory_types', 'DDR5'), module_type_id: L('memory_module_types', 'UDIMM'),
        total_capacity_gb: 32, module_count: 2, capacity_per_module_gb: 16, data_rate_mt_s: 6000,
        cas_latency_cl: 30, timing_t_r_c_d: 36, timing_t_r_p: 36, timing_t_r_a_s: 76, operating_voltage_v: 1.4,
        ecc_supported: false, registered_buffered: false, on_die_ecc: true, xmp_supported: true, expo_supported: true,
        heat_spreader_type_id: L('ram_heat_spreader_types', 'ALUMINUM'), rgb_supported: true,
        color_id: L('product_colors', 'BLACK'), module_height_mm: 44,
        memory_profiles: [
          { profile_type_id: L('ram_profile_types', 'JEDEC'), profile_name: 'JEDEC', data_rate_mt_s: 4800, voltage_v: 1.1, cas_latency_cl: 40, timing_t_r_c_d: 40, timing_t_r_p: 40, timing_t_r_a_s: 77 },
          { profile_type_id: L('ram_profile_types', 'XMP3'), profile_name: 'XMP Profile 1', data_rate_mt_s: 6000, voltage_v: 1.4, cas_latency_cl: 30, timing_t_r_c_d: 36, timing_t_r_p: 36, timing_t_r_a_s: 76 },
          { profile_type_id: L('ram_profile_types', 'EXPO'), profile_name: 'EXPO I', data_rate_mt_s: 6000, voltage_v: 1.4, cas_latency_cl: 30, timing_t_r_c_d: 36, timing_t_r_p: 36, timing_t_r_a_s: 76 },
        ],
      },
    },
    {
      type: 'RAM', manufacturer: 'KINGSTON', model_name: 'FURY Beast DDR4 16GB (2x8GB) 3200MHz CL16', release_date: '2021-06-01',
      description: 'Budget DDR4 kit.',
      sku: 'RAM-KST-FB-16-3200', price: 990_000, cost: 820_000, stock: 55,
      spec: {
        memory_type_id: L('memory_types', 'DDR4'), module_type_id: L('memory_module_types', 'UDIMM'),
        total_capacity_gb: 16, module_count: 2, capacity_per_module_gb: 8, data_rate_mt_s: 3200,
        cas_latency_cl: 16, timing_t_r_c_d: 18, timing_t_r_p: 18, operating_voltage_v: 1.35,
        ecc_supported: false, registered_buffered: false, on_die_ecc: false, xmp_supported: true, expo_supported: false,
        heat_spreader_type_id: L('ram_heat_spreader_types', 'LOW_PROFILE'), rgb_supported: false,
        color_id: L('product_colors', 'BLACK'), module_height_mm: 34.1,
        memory_profiles: [
          { profile_type_id: L('ram_profile_types', 'XMP3'), profile_name: 'XMP', data_rate_mt_s: 3200, voltage_v: 1.35, cas_latency_cl: 16, timing_t_r_c_d: 18, timing_t_r_p: 18, timing_t_r_a_s: null },
        ],
      },
    },
    {
      type: 'RAM', manufacturer: 'GSKILL', model_name: 'Trident Z5 RGB DDR5 64GB (2x32GB) 6400MHz CL32', release_date: '2023-05-15',
      description: 'High-capacity performance DDR5 kit.',
      sku: 'RAM-GSK-TZ5-64-6400', price: 6_490_000, cost: 5_700_000, stock: 9,
      spec: {
        memory_type_id: L('memory_types', 'DDR5'), module_type_id: L('memory_module_types', 'UDIMM'),
        total_capacity_gb: 64, module_count: 2, capacity_per_module_gb: 32, data_rate_mt_s: 6400,
        cas_latency_cl: 32, timing_t_r_c_d: 39, timing_t_r_p: 39, timing_t_r_a_s: 102, operating_voltage_v: 1.4,
        on_die_ecc: true, xmp_supported: true, expo_supported: false,
        heat_spreader_type_id: L('ram_heat_spreader_types', 'ALUMINUM'), rgb_supported: true,
        color_id: L('product_colors', 'SILVER'), module_height_mm: 44,
      },
    },

    // ---------------- STORAGE ----------------
    {
      type: 'STORAGE', manufacturer: 'SAMSUNG', model_name: '990 PRO 2TB', release_date: '2022-10-18',
      description: 'PCIe 4.0 NVMe flagship SSD.',
      sku: 'SSD-SAM-990PRO-2T', price: 4_590_000, cost: 3_950_000, stock: 30,
      spec: {
        device_type_id: L('storage_device_types', 'NVME_SSD'), capacity_gb: 2000,
        form_factor_id: L('storage_form_factors', 'M2_2280'), interface_id: L('storage_interfaces', 'M2_PCIE'),
        protocol_id: L('storage_protocols', 'NVME_2_0'), pcie_generation_id: L('pcie_generations', 'PCIE4'),
        sequential_read_mbps: 7450, sequential_write_mbps: 6900, random_read_iops: 1400000, random_write_iops: 1550000,
        nand_type_id: L('nand_types', 'TLC'), dram_cache_present: true, endurance_tbw: 1200, mtbf_hours: 1500000,
        warranty_years: 5, power_loss_protection: false,
        encryption_standard_ids: ids('encryption_standards', 'AES_256', 'TCG_OPAL'),
      },
    },
    {
      type: 'STORAGE', manufacturer: 'WD', model_name: 'Blue SA510 1TB 2.5"', release_date: '2022-03-01',
      description: 'SATA SSD for upgrades.',
      sku: 'SSD-WD-SA510-1T', price: 1_590_000, cost: 1_300_000, stock: 18,
      spec: {
        device_type_id: L('storage_device_types', 'SATA_SSD'), capacity_gb: 1000,
        form_factor_id: L('storage_form_factors', 'SFF_2_5'), interface_id: L('storage_interfaces', 'SATA3'),
        protocol_id: L('storage_protocols', 'AHCI'), sequential_read_mbps: 560, sequential_write_mbps: 520,
        random_read_iops: 90000, random_write_iops: 82000, nand_type_id: L('nand_types', 'TLC'),
        dram_cache_present: false, endurance_tbw: 400, mtbf_hours: 1750000, warranty_years: 3,
      },
    },
    {
      type: 'STORAGE', manufacturer: 'SEAGATE', model_name: 'BarraCuda 2TB 7200RPM', release_date: '2019-01-01',
      description: '3.5" desktop hard drive.',
      sku: 'HDD-SEA-BC-2T', price: 1_490_000, cost: 1_200_000, stock: 3,
      spec: {
        device_type_id: L('storage_device_types', 'HDD'), capacity_gb: 2000,
        form_factor_id: L('storage_form_factors', 'LFF_3_5'), interface_id: L('storage_interfaces', 'SATA3'),
        sequential_read_mbps: 220, sequential_write_mbps: 220, rpm: 7200, cache_mb: 256, warranty_years: 2,
      },
    },

    // ---------------- PSU ----------------
    {
      type: 'PSU', manufacturer: 'CORSAIR', model_name: 'RM850x (2024) 850W 80+ Gold', release_date: '2024-02-01',
      description: 'ATX 3.1 fully modular PSU.',
      sku: 'PSU-COR-RM850X-24', price: 3_690_000, cost: 3_150_000, stock: 14,
      spec: {
        rated_power_w: 850, efficiency_certification_id: L('psu_efficiency_certifications', '80P_GOLD'),
        psu_standard_id: L('psu_standards', 'ATX_3_1'), form_factor_id: L('psu_form_factors', 'ATX'),
        modular_type_id: L('psu_modular_types', 'FULL'), input_voltage_min_v: 100, input_voltage_max_v: 240,
        frequency_min_hz: 47, frequency_max_hz: 63, twelve_v_output_w: 850, sata_power_connector_count: 14,
        molex_connector_count: 4, fan_size_id: L('fan_sizes', 'FAN_140'), zero_rpm_mode: true,
        length_mm: 160, width_mm: 150, height_mm: 86, warranty_years: 10,
        power_connectors: [
          { power_connector_type_id: L('power_connector_types', 'ATX_24PIN'), connector_count: 1 },
          { power_connector_type_id: L('power_connector_types', 'EPS_8PIN'), connector_count: 2 },
          { power_connector_type_id: L('power_connector_types', 'PCIE_8PIN'), connector_count: 4 },
          { power_connector_type_id: L('power_connector_types', '12V_2X6'), connector_count: 1 },
        ],
        protection_type_ids: ids('psu_protection_types', 'OVP', 'UVP', 'OCP', 'OPP', 'OTP', 'SCP'),
      },
    },
    {
      type: 'PSU', manufacturer: 'SEASONIC', model_name: 'FOCUS GX-750 ATX 3.0', release_date: '2023-04-01',
      description: '750W fully modular, 80+ Gold.',
      sku: 'PSU-SSN-GX750-A3', price: 2_990_000, cost: 2_550_000, stock: 11,
      spec: {
        rated_power_w: 750, efficiency_certification_id: L('psu_efficiency_certifications', '80P_GOLD'),
        psu_standard_id: L('psu_standards', 'ATX_3_0'), form_factor_id: L('psu_form_factors', 'ATX'),
        modular_type_id: L('psu_modular_types', 'FULL'), twelve_v_output_w: 744,
        sata_power_connector_count: 8, molex_connector_count: 3, fan_size_id: L('fan_sizes', 'FAN_120'),
        zero_rpm_mode: true, length_mm: 140, width_mm: 150, height_mm: 86, warranty_years: 10,
        power_connectors: [
          { power_connector_type_id: L('power_connector_types', 'ATX_24PIN'), connector_count: 1 },
          { power_connector_type_id: L('power_connector_types', 'EPS_8PIN'), connector_count: 2 },
          { power_connector_type_id: L('power_connector_types', '12VHPWR'), connector_count: 1 },
        ],
        protection_type_ids: ids('psu_protection_types', 'OVP', 'UVP', 'OCP', 'OPP', 'OTP', 'SCP'),
      },
    },
    {
      type: 'PSU', manufacturer: 'MSI', model_name: 'MAG A650BN 650W 80+ Bronze', release_date: '2020-08-01',
      description: 'Entry-level non-modular PSU.',
      sku: 'PSU-MSI-A650BN', price: 1_190_000, cost: 980_000, stock: 26,
      spec: {
        rated_power_w: 650, efficiency_certification_id: L('psu_efficiency_certifications', '80P_BRONZE'),
        psu_standard_id: L('psu_standards', 'ATX_2'), form_factor_id: L('psu_form_factors', 'ATX'),
        modular_type_id: L('psu_modular_types', 'NON'), sata_power_connector_count: 5, molex_connector_count: 2,
        fan_size_id: L('fan_sizes', 'FAN_120'), zero_rpm_mode: false, warranty_years: 5,
        protection_type_ids: ids('psu_protection_types', 'OVP', 'UVP', 'OCP', 'OPP', 'SCP'),
      },
    },

    // ---------------- CASE ----------------
    {
      type: 'CASE', manufacturer: 'LIAN_LI', model_name: 'O11 Dynamic EVO', release_date: '2022-03-01',
      description: 'Dual-chamber mid tower with reversible layout.',
      sku: 'CASE-LL-O11-EVO-BK', price: 3_890_000, cost: 3_300_000, stock: 7,
      spec: {
        case_type_id: L('case_types', 'MID_TOWER'), color_id: L('product_colors', 'BLACK'),
        height_mm: 459, width_mm: 285, depth_mm: 465, volume_l: 60.8, max_gpu_length_mm: 426,
        max_cpu_cooler_height_mm: 167, max_psu_length_mm: 220, expansion_slot_count: 7, included_fan_count: 0,
        max_fan_count: 10, dust_filters_present: true, rgb_supported: true, vertical_gpu_mount_supported: true,
        side_panel_material_id: L('case_panel_materials', 'TEMPERED_GLASS'),
        front_panel_material_id: L('case_panel_materials', 'TEMPERED_GLASS'),
        supported_motherboard_form_factor_ids: ids('motherboard_form_factors', 'EATX', 'ATX', 'MATX', 'MITX'),
        fan_support: [
          { fan_size_id: L('fan_sizes', 'FAN_120'), supported_count: 10 },
          { fan_size_id: L('fan_sizes', 'FAN_140'), supported_count: 6 },
        ],
        radiator_support: [
          { radiator_size_id: L('radiator_sizes', 'RAD_360'), supported_count: 3 },
          { radiator_size_id: L('radiator_sizes', 'RAD_280'), supported_count: 2 },
        ],
        drive_bays: [
          { drive_bay_type_id: L('drive_bay_types', 'BAY_2_5'), bay_count: 4 },
          { drive_bay_type_id: L('drive_bay_types', 'BAY_3_5'), bay_count: 2 },
        ],
        front_io_ports: [
          { io_port_type_id: L('io_port_types', 'USB_A_3_2'), port_count: 2 },
          { io_port_type_id: L('io_port_types', 'USB_C_3_2'), port_count: 1 },
          { io_port_type_id: L('io_port_types', 'AUDIO_COMBO'), port_count: 1 },
        ],
      },
    },
    {
      type: 'CASE', manufacturer: 'NZXT', model_name: 'H5 Flow (2024) White', release_date: '2024-05-01',
      description: 'Compact airflow mid tower.',
      sku: 'CASE-NZXT-H5F-WH', price: 2_290_000, cost: 1_900_000, stock: 12,
      spec: {
        case_type_id: L('case_types', 'MID_TOWER'), color_id: L('product_colors', 'WHITE'),
        height_mm: 464, width_mm: 227, depth_mm: 446, max_gpu_length_mm: 365, max_cpu_cooler_height_mm: 165,
        max_psu_length_mm: 200, expansion_slot_count: 7, included_fan_count: 2, max_fan_count: 7,
        dust_filters_present: true, rgb_supported: false, vertical_gpu_mount_supported: false,
        side_panel_material_id: L('case_panel_materials', 'TEMPERED_GLASS'),
        front_panel_material_id: L('case_panel_materials', 'MESH'),
        supported_motherboard_form_factor_ids: ids('motherboard_form_factors', 'ATX', 'MATX', 'MITX'),
        fan_support: [{ fan_size_id: L('fan_sizes', 'FAN_120'), supported_count: 7 }],
        radiator_support: [{ radiator_size_id: L('radiator_sizes', 'RAD_360'), supported_count: 1 }],
        front_io_ports: [
          { io_port_type_id: L('io_port_types', 'USB_A_3_2'), port_count: 1 },
          { io_port_type_id: L('io_port_types', 'USB_C_3_2'), port_count: 1 },
        ],
      },
    },
    {
      type: 'CASE', manufacturer: 'FRACTAL', model_name: 'North Charcoal Black', release_date: '2022-10-01',
      description: 'Mid tower with walnut front panel.',
      sku: 'CASE-FD-NORTH-BK', price: 3_590_000, cost: 3_050_000, stock: 1,
      spec: {
        case_type_id: L('case_types', 'MID_TOWER'), color_id: L('product_colors', 'BLACK'),
        max_gpu_length_mm: 355, max_cpu_cooler_height_mm: 170, included_fan_count: 2,
        side_panel_material_id: L('case_panel_materials', 'MESH'),
        front_panel_material_id: L('case_panel_materials', 'WOOD'),
        dust_filters_present: true, rgb_supported: false, vertical_gpu_mount_supported: false,
        supported_motherboard_form_factor_ids: ids('motherboard_form_factors', 'ATX', 'MATX', 'MITX'),
      },
    },

    // ---------------- COOLING ----------------
    {
      type: 'COOLING', manufacturer: 'NOCTUA', model_name: 'NH-D15 G2', release_date: '2024-06-01',
      description: 'Dual-tower premium air cooler.',
      sku: 'COOL-NOC-NHD15G2', price: 3_990_000, cost: 3_400_000, stock: 5,
      spec: {
        cooling_type_id: L('cooling_types', 'AIR'), rated_tdp_w: 250, cooler_height_mm: 168,
        fan_size_mm: 140, fan_count: 2, fan_speed_min_rpm: 300, fan_speed_max_rpm: 1500, max_airflow_cfm: 82.5,
        static_pressure_mmh2o: 2.6, noise_level_dba: 24.8, fan_bearing_type_id: L('fan_bearing_types', 'SSO2'),
        fan_connector_type_id: L('electrical_connector_types', 'FAN_4PIN_PWM'),
        lighting_type_id: L('lighting_types', 'NONE'), thermal_paste_included: true, warranty_years: 6,
        supported_cpu_socket_ids: ids('cpu_sockets', 'LGA1700', 'LGA1851', 'AM4', 'AM5'),
        fans: [
          { fan_size_id: L('fan_sizes', 'FAN_140'), fan_count: 2, speed_min_rpm: 300, speed_max_rpm: 1500, max_airflow_cfm: 82.5, static_pressure_mmh2o: 2.6, noise_level_dba: 24.8, fan_bearing_type_id: L('fan_bearing_types', 'SSO2'), fan_connector_type_id: L('electrical_connector_types', 'FAN_4PIN_PWM') },
        ],
      },
    },
    {
      type: 'COOLING', manufacturer: 'DEEPCOOL', model_name: 'AK400', release_date: '2022-05-01',
      description: 'Single-tower budget air cooler.',
      sku: 'COOL-DC-AK400', price: 690_000, cost: 520_000, stock: 33,
      spec: {
        cooling_type_id: L('cooling_types', 'AIR'), rated_tdp_w: 220, cooler_height_mm: 155,
        fan_size_mm: 120, fan_count: 1, fan_speed_min_rpm: 500, fan_speed_max_rpm: 1850, max_airflow_cfm: 66.47,
        static_pressure_mmh2o: 2.04, noise_level_dba: 29, fan_bearing_type_id: L('fan_bearing_types', 'FDB'),
        fan_connector_type_id: L('electrical_connector_types', 'FAN_4PIN_PWM'),
        lighting_type_id: L('lighting_types', 'NONE'), thermal_paste_included: true, warranty_years: 3,
        supported_cpu_socket_ids: ids('cpu_sockets', 'LGA1700', 'AM4', 'AM5'),
      },
    },
    {
      type: 'COOLING', manufacturer: 'NZXT', model_name: 'Kraken 360 RGB', release_date: '2023-03-01',
      description: '360mm AIO with LCD pump head.',
      sku: 'COOL-NZXT-KRK360', price: 4_790_000, cost: 4_100_000, stock: 0,
      spec: {
        cooling_type_id: L('cooling_types', 'AIO'), rated_tdp_w: 300, radiator_size_id: L('radiator_sizes', 'RAD_360'),
        fan_size_mm: 120, fan_count: 3, fan_speed_min_rpm: 500, fan_speed_max_rpm: 1800,
        fan_bearing_type_id: L('fan_bearing_types', 'FDB'),
        fan_connector_type_id: L('electrical_connector_types', 'FAN_4PIN_PWM'),
        pump_speed_min_rpm: 800, pump_speed_max_rpm: 2800,
        pump_connector_type_id: L('electrical_connector_types', 'USB2_HEADER'),
        lighting_type_id: L('lighting_types', 'ARGB'), thermal_paste_included: true, warranty_years: 6,
        supported_cpu_socket_ids: ids('cpu_sockets', 'LGA1700', 'LGA1851', 'AM4', 'AM5'),
        fans: [
          { fan_size_id: L('fan_sizes', 'FAN_120'), fan_count: 3, speed_min_rpm: 500, speed_max_rpm: 1800, max_airflow_cfm: null, static_pressure_mmh2o: null, noise_level_dba: null, fan_bearing_type_id: L('fan_bearing_types', 'FDB'), fan_connector_type_id: L('electrical_connector_types', 'FAN_4PIN_PWM') },
        ],
      },
    },
  ]
}
