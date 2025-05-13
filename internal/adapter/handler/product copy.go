package handler

import (
	"category/internal/core/service"
	pb "category/protos/product/v1/pb"
)

type CategoryHandler struct {
	pb.UnimplementedCategoryServer
	svc service.CategoryService
}

func NewCategoryHandler(svc service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (h *CategoryHandler) Run() {
	go h.svc.Run()
}

func (h *CategoryHandler) Status() int {
	return h.svc.Status()
}

func (h *CategoryHandler) Stop() {
	h.svc.Stop()
}
