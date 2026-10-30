package cpu

import (
	"product-service/internal/modules/specs"
	"product-service/internal/platform/httpx"
	"product-service/internal/platform/optional"
)

// UpsertSpecRequest is a partial write: omitted fields are unchanged and an
// explicit null clears them. Upper bounds match the NUMERIC/SMALLINT columns
// so out-of-range values are rejected before reaching the database.
type UpsertSpecRequest struct {
	ArchitectureID           optional.Field[int64]   `json:"architecture_id" binding:"omitempty,gt=0"`
	CPUSeriesID              optional.Field[int64]   `json:"cpu_series_id" binding:"omitempty,gt=0"`
	SocketID                 optional.Field[int64]   `json:"socket_id" binding:"omitempty,gt=0"`
	CoreCount                optional.Field[int32]   `json:"core_count" binding:"omitempty,gt=0"`
	ThreadCount              optional.Field[int32]   `json:"thread_count" binding:"omitempty,gt=0"`
	BaseClockMHz             optional.Field[float64] `json:"base_clock_mhz" binding:"omitempty,gt=0,lte=999999.99"`
	MaxBoostClockMHz         optional.Field[float64] `json:"max_boost_clock_mhz" binding:"omitempty,gt=0,lte=999999.99"`
	L1CacheKB                optional.Field[float64] `json:"l1_cache_kb" binding:"omitempty,gte=0,lte=99999999.99"`
	L2CacheMB                optional.Field[float64] `json:"l2_cache_mb" binding:"omitempty,gte=0,lte=99999999.99"`
	L3CacheMB                optional.Field[float64] `json:"l3_cache_mb" binding:"omitempty,gte=0,lte=99999999.99"`
	TDPW                     optional.Field[float64] `json:"tdp_w" binding:"omitempty,gte=0,lte=999999.99"`
	HasIntegratedGraphics    optional.Field[bool]    `json:"has_integrated_graphics"`
	IntegratedGPUChipsetID   optional.Field[int64]   `json:"integrated_gpu_chipset_id" binding:"omitempty,gt=0"`
	MaxMemoryCapacityGB      optional.Field[float64] `json:"max_memory_capacity_gb" binding:"omitempty,gt=0,lte=99999999.99"`
	MemoryChannelCount       optional.Field[int32]   `json:"memory_channel_count" binding:"omitempty,gt=0,lte=32767"`
	PCIeGenerationID         optional.Field[int64]   `json:"pcie_generation_id" binding:"omitempty,gt=0"`
	UnlockedMultiplier       optional.Field[bool]    `json:"unlocked_multiplier"`
	CoolerIncluded           optional.Field[bool]    `json:"cooler_included"`
	MaxOperatingTemperatureC optional.Field[float64] `json:"max_operating_temperature_c" binding:"omitempty,gte=-999.99,lte=999.99"`

	MemoryTypeIDs     optional.Field[[]int64] `json:"memory_type_ids" binding:"omitempty,dive,gt=0"`
	InstructionSetIDs optional.Field[[]int64] `json:"instruction_set_ids" binding:"omitempty,dive,gt=0"`
}

func (r UpsertSpecRequest) Input() Patch {
	return Patch{
		ArchitectureID:           r.ArchitectureID,
		CPUSeriesID:              r.CPUSeriesID,
		SocketID:                 r.SocketID,
		CoreCount:                r.CoreCount,
		ThreadCount:              r.ThreadCount,
		BaseClockMHz:             r.BaseClockMHz,
		MaxBoostClockMHz:         r.MaxBoostClockMHz,
		L1CacheKB:                r.L1CacheKB,
		L2CacheMB:                r.L2CacheMB,
		L3CacheMB:                r.L3CacheMB,
		TDPW:                     r.TDPW,
		HasIntegratedGraphics:    r.HasIntegratedGraphics,
		IntegratedGPUChipsetID:   r.IntegratedGPUChipsetID,
		MaxMemoryCapacityGB:      r.MaxMemoryCapacityGB,
		MemoryChannelCount:       r.MemoryChannelCount,
		PCIeGenerationID:         r.PCIeGenerationID,
		UnlockedMultiplier:       r.UnlockedMultiplier,
		CoolerIncluded:           r.CoolerIncluded,
		MaxOperatingTemperatureC: r.MaxOperatingTemperatureC,
		MemoryTypeIDs:            r.MemoryTypeIDs,
		InstructionSetIDs:        r.InstructionSetIDs,
	}
}

func decodeRequest(body []byte) (specs.Input, error) {
	var request UpsertSpecRequest
	if err := httpx.DecodeJSON(body, &request); err != nil {
		return nil, err
	}
	return request.Input(), nil
}
