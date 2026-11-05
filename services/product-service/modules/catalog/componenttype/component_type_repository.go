package componenttype

import (
	"context"
	"product-service/internal/infrastructure"

	"github.com/jackc/pgx/v5"
)

type ComponentTypeRepository struct {
	db *infrastructure.DB
}

func NewComponentTypeRepository(db *infrastructure.DB) *ComponentTypeRepository {
	return &ComponentTypeRepository{db: db}
}

// A helper function to scan a component_types row to ComponentType struct
func scanComponentType(row pgx.Row) (ComponentType, error) {
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
		ct, err := scanComponentType(rows)
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
