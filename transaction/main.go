package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"
	"transaction/handlers"
	"transaction/repositories"
	"transaction/services"
)

func main() {
	// Database connection setup
	connStr := "user=username dbname=walletdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	// Инициализация репозиториев
	walletRepo := repositories.NewWalletRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	// Инициализация сервисов
	walletService := services.NewWalletService(walletRepo)
	orderService := services.NewOrderService(orderRepo, walletService)

	// Инициализация хендлеров
	walletHandler := handlers.NewWalletHandler(walletService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Настройка маршрутов
	http.HandleFunc("/wallet", walletHandler.GetUserWallet)
	http.HandleFunc("/wallet/deposit", walletHandler.Deposit)
	http.HandleFunc("/wallet/withdraw", walletHandler.Withdraw)
	http.HandleFunc("/wallet/transfer", walletHandler.Transfer)

	http.HandleFunc("/orders", orderHandler.FindOrders)
	http.HandleFunc("/orders/create", orderHandler.CreateOrder)
	http.HandleFunc("/orders/by-currency", orderHandler.FindOrdersByCurrency)
	http.HandleFunc("/orders/by-seller", orderHandler.FindOrdersBySellerUsername)
	http.HandleFunc("/orders/purchase", orderHandler.PurchaseOrder)

	// Запуск HTTP-сервера
	server := &http.Server{
		Addr:         ":8080",
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on :8080: %v\n", err)
	}
}
