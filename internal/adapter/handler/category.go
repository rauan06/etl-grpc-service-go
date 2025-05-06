package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"category/internal/core/port"
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
	h.logger.InfoContext(r.Context(), "recieved ListCategories() request")

	params := r.URL.Query()
	if err := validatePaginationParams(params); err != nil {
		h.logger.ErrorContext(r.Context(), "error while parsing parametres for ListCategories()", "error", err)

		WriteError(w, 400, err, "error while parsing parametres")

		return
	}

	categories, err := h.svc.ListCategories(r.Context(), params)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "error while listing categories", "error", err)

		WriteError(w, 500, err, "error while listing categories")

		return
	}

	h.logger.InfoContext(r.Context(), "successfully handled ListCategories() request")
	WriteResponse(w, 200, categories)
}

func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	h.logger.InfoContext(r.Context(), "recieved GetCategory() request")
	rawId := r.PathValue("id")

	id, err := strconv.ParseInt(rawId, 10, 64)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "error while listing categories", "error", err)

		WriteError(w, 400, err, "error while listing categories")

		return
	}

	category, err := h.svc.GetCategory(r.Context(), id)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "error while listing categories", "error", err)

		WriteError(w, 400, err, "error while listing categories")

		return
	}

	h.logger.InfoContext(r.Context(), "successfully handled GetCategory() request")
	WriteResponse(w, 200, category)
}
