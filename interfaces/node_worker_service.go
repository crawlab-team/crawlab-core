package interfaces

import "time"

type NodeWorkerService interface {
	NodeService
	Recv()
	ReportStatus()
	SetHeartbeatInterval(duration time.Duration)
}
