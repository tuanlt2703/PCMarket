CREATE TABLE cpu_supported_instruction_sets (
    cpu_model_id BIGINT NOT NULL REFERENCES cpu_specs(model_id) ON DELETE CASCADE,
    instruction_set_id BIGINT NOT NULL REFERENCES cpu_instruction_sets(id),
    PRIMARY KEY (cpu_model_id, instruction_set_id)
);
