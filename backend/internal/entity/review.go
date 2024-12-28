package entity

import "time"

type Review struct {
	Id         int       `json:"review_id" db:"ReviewID"`
	UserId     int       `json:"-" db:"UserId"`
	ItemId     int       `json:"itemId" db:"ItemID"`
	Name       string    `json:"name" db:"Name"`
	Rating     int       `json:"rating" db:"Rating" binding:"required"`
	Comment    string    `json:"comment" db:"Comment" binding:"required"`
	ReviewDate time.Time `json:"review_date" db:"ReviewDate"`
}
