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

CREATE INDEX idx_ram_specs_memory_type_id ON ram_specs(memory_type_id);
