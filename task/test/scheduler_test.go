package test

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSchedulerService_Fetch(t *testing.T) {
	var err error
	T.Setup(t)

	//ntest.T.StartMasterWorker()
	//defer ntest.T.StopMasterWorker()
	//time.Sleep(3 * time.Second)

	T.managerSvc.Start()
	go T.schedulerSvc.Fetch()
	time.Sleep(500 * time.Millisecond)

	// enqueue
	go func() {
		time.Sleep(500 * time.Millisecond)
		err = T.managerSvc.Enqueue(T.TestTask)
		require.Nil(t, err)
	}()

	// validate
	isErr := true
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	go func(t *testing.T) {
		ch := T.schedulerSvc.GetTaskChannel()
		tasks := <-ch
		ctx.Done()
		require.NotNil(t, tasks)
		require.Len(t, tasks, 1)
		require.Equal(t, tasks[0].GetId(), T.TestTask.GetId())
		isErr = false
	}(t)
	time.Sleep(2 * time.Second)
	require.Nil(t, ctx.Err())
	require.False(t, isErr)
}
