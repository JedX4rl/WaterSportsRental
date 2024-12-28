package entity

import "time"

type Payment struct {
	Id            int       `json:"id" db:"PaymentID"`
	Amount        int       `json:"amount" db:"Amount"`
	PaymentDate   time.Time `json:"payment_date" db:"PaymentDate"`
	PaymentMethod string    `json:"payment_method" db:"PaymentMethod"`
	Status        bool      `json:"status" db:"Status"`
}
