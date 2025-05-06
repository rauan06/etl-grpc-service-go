package handler

import (
	"log"
	"log/slog"
	"net/http"
)

func Routes() *http.ServeMux {
	mux := http.NewServeMux()

	// mux.HandleFunc("/", middleware.Middleware(http.NotFoundHandler().ServeHTTP))

	return mux
}

type APIServer struct {
	address string
	mux     *http.ServeMux
	url     string
	logger  *slog.Logger
}

func NewAPIServer(address string, url string, logger *slog.Logger) *APIServer {
	return &APIServer{
		address: address,
		url:     url,
		mux:     Routes(),
		logger:  logger,
	}
}

func (s *APIServer) Run() {
	// Logging http server initialization
	s.logger.Info("API server listening on " + s.address)

	// #######################
	// Repository Layer
	// #######################
	// categoryRepository := repository.NewCategoryRepository(s.url)

	// #######################
	// Service Layer
	// #######################
	// categoryService := service.NewCategoryService()

	// #######################
	// Presentation Layer
	// #######################
	// categoryHandler := NewCategoryeHandler()

	// #######################
	// Registering Endpoints
	// #######################

	s.logger.Info("API server listening on " + s.address)
	log.Fatal(http.ListenAndServe(s.address, s.mux))
}
