package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	ForeignKeyViolationErrorCode = pq.ErrorCode("23503")
	UniqueViolationErrorCode     = pq.ErrorCode("23505")
)

func Connect(user string, password string, dbName string, dbPort string) (*sqlx.DB, error) {
	// Disable SSL, not needed for this project
	connStr := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", user, password, dbPort, dbName)

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
	var count int

	err := db.Get(&count, "SELECT COUNT(*) FROM bank")
	if err != nil {
		return false, err
	}

	return count == 0, nil
}
