package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"wallet/models"
)

type AccountRepository struct {
	DB *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{DB: db}
}

func (repo *AccountRepository) CreateAccount(ctx context.Context, currency string) (*models.Account, error) {
	accountNumber := generateAccountNumber(currency)
	query := "INSERT INTO accounts (account_number, balance, active) VALUES ($1, $2, $3) RETURNING id, account_number, balance, active"

	var account models.Account
	err := repo.DB.QueryRowContext(ctx, query, accountNumber, 0.00, true).Scan(&account.ID, &account.AccountNumber, &account.Balance, &account.Active)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (repo *AccountRepository) GetAccountByNumber(ctx context.Context, accountNumber string) (*models.Account, error) {
	query := "SELECT id, account_number, balance, active FROM accounts WHERE account_number = $1"

	var account models.Account
	err := repo.DB.QueryRowContext(ctx, query, accountNumber).Scan(&account.ID, &account.AccountNumber, &account.Balance, &account.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, err
	}

	return &account, nil
}

func (repo *AccountRepository) UpdateAccount(ctx context.Context, account *models.Account) error {
	query := "UPDATE accounts SET balance = $1, active = $2 WHERE account_number  = $3"
	_, err := repo.DB.ExecContext(ctx, query, account.Balance, account.Active, account.AccountNumber)
	return err
}

func generateAccountNumber(currency string) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 29)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return currency + string(b)
}
