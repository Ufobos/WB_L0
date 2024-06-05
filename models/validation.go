package models

import (
	"errors"
)

// Функция для валидации данных.
func ValidateOrder(order Order) error {
	if order.OrderUID == "" {
		return errors.New("order_uid is required")
	}
	if order.Delivery.Name == "" {
		return errors.New("delivery name is required")
	}
	if order.Payment.Transaction == "" {
		return errors.New("payment transaction is required")
	}
	if len(order.Items) == 0 {
		return errors.New("at least one item is required")
	}
	for _, item := range order.Items {
		if item.ChrtID == 0 {
			return errors.New("item chrt_id is required")
		}
		if item.Name == "" {
			return errors.New("item name is required")
		}
	}
	return nil
}
