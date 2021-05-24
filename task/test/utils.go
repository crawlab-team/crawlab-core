package test

func startAllServices() {
	go T.managerSvc.Start()
	go T.schedulerSvc.Start()
	go T.handlerSvc.Start()
}

func stopAllServices() {
	go T.managerSvc.Stop()
	go T.schedulerSvc.Stop()
	go T.handlerSvc.Stop()
}
