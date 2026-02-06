package handler

import (
	"belajar-go/helper"
	"belajar-go/models"
	"belajar-go/services"
	"encoding/json"
	"net/http"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (handler CategoryHandler) HandleCategory(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handler.GetAll(w, r)
	case http.MethodPost:
		handler.Create(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
func (handler CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handler.GetByID(w, r)
	case http.MethodPut:
		handler.Update(w, r)
	case http.MethodDelete:
		handler.Delete(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (handler CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, limit, err := helper.PaginationConverter(r)

	if err != nil {
		helper.ResponseJson(w, false, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	categories, totalData, err := handler.service.GetAll(page, limit)
	if err != nil {
		helper.ResponseJson(w, false, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	if totalData == 0 {
		helper.ResponseJson(w, true, http.StatusOK, "No Categories Found", nil, nil)
		return
	}

	pagination := helper.BuildPagination(totalData, limit, page)

	helper.ResponseJson(w, true, http.StatusOK, "Categories Found", categories, &pagination)
}

func (handler CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := helper.UriIdConverter(r.URL.Path, "/api/categories/")
	if err != nil {
		helper.ResponseJson(w, false, http.StatusBadRequest, "Invalid Category ID", nil, nil)
		return
	}

	category, err := handler.service.GetByID(id)
	if err != nil {
		helper.ResponseJson(w, false, http.StatusNotFound, err.Error(), nil, nil)
		return
	}

	helper.ResponseJson(w, true, http.StatusOK, "Category Found", category, nil)
}

func (handler CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var categoryData models.Category
	err := json.NewDecoder(r.Body).Decode(&categoryData)
	if err != nil {
		helper.ResponseJson(w, false, http.StatusBadRequest, "Invalid request payload", nil, nil)
		return
	}

	err = handler.service.Create(&categoryData)
	if err != nil {
		helper.ResponseJson(w, false, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	helper.ResponseJson(w, true, http.StatusCreated, "Category created successfully", categoryData, nil)
}

func (handler CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := helper.UriIdConverter(r.URL.Path, "/api/categories/")
	if err != nil {
		helper.ResponseJson(w, false, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	var categoryData models.Category
	err = json.NewDecoder(r.Body).Decode(&categoryData)
	if err != nil {
		helper.ResponseJson(w, false, http.StatusBadRequest, "Invalid request payload", nil, nil)
		return
	}

	categoryData.ID = id

	err = handler.service.Update(&categoryData)
	if err != nil {
		helper.ResponseJson(w, false, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	helper.ResponseJson(w, true, http.StatusOK, "Updated successfully", categoryData, nil)
}

func (handler CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := helper.UriIdConverter(r.URL.Path, "/api/categories/")
	if err != nil {
		helper.ResponseJson(w, false, http.StatusBadRequest, "Invalid Category ID", nil, nil)
		return
	}

	err = handler.service.Delete(id)
	if err != nil {
		helper.ResponseJson(w, false, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	helper.ResponseJson(w, true, http.StatusOK, "Category deleted successfully", nil, nil)
}
