package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"kasir-api/utils"
	"net/http"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// Endpoint /api/products
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusBadRequest)
	}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	utils.ReturnJsonResponse(w, products)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	err = h.service.Create(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	utils.ReturnJsonResponse(w, product)
}

// Endpoint /api/products/
func (h *ProductHandler) HandleProductById(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetById(w, r)
	case http.MethodDelete:
		h.DeleteById(w, r)
	case http.MethodPut:
		h.Update(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusBadRequest)
	}

}

func (h *ProductHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id := utils.GetIdFromRequest(w, r, "/api/products/")
	product, err := h.service.GetById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	utils.ReturnJsonResponse(w, product)
}

func (h *ProductHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	id := utils.GetIdFromRequest(w, r, "/api/products/")
	err := h.service.DeleteById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	utils.ReturnJsonResponse(w, map[string]string{
		"message": "Successfully deleted product",
	})
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	id := utils.GetIdFromRequest(w, r, "/api/products/")
	product.ID = id
	err = h.service.Update(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	utils.ReturnJsonResponse(w, product)
}
