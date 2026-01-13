package repository

import (
	"OrderManagementSystem/internal/contracts"
	"OrderManagementSystem/internal/models"
	"database/sql"
)

var _ contracts.OrderRepository = (*dbQuery)(nil)

type dbQuery struct {
	db *sql.DB
}

func NewDBOrderRepository(db *sql.DB) contracts.OrderRepository {
	return &dbQuery{db: db}
}
func (s *dbQuery) CreateOrder(order models.Orders) (models.CommonResponse, error) {
	res, err := s.InserIntoTable(order)
	if err != nil {
		return models.CommonResponse{}, err
	}
	if res {
		return models.CommonResponse{
			Message: "Success",
			Status:  1,
		}, nil
	}
	return models.CommonResponse{
		Message: "failed",
		Status:  0,
	}, nil
}
func (s *dbQuery) InserIntoTable(order models.Orders) (bool, error) {
	_, err := s.db.Exec("Inser into tble.name(id,amount,status)value($1,$2,$3)", order.ID, order.Amount, order.Status)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (s *dbQuery) GetAllOrder() ([]models.Orders, error) {
	row, err := s.db.Query("Select * from tbl")
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var orders []models.Orders
	for row.Next() {
		var o models.Orders
		if err := row.Scan(&o.ID, &o.Amount, &o.Status); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
