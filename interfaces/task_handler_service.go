package interfaces

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TaskHandlerService interface {
	TaskBaseService
	// Run task and execute locally
	Run(taskId primitive.ObjectID) (err error)
	// Cancel task locally
	Cancel(taskId primitive.ObjectID) (err error)
	// ReportStatus periodically report handler status to master
	ReportStatus()
	// Reset reset internals to default
	Reset()
	// GetMaxRunners get max runners
	GetMaxRunners() (maxRunners int)
	// SetMaxRunners set max runners
	SetMaxRunners(maxRunners int)
	// GetExitWatchDuration get max runners
	GetExitWatchDuration() (duration time.Duration)
	// SetExitWatchDuration set max runners
	SetExitWatchDuration(duration time.Duration)
	// GetReportInterval get report interval
	GetReportInterval() (interval time.Duration)
	// SetReportInterval set report interval
	SetReportInterval(interval time.Duration)
	// GetModelService get model service
	GetModelService() (modelSvc GrpcClientModelService)
	// GetModelSpiderService get model spider service
	GetModelSpiderService() (modelSpiderSvc GrpcClientModelSpiderService)
	// GetModelTaskService get model task service
	GetModelTaskService() (modelTaskSvc GrpcClientModelTaskService)
}
