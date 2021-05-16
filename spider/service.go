package spider

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/task"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewSpiderService(opts *interfaces.ServiceOptions) (svc *Service, err error) {
	svc = &Service{}
	return svc, nil
}

type Service struct {
	modelSvc service.ModelService
	fsSvc    interfaces.FsService
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
	if fsSvc, err := svc.GetFs(id); err == nil {
		return svc.fsSvc.SyncToWorkspace()
	}
	return err
}

func (svc *Service) GetFs(id primitive.ObjectID) (fsSvc *FsService, err error) {
	fsSvc, err = NewSpiderFsService(&SpiderFsServiceOptions{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return fsSvc, nil
}

func (svc *Service) assignTasks(s *models.Spider, opts *interfaces.RunOptions) (err error) {
	// main task
	mainTask := models.Task{
		SpiderId:   s.Id,
		Mode:       opts.Mode,
		Cmd:        s.Cmd,
		Param:      opts.Param,
		Status:     constants.StatusPending,
		ScheduleId: opts.ScheduleId,
	}
	if err := mainTask.Add(); err != nil {
		return err
	}

	if svc.isMultiTask(opts) {
		// multi tasks
		nodeIds, err := svc.getNodeIds(opts)
		if err != nil {
			return err
		}
		for _, nodeId := range nodeIds {
			t := models.Task{
				SpiderId: s.Id,
				ParentId: mainTask.Id,
				Mode:     opts.Mode,
				Cmd:      s.Cmd,
				Param:    opts.Param,
				NodeId:   nodeId,
				Status:   constants.StatusPending,
			}
			if err := task.TaskService.Assign(&t); err != nil {
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
		if err := task.TaskService.Assign(&mainTask); err != nil {
			return err
		}
	}

	return nil
}

func (svc *Service) getNodeIds(opts *interfaces.RunOptions) (nodeIds []primitive.ObjectID, err error) {
	if opts.Mode == constants.RunTypeAllNodes {
		nodeIds, err = NodeService.GetAllNodeIds()
	} else if opts.Mode == constants.RunTypeSelectedNodes {
		nodeIds = opts.NodeIds
	}
	if err != nil {
		return nil, err
	}
	return nodeIds, nil
}

func (svc *Service) isMultiTask(opts *interfaces.RunOptions) (res bool) {
	if opts.Mode == constants.RunTypeAllNodes {
		nodeIds, _ := NodeService.GetAllNodeIds()
		return len(nodeIds) > 1
	} else if opts.Mode == constants.RunTypeRandom {
		return false
	} else if opts.Mode == constants.RunTypeSelectedNodes {
		return len(opts.NodeIds) > 1
	} else {
		return false
	}
}
