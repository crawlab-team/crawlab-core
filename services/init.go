package services

import (
	"github.com/crawlab-team/crawlab-db/mongo"
)

func InitAll() (err error) {
	if mongo.Client == nil {
		if err := mongo.InitMongo(); err != nil {
			return err
		}
	}
	return nil
}
