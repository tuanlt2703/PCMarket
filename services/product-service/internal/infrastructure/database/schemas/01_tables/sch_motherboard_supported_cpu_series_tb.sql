CREATE TABLE motherboard_supported_cpu_series (
    motherboard_model_id BIGINT NOT NULL REFERENCES motherboard_specs(model_id) ON DELETE CASCADE,
    cpu_series_id BIGINT NOT NULL REFERENCES cpu_series(id),
    PRIMARY KEY (motherboard_model_id, cpu_series_id)
);
