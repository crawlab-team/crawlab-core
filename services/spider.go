package services

import (
	"fmt"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/model"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpiderServiceInterface interface {
	Run(id primitive.ObjectID, opts *SpiderRunOptions) (err error)
	Clone(id primitive.ObjectID, opts *SpiderCloneOptions) (err error)
	Delete(id primitive.ObjectID) (err error)
	Sync(id primitive.ObjectID) (err error)
}

type SpiderServiceOptions struct {
}

type SpiderRunOptions struct {
	Mode    string
	NodeIds []primitive.ObjectID
	Param   string
}

type SpiderCloneOptions struct {
	Name string
}

func NewSpiderService(opts *SpiderServiceOptions) (svc *spiderService, err error) {
	svc = &spiderService{}
	return svc, nil
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
	panic("implement me")
}

func (svc *spiderService) Delete(id primitive.ObjectID) (err error) {
	panic("implement me")
}

func (svc *spiderService) Sync(id primitive.ObjectID) (err error) {
	fsSvc, err := NewFileSystemService(&FileSystemServiceOptions{
		IsMaster:      viper.GetBool("server.master"),
		FsPath:        svc.getFsPath(id),
		WorkspacePath: svc.getWorkspacePath(id),
		RepoPath:      svc.getRepoPath(id),
	})
	if err != nil {
		return err
	}
	return fsSvc.SyncToWorkspace()
}

func (svc *spiderService) getFsPath(id primitive.ObjectID) (res string) {
	return fmt.Sprintf("%s/%s", viper.GetString("spider.fs"), id.Hex())
}

func (svc *spiderService) getWorkspacePath(id primitive.ObjectID) (res string) {
	return fmt.Sprintf("%s/%s", viper.GetString("spider.workspace"), id.Hex())
}

func (svc *spiderService) getRepoPath(id primitive.ObjectID) (res string) {
	return fmt.Sprintf("%s/%s", viper.GetString("spider.repo"), id.Hex())
}

func (svc *spiderService) assignTasks(s *model.Spider, opts *SpiderRunOptions) (err error) {
	// main task
	mainTask := model.Task{
		SpiderId: s.Id,
		Mode:     opts.Mode,
		Cmd:      s.Cmd,
		Param:    opts.Param,
		Status:   constants.StatusPending,
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

var SpiderService, _ = NewSpiderService(nil)
