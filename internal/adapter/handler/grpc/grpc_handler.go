package handler

import (
	"category/internal/core/domain"
	"category/internal/core/port"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"sync"
)

type GrpcCategoryHandler struct {
	svc port.CategoryService
}

func NewGrpcCategoryHandler(svc port.CategoryService) *GrpcCategoryHandler {
	return &GrpcCategoryHandler{
		svc,
	}
}

func (h *GrpcCategoryHandler) Run(workers int) {
	limits := make(chan bool, workers)
	results := make(chan domain.CategoryMain)
	wg := new(sync.WaitGroup)

	// Example: simulate reading from 10 endpoints
	var i int64 = 1
	var total int64 = 1000

	for ; i <= total; i++ {
		wg.Add(1)
		limits <- true // throttle

		go func(id int64) {
			defer wg.Done()
			defer func() { <-limits }()

			category, err := h.svc.GetCategory(context.Background(), id)
			if err != nil {
				log.Printf("Error fetching %s: %v\n", url, err)
				return
			}

			body, _ := ioutil.ReadAll(resp.Body)

			// Simulate unmarshalling to CategoryMain
			results <- domain.CategoryMain{
				ID:   strconv.FormatInt(i, 10),
				Name: string(body),
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// Process results
	for category := range results {
		fmt.Printf("Received: %+v\n", category)
	}
}

func (h *GrpcCategoryHandler) Read(workers int) error {
	return nil
}

// Assuming it will write to the database, will be implemented soon
func (h *GrpcCategoryHandler) Write() error { return nil }
