package main

import (
	"log"
	"net/http"
	"orders_service/database"
	"orders_service/handlers"

	"github.com/gorilla/mux"
	stan "github.com/nats-io/stan.go"
)

func main() {
	database.InitDB()
	database.RestoreCache()
	go subscribeToNATS()

	r := mux.NewRouter()
	r.HandleFunc("/order/{id}", handlers.GetOrderHandler).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))
	http.Handle("/", r)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func subscribeToNATS() {
	log.Println("Подключение к NATS-Streaming...")
	sc, err := stan.Connect("test-cluster", "client-123")
	if err != nil {
		log.Fatalf("Ошибка при подключении к NATS-Streaming: %v", err)
	}
	defer sc.Close()
	log.Println("Успешное подключение к NATS-Streaming")

	_, err = sc.Subscribe("order_channel", func(m *stan.Msg) {
		handlers.HandleMessage(m.Data)
	}, stan.DurableName("durable"))
	if err != nil {
		log.Fatalf("Ошибка при подписке на канал: %v", err)
	}
	log.Println("Успешная подписка на канал order_channel")
}
