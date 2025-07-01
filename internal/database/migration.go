package database

import (
	"embed"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/epistax1s/gomer/internal/log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
)

//go:embed migrations/*.sql
var migrations embed.FS

func RunMigrations(dbURL string) error {
	entries, err := migrations.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Filter .up.sql migration files
	var migrationFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			// Check if the file name starts with a number (exampel 001)
			if !regexp.MustCompile(`^\d+_.+\.up\.sql$`).MatchString(entry.Name()) {
				return fmt.Errorf(
					"invalid migration file name: %s (must start with number, e.g., 001_name.up.sql)",
					entry.Name())
			}
			migrationFiles = append(migrationFiles, entry.Name())
		}
	}

	if len(migrationFiles) == 0 {
		return fmt.Errorf("no .up.sql migration files found in migrations directory")
	}

	sort.Strings(migrationFiles)

	source, err := bindata.WithInstance(bindata.Resource(
		migrationFiles,
		func(name string) ([]byte, error) {
			log.Info("Reading migration file: migrations/" + name)
			return migrations.ReadFile("migrations/" + name)
		},
	))
	if err != nil {
		return fmt.Errorf("failed to initialize migration source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance(
		"go-bindata",
		source,
		"sqlite3://"+dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Info("Migrations completed successfully")
	return nil
}
