package entity

import (
	"time"
)

type Agreement struct {
	Id              int       `json:"-" db:"RentalAgreementID"`
	PaymentID       *int      `json:"payment_id" db:"PaymentID" binding:"required"`
	CustomerId      int       `json:"customerId" db:"CustomerId" binding:"required"` //TODO change name
	DateOfAgreement time.Time `json:"date" db:"DateOfAgreement" binding:"required"`
	Status          bool      `json:"status" db:"Status" binding:"required"`
}

type AgreementResponse struct {
	DateOfAgreement time.Time `json:"date"`
	ItemID          int       `json:"item_id"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	Status          bool      `json:"status"`
}
