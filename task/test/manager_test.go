package test

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-db/redis"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestManagerService_Enqueue(t *testing.T) {
	var err error
	T.Setup(t)

	go T.managerSvc.Start()
	time.Sleep(1 * time.Second)

	err = T.managerSvc.Enqueue(T.TestTask)
	require.Nil(t, err)
	n, err := redis.RedisClient.LLen("tasks:public")
	require.Nil(t, err)
	require.Equal(t, 1, n)
	data, err := redis.RedisClient.LPop("tasks:public")
	require.Nil(t, err)
	err = json.Unmarshal([]byte(data), &T.TestTaskMessage)
	require.Nil(t, err)
	require.Equal(t, T.TestTask.GetId(), T.TestTaskMessage.Id)
	task, err := T.modelSvc.GetTaskById(T.TestTask.GetId())
	require.Nil(t, err)
	require.Equal(t, T.TestTask.GetId(), task.Id)
	require.Equal(t, constants.TaskStatusPending, task.Status)
}

func TestManagerService_Enqueue_WithNodeId(t *testing.T) {
	var err error
	T.Setup(t)

	go T.managerSvc.Start()
	time.Sleep(1 * time.Second)

	err = T.managerSvc.Enqueue(T.TestTaskWithNodeId)
	require.Nil(t, err)
	n, err := redis.RedisClient.LLen("tasks:node:" + T.TestNode.GetId().Hex())
	require.Nil(t, err)
	require.Equal(t, 1, n)
	data, err := redis.RedisClient.LPop("tasks:node:" + T.TestNode.GetId().Hex())
	require.Nil(t, err)
	err = json.Unmarshal([]byte(data), &T.TestTaskMessage)
	require.Nil(t, err)
	require.Equal(t, T.TestTaskWithNodeId.GetId(), T.TestTaskMessage.Id)
	task, err := T.modelSvc.GetTaskById(T.TestTask.GetId())
	require.Nil(t, err)
	require.Equal(t, T.TestTask.GetId(), task.Id)
	require.Equal(t, constants.TaskStatusPending, task.Status)
}
