package storage

import (
	"context"
	"sync"

	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
)

type Link struct {
	OriginalURL string
	UserID      string
	IsDeleted   bool
}

type MemStorage struct {
	links map[string]Link
	mu    sync.RWMutex
}

func (m *MemStorage) BatchDelete(_ context.Context, userID string, links []string) error {
	if len(links) == 0 {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, shortURL := range links {
		link, ok := m.links[shortURL]
		if !ok || link.UserID != userID {
			continue
		}

		link.IsDeleted = true
		m.links[shortURL] = link
	}

	return nil
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

func (m *MemStorage) Get(_ context.Context, key string) (repository.UserLink, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.links[key]
	if !ok {
		return repository.UserLink{}, false
	}

	return repository.UserLink{
		ShortURL:    key,
		OriginalURL: val.OriginalURL,
		IsDeleted:   val.IsDeleted,
	}, true
}

func (m *MemStorage) GetByUserID(_ context.Context, userID string) ([]repository.UserLink, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]repository.UserLink, 0)
	for shortURL, link := range m.links {
		if link.UserID != userID {
			continue
		}
		if link.IsDeleted {
			continue
		}

		result = append(result, repository.UserLink{
			ShortURL:    shortURL,
			OriginalURL: link.OriginalURL,
		})
	}

	return result, nil
}
