CREATE TABLE cpu_specs (
    model_id BIGINT PRIMARY KEY REFERENCES models(id) ON DELETE CASCADE,
    architecture_id BIGINT REFERENCES cpu_architectures(id),
    cpu_series_id BIGINT REFERENCES cpu_series(id),
    socket_id BIGINT REFERENCES cpu_sockets(id),
    core_count INTEGER CHECK (core_count IS NULL OR core_count > 0),
    thread_count INTEGER CHECK (thread_count IS NULL OR thread_count > 0),
    base_clock_mhz NUMERIC(8,2) CHECK (base_clock_mhz IS NULL OR base_clock_mhz > 0),
    max_boost_clock_mhz NUMERIC(8,2) CHECK (max_boost_clock_mhz IS NULL OR max_boost_clock_mhz > 0),
    l1_cache_kb NUMERIC(10,2) CHECK (l1_cache_kb IS NULL OR l1_cache_kb >= 0),
    l2_cache_mb NUMERIC(10,2) CHECK (l2_cache_mb IS NULL OR l2_cache_mb >= 0),
    l3_cache_mb NUMERIC(10,2) CHECK (l3_cache_mb IS NULL OR l3_cache_mb >= 0),
    tdp_w NUMERIC(8,2) CHECK (tdp_w IS NULL OR tdp_w >= 0),
    has_integrated_graphics BOOLEAN,
    integrated_gpu_chipset_id BIGINT REFERENCES gpu_chipsets(id),
    max_memory_capacity_gb NUMERIC(10,2) CHECK (max_memory_capacity_gb IS NULL OR max_memory_capacity_gb > 0),
    memory_channel_count SMALLINT CHECK (memory_channel_count IS NULL OR memory_channel_count > 0),
    pcie_generation_id BIGINT REFERENCES pcie_generations(id),
    unlocked_multiplier BOOLEAN,
    cooler_included BOOLEAN,
    max_operating_temperature_c NUMERIC(5,2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_cpu_specs_socket_id ON cpu_specs(socket_id);
CREATE INDEX idx_cpu_specs_architecture_id ON cpu_specs(architecture_id);
