package models

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"testing"
)

func setupTest(t *testing.T, cleanup func()) {
	_ = mongo.InitMongo()
	t.Cleanup(cleanup)
}
