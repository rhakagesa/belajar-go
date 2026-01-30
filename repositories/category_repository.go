package repositories

import (
	"belajar-go/models"
	"database/sql"
	"errors"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (repo *CategoryRepository) GetAll(page int, limit int) ([]models.Category, int, error) {
	offset := (page - 1) * limit

	var totalData int
	err := repo.db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&totalData)

	if err != nil {
		return nil, 0, err
	}

	rows, err := repo.db.Query("SELECT id, name, description FROM categories LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description)
		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, category)
	}

	return categories, totalData, nil
}

func (repo *CategoryRepository) GetByID(id int) (*models.Category, error) {
	var category models.Category
	err := repo.db.QueryRow("SELECT id, name, description FROM categories WHERE id = $1", id).Scan(&category.ID, &category.Name, &category.Description)

	if err == sql.ErrNoRows {
		return nil, errors.New("Category not found.")
	}

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (repo *CategoryRepository) Create(category *models.Category) error {
	err := repo.db.QueryRow("INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id",
		category.Name, category.Description).Scan(&category.ID)

	return err
}

func (repo *CategoryRepository) Update(category *models.Category) error {
	result, err := repo.db.Exec("UPDATE categories SET name = $1, description = $2 WHERE id = $3",
		category.Name, category.Description, category.ID)

	row, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if row == 0 {
		return errors.New("Category not found.")
	}

	return nil
}

func (repo *CategoryRepository) Delete(id int) error {
	result, err := repo.db.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}

	row, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if row == 0 {
		return errors.New("Category not found.")
	}

	return nil
}
