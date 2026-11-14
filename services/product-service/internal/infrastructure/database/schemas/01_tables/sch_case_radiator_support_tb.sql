CREATE TABLE case_radiator_support (
    case_model_id BIGINT NOT NULL REFERENCES case_specs(model_id) ON DELETE CASCADE,
    radiator_size_id BIGINT NOT NULL REFERENCES radiator_sizes(id),
    supported_count SMALLINT NOT NULL DEFAULT 1 CHECK (supported_count > 0),
    PRIMARY KEY (case_model_id, radiator_size_id)
);
