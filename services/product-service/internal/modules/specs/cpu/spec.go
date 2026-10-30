package cpu

import (
	"time"

	"product-service/internal/platform/optional"
)

// Spec mirrors cpu_specs plus its junction tables.
type Spec struct {
	ModelID int64 `json:"model_id"`
	Attributes
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Attributes are the writable CPU attributes; foreign keys reference lookup tables.
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

	MemoryTypeIDs     []int64 `json:"memory_type_ids"`
	InstructionSetIDs []int64 `json:"instruction_set_ids"`
}

// Patch is a partial CPU specification write: omitted fields keep their stored
// value (or NULL on first write) and an explicit null clears them.
type Patch struct {
	ArchitectureID           optional.Field[int64]
	CPUSeriesID              optional.Field[int64]
	SocketID                 optional.Field[int64]
	CoreCount                optional.Field[int32]
	ThreadCount              optional.Field[int32]
	BaseClockMHz             optional.Field[float64]
	MaxBoostClockMHz         optional.Field[float64]
	L1CacheKB                optional.Field[float64]
	L2CacheMB                optional.Field[float64]
	L3CacheMB                optional.Field[float64]
	TDPW                     optional.Field[float64]
	HasIntegratedGraphics    optional.Field[bool]
	IntegratedGPUChipsetID   optional.Field[int64]
	MaxMemoryCapacityGB      optional.Field[float64]
	MemoryChannelCount       optional.Field[int32]
	PCIeGenerationID         optional.Field[int64]
	UnlockedMultiplier       optional.Field[bool]
	CoolerIncluded           optional.Field[bool]
	MaxOperatingTemperatureC optional.Field[float64]

	MemoryTypeIDs     optional.Field[[]int64]
	InstructionSetIDs optional.Field[[]int64]
}

func (Patch) ComponentTypeCode() string { return ComponentTypeCode }

// Apply returns a copy of a with the present patch fields applied.
func (p Patch) Apply(a Attributes) Attributes {
	p.ArchitectureID.ApplyTo(&a.ArchitectureID)
	p.CPUSeriesID.ApplyTo(&a.CPUSeriesID)
	p.SocketID.ApplyTo(&a.SocketID)
	p.CoreCount.ApplyTo(&a.CoreCount)
	p.ThreadCount.ApplyTo(&a.ThreadCount)
	p.BaseClockMHz.ApplyTo(&a.BaseClockMHz)
	p.MaxBoostClockMHz.ApplyTo(&a.MaxBoostClockMHz)
	p.L1CacheKB.ApplyTo(&a.L1CacheKB)
	p.L2CacheMB.ApplyTo(&a.L2CacheMB)
	p.L3CacheMB.ApplyTo(&a.L3CacheMB)
	p.TDPW.ApplyTo(&a.TDPW)
	p.HasIntegratedGraphics.ApplyTo(&a.HasIntegratedGraphics)
	p.IntegratedGPUChipsetID.ApplyTo(&a.IntegratedGPUChipsetID)
	p.MaxMemoryCapacityGB.ApplyTo(&a.MaxMemoryCapacityGB)
	p.MemoryChannelCount.ApplyTo(&a.MemoryChannelCount)
	p.PCIeGenerationID.ApplyTo(&a.PCIeGenerationID)
	p.UnlockedMultiplier.ApplyTo(&a.UnlockedMultiplier)
	p.CoolerIncluded.ApplyTo(&a.CoolerIncluded)
	p.MaxOperatingTemperatureC.ApplyTo(&a.MaxOperatingTemperatureC)
	p.MemoryTypeIDs.ApplyValue(&a.MemoryTypeIDs)
	p.InstructionSetIDs.ApplyValue(&a.InstructionSetIDs)
	return a
}
