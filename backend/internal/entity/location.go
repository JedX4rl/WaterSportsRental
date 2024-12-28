package entity

type Location struct {
	Id          int     `json:"id" db:"LocationID"`
	Country     *string `json:"country" db:"Country" validate:"required"`
	City        *string `json:"city" db:"City" validate:"required"`
	Address     *string `json:"address" db:"Address" validate:"required"`
	OpeningTime *string `json:"opening_time" db:"OpeningTime" validate:"required"`
	ClosingTime *string `json:"closing_time" db:"ClosingTime" validate:"required"`
	PhoneNumber *string `json:"phone_number" db:"PhoneNumber" validate:"required"`
}
