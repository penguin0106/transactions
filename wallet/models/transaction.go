package models

type Transaction struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	Amount    float64 `json:"amount"`
	Type      string  `json:"type"`
	Timestamp string  `json:"timestamp"`
}
