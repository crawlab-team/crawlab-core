package models

import (
	"github.com/crawlab-team/crawlab-core/color"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/node"
	"github.com/crawlab-team/crawlab-core/store"
	"github.com/crawlab-team/crawlab-core/utils"
)

func initModelService(id interfaces.ModelId) (err error) {
	if err = store.ModelServiceStore.Set(id, NewBaseService(id)); err != nil {
		return err
	}
	return nil
}

func initModels() (err error) {
	// color service
	if err := color.InitColor(); err != nil {
		return err
	}

	// node service
	if err := node.InitNode(); err != nil {
		return err
	}

	// default model service
	store.RootModelService = NewService()

	// model service store
	store.ModelServiceStore = store.NewModelServiceStore()

	// set root model service into the store
	if err = store.ModelServiceStore.Set("", store.RootModelService); err != nil {
		return err
	}
	for _, item := range ModelInfoList {
		if err := initModelService(item.Id); err != nil {
			return err
		}
	}

	return nil
}

func InitModels() (err error) {
	return utils.InitModule(interfaces.ModuleIdModels, initModels)
}
