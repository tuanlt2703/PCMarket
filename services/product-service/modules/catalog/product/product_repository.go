package product

import (
	"context"
	"fmt"
	"product-service/internal/infrastructure/database"

	"github.com/jackc/pgx/v5"
)

type ProductRepository struct {
	db *database.DB
}

func NewProductRepository(db *database.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// A helper function to parse a model row to ProductModel struct
func parseProductModel(row pgx.Row) (ProductModel, error) {
	var prd ProductModel
	err := row.Scan(&prd.ID,
		&prd.ComponentType.ID, &prd.ComponentType.Name,
		&prd.Manufacturer.ID, &prd.Manufacturer.Name,
		&prd.ModelName, &prd.Description, &prd.ReleaseDate, &prd.IsActive, &prd.CreatedAt, &prd.UpdatedAt)
	return prd, err
}

func (repo *ProductRepository) GetAll() ([]ProductModel, error) {
	query := ` select md.id,
                      ct.id, ct.name,
                      mf.id, mf.name,
                      md.model_name, md.description, md.release_date, md.is_active, md.created_at, md.updated_at
                 from models md
                 join component_types ct on ct.id = md.component_type_id
                 join manufacturers mf on mf.id = md.manufacturer_id
             order by m.model_name, m.id`

	rows, err := repo.db.Conn().Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prdList []ProductModel
	for rows.Next() {
		ct, err := parseProductModel(rows)
		if err != nil {
			return nil, err
		}
		prdList = append(prdList, ct)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return prdList, nil
}

func (repo *ProductRepository) GetByID(id int) (ProductModel, error) {
	query := ` select md.id,
                      ct.id, ct.name,
                      mf.id, mf.name,
                      md.model_name, md.description, md.release_date, md.is_active, md.created_at, md.updated_at
                 from models md
                 join component_types ct on ct.id = md.component_type_id
                 join manufacturers mf on mf.id = md.manufacturer_id
                where md.id = $1"`

	row := repo.db.Conn().QueryRow(context.Background(), query, id)

	prd, err := parseProductModel(row)
	if err != nil {
		return ProductModel{}, err
	}
	return prd, nil
}

func (repo *ProductRepository) Create(createProductModel CreateProductModel) (ProductModel, error) {
	query := ` with new_model as (
	           insert into models (component_type_id, manufacturer_id, model_name, description, release_date, is_active)
                           values ($1, $2, $3, $4, $5::date, COALESCE($6, TRUE))
               returning id
			   )
			   select md.id,
                      ct.id, ct.name,
                      mf.id, mf.name,
                      md.model_name, md.description, md.release_date, md.is_active, md.created_at, md.updated_at
                 from models md
                 join component_types ct on ct.id = md.component_type_id
                 join manufacturers mf on mf.id = md.manufacturer_id
                where md.id = (select id from new_model)`

	row := repo.db.Conn().QueryRow(context.Background(), query,
		createProductModel.ComponentTypeID,
		createProductModel.ManufacturerID,
		createProductModel.ModelName,
		createProductModel.Description,
		createProductModel.ReleaseDate,
		createProductModel.IsActive)

	prd, err := parseProductModel(row)
	if err != nil {
		return ProductModel{}, err
	}
	return prd, nil
}

func (repo *ProductRepository) Update(id int, updateProductModel UpdateProductModel) error {
	query := `update models
	             set model_name = $1,
			         description = $2,
			         release_date = $3,
			         is_active = $4,
	           where id = $5 `

	result, err := repo.db.Conn().Exec(context.Background(), query,
		updateProductModel.ModelName,
		updateProductModel.Description,
		updateProductModel.ReleaseDate,
		updateProductModel.IsActive,
		id)

	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("Not found any row matches id %v", id)
	}

	return nil
}

func (repo *ProductRepository) Delete(id int) error {
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
