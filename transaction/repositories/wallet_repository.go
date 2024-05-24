package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"transaction/models"
)

type WalletRepository struct {
	DB *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{DB: db}
}

func (repo *WalletRepository) GetWalletByUserID(ctx context.Context, userID int) (*models.Wallet, error) {
	query := "SELECT user_id, accounts FROM wallets WHERE user_id = $1"
	row := repo.DB.QueryRowContext(ctx, query, userID)

	var wallet models.Wallet
	var accountsJSON []byte

	err := row.Scan(&wallet.UserID, &accountsJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var accounts []string
	if err := json.Unmarshal(accountsJSON, &accounts); err != nil {
		return nil, err
	}

	wallet.Accounts = accounts

	return &wallet, nil
}

func (repo *WalletRepository) Deposit(ctx context.Context, userID int, accountNumber string, amount float64) error {
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
			SET accounts = array_append(accounts, $1)
			WHERE user_id = $2
	`

	_, err = tx.ExecContext(ctx, query, accountNumber, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
