package cpu

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"product-service/internal/platform/database"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// specSelect reads a spec row. $2 (include inactive) false hides references to
// inactive lookup values, so the public API never exposes them.
var specSelect = fmt.Sprintf(`
SELECT s.model_id, %s, %s, %s, s.core_count, s.thread_count,
       s.base_clock_mhz, s.max_boost_clock_mhz, s.l1_cache_kb, s.l2_cache_mb, s.l3_cache_mb, s.tdp_w,
       s.has_integrated_graphics, %s, s.max_memory_capacity_gb,
       s.memory_channel_count, %s, s.unlocked_multiplier, s.cooler_included,
       s.max_operating_temperature_c, s.created_at, s.updated_at,
       %s,
       %s
FROM cpu_specs s
WHERE s.model_id = $1`,
	activeRef("cpu_architectures", "architecture_id"),
	activeRef("cpu_series", "cpu_series_id"),
	activeRef("cpu_sockets", "socket_id"),
	activeRef("gpu_chipsets", "integrated_gpu_chipset_id"),
	activeRef("pcie_generations", "pcie_generation_id"),
	activeLinks("cpu_supported_memory_types", "memory_type_id", "memory_types"),
	activeLinks("cpu_supported_instruction_sets", "instruction_set_id", "cpu_instruction_sets"),
)

func activeRef(lookupTable, column string) string {
	return fmt.Sprintf(
		"CASE WHEN $2::boolean OR EXISTS (SELECT 1 FROM %s l WHERE l.id = s.%s AND l.is_active) THEN s.%s END",
		lookupTable, column, column)
}

func activeLinks(junction, column, lookupTable string) string {
	return fmt.Sprintf(`ARRAY(SELECT j.%[2]s FROM %[1]s j JOIN %[3]s l ON l.id = j.%[2]s
             WHERE j.cpu_model_id = s.model_id AND ($2::boolean OR l.is_active) ORDER BY j.%[2]s)`,
		junction, column, lookupTable)
}

// Get returns the spec including references to inactive lookup values (admin view).
func (r *Repository) Get(ctx context.Context, modelID int64) (Spec, error) {
	return r.get(ctx, modelID, true, false)
}

// GetPublic returns the spec with inactive lookup references removed.
func (r *Repository) GetPublic(ctx context.Context, modelID int64) (Spec, error) {
	return r.get(ctx, modelID, false, false)
}

// GetForUpdate locks the spec row; it must run inside a transaction.
func (r *Repository) GetForUpdate(ctx context.Context, modelID int64) (Spec, error) {
	return r.get(ctx, modelID, true, true)
}

func (r *Repository) get(ctx context.Context, modelID int64, includeInactive, lock bool) (Spec, error) {
	sql := specSelect
	if lock {
		sql += " FOR UPDATE OF s"
	}
	var s Spec
	err := r.db.Q(ctx).QueryRow(ctx, sql, modelID, includeInactive).Scan(
		&s.ModelID, &s.ArchitectureID, &s.CPUSeriesID, &s.SocketID, &s.CoreCount, &s.ThreadCount,
		&s.BaseClockMHz, &s.MaxBoostClockMHz, &s.L1CacheKB, &s.L2CacheMB, &s.L3CacheMB, &s.TDPW,
		&s.HasIntegratedGraphics, &s.IntegratedGPUChipsetID, &s.MaxMemoryCapacityGB,
		&s.MemoryChannelCount, &s.PCIeGenerationID, &s.UnlockedMultiplier, &s.CoolerIncluded,
		&s.MaxOperatingTemperatureC, &s.CreatedAt, &s.UpdatedAt,
		&s.MemoryTypeIDs, &s.InstructionSetIDs)
	return s, database.MapError(err, notFound(modelID))
}

