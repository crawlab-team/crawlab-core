package server

import (
	"context"
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/entity"
	grpc2 "github.com/crawlab-team/crawlab-core/grpc"
	"github.com/crawlab-team/crawlab-core/interfaces"
	models2 "github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/utils"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestModelDelegateServer_Do_Add(t *testing.T) {
	name := "test-project"
	p := &models2.Project{
		Name: name,
	}
	data, err := json.Marshal(p)
	require.Nil(t, err)
	msg := entity.DelegateMessage{
		ModelId: interfaces.ModelIdProject,
		Method:  interfaces.ModelDelegateMethodAdd,
		Data:    data,
	}

	client := grpc2.TestServiceWorker.GetClient()
	_, err = client.GetModelDelegateClient().Do(context.Background(), &grpc.Request{
		NodeKey: grpc2.TestServiceWorker.nodeSvc.GetNodeKey(),
		Data:    msg.ToBytes(),
	})
	require.Nil(t, err)

	p, err = modelSvc.GetProject(bson.M{"name": name}, nil)
	require.Nil(t, err)
	require.False(t, p.Id.IsZero())
	require.Equal(t, name, p.Name)

	a, err := modelSvc.GetArtifactById(p.Id)
	require.Nil(t, err)
	require.Equal(t, p.Id, a.Id)
}

func TestModelDelegateServer_Do_Save(t *testing.T) {
	var modelSvc service.ModelService
	utils.MustResolveModule("", modelSvc)

	var err error
	grpc2.setupTest(t)

	oldName := "old-name"
	p := &models2.Project{
		Name: oldName,
	}
	err = p.Add()
	require.Nil(t, err)

	newName := "new-name"
	newTags := []models2.Tag{
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

	client := grpc2.TestServiceWorker.GetClient()
	_, err = client.GetModelDelegateClient().Do(context.Background(), &grpc.Request{
		NodeKey: grpc2.TestServiceWorker.nodeSvc.GetNodeKey(),
		Data:    msg.ToBytes(),
	})
	require.Nil(t, err)

	p, err = modelSvc.GetProjectById(p.Id)
	require.Nil(t, err)
	require.False(t, p.Id.IsZero())
	require.NotEmpty(t, p.Tags)
	require.False(t, p.Tags[0].GetId().IsZero())
}
