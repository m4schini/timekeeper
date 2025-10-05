package query

import "database/sql"

type Queries struct {
	DB *sql.DB
}

func NewQueries(db *sql.DB) *Queries {
	return &Queries{DB: db}
}
