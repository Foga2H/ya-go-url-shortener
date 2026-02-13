package storage

import (
	"context"
	"sync"

	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
)

type Link struct {
	OriginalURL string
	UserID      string
}

type MemStorage struct {
	links map[string]Link
	mu    sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		links: make(map[string]Link),
	}
}

func (m *MemStorage) Set(_ context.Context, userID, key, value string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.links[key] = Link{
		OriginalURL: value,
		UserID:      userID,
	}
	return key, nil
}

func (m *MemStorage) Get(_ context.Context, key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.links[key]
	return val.OriginalURL, ok
}

func (m *MemStorage) GetByUserID(_ context.Context, userID string) ([]repository.UserLink, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]repository.UserLink, 0)
	for shortURL, link := range m.links {
		if link.UserID != userID {
			continue
		}

		result = append(result, repository.UserLink{
			ShortURL:    shortURL,
			OriginalURL: link.OriginalURL,
		})
	}

	return result, nil
}
