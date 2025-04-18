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

const Port = "5432"

func Connect(user string, password string, dbName string, host string, port string) (*sqlx.DB, error) {
	// Disable SSL, not needed for this project
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)

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
