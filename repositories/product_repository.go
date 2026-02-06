package repositories

import (
	"belajar-go/models"
	"database/sql"
	"errors"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll(page int, limit int, name string) ([]models.Product, int, error) {
	offset := (page - 1) * limit

	var totalData int
	err := repo.db.QueryRow("SELECT COUNT(*) FROM products").Scan(&totalData)

	if err != nil {
		return nil, 0, err
	}

	where := ""
	args := []any{limit, offset}
	if name != "" {
		where += " WHERE p.name ILIKE $3"
		args = append(args, "%"+name+"%")
	}

	query := "SELECT p.id, p.name, p.description, p.price, p.stock, c.id, c.name, c.description FROM products p JOIN categories c ON c.id = p.category_id " + where + " LIMIT $1 OFFSET $2"

	rows, err := repo.db.Query(query, args...)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		var category models.Category
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock,
			&category.ID, &category.Name, &category.Description)
		if err != nil {
			return nil, 0, err
		}
		product.Category = &category
		products = append(products, product)
	}

	if len(products) == 0 {
		return nil, 0, nil
	}

	return products, totalData, nil
}

func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	var product models.Product
	var category models.Category
	err := repo.db.QueryRow(`
	SELECT p.id, p.name, p.description, p.price, p.stock, c.id, c.name, c.description 
	FROM products p 
	JOIN categories c ON c.id = p.category_id 
	WHERE p.id = $1`, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock,
		&category.ID, &category.Name, &category.Description)

	if err == sql.ErrNoRows {
		return nil, errors.New("Product not found.")
	}

	if err != nil {
		return nil, err
	}

	product.Category = &category
	return &product, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	err := repo.db.QueryRow("INSERT INTO products (name, description, price, stock, category_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		product.Name, product.Description, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)

	return err
}

func (repo *ProductRepository) Update(product *models.Product) error {
	result, err := repo.db.Exec("UPDATE products SET name = $1, description = $2, price = $3, stock = $4, category_id = $5 WHERE id = $6",
		product.Name, product.Description, product.Price, product.Stock, product.CategoryID, product.ID)

	row, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if row == 0 {
		return errors.New("Product not found.")
	}

	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	result, err := repo.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	row, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if row == 0 {
		return errors.New("Product not found.")
	}

	return nil
}
