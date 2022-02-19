package config

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func init() {
	// config instance
	c := Config{Name: ""}

	// init config file
	if err := c.Init(); err != nil {
		log.Warn("unable to init config")
		return
	}

	// watch config change and load responsively
	c.WatchConfig()
}

type Config struct {
	Name string
}

type InitConfigOptions struct {
	Name string
}

func (c *Config) WatchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})
}

func (c *Config) Init() error {
	// config
	if c.Name != "" {
		viper.SetConfigFile(c.Name) // if config file is set, load it accordingly
	} else {
		viper.AddConfigPath("./conf") // if no config file is set, load by default
		viper.SetConfigName("config")
	}

	// config type as yaml
	viper.SetConfigType("yaml") // default yaml

	// auto env
	viper.AutomaticEnv() // load matched environment variables

	// env prefix
	viper.SetEnvPrefix("CRAWLAB") // environment variable prefix as CRAWLAB

	// replacer
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// create conf directory if not exists
	confDirName := "./conf"
	if !utils.Exists(confDirName) {
		if err := os.MkdirAll("./conf", os.ModeDir|os.ModePerm); err != nil {
			return err
		}
	}

	// add default conf/config.yaml if not exists
	confFilePath := path.Join(confDirName, "config.yml")
	if !utils.Exists(confFilePath) {
		if err := ioutil.WriteFile(confFilePath, []byte(DefaultConfigYaml), os.FileMode(0766)); err != nil {
			return err
		}
	}

	// read config
	if err := viper.ReadInConfig(); err != nil { // viper parsing config file
		return err
	}

	return nil
}
