package models

import (
	"time"
)

type OrderStatus string

const (
	StatusNew       OrderStatus = "new"
	StatusConfirmed OrderStatus = "confirmed"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

type OrderSource string

const (
	SourceOne     OrderSource = "źródło_jeden"
	SourceTwo     OrderSource = "źródło_dwa"
	SourceWebsite OrderSource = "website"
	SourceManual  OrderSource = "manual"
)

type Order struct {
	ID            int         `json:"id" db:"id"`
	CustomerName  string      `json:"customer_name" db:"customer_name"`
	CustomerEmail string      `json:"customer_email" db:"customer_email"`
	Source        OrderSource `json:"source" db:"source"`
	Status        OrderStatus `json:"status" db:"status"`
	TotalAmount   float64     `json:"total_amount" db:"total_amount"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
	Items         []OrderItem `json:"items,omitempty" db:"-"`
}

type OrderItem struct {
	ID          int     `json:"id" db:"id"`
	OrderID     int     `json:"order_id" db:"order_id"`
	ProductName string  `json:"product_name" db:"product_name"`
	Quantity    int     `json:"quantity" db:"quantity"`
	Price       float64 `json:"price" db:"price"`
}
