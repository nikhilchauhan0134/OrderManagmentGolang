package handler

import (
	"OrderManagementSystem/internal/models"
	"OrderManagementSystem/internal/service"
	"encoding/json"
	"net/http"
)

type OrderHandler struct {
	ser service.OrderService
}

func NewHandler(ser *service.OrderService) *OrderHandler {
	return &OrderHandler{ser: *ser}
}
func (ser *OrderHandler) SaveOrder(dtoResponse http.ResponseWriter, dtoRequest *http.Request) {
	var order models.Orders
	err := json.NewDecoder(dtoRequest.Body).Decode(&order)
	if err != nil {
		http.Error(dtoResponse, "Invalid json", http.StatusBadRequest)
	}
	res, err := ser.ser.CreateOrder(order)
	if err != nil {
		http.Error(dtoResponse, err.Error(), http.StatusInternalServerError)
	}
	dtoResponse.Header().Set("Content-Type", "application/json")
	//dtoResponse.Header("content-type","application/json")
	dtoResponse.WriteHeader(http.StatusCreated)
	json.NewEncoder(dtoResponse).Encode(res)
	//return res
}
func (ser *OrderHandler) GetAllOrderDetails(w http.ResponseWriter, r *http.Request) {
	orderlist, err := ser.ser.GetAllOrder()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderlist)
}
