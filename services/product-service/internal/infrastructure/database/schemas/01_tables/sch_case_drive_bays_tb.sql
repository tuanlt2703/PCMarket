CREATE TABLE case_drive_bays (
    case_model_id BIGINT NOT NULL REFERENCES case_specs(model_id) ON DELETE CASCADE,
    drive_bay_type_id BIGINT NOT NULL REFERENCES drive_bay_types(id),
    bay_count SMALLINT NOT NULL DEFAULT 1 CHECK (bay_count > 0),
    PRIMARY KEY (case_model_id, drive_bay_type_id)
);
