package repository

type StorageRepo interface {
	Set(key string, value string) (string, error)
	Get(key string) (string, bool)
}
