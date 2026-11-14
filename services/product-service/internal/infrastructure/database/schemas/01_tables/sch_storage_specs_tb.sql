CREATE TABLE storage_specs (
    model_id BIGINT PRIMARY KEY REFERENCES models(id) ON DELETE CASCADE,
    device_type_id BIGINT NOT NULL REFERENCES storage_device_types(id),
    capacity_gb NUMERIC(12,2) CHECK (capacity_gb IS NULL OR capacity_gb > 0),
    form_factor_id BIGINT REFERENCES storage_form_factors(id),
    interface_id BIGINT REFERENCES storage_interfaces(id),
    protocol_id BIGINT REFERENCES storage_protocols(id),
    pcie_generation_id BIGINT REFERENCES pcie_generations(id),
    sequential_read_mbps INTEGER CHECK (sequential_read_mbps IS NULL OR sequential_read_mbps > 0),
    sequential_write_mbps INTEGER CHECK (sequential_write_mbps IS NULL OR sequential_write_mbps > 0),
    random_read_iops BIGINT CHECK (random_read_iops IS NULL OR random_read_iops > 0),
    random_write_iops BIGINT CHECK (random_write_iops IS NULL OR random_write_iops > 0),
    nand_type_id BIGINT REFERENCES nand_types(id),
    dram_cache_present BOOLEAN,
    endurance_tbw NUMERIC(12,2) CHECK (endurance_tbw IS NULL OR endurance_tbw >= 0),
    mtbf_hours BIGINT CHECK (mtbf_hours IS NULL OR mtbf_hours > 0),
    warranty_years NUMERIC(4,1) CHECK (warranty_years IS NULL OR warranty_years >= 0),
    rpm INTEGER CHECK (rpm IS NULL OR rpm > 0),
    cache_mb INTEGER CHECK (cache_mb IS NULL OR cache_mb >= 0),
    power_loss_protection BOOLEAN,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_storage_specs_device_type_id ON storage_specs(device_type_id);
CREATE INDEX idx_storage_specs_capacity_gb ON storage_specs(capacity_gb);
