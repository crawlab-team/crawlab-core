package utils

import "github.com/spf13/viper"

func IsDocker() (ok bool) {
	return viper.GetBool("docker")
}
