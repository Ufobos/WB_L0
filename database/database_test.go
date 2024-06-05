package database

import (
	"orders_service/models"
	"testing"
	"time"
)

func TestCacheOrder(t *testing.T) {
	order := models.Order{
		OrderUID: "test123",
		Delivery: models.Delivery{
			Name: "Test Name",
		},
		Payment: models.Payment{
			Transaction: "testtransaction",
		},
		Items: []models.Item{
			{ChrtID: 123, Name: "Test Item"},
		},
		Locale:      "en",
		CustomerID:  "testcustomer",
		DateCreated: time.Now(),
	}

	CacheOrder(order)
	cachedOrder, exists := GetOrderFromCache("test123")
	if !exists {
		t.Fatal("order not found in cache")
	}
	if cachedOrder.OrderUID != "test123" {
		t.Errorf("expected order UID %v, got %v", "test123", cachedOrder.OrderUID)
	}
}
