package utils

import "github.com/spf13/viper"

func EnvIsTrue(key string) bool {
	isTrueBool := viper.GetBool(key)
	isTrueString := viper.GetString(key)
	return isTrueBool || isTrueString == "Y"
}
