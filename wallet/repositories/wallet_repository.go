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

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{DB: db}
}

func (repo *WalletRepository) GetWalletByUserID(ctx context.Context, userID int) (*models.Wallet, error) {
	query := "SELECT user_id, usd, cryptocurrencies FROM wallets WHERE user_id = $1"

	row := repo.DB.QueryRowContext(ctx, query, userID)

	var wallet models.Wallet
	var usd float64
	var cryptocurrenciesJSON []byte

	err := row.Scan(&wallet.UserID, &usd, &cryptocurrenciesJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("wallet not found for the user")
		}
		return nil, err
	}

	var cryptocurrencies map[string]float64
	if err := json.Unmarshal(cryptocurrenciesJSON, &cryptocurrencies); err != nil {
		return nil, err
	}

	wallet.USD = usd
	wallet.Cryptocurrencies = cryptocurrencies

	return &wallet, nil
}

func (repo *WalletRepository) Deposit(ctx context.Context, userID int, amount float64) error {
	tx, err := repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query := `
        UPDATE wallets
        SET usd = usd + $1
        WHERE user_id = $2
    `

	_, err = tx.ExecContext(ctx, query, amount, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo *WalletRepository) UpdateWallet(ctx context.Context, wallet *models.Wallet) error {
	tx, err := repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	cryptocurrenciesJSON, err := json.Marshal(wallet.Cryptocurrencies)
	if err != nil {
		return err
	}

	query := `
        UPDATE wallets
        SET usd = $1, cryptocurrencies = $2
        WHERE user_id = $3
    `

	_, err = tx.ExecContext(ctx, query, wallet.USD, cryptocurrenciesJSON, wallet.UserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
