package handlers

import (
	"encoding/json"
	"net/http"
	"wallet/services"
)

type DepositRequest struct {
	UserID int     `json:"user_id"`
	Amount float64 `json:"amount"`
}

func DepositHandler(service *services.WalletService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DepositRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Amount <= 0 {
			http.Error(w, "Invalid deposit amount", http.StatusBadRequest)
			return
		}

		err := service.Deposit(r.Context(), req.UserID, req.Amount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}
