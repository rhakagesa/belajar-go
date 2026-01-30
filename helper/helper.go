package helper

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func UriIdConverter(path string, prefix string) (int, error) {
	var idStr = strings.TrimPrefix(path, prefix)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func PaginationConverter(r *http.Request) (int, int, error) {
	query := r.URL.Query()

	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	if pageStr == "" || limitStr == "" {
		return 0, 0, errors.New("Missing pagination parameters")
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	return page, limit, nil
}

type Pagination struct {
	CurrentPage int  `json:"current_page"`
	NextPage    *int `json:"next_page"`
	PrevPage    *int `json:"prev_page"`
	TotalPage   int  `json:"total_page"`
	TotalData   int  `json:"total_data"`
	PerPage     int  `json:"per_page"`
}

type BaseResponse struct {
	Status     bool        `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

func ResponseJson(w http.ResponseWriter, status bool, statusCode int, message string, data interface{}, pagination *Pagination) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := BaseResponse{
		Status:     status,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}

	encode := json.NewEncoder(w)
	encode.SetIndent("", " ")
	encode.Encode(response)
}

func BuildPagination(totalData int, limit int, page int) Pagination {
	totalPage := (totalData + limit - 1) / limit

	var next, prev *int
	if page < totalPage {
		nextVal := page + 1
		next = &nextVal
	}
	if page > 1 {
		prevVal := page - 1
		prev = &prevVal
	}

	return Pagination{
		CurrentPage: page,
		NextPage:    next,
		PrevPage:    prev,
		TotalPage:   totalPage,
		TotalData:   totalData,
		PerPage:     limit,
	}
}
