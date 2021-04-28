package grpc

import (
	"context"
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/models"
	grpc2 "github.com/crawlab-team/crawlab-grpc"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNodeServer_Register(t *testing.T) {
	setupTest(t)

	client, err := TestWorkerService.GetDefaultClient()
	require.Nil(t, err)

	nodeKey := "test-node-key"
	res, err := client.NodeClient.Register(context.Background(), &grpc2.Request{
		NodeKey: nodeKey,
	})
	require.Nil(t, err)
	require.Equal(t, grpc2.ResponseCode_OK, res.Code)

	var node models.Node
	err = json.Unmarshal(res.Data, &node)
	require.Nil(t, err)
	require.Equal(t, nodeKey, node.Key)
	require.Equal(t, constants.NodeStatusRegistered, node.Status)
	require.False(t, node.Id.IsZero())
}

func TestNodeServer_Ping(t *testing.T) {
	setupTest(t)

	client, err := TestWorkerService.GetDefaultClient()
	require.Nil(t, err)

	nodeKey := "test-node-key"
	res, err := client.NodeClient.Register(context.Background(), &grpc2.Request{
		NodeKey: nodeKey,
	})
	require.Nil(t, err)
	require.Equal(t, grpc2.ResponseCode_OK, res.Code)

	tic := time.Now()
	res, err = client.NodeClient.Ping(context.Background(), &grpc2.Request{
		NodeKey: nodeKey,
	})
	require.Nil(t, err)
	var node models.Node
	err = json.Unmarshal(res.Data, &node)
	require.Nil(t, err)
	require.Equal(t, constants.NodeStatusOnline, node.Status)
	require.Equal(t, nodeKey, node.Key)
	require.False(t, node.Id.IsZero())
	toc := node.ActiveTs
	require.LessOrEqual(t, tic.Unix(), toc.Unix())
	require.True(t, toc.Sub(tic) < 1*time.Second)
}
