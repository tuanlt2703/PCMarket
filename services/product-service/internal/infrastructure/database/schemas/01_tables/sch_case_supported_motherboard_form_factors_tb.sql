CREATE TABLE case_supported_motherboard_form_factors (
    case_model_id BIGINT NOT NULL REFERENCES case_specs(model_id) ON DELETE CASCADE,
    motherboard_form_factor_id BIGINT NOT NULL REFERENCES motherboard_form_factors(id),
    PRIMARY KEY (case_model_id, motherboard_form_factor_id)
);
