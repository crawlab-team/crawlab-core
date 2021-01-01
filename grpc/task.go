package grpc

import (
	"context"
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/model"
	pb "github.com/crawlab-team/crawlab-grpc"
)

type taskService struct {
	pb.UnimplementedTaskServiceServer
}

var TaskService = taskService{}

func (s taskService) GetTaskInfo(ctx context.Context, req *pb.TaskServiceRequest) (res *pb.TaskServiceResponse, err error) {
	// get task
	t, err := model.GetTask(req.TaskId)
	if err != nil {
		return handleTaskServiceError(err)
	}

	// get spider
	sp, err := model.GetSpider(t.SpiderId)
	if err != nil {
		return handleTaskServiceError(err)
	}

	// get node
	var n model.Node
	if t.NodeId != constants.ObjectIdNull {
		n, err = model.GetNode(t.NodeId)
		if err != nil {
			return handleTaskServiceError(err)
		}
	}

	// data
	data := make(map[string]interface{})
	data["task"] = t
	data["spider"] = sp
	data["node"] = n
	dataJsonBytes, err := json.Marshal(&data)
	if err != nil {
		return handleTaskServiceError(err)
	}
	dataJson := string(dataJsonBytes)

	res = &pb.TaskServiceResponse{
		Code:   pb.ResponseCode_OK,
		Status: constants.GrpcSuccess,
		Data:   dataJson,
	}
	return res, nil
}

func (s taskService) SaveItem(ctx context.Context, req *pb.TaskServiceRequest) (res *pb.TaskServiceResponse, err error) {
	// TODO: implement
	return res, nil
}

func (s taskService) SaveItems(ctx context.Context, req *pb.TaskServiceRequest) (res *pb.TaskServiceResponse, err error) {
	// TODO: implement
	return res, nil
}

func handleTaskServiceError(e error) (res *pb.TaskServiceResponse, err error) {
	res = &pb.TaskServiceResponse{
		Code:   pb.ResponseCode_ERROR,
		Status: constants.GrpcError,
		Error:  e.Error(),
	}
	return res, e
}
