package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll() ([]models.Product, error) {
	query := "SELECT id, name, price, stock FROM products"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id"
	// QueryRow executes the SQL query with the provided parameters (product.Name, product.Price, product.Stock)
	// and expects at most one row in the result. It returns a *sql.Row which can be scanned.
	// Scan then copies the columns from the matched row into the destination (&product.ID),
	// populating the product's ID field with the database-generated value.
	// QueryRow is used instead of Query because we expect exactly one row (the inserted product's ID),
	// and QueryRow is more efficient for single-row results as it automatically closes the underlying
	// resources after Scan is called. The Product ID is also automatically filled in the variable
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock).Scan(&product.ID)
	return err
}

func (repo *ProductRepository) GetById(id string) (*models.Product, error) {
	query := "SELECT id, name, price, stock FROM products WHERE id = $1"
	var product models.Product
	err := repo.db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Price, &product.Stock)
	if err == sql.ErrNoRows {
		return nil, errors.New("Product Not Found")
	}
	if err != nil {
		return nil, err
	}
	return &product, err
}

func (repo *ProductRepository) DeleteById(id string) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("Product Not Found")
	}
	return err
}

func (repo *ProductRepository) UpdateProduct(product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3 WHERE id = $4"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Product Not Found")
	}
	return err
}
