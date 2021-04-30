package node

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/store"
	"github.com/crawlab-team/crawlab-core/utils"
)

func initNode() (err error) {
	if store.NodeService, err = NewService(nil); err != nil {
		return err
	}
	store.NodeServiceStore = NewServiceStore()
	return nil
}

func InitNode() (err error) {
	return utils.InitModule(interfaces.ModuleIdNode, initNode)
}
