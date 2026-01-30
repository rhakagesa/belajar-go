package services

import (
	"belajar-go/models"
	"belajar-go/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (service *CategoryService) GetAll(page int, limit int) ([]models.Category, int, error) {
	return service.repo.GetAll(page, limit)
}

func (service *CategoryService) GetByID(id int) (*models.Category, error) {
	return service.repo.GetByID(id)
}

func (service *CategoryService) Create(category *models.Category) error {
	return service.repo.Create(category)
}

func (service *CategoryService) Update(category *models.Category) error {
	return service.repo.Update(category)
}

func (service *CategoryService) Delete(id int) error {
	return service.repo.Delete(id)
}
