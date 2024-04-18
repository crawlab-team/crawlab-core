package fs

import (
	"github.com/apex/log"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"path"
)

func init() {
	rootDir, err := homedir.Dir()
	if err != nil {
		log.Warnf("cannot find home directory: %v", err)
		return
	}
	DefaultWorkspacePath = path.Join(rootDir, "crawlab_workspace")
	DefaultRepoPath = path.Join(rootDir, "crawlab_repo")

	workspacePath := viper.GetString("fs.workspace")
	if workspacePath != "" {
		viper.Set("fs.workspace", DefaultWorkspacePath)
	}
}

const DefaultFsPath = "/fs"

var DefaultWorkspacePath string
var DefaultRepoPath string
