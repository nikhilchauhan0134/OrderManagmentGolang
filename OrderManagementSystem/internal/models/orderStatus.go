package models

type OrderStatus string

const (
	OrderCreated OrderStatus = "CREATE"
	OrderPaid    OrderStatus = "PAID"
	OrderFailed  OrderStatus = "FAILED"
)
