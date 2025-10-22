package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"wb_l0/internal/order"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type MockCache struct {
	store map[string]order.Order
}

func (m *MockCache) Get(key string) (order.Order, bool) {
	val, found := m.store[key]
	return val, found
}
func (m *MockCache) Set(key string, value order.Order) {
	if m.store == nil {
		m.store = make(map[string]order.Order)
	}
	m.store[key] = value
}

func TestOrderHandler_GetByUID(t *testing.T) {

	testOrder := order.Order{
		OrderUID:        "test",
		TrackNumber:     "WBILMTESTTRACK",
		Entry:           "WBIL",
		Delivery:        order.Delivery{Name: "Test", Phone: "+123", Zip: "12345", City: "Test City", Address: "Test Address", Email: "test@test.com"},
		Payment:         order.Payment{Transaction: "b563feb7b2b84b6test", Currency: "USD", Provider: "wbpay", Amount: 1, PaymentDt: 1, Bank: "alpha", DeliveryCost: 1, GoodsTotal: 1},
		Items:           []order.Item{{ChrtID: 1, TrackNumber: "WBILMTESTTRACK", Price: 1, Rid: "rid", Name: "name", Sale: 0, Size: "0", TotalPrice: 1, NmID: 1, Status: 200, Brand: "brand"}},
		Locale:          "en",
		CustomerID:      "test",
		DeliveryService: "meest",
		DateCreated:     time.Now(),
	}

	mockCache := &MockCache{}
	mockCache.Set(testOrder.OrderUID, testOrder)

	handler := NewOrderHandler(mockCache)

	router := chi.NewRouter()
	router.Get("/orders/{orderUID}", handler.GetByUID)

	t.Run("should return order when found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/orders/test", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "The code should be 200 OK")

		var returnOrder order.Order
		err := json.Unmarshal(rr.Body.Bytes(), &returnOrder)
		assert.NoError(t, err, "The body should be valid of JSON")
		assert.Equal(t, testOrder.OrderUID, returnOrder.OrderUID, "Order UID don't match with expected")
	})

	t.Run("Should return 404 when not found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/orders/te", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code, "The code should be 404 not found ")
	})
}
