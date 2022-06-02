package utils

import "github.com/spf13/viper"

func IsDocker() (ok bool) {
	isDockerBool := viper.GetBool("docker")
	isDockerString := viper.GetString("docker")
	return isDockerBool || isDockerString == "Y"
}
