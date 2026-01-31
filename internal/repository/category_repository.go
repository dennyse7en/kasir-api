package repository

import (
	"database/sql"
	"kasir-api/internal/model"
)

type CategoryRepository interface {
	Create(category model.Category) (model.Category, error)
	GetAll() ([]model.Category, error)
	GetByID(id int) (model.Category, error)
	Update(id int, category model.Category) (model.Category, error)
	Delete(id int) error
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category model.Category) (model.Category, error) {
	query := `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(query, category.Name, category.Description).Scan(&category.ID)
	if err != nil {
		return model.Category{}, err
	}
	return category, nil
}

func (r *categoryRepository) GetAll() ([]model.Category, error) {
	query := `SELECT id, name, description FROM categories`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *categoryRepository) GetByID(id int) (model.Category, error) {
	query := `SELECT id, name, description FROM categories WHERE id = $1`
	var c model.Category
	err := r.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		return model.Category{}, err
	}
	return c, nil
}

func (r *categoryRepository) Update(id int, category model.Category) (model.Category, error) {
	query := `UPDATE categories SET name = $1, description = $2 WHERE id = $3 RETURNING id, name, description`
	var updatedCategory model.Category
	err := r.db.QueryRow(query, category.Name, category.Description, id).Scan(&updatedCategory.ID, &updatedCategory.Name, &updatedCategory.Description)
	if err != nil {
		return model.Category{}, err
	}
	return updatedCategory, nil
}

func (r *categoryRepository) Delete(id int) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
