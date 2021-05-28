package test

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestHandlerService_Run(t *testing.T) {
	var err error
	T.Setup(t)

	err = T.schedulerSvc.Enqueue(T.TestTask)
	require.Nil(t, err)

	err = T.handlerSvc.Run(T.TestTask.GetId())
	require.Nil(t, err)
	time.Sleep(1 * time.Second)

	task, err := T.modelSvc.GetTaskById(T.TestTask.GetId())
	require.Nil(t, err)
	require.Equal(t, constants.TaskStatusFinished, task.Status)
}
