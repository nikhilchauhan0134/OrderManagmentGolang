package models

type Orders struct {
	ID     string      `json:"id"`
	Amount float64     `json:"amount"`
	Status OrderStatus `json:"status"`
}
