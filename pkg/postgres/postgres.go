package postgres

import (
	"database/sql"
	"fmt"
	"people-finder/config"

	_ "github.com/lib/pq"
)

func New(sc config.StorageConfig) (*sql.DB, error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open("postgres", sc.URL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db, nil
}
