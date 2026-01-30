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
	`
	_, err := db.Exec(query)

	if err != nil {
		return nil, err
	}

	log.Println("Migrated successfully.")
	return db, nil
}
