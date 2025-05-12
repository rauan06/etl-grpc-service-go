package service

import (
	"category/internal/core/port"
)

type Pipeline struct {
	grpcClient port.CategoryClient
	httpClient port.CategoryClient
}

func NewPipeline(grpcClient port.CategoryClient, httpClient port.CategoryClient) *Pipeline {
	return &Pipeline{
		grpcClient,
		httpClient,
	}
}

func (p *Pipeline) Run() error { return nil }

func (p *Pipeline) Status() {}

func (p *Pipeline) Stop() {}
