package test

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSchedulerService_Enqueue(t *testing.T) {
	var err error
	T.Setup(t)

	task := T.NewTask()

	_, err = T.schedulerSvc.Enqueue(task)
	require.Nil(t, err)

	task, err = T.modelSvc.GetTask(nil, nil)
	require.Nil(t, err)
	require.False(t, task.GetId().IsZero())

	tq, err := T.modelSvc.GetTaskQueueItemById(task.GetId())
	require.Nil(t, err)
	require.Equal(t, task.GetId(), tq.Id)
}
