package cpu

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"product-service/internal/platform/apperror"
	"product-service/internal/platform/database"
)

// Repository reads and writes cpu_specs and its junction tables.
type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// selectSpec reads one spec row. Parameters:
//
//	$1 = model id
//	$2 = include inactive lookups. When false, a reference to an inactive
//	     lookup value is returned as NULL (or left out of the id lists), so the
//	     public API never shows it.
//
// Column order must match the Scan call in Repository.get.
var selectSpec = `
SELECT s.model_id,
       ` + lookupColumn("cpu_architectures", "architecture_id") + `,
       ` + lookupColumn("cpu_series", "cpu_series_id") + `,
       ` + lookupColumn("cpu_sockets", "socket_id") + `,
       s.core_count, s.thread_count,
       s.base_clock_mhz, s.max_boost_clock_mhz,
       s.l1_cache_kb, s.l2_cache_mb, s.l3_cache_mb, s.tdp_w,
       s.has_integrated_graphics,
       ` + lookupColumn("gpu_chipsets", "integrated_gpu_chipset_id") + `,
       s.max_memory_capacity_gb, s.memory_channel_count,
       ` + lookupColumn("pcie_generations", "pcie_generation_id") + `,
       s.unlocked_multiplier, s.cooler_included, s.max_operating_temperature_c,
       s.created_at, s.updated_at,
       ` + junctionIDs("cpu_supported_memory_types", "memory_type_id", "memory_types") + `,
       ` + junctionIDs("cpu_supported_instruction_sets", "instruction_set_id", "cpu_instruction_sets") + `
FROM cpu_specs s
WHERE s.model_id = $1`

// lookupColumn returns SQL that reads s.<column>, or NULL when the referenced
// lookup row is inactive and $2 (include inactive) is false.
func lookupColumn(lookupTable, column string) string {
	return fmt.Sprintf(
		"CASE WHEN $2::boolean OR EXISTS (SELECT 1 FROM %s l WHERE l.id = s.%s AND l.is_active) THEN s.%s END",
		lookupTable, column, column)
}

// junctionIDs returns SQL that reads the ids stored in a junction table as an
// array, skipping inactive lookup rows when $2 (include inactive) is false.
func junctionIDs(junctionTable, column, lookupTable string) string {
	return fmt.Sprintf(`ARRAY(
           SELECT j.%s FROM %s j JOIN %s l ON l.id = j.%s
           WHERE j.cpu_model_id = s.model_id AND ($2::boolean OR l.is_active)
           ORDER BY j.%s)`,
		column, junctionTable, lookupTable, column, column)
}

// Get returns the spec including references to inactive lookup values (admin view).
func (r *Repository) Get(ctx context.Context, modelID int64) (Spec, error) {
	return r.get(ctx, modelID, true, false)
}

// GetPublic returns the spec without references to inactive lookup values.
func (r *Repository) GetPublic(ctx context.Context, modelID int64) (Spec, error) {
	return r.get(ctx, modelID, false, false)
}

// GetForUpdate locks the spec row until the transaction ends.
// It must be called inside a transaction.
func (r *Repository) GetForUpdate(ctx context.Context, modelID int64) (Spec, error) {
	return r.get(ctx, modelID, true, true)
}

func (r *Repository) get(ctx context.Context, modelID int64, includeInactive, lockRow bool) (Spec, error) {
	sql := selectSpec
	if lockRow {
		sql += " FOR UPDATE OF s"
	}

	var s Spec
	err := r.db.Conn(ctx).QueryRow(ctx, sql, modelID, includeInactive).Scan(
		&s.ModelID,
		&s.ArchitectureID, &s.CPUSeriesID, &s.SocketID,
		&s.CoreCount, &s.ThreadCount,
		&s.BaseClockMHz, &s.MaxBoostClockMHz,
		&s.L1CacheKB, &s.L2CacheMB, &s.L3CacheMB, &s.TDPW,
		&s.HasIntegratedGraphics,
		&s.IntegratedGPUChipsetID,
		&s.MaxMemoryCapacityGB, &s.MemoryChannelCount,
		&s.PCIeGenerationID,
		&s.UnlockedMultiplier, &s.CoolerIncluded, &s.MaxOperatingTemperatureC,
		&s.CreatedAt, &s.UpdatedAt,
		&s.MemoryTypeIDs,
		&s.InstructionSetIDs)
	if err != nil {
		return Spec{}, database.MapError(err, notFoundMessage(modelID))
	}
	return s, nil
}

// Upsert inserts or overwrites the spec row and replaces its junction rows.
// It must be called inside a transaction so all three writes succeed or fail together.
func (r *Repository) Upsert(ctx context.Context, modelID int64, attributes Attributes) error {
	const sql = `
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
    updated_at = CURRENT_TIMESTAMP`

	a := attributes
	conn := r.db.Conn(ctx)
	_, err := conn.Exec(ctx, sql,
		modelID, a.ArchitectureID, a.CPUSeriesID, a.SocketID, a.CoreCount, a.ThreadCount,
		a.BaseClockMHz, a.MaxBoostClockMHz, a.L1CacheKB, a.L2CacheMB, a.L3CacheMB, a.TDPW,
		a.HasIntegratedGraphics, a.IntegratedGPUChipsetID, a.MaxMemoryCapacityGB,
		a.MemoryChannelCount, a.PCIeGenerationID, a.UnlockedMultiplier, a.CoolerIncluded,
		a.MaxOperatingTemperatureC)
	if err != nil {
		return database.MapError(err, "")
	}

	err = replaceJunctionRows(ctx, conn, "cpu_supported_memory_types", "memory_type_id", modelID, a.MemoryTypeIDs)
	if err != nil {
		return err
	}
	return replaceJunctionRows(ctx, conn, "cpu_supported_instruction_sets", "instruction_set_id", modelID, a.InstructionSetIDs)
}

func (r *Repository) Delete(ctx context.Context, modelID int64) error {
	result, err := r.db.Conn(ctx).Exec(ctx, "DELETE FROM cpu_specs WHERE model_id = $1", modelID)
	if err != nil {
		return database.MapError(err, "")
	}
	if result.RowsAffected() == 0 {
		return apperror.NotFound(notFoundMessage(modelID))
	}
	return nil
}

// replaceJunctionRows deletes all rows of a model in a junction table and
// inserts the given ids (duplicates are ignored).
// table and column must be constants written in code.
func replaceJunctionRows(ctx context.Context, conn database.Querier, table, column string, modelID int64, ids []int64) error {
	quotedTable := pgx.Identifier{table}.Sanitize()
	quotedColumn := pgx.Identifier{column}.Sanitize()

	deleteSQL := fmt.Sprintf("DELETE FROM %s WHERE cpu_model_id = $1", quotedTable)
	if _, err := conn.Exec(ctx, deleteSQL, modelID); err != nil {
		return database.MapError(err, "")
	}
	if len(ids) == 0 {
		return nil
	}

	// unnest turns the id array into rows, so all ids are inserted in one statement.
	insertSQL := fmt.Sprintf(
		"INSERT INTO %s (cpu_model_id, %s) SELECT DISTINCT $1::bigint, u.id FROM unnest($2::bigint[]) AS u(id)",
		quotedTable, quotedColumn)
	if _, err := conn.Exec(ctx, insertSQL, modelID, ids); err != nil {
		return database.MapError(err, "")
	}
	return nil
}

func notFoundMessage(modelID int64) string {
	return fmt.Sprintf("cpu specification for model %d not found", modelID)
}
