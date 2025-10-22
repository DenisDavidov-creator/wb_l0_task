package main

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"

	"wb_l0/internal/cache"
	"wb_l0/internal/handler"
	"wb_l0/internal/order"
)

func main() {
	// ---------------------------Подключение к PostgreSQL---------------------------
	db := connectToDB()
	defer db.Close()
	log.Println("Успешно подключено к PostgreSQL")

	// -------------------------Подключение к NATS Streaming-------------------------
	sc := conntectToStan()
	defer sc.Close()

	// ------------------------------------------------------------------------------
	orderRepo := order.NewRepository(db)
	orderCache := cache.NewCache()

	// --------------------------------Обновление кэша-------------------------------
	err := restoreCache(orderCache, orderRepo)
	if err != nil {
		log.Fatalf("Не удалось востановить кэш: %v", err)
	}

	// ------------------------------------------------------------------------------
	subscriber := handler.NewSubscriber(orderRepo, orderCache)

	orderHTTPHandler := handler.NewOrderHandler(orderCache)

	sub, err := subscriber.Subscribe(sc)
	if err != nil {
		log.Fatalf("Не удалось запустить subscriber: %v", err)
	}
	defer sub.Unsubscribe()

	log.Println("Подписка на канал активна")

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// Укажи здесь адрес своего Live Server
		AllowedOrigins: []string{"http://localhost:5500", "http://127.0.0.1:5500"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		// Внутри группы путь уже относительный
		r.Get("/orders/{orderUID}", orderHTTPHandler.GetByUID)
	})

	log.Println("Запускаем сервер на порту 8080...")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Ошибка запуска сервера %v", err)
	}

	//---------------------------
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Получен сигнал завершения, отписка от канала")
}

func connectToDB() *sql.DB {

	connStr := "user=postgres password=123 dbname=wb_orders sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Не получилось подключиться к БД: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Не удалось проверить соединение с БД: %v", err)
	}

	return db
}

func conntectToStan() stan.Conn {
	clusterID := "test-cluster"
	clientID := "ImServer"
	natsURL := "nats://localhost:4222"

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Не удалось подключится к Nats: %v", err)
	}

	log.Printf("Подключились к %s c cliendID %s", clusterID, clientID)

	return sc
}

func restoreCache(c *cache.Cache, r *order.Repository) error {

	allOrders, err := r.GetAll()
	if err != nil {

		return err
	}

	for _, order := range allOrders {
		c.Set(order.OrderUID, order)
	}

	log.Println("Кэш востановлен")
	return nil
}
