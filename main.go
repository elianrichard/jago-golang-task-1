package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"categories-api/model"
)

func returnJsonResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func getCategoryRequestPayload(w http.ResponseWriter, r *http.Request) (model.Category, bool) {
	var categoryPayload model.Category
	err := json.NewDecoder(r.Body).Decode(&categoryPayload)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return model.Category{}, false
	}
	return categoryPayload, true
}

func getCategoryById(w http.ResponseWriter, r *http.Request) (model.Category, int) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/category/")
	categoryId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return model.Category{}, -1
	}

	for i, category := range DefaultCategories {
		if categoryId == category.ID {
			return category, i
		}
	}
	http.Error(w, "Category Not Found", http.StatusBadRequest)
	return model.Category{}, -1
}

var DefaultCategories = []model.Category{
	{
		ID:          1,
		Name:        "Main Course",
		Description: "This is the main course",
	},
	{
		ID:          2,
		Name:        "Beverage",
		Description: "This is the beverages",
	},
}

func main() {
	http.HandleFunc("/api/v1/category/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			categoryItem, i := getCategoryById(w, r)
			if i < 0 {
				return
			}
			returnJsonResponse(w, categoryItem)
		case "PUT":
			newData, ok := getCategoryRequestPayload(w, r)
			if !ok {
				return
			}
			categoryItem, i := getCategoryById(w, r)
			if i < 0 {
				return
			}
			newData.ID = categoryItem.ID
			DefaultCategories[i] = newData
			returnJsonResponse(w, newData)
		case "DELETE":
			_, i := getCategoryById(w, r)
			if i < 0 {
				return
			}
			DefaultCategories = append(DefaultCategories[:i], DefaultCategories[i+1:]...)
			returnJsonResponse(w, map[string]string{
				"message": "Successfully deleted category",
			})

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/v1/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			returnJsonResponse(w, DefaultCategories)
		case "POST":
			newCategory, ok := getCategoryRequestPayload(w, r)
			if !ok {
				return
			}
			for _, category := range DefaultCategories {
				if newCategory.ID <= category.ID {
					newCategory.ID = category.ID + 1
				}
			}
			DefaultCategories = append(DefaultCategories, newCategory)
			returnJsonResponse(w, newCategory)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		returnJsonResponse(w, map[string]string{
			"status":  "OK",
			"message": "Server Running",
		})
	})

	fmt.Println("Server running in http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Internal Server Error")
		return
	}
}
