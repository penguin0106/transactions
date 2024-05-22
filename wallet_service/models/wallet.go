package models

type Wallet struct {
	UserID           int                `json:"user_id"`
	USD              float64            `json:"usd"`
	Cryptocurrencies map[string]float64 `json:"cryptocurrencies"`
}
