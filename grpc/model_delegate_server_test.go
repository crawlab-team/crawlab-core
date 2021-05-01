package grpc

import (
	"context"
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models"
	node2 "github.com/crawlab-team/crawlab-core/node"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestModelDelegateServer_Do_Add(t *testing.T) {
	setupTest(t)

	name := "test-project"
	p := &models.Project{
		Name: name,
	}
	data, err := json.Marshal(p)
	require.Nil(t, err)
	msg := entity.DelegateMessage{
		ModelId: interfaces.ModelIdProject,
		Method:  interfaces.ModelDelegateMethodAdd,
		Data:    data,
	}

	client := TestServiceWorker.MustGetDefaultClient()
	_, err = client.GetModelDelegateClient().Do(context.Background(), &grpc.Request{
		NodeKey: node2.MustGetNodeKey(),
		Data:    msg.ToBytes(),
	})
	require.Nil(t, err)

	p, err = models.MustGetRootService().GetProject(bson.M{"name": name}, nil)
	require.Nil(t, err)
	require.False(t, p.Id.IsZero())
	require.Equal(t, name, p.Name)
}
