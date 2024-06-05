package client

import (
	"encoding/json"
	"fmt"
	"github.com/crawlab-team/crawlab-core/grpc/client"
	"github.com/crawlab-team/crawlab-core/grpc/payload"
	"github.com/crawlab-team/crawlab-core/interfaces"
	nodeconfig "github.com/crawlab-team/crawlab-core/node/config"
	"github.com/crawlab-team/crawlab-db/mongo"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"sync"
)

var (
	instanceMap = make(map[string]interface{})
	onceMap     = make(map[string]*sync.Once)
	mu          sync.Mutex
)

type ModelServiceV2[T any] struct {
	cfg interfaces.NodeConfigService
	c   *client.GrpcClientV2
}

func (svc *ModelServiceV2[T]) GetById(id primitive.ObjectID) (model *T, err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	req, err := svc.newRequest(payload.ModelServiceV2Payload[T]{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	res, err := svc.c.ModelBaseServiceV2Client.GetById(ctx, req)
	if err != nil {
		return nil, err
	}
	return svc.deserializeOne(res)
}

func (svc *ModelServiceV2[T]) Get(query bson.M, options *mongo.FindOptions) (model *T, err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	req, err := svc.newRequest(payload.ModelServiceV2Payload[T]{
		Query:       query,
		FindOptions: options,
	})
	if err != nil {
		return nil, err
	}
	res, err := svc.c.ModelBaseServiceV2Client.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return svc.deserializeOne(res)
}

func (svc *ModelServiceV2[T]) GetList(query bson.M, options *mongo.FindOptions) (models []T, err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	req, err := svc.newRequest(payload.ModelServiceV2Payload[T]{
		Query:       query,
		FindOptions: options,
	})
	if err != nil {
		return nil, err
	}
	res, err := svc.c.ModelBaseServiceV2Client.GetList(ctx, req)
	if err != nil {
		return nil, err
	}
	return svc.deserializeMany(res)
}

func (svc *ModelServiceV2[T]) DeleteById(id primitive.ObjectID) (err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	req, err := svc.newRequest(payload.ModelServiceV2Payload[T]{
		Id: id,
	})
	if err != nil {
		return err
	}
	_, err = svc.c.ModelBaseServiceV2Client.DeleteById(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (svc *ModelServiceV2[T]) DeleteList(query bson.M) (err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	req, err := svc.newRequest(payload.ModelServiceV2Payload[T]{
		Query: query,
	})
	if err != nil {
		return err
	}
	_, err = svc.c.ModelBaseServiceV2Client.DeleteList(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (svc *ModelServiceV2[T]) UpdateById(id primitive.ObjectID, update bson.M) (err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	req, err := svc.newRequest(payload.ModelServiceV2Payload[T]{
		Id:     id,
		Update: update,
	})
	if err != nil {
		return err
	}
	_, err = svc.c.ModelBaseServiceV2Client.UpdateById(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (svc *ModelServiceV2[T]) UpdateOne(query bson.M, update bson.M) (err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	req, err := svc.newRequest(payload.ModelServiceV2Payload[T]{
		Query:  query,
		Update: update,
	})
	if err != nil {
		return err
	}
	_, err = svc.c.ModelBaseServiceV2Client.UpdateOne(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (svc *ModelServiceV2[T]) UpdateMany(query bson.M, update bson.M) (err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	req, err := svc.newRequest(payload.ModelServiceV2Payload[T]{
		Query:  query,
		Update: update,
	})
	if err != nil {
		return err
	}
	_, err = svc.c.ModelBaseServiceV2Client.UpdateMany(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (svc *ModelServiceV2[T]) ReplaceById(id primitive.ObjectID, model T) (err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	req, err := svc.newRequest(payload.ModelServiceV2Payload[T]{
		Id:    id,
		Model: model,
	})
	if err != nil {
		return err
	}
	_, err = svc.c.ModelBaseServiceV2Client.ReplaceById(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (svc *ModelServiceV2[T]) Replace(query bson.M, model T) (err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	req, err := svc.newRequest(payload.ModelServiceV2Payload[T]{
		Query: query,
		Model: model,
	})
	if err != nil {
		return err
	}
	_, err = svc.c.ModelBaseServiceV2Client.Replace(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (svc *ModelServiceV2[T]) InsertOne(model T) (id primitive.ObjectID, err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	req, err := svc.newRequest(payload.ModelServiceV2Payload[T]{
		Model: model,
	})
	if err != nil {
		return primitive.NilObjectID, err
	}
	res, err := svc.c.ModelBaseServiceV2Client.InsertOne(ctx, req)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return deserialize[primitive.ObjectID](res)
}

func (svc *ModelServiceV2[T]) InsertMany(models []T) (ids []primitive.ObjectID, err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	req, err := svc.newRequest(payload.ModelServiceV2Payload[T]{
		Models: models,
	})
	if err != nil {
		return nil, err
	}
	res, err := svc.c.ModelBaseServiceV2Client.InsertOne(ctx, req)
	if err != nil {
		return nil, err
	}
	return deserialize[[]primitive.ObjectID](res)
}

func (svc *ModelServiceV2[T]) Count(query bson.M) (total int, err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	req, err := svc.newRequest(payload.ModelServiceV2Payload[T]{
		Query: query,
	})
	if err != nil {
		return 0, err
	}
	res, err := svc.c.ModelBaseServiceV2Client.InsertOne(ctx, req)
	if err != nil {
		return 0, err
	}
	return deserialize[int](res)
}

func (svc *ModelServiceV2[T]) newRequest(p payload.ModelServiceV2Payload[T]) (req *grpc.Request, err error) {
	var v T
	p.Type = reflect.TypeOf(v).Name()
	d, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	return &grpc.Request{
		NodeKey: svc.cfg.GetNodeKey(),
		Data:    d,
	}, nil
}

func (svc *ModelServiceV2[T]) deserializeOne(res *grpc.Response) (result *T, err error) {
	r, err := deserialize[T](res)
	if err != nil {
		return nil, err
	}
	return &r, err
}

func (svc *ModelServiceV2[T]) deserializeMany(res *grpc.Response) (results []T, err error) {
	return deserialize[[]T](res)
}

func deserialize[T any](res *grpc.Response) (result T, err error) {
	err = json.Unmarshal(res.Data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func NewModelServiceV2[T any]() *ModelServiceV2[T] {
	typeName := fmt.Sprintf("%T", *new(T))

	mu.Lock()
	defer mu.Unlock()

	if _, exists := onceMap[typeName]; !exists {
		onceMap[typeName] = &sync.Once{}
	}

	var instance *ModelServiceV2[T]

	c, err := client.GetGrpcClientV2()
	if err != nil {
		panic(err)
	}
	if !c.IsStarted() {
		err = c.Start()
		if err != nil {
			panic(err)
		}
	}

	onceMap[typeName].Do(func() {
		instance = &ModelServiceV2[T]{
			cfg: nodeconfig.GetNodeConfigService(),
			c:   c,
		}
		instanceMap[typeName] = instance
	})

	return instanceMap[typeName].(*ModelServiceV2[T])
}
