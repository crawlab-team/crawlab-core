package service_test

import (
	"fmt"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	models2 "github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestNodeService_GetModelById(t *testing.T) {
	SetupTest(t)

	node := &models2.Node{
		Name:     "test node",
		IsMaster: true,
		Tags: []models2.Tag{
			{Name: "tag 1", Color: "red"},
		},
	}
	err := delegate.NewModelDelegate(node).Add()
	require.Nil(t, err)

	svc, err := service.NewService()
	require.Nil(t, err)

	node, err = svc.GetNodeById(node.Id)
	require.Nil(t, err)
	require.NotNil(t, node.Tags)
	require.False(t, node.Id.IsZero())
	require.False(t, node.Tags[0].Id.IsZero())
	require.Equal(t, "tag 1", node.Tags[0].Name)
	require.Equal(t, "red", node.Tags[0].Color)
	require.Equal(t, interfaces.ModelColNameNode, node.Tags[0].Col)
}

func TestNodeService_GetModel(t *testing.T) {
	SetupTest(t)

	node := &models2.Node{
		Name:     "test node",
		IsMaster: true,
		Tags: []models2.Tag{
			{Name: "tag 1", Color: "red"},
		},
	}
	err := delegate.NewModelDelegate(node).Add()
	require.Nil(t, err)

	svc, err := service.NewService()
	require.Nil(t, err)

	node, err = svc.GetNode(bson.M{"name": "test node"}, nil)
	require.Nil(t, err)
	require.False(t, node.Id.IsZero())
	require.NotNil(t, node.Tags)
	require.Equal(t, "tag 1", node.Tags[0].Name)
	require.Equal(t, "red", node.Tags[0].Color)
	require.Equal(t, interfaces.ModelColNameNode, node.Tags[0].Col)
}

func TestNodeService_GetModelList(t *testing.T) {
	SetupTest(t)

	n := 10
	for i := 0; i < n; i++ {
		node := &models2.Node{
			Name:     fmt.Sprintf("test node %d", i),
			IsMaster: true,
			Tags: []models2.Tag{
				{Name: fmt.Sprintf("tag %d", i)},
			},
		}
		err := delegate.NewModelDelegate(node).Add()
		require.Nil(t, err)
	}

	svc, err := service.NewService()
	require.Nil(t, err)

	nodes, err := svc.GetNodeList(nil, nil)
	require.Nil(t, err)

	for i := 0; i < n; i++ {
		node := nodes[i]
		require.False(t, node.Id.IsZero())
		require.False(t, node.Tags[0].Id.IsZero())
		require.Equal(t, fmt.Sprintf("tag %d", i), node.Tags[0].Name)
		require.Equal(t, interfaces.ModelColNameNode, node.Tags[0].Col)
	}
}
