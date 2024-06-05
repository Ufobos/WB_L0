package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"orders_service/database"
	"orders_service/models"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
)

// Тесты для HTTP-обработчиков.
func TestGetOrderHandler(t *testing.T) {
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
	database.CacheOrder(order)

	req, err := http.NewRequest("GET", "/order/test123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/order/{id}", GetOrderHandler).Methods("GET")
	router.ServeHTTP(rr, req)

	// Проверка респонда
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var got, want models.Order
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	want = order
	want.DateCreated = got.DateCreated
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("handler returned unexpected body (-got +want):\n%s", diff)
	}
}
