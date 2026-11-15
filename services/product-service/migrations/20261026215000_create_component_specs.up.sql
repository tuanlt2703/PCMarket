CREATE TABLE cpu_specs (
    model_id BIGINT PRIMARY KEY REFERENCES models(id) ON DELETE CASCADE,
    architecture_id BIGINT REFERENCES cpu_architectures(id),
    cpu_series_id BIGINT REFERENCES cpu_series(id),
    socket_id BIGINT REFERENCES cpu_sockets(id),
    core_count INTEGER CHECK (core_count IS NULL OR core_count > 0),
    thread_count INTEGER CHECK (thread_count IS NULL OR thread_count > 0),
    base_clock_mhz NUMERIC(8,2) CHECK (base_clock_mhz IS NULL OR base_clock_mhz > 0),
    max_boost_clock_mhz NUMERIC(8,2) CHECK (max_boost_clock_mhz IS NULL OR max_boost_clock_mhz > 0),
    l1_cache_kb NUMERIC(10,2) CHECK (l1_cache_kb IS NULL OR l1_cache_kb >= 0),
    l2_cache_mb NUMERIC(10,2) CHECK (l2_cache_mb IS NULL OR l2_cache_mb >= 0),
    l3_cache_mb NUMERIC(10,2) CHECK (l3_cache_mb IS NULL OR l3_cache_mb >= 0),
    tdp_w NUMERIC(8,2) CHECK (tdp_w IS NULL OR tdp_w >= 0),
    has_integrated_graphics BOOLEAN,
    integrated_gpu_chipset_id BIGINT REFERENCES gpu_chipsets(id),
    max_memory_capacity_gb NUMERIC(10,2) CHECK (max_memory_capacity_gb IS NULL OR max_memory_capacity_gb > 0),
    memory_channel_count SMALLINT CHECK (memory_channel_count IS NULL OR memory_channel_count > 0),
    pcie_generation_id BIGINT REFERENCES pcie_generations(id),
    unlocked_multiplier BOOLEAN,
    cooler_included BOOLEAN,
    max_operating_temperature_c NUMERIC(5,2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cpu_supported_memory_types (
    cpu_model_id BIGINT NOT NULL REFERENCES cpu_specs(model_id) ON DELETE CASCADE,
    memory_type_id BIGINT NOT NULL REFERENCES memory_types(id),
    PRIMARY KEY (cpu_model_id, memory_type_id)
);

CREATE TABLE cpu_supported_instruction_sets (
    cpu_model_id BIGINT NOT NULL REFERENCES cpu_specs(model_id) ON DELETE CASCADE,
    instruction_set_id BIGINT NOT NULL REFERENCES cpu_instruction_sets(id),
    PRIMARY KEY (cpu_model_id, instruction_set_id)
);

CREATE TABLE gpu_specs (
    model_id BIGINT PRIMARY KEY REFERENCES models(id) ON DELETE CASCADE,
    chipset_id BIGINT REFERENCES gpu_chipsets(id),
    architecture_id BIGINT REFERENCES gpu_architectures(id),
    vram_capacity_gb NUMERIC(10,2) CHECK (vram_capacity_gb IS NULL OR vram_capacity_gb > 0),
    vram_type_id BIGINT REFERENCES vram_types(id),
    memory_bus_width_bits INTEGER CHECK (memory_bus_width_bits IS NULL OR memory_bus_width_bits > 0),
    memory_bandwidth_gbps NUMERIC(10,2) CHECK (memory_bandwidth_gbps IS NULL OR memory_bandwidth_gbps > 0),
    compute_unit_count INTEGER CHECK (compute_unit_count IS NULL OR compute_unit_count > 0),
    base_clock_mhz NUMERIC(8,2) CHECK (base_clock_mhz IS NULL OR base_clock_mhz > 0),
    boost_clock_mhz NUMERIC(8,2) CHECK (boost_clock_mhz IS NULL OR boost_clock_mhz > 0),
    pcie_generation_id BIGINT REFERENCES pcie_generations(id),
    pcie_lane_count SMALLINT CHECK (pcie_lane_count IS NULL OR pcie_lane_count > 0),
    recommended_psu_w INTEGER CHECK (recommended_psu_w IS NULL OR recommended_psu_w > 0),
    total_graphics_power_w NUMERIC(8,2) CHECK (total_graphics_power_w IS NULL OR total_graphics_power_w >= 0),
    max_resolution_width_px INTEGER CHECK (max_resolution_width_px IS NULL OR max_resolution_width_px > 0),
    max_resolution_height_px INTEGER CHECK (max_resolution_height_px IS NULL OR max_resolution_height_px > 0),
    max_displays_supported SMALLINT CHECK (max_displays_supported IS NULL OR max_displays_supported > 0),
    length_mm NUMERIC(8,2) CHECK (length_mm IS NULL OR length_mm > 0),
    thickness_slots NUMERIC(4,2) CHECK (thickness_slots IS NULL OR thickness_slots > 0),
    cooling_type_id BIGINT REFERENCES cooling_types(id),
    graphics_api_version_id BIGINT REFERENCES graphics_api_versions(id),
    ray_tracing_supported BOOLEAN,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE gpu_power_connectors (
    gpu_model_id BIGINT NOT NULL REFERENCES gpu_specs(model_id) ON DELETE CASCADE,
    power_connector_type_id BIGINT NOT NULL REFERENCES power_connector_types(id),
    connector_count SMALLINT NOT NULL DEFAULT 1 CHECK (connector_count > 0),
    PRIMARY KEY (gpu_model_id, power_connector_type_id)
);

CREATE TABLE gpu_display_outputs (
    gpu_model_id BIGINT NOT NULL REFERENCES gpu_specs(model_id) ON DELETE CASCADE,
    display_connector_type_id BIGINT NOT NULL REFERENCES display_connector_types(id),
    port_count SMALLINT NOT NULL DEFAULT 1 CHECK (port_count > 0),
    PRIMARY KEY (gpu_model_id, display_connector_type_id)
);

CREATE TABLE gpu_upscaling_technologies (
    gpu_model_id BIGINT NOT NULL REFERENCES gpu_specs(model_id) ON DELETE CASCADE,
    upscaling_technology_id BIGINT NOT NULL REFERENCES upscaling_technologies(id),
    PRIMARY KEY (gpu_model_id, upscaling_technology_id)
);

CREATE TABLE motherboard_specs (
    model_id BIGINT PRIMARY KEY REFERENCES models(id) ON DELETE CASCADE,
    chipset_id BIGINT REFERENCES motherboard_chipsets(id),
    socket_id BIGINT REFERENCES cpu_sockets(id),
    form_factor_id BIGINT REFERENCES motherboard_form_factors(id),
    memory_slot_count SMALLINT CHECK (memory_slot_count IS NULL OR memory_slot_count >= 0),
    max_memory_capacity_gb NUMERIC(10,2) CHECK (max_memory_capacity_gb IS NULL OR max_memory_capacity_gb > 0),
    max_memory_data_rate_mt_s INTEGER CHECK (max_memory_data_rate_mt_s IS NULL OR max_memory_data_rate_mt_s > 0),
    memory_channel_count SMALLINT CHECK (memory_channel_count IS NULL OR memory_channel_count > 0),
    sata_port_count SMALLINT CHECK (sata_port_count IS NULL OR sata_port_count >= 0),
    ethernet_speed_mbps INTEGER CHECK (ethernet_speed_mbps IS NULL OR ethernet_speed_mbps > 0),
    audio_codec_id BIGINT REFERENCES audio_codecs(id),
    vrm_phase_description VARCHAR(100),
    bios_flashback_supported BOOLEAN,
    rgb_header_count SMALLINT CHECK (rgb_header_count IS NULL OR rgb_header_count >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE motherboard_supported_cpu_series (
    motherboard_model_id BIGINT NOT NULL REFERENCES motherboard_specs(model_id) ON DELETE CASCADE,
    cpu_series_id BIGINT NOT NULL REFERENCES cpu_series(id),
    PRIMARY KEY (motherboard_model_id, cpu_series_id)
);

CREATE TABLE motherboard_supported_memory_types (
    motherboard_model_id BIGINT NOT NULL REFERENCES motherboard_specs(model_id) ON DELETE CASCADE,
    memory_type_id BIGINT NOT NULL REFERENCES memory_types(id),
    PRIMARY KEY (motherboard_model_id, memory_type_id)
);

CREATE TABLE motherboard_expansion_slots (
    id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    motherboard_model_id BIGINT NOT NULL REFERENCES motherboard_specs(model_id) ON DELETE CASCADE,
    pcie_generation_id BIGINT REFERENCES pcie_generations(id),
    lane_count SMALLINT CHECK (lane_count IS NULL OR lane_count > 0),
    slot_count SMALLINT NOT NULL DEFAULT 1 CHECK (slot_count > 0),
    description VARCHAR(200)
);

CREATE TABLE motherboard_storage_slots (
    id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    motherboard_model_id BIGINT NOT NULL REFERENCES motherboard_specs(model_id) ON DELETE CASCADE,
    storage_interface_id BIGINT NOT NULL REFERENCES storage_interfaces(id),
    pcie_generation_id BIGINT REFERENCES pcie_generations(id),
    slot_count SMALLINT NOT NULL DEFAULT 1 CHECK (slot_count > 0),
    form_factor_id BIGINT REFERENCES storage_form_factors(id),
    description VARCHAR(200)
);

CREATE TABLE motherboard_usb_ports (
    motherboard_model_id BIGINT NOT NULL REFERENCES motherboard_specs(model_id) ON DELETE CASCADE,
    usb_port_type_id BIGINT NOT NULL REFERENCES usb_port_types(id),
    port_count SMALLINT NOT NULL DEFAULT 1 CHECK (port_count > 0),
    PRIMARY KEY (motherboard_model_id, usb_port_type_id)
);

CREATE TABLE motherboard_display_outputs (
    motherboard_model_id BIGINT NOT NULL REFERENCES motherboard_specs(model_id) ON DELETE CASCADE,
    display_connector_type_id BIGINT NOT NULL REFERENCES display_connector_types(id),
    port_count SMALLINT NOT NULL DEFAULT 1 CHECK (port_count > 0),
    PRIMARY KEY (motherboard_model_id, display_connector_type_id)
);

CREATE TABLE motherboard_wireless_standards (
    motherboard_model_id BIGINT NOT NULL REFERENCES motherboard_specs(model_id) ON DELETE CASCADE,
    wireless_standard_id BIGINT NOT NULL REFERENCES wireless_standards(id),
    PRIMARY KEY (motherboard_model_id, wireless_standard_id)
);

CREATE TABLE motherboard_bluetooth_versions (
    motherboard_model_id BIGINT NOT NULL REFERENCES motherboard_specs(model_id) ON DELETE CASCADE,
    bluetooth_version_id BIGINT NOT NULL REFERENCES bluetooth_versions(id),
    PRIMARY KEY (motherboard_model_id, bluetooth_version_id)
);

CREATE TABLE ram_specs (
    model_id BIGINT PRIMARY KEY REFERENCES models(id) ON DELETE CASCADE,
    memory_type_id BIGINT NOT NULL REFERENCES memory_types(id),
    module_type_id BIGINT REFERENCES memory_module_types(id),
    total_capacity_gb NUMERIC(10,2) CHECK (total_capacity_gb IS NULL OR total_capacity_gb > 0),
    module_count SMALLINT CHECK (module_count IS NULL OR module_count > 0),
    capacity_per_module_gb NUMERIC(10,2) CHECK (capacity_per_module_gb IS NULL OR capacity_per_module_gb > 0),
    data_rate_mt_s INTEGER CHECK (data_rate_mt_s IS NULL OR data_rate_mt_s > 0),
    cas_latency_cl SMALLINT CHECK (cas_latency_cl IS NULL OR cas_latency_cl > 0),
    timing_t_r_c_d SMALLINT CHECK (timing_t_r_c_d IS NULL OR timing_t_r_c_d > 0),
    timing_t_r_p SMALLINT CHECK (timing_t_r_p IS NULL OR timing_t_r_p > 0),
    timing_t_r_a_s SMALLINT CHECK (timing_t_r_a_s IS NULL OR timing_t_r_a_s > 0),
    operating_voltage_v NUMERIC(5,3) CHECK (operating_voltage_v IS NULL OR operating_voltage_v > 0),
    ecc_supported BOOLEAN,
    registered_buffered BOOLEAN,
    on_die_ecc BOOLEAN,
    xmp_supported BOOLEAN,
    expo_supported BOOLEAN,
    heat_spreader_type_id BIGINT REFERENCES ram_heat_spreader_types(id),
    rgb_supported BOOLEAN,
    color_id BIGINT REFERENCES product_colors(id),
    module_height_mm NUMERIC(7,2) CHECK (module_height_mm IS NULL OR module_height_mm > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ram_memory_profiles (
    id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    ram_model_id BIGINT NOT NULL REFERENCES ram_specs(model_id) ON DELETE CASCADE,
    profile_type_id BIGINT NOT NULL REFERENCES ram_profile_types(id),
    profile_name VARCHAR(100),
    data_rate_mt_s INTEGER NOT NULL CHECK (data_rate_mt_s > 0),
    voltage_v NUMERIC(5,3) CHECK (voltage_v IS NULL OR voltage_v > 0),
    cas_latency_cl SMALLINT CHECK (cas_latency_cl IS NULL OR cas_latency_cl > 0),
    timing_t_r_c_d SMALLINT,
    timing_t_r_p SMALLINT,
    timing_t_r_a_s SMALLINT
);

CREATE TABLE storage_specs (
    model_id BIGINT PRIMARY KEY REFERENCES models(id) ON DELETE CASCADE,
    device_type_id BIGINT NOT NULL REFERENCES storage_device_types(id),
    capacity_gb NUMERIC(12,2) CHECK (capacity_gb IS NULL OR capacity_gb > 0),
    form_factor_id BIGINT REFERENCES storage_form_factors(id),
    interface_id BIGINT REFERENCES storage_interfaces(id),
    protocol_id BIGINT REFERENCES storage_protocols(id),
    pcie_generation_id BIGINT REFERENCES pcie_generations(id),
    sequential_read_mbps INTEGER CHECK (sequential_read_mbps IS NULL OR sequential_read_mbps > 0),
    sequential_write_mbps INTEGER CHECK (sequential_write_mbps IS NULL OR sequential_write_mbps > 0),
    random_read_iops BIGINT CHECK (random_read_iops IS NULL OR random_read_iops > 0),
    random_write_iops BIGINT CHECK (random_write_iops IS NULL OR random_write_iops > 0),
    nand_type_id BIGINT REFERENCES nand_types(id),
    dram_cache_present BOOLEAN,
    endurance_tbw NUMERIC(12,2) CHECK (endurance_tbw IS NULL OR endurance_tbw >= 0),
    mtbf_hours BIGINT CHECK (mtbf_hours IS NULL OR mtbf_hours > 0),
    warranty_years NUMERIC(4,1) CHECK (warranty_years IS NULL OR warranty_years >= 0),
    rpm INTEGER CHECK (rpm IS NULL OR rpm > 0),
    cache_mb INTEGER CHECK (cache_mb IS NULL OR cache_mb >= 0),
    power_loss_protection BOOLEAN,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE storage_encryption_standards (
    storage_model_id BIGINT NOT NULL REFERENCES storage_specs(model_id) ON DELETE CASCADE,
    encryption_standard_id BIGINT NOT NULL REFERENCES encryption_standards(id),
    PRIMARY KEY (storage_model_id, encryption_standard_id)
);

CREATE TABLE psu_specs (
    model_id BIGINT PRIMARY KEY REFERENCES models(id) ON DELETE CASCADE,
    rated_power_w INTEGER CHECK (rated_power_w IS NULL OR rated_power_w > 0),
    efficiency_certification_id BIGINT REFERENCES psu_efficiency_certifications(id),
    psu_standard_id BIGINT REFERENCES psu_standards(id),
    form_factor_id BIGINT REFERENCES psu_form_factors(id),
    modular_type_id BIGINT REFERENCES psu_modular_types(id),
    input_voltage_min_v NUMERIC(7,2),
    input_voltage_max_v NUMERIC(7,2),
    frequency_min_hz NUMERIC(7,2),
    frequency_max_hz NUMERIC(7,2),
    twelve_v_output_w NUMERIC(8,2) CHECK (twelve_v_output_w IS NULL OR twelve_v_output_w >= 0),
    sata_power_connector_count SMALLINT CHECK (sata_power_connector_count IS NULL OR sata_power_connector_count >= 0),
    molex_connector_count SMALLINT CHECK (molex_connector_count IS NULL OR molex_connector_count >= 0),
    fan_size_id BIGINT REFERENCES fan_sizes(id),
    zero_rpm_mode BOOLEAN,
    length_mm NUMERIC(7,2),
    width_mm NUMERIC(7,2),
    height_mm NUMERIC(7,2),
    warranty_years NUMERIC(4,1) CHECK (warranty_years IS NULL OR warranty_years >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE psu_power_connectors (
    psu_model_id BIGINT NOT NULL REFERENCES psu_specs(model_id) ON DELETE CASCADE,
    power_connector_type_id BIGINT NOT NULL REFERENCES power_connector_types(id),
    connector_count SMALLINT NOT NULL DEFAULT 1 CHECK (connector_count > 0),
    PRIMARY KEY (psu_model_id, power_connector_type_id)
);

CREATE TABLE psu_protection_features (
    psu_model_id BIGINT NOT NULL REFERENCES psu_specs(model_id) ON DELETE CASCADE,
    protection_type_id BIGINT NOT NULL REFERENCES psu_protection_types(id),
    PRIMARY KEY (psu_model_id, protection_type_id)
);

CREATE TABLE case_specs (
    model_id BIGINT PRIMARY KEY REFERENCES models(id) ON DELETE CASCADE,
    case_type_id BIGINT REFERENCES case_types(id),
    height_mm NUMERIC(8,2),
    width_mm NUMERIC(8,2),
    depth_mm NUMERIC(8,2),
    volume_l NUMERIC(8,2),
    max_gpu_length_mm NUMERIC(8,2),
    max_cpu_cooler_height_mm NUMERIC(8,2),
    max_psu_length_mm NUMERIC(8,2),
    expansion_slot_count SMALLINT CHECK (expansion_slot_count IS NULL OR expansion_slot_count >= 0),
    included_fan_count SMALLINT CHECK (included_fan_count IS NULL OR included_fan_count >= 0),
    max_fan_count SMALLINT CHECK (max_fan_count IS NULL OR max_fan_count >= 0),
    dust_filters_present BOOLEAN,
    rgb_supported BOOLEAN,
    vertical_gpu_mount_supported BOOLEAN,
    cable_management_space_mm NUMERIC(7,2),
    side_panel_material_id BIGINT REFERENCES case_panel_materials(id),
    front_panel_material_id BIGINT REFERENCES case_panel_materials(id),
    color_id BIGINT REFERENCES product_colors(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE case_supported_motherboard_form_factors (
    case_model_id BIGINT NOT NULL REFERENCES case_specs(model_id) ON DELETE CASCADE,
    motherboard_form_factor_id BIGINT NOT NULL REFERENCES motherboard_form_factors(id),
    PRIMARY KEY (case_model_id, motherboard_form_factor_id)
);

CREATE TABLE case_drive_bays (
    case_model_id BIGINT NOT NULL REFERENCES case_specs(model_id) ON DELETE CASCADE,
    drive_bay_type_id BIGINT NOT NULL REFERENCES drive_bay_types(id),
    bay_count SMALLINT NOT NULL DEFAULT 1 CHECK (bay_count > 0),
    PRIMARY KEY (case_model_id, drive_bay_type_id)
);

CREATE TABLE case_fan_support (
    case_model_id BIGINT NOT NULL REFERENCES case_specs(model_id) ON DELETE CASCADE,
    fan_size_id BIGINT NOT NULL REFERENCES fan_sizes(id),
    supported_count SMALLINT NOT NULL DEFAULT 1 CHECK (supported_count > 0),
    PRIMARY KEY (case_model_id, fan_size_id)
);

CREATE TABLE case_radiator_support (
    case_model_id BIGINT NOT NULL REFERENCES case_specs(model_id) ON DELETE CASCADE,
    radiator_size_id BIGINT NOT NULL REFERENCES radiator_sizes(id),
    supported_count SMALLINT NOT NULL DEFAULT 1 CHECK (supported_count > 0),
    PRIMARY KEY (case_model_id, radiator_size_id)
);

CREATE TABLE case_front_io_ports (
    case_model_id BIGINT NOT NULL REFERENCES case_specs(model_id) ON DELETE CASCADE,
    io_port_type_id BIGINT NOT NULL REFERENCES io_port_types(id),
    port_count SMALLINT NOT NULL DEFAULT 1 CHECK (port_count > 0),
    PRIMARY KEY (case_model_id, io_port_type_id)
);

CREATE TABLE cooling_specs (
    model_id BIGINT PRIMARY KEY REFERENCES models(id) ON DELETE CASCADE,
    cooling_type_id BIGINT NOT NULL REFERENCES cooling_types(id),
    rated_tdp_w INTEGER CHECK (rated_tdp_w IS NULL OR rated_tdp_w > 0),
    radiator_size_id BIGINT REFERENCES radiator_sizes(id),
    fan_size_mm NUMERIC(6,2),
    fan_count SMALLINT CHECK (fan_count IS NULL OR fan_count >= 0),
    fan_speed_min_rpm INTEGER,
    fan_speed_max_rpm INTEGER,
    max_airflow_cfm NUMERIC(8,2),
    static_pressure_mmh2o NUMERIC(8,3),
    noise_level_dba NUMERIC(6,2),
    fan_bearing_type_id BIGINT REFERENCES fan_bearing_types(id),
    pump_speed_min_rpm INTEGER,
    pump_speed_max_rpm INTEGER,
    pump_connector_type_id BIGINT REFERENCES electrical_connector_types(id),
    fan_connector_type_id BIGINT REFERENCES electrical_connector_types(id),
    lighting_type_id BIGINT REFERENCES lighting_types(id),
    cooler_height_mm NUMERIC(7,2),
    thermal_paste_included BOOLEAN,
    warranty_years NUMERIC(4,1) CHECK (warranty_years IS NULL OR warranty_years >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cooling_supported_cpu_sockets (
    cooling_model_id BIGINT NOT NULL REFERENCES cooling_specs(model_id) ON DELETE CASCADE,
    cpu_socket_id BIGINT NOT NULL REFERENCES cpu_sockets(id),
    PRIMARY KEY (cooling_model_id, cpu_socket_id)
);

CREATE TABLE cooling_fans (
    id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    cooling_model_id BIGINT NOT NULL REFERENCES cooling_specs(model_id) ON DELETE CASCADE,
    fan_size_id BIGINT REFERENCES fan_sizes(id),
    fan_count SMALLINT NOT NULL DEFAULT 1 CHECK (fan_count > 0),
    speed_min_rpm INTEGER,
    speed_max_rpm INTEGER,
    max_airflow_cfm NUMERIC(8,2),
    static_pressure_mmh2o NUMERIC(8,3),
    noise_level_dba NUMERIC(6,2),
    fan_bearing_type_id BIGINT REFERENCES fan_bearing_types(id),
    fan_connector_type_id BIGINT REFERENCES electrical_connector_types(id)
);

-- Helpful indexes for common foreign-key joins and filters
CREATE INDEX idx_cpu_specs_socket_id ON cpu_specs(socket_id);
CREATE INDEX idx_cpu_specs_architecture_id ON cpu_specs(architecture_id);
CREATE INDEX idx_gpu_specs_chipset_id ON gpu_specs(chipset_id);
CREATE INDEX idx_gpu_specs_vram_type_id ON gpu_specs(vram_type_id);
CREATE INDEX idx_gpu_specs_pcie_generation_id ON gpu_specs(pcie_generation_id);
CREATE INDEX idx_motherboard_specs_socket_id ON motherboard_specs(socket_id);
CREATE INDEX idx_motherboard_specs_chipset_id ON motherboard_specs(chipset_id);
CREATE INDEX idx_ram_specs_memory_type_id ON ram_specs(memory_type_id);
CREATE INDEX idx_storage_specs_device_type_id ON storage_specs(device_type_id);
CREATE INDEX idx_storage_specs_capacity_gb ON storage_specs(capacity_gb);
CREATE INDEX idx_psu_specs_rated_power_w ON psu_specs(rated_power_w);
CREATE INDEX idx_case_specs_case_type_id ON case_specs(case_type_id);
CREATE INDEX idx_cooling_specs_cooling_type_id ON cooling_specs(cooling_type_id);
