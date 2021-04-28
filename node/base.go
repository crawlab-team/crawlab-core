package node

import (
	"github.com/mitchellh/go-homedir"
	"path"
)

var homeDirPath, _ = homedir.Dir()

var configDirName = ".crawlab"

var defaultConfigDirPath = path.Join(homeDirPath, configDirName)

var configName = "config.json"

var defaultConfigPath = path.Join(homeDirPath, configDirName, configName)
