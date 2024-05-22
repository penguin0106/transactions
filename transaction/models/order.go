package models

import "time"

type Order struct {
	ID              int       `json:"id"`
	SellerID        int       `json:"seller_id"`
	BuyerID         int       `json:"buyer_id"`
	Cryptocurrency  string    `json:"cryptocurrency"`
	Amount          float64   `json:"amount"`
	DesiredCurrency string    `json:"desired_currency"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}
