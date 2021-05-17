package client

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
)

type SpiderServiceDelegate struct {
	interfaces.GrpcClientModelBaseService
}

func (svc *SpiderServiceDelegate) GetSpiderById(id primitive.ObjectID) (s interfaces.Spider, err error) {
	res, err := svc.GetById(id)
	if err != nil {
		return nil, err
	}
	s, ok := res.(interfaces.Spider)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return s, nil
}

func (svc *SpiderServiceDelegate) GetSpider(query bson.M, opts *mongo.FindOptions) (s interfaces.Spider, err error) {
	res, err := svc.Get(query, opts)
	if err != nil {
		return nil, err
	}
	s, ok := res.(interfaces.Spider)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return s, nil
}

func (svc *SpiderServiceDelegate) GetSpiderList(query bson.M, opts *mongo.FindOptions) (res []interfaces.Spider, err error) {
	list, err := svc.GetList(query, opts)
	if err != nil {
		return nil, err
	}
	for _, item := range list.Values() {
		s, ok := item.(interfaces.Spider)
		if !ok {
			return nil, errors.ErrorModelInvalidType
		}
		res = append(res, s)
	}
	return res, nil
}

func NewSpiderServiceDelegate() (svc2 interfaces.GrpcClientModelSpiderService, err error) {
	svc := &SpiderServiceDelegate{}

	// dependency injection
	c := dig.New()
	if err := c.Provide(ProvideBaseServiceDelegate(interfaces.ModelIdSpider)); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(baseSvc interfaces.GrpcClientModelBaseService) {
		svc.GrpcClientModelBaseService = baseSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	return svc, nil
}
