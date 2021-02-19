package services

import (
	"github.com/crawlab-team/crawlab-core/model"
	"github.com/crawlab-team/crawlab-db/redis"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

func setupTestSchedule() (err error) {
	viper.Set("server.master", true)
	viper.Set("spider.fs", "spiders")
	viper.Set("spider.workspace", "/tmp/crawlab/workspace")
	viper.Set("spider.repo", "/tmp/crawlab/repo")
	viper.Set("schedule.monitorIntervalSeconds", 60)
	if err := InitAll(); err != nil {
		return err
	}
	cleanupTestSchedule()
	return nil
}

func cleanupTestSchedule() {
	_ = redis.RedisClient.Del("tasks:public")
	_ = model.NodeService.DeleteList(nil)
	_ = model.SpiderService.DeleteList(nil)
	_ = model.TaskService.DeleteList(nil)
	_ = model.ScheduleService.DeleteList(nil)
}

func TestScheduleService_AddSchedule(t *testing.T) {
	err := setupTestSchedule()
	require.Nil(t, err)

	// spider
	spider := model.Spider{
		Name: "test_spider",
		Cmd:  "python main.py",
	}
	err = spider.Add()
	require.Nil(t, err)
	require.False(t, spider.Id.IsZero())

	// script
	script := `print('it works')`
	fsSvc, err := SpiderService.GetFs(spider.Id)
	require.Nil(t, err)
	err = fsSvc.Save("main.py", []byte(script), nil)
	require.Nil(t, err)

	// test method
	sch := model.Schedule{
		Name:     "test_schedule",
		SpiderId: spider.Id,
		Cron:     "* * * * *",
		Enabled:  true,
	}
	err = ScheduleService.AddSchedule(&sch)
	require.Nil(t, err)

	// wait until it triggers
	var task model.Task
	timeout := 60
	for i := 0; i < timeout; i++ {
		task, err = model.TaskService.Get(bson.M{"schedule_id": sch.Id}, nil)
		if err == mongo2.ErrNoDocuments {
			time.Sleep(1 * time.Second)
			continue
		}
		require.Nil(t, err)
		break
	}
	require.NotNil(t, task)
	require.False(t, task.Id.IsZero())
	require.Equal(t, sch.Id, task.ScheduleId)
	require.Equal(t, spider.Id, task.SpiderId)

	cleanupTestSchedule()
}
