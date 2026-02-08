package handler

import (
	"encoding/json"
	"kasir-api/internal/model"
	"kasir-api/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		nameQuery := r.URL.Query().Get("name")
		var products []model.Product
		var err error

		if nameQuery != "" {
			products, err = h.service.SearchByName(nameQuery)
		} else {
			products, err = h.service.GetAll()
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": products})
		return
	}

	if r.Method == http.MethodPost {
		var product model.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Invalid request body"})
			return
		}

		createdProduct, err := h.service.Create(product)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Product created successfully", "data": createdProduct})
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Method not allowed"})
}

func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := strings.TrimPrefix(r.URL.Path, "/products/")
	idCursor := strings.Split(path, "/")
	if len(idCursor) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idCursor[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Invalid product ID"})
		return
	}

	if r.Method == http.MethodGet {
		product, err := h.service.GetByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Product not found"})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": product})
		return
	}

	if r.Method == http.MethodPut {
		var product model.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Invalid request body"})
			return
		}

		updatedProduct, err := h.service.Update(id, product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Product updated successfully", "data": updatedProduct})
		return
	}

	if r.Method == http.MethodDelete {
		if err := h.service.Delete(id); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Product deleted successfully"})
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Method not allowed"})
}
