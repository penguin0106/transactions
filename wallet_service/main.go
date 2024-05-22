package main

import (
	"log"
	"net/http"
	"wallet/handlers"
	"wallet/repositories"
	"wallet/services"
)

func main() {

	walletRepo := &repositories.WalletRepository{DB: db}
	walletService := &services.WalletService{walletRepo}

	http.HandlerFunc("/balance", handlers.BalanceHandler(walletService))
	http.HandlerFunc("/deposit", handlers.DepositHandler(walletService))

	log.Fatal(http.ListenAndServe(":8081", nil))
}
