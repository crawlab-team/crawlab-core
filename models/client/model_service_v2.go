package client

import (
	"encoding/json"
	"fmt"
	"github.com/crawlab-team/crawlab-core/grpc/client"
	"github.com/crawlab-team/crawlab-core/grpc/payload"
	"github.com/crawlab-team/crawlab-core/interfaces"
	nodeconfig "github.com/crawlab-team/crawlab-core/node/config"
	grpc "github.com/crawlab-team/crawlab-grpc"
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

func (svc ModelServiceV2[T]) GetById(id primitive.ObjectID) (model *T, err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()

	var v T
	t := reflect.TypeOf(v)
	typeName := t.Name()
	//typeName := fmt.Sprintf("%T", *new(T))

	p := payload.ModelServiceV2Payload[T]{
		Type: typeName,
		Id:   id,
	}
	d, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	req := &grpc.Request{
		NodeKey: svc.cfg.GetNodeKey(),
		Data:    d,
	}

	res, err := svc.c.ModelBaseServiceV2Client.GetById(ctx, req)
	if err != nil {
		return nil, err
	}

	var result T
	err = json.Unmarshal(res.Data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
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
