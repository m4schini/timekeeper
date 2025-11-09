package adapters

import (
	"database/sql"
	"timekeeper/config"

	_ "github.com/lib/pq" //postgres driver
)

func NewPostgresqlDatabase() (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DatabaseConnectionString())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
