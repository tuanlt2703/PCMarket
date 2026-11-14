package componenttype

import (
	"context"
	"fmt"
	"product-service/internal/infrastructure/database"

	"github.com/jackc/pgx/v5"
)

type ComponentTypeRepository struct {
	db *database.DB
}

func NewComponentTypeRepository(db *database.DB) *ComponentTypeRepository {
	return &ComponentTypeRepository{db: db}
}

// A helper function to parse a component_types row to ComponentType struct
func parseComponentType(row pgx.Row) (ComponentType, error) {
	var ct ComponentType
	err := row.Scan(&ct.ID, &ct.Code, &ct.Name, &ct.Description, &ct.IsActive, &ct.CreatedAt, &ct.UpdatedAt)
	return ct, err
}

func (repo *ComponentTypeRepository) GetAll() ([]ComponentType, error) {
	query := "select * from component_types order by name"

	rows, err := repo.db.Conn().Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var componentList []ComponentType
	for rows.Next() {
		ct, err := parseComponentType(rows)
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

func (repo *ComponentTypeRepository) GetByID(id int) (ComponentType, error) {
	query := "select * from component_types where id = $1 order by name"

	row := repo.db.Conn().QueryRow(context.Background(), query, id)

	component, err := parseComponentType(row)
	if err != nil {
		return ComponentType{}, err
	}
	return component, nil
}

func (repo *ComponentTypeRepository) Create(newComponent CreateComponentType) (ComponentType, error) {
	query := `insert into component_types (code, name, description, is_active)
	          values ($1, $2, $3, coallesce($4, true))`

	row := repo.db.Conn().QueryRow(context.Background(), query,
		newComponent.Code,
		newComponent.Name,
		newComponent.Description,
		newComponent.IsActive)

	component, err := parseComponentType(row)
	if err != nil {
		return ComponentType{}, err
	}
	return component, nil
}

func (repo *ComponentTypeRepository) Update(id int, updateComponent UpdateComponentType) (ComponentType, error) {
	query := `update component_types
	             set description = coallesce($1, description),
			         is_active = coallesce($2, is_active)
	           where id = $3 `

	row := repo.db.Conn().QueryRow(context.Background(), query,
		updateComponent.Description,
		updateComponent.IsActive,
		id)

	component, err := parseComponentType(row)
	if err != nil {
		return ComponentType{}, err
	}
	return component, nil
}

func (repo *ComponentTypeRepository) Delete(id int) error {
	query := `delete from component_types where id = $1`

	result, err := repo.db.Conn().Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("Not found any row matches id %v", id)
	}

	return nil
}
