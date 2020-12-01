package local_node

import (
	"github.com/crawlab-team/crawlab-core/model"
	"github.com/spf13/viper"
)

func GetLocalNode() *LocalNode {
	return localNode
}
func CurrentNode() *model.Node {
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
