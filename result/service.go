package result

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

type Service struct {
	// dependencies
	modelSvc    service.ModelService
	modelColSvc interfaces.ModelBaseService

	// internals
	id primitive.ObjectID     // id of models.DataCollection
	dc *models.DataCollection // models.DataCollection
}

func (svc *Service) GetId() (id primitive.ObjectID) {
	return svc.id
}

func (svc *Service) SetId(id primitive.ObjectID) {
	svc.id = id
}

func (svc *Service) GetList(query bson.M, opts *mongo.FindOptions) (results []interfaces.Result, err error) {
	return svc.getList(query, opts)
}

func (svc *Service) Count(query bson.M) (total int, err error) {
	return svc.modelColSvc.Count(query)
}

func (svc *Service) Insert(docs ...interface{}) (err error) {
	_, err = mongo.GetMongoCol(svc.dc.Name).InsertMany(docs)
	if err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (svc *Service) getList(query bson.M, opts *mongo.FindOptions) (results []interfaces.Result, err error) {
	list, err := svc.modelColSvc.GetList(query, opts)
	if err != nil {
		return nil, err
	}
	for _, d := range list.Values() {
		r, ok := d.(interfaces.Result)
		if ok {
			results = append(results, r)
		}
	}
	return results, nil
}

func NewResultService(id primitive.ObjectID, opts ...Option) (svc2 interfaces.ResultService, err error) {
	// service
	svc := &Service{
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
