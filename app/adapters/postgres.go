package adapters

import (
	"database/sql"
	_ "github.com/lib/pq" //postgres driver
	"timekeeper/config"
)

func NewPostgresqlDatabase() (*sql.DB, error) {
	return sql.Open("postgres", config.DatabaseConnectionString())
}
