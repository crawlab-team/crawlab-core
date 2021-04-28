package grpc

import (
	"github.com/crawlab-team/crawlab-core/models"
	"testing"
)

func setupTest(t *testing.T) {
	_ = models.InitModelServices()
	t.Cleanup(cleanupTest)
}

func cleanupTest() {
}
