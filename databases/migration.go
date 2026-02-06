package databases

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func Migrate(db *sql.DB) (*sql.DB, error) {
	query := `
	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description VARCHAR(255)
	);

	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description VARCHAR(255),
		price NUMERIC(10,2) NOT NULL,
		stock INT NOT NULL,
		category_id INT REFERENCES categories(id)
	);

	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		total_amount INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS transaction_details (
		id SERIAL PRIMARY KEY,
		transaction_id INT REFERENCES transactions(id) ON DELETE CASCADE,
		product_id INT REFERENCES products(id),
		quantity INT NOT NULL,
		subtotal INT NOT NULL
	);

	ALTER TABLE products ALTER COLUMN price TYPE integer USING price::integer;
	`
	_, err := db.Exec(query)

	if err != nil {
		return nil, err
	}

	log.Println("Migrated successfully.")
	return db, nil
}
