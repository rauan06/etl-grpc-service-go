package handler

import (
	"category/internal/core/port"
)

type CategoryHandler struct {
	svc port.CategoryService
}

func NewCategoryeHandler(svc port.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		svc,
	}
}

// func (h *CategoryHandler) RegisterEndpoints(mux *http.ServeMux) {
// }
