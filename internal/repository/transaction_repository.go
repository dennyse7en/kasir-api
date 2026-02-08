package repository

import (
	"context"
	"database/sql"
	"fmt"
	"kasir-api/internal/model"
	"time"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction model.Transaction, details []model.TransactionDetail) (model.Transaction, error)
	GetDailyReport(ctx context.Context, date time.Time) (model.DailyReport, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(ctx context.Context, transaction model.Transaction, details []model.TransactionDetail) (model.Transaction, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Transaction{}, err
	}
	defer tx.Rollback()

	// Insert Transaction
	query := `INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at`
	if err := tx.QueryRowContext(ctx, query, transaction.TotalAmount).Scan(&transaction.ID, &transaction.CreatedAt); err != nil {
		return model.Transaction{}, fmt.Errorf("failed to insert transaction: %w", err)
	}

	// Insert Details and Update Stock
	detailsQuery := `INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)`
	updateStockQuery := `UPDATE products SET stock = stock - $1 WHERE id = $2`

	for _, detail := range details {
		_, err := tx.ExecContext(ctx, detailsQuery, transaction.ID, detail.ProductID, detail.Quantity, detail.Subtotal)
		if err != nil {
			return model.Transaction{}, fmt.Errorf("failed to insert detail: %w", err)
		}

		_, err = tx.ExecContext(ctx, updateStockQuery, detail.Quantity, detail.ProductID)
		if err != nil {
			return model.Transaction{}, fmt.Errorf("failed to update stock: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return model.Transaction{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	transaction.Details = details
	return transaction, nil
}

func (r *transactionRepository) GetDailyReport(ctx context.Context, date time.Time) (model.DailyReport, error) {
	query := `
		SELECT 
			COALESCE(SUM(total_amount), 0) as total_sales,
			COUNT(id) as transaction_count
		FROM transactions 
		WHERE DATE(created_at) = $1
	`

	var report model.DailyReport
	report.Date = date.Format("2006-01-02")

	if err := r.db.QueryRowContext(ctx, query, report.Date).Scan(&report.TotalSales, &report.TransactionCount); err != nil {
		return model.DailyReport{}, err
	}

	return report, nil
}
