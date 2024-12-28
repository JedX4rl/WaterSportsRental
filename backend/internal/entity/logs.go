package entity

import "time"

type Logs struct {
	ID        int       `json:"id"`
	ItemId    int       `json:"itemId"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}
