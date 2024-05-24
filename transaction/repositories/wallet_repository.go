package repositories

import (
	"context"
	"database/sql"
	"transaction/models"
)

type WalletRepository struct {
	DB *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{DB: db}
}

func (repo *WalletRepository) GetWalletByUserID(ctx context.Context, userID int) (*models.Wallet, error) {
	query := "SELECT accounts FROM wallets WHERE user_id = $1"
	row := repo.DB.QueryRowContext(ctx, query, userID)

	var wallet models.Wallet
	err := row.Scan(&wallet.Accounts)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	wallet.UserID = userID
	return &wallet, nil
}

func (repo *WalletRepository) Deposit(ctx context.Context, amount float64, accountNumber string) error {
	query := "UPDATE accounts SET balance = balance + $1 WHERE account_number = $2"
	_, err := repo.DB.ExecContext(ctx, query, amount, accountNumber)
	return err
}

func (repo *WalletRepository) Withdraw(ctx context.Context, amount float64, accountNumber string) error {
	query := "UPDATE accounts SET balance = balance - $1 WHERE account_number = $2"
	_, err := repo.DB.ExecContext(ctx, query, amount, accountNumber)
	return err
}

func (repo *WalletRepository) Transfer(ctx context.Context, amount float64, senderAccountNumber, receiverAccountNumber string) error {
	tx, err := repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Снять средства со счета отправителя
	query := "UPDATE accounts SET balance = balance - $1 WHERE account_number = $2"
	_, err = tx.ExecContext(ctx, query, amount, senderAccountNumber)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Зачислить средства на счет получателя
	query = "UPDATE accounts SET balance = balance + $1 WHERE account_number = $2"
	_, err = tx.ExecContext(ctx, query, amount, receiverAccountNumber)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
