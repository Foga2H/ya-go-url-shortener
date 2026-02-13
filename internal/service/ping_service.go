package service

import (
	"context"
	"errors"
	"fmt"
)

var ErrPingDatabase = errors.New("ping database")

type DBPinger interface {
	PingContext(ctx context.Context) error
}

type PingService struct {
	db DBPinger
}

func NewPingService(db DBPinger) *PingService {
	return &PingService{db: db}
}

func (s *PingService) Ping(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("%w: %v", ErrPingDatabase, err)
	}

	return nil
}
