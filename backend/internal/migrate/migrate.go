package migrate

import (
	"backend/internal/config"
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed sql/*.sql
var embedMigrations embed.FS

func RunMigrations(cfg *config.Config, db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(db, "sql"); err != nil && err != goose.ErrNoMigrationFiles {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
