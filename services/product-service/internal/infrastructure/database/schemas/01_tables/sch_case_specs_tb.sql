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

CREATE INDEX idx_case_specs_case_type_id ON case_specs(case_type_id);
