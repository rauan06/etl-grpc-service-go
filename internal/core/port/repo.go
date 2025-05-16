package port

import "github.com/rauan06/etl-grpc-service-go/internal/core/domain"

type Repository interface {
	SavePair(domain.MarketPair)
	ReadPair() (domain.MarketPair, bool)
}
