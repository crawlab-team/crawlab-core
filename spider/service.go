package spider

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	// dependencies
	modelSvc     service.ModelService
	fsSvc        interfaces.SpiderFsService
	schedulerSvc interfaces.TaskSchedulerService
}

func (svc *Service) Run(id primitive.ObjectID, opts *interfaces.RunOptions) (err error) {
	// spider
	s, err := svc.modelSvc.GetSpiderById(id)
	if err != nil {
		return err
	}

	// assign tasks
	if err := svc.assignTasks(s, opts); err != nil {
		return err
	}

	return nil
}

func (svc *Service) Clone(id primitive.ObjectID, opts *interfaces.CloneOptions) (err error) {
	// TODO: implement
	return nil
}

func (svc *Service) Delete(id primitive.ObjectID) (err error) {
	panic("implement me")
}

func (svc *Service) Sync(id primitive.ObjectID) (err error) {
	if fsSvc, err := GetSpiderFsService(id); err == nil {
		return fsSvc.GetFsService().SyncToWorkspace()
	}
	return nil
}

func (svc *Service) assignTasks(s *models.Spider, opts *interfaces.RunOptions) (err error) {
	// main task
	mainTask := &models.Task{
		SpiderId:   s.Id,
		Mode:       opts.Mode,
		Cmd:        s.Cmd,
		Param:      opts.Param,
		Status:     constants.StatusPending,
		ScheduleId: opts.ScheduleId,
	}
	if err := delegate.NewModelDelegate(mainTask).Add(); err != nil {
		return err
	}

	if svc.isMultiTask(opts) {
		// multi tasks
		nodeIds, err := svc.getNodeIds(opts)
		if err != nil {
			return err
		}
		for _, nodeId := range nodeIds {
			t := &models.Task{
				SpiderId: s.Id,
				ParentId: mainTask.Id,
				Mode:     opts.Mode,
				Cmd:      s.Cmd,
				Param:    opts.Param,
				NodeId:   nodeId,
				Status:   constants.StatusPending,
			}
			if err := svc.schedulerSvc.Assign(t); err != nil {
				return err
			}
		}
	} else {
		// single task
		nodeIds, err := svc.getNodeIds(opts)
		if err != nil {
			return err
		}
		if len(nodeIds) > 0 {
			mainTask.NodeId = nodeIds[0]
		}
		if err := svc.schedulerSvc.Assign(mainTask); err != nil {
			return err
		}
	}

	return nil
}

func (svc *Service) getNodeIds(opts *interfaces.RunOptions) (nodeIds []primitive.ObjectID, err error) {
	if opts.Mode == constants.RunTypeAllNodes {
		query := bson.M{
			"active":  true,
			"enabled": true,
			"status":  constants.NodeStatusOnline,
		}
		nodes, err := svc.modelSvc.GetNodeList(query, nil)
		if err != nil {
			return nil, err
		}
		for _, node := range nodes {
			nodeIds = append(nodeIds, node.GetId())
		}
	} else if opts.Mode == constants.RunTypeSelectedNodes {
		nodeIds = opts.NodeIds
	}
	return nodeIds, nil
}

func (svc *Service) isMultiTask(opts *interfaces.RunOptions) (res bool) {
	if opts.Mode == constants.RunTypeAllNodes {
		query := bson.M{
			"active":  true,
			"enabled": true,
			"status":  constants.NodeStatusOnline,
		}
		nodes, err := svc.modelSvc.GetNodeList(query, nil)
		if err != nil {
			trace.PrintError(err)
			return false
		}
		return len(nodes) > 1
	} else if opts.Mode == constants.RunTypeRandom {
		return false
	} else if opts.Mode == constants.RunTypeSelectedNodes {
		return len(opts.NodeIds) > 1
	} else {
		return false
	}
}

func NewSpiderService(opts ...interfaces.ServiceOption) (svc2 interfaces.SpiderService, err error) {
	svc := &Service{}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// dependency injection

	return svc, nil
}
