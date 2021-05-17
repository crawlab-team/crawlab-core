package test

//func TestSpiderService_Run(t *testing.T) {
//	var err error
//
//	// run
//	err = T.adminSvc.Run(T.s.Id, &interfaces.RunOptions{
//		Mode: constants.RunTypeRandom,
//	})
//	require.Nil(t, err)
//
//	// validate task status
//	time.Sleep(5 * time.Second)
//	task, err := T.modelSvc.GetTask(bson.M{"spider_id": T.s.Id}, nil)
//	require.Nil(t, err)
//	require.False(t, task.Id.IsZero())
//	require.Equal(t, constants.StatusFinished, task.Status)
//}
