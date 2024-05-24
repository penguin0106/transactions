package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"wallet/services"
)

type CreateAccountRequest struct {
	UserID   int    `json:"user_id"`
	Currency string `json:"currency"`
}

// BalanceHandler handles the request for getting the wallet balance of a user.
func BalanceHandler(service *services.WalletService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(userID)
		if err != nil {
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}

		wallet, err := service.GetWalletByUserId(context.Background(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(wallet)
	}
}

type UpdateBalanceRequest struct {
	AccountNumber string  `json:"account_number"`
	Amount        float64 `json:"amount"`
}

func UpdateBalanceHandler(walletService *services.WalletService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateBalanceRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = walletService.Deposit(context.Background(), req.AccountNumber, req.Amount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

func CreateAccountHandler(walletService *services.WalletService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateAccountRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		account, err := walletService.CreateAccount(context.Background(), req.UserID, req.Currency)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(account)
	}
}
