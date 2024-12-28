package repository

import (
	"WaterSportsRental/internal/entity"
	"context"
	"database/sql"
)

type LocationPostgres struct {
	db *sql.DB
}

func (l LocationPostgres) Create(c context.Context, location *entity.Location) error {
	tx, err := l.db.Begin()
	if err != nil {
		return err
	}
	query := "SELECT add_location($1, $2, $3, $4, $5, $6)"
	row := tx.QueryRowContext(c, query, location.Country, location.City, location.Address, location.OpeningTime, location.ClosingTime, location.PhoneNumber) //TODO time??
	if err = row.Err(); err != nil {
		tx.Rollback()
		return err
	}
	if err = row.Scan(&location.Id); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (l LocationPostgres) GetAll(c context.Context) ([]entity.Location, error) {
	tx, err := l.db.Begin()
	if err != nil {
		return nil, err
	}
	query := "SELECT * FROM get_all_locations()"
	rows, err := tx.QueryContext(c, query)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	var locations []entity.Location
	for rows.Next() {
		var location entity.Location
		if err = rows.Scan(&location.Id, &location.Country, &location.City, &location.Address, &location.OpeningTime, &location.ClosingTime, &location.PhoneNumber); err != nil {
			tx.Rollback()
			return nil, err
		}
		locations = append(locations, location)
	}
	return locations, tx.Commit()
}

func (l LocationPostgres) Update(c context.Context, location *entity.Location) error {
	tx, err := l.db.Begin()
	if err != nil {
		return err
	}
	query := "SELECT update_location($1, $2, $3, $4, $5, $6, $7)"
	row := tx.QueryRowContext(c, query, location.Id, location.Country, location.City, location.Address, location.OpeningTime, location.ClosingTime, location.PhoneNumber)
	if err = row.Err(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (l LocationPostgres) Delete(c context.Context, locationId int) error {
	tx, err := l.db.Begin()
	if err != nil {
		return err
	}
	query := "SELECT delete_location($1)"
	row := tx.QueryRowContext(c, query, locationId)
	if err = row.Err(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (l LocationPostgres) DeleteAll(c context.Context) error {
	tx, err := l.db.Begin()
	if err != nil {
		return err
	}
	query := "SELECT delete_all_locations()"
	row := tx.QueryRowContext(c, query)
	if err = row.Err(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func NewLocationPostgres(db *sql.DB) *LocationPostgres {
	return &LocationPostgres{db: db}
}
