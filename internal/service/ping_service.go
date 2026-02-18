package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
)

var ErrPingDatabase = errors.New("ping database")

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --dir . --name DBPinger --output ../mocks --outpkg mocks --filename db_pinger.go --with-expecter --disable-version-string --issue-845-fix
type DBPinger interface {
	PingContext(ctx context.Context) error
}

type PingService struct {
	db     DBPinger
	logger *logger.Logger
}

func NewPingService(db DBPinger, logger *logger.Logger) *PingService {
	return &PingService{db: db, logger: logger}
}

func (s *PingService) Ping(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		s.logger.Errorf("Failed to ping database: %v", err)
		return fmt.Errorf("%w: %v", ErrPingDatabase, err)
	}

	return nil
}
