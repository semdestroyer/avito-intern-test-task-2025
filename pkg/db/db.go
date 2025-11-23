package db

import (
	"avito-intern-test-task-2025/internal/config"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type DB struct {
	Pool *pgxpool.Pool
}

func InitDB(cfg *config.Config) (*DB, error) {
	dbUrl := "postgresql://" + cfg.DB_USERNAME + ":" + cfg.DB_PASSWORD + "@" + cfg.DB_HOST + ":" + cfg.DB_PORT + "/" + cfg.DB_NAME
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatal("Unable to connect to database: %v\n", err)
	}

	return &DB{Pool: pool}, err
}

func (db *DB) Close() {
	db.Pool.Close()
}
