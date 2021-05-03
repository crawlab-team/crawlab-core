package grpc

import (
	"context"
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models"
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

	client := TestServiceWorker.GetClient()
	_, err = client.GetModelDelegateClient().Do(context.Background(), &grpc.Request{
		NodeKey: TestServiceWorker.nodeSvc.GetNodeKey(),
		Data:    msg.ToBytes(),
	})
	require.Nil(t, err)

	p, err = models.MustGetRootService().GetProject(bson.M{"name": name}, nil)
	require.Nil(t, err)
	require.False(t, p.Id.IsZero())
	require.Equal(t, name, p.Name)

	a, err := models.MustGetRootService().GetArtifactById(p.Id)
	require.Nil(t, err)
	require.Equal(t, p.Id, a.Id)
}

func TestModelDelegateServer_Do_Save(t *testing.T) {
	var err error
	setupTest(t)

	oldName := "old-name"
	p := &models.Project{
		Name: oldName,
	}
	err = p.Add()
	require.Nil(t, err)

	newName := "new-name"
	newTags := []models.Tag{
		{Name: "new-tag", Color: "red"},
	}
	p.Name = newName
	p.Tags = newTags
	data, err := json.Marshal(p)
	require.Nil(t, err)
	msg := entity.DelegateMessage{
		ModelId: interfaces.ModelIdProject,
		Method:  interfaces.ModelDelegateMethodSave,
		Data:    data,
	}

	client := TestServiceWorker.GetClient()
	_, err = client.GetModelDelegateClient().Do(context.Background(), &grpc.Request{
		NodeKey: TestServiceWorker.nodeSvc.GetNodeKey(),
		Data:    msg.ToBytes(),
	})
	require.Nil(t, err)

	p, err = models.MustGetRootService().GetProjectById(p.Id)
	require.Nil(t, err)
	require.False(t, p.Id.IsZero())
	require.NotEmpty(t, p.Tags)
	require.False(t, p.Tags[0].GetId().IsZero())
}
