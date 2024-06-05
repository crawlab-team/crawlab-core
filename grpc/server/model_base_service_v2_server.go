package server

import (
	"context"
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/grpc/payload"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
)

var (
	typeNameColNameMap = make(map[string]string)
	typeInstances      = []interface{}{
		models.TestModel{},
		models.NodeV2{},
		models.UserV2{},
	}
)

func init() {
	for _, v := range typeInstances {
		t := reflect.TypeOf(v)
		typeName := t.Name()
		colName := service.GetCollectionNameByInstance(v)
		typeNameColNameMap[typeName] = colName
	}
}

type ModelBaseServiceV2Server struct {
	grpc.UnimplementedModelBaseServiceV2Server
}

func (svr ModelBaseServiceV2Server) GetById(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	var p payload.ModelServiceV2Payload[any]
	err = json.Unmarshal(req.Data, &p)
	if err != nil {
		return nil, err
	}
	data, err := svr.GetModelService(p.Type).GetById(p.Id)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccessWithData(data)
}

func (svr ModelBaseServiceV2Server) Get(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	var p payload.ModelServiceV2Payload[any]
	err = json.Unmarshal(req.Data, &p)
	if err != nil {
		return nil, err
	}
	data, err := svr.GetModelService(p.Type).Get(p.Query, p.FindOptions)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccessWithData(data)
}

func (svr ModelBaseServiceV2Server) GetList(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	var p payload.ModelServiceV2Payload[any]
	err = json.Unmarshal(req.Data, &p)
	if err != nil {
		return nil, err
	}
	data, err := svr.GetModelService(p.Type).GetList(p.Query, p.FindOptions)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccessWithData(data)
}

func (svr ModelBaseServiceV2Server) DeleteById(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	var p payload.ModelServiceV2Payload[any]
	err = json.Unmarshal(req.Data, &p)
	if err != nil {
		return nil, err
	}
	err = svr.GetModelService(p.Type).DeleteById(p.Id)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccess()
}

func (svr ModelBaseServiceV2Server) DeleteList(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	var p payload.ModelServiceV2Payload[any]
	err = json.Unmarshal(req.Data, &p)
	if err != nil {
		return nil, err
	}
	err = svr.GetModelService(p.Type).DeleteList(p.Query)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccess()
}

func (svr ModelBaseServiceV2Server) UpdateById(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	var p payload.ModelServiceV2Payload[any]
	err = json.Unmarshal(req.Data, &p)
	if err != nil {
		return nil, err
	}
	err = svr.GetModelService(p.Type).UpdateById(p.Id, p.Update)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccess()
}

func (svr ModelBaseServiceV2Server) UpdateOne(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	var p payload.ModelServiceV2Payload[any]
	err = json.Unmarshal(req.Data, &p)
	if err != nil {
		return nil, err
	}
	err = svr.GetModelService(p.Type).UpdateOne(p.Query, p.Update)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccess()
}

func (svr ModelBaseServiceV2Server) UpdateMany(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	var p payload.ModelServiceV2Payload[any]
	err = json.Unmarshal(req.Data, &p)
	if err != nil {
		return nil, err
	}
	err = svr.GetModelService(p.Type).UpdateMany(p.Query, p.Update)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccess()
}

func (svr ModelBaseServiceV2Server) ReplaceById(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	var p payload.ModelServiceV2Payload[bson.M]
	err = json.Unmarshal(req.Data, &p)
	if err != nil {
		return nil, err
	}
	err = svr.GetModelService(p.Type).ReplaceById(p.Id, p.Model)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccess()
}

func (svr ModelBaseServiceV2Server) Replace(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	var p payload.ModelServiceV2Payload[bson.M]
	err = json.Unmarshal(req.Data, &p)
	if err != nil {
		return nil, err
	}
	err = svr.GetModelService(p.Type).Replace(p.Query, p.Model)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccess()
}

func (svr ModelBaseServiceV2Server) InsertOne(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	var p payload.ModelServiceV2Payload[bson.M]
	err = json.Unmarshal(req.Data, &p)
	if err != nil {
		return nil, err
	}
	id, err := svr.GetModelService(p.Type).InsertOne(p.Model)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccessWithData(id)
}

func (svr ModelBaseServiceV2Server) InsertMany(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	var p payload.ModelServiceV2Payload[bson.M]
	err = json.Unmarshal(req.Data, &p)
	if err != nil {
		return nil, err
	}
	ids, err := svr.GetModelService(p.Type).InsertMany(p.Models)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccessWithData(ids)
}

func (svr ModelBaseServiceV2Server) Count(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	var p payload.ModelServiceV2Payload[bson.M]
	err = json.Unmarshal(req.Data, &p)
	if err != nil {
		return nil, err
	}
	count, err := svr.GetModelService(p.Type).Count(p.Query)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccessWithData(count)
}

func (svr ModelBaseServiceV2Server) GetModelService(typeName string) *service.ModelServiceV2[bson.M] {
	return service.NewModelServiceV2WithColName(typeNameColNameMap[typeName])
}

func NewModelBaseServiceV2Server() *ModelBaseServiceV2Server {
	return &ModelBaseServiceV2Server{}
}
