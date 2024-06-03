package server

import (
	"context"
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/grpc/payload"
	"github.com/crawlab-team/crawlab-core/models/service"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ModelBaseServiceV2Server struct {
	grpc.UnimplementedModelBaseServiceV2Server
}

func (svr ModelBaseServiceV2Server) GetById(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	var p payload.ModelServiceV2Payload[[]byte]
	err = json.Unmarshal(req.Data, &p)
	if err != nil {
		return nil, err
	}

	data, err := service.NewModelServiceV2[[]byte]().GetById(p.Id)
	if err != nil {
		return nil, err
	}

	return HandleSuccessWithData(data)
}

func (svr ModelBaseServiceV2Server) Get(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}

func (svr ModelBaseServiceV2Server) GetList(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetList not implemented")
}

func (svr ModelBaseServiceV2Server) DeleteById(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteById not implemented")
}

func (svr ModelBaseServiceV2Server) DeleteList(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteList not implemented")
}

func (svr ModelBaseServiceV2Server) UpdateById(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateById not implemented")
}

func (svr ModelBaseServiceV2Server) UpdateOne(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateOne not implemented")
}

func (svr ModelBaseServiceV2Server) UpdateMany(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMany not implemented")
}

func (svr ModelBaseServiceV2Server) ReplaceById(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReplaceById not implemented")
}

func (svr ModelBaseServiceV2Server) Replace(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method Replace not implemented")
}

func (svr ModelBaseServiceV2Server) InsertOne(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method InsertOne not implemented")
}

func (svr ModelBaseServiceV2Server) InsertMany(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method InsertMany not implemented")
}

func (svr ModelBaseServiceV2Server) Count(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method Count not implemented")
}

func NewModelBaseServiceV2Server() *ModelBaseServiceV2Server {
	return &ModelBaseServiceV2Server{}
}
