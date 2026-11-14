CREATE TABLE motherboard_wireless_standards (
    motherboard_model_id BIGINT NOT NULL REFERENCES motherboard_specs(model_id) ON DELETE CASCADE,
    wireless_standard_id BIGINT NOT NULL REFERENCES wireless_standards(id),
    PRIMARY KEY (motherboard_model_id, wireless_standard_id)
);
