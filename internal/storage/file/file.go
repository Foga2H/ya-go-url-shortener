package file

import (
	"encoding/json"
	"log"
	"os"

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
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (s *Storage) Set(key string, value string) (string, error) {
	storageFile, err := s.load()
	log.Print(err, err != nil)
	if err != nil {
		return "", err
	}

	storageFile = append(storageFile, StorageItem{
		UUID:        uuid.New().String(),
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

func (s *Storage) Get(key string) (string, bool) {
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
