package handler

import (
	"belajar-go/helper"
	"belajar-go/models"
	"belajar-go/services"
	"encoding/json"
	"net/http"
	"time"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (handler *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handler.Checkout(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (handler *TransactionHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handler.GetReportSales(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (handler *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var request models.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		helper.ResponseJson(w, false, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	transaction, err := handler.service.Checkout(request.Items, true)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case models.ErrProductNotFound:
			statusCode = http.StatusNotFound
		case models.ErrProductStockNotEnough:
			statusCode = http.StatusBadRequest
		case models.ErrQuantityNotValid:
			statusCode = http.StatusBadRequest
		}
		helper.ResponseJson(w, false, statusCode, err.Error(), nil, nil)
		return
	}

	helper.ResponseJson(w, true, 200, "Checkout items successfully", transaction, nil)
}

func (handler *TransactionHandler) GetReportSales(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	isToday := r.URL.Path == "/api/reports/today"

	layout := "2006-01-02"
	now := time.Now().UTC()

	if isToday {
		startDate = now.Format(layout)
		endDate = now.Format(layout)
	} else if startDate == "" && endDate == "" {
		startDate = now.AddDate(0, 0, -30).Format(layout)
		endDate = now.Format(layout)
	}

	if startDate > endDate {
		helper.ResponseJson(w, false, http.StatusBadRequest, "Start date must be less than end date", nil, nil)
		return
	}

	reports, err := handler.service.GetReportSales(startDate, endDate)
	if err != nil {
		helper.ResponseJson(w, false, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	helper.ResponseJson(w, true, http.StatusOK, "Report sales successfully", reports, nil)
}
