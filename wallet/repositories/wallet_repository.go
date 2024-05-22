package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"wallet/models"
)

type WalletRepository struct {
	DB *sql.DB
}

func (repo *WalletRepository) GetWalletByUserID(ctx context.Context, userID int) (*models.Wallet, error) {
	// SQL query to fetch wallet data for the given user ID
	query := `SELECT user_id, usd, cryptocurrencies FROM wallets WHERE user_id = $1`

	// Execute the query
	row := repo.DB.QueryRowContext(ctx, query, userID)

	// Initialize variables to store the retrieved data
	var wallet models.Wallet
	var usd float64
	var cryptocurrenciesJSON []byte

	// Scan the row to extract data into variables
	err := row.Scan(&wallet.UserID, &usd, &cryptocurrenciesJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("wallet not found for the user")
		}
		return nil, err
	}

	// Unmarshal cryptocurrencies JSON into a map
	var cryptocurrencies map[string]float64
	if err := json.Unmarshal(cryptocurrenciesJSON, &cryptocurrencies); err != nil {
		return nil, err
	}

	// Assign retrieved data to the wallet struct
	wallet.USD = usd
	wallet.Cryptocurrencies = cryptocurrencies

	return &wallet, nil
}

func (repo *WalletRepository) Deposit(ctx context.Context, userID int, amount float64) error {
	// Start a transaction
	tx, err := repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		// Rollback the transaction if an error occurs
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// SQL query to update wallet balance
	query := `
        UPDATE wallets 
        SET usd = usd + $1 
        WHERE user_id = $2
    `

	// Execute the query
	_, err = tx.ExecContext(ctx, query, amount, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
