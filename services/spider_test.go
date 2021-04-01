package services

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/models"
	"github.com/crawlab-team/crawlab-db/redis"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

func setupTestSpider() (err error) {
	viper.Set("server.master", true)
	viper.Set("spider.fs", "spiders")
	viper.Set("spider.workspace", "/tmp/crawlab/workspace")
	viper.Set("spider.repo", "/tmp/crawlab/repo")
	if err := InitAll(); err != nil {
		return err
	}
	cleanupTestSpider()
	return nil
}

func cleanupTestSpider() {
	_ = redis.RedisClient.Del("tasks:public")
	_ = models.NodeService.DeleteList(nil)
	_ = models.SpiderService.DeleteList(nil)
	_ = models.TaskService.DeleteList(nil)
}

func TestSpiderService_Run(t *testing.T) {
	err := setupTestSpider()
	require.Nil(t, err)

	// spider
	s := models.Spider{
		Name: "test_spider",
		Cmd:  "python main.py",
	}
	err = s.Add()
	require.Nil(t, err)
	require.False(t, s.Id.IsZero())

	// script
	script := `print('it works')`
	fsSvc, err := SpiderService.GetFs(s.Id)
	require.Nil(t, err)
	err = fsSvc.Save("main.py", []byte(script), nil)
	require.Nil(t, err)

	// run
	err = SpiderService.Run(s.Id, &SpiderRunOptions{
		Mode: constants.RunTypeRandom,
	})
	require.Nil(t, err)

	// validate task status
	time.Sleep(5 * time.Second)
	task, err := models.TaskService.GetModel(bson.M{"spider_id": s.Id}, nil)
	require.Nil(t, err)
	require.False(t, task.Id.IsZero())
	require.Equal(t, constants.StatusFinished, task.Status)

	cleanupTestSpider()
}
