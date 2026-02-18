package service

import (
	"context"
	"strings"
	"time"

	"github.com/Foga2H/ya-go-url-shortener/internal/config"
	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
)

type DeleteUserURLsService struct {
	storage    repository.StorageRepo
	config     *config.Config
	logger     *logger.Logger
	in         chan []string
	userID     string
	flushEvery time.Duration
	maxBatch   int
}

func NewDeleteUserURLsService(storage repository.StorageRepo, config *config.Config, logger *logger.Logger) *DeleteUserURLsService {
	return &DeleteUserURLsService{
		storage:    storage,
		config:     config,
		logger:     logger,
		in:         make(chan []string, 250),
		flushEvery: 300 * time.Millisecond,
		maxBatch:   300,
	}
}

func (d *DeleteUserURLsService) Enqueue(userID string, shortUrls []string) {
	d.userID = userID
	d.in <- shortUrls
}

func (d *DeleteUserURLsService) Start(ctx context.Context) {
	go d.run(ctx)
}

func (d *DeleteUserURLsService) run(ctx context.Context) {
	ticker := time.NewTicker(d.flushEvery)
	defer ticker.Stop()

	pending := make(map[string]struct{}, d.maxBatch)

	for {
		select {
		case <-ctx.Done():
			pending = d.flushPending(ctx, pending)
			return

		case urls := <-d.in:
			for _, u := range urls {
				pending[u] = struct{}{}
				if len(pending) >= d.maxBatch {
					pending = d.flushPending(ctx, pending)
				}
			}

		case <-ticker.C:
			pending = d.flushPending(ctx, pending)
		}
	}
}

func (d *DeleteUserURLsService) flushPending(ctx context.Context, pending map[string]struct{}) map[string]struct{} {
	if len(pending) == 0 {
		return pending
	}

	batch := make([]string, 0, len(pending))
	for u := range pending {
		batch = append(batch, u)
	}

	nextPending := make(map[string]struct{}, d.maxBatch)

	err := d.storage.BatchDelete(ctx, d.userID, batch)
	if err != nil {
		d.logger.Errorf("Error deleting URLs for user %s, urls=%s: %v", d.userID, strings.Join(batch, ","), err)
	}

	return nextPending
}
