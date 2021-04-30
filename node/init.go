package node

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/store"
	"github.com/crawlab-team/crawlab-core/utils"
)

func initNode() (err error) {
	// default node service
	if store.NodeService, err = NewService(nil); err != nil {
		return err
	}

	// node service store
	store.NodeServiceStore = store.NewServiceStore()

	// set default node service into the store
	if err = store.NodeServiceStore.Set("", store.NodeService); err != nil {
		return err
	}

	return nil
}

func InitNode() (err error) {
	return utils.InitModule(interfaces.ModuleIdNode, initNode)
}
