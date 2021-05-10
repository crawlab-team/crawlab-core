package interfaces

import "time"

type NodeWorkerService interface {
	NodeService
	Subscribe()
	Unsubscribe()
	Recv()
	ReportStatus()
	SetHeartbeatInterval(duration time.Duration)
}
