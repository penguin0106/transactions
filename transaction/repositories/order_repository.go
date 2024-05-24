package repositories

import (
	"context"
	"database/sql"
	"errors"
	"transaction/models"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (repo *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) (int, error) {
	var orderID int
	query := "INSERT INTO orders (seller_id, cryptocurrency, amount, price, status, exchange_to) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	err := repo.DB.QueryRowContext(ctx, query, order.SellerID, order.Cryptocurrency, order.Amount, order.Price, order.Status, order.ExchangeTo).Scan(&orderID)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

func (repo *OrderRepository) GetOrdersBySellerUsername(ctx context.Context, username string) ([]*models.Order, error) {
	query := "SELECT id, seller_id, cryptocurrency, amount, price, status, exchange_to FROM orders WHERE seller_id IN (SELECT id FROM users WHERE username = $1)"
	rows, err := repo.DB.QueryContext(ctx, query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.SellerID, &order.Cryptocurrency, &order.Amount, &order.Price, &order.Status, &order.ExchangeTo)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

func (repo *OrderRepository) GetOrdersByCurrency(ctx context.Context, currency string) ([]*models.Order, error) {
	query := "SELECT id, seller_id, cryptocurrency, amount, price, status, exchange_to FROM orders WHERE cryptocurrency = $1"
	rows, err := repo.DB.QueryContext(ctx, query, currency)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.SellerID, &order.Cryptocurrency, &order.Amount, &order.Price, &order.Status, &order.ExchangeTo)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

func (repo *OrderRepository) PurchaseOrder(ctx context.Context, buyerID, orderID int) error {
	order, err := repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != "PENDING" {
		return errors.New("order is not available for purchase")
	}

	order.BuyerID = buyerID
	order.Status = "COMPLETED"

	return repo.UpdateOrder(ctx, order)
}

func (repo *OrderRepository) GetOrderByID(ctx context.Context, orderID int) (*models.Order, error) {
	query := "SELECT seller_id, buyer_id, cryptocurrency, amount, price, status, exchange_to FROM orders WHERE id = $1"
	row := repo.DB.QueryRowContext(ctx, query, orderID)

	var order models.Order
	err := row.Scan(&order.SellerID, &order.BuyerID, &order.Cryptocurrency, &order.Amount, &order.Price, &order.Status, &order.ExchangeTo)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	order.ID = orderID
	return &order, nil
}

func (repo *OrderRepository) UpdateOrder(ctx context.Context, order *models.Order) error {
	query := "UPDATE orders SET seller_id = $1, buyer_id = $2, cryptocurrency = $3, amount = $4, price = $5, status = $6, exchange_to = $7 WHERE id = $8"
	_, err := repo.DB.ExecContext(ctx, query, order.SellerID, order.BuyerID, order.Cryptocurrency, order.Amount, order.Price, order.Status, order.ExchangeTo, order.ID)
	return err
}
