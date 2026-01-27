package database

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"

	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//go:embed migrations
var fs embed.FS

func InitSchema(db *sql.DB) (bool, error) {
	log := zap.L().Named("database")

	migrationFs, err := iofs.New(fs, "migrations")
	if err != nil {
		return false, err
	}
	defer migrationFs.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return false, err
	}

	m, err := migrate.NewWithInstance(
		"iofs", migrationFs,
		"postgres", driver)
	if err != nil {
		return false, err
	}

	err = m.Up()
	if err != nil {
		return false, nil
	}

	version, dirty, err := m.Version()
	log.Info("upgraded database schema", zap.Uint("version", version), zap.Bool("dirty", dirty), zap.Error(err))

	return true, nil
}
