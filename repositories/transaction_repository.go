package repositories

import (
	"belajar-go/models"
	"database/sql"
	"fmt"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		query := "SELECT name, price, stock FROM products WHERE id = $1"
		if useLock {
			query += " FOR UPDATE"
		}
		err := tx.QueryRow(query, item.ProductID).Scan(&productName, &productPrice, &stock)

		if err != nil {
			return nil, err
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)

		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	query := "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES "
	values := make([]any, 0)
	for i, detail := range details {
		n := i + 1
		query += fmt.Sprintf("($%d, $%d, $%d, $%d)", n*4-3, n*4-2, n*4-1, n*4)
		if i < len(details)-1 {
			query += ","
		}
		values = append(values, transactionID, detail.ProductID, detail.Quantity, detail.Subtotal)
	}
	query += " RETURNING id"

	_, err = tx.Exec(query, values...)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

func (repo *TransactionRepository) GetReportSales(startDate, endDate string) (*models.SalesSummary, error) {
	var salesSummary models.SalesSummary

	querySummary := `
        SELECT 
            COALESCE(SUM(total_amount), 0) AS total_revenue, 
            COUNT(*) AS total_transactions 
        FROM transactions 
        WHERE created_at::date BETWEEN $1 AND $2`

	err := repo.db.QueryRow(querySummary, startDate, endDate).Scan(
		&salesSummary.TotalRevenue,
		&salesSummary.TotalTransactions,
	)

	if err != nil {
		return nil, err
	}

	queryBestSelling := `
        SELECT 
            p.name as product_name, 
            SUM(td.quantity) as total_sales
        FROM transaction_details td
        JOIN products p ON td.product_id = p.id
        JOIN transactions t ON td.transaction_id = t.id
        WHERE t.created_at::date BETWEEN $1 AND $2
        GROUP BY p.name
        ORDER BY total_sales DESC`

	rows, err := repo.db.Query(queryBestSelling, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	salesSummary.BestSellingProducts = []models.SellingProduct{}

	for rows.Next() {
		var sellingProduct models.SellingProduct
		if err := rows.Scan(&sellingProduct.ProductName, &sellingProduct.TotalSales); err != nil {
			return nil, err
		}
		salesSummary.BestSellingProducts = append(salesSummary.BestSellingProducts, sellingProduct)
	}

	return &salesSummary, nil
}
