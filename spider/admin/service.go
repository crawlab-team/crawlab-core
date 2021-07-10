package admin

import (
	config2 "github.com/crawlab-team/crawlab-core/config"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/node/config"
	"github.com/crawlab-team/crawlab-core/task/scheduler"
	"github.com/crawlab-team/go-trace"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
)

type Service struct {
	// dependencies
	nodeCfgSvc   interfaces.NodeConfigService
	modelSvc     service.ModelService
	schedulerSvc interfaces.TaskSchedulerService

	// settings
	cfgPath string
}

func (svc *Service) GetConfigPath() (path string) {
	return svc.cfgPath
}

func (svc *Service) SetConfigPath(path string) {
	svc.cfgPath = path
}

func (svc *Service) Schedule(id primitive.ObjectID, opts *interfaces.SpiderRunOptions) (err error) {
	// spider
	s, err := svc.modelSvc.GetSpiderById(id)
	if err != nil {
		return err
	}

	// assign tasks
	if err := svc.scheduleTasks(s, opts); err != nil {
		return err
	}

	return nil
}

func (svc *Service) Clone(id primitive.ObjectID, opts *interfaces.SpiderCloneOptions) (err error) {
	// TODO: implement
	return nil
}

func (svc *Service) Delete(id primitive.ObjectID) (err error) {
	panic("implement me")
}

func (svc *Service) scheduleTasks(s *models.Spider, opts *interfaces.SpiderRunOptions) (err error) {
	// main task
	mainTask := &models.Task{
		SpiderId:   s.Id,
		Mode:       opts.Mode,
		NodeIds:    opts.NodeIds,
		NodeTags:   opts.NodeTags,
		Cmd:        opts.Cmd,
		Param:      opts.Param,
		ScheduleId: opts.ScheduleId,
		Priority:   opts.Priority,
	}

	if svc.isMultiTask(opts) {
		// multi tasks
		if err := delegate.NewModelDelegate(mainTask).Add(); err != nil {
			return err
		}
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
				Priority: opts.Priority,
			}
			if err := svc.schedulerSvc.Enqueue(t); err != nil {
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
		if err := svc.schedulerSvc.Enqueue(mainTask); err != nil {
			return err
		}
	}

	return nil
}

func (svc *Service) getNodeIds(opts *interfaces.SpiderRunOptions) (nodeIds []primitive.ObjectID, err error) {
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

func (svc *Service) isMultiTask(opts *interfaces.SpiderRunOptions) (res bool) {
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

func NewSpiderAdminService(opts ...Option) (svc2 interfaces.SpiderAdminService, err error) {
	svc := &Service{
		cfgPath: config2.DefaultConfigPath,
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(config.ProvideConfigService(svc.cfgPath)); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(service.NewService); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(scheduler.ProvideGetTaskSchedulerService(svc.cfgPath)); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(nodeCfgSvc interfaces.NodeConfigService, modelSvc service.ModelService, schedulerSvc interfaces.TaskSchedulerService) {
		svc.nodeCfgSvc = nodeCfgSvc
		svc.modelSvc = modelSvc
		svc.schedulerSvc = schedulerSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	// validate node type
	if !svc.nodeCfgSvc.IsMaster() {
		return nil, trace.TraceError(errors.ErrorSpiderForbidden)
	}

	return svc, nil
}

func ProvideSpiderAdminService(path string, opts ...Option) func() (svc interfaces.SpiderAdminService, err error) {
	if path != "" || path == config2.DefaultConfigPath {
		if viper.GetString("config.path") != "" {
			path = viper.GetString("config.path")
		} else {
			path = config2.DefaultConfigPath
		}
	}
	opts = append(opts, WithConfigPath(path))
	return func() (svc interfaces.SpiderAdminService, err error) {
		return NewSpiderAdminService(opts...)
	}
}
