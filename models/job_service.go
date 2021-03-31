package models

type jobService struct {
	*Service
}

func NewJobService() (svc *jobService) {
	return &jobService{
		NewService(ModelIdJob),
	}
}

var JobService = NewJobService()
