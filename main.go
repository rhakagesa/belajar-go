package main

import (
	"belajar-go/databases"
	"belajar-go/handler"
	"belajar-go/repositories"
	"belajar-go/services"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port    string `mapstructure:"PORT"`
	DB      string `mapstructure:"DB_CONN"`
	API_URL string `mapstructure:"API_URL"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port: viper.GetString("PORT"),
		DB:   viper.GetString("DB_CONN"),
	}

	db, err := databases.InitDB(config.DB)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	_, err = databases.Migrate(db)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	defer db.Close()

	categoryRepository := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepository)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	productRepository := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepository, categoryRepository)
	productHandler := handler.NewProductHandler(productService)

	transactionRepository := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepository, productRepository)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout)
	http.HandleFunc("/api/reports/today", transactionHandler.HandleReport)
	http.HandleFunc("/api/reports", transactionHandler.HandleReport)
	http.HandleFunc("/api/categories", categoryHandler.HandleCategory)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)
	http.HandleFunc("/api/products", productHandler.HandleProduct)
	http.HandleFunc("/api/products/", productHandler.HandleProductByID)

	fmt.Println("Server started on " + config.API_URL + ":" + config.Port)
	err = http.ListenAndServe(":"+config.Port, nil)

	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
