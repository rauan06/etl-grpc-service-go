package handler

import (
	"category/internal/core/service"
	"log/slog"
	"net/http"
)

type ETLHandler struct {
	logger *slog.Logger
	svc    service.Pipeline
}

func NewETLHandler(logger *slog.Logger, svc service.Pipeline) *ETLHandler {
	return &ETLHandler{
		logger,
		svc,
	}
}

func (h *ETLHandler) RegisterEndpoint(mux *http.ServeMux) {
	mux.HandleFunc("GET /run", h.Run)
	mux.HandleFunc("GET /status", h.Status)
	mux.HandleFunc("GET /stop", h.Stop)
}

func (h *ETLHandler) Run(w http.ResponseWriter, r *http.Request)    {}
func (h *ETLHandler) Status(w http.ResponseWriter, r *http.Request) {}
func (h *ETLHandler) Stop(w http.ResponseWriter, r *http.Request)   {}
