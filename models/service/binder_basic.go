package service

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
)

func NewBasicBinder(id interfaces.ModelId, fr *mongo.FindResult) (b interfaces.ModelBinder) {
	return &BasicBinder{
		id: id,
		fr: fr,
		m:  models.NewModelMap(),
	}
}

type BasicBinder struct {
	id interfaces.ModelId
	fr *mongo.FindResult
	m  *models.ModelMap
}

func (b *BasicBinder) Bind() (res interfaces.Model, err error) {
	m := b.m

	switch b.id {
	case interfaces.ModelIdArtifact:
		return b.Process(&m.Artifact)
	case interfaces.ModelIdTag:
		return b.Process(&m.Tag)
	case interfaces.ModelIdNode:
		return b.ProcessWithFieldIds(&m.Node, interfaces.ModelIdTag)
	case interfaces.ModelIdProject:
		return b.ProcessWithFieldIds(&m.Project, interfaces.ModelIdTag)
	case interfaces.ModelIdSpider:
		return b.ProcessWithFieldIds(&m.Spider, interfaces.ModelIdTag)
	case interfaces.ModelIdTask:
		return b.Process(&m.Task)
	case interfaces.ModelIdJob:
		return b.Process(&m.Job)
	case interfaces.ModelIdSchedule:
		return b.Process(&m.Schedule)
	case interfaces.ModelIdUser:
		return b.Process(&m.User)
	case interfaces.ModelIdSetting:
		return b.Process(&m.Setting)
	case interfaces.ModelIdToken:
		return b.Process(&m.Token)
	case interfaces.ModelIdVariable:
		return b.Process(&m.Variable)
	case interfaces.ModelIdTaskQueue:
		return b.Process(&m.TaskQueueItem)
	case interfaces.ModelIdTaskStat:
		return b.Process(&m.TaskStat)
	case interfaces.ModelIdPlugin:
		return b.Process(&m.Plugin)
	case interfaces.ModelIdSpiderStat:
		return b.Process(&m.SpiderStat)
	case interfaces.ModelIdDataSource:
		return b.Process(&m.DataSource)
	case interfaces.ModelIdDataCollection:
		return b.Process(&m.DataCollection)
	case interfaces.ModelIdResult:
		return b.Process(&m.Result)
	case interfaces.ModelIdPassword:
		return b.Process(&m.Password)
	default:
		return nil, errors.ErrorModelInvalidModelId
	}
}

func (b *BasicBinder) MustBind() (res interfaces.Model) {
	res, err := b.Bind()
	if err != nil {
		panic(err)
	}
	return res
}

func (b *BasicBinder) Process(d interfaces.Model) (res interfaces.Model, err error) {
	if err := b.fr.One(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (b *BasicBinder) ProcessWithFieldIds(d interfaces.Model, fieldIds ...interfaces.ModelId) (res interfaces.Model, err error) {
	if d, err = b.Process(d); err != nil {
		return nil, err
	}
	return b.AssignFields(d, fieldIds...)
}

func (b *BasicBinder) AssignFields(d interfaces.Model, fieldIds ...interfaces.ModelId) (res interfaces.Model, err error) {
	return b.assignFields(d, fieldIds...)
}

func (b *BasicBinder) assignFields(d interfaces.Model, fieldIds ...interfaces.ModelId) (res interfaces.Model, err error) {
	// model service
	modelSvc, err := NewService()
	if err != nil {
		return nil, err
	}

	// convert to model
	doc, ok := d.(interfaces.Model)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}

	// skip if no field ids
	if len(fieldIds) == 0 {
		return doc, nil
	}

	// iterate  field ids
	for _, fid := range fieldIds {
		switch fid {
		case interfaces.ModelIdTag:
			// convert interface
			d, ok := doc.(interfaces.ModelWithTags)
			if !ok {
				return nil, errors.ErrorModelInvalidType
			}

			// attempt to get artifact
			a, err := delegate.NewModelDelegate(doc).GetArtifact()
			if err != nil {
				return nil, err
			}

			// skip if no artifact found
			if a == nil {
				return d, nil
			}

			// skip if artifact has no tags
			if a.GetTagIds() == nil {
				return d, nil
			}

			// get tags
			modelTagSvc := modelSvc.GetBaseService(interfaces.ModelIdTag)
			list, err := modelTagSvc.GetList(bson.M{
				"_id": bson.M{
					"$in": a.GetTagIds(),
				},
			}, nil)
			if err != nil {
				if err == mongo2.ErrNoDocuments {
					return d, nil
				}
				return nil, err
			}
			var tags []interfaces.Tag
			_ = list.All(func(index int, value interface{}) bool {
				tag, ok := value.(interfaces.Tag)
				if !ok {
					_ = trace.TraceError(errors.ErrorModelInvalidType)
					return false
				}
				tags = append(tags, tag)
				return true
			})

			// assign tags
			d.SetTags(tags)

			return d, nil
		}
	}
	return doc, nil
}
