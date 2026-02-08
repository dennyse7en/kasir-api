package service

import (
	"context"
	"fmt"
	"kasir-api/internal/model"
	"kasir-api/internal/repository"
	"time"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, request model.TransactionRequest) (model.Transaction, error)
	GetDailyReport(ctx context.Context) (model.DailyReport, error)
}

type transactionService struct {
	repo        repository.TransactionRepository
	productRepo repository.ProductRepository
}

func NewTransactionService(repo repository.TransactionRepository, productRepo repository.ProductRepository) TransactionService {
	return &transactionService{repo: repo, productRepo: productRepo}
}

func (s *transactionService) CreateTransaction(ctx context.Context, request model.TransactionRequest) (model.Transaction, error) {
	var totalAmount float64
	var details []model.TransactionDetail

	for _, item := range request.Items {
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			return model.Transaction{}, fmt.Errorf("product not found: %d", item.ProductID)
		}

		if product.Stock < item.Quantity {
			return model.Transaction{}, fmt.Errorf("insufficient stock for product: %s", product.Name)
		}

		subtotal := product.Price * float64(item.Quantity)
		totalAmount += subtotal

		details = append(details, model.TransactionDetail{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Subtotal:  subtotal,
		})
	}

	transaction := model.Transaction{
		TotalAmount: totalAmount,
	}

	return s.repo.CreateTransaction(ctx, transaction, details)
}

func (s *transactionService) GetDailyReport(ctx context.Context) (model.DailyReport, error) {
	return s.repo.GetDailyReport(ctx, time.Now())
}
