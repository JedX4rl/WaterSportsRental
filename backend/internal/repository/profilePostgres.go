package repository

import (
	"WaterSportsRental/internal/entity"
	"database/sql"
	"golang.org/x/net/context"
	"log/slog"
)

type ProfilePostgres struct {
	db *sql.DB
}

func (p ProfilePostgres) Update(c context.Context, user *entity.User) error {
	tx, err := p.db.BeginTx(c, nil)
	if err != nil {
		return err
	}

	query := `SELECT update_user($1, $2, $3, $4, $5)`
	row := tx.QueryRowContext(c, query, user.Id, user.FirstName, user.LastName, user.Address, user.PhoneNumber)
	if err := row.Err(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (p ProfilePostgres) Get(c context.Context, userId int) (*entity.User, error) {
	tx, err := p.db.BeginTx(c, nil)
	if err != nil {
		return nil, err
	}

	query := `SELECT * FROM get_user($1)`
	row := p.db.QueryRowContext(c, query, userId)
	if err := row.Err(); err != nil {
		tx.Rollback()
		return nil, err
	}

	var user entity.User

	if err = row.Scan(&user.Id,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Address,
		&user.PhoneNumber,
		&user.RegistrationDate,
		&user.Role); err != nil {
		tx.Rollback()
		slog.Info(err.Error())
		return nil, err
	}
	user.Password = ""

	return &user, tx.Commit()
}

func NewProfilePostgres(db *sql.DB) *ProfilePostgres {
	return &ProfilePostgres{
		db: db,
	}
}
