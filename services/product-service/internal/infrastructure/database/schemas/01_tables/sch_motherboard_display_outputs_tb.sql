CREATE TABLE motherboard_display_outputs (
    motherboard_model_id BIGINT NOT NULL REFERENCES motherboard_specs(model_id) ON DELETE CASCADE,
    display_connector_type_id BIGINT NOT NULL REFERENCES display_connector_types(id),
    port_count SMALLINT NOT NULL DEFAULT 1 CHECK (port_count > 0),
    PRIMARY KEY (motherboard_model_id, display_connector_type_id)
);
