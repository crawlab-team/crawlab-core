package store

func InitStore() (err error) {
	if NodeService, err = initializeNodeService(); err != nil {
		return err
	}
	if GrpcService, err = initializeGrpcService(); err != nil {
		return err
	}

	return nil
}
