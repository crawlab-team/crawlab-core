package models

type spiderService struct {
	*Service
}

func NewSpiderService() (svc *spiderService) {
	return &spiderService{NewService(ModelIdSpider)}
}

var SpiderService = NewSpiderService()
