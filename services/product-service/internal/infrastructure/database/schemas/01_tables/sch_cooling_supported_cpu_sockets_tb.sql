CREATE TABLE cooling_supported_cpu_sockets (
    cooling_model_id BIGINT NOT NULL REFERENCES cooling_specs(model_id) ON DELETE CASCADE,
    cpu_socket_id BIGINT NOT NULL REFERENCES cpu_sockets(id),
    PRIMARY KEY (cooling_model_id, cpu_socket_id)
);
