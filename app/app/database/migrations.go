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

const LatestSchema = uint(0)

// InitSchema migrates the database to the specified version (if version = 0, migrates to latest version)
func InitSchema(db *sql.DB, version uint) (bool, error) {
	log := zap.L().Named("database")
	log.Info("syncing database schema")

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
	v, dirty, err := m.Version()
	log.Info("database schema version", zap.Uint("version", v), zap.Bool("dirty", dirty), zap.Error(err))

	if version <= 0 {
		err = m.Up()
	} else {
		err = m.Migrate(version)
	}
	if err != nil {
		log.Warn("migration failed", zap.Error(err))
		return false, nil
	}

	v, dirty, err = m.Version()
	log.Info("migrated database schema", zap.Uint("version", v), zap.Bool("dirty", dirty), zap.Error(err))

	return true, nil
}
