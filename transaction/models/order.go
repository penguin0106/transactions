package models

type Order struct {
	ID             int     `json:"id"`
	SellerID       int     `json:"seller_id"`
	BuyerID        int     `json:"buyer_id"`
	Cryptocurrency string  `json:"cryptocurrency"`
	Amount         float64 `json:"amount"`
	Price          float64 `json:"price"`
	Status         string  `json:"status"`
	ExchangeTo     string  `json:"exchange_to"`
}
