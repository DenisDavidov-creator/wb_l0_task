package handler

import (
	"encoding/json"
	"log"
	"time"
	"wb_l0/internal/cache"
	"wb_l0/internal/order"

	"github.com/nats-io/stan.go"
)

type Subscriber struct {
	repo  *order.Repository
	cache *cache.Cache
}

func NewSubscriber(repo *order.Repository, cache *cache.Cache) *Subscriber {
	return &Subscriber{repo: repo, cache: cache}
}

func (s *Subscriber) Subscribe(sc stan.Conn) (stan.Subscription, error) {

	messageHandler := func(msg *stan.Msg) {
		recieveMessage := string(msg.Data)

		var orderData order.Order
		if err := json.Unmarshal([]byte(recieveMessage), &orderData); err != nil {
			log.Printf("Ошибка декодирования")
			msg.Ack()
			return
		}

		if err := orderData.Validate(); err != nil {
			log.Printf("Данные не прошли валидацию, %v", err)
			msg.Ack()
			return
		}

		isCreated, err := s.repo.Create(orderData)
		if err != nil {
			log.Printf("Ошибка сохранения в БД %v", err)
			return
		}
		if isCreated {
			log.Printf("Заказ %s создан и сохранен", orderData.OrderUID)
			s.cache.Set(orderData.OrderUID, orderData)
		} else {
			log.Printf("Заказ %s уже существует", orderData.OrderUID)
		}

		msg.Ack()

	}
	channel := "orders-channel"
	durableName := "my-durable-name"

	return sc.Subscribe(
		channel,
		messageHandler,
		stan.SetManualAckMode(),
		stan.DurableName(durableName),
		stan.AckWait(30*time.Second),
		stan.DeliverAllAvailable(),
	)
}
