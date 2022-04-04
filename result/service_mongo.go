package result

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/crawlab-team/crawlab-db/generic"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceMongo struct {
	// dependencies
	modelSvc    service.ModelService
	modelColSvc interfaces.ModelBaseService

	// internals
	id primitive.ObjectID     // id of models.DataCollection
	dc *models.DataCollection // models.DataCollection
}

func (svc *ServiceMongo) List(query generic.ListQuery, opts *generic.ListOptions) (results []interfaces.Result, err error) {
	_query := svc.getQuery(query)
	_opts := svc.getOpts(opts)
	return svc.getList(_query, _opts)
}

func (svc *ServiceMongo) Count(query generic.ListQuery) (n int, err error) {
	_query := svc.getQuery(query)
	return svc.modelColSvc.Count(_query)
}

func (svc *ServiceMongo) GetId() (id primitive.ObjectID) {
	return svc.id
}

func (svc *ServiceMongo) SetId(id primitive.ObjectID) {
	svc.id = id
}

func (svc *ServiceMongo) Insert(docs ...interface{}) (err error) {
	_, err = mongo.GetMongoCol(svc.dc.Name).InsertMany(docs)
	if err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (svc *ServiceMongo) getList(query bson.M, opts *mongo.FindOptions) (results []interfaces.Result, err error) {
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

func (svc *ServiceMongo) getQuery(query generic.ListQuery) (res bson.M) {
	return utils.GetMongoQuery(query)
}

func (svc *ServiceMongo) getOpts(opts *generic.ListOptions) (res *mongo.FindOptions) {
	return utils.GetMongoOpts(opts)
}
