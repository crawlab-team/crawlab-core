package models

type projectService struct {
	*Service
}

func NewProjectService() (svc *projectService) {
	return &projectService{NewService(ModelIdProject)}
}

var ProjectService = NewProjectService()
