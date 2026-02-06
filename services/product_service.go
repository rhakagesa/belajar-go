package services

import (
	"belajar-go/models"
	"belajar-go/repositories"
)

type ProductService struct {
	repo         *repositories.ProductRepository
	categoryRepo *repositories.CategoryRepository
}

func NewProductService(repo *repositories.ProductRepository, categoryRepo *repositories.CategoryRepository) *ProductService {
	return &ProductService{repo: repo, categoryRepo: categoryRepo}
}

func (service *ProductService) GetAll(page int, limit int, name string) ([]models.Product, int, error) {
	return service.repo.GetAll(page, limit, name)
}

func (service *ProductService) GetByID(id int) (*models.Product, error) {
	return service.repo.GetByID(id)
}

func (service *ProductService) Create(product *models.Product) error {
	_, err := service.categoryRepo.GetByID(product.CategoryID)

	if err != nil {
		return err
	}

	return service.repo.Create(product)
}

func (service *ProductService) Update(product *models.Product) error {
	_, err := service.categoryRepo.GetByID(product.CategoryID)

	if err != nil {
		return err
	}

	return service.repo.Update(product)
}

func (service *ProductService) Delete(id int) error {
	return service.repo.Delete(id)
}
