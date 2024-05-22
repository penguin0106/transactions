package services

import (
	"context"
	"wallet/models"
	"wallet/repositories"
)

type WalletService struct {
	repo *repositories.WalletRepository
}

func NewWalletService(repo *repositories.WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

func (service *WalletService) GetWalletByUserID(ctx context.Context, userID int) (*models.Wallet, error) {
	return service.repo.GetWalletByUserID(ctx, userID)
}

func (service *WalletService) Deposit(ctx context.Context, userID int, amount float64) error {
	return service.repo.Deposit(ctx, userID, amount)
}
