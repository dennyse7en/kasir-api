package model

import "time"

type Transaction struct {
	ID          int                 `json:"id"`
	TotalAmount float64             `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	Details     []TransactionDetail `json:"details,omitempty"`
}

type TransactionDetail struct {
	ID            int     `json:"id"`
	TransactionID int     `json:"transaction_id"`
	ProductID     int     `json:"product_id"`
	Quantity      int     `json:"quantity"`
	Subtotal      float64 `json:"subtotal"`
}

type TransactionRequestItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type TransactionRequest struct {
	Items []TransactionRequestItem `json:"items"`
}

type DailyReport struct {
	Date             string  `json:"date"`
	TotalSales       float64 `json:"total_sales"`
	TransactionCount int     `json:"transaction_count"`
}
