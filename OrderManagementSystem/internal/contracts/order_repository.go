package contracts

import "OrderManagementSystem/internal/models"

type OrderRepository interface {
	CreateOrder(order models.Orders) (models.CommonResponse, error)
}
