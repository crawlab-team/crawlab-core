package grpc

var initialized = false

func InitGrpcServices() (err error) {
	// skip if already initialized
	if initialized {
		return nil
	}

	GrpcService, err = NewService(nil)
	if err != nil {
		return err
	}

	// mark as initialized
	initialized = true
	return nil
}
