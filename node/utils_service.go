package node

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/store"
)

func GetService(key string) (svc interfaces.NodeService, err error) {
	svc, err = store.NodeServiceStore.GetNodeService(key)
	if err != nil {
		return nil, err
	}
	return svc, nil
}

func MustGetService(key string) (svc interfaces.NodeService) {
	svc, err := GetService(key)
	if err != nil {
		panic(err)
	}
	return svc
}
