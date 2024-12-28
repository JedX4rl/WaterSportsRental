package entity

import "time"

type RentedEquipment struct {
	Id                int       `json:"id" db:"RentedID"`
	RentalAgreementId int       `json:"rental_agreement_id" db:"RentalAgreementId" binding:"required"`
	ReviewId          int       `json:"review_id" db:"ReviewID"`
	InventoryID       int       `json:"inventory_id" db:"InventoryID" binding:"required"`
	StartDate         time.Time `json:"start_date" db:"StartDate" binding:"required"`
	EndDate           time.Time `json:"end_date" db:"EndDate" binding:"required"`
}

type AvailableDatesRequest struct {
	LocationId int `json:"location_id"`
	ItemId     int `json:"item_id" db:"LocationIDs"`
}

type AvailableDatesResponse struct {
	StartDate time.Time `json:"start_date" db:"StartDate" binding:"required"`
	EndDate   time.Time `json:"end_date" db:"EndDate" binding:"required"`
}
