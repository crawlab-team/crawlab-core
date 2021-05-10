package test

import (
	"context"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/node/test"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGrpcServer_Register(t *testing.T) {
	var err error

	T, _ = NewTest()
	T.Setup(t)

	// register
	register(t)

	// validate
	workerNodeKey := T.WorkerNodeInfo.Key
	workerNode, err := test.T.ModelSvc.GetNodeByKey(workerNodeKey, nil)
	require.Nil(t, err)
	require.Equal(t, workerNodeKey, workerNode.Key)
	require.Equal(t, constants.NodeStatusRegistered, workerNode.Status)
}

func TestGrpcServer_SendHeartbeat(t *testing.T) {
	var err error

	T, _ = NewTest()
	T.Setup(t)

	// register
	register(t)

	// send heartbeat
	sendHeartbeat(t)

	// validate
	workerNodeKey := T.WorkerNodeInfo.Key
	workerNode, err := test.T.ModelSvc.GetNodeByKey(workerNodeKey, nil)
	require.Nil(t, err)
	require.Equal(t, workerNodeKey, workerNode.Key)
	require.Equal(t, constants.NodeStatusOnline, workerNode.Status)
}

func TestGrpcServer_Stream(t *testing.T) {
	var err error

	T, _ = NewTest()
	T.Setup(t)

	// register
	register(t)

	// client-side stream
	clientSideStream(t)

	// CONNECT
	err = T.ClientSideStream.Send(&grpc.StreamMessage{
		Code:    grpc.StreamMessageCode_CONNECT,
		NodeKey: T.WorkerNodeInfo.Key,
	})
	require.Nil(t, err)

	time.Sleep(1 * time.Second)

	// server-side stream
	serverSideStream(t)

	// handle client stream message
	go handleClientStreamMessage(t)

	time.Sleep(1 * time.Second)

	// PING
	err = T.ServerSideStream.Send(&grpc.StreamMessage{
		Code:    grpc.StreamMessageCode_PING,
		NodeKey: T.MasterNodeInfo.Key,
	})
	if err != nil {
		log.Errorf("stream ping error: %v", err)
	}
	require.Nil(t, err)

	time.Sleep(1 * time.Second)
}

func register(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := T.Client.GetNodeClient().Register(ctx, T.Client.NewRequest(T.WorkerNodeInfo))
	require.Nil(t, err)
	require.Equal(t, grpc.ResponseCode_OK, res.Code)
}

func sendHeartbeat(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := T.Client.GetNodeClient().SendHeartbeat(ctx, T.Client.NewRequest(T.WorkerNodeInfo))
	require.Nil(t, err)
	require.Equal(t, grpc.ResponseCode_OK, res.Code)
}

func serverSideStream(t *testing.T) {
	var err error
	T.ServerSideStream, err = T.Server.GetSubscribe(T.WorkerNodeInfo.Key)
	require.Nil(t, err)
}

func clientSideStream(t *testing.T) {
	var err error
	T.ClientSideStream, err = T.Client.GetNodeClient().Stream(context.Background())
	require.Nil(t, err)
}

func handleClientStreamMessage(t *testing.T) {
	for {
		msg, err := T.ClientSideStream.Recv()
		require.Nil(t, err)
		switch msg.Code {
		case grpc.StreamMessageCode_CONNECT:
			log.Infof("stream connect data: %v", msg.Data)
			//require.Nil(t, msg.Data)
		case grpc.StreamMessageCode_DISCONNECT:
			log.Infof("stream disconnect data: %v", msg.Data)
			//require.Nil(t, msg.Data)
		case grpc.StreamMessageCode_PING:
			log.Infof("stream ping data: %v", msg.Data)
			//require.Nil(t, msg.Data)
		}

		time.Sleep(10 * time.Millisecond)
	}
}
