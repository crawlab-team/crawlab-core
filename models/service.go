package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
)

type Service struct {
	interfaces.ModelService
}

func NewService() (svc *Service) {
	return &Service{}
}
