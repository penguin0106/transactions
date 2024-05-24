package models

type Account struct {
	ID            int     `json:"id"`
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
	Active        bool    `json:"active"`
}
