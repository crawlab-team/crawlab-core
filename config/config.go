package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"strings"
)

type Config struct {
	Name string
}

type InitConfigOptions struct {
	Name string
}

func (c *Config) WatchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
	})
}

func (c *Config) Init() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name) // if config file is set, load it accordingly
	} else {
		viper.AddConfigPath("./conf") // if no config file is set, load by default
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")   // default yaml
	viper.AutomaticEnv()          // load matched environment variables
	viper.SetEnvPrefix("CRAWLAB") // environment variable prefix as CRAWLAB
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil { // viper parsing config file
		return err
	}

	return nil
}

func InitConfig() error {
	// config instance
	c := Config{Name: ""}

	// init config file
	if err := c.Init(); err != nil {
		return err
	}

	// watch config change and load responsively
	c.WatchConfig()

	return nil
}
