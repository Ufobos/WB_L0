package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"orders_service/database"
	"orders_service/models"

	"github.com/gorilla/mux"
)

// Обработчик HTTP-запросов для получения заказа по его идентфикатору.
func GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderUID := vars["id"]
	order, exists := database.GetOrderFromCache(orderUID)
	if !exists {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// Обработчик сообщений из NATS, который валидирует данные и сохраняет их в базу данных.
func HandleMessage(data []byte) {
	var order models.Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		log.Println("Ошибка при парсинге JSON:", err)
		return
	}

	err = models.ValidateOrder(order)
	if err != nil {
		log.Println("Ошибка валидации заказа:", err)
		return
	}

	log.Println("Получено сообщение из NATS:", order.OrderUID)
	database.SaveOrderToDB(order)
}
