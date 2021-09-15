package event

import (
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/go-trace"
	"github.com/thoas/go-funk"
	"regexp"
)

var S interfaces.EventService

type Service struct {
	chs  []chan interfaces.EventData
	keys []string
}

func (svc *Service) Register(key string, ch chan interfaces.EventData) {
	svc.chs = append(svc.chs, ch)
	svc.keys = append(svc.keys, key)
}

func (svc *Service) Unregister(key string) {
	idx := funk.IndexOfString(svc.keys, key)
	if idx != -1 {
		svc.chs = append(svc.chs[:idx], svc.chs[(idx+1):]...)
		svc.keys = append(svc.keys[:idx], svc.keys[(idx+1):]...)
	}
}

func (svc *Service) SendEvent(eventName string, data ...interface{}) {
	for i, key := range svc.keys {
		matched, err := regexp.MatchString(key, eventName)
		if err != nil {
			trace.PrintError(err)
			continue
		}
		if matched {
			ch := svc.chs[i]
			for _, d := range data {
				ch <- &entity.EventData{
					Event: eventName,
					Data:  d,
				}
			}
		}
	}
}

func NewEventService() (svc interfaces.EventService) {
	if S != nil {
		return S
	}

	svc = &Service{
		chs:  []chan interfaces.EventData{},
		keys: []string{},
	}

	S = svc

	return svc
}
