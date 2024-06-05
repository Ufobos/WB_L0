package database

import (
	"database/sql"
	"log"
	"orders_service/models"
	"sync"

	_ "github.com/lib/pq"
)

var db *sql.DB

var cache = struct {
	sync.RWMutex
	orders map[string]models.Order
}{orders: make(map[string]models.Order)}

// Подключение к базе данных PostgreSQL.
func InitDB() {
	var err error
	connStr := "user=postgres password=176178 dbname=orders_db sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Ошибка при проверке соединения к базе данных: %v", err)
	}
	log.Println("Подключение к базе данных PostgreSQL установлено успешно!")
}

// Добавлние заказа в кэш.
func CacheOrder(order models.Order) {
	cache.Lock()
	cache.orders[order.OrderUID] = order
	cache.Unlock()
	log.Println("Заказ добавлен в кэш:", order.OrderUID)
}

// Получение заказа из кэша.
func GetOrderFromCache(orderUID string) (models.Order, bool) {
	cache.RLock()
	order, exists := cache.orders[orderUID]
	cache.RUnlock()
	return order, exists
}

// Восстановление кэша из базы данных.
func RestoreCache() {
	rows, err := db.Query("SELECT * FROM orders")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard)
		if err != nil {
			log.Fatal(err)
		}
		CacheOrder(order)
	}
	log.Println("Кэш восстановлен из базы данных")
}

// Сохраниение заказа в БД, обновляя его если он уже существует, или вставляя новый.
func SaveOrderToDB(order models.Order) {
	var existingOrderUID string
	err := db.QueryRow("SELECT order_uid FROM orders WHERE order_uid = $1", order.OrderUID).Scan(&existingOrderUID)

	if err == nil {
		log.Println("Заказ уже существует, обновляем данные:", order.OrderUID)
		updateOrderInDB(order)
	} else if err == sql.ErrNoRows {
		insertOrderToDB(order)
	} else {
		log.Fatal(err)
	}
}

// Заказ не существует, вставляем новый
func insertOrderToDB(order models.Order) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec("INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	_, err = tx.Exec("INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	_, err = tx.Exec("INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	for _, item := range order.Items {
		_, err = tx.Exec("INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Заказ успешно сохранен в базе данных:", order.OrderUID)
	CacheOrder(order)
}

// Заказ уже существует, обновляем его
func updateOrderInDB(order models.Order) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec("UPDATE orders SET track_number = $2, entry = $3, locale = $4, internal_signature = $5, customer_id = $6, delivery_service = $7, shardkey = $8, sm_id = $9, date_created = $10, oof_shard = $11 WHERE order_uid = $1",
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	_, err = tx.Exec("UPDATE delivery SET name = $2, phone = $3, zip = $4, city = $5, address = $6, region = $7, email = $8 WHERE order_uid = $1",
		order.Delivery.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	_, err = tx.Exec("UPDATE payment SET transaction = $2, request_id = $3, currency = $4, provider = $5, amount = $6, payment_dt = $7, bank = $8, delivery_cost = $9, goods_total = $10, custom_fee = $11 WHERE order_uid = $1",
		order.Payment.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	_, err = tx.Exec("DELETE FROM items WHERE order_uid = $1", order.OrderUID)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	for _, item := range order.Items {
		_, err = tx.Exec("INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Заказ успешно обновлен в базе данных:", order.OrderUID)
	CacheOrder(order)
}
