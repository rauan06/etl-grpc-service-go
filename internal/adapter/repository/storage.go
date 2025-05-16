package repository

import (
	"sync"

	"github.com/rauan06/etl-grpc-service-go/internal/core/domain"
)

type Repositry struct {
	storage chan domain.MarketPair

	results map[string]domain.FullProduct
	mu      sync.Mutex
}

func NewRepositry() *Repositry {
	return &Repositry{
		storage: make(chan domain.MarketPair),
		results: make(map[string]domain.FullProduct),
	}
}

func (s *Repositry) SavePair(pair domain.MarketPair) {
	s.storage <- pair
}

func (s *Repositry) ReadPair() (domain.MarketPair, bool) {
	pair, ok := <-s.storage
	return pair, ok
}

func (s *Repositry) SaveResult(key string, value domain.FullProduct) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.results[key] = value
}

func (s *Repositry) GetResults() []domain.FullProduct {
	s.mu.Lock()
	defer s.mu.Unlock()

	resp := []domain.FullProduct{}
	for _, val := range s.results {
		resp = append(resp, val)
	}

	return resp
}
