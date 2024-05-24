package repositories

import (
	"context"
	"database/sql"
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

func (repo *OrderRepository) GetOrders(ctx context.Context) ([]*models.Order, error) {
	query := "SELECT id, seller_id, cryptocurrency, amount, price, status, exchange_to FROM orders"
	rows, err := repo.DB.QueryContext(ctx, query)
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

func (repo *OrderRepository) GetOrderByID(ctx context.Context, orderID int) (*models.Order, error) {
	query := "SELECT id, seller_id, cryptocurrency, amount, price, status, exchange_to FROM orders WHERE id = $1"
	row := repo.DB.QueryRowContext(ctx, query, orderID)

	var order models.Order
	err := row.Scan(&order.ID, &order.SellerID, &order.Cryptocurrency, &order.Amount, &order.Price, &order.Status, &order.ExchangeTo)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &order, nil
}

func (repo *OrderRepository) UpdateOrder(ctx context.Context, order *models.Order) error {
	query := "UPDATE orders SET status = $1 WHERE id = $2"
	_, err := repo.DB.ExecContext(ctx, query, order.Status, order.ID)
	return err
}
