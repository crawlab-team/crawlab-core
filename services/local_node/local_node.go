package local_node

import (
	models2 "github.com/crawlab-team/crawlab-core/models/models"
	"github.com/spf13/viper"
)

func GetLocalNode() *LocalNode {
	return localNode
}
func CurrentNode() *models2.Node {
	return GetLocalNode().Current()
}

func InitLocalNode() (node *LocalNode, err error) {
	registerType := viper.GetString("server.register.type")
	ip := viper.GetString("server.register.ip")
	customNodeName := viper.GetString("server.register.customNodeName")

	localNode, err = NewLocalNode(ip, customNodeName, registerType)
	if err != nil {
		return nil, err
	}
	return localNode, err
}
