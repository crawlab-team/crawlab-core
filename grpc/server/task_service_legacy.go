package server

//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/apex/log"
//	"github.com/crawlab-team/crawlab-core/constants"
//	"github.com/crawlab-team/crawlab-core/entity"
//	"github.com/crawlab-team/crawlab-core/models"
//	"github.com/crawlab-team/crawlab-db/mongo"
//	pb "github.com/crawlab-team/crawlab-grpc"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"sync"
//	"time"
//)
//
//var itemQueueMap = map[string]*[]entity.ResultItem{}
//var itemQueueMapLock = sync.Mutex{}
//
//type TaskResultItemServiceOptions struct {
//	FlushWaitSeconds int
//}
//
//func NewTaskResultItemService(options *TaskResultItemServiceOptions) (rs *TaskResultItemService) {
//	// normalize FlushWaitSeconds
//	if options.FlushWaitSeconds == 0 {
//		options.FlushWaitSeconds = 1
//	}
//
//	// task result item service
//	rs = &TaskResultItemService{
//		finished: false,
//		flushing: false,
//		opts:     options,
//	}
//
//	// init
//	rs.Init()
//
//	return rs
//}
//
//type TaskResultItemServiceInterface interface {
//	Init()
//	Stop()
//	Flush() (err error)
//}
//
//type TaskResultItemService struct {
//	TaskResultItemServiceInterface
//	finished bool
//	flushing bool
//	opts     *TaskResultItemServiceOptions
//}
//
//func (s *TaskResultItemService) Init() {
//	go func() {
//		for {
//			// if finished flag is set to true, end
//			if s.finished {
//				return
//			}
//
//			// flush
//			if err := s.Flush(); err != nil {
//				log.Error(fmt.Sprintf("flush error: %s", err.Error()))
//			}
//
//			// wait for a period
//			time.Sleep(time.Duration(s.opts.FlushWaitSeconds) * time.Second)
//		}
//	}()
//}
//
//func (s *TaskResultItemService) Stop() {
//	s.finished = true
//}
//
//func (s *TaskResultItemService) Flush() (err error) {
//	if s.flushing {
//		return nil
//	}
//
//	// lock
//	itemQueueMapLock.Lock()
//
//	for colName, itemQueue := range itemQueueMap {
//		// skip if no item in queue
//		if len(*itemQueue) == 0 {
//			continue
//		}
//
//		// save items
//		if err := s.saveItems(colName, *itemQueue); err != nil {
//			itemQueueMapLock.Unlock()
//			return err
//		}
//
//		// reset queue
//		*itemQueue = []entity.ResultItem{}
//	}
//
//	// reset queue map
//	itemQueueMap = map[string]*[]entity.ResultItem{}
//
//	// unlock
//	itemQueueMapLock.Unlock()
//
//	return err
//}
//
//func (s *TaskResultItemService) saveItems(colName string, items []entity.ResultItem) (err error) {
//	col := mongo.GetMongoCol(colName)
//
//	var _items []interface{}
//	for _, item := range items {
//		_items = append(_items, item)
//	}
//
//	// TODO: dedup
//	_, err = col.InsertMany(_items)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//type taskService struct {
//	pb.UnimplementedTaskServiceServer
//}
//
//var TaskService = taskService{}
//
//func (s taskService) GetTaskInfo(ctx context.Context, req *pb.TaskServiceRequest) (res *pb.TaskServiceResponse, err error) {
//	// get task
//	id, err := primitive.ObjectIDFromHex(req.TaskId)
//	if err != nil {
//		return nil, err
//	}
//	t, err := models.TaskService.GetModelById(id)
//	if err != nil {
//		return nil, err
//	}
//
//	// get spider
//	sp, err := models.SpiderService.GetModelById(t.SpiderId)
//	if err != nil {
//		return nil, err
//	}
//
//	// get node
//	var n models.Node
//	if !t.NodeId.IsZero() {
//		n, err = models.NodeServer.GetModelById(t.NodeId)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	// response
//	res = &pb.TaskServiceResponse{
//		Code:   pb.ResponseCode_OK,
//		Status: constants.GrpcSuccess,
//		Task:   getPbTask(&t),
//		Spider: getPbSpider(&sp),
//		Node:   getPbNode(&n),
//	}
//
//	return res, nil
//}
//
//func (s taskService) addItemToQueue(colName string, item entity.ResultItem) {
//	// lock
//	itemQueueMapLock.Lock()
//
//	// attempt to get item queue
//	queue, ok := itemQueueMap[colName]
//
//	if ok {
//		// exists
//		*queue = append(*queue, item)
//	} else {
//		// not exists
//		queue = &[]entity.ResultItem{}
//		*queue = append(*queue, item)
//		itemQueueMap[colName] = queue
//	}
//
//	// unlock
//	itemQueueMapLock.Unlock()
//}
//
//func (s taskService) getColName(req *pb.TaskServiceRequest) (colName string, err error) {
//	// task
//	id, err := primitive.ObjectIDFromHex(req.TaskId)
//	if err != nil {
//		return "", err
//	}
//	t, err := models.TaskService.GetModelById(id)
//	if err != nil {
//		return "", err
//	}
//
//	// spider
//	sp, err := models.SpiderService.GetModelById(t.SpiderId)
//	if err != nil {
//		return "", err
//	}
//
//	// normalize collection name
//	colName = sp.Col
//	if colName == "" {
//		colName = fmt.Sprintf("results_%s", sp.Id.Hex())
//	}
//
//	return colName, nil
//}
//
//func (s taskService) SaveItem(ctx context.Context, req *pb.TaskServiceRequest) (res *pb.TaskServiceResponse, err error) {
//	// declare result item
//	var item entity.ResultItem
//
//	// deserialize
//	if err := json.Unmarshal(req.Data, &item); err != nil {
//		return nil, err
//	}
//
//	// collection name
//	colName, err := s.getColName(req)
//	if err != nil {
//		return nil, err
//	}
//
//	// add item to queue
//	s.addItemToQueue(colName, item)
//
//	// response
//	res = handleTaskServiceSuccessResponse()
//
//	return res, nil
//}
//
//func (s taskService) SaveItems(ctx context.Context, req *pb.TaskServiceRequest) (res *pb.TaskServiceResponse, err error) {
//	// declare result items
//	var items []entity.ResultItem
//
//	// deserialize
//	if err := json.Unmarshal(req.Data, &items); err != nil {
//		return nil, err
//	}
//
//	// collection name
//	colName, err := s.getColName(req)
//	if err != nil {
//		return nil, err
//	}
//
//	// add items to queue
//	for _, item := range items {
//		s.addItemToQueue(colName, item)
//	}
//
//	// response
//	res = handleTaskServiceSuccessResponse()
//
//	return res, nil
//}
//
//func getPbTask(t *models.Task) (pbT *pb.Task) {
//	if t == nil {
//		return nil
//	}
//	ta, err := t.GetArtifact()
//	if err != nil {
//		return nil
//	}
//	pbT = &pb.Task{
//		XId:             t.Id.Hex(),
//		SpiderId:        t.SpiderId.Hex(),
//		StartTs:         t.StartTs.String(),
//		FinishTs:        t.FinishTs.String(),
//		Status:          t.Status,
//		NodeId:          t.NodeId.Hex(),
//		Cmd:             t.Cmd,
//		Param:           t.Param,
//		Error:           t.Error,
//		ResultCount:     int32(t.ResultCount),
//		ErrorLogCount:   int32(t.ErrorLogCount),
//		WaitDuration:    int32(t.WaitDuration),
//		RuntimeDuration: int32(t.RuntimeDuration),
//		TotalDuration:   int32(t.TotalDuration),
//		Pid:             int32(t.Pid),
//		RunType:         t.RunType,
//		ScheduleId:      t.ScheduleId.Hex(),
//		Type:            t.Type,
//		UserId:          ta.CreateUid.Hex(),
//		CreateTs:        ta.CreateTs.String(),
//		UpdateTs:        ta.UpdateTs.String(),
//	}
//	return pbT
//}
//
//func getPbSpider(s *models.Spider) (pbSp *pb.Spider) {
//	if s == nil {
//		return nil
//	}
//	sa, err := s.GetArtifact()
//	if err != nil {
//		return nil
//	}
//	pbSp = &pb.Spider{
//		XId:         s.Id.Hex(),
//		Name:        s.Name,
//		DisplayName: s.DisplayName,
//		Type:        s.Type,
//		Col:         s.Col,
//		//Envs:        &[]pb.Env{},
//		Remark:      s.Description,
//		ProjectId:   s.ProjectId.Hex(),
//		IsPublic:    s.IsPublic,
//		Cmd:         s.Cmd,
//		IsScrapy:    s.IsScrapy,
//		Template:    s.Template,
//		IsDedup:     s.IsDedup,
//		DedupField:  s.DedupField,
//		DedupMethod: s.DedupMethod,
//		IsWebHook:   s.IsWebHook,
//		WebHookUrl:  s.WebHookUrl,
//		UserId:      sa.CreateUid.String(),
//		CreateTs:    sa.CreateTs.String(),
//		UpdateTs:    sa.UpdateTs.String(),
//	}
//	return pbSp
//}
//
//func getPbNode(n *models.Node) (pbN *pb.Node) {
//	if n == nil {
//		return nil
//	}
//	na, err := n.GetArtifact()
//	if err != nil {
//		return nil
//	}
//	settings := &pb.NodeSettings{
//		MaxRunners: int32(n.Settings.MaxRunners),
//	}
//	pbN = &pb.Node{
//		XId:  n.Id.Hex(),
//		Name: n.Name,
//		//Status:       n.Status,
//		Ip:           n.Ip,
//		Port:         n.Port,
//		Mac:          n.Mac,
//		Hostname:     n.Hostname,
//		Description:  n.Description,
//		Key:          n.Key,
//		Settings:     settings,
//		IsMaster:     n.IsMaster,
//		CreateTs:     na.CreateTs.String(),
//		UpdateTs:     na.UpdateTs.String(),
//		UpdateTsUnix: na.UpdateTs.Unix(),
//	}
//	return pbN
//}
//
//func handleTaskServiceSuccessResponse() (res *pb.TaskServiceResponse) {
//	res = &pb.TaskServiceResponse{
//		Code:   pb.ResponseCode_OK,
//		Status: constants.GrpcSuccess,
//	}
//	return res
//}
