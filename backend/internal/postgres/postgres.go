package postgres

import (
	cfg "WaterSportsRental/internal/configs/storageConfig"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func NewPostgresDb(dbConfig cfg.StorageConfig) (*sql.DB, error) {

	const op = "internal.repository.postgres.NewPostgresDb"

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.DBName, dbConfig.Password, dbConfig.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}

func CloseDb(db *sql.DB) error {
	return db.Close()
}
