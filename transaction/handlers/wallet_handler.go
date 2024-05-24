package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"transaction/services"
)

type WalletHandler struct {
	WalletService *services.WalletService
}

func NewWalletHandler(walletService *services.WalletService) *WalletHandler {
	return &WalletHandler{WalletService: walletService}
}

func (handler *WalletHandler) GetUserWallet(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "User ID parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID parameter", http.StatusBadRequest)
		return
	}

	wallet, err := handler.WalletService.GetUserWallet(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(wallet)
}

func (handler *WalletHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	var depositData struct {
		Amount        float64
		AccountNumber string
	}

	err := json.NewDecoder(r.Body).Decode(&depositData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = handler.WalletService.Deposit(r.Context(), depositData.Amount, depositData.AccountNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *WalletHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	var withdrawalData struct {
		Amount        float64
		AccountNumber string
	}

	err := json.NewDecoder(r.Body).Decode(&withdrawalData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = handler.WalletService.Withdraw(r.Context(), withdrawalData.Amount, withdrawalData.AccountNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *WalletHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var transferData struct {
		Amount                float64
		SenderAccountNumber   string
		ReceiverAccountNumber string
	}

	err := json.NewDecoder(r.Body).Decode(&transferData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = handler.WalletService.Transfer(r.Context(), transferData.Amount, transferData.SenderAccountNumber, transferData.ReceiverAccountNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
