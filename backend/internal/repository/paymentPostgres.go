package repository

import (
	"WaterSportsRental/internal/entity"
	"database/sql"
	"golang.org/x/net/context"
)

type PaymentPostgres struct {
	db *sql.DB
}

func (p PaymentPostgres) Create(c context.Context, payment *entity.Payment) error {
	return nil
}

func NewPaymentPostgres(db *sql.DB) *PaymentPostgres {
	return &PaymentPostgres{db: db}
}
