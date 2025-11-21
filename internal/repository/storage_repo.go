package repository

type StorageRepo interface {
	Set(key string, value string) error
	Get(key string) (string, bool)
}
