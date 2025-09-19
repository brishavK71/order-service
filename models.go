package main

import "time"

type OrderRequest struct {
	CustomerID string            `json:"customer_id" validate:"required"`
	Items      []OrderItem       `json:"items" validate:"required,dive"`
	Meta       map[string]string `json:"meta,omitempty"`
}

type OrderItem struct {
	SKU      string  `json:"sku" validate:"required"`
	Quantity int     `json:"quantity" validate:"required,min=1"`
	Price    float64 `json:"price" validate:"required"`
}

type Order struct {
	ID         string      `json:"id"`
	CustomerID string      `json:"customer_id"`
	Items      []OrderItem `json:"items"`
	Status     string      `json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
}

type OrderEvent struct {
	EventType string    `json:"event_type"`
	OrderID   string    `json:"order_id"`
	Payload   Order     `json:"payload"`
	Time      time.Time `json:"time"`
}