// Upsert writes the spec row and replaces its junction rows. Callers must run it in a transaction.
func (r *Repository) Upsert(ctx context.Context, modelID int64, f Attributes) error {
	q := r.db.Q(ctx)
	_, err := q.Exec(ctx, `
INSERT INTO cpu_specs (
    model_id, architecture_id, cpu_series_id, socket_id, core_count, thread_count,
    base_clock_mhz, max_boost_clock_mhz, l1_cache_kb, l2_cache_mb, l3_cache_mb, tdp_w,
    has_integrated_graphics, integrated_gpu_chipset_id, max_memory_capacity_gb,
    memory_channel_count, pcie_generation_id, unlocked_multiplier, cooler_included,
    max_operating_temperature_c)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
ON CONFLICT (model_id) DO UPDATE SET
    architecture_id = EXCLUDED.architecture_id,
    cpu_series_id = EXCLUDED.cpu_series_id,
    socket_id = EXCLUDED.socket_id,
    core_count = EXCLUDED.core_count,
    thread_count = EXCLUDED.thread_count,
    base_clock_mhz = EXCLUDED.base_clock_mhz,
    max_boost_clock_mhz = EXCLUDED.max_boost_clock_mhz,
    l1_cache_kb = EXCLUDED.l1_cache_kb,
    l2_cache_mb = EXCLUDED.l2_cache_mb,
    l3_cache_mb = EXCLUDED.l3_cache_mb,
    tdp_w = EXCLUDED.tdp_w,
    has_integrated_graphics = EXCLUDED.has_integrated_graphics,
    integrated_gpu_chipset_id = EXCLUDED.integrated_gpu_chipset_id,
    max_memory_capacity_gb = EXCLUDED.max_memory_capacity_gb,
    memory_channel_count = EXCLUDED.memory_channel_count,
    pcie_generation_id = EXCLUDED.pcie_generation_id,
    unlocked_multiplier = EXCLUDED.unlocked_multiplier,
    cooler_included = EXCLUDED.cooler_included,
    max_operating_temperature_c = EXCLUDED.max_operating_temperature_c,
    updated_at = CURRENT_TIMESTAMP`,
		modelID, f.ArchitectureID, f.CPUSeriesID, f.SocketID, f.CoreCount, f.ThreadCount,
		f.BaseClockMHz, f.MaxBoostClockMHz, f.L1CacheKB, f.L2CacheMB, f.L3CacheMB, f.TDPW,
		f.HasIntegratedGraphics, f.IntegratedGPUChipsetID, f.MaxMemoryCapacityGB,
		f.MemoryChannelCount, f.PCIeGenerationID, f.UnlockedMultiplier, f.CoolerIncluded,
		f.MaxOperatingTemperatureC)
	if err != nil {
		return database.MapError(err, "")
	}

	if err := replaceLinks(ctx, q, "cpu_supported_memory_types", "memory_type_id", modelID, f.MemoryTypeIDs); err != nil {
		return err
	}
	return replaceLinks(ctx, q, "cpu_supported_instruction_sets", "instruction_set_id", modelID, f.InstructionSetIDs)
}

func (r *Repository) Delete(ctx context.Context, modelID int64) error {
	tag, err := r.db.Q(ctx).Exec(ctx, "DELETE FROM cpu_specs WHERE model_id = $1", modelID)
	if err != nil {
		return database.MapError(err, "")
	}
	if tag.RowsAffected() == 0 {
		return database.MapError(pgx.ErrNoRows, notFound(modelID))
	}
	return nil
}

// replaceLinks replaces junction rows; table and column are compile-time constants.
func replaceLinks(ctx context.Context, q database.Querier, table, column string, modelID int64, ids []int64) error {
	t, col := pgx.Identifier{table}.Sanitize(), pgx.Identifier{column}.Sanitize()
	if _, err := q.Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE cpu_model_id = $1", t), modelID); err != nil {
		return database.MapError(err, "")
	}
	if len(ids) == 0 {
		return nil
	}
	_, err := q.Exec(ctx, fmt.Sprintf(
		"INSERT INTO %s (cpu_model_id, %s) SELECT DISTINCT $1::bigint, u.id FROM unnest($2::bigint[]) AS u(id)",
		t, col), modelID, ids)
	return database.MapError(err, "")
}

func notFound(modelID int64) string {
	return fmt.Sprintf("cpu specification for model %d not found", modelID)
}
