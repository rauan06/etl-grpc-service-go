package handler

import (
	"category/internal/core/service"
	pb "category/protos/product/v1/pb"
)

type ProductHandler struct {
	pb.UnimplementedProductServer
	svc service.ProductService
}

func NewProductHandler(svc service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) Run() {
	go h.svc.Run()
}

func (h *ProductHandler) Status() int {
	return h.svc.Status()
}

func (h *ProductHandler) Stop() {
	h.svc.Stop()
}
