package database

import (
	"database/sql"
	"raumzeitalpaka/app/database/command"
	"raumzeitalpaka/app/database/query"
)

type Database struct {
	Queries  query.Queries
	Commands command.Commands
}

func New(db *sql.DB) *Database {
	return &Database{
		Queries:  query.NewQueries(db),
		Commands: command.NewCommands(db),
	}
}
