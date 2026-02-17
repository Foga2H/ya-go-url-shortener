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

func (s *Storage) BatchDelete(_ context.Context, userID string, links []string) error {
	if len(links) == 0 {
		return nil
	}

	items, err := s.load()
	if err != nil {
		return err
	}

	toDelete := make(map[string]struct{}, len(links))
	for _, shortURL := range links {
		toDelete[shortURL] = struct{}{}
	}

	for i := range items {
		if items[i].UserID != userID {
			continue
		}
		if _, ok := toDelete[items[i].ShortURL]; !ok {
			continue
		}

		items[i].IsDeleted = true
	}

	data, err := json.Marshal(items)
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, data, 0666)
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
	IsDeleted   bool   `json:"is_deleted"`
}

func (s *Storage) Set(_ context.Context, userID, key, value string) (string, error) {
	storageFile, err := s.load()
	log.Print(err, err != nil)
	if err != nil {
		return "", err
	}

	storageFile = append(storageFile, StorageItem{
		UUID:        uuid.NewString(),
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

func (s *Storage) Get(_ context.Context, key string) (repository.UserLink, bool) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return repository.UserLink{}, false
	}
	var items []StorageItem
	if err := json.Unmarshal(data, &items); err != nil {
		return repository.UserLink{}, false
	}

	for _, item := range items {
		if item.ShortURL == key {
			return repository.UserLink{
				ShortURL:    item.ShortURL,
				OriginalURL: item.OriginalURL,
				IsDeleted:   item.IsDeleted,
			}, true
		}
	}

	return repository.UserLink{}, false
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
		if item.IsDeleted {
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
