package test

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestManagerService_Enqueue(t *testing.T) {
	var err error
	T.Setup(t)

	queue := T.managerSvc.GetQueue(T.TestTask.GetNodeId())

	go T.managerSvc.Start()
	defer T.managerSvc.Stop()
	time.Sleep(1 * time.Second)

	err = T.managerSvc.Enqueue(T.TestTask)
	require.Nil(t, err)
	time.Sleep(500 * time.Millisecond)
	n, err := T.redis.ZCountAll(queue)
	require.Nil(t, err)
	require.Equal(t, 1, n)
	data, err := T.redis.ZPopMaxOne(queue)
	require.Nil(t, err)
	time.Sleep(500 * time.Millisecond)
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

	queue := T.managerSvc.GetQueue(T.TestTaskWithNodeId.GetNodeId())

	go T.managerSvc.Start()
	defer T.managerSvc.Stop()
	time.Sleep(1 * time.Second)

	err = T.managerSvc.Enqueue(T.TestTaskWithNodeId)
	require.Nil(t, err)
	time.Sleep(500 * time.Millisecond)
	n, err := T.redis.ZCountAll(queue)
	require.Nil(t, err)
	require.Equal(t, 1, n)
	data, err := T.redis.ZPopMaxOne(queue)
	require.Nil(t, err)
	time.Sleep(500 * time.Millisecond)
	err = json.Unmarshal([]byte(data), &T.TestTaskMessage)
	require.Nil(t, err)
	require.Equal(t, T.TestTaskWithNodeId.GetId(), T.TestTaskMessage.Id)
	task, err := T.modelSvc.GetTaskById(T.TestTask.GetId())
	require.Nil(t, err)
	require.Equal(t, T.TestTask.GetId(), task.Id)
	require.Equal(t, constants.TaskStatusPending, task.Status)
}
