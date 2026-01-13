package service

import (
	"OrderManagementSystem/internal/contracts"
	"OrderManagementSystem/internal/models"
)

type OrderService struct {
	repo contracts.OrderRepository
}

func NewOrderService(repo contracts.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}
func (s *OrderService) CreateOrder(order models.Orders) (models.CommonResponse, error) {
	if order.Amount == 0 {
		return models.CommonResponse{
			Message: "Amount must not 0",
			Status:  0,
		}, nil
	}

	if order.Status == "" {
		order.Status = "CREATE"
	}
	s.repo.CreateOrder(order)
	return models.CommonResponse{
		Message: "Success",
		Status:  1,
	}, nil

}
func (s *OrderService) GetAllOrder() ([]models.Orders, error) {
	return s.repo.GetAllOrder()
}
