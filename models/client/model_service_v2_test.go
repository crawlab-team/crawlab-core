package client_test

import (
	"context"
	"github.com/crawlab-team/crawlab-core/grpc/server"
	"github.com/crawlab-team/crawlab-core/models/client"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

type TestModel struct {
	Id                            primitive.ObjectID `json:"_id" bson:"_id" collection:"testmodels"`
	models.BaseModelV2[TestModel] `bson:",inline"`
	Name                          string `json:"name" bson:"name"`
}

func setupTestDB() {
	viper.Set("mongo.db", "testdb")
}

func teardownTestDB() {
	db := mongo.GetMongoDb("testdb")
	db.Drop(context.Background())
}

func TestModelServiceV2_GetById(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	m := TestModel{
		Id:   primitive.NewObjectID(),
		Name: "Test Name",
	}
	modelSvc := service.NewModelServiceV2[TestModel]()
	_, err := modelSvc.InsertOne(m)
	require.Nil(t, err)

	go func() {
		svr, err := server.GetServerV2()
		require.Nil(t, err)
		svr.Start()
	}()
	time.Sleep(1 * time.Second)

	c, err := grpc.Dial("localhost:9666", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	c.Connect()

	clientSvc := client.NewModelServiceV2[TestModel]()
	res, err := clientSvc.GetById(m.Id)
	require.Nil(t, err)
	assert.Equal(t, res.Id, m.Id)
	assert.Equal(t, res.Name, m.Name)
}
