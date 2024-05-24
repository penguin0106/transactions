package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"wallet/models"
)

type WalletRepository struct {
	DB *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{DB: db}
}

func (repo *WalletRepository) GetWalletByUserID(ctx context.Context, userID int) (*models.Wallet, error) {
	query := "SELECT accounts FROM wallets WHERE user_id = $1"

	var accountsJSON []byte
	err := repo.DB.QueryRowContext(ctx, query, userID).Scan(&accountsJSON)
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

	return &models.Wallet{
		UserID:   userID,
		Accounts: accounts,
	}, nil
}

func (repo *WalletRepository) CreateWallet(ctx context.Context, userID int, usdAccount string) error {
	accounts := []string{usdAccount}
	accountsJson, err := json.Marshal(accounts)
	if err != nil {
		return err
	}

	query := "INSERT INTO wallets (user_id, accounts) VALUES ($1, $2)"
	_, err = repo.DB.ExecContext(ctx, query, userID, accountsJson)
	return err
}

func (repo *WalletRepository) AddAccountToWallet(ctx context.Context, userID int, accountNumber string) error {
	wallet, err := repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if wallet == nil {
		return sql.ErrNoRows
	}

	wallet.Accounts = append(wallet.Accounts, accountNumber)
	accountsJSON, err := json.Marshal(wallet.Accounts)
	if err != nil {
		return err
	}

	query := "UPDATE wallets SET accounts = $1 WHERE user_id = $2"
	_, err = repo.DB.ExecContext(ctx, query, accountsJSON, userID)
	return err
}

func (repo *WalletRepository) Deposit(ctx context.Context, accountNumber string, amount float64) error {
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
			UPDATE accounts
			SET balance = balance + $1
			WHERE account_number = $2
	`

	_, err = tx.ExecContext(ctx, query, amount, accountNumber)
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

	accountsJSON, err := json.Marshal(wallet.Accounts)
	if err != nil {
		return err
	}

	query := `
			UPDATE wallets
			SET accounts = $1
			WHERE user_id = $2
	`

	_, err = tx.ExecContext(ctx, query, accountsJSON, wallet.UserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
