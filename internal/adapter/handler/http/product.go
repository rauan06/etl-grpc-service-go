package handler

import (
	"category/internal/core/port"
	"log/slog"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	svc    port.ProductService
	logger *slog.Logger
}

func NewProducteHandler(svc port.ProductService, logger *slog.Logger) *ProductHandler {
	return &ProductHandler{
		svc,
		logger,
	}
}

func (h *ProductHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("GET /product", h.ListPRoducts)
	mux.HandleFunc("GET /product/{id}", h.GetProduct)
}

func (h *ProductHandler) ListPRoducts(w http.ResponseWriter, r *http.Request) {
	h.logger.InfoContext(r.Context(), "recieved ListProducts() request")

	params := r.URL.Query()
	if err := validatePaginationParams(params); err != nil {
		h.logger.ErrorContext(r.Context(), "error while parsing parametres for ListProducts()", "error", err)

		WriteError(w, 400, err, "error while parsing parametres")

		return
	}

	categoryBool := params.Get("with_category")
	if categoryBool != "true" {
		params.Add("with_category", "false")
	}

	products, err := h.svc.ListProducts(r.Context(), params)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "error while listing products", "error", err)

		WriteError(w, 500, err, "error while listing products")

		return
	}

	h.logger.InfoContext(r.Context(), "successfully handled ListProducts() request")
	WriteResponse(w, 200, products)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	h.logger.InfoContext(r.Context(), "recieved GetProduct() request")
	rawId := r.PathValue("id")

	id, err := strconv.ParseInt(rawId, 10, 64)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "error while listing products", "error", err)

		WriteError(w, 400, err, "error while listing products")

		return
	}

	Product, err := h.svc.GetProduct(r.Context(), id)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "error while listing products", "error", err)

		WriteError(w, 400, err, "error while listing products")

		return
	}

	h.logger.InfoContext(r.Context(), "successfully handled GetProduct() request")
	WriteResponse(w, 200, Product)
}
