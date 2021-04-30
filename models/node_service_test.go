package models

import (
	"fmt"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestNodeService_GetModelById(t *testing.T) {
	setupTest(t)

	node := Node{
		Name:     "test node",
		IsMaster: true,
		Tags: []Tag{
			{Name: "tag 1", Color: "red"},
		},
	}
	err := node.Add()
	require.Nil(t, err)

	node, err = NodeService.GetModelById(node.Id)
	require.Nil(t, err)
	require.False(t, node.Id.IsZero())
	require.False(t, node.Tags[0].Id.IsZero())
	require.Equal(t, "tag 1", node.Tags[0].Name)
	require.Equal(t, "red", node.Tags[0].Color)
	require.Equal(t, interfaces.ModelColNameNode, node.Tags[0].Col)
}

func TestNodeService_GetModel(t *testing.T) {
	setupTest(t)

	node := Node{
		Name:     "test node",
		IsMaster: true,
		Tags: []Tag{
			{Name: "tag 1", Color: "red"},
		},
	}
	err := node.Add()
	require.Nil(t, err)

	node, err = NodeService.GetModel(bson.M{"name": "test node"}, nil)
	require.Nil(t, err)
	require.False(t, node.Id.IsZero())
	require.Equal(t, "tag 1", node.Tags[0].Name)
	require.Equal(t, "red", node.Tags[0].Color)
	require.Equal(t, interfaces.ModelColNameNode, node.Tags[0].Col)
}

func TestNodeService_GetModelList(t *testing.T) {
	setupTest(t)

	n := 10
	for i := 0; i < n; i++ {
		node := Node{
			Name:     fmt.Sprintf("test node %d", i),
			IsMaster: true,
			Tags: []Tag{
				{Name: fmt.Sprintf("tag %d", i)},
			},
		}
		err := node.Add()
		require.Nil(t, err)
	}

	nodes, err := NodeService.GetModelList(nil, nil)
	require.Nil(t, err)

	for i := 0; i < n; i++ {
		node := nodes[i]
		require.False(t, node.Id.IsZero())
		require.False(t, node.Tags[0].Id.IsZero())
		require.Equal(t, fmt.Sprintf("tag %d", i), node.Tags[0].Name)
		require.Equal(t, interfaces.ModelColNameNode, node.Tags[0].Col)
	}
}
