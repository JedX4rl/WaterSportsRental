package entity

import "time"

type Equipment struct {
	Id    int     `json:"id" db:"ID"`
	Type  string  `json:"type" db:"type" binding:"required"`
	Brand string  `json:"brand" db:"brand" binding:"required"`
	Model string  `json:"model" db:"model" binding:"required"`
	Year  int     `json:"year" db:"year" binding:"required"`
	Price float64 `json:"price" db:"price" binding:"required"`
	Image string  `json:"image" db:"imageUrl"`
}

type EquipmentRequest struct {
	IDs           []int     `json:"ids" binding:"required"`
	LocationId    int       `json:"location_id" binding:"required"`
	StartDate     time.Time `json:"start_date" binding:"required"`
	EndDate       time.Time `json:"end_date" binding:"required"`
	PaymentMethod string    `json:"payment_method" binding:"required"`
}

type EquipmentResponse struct {
	Id         int       `json:"id" db:"ID"`
	Type       string    `json:"type" db:"type" binding:"required"`
	Brand      string    `json:"brand" db:"brand" binding:"required"`
	Model      string    `json:"model" db:"model" binding:"required"`
	Year       int       `json:"year" db:"year" binding:"required"`
	Image      string    `json:"image" db:"imageUrl"`
	TotalPrice float64   `json:"total_price" db:"total_price" binding:"required"`
	StartDate  time.Time `json:"start_date" binding:"required"`
	EndDate    time.Time `json:"end_date" binding:"required"`
}
