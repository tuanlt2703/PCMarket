CREATE TABLE motherboard_bluetooth_versions (
    motherboard_model_id BIGINT NOT NULL REFERENCES motherboard_specs(model_id) ON DELETE CASCADE,
    bluetooth_version_id BIGINT NOT NULL REFERENCES bluetooth_versions(id),
    PRIMARY KEY (motherboard_model_id, bluetooth_version_id)
);
