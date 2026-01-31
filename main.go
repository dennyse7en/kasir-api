package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/config"
	"kasir-api/internal/database"
	"kasir-api/internal/handler"
	"kasir-api/internal/repository"
	"kasir-api/internal/service"
	"log"
	"net/http"
	"os"
)

func main() {
	// Load Configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	if port := os.Getenv("PORT"); port != "" {
		cfg.ServerAddress = ":" + port
	}

	// Initialize Database
	db, err := database.NewPostgresDB(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Layers
	repo := repository.NewCategoryRepository(db)
	svc := service.NewCategoryService(repo)
	h := handler.NewCategoryHandler(svc)

	productRepo := repository.NewProductRepository(db)
	productSvc := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productSvc)

	// Routes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	http.HandleFunc("/categories", h.HandleCategories)
	http.HandleFunc("/categories/", h.HandleCategoryByID)

	http.HandleFunc("/products", productHandler.HandleProducts)
	http.HandleFunc("/products/", productHandler.HandleProductByID)

	// Start Server
	fmt.Printf("Server running on port %s\n", cfg.ServerAddress)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, nil))
}
