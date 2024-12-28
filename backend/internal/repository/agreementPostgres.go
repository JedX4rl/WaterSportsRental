package repository

import (
	"WaterSportsRental/internal/entity"
	"database/sql"
	"golang.org/x/net/context"
	"time"
)

type AgreementPostgres struct {
	db *sql.DB
}

func (a AgreementPostgres) Create(c context.Context, userId int) (entity.Agreement, error) {
	tx, err := a.db.BeginTx(c, nil)
	if err != nil {
		return entity.Agreement{}, err
	}
	rentalAgreement := entity.Agreement{
		CustomerId:      userId,
		DateOfAgreement: time.Now(),
		Status:          false,
	}

	query := "SELECT create_rental_agreement($1, $2, $3)"

	err = tx.QueryRowContext(c, query, rentalAgreement.CustomerId, rentalAgreement.DateOfAgreement, rentalAgreement.Status).Scan(&rentalAgreement.Id)

	if err != nil {
		tx.Rollback()
		return entity.Agreement{}, err
	}

	if err = tx.Commit(); err != nil {
		return entity.Agreement{}, err
	}

	return rentalAgreement, nil

}

func (a AgreementPostgres) Delete(c context.Context, agreementId int) error {
	tx, err := a.db.BeginTx(c, nil)
	if err != nil {
		return err
	}
	query := "SELECT delete_rental_agreement($1)"

	_, err = tx.ExecContext(c, query, agreementId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()

}

func (a AgreementPostgres) GetAll(c context.Context, userId int) ([]entity.AgreementResponse, error) {
	tx, err := a.db.BeginTx(c, nil)
	if err != nil {
		return nil, err
	}

	query := "SELECT * FROM get_all_rental_agreements($1)"
	rows, err := tx.QueryContext(c, query, userId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	defer rows.Close()

	var agreements []entity.AgreementResponse

	for rows.Next() {
		var agreement entity.AgreementResponse
		err = rows.Scan(
			&agreement.DateOfAgreement,
			&agreement.ItemID,
			&agreement.StartDate,
			&agreement.EndDate,
			&agreement.Status,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		agreements = append(agreements, agreement)
	}
	if err = rows.Err(); err != nil {
		tx.Rollback()
		return nil, err
	}
	return agreements, tx.Commit()
}

func NewAgreementPostgres(db *sql.DB) *AgreementPostgres {
	return &AgreementPostgres{db: db}
}
