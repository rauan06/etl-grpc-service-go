package storage

import "github.com/rauan06/etl-grpc-service-go/internal/core/domain"

type Storage struct {
	storage chan domain.MarketPair
}

func NewStorage() *Storage {
	return &Storage{
		storage: make(chan domain.MarketPair),
	}
}

func (s *Storage) SavePair(domain.MarketPair) {

}

func (s *Storage) ReadPair() (domain.MarketPair, bool) {
	pair, ok := <-s.storage
	return pair, ok
}
