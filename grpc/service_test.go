package grpc

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func setupService() (err error) {
	return nil
}

func cleanupService() {
}

func TestNewCrawlabGrpcService(t *testing.T) {
	err := setupService()
	require.Nil(t, err)

	s, err := NewCrawlabGrpcService()
	require.Nil(t, err)
	require.NotNil(t, s)

	cleanupService()
}

func TestCrawlabGrpcService_Init(t *testing.T) {
	err := setupService()
	require.Nil(t, err)

	s, err := NewCrawlabGrpcService()
	require.Nil(t, err)

	// test init
	isStopped := false
	go func() {
		err = s.Init()
		require.Nil(t, err)
		isStopped = true
	}()

	// stop
	s.Stop()

	time.Sleep(1 * time.Second)
	require.True(t, isStopped)

	cleanupService()
}
