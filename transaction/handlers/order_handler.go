package handlers

import (
	"encoding/json"
	"net/http"
	"transaction/services"
)

type OrderHandler struct {
	OrderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{OrderService: orderService}
}

func (handler *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderData struct {
		SellerID       int
		Cryptocurrency string
		Amount         float64
		Price          float64
		ExchangeTo     string
	}

	err := json.NewDecoder(r.Body).Decode(&orderData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = handler.OrderService.CreateOrder(r.Context(), orderData.SellerID, orderData.Cryptocurrency, orderData.Amount, orderData.Price, orderData.ExchangeTo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *OrderHandler) FindOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := handler.OrderService.FindOrders(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}

func (handler *OrderHandler) FindOrdersByCurrency(w http.ResponseWriter, r *http.Request) {
	currency := r.URL.Query().Get("currency")
	if currency == "" {
		http.Error(w, "Currency parameter is required", http.StatusBadRequest)
		return
	}

	orders, err := handler.OrderService.FindOrdersByCurrency(r.Context(), currency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}

func (handler *OrderHandler) FindOrdersBySellerUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username parameter is required", http.StatusBadRequest)
		return
	}

	orders, err := handler.OrderService.FindOrdersBySellerUsername(r.Context(), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}

func (handler *OrderHandler) PurchaseOrder(w http.ResponseWriter, r *http.Request) {
	var purchaseData struct {
		BuyerID int
		OrderID int
	}

	err := json.NewDecoder(r.Body).Decode(&purchaseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = handler.OrderService.PurchaseOrder(r.Context(), purchaseData.BuyerID, purchaseData.OrderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
