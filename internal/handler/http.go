package handler

import (
	"encoding/json"
	"net/http"
	"wb_l0/internal/cache"

	"github.com/go-chi/chi/v5"
)

type OrderHandler struct {
	cache cache.OrderCache
}

func NewOrderHandler(c cache.OrderCache) *OrderHandler {
	return &OrderHandler{cache: c}
}

func (h *OrderHandler) GetByUID(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "orderUID")
	if uid == "" {
		http.Error(w, "Order is empty", http.StatusNotFound)
		return
	}
	order, found := h.cache.Get(uid)
	if !found {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}
