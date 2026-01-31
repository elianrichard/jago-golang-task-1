package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"kasir-api/utils"
	"net/http"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// Endpoint /api/categories
func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusBadRequest)
	}
}

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAll()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	utils.ReturnJsonResponse(w, categories)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	err = h.service.Create(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	utils.ReturnJsonResponse(w, category)
}

// Endpoint /api/categories/
func (h *CategoryHandler) HandleCategoryById(w http.ResponseWriter, r *http.Request) {
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

func (h *CategoryHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id := utils.GetIdFromRequest(w, r, "/api/categories/")
	category, err := h.service.GetById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	utils.ReturnJsonResponse(w, category)
}

func (h *CategoryHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	id := utils.GetIdFromRequest(w, r, "/api/categories/")
	err := h.service.DeleteById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	utils.ReturnJsonResponse(w, map[string]string{
		"message": "Successfully deleted category",
	})
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	id := utils.GetIdFromRequest(w, r, "/api/categories/")
	category.ID = id
	err = h.service.Update(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	utils.ReturnJsonResponse(w, category)
}
