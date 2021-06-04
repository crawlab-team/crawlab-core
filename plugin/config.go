package plugin

import (
	"github.com/crawlab-team/crawlab-core/config"
	"path"
)

const DefaultPluginFsPathBase = "plugins"
const DefaultPluginDirName = "plugins"

var DefaultPluginDirPath = path.Join(config.DefaultConfigDirPath, DefaultPluginDirName)
