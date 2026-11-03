package cpu

import (
	"product-service/internal/platform/httpx"
	"product-service/internal/platform/optional"
)

// Upper limits of the database columns, so out-of-range values get a clear
// 422 instead of a database error.
const (
	maxClockMHz           = 999999.99   // NUMERIC(8,2)
	maxPower              = 999999.99   // NUMERIC(8,2)
	maxSize               = 99999999.99 // NUMERIC(10,2)
	maxTemperatureC       = 999.99      // NUMERIC(5,2)
	maxSmallInt     int32 = 32767       // SMALLINT
)

// UpsertRequest is the JSON body of PUT /admin/models/:id/specs for a CPU model.
// It is a partial write: omitted fields are unchanged and an explicit null
// clears them.
type UpsertRequest struct {
	ArchitectureID           optional.Field[int64]   `json:"architecture_id"`
	CPUSeriesID              optional.Field[int64]   `json:"cpu_series_id"`
	SocketID                 optional.Field[int64]   `json:"socket_id"`
	CoreCount                optional.Field[int32]   `json:"core_count"`
	ThreadCount              optional.Field[int32]   `json:"thread_count"`
	BaseClockMHz             optional.Field[float64] `json:"base_clock_mhz"`
	MaxBoostClockMHz         optional.Field[float64] `json:"max_boost_clock_mhz"`
	L1CacheKB                optional.Field[float64] `json:"l1_cache_kb"`
	L2CacheMB                optional.Field[float64] `json:"l2_cache_mb"`
	L3CacheMB                optional.Field[float64] `json:"l3_cache_mb"`
	TDPW                     optional.Field[float64] `json:"tdp_w"`
	HasIntegratedGraphics    optional.Field[bool]    `json:"has_integrated_graphics"`
	IntegratedGPUChipsetID   optional.Field[int64]   `json:"integrated_gpu_chipset_id"`
	MaxMemoryCapacityGB      optional.Field[float64] `json:"max_memory_capacity_gb"`
	MemoryChannelCount       optional.Field[int32]   `json:"memory_channel_count"`
	PCIeGenerationID         optional.Field[int64]   `json:"pcie_generation_id"`
	UnlockedMultiplier       optional.Field[bool]    `json:"unlocked_multiplier"`
	CoolerIncluded           optional.Field[bool]    `json:"cooler_included"`
	MaxOperatingTemperatureC optional.Field[float64] `json:"max_operating_temperature_c"`

	MemoryTypeIDs     optional.Field[[]int64] `json:"memory_type_ids"`
	InstructionSetIDs optional.Field[[]int64] `json:"instruction_set_ids"`
}

// Validate checks each field on its own. Rules that compare several fields
// (e.g. thread_count >= core_count) are checked by the service after the
// patch is merged with the stored values.
func (r UpsertRequest) Validate() error {
	var v httpx.Validator

	httpx.CheckGreaterThan(&v, "architecture_id", r.ArchitectureID.Value, 0)
	httpx.CheckGreaterThan(&v, "cpu_series_id", r.CPUSeriesID.Value, 0)
	httpx.CheckGreaterThan(&v, "socket_id", r.SocketID.Value, 0)
	httpx.CheckGreaterThan(&v, "core_count", r.CoreCount.Value, 0)
	httpx.CheckGreaterThan(&v, "thread_count", r.ThreadCount.Value, 0)
	httpx.CheckPositiveMax(&v, "base_clock_mhz", r.BaseClockMHz.Value, maxClockMHz)
	httpx.CheckPositiveMax(&v, "max_boost_clock_mhz", r.MaxBoostClockMHz.Value, maxClockMHz)
	httpx.CheckRange(&v, "l1_cache_kb", r.L1CacheKB.Value, 0, maxSize)
	httpx.CheckRange(&v, "l2_cache_mb", r.L2CacheMB.Value, 0, maxSize)
	httpx.CheckRange(&v, "l3_cache_mb", r.L3CacheMB.Value, 0, maxSize)
	httpx.CheckRange(&v, "tdp_w", r.TDPW.Value, 0, maxPower)
	httpx.CheckGreaterThan(&v, "integrated_gpu_chipset_id", r.IntegratedGPUChipsetID.Value, 0)
	httpx.CheckPositiveMax(&v, "max_memory_capacity_gb", r.MaxMemoryCapacityGB.Value, maxSize)
	httpx.CheckPositiveMax(&v, "memory_channel_count", r.MemoryChannelCount.Value, maxSmallInt)
	httpx.CheckGreaterThan(&v, "pcie_generation_id", r.PCIeGenerationID.Value, 0)
	httpx.CheckRange(&v, "max_operating_temperature_c", r.MaxOperatingTemperatureC.Value, -maxTemperatureC, maxTemperatureC)
	v.PositiveIDs("memory_type_ids", r.MemoryTypeIDs.Value)
	v.PositiveIDs("instruction_set_ids", r.InstructionSetIDs.Value)

	return v.Err()
}

func (r UpsertRequest) Patch() Patch {
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

// DecodeUpsertRequest parses and validates the JSON body into a Patch.
func DecodeUpsertRequest(body []byte) (Patch, error) {
	var request UpsertRequest
	if err := httpx.DecodeJSON(body, &request); err != nil {
		return Patch{}, err
	}
	return request.Patch(), nil
}
