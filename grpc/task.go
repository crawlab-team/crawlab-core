package grpc

import (
	"context"
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

	// response
	res = &pb.TaskServiceResponse{
		Code:   pb.ResponseCode_OK,
		Status: constants.GrpcSuccess,
		Task:   getPbTask(&t),
		Spider: getPbSpider(&sp),
		Node:   getPbNode(&n),
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

func getPbTask(t *model.Task) (pbT *pb.Task) {
	if t == nil {
		return nil
	}
	pbT = &pb.Task{
		XId:             t.Id,
		SpiderId:        t.SpiderId.Hex(),
		StartTs:         t.StartTs.String(),
		FinishTs:        t.FinishTs.String(),
		Status:          t.Status,
		NodeId:          t.NodeId.Hex(),
		Cmd:             t.Cmd,
		Param:           t.Param,
		Error:           t.Error,
		ResultCount:     int32(t.ResultCount),
		ErrorLogCount:   int32(t.ErrorLogCount),
		WaitDuration:    int32(t.WaitDuration),
		RuntimeDuration: int32(t.RuntimeDuration),
		TotalDuration:   int32(t.TotalDuration),
		Pid:             int32(t.Pid),
		RunType:         t.RunType,
		ScheduleId:      t.ScheduleId.Hex(),
		Type:            t.Type,
		UserId:          t.UserId.Hex(),
		CreateTs:        t.CreateTs.String(),
		UpdateTs:        t.UpdateTs.String(),
	}
	return pbT
}

func getPbSpider(sp *model.Spider) (pbSp *pb.Spider) {
	if sp == nil {
		return nil
	}
	pbSp = &pb.Spider{
		XId:         sp.Id.Hex(),
		Name:        sp.Name,
		DisplayName: sp.DisplayName,
		Type:        sp.Type,
		Col:         sp.Col,
		//Envs:        &[]pb.Env{},
		Remark:      sp.Remark,
		ProjectId:   sp.ProjectId.Hex(),
		IsPublic:    sp.IsPublic,
		Cmd:         sp.Cmd,
		IsScrapy:    sp.IsScrapy,
		Template:    sp.Template,
		IsDedup:     sp.IsDedup,
		DedupField:  sp.DedupField,
		DedupMethod: sp.DedupMethod,
		IsWebHook:   sp.IsWebHook,
		WebHookUrl:  sp.WebHookUrl,
		UserId:      sp.UserId.Hex(),
		CreateTs:    sp.CreateTs.String(),
		UpdateTs:    sp.UpdateTs.String(),
	}
	return pbSp
}

func getPbNode(n *model.Node) (pbN *pb.Node) {
	if n == nil {
		return nil
	}
	settings := &pb.NodeSettings{
		MaxRunners: int32(n.Settings.MaxRunners),
	}
	pbN = &pb.Node{
		XId:          n.Id.Hex(),
		Name:         n.Name,
		Status:       n.Status,
		Ip:           n.Ip,
		Port:         n.Port,
		Mac:          n.Mac,
		Hostname:     n.Hostname,
		Description:  n.Description,
		Key:          n.Key,
		Settings:     settings,
		IsMaster:     n.IsMaster,
		CreateTs:     n.CreateTs.String(),
		UpdateTs:     n.UpdateTs.String(),
		UpdateTsUnix: n.UpdateTsUnix,
	}
	return pbN
}
