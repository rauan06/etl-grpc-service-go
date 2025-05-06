package handler

import (
	"category/internal/core/port"
	"net/http"
)

type CategoryHandler struct {
	svc port.CategoryService
}

func NewCategoryHandler(svc port.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		svc,
	}
}

func (h *CategoryHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("GET /category", h.ListCategories)
	mux.HandleFunc("GET /category/{id}", h.GetCategory)
}

func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
}

func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {

}
