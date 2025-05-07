package main

import (
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"category/pkg/lib/logger"
)

const (
	apiURL = "apiLink.com"
)

func main() {
	// cfg := config.LoadConfig()

	logger := logger.SetupPrettySlog(os.Stdout)
	slog.SetDefault(logger)

	URL, err := url.Parse(apiURL)
	if err != nil {
		log.Fatalf("Error parsing url: %v", err)
	}

	httpSrv := NewAPIServer(
		"0.0.0.0:8080",
		URL,
		logger,
	)
	httpSrv.Run()
}

func Routes() *http.ServeMux {
	mux := http.NewServeMux()

	// mux.HandleFunc("/", middleware.Middleware(http.NotFoundHandler().ServeHTTP))

	return mux
}

type APIServer struct {
	address string
	mux     *http.ServeMux
	URL     *url.URL
	logger  *slog.Logger
}

func NewAPIServer(address string, url *url.URL, logger *slog.Logger) *APIServer {
	return &APIServer{
		address: address,
		URL:     url,
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
