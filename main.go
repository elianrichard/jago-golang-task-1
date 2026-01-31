package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/models"
	"kasir-api/repositories"
	"kasir-api/services"
	"kasir-api/utils"

	"github.com/spf13/viper"
)

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

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryById)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.ReturnJsonResponse(w, map[string]string{
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
