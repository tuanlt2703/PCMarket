package repositories

import (
	"database/sql"
	"fmt"

	"product-service/internal/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) Create(req models.CreateProductRequest) (*models.Product, error) {

	query := `
	INSERT INTO products(name, description, price, quantity, created_at, updated_at)
	VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	result, err := r.db.Exec(
		query,
		req.Name,
		req.Description,
		req.Price,
		req.Quantity,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetByID(int(id))
}

func (r *ProductRepository) GetAll(keyword string) ([]models.Product, error) {

	var rows *sql.Rows
	var err error

	if keyword != "" {
		query := `
		SELECT id, name, description, price, quantity, created_at, updated_at
		FROM products
		WHERE name LIKE ?
		ORDER BY id DESC
		`

		rows, err = r.db.Query(query, "%"+keyword+"%")
	} else {

		query := `
		SELECT id, name, description, price, quantity, created_at, updated_at
		FROM products
		ORDER BY id DESC
		`

		rows, err = r.db.Query(query)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []models.Product

	for rows.Next() {

		var product models.Product

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Quantity,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) GetByID(id int) (*models.Product, error) {

	query := `
	SELECT id, name, description, price, quantity, created_at, updated_at
	FROM products
	WHERE id = ?
	`

	var product models.Product

	err := r.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) Update(id int, req models.UpdateProductRequest) (*models.Product, error) {

	query := `
	UPDATE products
	SET
		name = ?,
		description = ?,
		price = ?,
		quantity = ?,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`

	result, err := r.db.Exec(
		query,
		req.Name,
		req.Description,
		req.Price,
		req.Quantity,
		id,
	)

	if err != nil {
		return nil, err
	}

	rows, _ := result.RowsAffected()

	if rows == 0 {
		return nil, fmt.Errorf("product not found")
	}

	return r.GetByID(id)
}

func (r *ProductRepository) Delete(id int) error {

	query := `
	DELETE FROM products
	WHERE id = ?
	`

	result, err := r.db.Exec(query, id)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()

	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
