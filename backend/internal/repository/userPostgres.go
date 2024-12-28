package repository

import (
	"WaterSportsRental/internal/entity"
	"context"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func (u UserRepository) Create(c context.Context, user *entity.User) error {
	tx, err := u.db.BeginTx(c, nil)
	if err != nil {
		return err
	}

	query := "SELECT create_user($1, $2, $3)"
	row := tx.QueryRowContext(c, query, user.Email, user.Password, user.RegistrationDate)
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

func (u UserRepository) Get(c context.Context, userId int) (*entity.User, error) {
	tx, err := u.db.BeginTx(c, nil)
	if err != nil {
		return nil, err
	}
	query := "SELECT * FROM get_user($1)"
	row := tx.QueryRowContext(c, query, userId)
	if err := row.Err(); err != nil {
		tx.Rollback()
		return nil, err
	}
	var user entity.User
	var fName, lName, Addr, pNumber, role, mail sql.NullString
	var regDate sql.NullTime

	if err = row.Scan(&user.Id, &mail, &fName, &lName, &Addr, &pNumber, &regDate, &role); err != nil {
		tx.Rollback()
		return nil, err
	}
	user.FirstName = fName.String
	user.LastName = lName.String
	user.Address = Addr.String
	user.PhoneNumber = pNumber.String
	user.RegistrationDate = regDate.Time
	user.Role = role.String
	user.Email = mail.String

	return &user, tx.Commit()
}

func (u UserRepository) UpdateWithRole(c context.Context, user *entity.User) error {
	tx, err := u.db.BeginTx(c, nil)
	if err != nil {
		return err
	}
	query := "SELECT update_user($1, $2, $3, $4, $5, $6)"
	row := tx.QueryRowContext(c, query, user.Id, user.FirstName, user.LastName, user.Address, user.PhoneNumber, user.Role)
	if err := row.Err(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (u UserRepository) GiveRole(c context.Context, user *entity.User) error {
	tx, err := u.db.BeginTx(c, nil)
	if err != nil {
		return err
	}
	query := "SELECT give_role($1, $2, $3, $4)"
	row := tx.QueryRowContext(c, query, user.Id, user.Role)
	if err := row.Err(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (u UserRepository) GetAll(c context.Context) ([]entity.User, error) {

	tx, err := u.db.BeginTx(c, nil)
	if err != nil {
		return nil, err
	}
	query := "SELECT * FROM get_all_users()"
	rows, err := tx.QueryContext(c, query)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	var users []entity.User
	defer rows.Close()
	for rows.Next() {
		var user entity.User
		var fName, lName, Addr, pNumber, role, mail sql.NullString
		var regDate sql.NullTime

		if err = rows.Scan(&user.Id, &mail, &fName, &lName, &Addr, &pNumber, &regDate, &role); err != nil {
			tx.Rollback()
			return nil, err
		}
		user.FirstName = fName.String
		user.LastName = lName.String
		user.Address = Addr.String
		user.PhoneNumber = pNumber.String
		user.RegistrationDate = regDate.Time
		user.Role = role.String
		user.Email = mail.String
		users = append(users, user)
	}
	return users, tx.Commit()
}

func (u UserRepository) Delete(c context.Context, userId int) error {
	tx, err := u.db.BeginTx(c, nil)
	if err != nil {
		return err
	}
	query := "SELECT delete_user($1)"
	row := tx.QueryRowContext(c, query, userId)
	if err := row.Err(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (u UserRepository) GetAllAgreementsById(c context.Context, agreementId int) ([]entity.AgreementResponse, error) {
	tx, err := u.db.BeginTx(c, nil)
	if err != nil {
		return nil, err
	}
	query := "SELECT * FROM get_all_rental_agreements($1)"
	rows, err := tx.QueryContext(c, query, agreementId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	var response []entity.AgreementResponse
	for rows.Next() {
		var agreement entity.AgreementResponse
		if err := rows.Scan(&agreement.DateOfAgreement, &agreement.ItemID, &agreement.StartDate, &agreement.EndDate, &agreement.Status); err != nil {
			tx.Rollback()
			return nil, err
		}
		response = append(response, agreement)
	}
	return response, tx.Commit()
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
