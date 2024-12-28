package repository

import (
	"WaterSportsRental/internal/entity"
	"database/sql"
	"golang.org/x/net/context"
	"log/slog"
)

type AuthPostgres struct {
	db *sql.DB
}

func (a AuthPostgres) Create(ctx context.Context, user *entity.User) error {

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := "SELECT create_user($1, $2, $3)"
	row := tx.QueryRowContext(ctx, query, user.Email, user.Password, user.RegistrationDate)
	if err := row.Err(); err != nil {
		tx.Rollback()
		return err
	}
	if err := row.Scan(&user.Id); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (a AuthPostgres) GetByEmail(ctx context.Context, email string) (entity.User, error) {

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return entity.User{}, err
	}

	query := "SELECT * FROM get_user_by_email($1)"
	row := tx.QueryRowContext(ctx, query, email)

	if err := row.Err(); err != nil {
		tx.Rollback()
		return entity.User{}, err
	}

	var user entity.User
	var fName, lName, Addr, pNumber, role, mail, pwd sql.NullString
	var regDate sql.NullTime

	err = row.Scan(&user.Id, &fName, &lName, &Addr, &pNumber, &regDate, &role, &mail, &pwd)

	if err != nil {
		slog.Info(err.Error())
		tx.Rollback()
		return entity.User{}, err
	}

	user.FirstName = fName.String
	user.LastName = lName.String
	user.Address = Addr.String
	user.PhoneNumber = pNumber.String
	user.RegistrationDate = regDate.Time
	user.Role = role.String
	user.Email = mail.String
	user.Password = pwd.String

	return user, tx.Commit()

}

func (a AuthPostgres) GetById(c context.Context, id int) (entity.User, error) {
	tx, err := a.db.BeginTx(c, nil)
	if err != nil {
		return entity.User{}, err
	}
	query := "SELECT * FROM get_user_by_id($1)"
	row := tx.QueryRowContext(c, query, id)
	if err := row.Err(); err != nil {
		tx.Rollback()
		return entity.User{}, err
	}
	var user entity.User
	var fName, lName, Addr, pNumber, role, mail, pwd sql.NullString
	var regDate sql.NullTime

	err = row.Scan(&user.Id, &fName, &lName, &Addr, &pNumber, &regDate, &role, &mail, &pwd) //TODO maybe no need in password?

	if err != nil {
		tx.Rollback()
		return entity.User{}, err
	}
	user.FirstName = fName.String
	user.LastName = lName.String
	user.Address = Addr.String
	user.PhoneNumber = pNumber.String
	user.RegistrationDate = regDate.Time
	user.Role = role.String
	user.Email = mail.String
	user.Password = pwd.String

	return user, tx.Commit()
}
func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{
		db: db,
	}
}
