CREATE TABLE cpu_supported_memory_types (
    cpu_model_id BIGINT NOT NULL REFERENCES cpu_specs(model_id) ON DELETE CASCADE,
    memory_type_id BIGINT NOT NULL REFERENCES memory_types(id),
    PRIMARY KEY (cpu_model_id, memory_type_id)
);
