package controllers

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/crawlab-team/crawlab-db/redis"
	"github.com/stretchr/testify/require"
	"testing"
)

func setupTest(t *testing.T, cleanup func()) {
	err := mongo.InitMongo()
	require.Nil(t, err)
	err = redis.InitRedis()
	require.Nil(t, err)
	t.Cleanup(cleanup)
}
