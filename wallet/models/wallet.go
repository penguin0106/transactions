package models

type Wallet struct {
	UserID   int      `json:"user_id"`
	Accounts []string `json:"accounts"`
}
