package main

import (
	"encoding/json"
	"log"
	"time"

	stan "github.com/nats-io/stan.go"
)

type Delivery struct {
	ID       int    `json:"id"`
	OrderUID string `json:"order_uid"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Zip      string `json:"zip"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Region   string `json:"region"`
	Email    string `json:"email"`
}

type Payment struct {
	ID           int    `json:"id"`
	OrderUID     string `json:"order_uid"`
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ID          int    `json:"id"`
	OrderUID    string `json:"order_uid"`
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

type Order struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

func main() {
	// Подключение к NATS-Streaming серверу
	sc, err := stan.Connect("test-cluster", "publisher-123")
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	// Создание тестового заказа
	order := Order{
		OrderUID:    "b563feb7b2b84b6test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: Payment{
			Transaction:  "b563feb7b2b84b6test",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}

	// Преобразование заказа в JSON
	data, err := json.Marshal(order)
	if err != nil {
		log.Fatal(err)
	}

	// Публикация сообщения в канал order_channel
	err = sc.Publish("order_channel", data)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Сообщение успешно отправлено:", string(data))
}
