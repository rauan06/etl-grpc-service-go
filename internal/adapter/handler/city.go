package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"category/internal/core/port"
)

type CityHandler struct {
	svc    port.CityService
	logger *slog.Logger
}

func NewCityHandler(svc port.CityService, logger *slog.Logger) *CityHandler {
	return &CityHandler{
		svc,
		logger,
	}
}

func (h *CityHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("GET /City", h.ListCities)
	mux.HandleFunc("GET /City/{id}", h.GetCity)
}

func (h *CityHandler) ListCities(w http.ResponseWriter, r *http.Request) {
	h.logger.InfoContext(r.Context(), "recieved ListCities() request")

	params := r.URL.Query()
	if err := validatePaginationParams(params); err != nil {
		h.logger.ErrorContext(r.Context(), "error while parsing parametres for ListCities()", "error", err)

		WriteError(w, 400, err, "error while parsing parametres")

		return
	}

	cities, err := h.svc.ListCities(r.Context(), params)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "error while listing cities", "error", err)

		WriteError(w, 500, err, "error while listing cities")

		return
	}

	h.logger.InfoContext(r.Context(), "successfully handled ListCities() request")
	WriteResponse(w, 200, cities)
}

func (h *CityHandler) GetCity(w http.ResponseWriter, r *http.Request) {
	h.logger.InfoContext(r.Context(), "recieved GetCity() request")
	rawId := r.PathValue("id")

	id, err := strconv.ParseInt(rawId, 10, 64)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "error while listing cities", "error", err)

		WriteError(w, 400, err, "error while listing cities")

		return
	}

	City, err := h.svc.GetCity(r.Context(), id)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "error while listing cities", "error", err)

		WriteError(w, 400, err, "error while listing cities")

		return
	}

	h.logger.InfoContext(r.Context(), "successfully handled GetCity() request")
	WriteResponse(w, 200, City)
}
