package entity

type AvailableProducts struct {
	//Id         int    `json:"id" db:"ID"`
	ProductId  string `json:"product_id" db:"Product_ID"`
	LocationId string `json:"location_id"`
	Location   string `json:"location"`
	Number     string `json:"number" db:"Number"`
}

type CreateAvailableProductsInput struct {
	ProductId  int `json:"product_id" `
	LocationId int `json:"location_id"`
	Number     int `json:"number"`
}
