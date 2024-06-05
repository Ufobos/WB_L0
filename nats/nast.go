package nats

import (
	"encoding/json"
	"log"
	"orders_service/database"
	"orders_service/models"

	stan "github.com/nats-io/stan.go"
)

func SubscribeToNATS() {
	log.Println("Подключение к NATS-Streaming...")
	sc, err := stan.Connect("test-cluster", "client-123")
	if err != nil {
		log.Fatalf("Ошибка при подключении к NATS-Streaming: %v", err)
	}
	defer sc.Close()
	log.Println("Успешное подключение к NATS-Streaming")

	_, err = sc.Subscribe("order_channel", func(m *stan.Msg) {
		handleMessage(m.Data)
	}, stan.DurableName("durable"))
	if err != nil {
		log.Fatalf("Ошибка при подписке на канал: %v", err)
	}
	log.Println("Успешная подписка на канал order_channel")
}

func handleMessage(data []byte) {
	var order models.Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		log.Println("Ошибка при парсинге JSON:", err)
		return
	}
	log.Println("Получено сообщение из NATS:", order.OrderUID)
	database.SaveOrderToDB(order)
}
