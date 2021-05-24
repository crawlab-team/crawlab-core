package test

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/crawlab-team/crawlab-db/redis"
	"testing"
)

func TestMain(m *testing.M) {
	if err := redis.InitRedis(); err != nil {
		panic(err)
	}
	if err := mongo.InitMongo(); err != nil {
		panic(err)
	}
	m.Run()
}
