package handler

import (
	"belajar-go/helper"
	"belajar-go/models"
	"belajar-go/services"
	"encoding/json"
	"net/http"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (handler ProductHandler) HandleProduct(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handler.GetAll(w, r)
	case http.MethodPost:
		handler.Create(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
func (handler ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
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

func (handler ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, limit, err := helper.PaginationConverter(r)

	if err != nil {
		helper.ResponseJson(w, false, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	products, totalData, err := handler.service.GetAll(page, limit)
	if err != nil {
		helper.ResponseJson(w, false, http.StatusInternalServerError, err.Error(), nil, nil)
	}

	if totalData == 0 {
		helper.ResponseJson(w, true, http.StatusOK, "No Products Found", nil, nil)
		return
	}

	pagination := helper.BuildPagination(totalData, limit, page)

	helper.ResponseJson(w, true, http.StatusOK, "Products Found", products, &pagination)
}

func (handler ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := helper.UriIdConverter(r.URL.Path, "/api/products/")
	if err != nil {
		helper.ResponseJson(w, false, http.StatusBadRequest, "Invalid Product ID", nil, nil)
		return
	}

	product, err := handler.service.GetByID(id)
	if err != nil {
		helper.ResponseJson(w, false, http.StatusNotFound, err.Error(), nil, nil)
		return
	}

	helper.ResponseJson(w, true, http.StatusOK, "Product Found", product, nil)
}

func (handler ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var productData models.Product
	err := json.NewDecoder(r.Body).Decode(&productData)
	if err != nil {
		helper.ResponseJson(w, false, http.StatusBadRequest, "Invalid request payload", nil, nil)
		return
	}

	err = handler.service.Create(&productData)
	if err != nil {
		helper.ResponseJson(w, false, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	helper.ResponseJson(w, true, http.StatusCreated, "Product created successfully", productData, nil)
}

func (handler ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := helper.UriIdConverter(r.URL.Path, "/api/products/")
	if err != nil {
		helper.ResponseJson(w, false, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	var productData models.Product
	err = json.NewDecoder(r.Body).Decode(&productData)
	if err != nil {
		helper.ResponseJson(w, false, http.StatusBadRequest, "Invalid request payload", nil, nil)
		return
	}

	productData.ID = id

	err = handler.service.Update(&productData)
	if err != nil {
		helper.ResponseJson(w, false, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	helper.ResponseJson(w, true, http.StatusOK, "Updated successfully", productData, nil)
}

func (handler ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := helper.UriIdConverter(r.URL.Path, "/api/products/")
	if err != nil {
		helper.ResponseJson(w, false, http.StatusBadRequest, "Invalid Product ID", nil, nil)
		return
	}

	err = handler.service.Delete(id)
	if err != nil {
		helper.ResponseJson(w, false, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	helper.ResponseJson(w, true, http.StatusOK, "Product deleted successfully", nil, nil)
}
