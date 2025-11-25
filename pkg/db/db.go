package db

import (
	"avito-intern-test-task-2025/internal/config"
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	Pool  *pgxpool.Pool
	dbUrl string
}

func InitDB(cfg *config.Config) (*DB, error) {

	dbUrl := "postgresql://" + cfg.DB_USER + ":" + cfg.DB_PASSWORD + "@" + cfg.DB_HOST + ":" + cfg.DB_PORT + "/" + cfg.DB_NAME

	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatal("Unable to connect to database: %v\n", err)
	}

	return &DB{Pool: pool, dbUrl: dbUrl}, err
}

func (db *DB) RunMigrations() {

	driver, err := postgres.WithInstance(sql.OpenDB(stdlib.GetPoolConnector(db.Pool)), &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to initialize migrate driver: %v", err)
	}

	_, filename, _, _ := runtime.Caller(0)
	migrationsPath := filepath.Join(filepath.Dir(filename), "migrations")

	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		if cwd, cwdErr := os.Getwd(); cwdErr == nil {
			altPath := filepath.Join(cwd, "pkg", "db", "migrations")
			if _, altErr := os.Stat(altPath); altErr == nil {
				migrationsPath = altPath
			} else {
				log.Fatalf("Migrations directory not found (checked %s and %s)", migrationsPath, altPath)
			}
		} else {
			log.Fatalf("Unable to determine working directory for migrations: %v", cwdErr)
		}
	}

	migrationsURL := "file://" + filepath.ToSlash(migrationsPath)

	m, err := migrate.NewWithDatabaseInstance(migrationsURL, "postgres", driver)
	if err != nil {
		log.Fatalf("Migration setup failed: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("No new migrations to apply")
	} else {
		version, dirty, _ := m.Version()
		log.Printf("Migrations applied successfully! Current version: %d (dirty: %v)", version, dirty)
	}
}

func (db *DB) Close() {
	db.Pool.Close()
}
