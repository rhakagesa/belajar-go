package services

import (
	"belajar-go/models"
	"belajar-go/repositories"
)

type TransactionService struct {
	repo        *repositories.TransactionRepository
	productRepo *repositories.ProductRepository
}

func NewTransactionService(repo *repositories.TransactionRepository, productRepo *repositories.ProductRepository) *TransactionService {
	return &TransactionService{repo: repo, productRepo: productRepo}
}

func (service *TransactionService) Checkout(items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {

	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, models.ErrQuantityNotValid
		}

		product, err := service.productRepo.GetByID(item.ProductID)
		if err != nil {
			return nil, models.ErrProductNotFound
		}

		if product.Stock < item.Quantity {
			return nil, models.ErrProductStockNotEnough
		}
	}

	return service.repo.CreateTransaction(items, true)
}

func (service *TransactionService) GetReportSales(startDate, endDate string) (*models.SalesSummary, error) {
	return service.repo.GetReportSales(startDate, endDate)
}
