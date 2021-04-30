package controllers

import (
	"github.com/crawlab-team/crawlab-core/models"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/crawlab-team/crawlab-db/redis"
	"github.com/stretchr/testify/require"
	"testing"
)

func setupTest(t *testing.T, cleanup func()) {
	// init mongo
	err := mongo.InitMongo()
	require.Nil(t, err)

	// init redis
	err = redis.InitRedis()
	require.Nil(t, err)

	// init model services
	err = models.InitModels()
	require.Nil(t, err)

	// init controllers
	err = InitControllers()
	require.Nil(t, err)

	// cleanup hook
	t.Cleanup(cleanup)
}
