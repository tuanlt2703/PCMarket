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

CREATE INDEX idx_cooling_specs_cooling_type_id ON cooling_specs(cooling_type_id);
