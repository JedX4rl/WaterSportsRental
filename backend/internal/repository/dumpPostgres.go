package repository

import (
	"WaterSportsRental/internal/entity"
	"context"
	"database/sql"
)

type DumpPostgres struct {
	db *sql.DB
}

func (d DumpPostgres) InsertDump(c context.Context, filePath string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	query := "SELECT insert_dump($1)"
	row := tx.QueryRowContext(c, query, filePath)
	if err := row.Err(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (d DumpPostgres) GetAllDumps(c context.Context) ([]entity.Dump, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	query := "SELECT * FROM get_all_dumps()"
	rows, err := tx.QueryContext(c, query)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	var dumps []entity.Dump
	for rows.Next() {
		var dump entity.Dump
		if err := rows.Scan(&dump.Filename); err != nil {
			tx.Rollback()
			return nil, err
		}
		dumps = append(dumps, dump)
	}
	return dumps, tx.Commit()
}

func NewDumpPostgres(db *sql.DB) *DumpPostgres {
	return &DumpPostgres{db: db}
}
