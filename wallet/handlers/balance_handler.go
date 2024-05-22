package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wallet/services"
)

// BalanceHandler handles the request for getting the wallet balance of a user.
func BalanceHandler(service *services.WalletService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from query parameters
		userIDStr := r.URL.Query().Get("user_id")
		if userIDStr == "" {
			http.Error(w, "Missing user ID", http.StatusBadRequest)
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Get the wallet balance using the service
		wallet, err := service.GetWalletByUserID(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Encode the wallet balance to JSON and write it to the response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(wallet); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
