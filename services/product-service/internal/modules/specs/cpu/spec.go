package cpu

import "time"

// Spec is one row of cpu_specs plus the ids from its two junction tables.
//
// Attributes is embedded: its fields can be used as spec.CoreCount, and the
// JSON output is flat ({"model_id":1,"core_count":8,...}), not nested.
type Spec struct {
	ModelID int64 `json:"model_id"`
	Attributes
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Attributes are the writable CPU columns. A nil pointer is NULL in the
// database. Fields ending in ID reference a lookup table.
type Attributes struct {
	ArchitectureID           *int64   `json:"architecture_id"`
	CPUSeriesID              *int64   `json:"cpu_series_id"`
	SocketID                 *int64   `json:"socket_id"`
	CoreCount                *int32   `json:"core_count"`
	ThreadCount              *int32   `json:"thread_count"`
	BaseClockMHz             *float64 `json:"base_clock_mhz"`
	MaxBoostClockMHz         *float64 `json:"max_boost_clock_mhz"`
	L1CacheKB                *float64 `json:"l1_cache_kb"`
	L2CacheMB                *float64 `json:"l2_cache_mb"`
	L3CacheMB                *float64 `json:"l3_cache_mb"`
	TDPW                     *float64 `json:"tdp_w"`
	HasIntegratedGraphics    *bool    `json:"has_integrated_graphics"`
	IntegratedGPUChipsetID   *int64   `json:"integrated_gpu_chipset_id"`
	MaxMemoryCapacityGB      *float64 `json:"max_memory_capacity_gb"`
	MemoryChannelCount       *int32   `json:"memory_channel_count"`
	PCIeGenerationID         *int64   `json:"pcie_generation_id"`
	UnlockedMultiplier       *bool    `json:"unlocked_multiplier"`
	CoolerIncluded           *bool    `json:"cooler_included"`
	MaxOperatingTemperatureC *float64 `json:"max_operating_temperature_c"`

	// Stored in the junction tables cpu_supported_memory_types and
	// cpu_supported_instruction_sets.
	MemoryTypeIDs     []int64 `json:"memory_type_ids"`
	InstructionSetIDs []int64 `json:"instruction_set_ids"`
}
