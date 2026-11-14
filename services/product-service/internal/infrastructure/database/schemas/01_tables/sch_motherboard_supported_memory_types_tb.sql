CREATE TABLE motherboard_supported_memory_types (
    motherboard_model_id BIGINT NOT NULL REFERENCES motherboard_specs(model_id) ON DELETE CASCADE,
    memory_type_id BIGINT NOT NULL REFERENCES memory_types(id),
    PRIMARY KEY (motherboard_model_id, memory_type_id)
);
