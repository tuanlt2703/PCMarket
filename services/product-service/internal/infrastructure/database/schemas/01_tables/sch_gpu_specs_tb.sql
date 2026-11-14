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

CREATE INDEX idx_gpu_specs_chipset_id ON gpu_specs(chipset_id);
CREATE INDEX idx_gpu_specs_vram_type_id ON gpu_specs(vram_type_id);
CREATE INDEX idx_gpu_specs_pcie_generation_id ON gpu_specs(pcie_generation_id);
