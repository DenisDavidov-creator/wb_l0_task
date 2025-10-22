package cache

import (
	"fmt"
	"sync"
	"wb_l0/internal/order"
)

type OrderCache interface {
	Set(key string, value order.Order)
	Get(key string) (order.Order, bool)
}

type Cache struct {
	mu    sync.RWMutex
	items map[string]order.Order
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]order.Order),
	}
}

func (c *Cache) Set(key string, value order.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = value
	fmt.Printf("КЭШ [ЗАПИСЬ] ключ '%s', значение '%s'", key, value.OrderUID)
}

func (c *Cache) Get(key string) (order.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if found {
		fmt.Printf("КЭШ [ЧТЕНИЕ] Ключ '%s'найдено. Значение '%s'\n", key, item.OrderUID)
	} else {
		fmt.Printf("КЭШ [ЧТЕНИЕ] Ключ '%s' не найден\n", key)
	}
	return item, found
}
