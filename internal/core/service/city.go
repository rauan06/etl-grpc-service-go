package service

import (
	"category/internal/core/port"
)

type CityService struct {
	client *port.CityClient
}

func NewCityService(client *port.CityClient) *CityService {
	return &CityService{
		client,
	}
}
