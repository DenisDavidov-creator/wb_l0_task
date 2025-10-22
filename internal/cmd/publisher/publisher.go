package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nats-io/stan.go"
)

func main() {
	file, err := os.ReadFile("./orders.json")
	if err != nil {
		log.Fatalf("Не удалось прочитать файл с данными: %v", err)
	}

	var orders []map[string]interface{}
	if err := json.Unmarshal(file, &orders); err != nil {
		log.Fatalf("Не удалось декодировать JSON: %v", err)
	}

	channel, sc := connectToStan()
	defer sc.Close()

	order := orders[0]

	message, err := json.Marshal(order)
	if err != nil {
		log.Printf("Ошибка кодирования заказа в JSON: %v", err)
	}

	err = sc.Publish(channel, message)
	if err != nil {
		log.Printf("Ошибка публикации: %v", err)
	} else {
		log.Printf("Сообщение для UID %s успешно отправлено.", order["order_uid"])
	}
}

func connectToStan() (string, stan.Conn) {
	clusterID := "test-cluster"
	clientID := "ImPublisher"
	natsURL := "nats://localhost:4222"
	channel := "orders-channel"

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Не удалось подключиться к Nats: %v", err)
	}

	log.Printf("Подключились к '%s' c clientID '%s'", clusterID, clientID)

	return channel, sc
}
