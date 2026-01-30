package databases

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	//Buka Koneksi ke DB
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	//Test Koneksi
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	//Setting Connection Pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connected successfully")
	return db, nil
}
