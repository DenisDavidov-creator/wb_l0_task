package cache

import (
	"testing"
	"wb_l0/internal/order"
)

func TestCache(t *testing.T) {

	cache := NewCache()

	testKey := "Test_one"

	testOrder := order.Order{
		OrderUID:    testKey,
		Entry:       "sdasd",
		TrackNumber: "sdasdsadasd",
	}

	cache.Set(testKey, testOrder)

	getOrder, found := cache.Get(testKey)

	if !found {
		t.Fatalf("set/get faild. Expected to find order with this uid %s", testKey)
	}

	if getOrder.OrderUID != testOrder.OrderUID {
		t.Errorf("Excpected orderUID: %s, we get %s", testOrder.OrderUID, getOrder.OrderUID)
	}
}
