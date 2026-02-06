package models

import (
	"errors"
	"time"
)

type Transaction struct {
	ID          int                 `json:"id"`
	TotalAmount int                 `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	Details     []TransactionDetail `json:"details"`
}

type TransactionDetail struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name,omitempty"`
	Quantity      int    `json:"quantity"`
	Subtotal      int    `json:"subtotal"`
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

var (
	ErrProductNotFound       = errors.New("Product not found")
	ErrProductStockNotEnough = errors.New("Product stock not enough")
	ErrQuantityNotValid      = errors.New("Quantity not valid")
)

type SalesSummary struct {
	TotalRevenue        int              `json:"total_revenue"`
	TotalTransactions   int              `json:"total_transactions"`
	BestSellingProducts []SellingProduct `json:"best_selling_products"`
}

type SellingProduct struct {
	ProductName string `json:"product_name"`
	TotalSales  int    `json:"quantity_sales"`
}
