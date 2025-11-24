package db

import (
	"avito-intern-test-task-2025/internal/config"
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"log"
	"path/filepath"
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

// TODO: провести обязательный рефакторинг этого места после того как сделаю базу. + пробросить пути и переместить migrations в корень
func (db *DB) RunMigrations() {

	driver, err := postgres.WithInstance(sql.OpenDB(stdlib.GetPoolConnector(db.Pool)), &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to initialize migrate driver: %v", err)
	}

	//migrationsRelPath := "D:\\GolangProjects\\avito-intern-test-task-2025\\pkg\\db\\migrations"
	migrationsRelPath := "./migrations"
	absPath, err := filepath.Abs(migrationsRelPath)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return
	}

	urlPath := filepath.ToSlash(absPath)

	m, err := migrate.NewWithDatabaseInstance("file://"+urlPath, "postgres", driver)
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
