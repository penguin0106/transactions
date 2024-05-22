package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"transaction/clients"
	"transaction/handlers"
	"transaction/repositories"
	"transaction/services"
)

func main() {
	// Database connection setup
	connStr := "user=username dbname=transactiondb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Repository and Service setup
	orderRepo := repositories.NewOrderRepository(db)
	walletClient := clients.NewWalletClient("http://wallet_service_url")
	orderService := services.NewOrderService(orderRepo, walletClient)

	// Handlers setup
	http.HandleFunc("/create_order", handlers.CreateOrderHandler(orderService))
	http.HandleFunc("/get_order", handlers.GetOrderHandler(orderService))
	http.HandleFunc("/purchase_order", handlers.PurchaseOrderHandler(orderService))

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":8081", nil))
}
