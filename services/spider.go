package services

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpiderServiceInterface interface {
	// basic operations
	Run(id primitive.ObjectID, opts *SpiderRunOptions) (err error)
	Clone(id primitive.ObjectID, opts *SpiderCloneOptions) (err error)
	Delete(id primitive.ObjectID) (err error)
	Sync(id primitive.ObjectID) (err error)

	// get *spiderFsService
	GetFs(id primitive.ObjectID) (fsSvc *spiderFsService, err error)
}

type SpiderServiceOptions struct {
}

type SpiderRunOptions struct {
	Mode       string
	NodeIds    []primitive.ObjectID
	Param      string
	ScheduleId primitive.ObjectID
}

type SpiderCloneOptions struct {
	Name string
}

func NewSpiderService(opts *SpiderServiceOptions) (svc *spiderService, err error) {
	svc = &spiderService{}
	return svc, nil
}

func InitSpiderService() (err error) {
	SpiderService, err = NewSpiderService(nil)
	return err
}

type spiderService struct {
}

func (svc *spiderService) Run(id primitive.ObjectID, opts *SpiderRunOptions) (err error) {
	// sync to workspace
	if err := svc.Sync(id); err != nil {
		return err
	}

	// spider
	s, err := model.SpiderService.GetById(id)
	if err != nil {
		return err
	}

	// assign tasks
	if err := svc.assignTasks(&s, opts); err != nil {
		return err
	}

	return nil
}

func (svc *spiderService) Clone(id primitive.ObjectID, opts *SpiderCloneOptions) (err error) {
	// TODO: implement
	return nil
}

func (svc *spiderService) Delete(id primitive.ObjectID) (err error) {
	panic("implement me")
}

func (svc *spiderService) Sync(id primitive.ObjectID) (err error) {
	if fsSvc, err := svc.GetFs(id); err == nil {
		return fsSvc.SyncToWorkspace()
	}
	return err
}

func (svc *spiderService) GetFs(id primitive.ObjectID) (fsSvc *spiderFsService, err error) {
	fsSvc, err = NewSpiderFsService(&SpiderFsServiceOptions{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return fsSvc, nil
}

func (svc *spiderService) assignTasks(s *model.Spider, opts *SpiderRunOptions) (err error) {
	// main task
	mainTask := model.Task{
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
			t := model.Task{
				SpiderId: s.Id,
				ParentId: mainTask.Id,
				Mode:     opts.Mode,
				Cmd:      s.Cmd,
				Param:    opts.Param,
				NodeId:   nodeId,
				Status:   constants.StatusPending,
			}
			if err := TaskService.Assign(&t); err != nil {
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
		if err := TaskService.Assign(&mainTask); err != nil {
			return err
		}
	}

	return nil
}

func (svc *spiderService) getNodeIds(opts *SpiderRunOptions) (nodeIds []primitive.ObjectID, err error) {
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

func (svc *spiderService) isMultiTask(opts *SpiderRunOptions) (res bool) {
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

var SpiderService *spiderService
