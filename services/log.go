package services

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/lib/cron"
	"github.com/crawlab-team/crawlab-core/model"
	"github.com/crawlab-team/crawlab-core/utils"
	database "github.com/crawlab-team/crawlab-db"
	"github.com/globalsign/mgo"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// 任务日志频道映射
var TaskLogChanMap = utils.NewChanMap()

// 定时删除日志
func DeleteLogPeriodically() {
	logDir := viper.GetString("log.path")
	if !utils.Exists(logDir) {
		log.Error("Can Not Set Delete Logs Periodically,No Log Dir")
		return
	}
	rd, err := ioutil.ReadDir(logDir)
	if err != nil {
		log.Error("Read Log Dir Failed")
		return
	}

	for _, fi := range rd {
		if fi.IsDir() {
			log.Info(filepath.Join(logDir, fi.Name()))
			_ = os.RemoveAll(filepath.Join(logDir, fi.Name()))
			log.Info("Delete Log File Success")
		}
	}

}

// 初始化定时删除日志
func InitDeleteLogPeriodically() error {
	c := cron.New(cron.WithSeconds())
	if _, err := c.AddFunc(viper.GetString("log.deleteFrequency"), DeleteLogPeriodically); err != nil {
		return err
	}

	c.Start()
	return nil
}

func InitLogIndexes() error {
	s, c := database.GetCol("logs")
	defer s.Close()
	se, ce := database.GetCol("error_logs")
	defer se.Close()

	_ = c.EnsureIndex(mgo.Index{
		Key: []string{"task_id", "seq"},
	})
	_ = c.EnsureIndex(mgo.Index{
		Key: []string{"task_id", "msg"},
	})
	_ = c.EnsureIndex(mgo.Index{
		Key:         []string{"expire_ts"},
		Sparse:      true,
		ExpireAfter: 1 * time.Second,
	})
	_ = ce.EnsureIndex(mgo.Index{
		Key: []string{"task_id"},
	})
	_ = ce.EnsureIndex(mgo.Index{
		Key: []string{"log_id"},
	})
	_ = ce.EnsureIndex(mgo.Index{
		Key:         []string{"expire_ts"},
		Sparse:      true,
		ExpireAfter: 1 * time.Second,
	})

	return nil
}

func InitLogService() error {
	logLevel := viper.GetString("log.level")
	if logLevel != "" {
		log.SetLevelFromString(logLevel)
	}
	log.Info("initialized log config successfully")
	if viper.GetString("log.isDeletePeriodically") == "Y" {
		if err := InitDeleteLogPeriodically(); err != nil {
			log.Error("init DeletePeriodically failed")
			return err
		}
		log.Info("initialized periodically cleaning log successfully")
	} else {
		log.Info("periodically cleaning log is switched off")
	}

	if model.IsMaster() {
		if err := InitLogIndexes(); err != nil {
			log.Errorf(err.Error())
			return err
		}
	}

	return nil
}
