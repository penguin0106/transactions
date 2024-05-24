package services

import (
	"context"
	"wallet/models"
	"wallet/repositories"
)

type WalletService struct {
	walletRepo  *repositories.WalletRepository
	accountRepo *repositories.AccountRepository
}

func NewWalletService(walletRepo *repositories.WalletRepository, accountRepo *repositories.AccountRepository) *WalletService {
	return &WalletService{walletRepo: walletRepo, accountRepo: accountRepo}
}

func (service *WalletService) GetWalletByUserId(ctx context.Context, userID int) (*models.Wallet, error) {
	return service.walletRepo.GetWalletByUserID(ctx, userID)
}

func (service *WalletService) CreateWallet(ctx context.Context, userID int) error {
	account, err := service.accountRepo.CreateAccount(ctx, "USD")
	if err != nil {
		return err
	}
	return service.walletRepo.CreateWallet(ctx, userID, account.AccountNumber)
}

func (service *WalletService) CreateAccount(ctx context.Context, userID int, currency string) (*models.Account, error) {
	account, err := service.accountRepo.CreateAccount(ctx, currency)
	if err != nil {
		return nil, err
	}

	err = service.walletRepo.AddAccountToWallet(ctx, userID, account.AccountNumber)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (service *WalletService) UpdateWallet(ctx context.Context, wallet *models.Wallet) error {
	return service.walletRepo.UpdateWallet(ctx, wallet)
}

func (service *WalletService) Deposit(ctx context.Context, accountNumber string, amount float64) error {
	return service.walletRepo.Deposit(ctx, accountNumber, amount)
}
