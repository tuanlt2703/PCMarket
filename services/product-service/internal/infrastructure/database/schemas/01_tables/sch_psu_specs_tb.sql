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

CREATE INDEX idx_psu_specs_rated_power_w ON psu_specs(rated_power_w);
