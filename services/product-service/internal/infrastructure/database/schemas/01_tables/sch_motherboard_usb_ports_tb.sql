CREATE TABLE motherboard_usb_ports (
    motherboard_model_id BIGINT NOT NULL REFERENCES motherboard_specs(model_id) ON DELETE CASCADE,
    usb_port_type_id BIGINT NOT NULL REFERENCES usb_port_types(id),
    port_count SMALLINT NOT NULL DEFAULT 1 CHECK (port_count > 0),
    PRIMARY KEY (motherboard_model_id, usb_port_type_id)
);
