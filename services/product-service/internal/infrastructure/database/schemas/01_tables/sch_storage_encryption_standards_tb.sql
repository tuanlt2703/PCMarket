CREATE TABLE storage_encryption_standards (
    storage_model_id BIGINT NOT NULL REFERENCES storage_specs(model_id) ON DELETE CASCADE,
    encryption_standard_id BIGINT NOT NULL REFERENCES encryption_standards(id),
    PRIMARY KEY (storage_model_id, encryption_standard_id)
);
