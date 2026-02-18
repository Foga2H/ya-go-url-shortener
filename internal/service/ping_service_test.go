package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	"github.com/Foga2H/ya-go-url-shortener/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPingService_PingError(t *testing.T) {
	pinger := mocks.NewDBPinger(t)
	pinger.EXPECT().PingContext(mock.Anything).Return(errors.New("unavailable"))
	svc := NewPingService(pinger, logger.NewLogger())

	err := svc.Ping(context.Background())
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrPingDatabase)
}
