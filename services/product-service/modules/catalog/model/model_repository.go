package model

import (
	"context"
	"fmt"
	"product-service/internal/infrastructure/database"

	"github.com/jackc/pgx/v5"
)

type ModelRepository struct {
	db *database.DB
}

func NewModelRepository(db *database.DB) *ModelRepository {
	return &ModelRepository{db: db}
}

// A helper function to parse a model row to Model struct
func parseModel(row pgx.Row) (Model, error) {
	var md Model
	err := row.Scan(&md.ID,
		&md.ComponentType.ID, &md.ComponentType.Name,
		&md.Manufacturer.ID, &md.Manufacturer.Name,
		&md.ModelName, &md.Description, &md.ReleaseDate, &md.IsActive, &md.CreatedAt, &md.UpdatedAt)
	return md, err
}

func (repo *ModelRepository) GetAll() ([]Model, error) {
	query := `select * from model_detail order by model_name, id`

	rows, err := repo.db.Conn().Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mdList []Model
	for rows.Next() {
		md, err := parseModel(rows)
		if err != nil {
			return nil, err
		}
		mdList = append(mdList, md)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return mdList, nil
}

func (repo *ModelRepository) GetByID(id int) (Model, error) {
	query := `select * from model_detail where id = $1`

	row := repo.db.Conn().QueryRow(context.Background(), query, id)

	md, err := parseModel(row)
	if err != nil {
		return Model{}, err
	}
	return md, nil
}

func (repo *ModelRepository) Create(createModel CreateModel) (Model, error) {
	ctx := context.Background()
	var md Model
	err := pgx.BeginFunc(ctx, repo.db.Conn(), func(tx pgx.Tx) error {
		query := `insert into models (component_type_id, manufacturer_id, model_name, description, release_date, is_active)
		          values ($1, $2, $3, $4, $5::date, coalesce($6, TRUE))
		          returning id`

		var id int64
		err := tx.QueryRow(ctx, query,
			createModel.ComponentTypeID,
			createModel.ManufacturerID,
			createModel.ModelName,
			createModel.Description,
			createModel.ReleaseDate,
			createModel.IsActive).Scan(&id)
		if err != nil {
			return err
		}

		md, err = parseModel(tx.QueryRow(ctx, `select * from model_detail where id = $1`, id))
		return err
	})
	if err != nil {
		return Model{}, err
	}
	return md, nil
}

func (repo *ModelRepository) Update(id int, updateModel UpdateModel) error {
	query := `update models
	             set model_name = $1,
			         description = $2,
			         release_date = $3,
			         is_active = $4
	           where id = $5 `

	result, err := repo.db.Conn().Exec(context.Background(), query,
		updateModel.ModelName,
		updateModel.Description,
		updateModel.ReleaseDate,
		updateModel.IsActive,
		id)

	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("Not found any row matches id %v", id)
	}

	return nil
}

func (repo *ModelRepository) Delete(id int) error {
	query := `delete from models where id = $1`

	result, err := repo.db.Conn().Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("Not found any row matches id %v", id)
	}

	return nil
}
