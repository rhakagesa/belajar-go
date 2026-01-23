package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Categories struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var categories = []Categories{
	{ID: 1, Name: "Category 1", Description: "Description 1"},
	{ID: 2, Name: "Category 2", Description: "Description 2"},
	{ID: 3, Name: "Category 3", Description: "Description 3"},
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory Categories
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	newCategory.ID = len(categories) + 1
	categories = append(categories, newCategory)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCategory)
}

func getCategory(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func getCategoryById(w http.ResponseWriter, r *http.Request) {
	var idStr = strings.TrimPrefix(r.URL.Path, "/api/categories/")
	var id, err = strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	for _, category := range categories {
		if category.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(category)
			return
		}
	}

	http.Error(w, "Category Not Found", http.StatusNotFound)
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	var idStr = strings.TrimPrefix(r.URL.Path, "/api/categories/")
	var id, err = strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	var updatedCategory Categories
	err = json.NewDecoder(r.Body).Decode(&updatedCategory)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	for i := range categories {
		if categories[i].ID == id {
			updatedCategory.ID = id
			categories[i] = updatedCategory

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories[i])
			return
		}
	}

	http.Error(w, "Category Not Found", http.StatusNotFound)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	var idStr = strings.TrimPrefix(r.URL.Path, "/api/categories/")
	var id, err = strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	for i := range categories {
		if categories[i].ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Category Deleted"})
			return
		}
	}
}

func main() {
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getCategory(w)
		case "POST":
			createCategory(w, r)
		default:
			getCategory(w)
		}
	})

	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getCategoryById(w, r)
		case "PUT":
			updateCategory(w, r)
		case "DELETE":
			deleteCategory(w, r)
		}
	})

	fmt.Println("Server started on port http://localhost:9000")
	http.ListenAndServe(":9000", nil)
}
