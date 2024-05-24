package services

import (
	"context"
	"errors"
	"transaction/models"
	"transaction/repositories"
)

type WalletService struct {
	repo *repositories.WalletRepository
}

func NewWalletService(repo *repositories.WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

func (service *WalletService) GetUserWallet(ctx context.Context, userID int) (*models.Wallet, error) {
	wallet, err := service.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (service *WalletService) Deposit(ctx context.Context, amount float64, accountNumber string) error {
	// Проверяем, что сумма депозита положительная
	if amount <= 0 {
		return errors.New("deposit amount must be positive")
	}
	// Вызываем метод репозитория для выполнения операции депозита
	err := service.repo.Deposit(ctx, amount, accountNumber)
	if err != nil {
		return err
	}
	return nil
}

func (service *WalletService) Withdraw(ctx context.Context, amount float64, accountNumber string) error {
	// Проверяем, что сумма снятия положительная
	if amount <= 0 {
		return errors.New("withdrawal amount must be positive")
	}
	// Вызываем метод репозитория для выполнения операции снятия
	err := service.repo.Withdraw(ctx, amount, accountNumber)
	if err != nil {
		return err
	}
	return nil
}

func (service *WalletService) Transfer(ctx context.Context, amount float64, senderAccountNumber, receiverAccountNumber string) error {
	err := service.repo.Transfer(ctx, amount, senderAccountNumber, receiverAccountNumber)
	if err != nil {
		return err
	}
	return nil
}
