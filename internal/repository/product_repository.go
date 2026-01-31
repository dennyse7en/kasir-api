package repository

import (
	"database/sql"
	"kasir-api/internal/model"
)

type ProductRepository interface {
	Create(product model.Product) (model.Product, error)
	GetAll() ([]model.Product, error)
	GetByID(id int) (model.Product, error)
	Update(id int, product model.Product) (model.Product, error)
	Delete(id int) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product model.Product) (model.Product, error) {
	query := `INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	if err != nil {
		return model.Product{}, err
	}
	return product, nil
}

func (r *productRepository) GetAll() ([]model.Product, error) {
	query := `SELECT id, name, price, stock, category_id FROM products`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *productRepository) GetByID(id int) (model.Product, error) {
	query := `SELECT id, name, price, stock, category_id FROM products WHERE id = $1`
	var p model.Product
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID)
	if err != nil {
		return model.Product{}, err
	}
	return p, nil
}

func (r *productRepository) Update(id int, product model.Product) (model.Product, error) {
	query := `UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5 RETURNING id, name, price, stock, category_id`
	var updatedProduct model.Product
	err := r.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID, id).Scan(&updatedProduct.ID, &updatedProduct.Name, &updatedProduct.Price, &updatedProduct.Stock, &updatedProduct.CategoryID)
	if err != nil {
		return model.Product{}, err
	}
	return updatedProduct, nil
}

func (r *productRepository) Delete(id int) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
