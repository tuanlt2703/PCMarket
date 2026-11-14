CREATE TABLE gpu_display_outputs (
    gpu_model_id BIGINT NOT NULL REFERENCES gpu_specs(model_id) ON DELETE CASCADE,
    display_connector_type_id BIGINT NOT NULL REFERENCES display_connector_types(id),
    port_count SMALLINT NOT NULL DEFAULT 1 CHECK (port_count > 0),
    PRIMARY KEY (gpu_model_id, display_connector_type_id)
);
