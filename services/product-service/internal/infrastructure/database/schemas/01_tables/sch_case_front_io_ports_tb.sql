CREATE TABLE case_front_io_ports (
    case_model_id BIGINT NOT NULL REFERENCES case_specs(model_id) ON DELETE CASCADE,
    io_port_type_id BIGINT NOT NULL REFERENCES io_port_types(id),
    port_count SMALLINT NOT NULL DEFAULT 1 CHECK (port_count > 0),
    PRIMARY KEY (case_model_id, io_port_type_id)
);
