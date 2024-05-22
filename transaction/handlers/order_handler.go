package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"transaction/models"
	"transaction/services"
)

func CreateOrderHandler(service *services.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var order models.Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := service.CreateOrder(r.Context(), &order); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(order)
	}
}

func GetOrderHandler(service *services.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderIDStr := r.URL.Query().Get("order_id")
		orderID, err := strconv.Atoi(orderIDStr)
		if err != nil {
			http.Error(w, "Invalid order ID", http.StatusBadRequest)
			return
		}

		order, err := service.GetOrderByID(r.Context(), orderID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(order)
	}
}

// PurchaseRequest represents the request body for purchasing an order.
type PurchaseRequest struct {
	BuyerID int `json:"buyer_id"`
}

// PurchaseOrderHandler handles the request for purchasing an order.
func PurchaseOrderHandler(service *services.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderIDStr := r.URL.Query().Get("order_id")
		orderID, err := strconv.Atoi(orderIDStr)
		if err != nil {
			http.Error(w, "Invalid order ID", http.StatusBadRequest)
			return
		}

		var req PurchaseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.BuyerID <= 0 {
			http.Error(w, "Invalid buyer ID", http.StatusBadRequest)
			return
		}

		err = service.PurchaseOrder(r.Context(), req.BuyerID, orderID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

// DeleteOrderHandler handles the request for deleting an order.
func DeleteOrderHandler(service *services.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderIDStr := r.URL.Query().Get("order_id")
		orderID, err := strconv.Atoi(orderIDStr)
		if err != nil {
			http.Error(w, "Invalid order ID", http.StatusBadRequest)
			return
		}

		err = service.DeleteOrder(r.Context(), orderID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

// FindOrdersBySellerUsernameHandler handles the request for finding orders by seller's username.
func FindOrdersBySellerUsernameHandler(service *services.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		if username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}

		orders, err := service.FindOrdersBySellerUsername(r.Context(), username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(orders)
	}
}

// FindOrdersByCryptocurrencyHandler handles the request for finding orders by cryptocurrency.
func FindOrdersByCryptocurrencyHandler(service *services.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cryptocurrency := r.URL.Query().Get("cryptocurrency")
		if cryptocurrency == "" {
			http.Error(w, "Cryptocurrency is required", http.StatusBadRequest)
			return
		}

		orders, err := service.FindOrdersByCryptocurrency(r.Context(), cryptocurrency)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(orders)
	}
}
