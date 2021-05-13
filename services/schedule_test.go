package services

import (
	"github.com/crawlab-team/crawlab-core/models"
	models2 "github.com/crawlab-team/crawlab-core/models/models"
	spider2 "github.com/crawlab-team/crawlab-core/spider"
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
}

func TestScheduleService_AddSchedule(t *testing.T) {
	err := setupTestSchedule()
	require.Nil(t, err)

	// spider
	spider := models2.Spider{
		Name: "test_spider",
		Cmd:  "python main.py",
	}
	err = spider.Add()
	require.Nil(t, err)
	require.False(t, spider.Id.IsZero())

	// script
	script := `print('it works')`
	fsSvc, err := spider2.SpiderService.GetFs(spider.Id)
	require.Nil(t, err)
	err = fsSvc.Save("main.py", []byte(script), nil)
	require.Nil(t, err)

	// test method
	sch := models2.Schedule{
		Name:     "test_schedule",
		SpiderId: spider.Id,
		Cron:     "* * * * *",
		Enabled:  true,
	}
	err = ScheduleService.Add(&sch)
	require.Nil(t, err)

	// wait until it triggers
	var task *models2.Task
	timeout := 60
	for i := 0; i < timeout; i++ {
		task, err = models.MustGetRootService().GetTask(bson.M{"schedule_id": sch.Id}, nil)
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

func TestScheduleService_UpdateSchedule(t *testing.T) {
	err := setupTestSchedule()
	require.Nil(t, err)

	// spider
	spider := models2.Spider{
		Name: "test_spider",
		Cmd:  "python main.py",
	}
	err = spider.Add()
	require.Nil(t, err)
	require.False(t, spider.Id.IsZero())

	// script
	script := `print('it works')`
	fsSvc, err := spider2.SpiderService.GetFs(spider.Id)
	require.Nil(t, err)
	err = fsSvc.Save("main.py", []byte(script), nil)
	require.Nil(t, err)

	// add schedule
	sch := models2.Schedule{
		Name:     "test_schedule",
		SpiderId: spider.Id,
		Cron:     "* * * * *",
		Enabled:  false,
	}
	err = ScheduleService.Add(&sch)
	require.Nil(t, err)

	// test method
	sch.Enabled = true
	err = ScheduleService.Update(&sch)
	require.Nil(t, err)

	// wait until it triggers
	var task *models2.Task
	timeout := 60
	for i := 0; i < timeout; i++ {
		task, err = models.MustGetRootService().GetTask(bson.M{"schedule_id": sch.Id}, nil)
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

func TestScheduleService_DeleteSchedule(t *testing.T) {
	err := setupTestSchedule()
	require.Nil(t, err)

	// spider
	spider := models2.Spider{
		Name: "test_spider",
		Cmd:  "python main.py",
	}
	err = spider.Add()
	require.Nil(t, err)
	require.False(t, spider.Id.IsZero())

	// script
	script := `print('it works')`
	fsSvc, err := spider2.SpiderService.GetFs(spider.Id)
	require.Nil(t, err)
	err = fsSvc.Save("main.py", []byte(script), nil)
	require.Nil(t, err)

	// add schedule
	sch := models2.Schedule{
		Name:     "test_schedule",
		SpiderId: spider.Id,
		Cron:     "* * * * *",
		Enabled:  false,
	}
	err = ScheduleService.Add(&sch)
	require.Nil(t, err)

	// test method
	sch.Enabled = true
	err = ScheduleService.Delete(sch.Id)
	require.Nil(t, err)

	// validate
	entry := ScheduleService.c.Entry(sch.EntryId)
	require.False(t, entry.Valid())
	_, err = models.MustGetRootService().GetScheduleById(sch.Id)
	require.NotNil(t, err)
	require.Equal(t, mongo2.ErrNoDocuments.Error(), err.Error())

	cleanupTestSchedule()
}

func TestScheduleService_ParseCronSpec(t *testing.T) {
	err := setupTestSchedule()
	require.Nil(t, err)

	// test method
	_, err = ScheduleService.ParseCronSpec("* * * * *")
	require.Nil(t, err)

	cleanupTestSchedule()
}
