CREATE TABLE psu_power_connectors (
    psu_model_id BIGINT NOT NULL REFERENCES psu_specs(model_id) ON DELETE CASCADE,
    power_connector_type_id BIGINT NOT NULL REFERENCES power_connector_types(id),
    connector_count SMALLINT NOT NULL DEFAULT 1 CHECK (connector_count > 0),
    PRIMARY KEY (psu_model_id, power_connector_type_id)
);
