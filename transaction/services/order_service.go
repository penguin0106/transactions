package services

import (
	"context"
	"errors"
	"transaction/clients"
	"transaction/models"
	"transaction/repositories"
)

type OrderService struct {
	repo         *repositories.OrderRepository
	walletClient *clients.WalletClient
}

func NewOrderService(repo *repositories.OrderRepository, walletClient *clients.WalletClient) *OrderService {
	return &OrderService{
		repo:         repo,
		walletClient: walletClient,
	}
}

func (service *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	order.Status = "PENDING"
	return service.repo.CreateOrder(ctx, order)
}

func (service *OrderService) GetOrderByID(ctx context.Context, orderID int) (*models.Order, error) {
	return service.repo.GetOrderByID(ctx, orderID)
}

func (service *OrderService) UpdateOrderStatus(ctx context.Context, orderID int, status string) error {
	return service.repo.UpdateOrderStatus(ctx, orderID, status)
}

func (service *OrderService) PurchaseOrder(ctx context.Context, buyerID, orderID int) error {
	order, err := service.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != "PENDING" {
		return errors.New("order is not available for purchase")
	}

	// Perform transaction: deduct from buyer and credit to seller
	if order.DesiredCurrency == "USD" {
		err = service.walletClient.UpdateBalance(ctx, buyerID, -order.Amount)
		if err != nil {
			return err
		}
		err = service.walletClient.UpdateBalance(ctx, order.SellerID, order.Amount)
		if err != nil {
			return err
		}
	} else {
		// Deduct desired currency from buyer
		err = service.walletClient.UpdateBalance(ctx, buyerID, -order.Amount)
		if err != nil {
			return err
		}
		// Credit cryptocurrency to buyer
		err = service.walletClient.UpdateBalance(ctx, order.SellerID, order.Amount)
		if err != nil {
			return err
		}
	}

	// Update order status and buyer ID
	order.BuyerID = buyerID
	order.Status = "COMPLETED"
	return service.repo.UpdateOrder(ctx, order)
}
