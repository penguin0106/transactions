package wallet

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"wallet/handlers"
	"wallet/repositories"
	"wallet/services"
)

func main() {
	// Database connection setup
	connStr := "user=username dbname=walletdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Repository and Service setup
	walletRepo := &repositories.WalletRepository{DB: db}
	walletService := services.NewWalletService(walletRepo)

	// Handlers setup
	http.HandleFunc("/balance", handlers.BalanceHandler(walletService))
	http.HandleFunc("/deposit", handlers.DepositHandler(walletService))

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
