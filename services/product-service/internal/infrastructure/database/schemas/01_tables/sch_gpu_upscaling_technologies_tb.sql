CREATE TABLE gpu_upscaling_technologies (
    gpu_model_id BIGINT NOT NULL REFERENCES gpu_specs(model_id) ON DELETE CASCADE,
    upscaling_technology_id BIGINT NOT NULL REFERENCES upscaling_technologies(id),
    PRIMARY KEY (gpu_model_id, upscaling_technology_id)
);
