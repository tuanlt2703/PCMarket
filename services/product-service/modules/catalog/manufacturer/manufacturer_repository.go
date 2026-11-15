package manufacturer

import (
	"context"
	"fmt"
	"product-service/internal/infrastructure/database"

	"github.com/jackc/pgx/v5"
)

type ManufacturerRepository struct {
	db *database.DB
}

func NewManufacturerRepository(db *database.DB) *ManufacturerRepository {
	return &ManufacturerRepository{db: db}
}

// A helper function to parse a manufacturer row to Manufacturer struct
func parseManufacturer(row pgx.Row) (Manufacturer, error) {
	var ct Manufacturer
	err := row.Scan(&ct.ID, &ct.Name, &ct.Description, &ct.IsActive, &ct.CreatedAt, &ct.UpdatedAt)
	return ct, err
}

func (repo *ManufacturerRepository) GetAll() ([]Manufacturer, error) {
	query := "select * from manufacturer order by name"

	rows, err := repo.db.Conn().Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var componentList []Manufacturer
	for rows.Next() {
		ct, err := parseManufacturer(rows)
		if err != nil {
			return nil, err
		}
		componentList = append(componentList, ct)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return componentList, nil
}

func (repo *ManufacturerRepository) GetByID(id int) (Manufacturer, error) {
	query := "select * from manufacturer where id = $1 order by name"

	row := repo.db.Conn().QueryRow(context.Background(), query, id)

	component, err := parseManufacturer(row)
	if err != nil {
		return Manufacturer{}, err
	}
	return component, nil
}

func (repo *ManufacturerRepository) Create(newComponent CreateManufacturer) (Manufacturer, error) {
	query := `insert into manufacturer (name, description, is_active)
	          values ($1, $2, coallesce($3, true))`

	row := repo.db.Conn().QueryRow(context.Background(), query,
		newComponent.Name,
		newComponent.Description,
		newComponent.IsActive)

	component, err := parseManufacturer(row)
	if err != nil {
		return Manufacturer{}, err
	}
	return component, nil
}

func (repo *ManufacturerRepository) Update(id int, updateComponent UpdateManufacturer) (Manufacturer, error) {
	query := `update manufacturer
	             set description = coallesce($1, description),
			         is_active = coallesce($2, is_active)
	           where id = $3 `

	row := repo.db.Conn().QueryRow(context.Background(), query,
		updateComponent.Description,
		updateComponent.IsActive,
		id)

	component, err := parseManufacturer(row)
	if err != nil {
		return Manufacturer{}, err
	}
	return component, nil
}

func (repo *ManufacturerRepository) Delete(id int) error {
	query := `delete from manufacturer where id = $1`

	result, err := repo.db.Conn().Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("Not found any row matches id %v", id)
	}

	return nil
}
