package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"transaction/models"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (repo *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	query := `
        INSERT INTO orders (seller_id, cryptocurrency, amount, desired_currency, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `

	err := repo.DB.QueryRowContext(ctx, query, order.SellerID, order.Cryptocurrency, order.Amount, order.DesiredCurrency, order.Status, time.Now()).Scan(&order.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *OrderRepository) GetOrderByID(ctx context.Context, orderID int) (*models.Order, error) {
	query := "SELECT id, seller_id, buyer_id, cryptocurrency, amount, desired_currency, status, created_at FROM orders WHERE id = $1"

	row := repo.DB.QueryRowContext(ctx, query, orderID)

	var order models.Order
	err := row.Scan(&order.ID, &order.SellerID, &order.BuyerID, &order.Cryptocurrency, &order.Amount, &order.DesiredCurrency, &order.Status, &order.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	return &order, nil
}

func (repo *OrderRepository) UpdateOrderStatus(ctx context.Context, orderID int, status string) error {
	query := "UPDATE orders SET status = $1 WHERE id = $2"

	_, err := repo.DB.ExecContext(ctx, query, status, orderID)
	return err
}

func (repo *OrderRepository) UpdateOrder(ctx context.Context, order *models.Order) error {
	query := `
        UPDATE orders
        SET buyer_id = $1, status = $2
        WHERE id = $3
    `
	_, err := repo.DB.ExecContext(ctx, query, order.BuyerID, order.Status, order.ID)
	return err
}

func (repo *OrderRepository) DeleteOrder(ctx context.Context, orderID int) error {
	query := "DELETE FROM orders WHERE id = $1"
	_, err := repo.DB.ExecContext(ctx, query, orderID)
	return err
}

func (repo *OrderRepository) FindOrdersBySellerUsername(ctx context.Context, username string) ([]*models.Order, error) {
	query := `
        SELECT o.id, o.seller_id, o.buyer_id, o.cryptocurrency, o.amount, o.desired_currency, o.status, o.created_at
        FROM orders o
        JOIN users u ON o.seller_id = u.id
        WHERE u.username = $1
    `
	rows, err := repo.DB.QueryContext(ctx, query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.SellerID, &order.BuyerID, &order.Cryptocurrency, &order.Amount, &order.DesiredCurrency, &order.Status, &order.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	return orders, nil
}

func (repo *OrderRepository) FindOrdersByCryptocurrency(ctx context.Context, cryptocurrency string) ([]*models.Order, error) {
	query := `
        SELECT id, seller_id, buyer_id, cryptocurrency, amount, desired_currency, status, created_at
        FROM orders
        WHERE cryptocurrency = $1
    `
	rows, err := repo.DB.QueryContext(ctx, query, cryptocurrency)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.SellerID, &order.BuyerID, &order.Cryptocurrency, &order.Amount, &order.DesiredCurrency, &order.Status, &order.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	return orders, nil
}
