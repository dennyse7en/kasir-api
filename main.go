package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
	{ID: 1, Name: "Makanan", Description: "Kategori untuk makanan"},
	{ID: 2, Name: "Minuman", Description: "Kategori untuk minuman"},
	{ID: 3, Name: "Lainnya", Description: "Kategori untuk item lainnya"},
}

var nextCategoryID = 4

func main() {
	// ===== Health Check =====
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	// ===== CATEGORY ENDPOINTS =====

	// GET /categories - Ambil semua kategori
	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data":    categories,
			})
			return
		}
		// POST /categories - Tambah kategori
		if r.Method == http.MethodPost {
			var category Category
			err := json.NewDecoder(r.Body).Decode(&category)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"success": "false",
					"error":   "Invalid request body",
				})
				return
			}

			if category.Name == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"success": "false",
					"error":   "Name is required",
				})
				return
			}

			category.ID = nextCategoryID
			nextCategoryID++
			categories = append(categories, category)

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"message": "Category created successfully",
				"data":    category,
			})
			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"success": "false",
			"error":   "Method not allowed",
		})
	})

	// GET /categories/{id} - Ambil detail satu kategori
	// PUT /categories/{id} - Update kategori
	// DELETE /categories/{id} - Hapus kategori
	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		path := strings.TrimPrefix(r.URL.Path, "/categories/")
		id, err := strconv.Atoi(path)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"success": "false",
				"error":   "Invalid category ID",
			})
			return
		}

		// GET /categories/{id}
		if r.Method == http.MethodGet {
			for _, category := range categories {
				if category.ID == id {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"success": true,
						"data":    category,
					})
					return
				}
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"success": "false",
				"error":   "Category not found",
			})
			return
		}

		// PUT /categories/{id}
		if r.Method == http.MethodPut {
			var updateData Category
			err := json.NewDecoder(r.Body).Decode(&updateData)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"success": "false",
					"error":   "Invalid request body",
				})
				return
			}

			for i, category := range categories {
				if category.ID == id {
					if updateData.Name != "" {
						categories[i].Name = updateData.Name
					}
					if updateData.Description != "" {
						categories[i].Description = updateData.Description
					}

					json.NewEncoder(w).Encode(map[string]interface{}{
						"success": true,
						"message": "Category updated successfully",
						"data":    categories[i],
					})
					return
				}
			}

			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"success": "false",
				"error":   "Category not found",
			})
			return
		}

		// DELETE /categories/{id}
		if r.Method == http.MethodDelete {
			for i, category := range categories {
				if category.ID == id {
					categories = append(categories[:i], categories[i+1:]...)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"success": true,
						"message": "Category deleted successfully",
					})
					return
				}
			}

			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"success": "false",
				"error":   "Category not found",
			})
			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"success": "false",
			"error":   "Method not allowed",
		})
	})

	// ===== START SERVER =====

	port := ":8080"
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
