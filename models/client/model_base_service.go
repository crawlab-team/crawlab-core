package client

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/grpc/client"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/emirpasic/gods/lists/arraylist"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
)

type BaseServiceDelegate struct {
	// settings variables
	cfgPath string

	// internals
	id interfaces.ModelId
	c  interfaces.GrpcClient
}

func (d *BaseServiceDelegate) GetModelId() (id interfaces.ModelId) {
	return d.id
}

func (d *BaseServiceDelegate) SetModelId(id interfaces.ModelId) {
	d.id = id
}

func (d *BaseServiceDelegate) GetConfigPath() (path string) {
	return d.cfgPath
}

func (d *BaseServiceDelegate) SetConfigPath(path string) {
	d.cfgPath = path
}

func (d *BaseServiceDelegate) GetById(id primitive.ObjectID) (doc interfaces.Model, err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Id: id})
	res, err := d.c.GetModelBaseServiceClient().GetById(ctx, req)
	if err != nil {
		return nil, err
	}
	return NewBasicBinder(d.id, res).Bind()
}

func (d *BaseServiceDelegate) Get(query bson.M, opts *mongo.FindOptions) (doc interfaces.Model, err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Query: query, FindOptions: opts})
	res, err := d.c.GetModelBaseServiceClient().Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return NewBasicBinder(d.id, res).Bind()
}

func (d *BaseServiceDelegate) GetList(query bson.M, opts *mongo.FindOptions) (list arraylist.List, err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Query: query, FindOptions: opts})
	res, err := d.c.GetModelBaseServiceClient().Get(ctx, req)
	if err != nil {
		return list, err
	}
	return NewListBinder(d.id, res).Bind()
}

func (d *BaseServiceDelegate) DeleteById(id primitive.ObjectID) (err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Id: id})
	_, err = d.c.GetModelBaseServiceClient().DeleteById(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (d *BaseServiceDelegate) Delete(query bson.M) (err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Query: query})
	_, err = d.c.GetModelBaseServiceClient().Delete(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (d *BaseServiceDelegate) DeleteList(query bson.M) (err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Query: query})
	_, err = d.c.GetModelBaseServiceClient().DeleteList(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (d *BaseServiceDelegate) ForceDeleteList(query bson.M) (err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Query: query})
	_, err = d.c.GetModelBaseServiceClient().ForceDeleteList(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (d *BaseServiceDelegate) UpdateById(id primitive.ObjectID, update bson.M) (err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Id: id, Update: update})
	_, err = d.c.GetModelBaseServiceClient().UpdateById(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (d *BaseServiceDelegate) Update(query bson.M, update bson.M, fields []string) (err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Query: query, Update: update, Fields: fields})
	_, err = d.c.GetModelBaseServiceClient().Update(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (d *BaseServiceDelegate) Insert(docs ...interfaces.Model) (err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Docs: docs})
	_, err = d.c.GetModelBaseServiceClient().Insert(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (d *BaseServiceDelegate) Count(query bson.M) (total int, err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Query: query})
	res, err := d.c.GetModelBaseServiceClient().Insert(ctx, req)
	if err != nil {
		return total, err
	}
	if err := json.Unmarshal(res.Data, &total); err != nil {
		return total, err
	}
	return total, nil
}

func (d *BaseServiceDelegate) newRequest(params interfaces.GrpcBaseServiceParams) (req *grpc.Request, err error) {
	return d.c.NewModelBaseServiceRequest(d.id, params)
}
func (d *BaseServiceDelegate) mustNewRequest(params *entity.GrpcBaseServiceParams) (req *grpc.Request) {
	req, err := d.newRequest(params)
	if err != nil {
		panic(err)
	}
	return req
}

func NewBaseServiceDelegate(opts ...ModelBaseServiceDelegateOption) (svc2 interfaces.GrpcClientModelBaseService, err error) {
	// mongo
	if mongo.Client == nil {
		_ = mongo.InitMongo()
	}

	// base service
	svc := &BaseServiceDelegate{}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(client.ProvideGetClient(svc.GetConfigPath())); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(client interfaces.GrpcClient) {
		svc.c = client
	}); err != nil {
		return nil, err
	}

	return svc, nil
}
