package test

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"testing"
)

func TestMain(m *testing.M) {
	if err := mongo.InitMongo(); err != nil {
		panic(err)
	}

	m.Run()
}
