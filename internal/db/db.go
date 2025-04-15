package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(user string, password string, dbName string) (*sqlx.DB, error) {
	// Disable SSL, not needed for this project
	connStr := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", user, password, dbName)

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func IsEmpty(db *sqlx.DB) (bool, error) {
	rows, err := db.Query("SELECT COUNT(*) FROM bank, country")
	if err != nil {
		return false, err
	}

	var count int
	rows.Next()
	err = rows.Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}
