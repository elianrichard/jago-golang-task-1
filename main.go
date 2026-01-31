package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/models"
	"kasir-api/repositories"
	"kasir-api/services"

	"github.com/spf13/viper"
)

func returnJsonResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func getCategoryRequestPayload(w http.ResponseWriter, r *http.Request) (models.Category, bool) {
	var categoryPayload models.Category
	err := json.NewDecoder(r.Body).Decode(&categoryPayload)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return models.Category{}, false
	}
	return categoryPayload, true
}

func getCategoryById(w http.ResponseWriter, r *http.Request) (models.Category, int) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	categoryId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return models.Category{}, -1
	}

	for i, category := range DefaultCategories {
		if categoryId == category.ID {
			return category, i
		}
	}
	http.Error(w, "Category Not Found", http.StatusBadRequest)
	return models.Category{}, -1
}

var DefaultCategories = []models.Category{
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
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := models.Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DATABASE_URL"),
	}

	log.Println("Connecting to database...")
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	http.HandleFunc("/api/products", productHandler.HandleProducts)
	http.HandleFunc("/api/products/", productHandler.HandleProductById)

	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
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

	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
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

	addr := ":" + config.Port
	fmt.Println("Server running in http", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Internal Server Error")
		return
	}
}
