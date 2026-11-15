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

// A helper function to parse a manufacturers row to manufacturers struct
func parseManufacturer(row pgx.Row) (Manufacturer, error) {
	var mf Manufacturer
	err := row.Scan(&mf.ID, &mf.Name, &mf.Description, &mf.IsActive, &mf.CreatedAt, &mf.UpdatedAt)
	return mf, err
}

func (repo *ManufacturerRepository) GetAll() ([]Manufacturer, error) {
	query := "select * from manufacturers order by name"

	rows, err := repo.db.Conn().Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mfList []Manufacturer
	for rows.Next() {
		mt, err := parseManufacturer(rows)
		if err != nil {
			return nil, err
		}
		mfList = append(mfList, mt)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return mfList, nil
}

func (repo *ManufacturerRepository) GetByID(id int) (Manufacturer, error) {
	query := "select * from manufacturers where id = $1 order by name"

	row := repo.db.Conn().QueryRow(context.Background(), query, id)

	mf, err := parseManufacturer(row)
	if err != nil {
		return Manufacturer{}, err
	}
	return mf, nil
}

func (repo *ManufacturerRepository) Create(createManufacturer CreateManufacturer) (Manufacturer, error) {
	query := `insert into manufacturers (name, description, is_active)
	          values ($1, $2, coallesce($3, true))`

	row := repo.db.Conn().QueryRow(context.Background(), query,
		createManufacturer.Name,
		createManufacturer.Description,
		createManufacturer.IsActive)

	mf, err := parseManufacturer(row)
	if err != nil {
		return Manufacturer{}, err
	}
	return mf, nil
}

func (repo *ManufacturerRepository) Update(id int, updateManufacturer UpdateManufacturer) (Manufacturer, error) {
	query := `update manufacturer
	             set description = coallesce($1, description),
			         is_active = coallesce($2, is_active)
	           where id = $3 `

	row := repo.db.Conn().QueryRow(context.Background(), query,
		updateManufacturer.Description,
		updateManufacturer.IsActive,
		id)

	mf, err := parseManufacturer(row)
	if err != nil {
		return Manufacturer{}, err
	}
	return mf, nil
}

func (repo *ManufacturerRepository) Delete(id int) error {
	query := `delete from manufacturers where id = $1`

	result, err := repo.db.Conn().Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("Not found any row matches id %v", id)
	}

	return nil
}
