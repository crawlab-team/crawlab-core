package services

import (
	"github.com/crawlab-team/crawlab-core/spider"
	"github.com/crawlab-team/crawlab-core/task"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/crawlab-team/crawlab-db/redis"
)

func InitAll() (err error) {
	if redis.RedisClient == nil {
		if err := redis.InitRedis(); err != nil {
			return err
		}
	}
	if mongo.Client == nil {
		if err := mongo.InitMongo(); err != nil {
			return err
		}
	}
	if err := InitNodeService(); err != nil {
		return err
	}
	if err := spider.InitSpiderService(); err != nil {
		return err
	}
	if err := task.InitTaskService(); err != nil {
		return err
	}
	if err := InitScheduleService(); err != nil {
		return err
	}
	return nil
}
