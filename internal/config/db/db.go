package db

import (
	"flag"
	"os"
)

var (
	flagDatabaseDsn = flag.String("d", "postgres://foga:foga@localhost:5432/foga?sslmode=disable", "database connection string")
)

type Config struct {
	DatabaseDSN string
}

func NewConfig() *Config {
	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		*flagDatabaseDsn = envDatabaseDSN
	}

	return &Config{
		DatabaseDSN: *flagDatabaseDsn,
	}
}

func NewConfigFrom(databaseDSN string) *Config {
	return &Config{
		DatabaseDSN: databaseDSN,
	}
}
