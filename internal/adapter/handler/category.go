package handler

import (
	"category/internal/core/domain"
	"category/internal/core/port"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type CategoryHandler struct {
	svc    port.CategoryService
	logger *slog.Logger
}

func NewCategoryHandler(svc port.CategoryService, logger *slog.Logger) *CategoryHandler {
	return &CategoryHandler{
		svc,
		logger,
	}
}

func (h *CategoryHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("GET /category", h.ListCategories)
	mux.HandleFunc("GET /category/{id}", h.GetCategory)
}

func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	pageStr := query.Get("list_params.page")
	var page int64
	if pageStr != "" {
		p, err := strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid list_params.page", http.StatusBadRequest)
			return
		}
		page = p
	}

	pageSizeStr := query.Get("list_params.page_size")
	var pageSize int64
	if pageSizeStr != "" {
		ps, err := strconv.ParseInt(pageSizeStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid list_params.page_size", http.StatusBadRequest)
			return
		}
		pageSize = ps
	}

	sortParams := query["list_params.sort"]

	// Parse ids (multi)
	ids := query["ids"]

	params := domain.ProductListParamsSt{
		Page:     page,
		PageSize: pageSize,
		Sort:     sortParams,
	}

	// Call service/repository layer
	categories, err := h.svc.ListCategories(r.Context(), params, ids)
	if err != nil {
		http.Error(w, "Failed to get categories", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	rawId := r.PathValue("id")
	id, err := strconv.ParseInt(rawId, 10, 64)
	if err != nil {
		h.logger.Error("Error converting to int", "error", err)
		http.Error(w, "Error converting to int", http.StatusInternalServerError)
		return
	}

	category, err := h.svc.GetCategory(context.Background(), int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(category); err != nil {
		h.logger.Error("Error encoding response to json", "error", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
