package result

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

func NewResultService(id primitive.ObjectID, opts ...Option) (svc2 interfaces.ResultService, err error) {
	// service
	svc := &ServiceMongo{
		id: id,
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// dependency injection
	svc.modelSvc, err = service.GetService()
	if err != nil {
		return nil, err
	}

	// data collection
	svc.dc, err = svc.modelSvc.GetDataCollectionById(id)
	if err != nil {
		return nil, err
	}

	// data collection model service
	svc.modelColSvc = service.GetBaseServiceByColName(interfaces.ModelIdResult, svc.dc.Name)

	return svc, nil
}

var store = sync.Map{}

func GetResultService(id primitive.ObjectID, opts ...Option) (svc interfaces.ResultService, err error) {
	res, ok := store.Load(id)
	if ok {
		svc, ok = res.(interfaces.ResultService)
		if ok {
			return svc, nil
		}
	}
	svc, err = NewResultService(id, opts...)
	if err != nil {
		return nil, err
	}
	store.Store(id, svc)
	return svc, nil
}

func ProvideGetResultService(id primitive.ObjectID, opts ...Option) func() (svc interfaces.ResultService, err error) {
	return func() (svc interfaces.ResultService, err error) {
		return GetResultService(id, opts...)
	}
}
