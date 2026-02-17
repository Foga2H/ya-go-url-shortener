package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Foga2H/ya-go-url-shortener/internal/config"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
)

type DeleteUserURLsService struct {
	storage    repository.StorageRepo
	config     *config.Config
	in         chan []string
	userID     string
	flushEvery time.Duration
	maxBatch   int
}

func NewDeleteUserURLsService(storage repository.StorageRepo, config *config.Config) *DeleteUserURLsService {
	return &DeleteUserURLsService{
		storage:    storage,
		config:     config,
		in:         make(chan []string, 250),
		flushEvery: 300 * time.Millisecond,
		maxBatch:   300,
	}
}

func (d *DeleteUserURLsService) Enqueue(userId string, shortUrls []string) {
	d.userID = userId
	d.in <- shortUrls
}

func (d *DeleteUserURLsService) Start(ctx context.Context) {
	go d.run(ctx)
}

func (d *DeleteUserURLsService) run(ctx context.Context) {
	ticker := time.NewTicker(d.flushEvery)
	defer ticker.Stop()

	pending := make(map[string]struct{}, d.maxBatch)

	flush := func() {
		if len(pending) == 0 {
			return
		}

		batch := make([]string, 0, len(pending))

		for u := range pending {
			batch = append(batch, u)
		}

		pending = make(map[string]struct{}, d.maxBatch)

		err := d.storage.BatchDelete(ctx, d.userID, batch)
		if err != nil {
			fmt.Printf("Error deleting URLs for user %s: %v\n", d.userID, strings.Join(batch, ","))
		}
	}

	for {
		select {
		case <-ctx.Done():
			flush()
			return

		case urls := <-d.in:
			for _, u := range urls {
				pending[u] = struct{}{}
				if len(pending) >= d.maxBatch {
					flush()
				}
			}

		case <-ticker.C:
			flush()
		}
	}
}
