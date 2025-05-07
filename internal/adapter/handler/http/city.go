package handler

import "category/internal/core/port"

type CityHandler struct {
	svc port.CityService
}

func NewCityeHandler(svc port.CityService) *CityHandler {
	return &CityHandler{
		svc,
	}
}
