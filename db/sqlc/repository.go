package db

import (
	"database/sql"
)

// Repository provides all functions to execute db queries
type Repository interface {
	Querier
}

// SQLRepository provides all functions to execute SQL queries
type SQLRepository struct {
	db *sql.DB
	*Queries
}

// NewRepository creates a new Repository
func NewRepository(db *sql.DB) Repository {
	return &SQLRepository{
		db:      db,
		Queries: New(db),
	}
}
