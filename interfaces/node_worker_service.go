package interfaces

type NodeWorkerService interface {
	NodeService
	Subscribe()
	Unsubscribe()
	Recv()
	ReportStatus()
}
