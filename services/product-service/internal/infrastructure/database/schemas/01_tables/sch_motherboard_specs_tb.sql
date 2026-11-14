CREATE TABLE motherboard_specs (
    model_id BIGINT PRIMARY KEY REFERENCES models(id) ON DELETE CASCADE,
    chipset_id BIGINT REFERENCES motherboard_chipsets(id),
    socket_id BIGINT REFERENCES cpu_sockets(id),
    form_factor_id BIGINT REFERENCES motherboard_form_factors(id),
    memory_slot_count SMALLINT CHECK (memory_slot_count IS NULL OR memory_slot_count >= 0),
    max_memory_capacity_gb NUMERIC(10,2) CHECK (max_memory_capacity_gb IS NULL OR max_memory_capacity_gb > 0),
    max_memory_data_rate_mt_s INTEGER CHECK (max_memory_data_rate_mt_s IS NULL OR max_memory_data_rate_mt_s > 0),
    memory_channel_count SMALLINT CHECK (memory_channel_count IS NULL OR memory_channel_count > 0),
    sata_port_count SMALLINT CHECK (sata_port_count IS NULL OR sata_port_count >= 0),
    ethernet_speed_mbps INTEGER CHECK (ethernet_speed_mbps IS NULL OR ethernet_speed_mbps > 0),
    audio_codec_id BIGINT REFERENCES audio_codecs(id),
    vrm_phase_description VARCHAR(100),
    bios_flashback_supported BOOLEAN,
    rgb_header_count SMALLINT CHECK (rgb_header_count IS NULL OR rgb_header_count >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_motherboard_specs_socket_id ON motherboard_specs(socket_id);
CREATE INDEX idx_motherboard_specs_chipset_id ON motherboard_specs(chipset_id);
