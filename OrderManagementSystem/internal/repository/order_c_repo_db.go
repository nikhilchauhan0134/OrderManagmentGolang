package repository

import (
	"OrderManagementSystem/internal/contracts"
	"OrderManagementSystem/internal/models"
	"context"
	"database/sql"
)

var _ contracts.OrderCurrencyRespository = (*dbQuery)(nil)

func NewDBOrderCurrencyRespository(db *sql.DB) contracts.OrderCurrencyRespository {
	return &dbQuery{db: db}
}
func (s *dbQuery) CreateOrderConcurrent(
	ctx context.Context,
	order models.Orders,
) (models.CommonResponse, error) {

	query := `
		INSERT INTO orders (id, amount, status)
		VALUES ($1, $2, $3)
	`

	_, err := s.db.ExecContext(
		ctx,
		query,
		order.ID,
		order.Amount,
		order.Status,
	)
	if err != nil {
		return models.CommonResponse{
			Status:  0,
			Message: "failed to create order",
		}, err
	}

	return models.CommonResponse{
		Status:  1,
		Message: "order created successfully",
	}, nil
}

func (s *dbQuery) GetOrderSummary(
	ctx context.Context,
) (map[string]interface{}, error) {

	query := `
		SELECT 
			COUNT(*) AS total_orders,
			COALESCE(SUM(amount), 0) AS total_amount
		FROM orders
	`

	row := s.db.QueryRowContext(ctx, query)

	var totalOrders int64
	var totalAmount float64

	if err := row.Scan(&totalOrders, &totalAmount); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_orders": totalOrders,
		"total_amount": totalAmount,
	}, nil
}

func (s *dbQuery) BulkOrderCreation(
	ctx context.Context,
	orders []models.Orders,
) (models.CommonResponse, error) {

	if len(orders) == 0 {
		return models.CommonResponse{}, nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.CommonResponse{}, err
	}

	query := `
		INSERT INTO orders (id, amount, status)
		VALUES ($1, $2, $3)
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		tx.Rollback()
		return models.CommonResponse{}, err
	}
	defer stmt.Close()

	for _, order := range orders {
		_, err := stmt.ExecContext(
			ctx,
			order.ID,
			order.Amount,
			order.Status,
		)
		if err != nil {
			tx.Rollback()
			return models.CommonResponse{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return models.CommonResponse{}, err
	}

	return models.CommonResponse{
		Status:  1,
		Message: "bulk order creation successful",
	}, nil
}
