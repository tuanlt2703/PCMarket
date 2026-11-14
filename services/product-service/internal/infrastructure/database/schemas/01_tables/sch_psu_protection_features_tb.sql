CREATE TABLE psu_protection_features (
    psu_model_id BIGINT NOT NULL REFERENCES psu_specs(model_id) ON DELETE CASCADE,
    protection_type_id BIGINT NOT NULL REFERENCES psu_protection_types(id),
    PRIMARY KEY (psu_model_id, protection_type_id)
);
