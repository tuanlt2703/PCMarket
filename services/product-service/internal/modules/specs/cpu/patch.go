package cpu

import "product-service/internal/platform/optional"

// Patch is a partial write of a CPU specification. For each field:
// omitted keeps the stored value (NULL on the first write), null clears it,
// and a value replaces it.
//
// Patch implements specs.Input.
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

func (Patch) ComponentTypeCode() string {
	return ComponentTypeCode
}

// Apply returns current with the present patch fields applied.
// current is passed by value, so the caller's copy is not modified.
func (p Patch) Apply(current Attributes) Attributes {
	result := current
	p.ArchitectureID.ApplyToPointer(&result.ArchitectureID)
	p.CPUSeriesID.ApplyToPointer(&result.CPUSeriesID)
	p.SocketID.ApplyToPointer(&result.SocketID)
	p.CoreCount.ApplyToPointer(&result.CoreCount)
	p.ThreadCount.ApplyToPointer(&result.ThreadCount)
	p.BaseClockMHz.ApplyToPointer(&result.BaseClockMHz)
	p.MaxBoostClockMHz.ApplyToPointer(&result.MaxBoostClockMHz)
	p.L1CacheKB.ApplyToPointer(&result.L1CacheKB)
	p.L2CacheMB.ApplyToPointer(&result.L2CacheMB)
	p.L3CacheMB.ApplyToPointer(&result.L3CacheMB)
	p.TDPW.ApplyToPointer(&result.TDPW)
	p.HasIntegratedGraphics.ApplyToPointer(&result.HasIntegratedGraphics)
	p.IntegratedGPUChipsetID.ApplyToPointer(&result.IntegratedGPUChipsetID)
	p.MaxMemoryCapacityGB.ApplyToPointer(&result.MaxMemoryCapacityGB)
	p.MemoryChannelCount.ApplyToPointer(&result.MemoryChannelCount)
	p.PCIeGenerationID.ApplyToPointer(&result.PCIeGenerationID)
	p.UnlockedMultiplier.ApplyToPointer(&result.UnlockedMultiplier)
	p.CoolerIncluded.ApplyToPointer(&result.CoolerIncluded)
	p.MaxOperatingTemperatureC.ApplyToPointer(&result.MaxOperatingTemperatureC)
	p.MemoryTypeIDs.ApplyToValue(&result.MemoryTypeIDs)
	p.InstructionSetIDs.ApplyToValue(&result.InstructionSetIDs)
	return result
}
