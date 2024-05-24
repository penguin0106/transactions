package services

import (
	"context"
	"errors"
	"transaction/models"
	"transaction/repositories"
)

type OrderService struct {
	orderRepo     *repositories.OrderRepository
	walletService *WalletService
}

func NewOrderService(repo *repositories.OrderRepository, walletService *WalletService) *OrderService {
	return &OrderService{
		orderRepo:     repo,
		walletService: walletService,
	}
}

func (service *OrderService) CreateOrder(ctx context.Context, sellerID int, cryptocurrency string, amount, price float64, exchangeTo string) error {
	// Создаем новый заказ
	order := &models.Order{
		SellerID:       sellerID,
		Cryptocurrency: cryptocurrency,
		Amount:         amount,
		Price:          price,
		ExchangeTo:     exchangeTo,
		Status:         "PENDING",
	}

	// Добавляем заказ в базу данных
	_, err := service.orderRepo.CreateOrder(ctx, order)
	return err
}

func (service *OrderService) FindOrders(ctx context.Context) ([]*models.Order, error) {
	// Получаем все заказы с торговой площадки
	orders, err := service.orderRepo.GetOrders(ctx)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (service *OrderService) FindOrdersByCurrency(ctx context.Context, currency string) ([]*models.Order, error) {
	// Получаем все заказы на торговой площадке по заданной валюте
	orders, err := service.orderRepo.GetOrdersByCurrency(ctx, currency)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (service *OrderService) FindOrdersBySellerUsername(ctx context.Context, username string) ([]*models.Order, error) {
	// Получаем все заказы на торговой площадке от продавца с заданным именем пользователя
	orders, err := service.orderRepo.GetOrdersBySellerUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (service *OrderService) PurchaseOrder(ctx context.Context, buyerID, orderID int) error {
	// Получаем информацию о заказе
	order, err := service.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Проверяем, что заказ существует и его статус "PENDING"
	if order == nil || order.Status != "PENDING" {
		return errors.New("order is not available for purchase")
	}

	// Получаем номера счетов продавца и покупателя
	sellerWallet, err := service.walletService.GetUserWallet(ctx, order.SellerID)
	if err != nil {
		return err
	}

	buyerWallet, err := service.walletService.GetUserWallet(ctx, buyerID)
	if err != nil {
		return err
	}

	sellerAccountNumber := ""
	buyerAccountNumber := ""

	// Ищем счета продавца и покупателя
	for _, account := range sellerWallet.Accounts {
		if account[:3] == order.Cryptocurrency {
			sellerAccountNumber = account
			break
		}
	}

	for _, account := range buyerWallet.Accounts {
		if account[:3] == order.ExchangeTo {
			buyerAccountNumber = account
			break
		}
	}

	if sellerAccountNumber == "" || buyerAccountNumber == "" {
		return errors.New("appropriate accounts not found for transaction")
	}

	// Определяем количество валюты для обмена
	exchangeAmount := order.Price * order.Amount

	// Выполняем перевод средств
	err = service.walletService.Transfer(ctx, exchangeAmount, buyerAccountNumber, sellerAccountNumber)
	if err != nil {
		return err
	}

	// Теперь переводим криптовалюту от продавца к покупателю
	err = service.walletService.Transfer(ctx, order.Amount, sellerAccountNumber, buyerAccountNumber)
	if err != nil {
		// В случае ошибки возвращаем первый перевод
		_ = service.walletService.Transfer(ctx, exchangeAmount, sellerAccountNumber, buyerAccountNumber)
		return err
	}

	// Обновляем статус заказа
	order.Status = "COMPLETED"
	err = service.orderRepo.UpdateOrder(ctx, order)
	if err != nil {
		return err
	}

	return nil
}
