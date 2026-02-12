package file

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	"github.com/google/uuid"
)

type Storage struct {
	path string
}

func NewStorage(path string) *Storage {
	return &Storage{
		path: path,
	}
}

type StorageItem struct {
	UUID        string `json:"uuid"`
	UserID      string `json:"user_id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (s *Storage) Set(_ context.Context, userID, key, value string) (string, error) {
	storageFile, err := s.load()
	log.Print(err, err != nil)
	if err != nil {
		return "", err
	}

	storageFile = append(storageFile, StorageItem{
		UUID:        uuid.New().String(),
		UserID:      userID,
		ShortURL:    key,
		OriginalURL: value,
	})

	data, err2 := json.Marshal(storageFile)

	if err2 != nil {
		return "", err2
	}

	err = os.WriteFile(s.path, data, 0666)
	if err != nil {
		return "", err
	}

	return key, nil
}

func (s *Storage) Get(_ context.Context, key string) (string, bool) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return "", false
	}
	var items []StorageItem
	if err := json.Unmarshal(data, &items); err != nil {
		return "", false
	}

	for _, item := range items {
		if item.ShortURL == key {
			return item.OriginalURL, true
		}
	}

	return "", false
}

func (s *Storage) GetByUserID(_ context.Context, userID string) ([]repository.UserLink, error) {
	items, err := s.load()
	if err != nil {
		return nil, err
	}

	result := make([]repository.UserLink, 0)
	for _, item := range items {
		if item.UserID != userID {
			continue
		}

		result = append(result, repository.UserLink{
			ShortURL:    item.ShortURL,
			OriginalURL: item.OriginalURL,
		})
	}

	return result, nil
}

func (s *Storage) load() ([]StorageItem, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []StorageItem{}, nil
		}

		return nil, err
	}

	var items []StorageItem
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}

	return items, nil
}
