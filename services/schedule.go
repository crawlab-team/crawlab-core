package services

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/model"
	"github.com/crawlab-team/go-trace"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ScheduleServiceInterface interface {
	Init() (err error)
	Close() (err error)
	Add(s *model.Schedule) (err error)
	Update(s *model.Schedule) (err error)
	Delete(id primitive.ObjectID) (err error)
	ParseCronSpec(spec string) (s *cron.SpecSchedule, err error)
}

type ScheduleServiceOptions struct {
	MonitorIntervalSeconds int
}

func NewScheduleService(opts *ScheduleServiceOptions) (svc *scheduleService, err error) {
	// normalize options
	if opts == nil {
		opts = &ScheduleServiceOptions{}
	}
	if opts.MonitorIntervalSeconds == 0 {
		opts.MonitorIntervalSeconds = 600
	}

	// service
	svc = &scheduleService{
		c:      cron.New(),
		active: false,
		opts:   opts,
	}

	// initialize
	if err := svc.Init(); err != nil {
		return nil, err
	}

	return svc, nil
}

func InitScheduleService() (err error) {
	ScheduleService, err = NewScheduleService(&ScheduleServiceOptions{
		MonitorIntervalSeconds: viper.GetInt("schedule.monitorIntervalSeconds"),
	})
	if err != nil {
		return err
	}
	return ScheduleService.Init()
}

type scheduleService struct {
	c      *cron.Cron
	active bool
	opts   *ScheduleServiceOptions
}

func (svc *scheduleService) Init() (err error) {
	svc.c.Start()
	go svc.monitorAndUpdateCron()
	return nil
}

func (svc *scheduleService) Close() (err error) {
	svc.c.Stop()
	svc.active = false
	return nil
}

func (svc *scheduleService) Add(s *model.Schedule) (err error) {
	// validate schedule
	if !svc.isValidCronSpec(s.Cron) {
		return trace.TraceError(constants.ErrInvalidCronSpec)
	}

	// add to database
	if s.Id.IsZero() {
		if err := s.Add(); err != nil {
			return err
		}
	} else {
		_, err = model.ScheduleService.GetById(s.Id)
		if err == nil {
			return constants.ErrAlreadyExists
		}
		if err != mongo2.ErrNoDocuments {
			return err
		}
		if err := s.Add(); err != nil {
			return err
		}
	}

	// update schedule
	return svc.Update(s)
}

func (svc *scheduleService) Update(s *model.Schedule) (err error) {
	// validate
	if s.Id.IsZero() {
		return trace.TraceError(constants.ErrMissingId)
	}

	// delete old from cron
	if s.EntryId > 0 {
		svc.c.Remove(s.EntryId)
	}

	// add new to cron
	if s.Enabled {
		if s.EntryId, err = svc.addFunc(s); err != nil {
			return err
		}
	}

	// update in database
	return s.Save()
}

func (svc *scheduleService) Delete(id primitive.ObjectID) (err error) {
	// schedule
	s, err := model.ScheduleService.GetById(id)
	if err != nil {
		return err
	}

	// delete from cron
	svc.c.Remove(s.EntryId)

	// delete from database
	return s.Delete()
}

func (svc *scheduleService) ParseCronSpec(spec string) (s *cron.SpecSchedule, err error) {
	sch, err := cron.ParseStandard(spec)
	if err != nil {
		return nil, err
	}
	s, ok := sch.(*cron.SpecSchedule)
	if !ok {
		return nil, constants.ErrInvalidType
	}
	return s, nil
}

func (svc *scheduleService) addFunc(s *model.Schedule) (entryId cron.EntryID, err error) {
	return svc.c.AddFunc(s.Cron, func() {
		_ = SpiderService.Run(s.SpiderId, &SpiderRunOptions{
			Mode:       s.Mode,
			NodeIds:    s.NodeIds,
			Param:      s.Param,
			ScheduleId: s.Id,
		})
	})
}

func (svc *scheduleService) monitorAndUpdateCron() {
	for {
		// all schedules
		schedulesMap := map[cron.EntryID]*model.Schedule{}
		schedules, err := model.ScheduleService.GetList(nil, nil)
		if err != nil {
			if err != mongo2.ErrNoDocuments {
				trace.PrintError(err)
			}
			continue
		}
		for _, sch := range schedules {
			// validate entry id duplication
			if _, ok := schedulesMap[sch.EntryId]; ok {
				_ = svc.Update(&sch)
			}

			// assign to map
			schedulesMap[sch.EntryId] = &sch
		}

		// all entries
		entriesMap := map[cron.EntryID]*cron.Entry{}
		for _, entry := range svc.c.Entries() {
			if _, ok := schedulesMap[entry.ID]; !ok {
				// remove if not exists
				svc.c.Remove(entry.ID)
				continue
			}

			// assign to map
			entriesMap[entry.ID] = &entry
		}

		// iterate schedules map
		for _, sch := range schedulesMap {
			// skip disabled or those with invalid entry id schedule
			if !sch.Enabled || sch.EntryId == 0 {
				continue
			}

			// add to cron if not in entries
			if _, ok := entriesMap[sch.EntryId]; !ok {
				_ = svc.Add(sch)
			}
		}

		// break if stopped
		if !svc.active {
			break
		}

		// wait
		time.Sleep(time.Duration(svc.opts.MonitorIntervalSeconds) * time.Second)
	}
}

func (svc *scheduleService) isValidCronSpec(spec string) (res bool) {
	_, err := cron.ParseStandard(spec)
	return err == nil
}

var ScheduleService *scheduleService
