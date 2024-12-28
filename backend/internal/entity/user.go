package entity

import "time"

type User struct {
	Id               int       `json:"id" db:"ID"`
	Email            string    `json:"email" db:"email" binding:"required"`
	Password         string    `json:"password" db:"password" binding:"required"`
	FirstName        string    `json:"first_name" db:"FirstName"`
	LastName         string    `json:"last_name" db:"LastName"`
	Address          string    `json:"address" db:"Address"`
	PhoneNumber      string    `json:"phone_number" db:"PhoneNumber"`
	RegistrationDate time.Time `json:"registration_date" db:"RegistrationDate"`
	Role             string    `json:"role" db:"Role"`
}
