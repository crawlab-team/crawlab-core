package models

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceInterface interface {
	findId(id primitive.ObjectID) (fr *mongo.FindResult)
	find(query bson.M, opts *mongo.FindOptions) (fr *mongo.FindResult)
	deleteId(id primitive.ObjectID) (err error)
	delete(query bson.M) (err error)
	count(query bson.M) (total int, err error)
	update(query bson.M, update interface{}) (err error)
	updateId(id primitive.ObjectID, update interface{}) (err error)
	//insert(docs ...interface{}) (err error) // TODO: implement
}

type PublicServiceInterface interface {
	GetById(id primitive.ObjectID) (res interface{}, err error)
	Get(query bson.M, opts *mongo.FindOptions) (res interface{}, err error)
	GetList(query bson.M, opts *mongo.FindOptions) (res interface{}, err error)
	DeleteById(id primitive.ObjectID) (err error)
	Delete(query bson.M) (err error)
	DeleteList(query bson.M) (err error)
	UpdateById(id primitive.ObjectID, update interface{}) (err error)
	Update(query bson.M, update interface{}) (err error)
	Count(query bson.M) (total int, err error)
}

func NewService(id ModelId) (svc *Service) {
	if mongo.Client == nil {
		_ = mongo.InitMongo()
	}
	colName := getModelColName(id)
	col := mongo.GetMongoCol(colName)
	return &Service{
		id:  id,
		col: col,
	}
}

func getModelColName(id ModelId) (colName string) {
	switch id {
	case ModelIdNode:
		return ModelColNameNode
	case ModelIdProject:
		return ModelColNameProject
	case ModelIdSpider:
		return ModelColNameSpider
	case ModelIdTask:
		return ModelColNameTask
	case ModelIdSchedule:
		return ModelColNameSchedule
	case ModelIdUser:
		return ModelColNameUser
	case ModelIdSetting:
		return ModelColNameSetting
	case ModelIdToken:
		return ModelColNameToken
	case ModelIdVariable:
		return ModelColNameVariable
	default:
		panic(errors.ErrorModelNotImplemented)
	}
}

type Service struct {
	id  ModelId
	col *mongo.Col
	ServiceInterface
	PublicServiceInterface
}

func (svc *Service) findId(id primitive.ObjectID) (fr *mongo.FindResult) {
	if svc.col == nil {
		return mongo.NewFindResultWithError(constants.ErrMissingCol)
	}
	return svc.col.FindId(id)
}

func (svc *Service) find(query bson.M, opts *mongo.FindOptions) (fr *mongo.FindResult) {
	if svc.col == nil {
		return mongo.NewFindResultWithError(constants.ErrMissingCol)
	}
	return svc.col.Find(query, opts)
}

func (svc *Service) deleteId(id primitive.ObjectID) (err error) {
	if svc.col == nil {
		return trace.TraceError(constants.ErrMissingCol)
	}
	var doc BaseModel
	if err := svc.findId(id).One(&doc); err != nil {
		return err
	}
	d := NewDelegate(svc.col.GetName(), &doc)
	return d.Delete()
}

func (svc *Service) delete(query bson.M) (err error) {
	if svc.col == nil {
		return trace.TraceError(constants.ErrMissingCol)
	}
	var docs []BaseModel
	if err := svc.find(query, nil).All(&docs); err != nil {
		return err
	}
	for _, doc := range docs {
		if err := svc.deleteId(doc.Id); err != nil {
			return err
		}
	}
	return nil
}

func (svc *Service) count(query bson.M) (total int, err error) {
	if svc.col == nil {
		return total, trace.TraceError(constants.ErrMissingCol)
	}
	return svc.col.Count(query)
}

func (svc *Service) update(query bson.M, update interface{}) (err error) {
	return svc.col.Update(query, bson.M{
		"$set": update,
	})
}

func (svc *Service) updateId(id primitive.ObjectID, update interface{}) (err error) {
	if svc.col == nil {
		return trace.TraceError(constants.ErrMissingCol)
	}
	return svc.col.UpdateId(id, update)
}

func (svc *Service) GetById(id primitive.ObjectID) (res interface{}, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *Service) Get(query bson.M, opts *mongo.FindOptions) (res interface{}, err error) {
	// declare
	var n Node
	var p Project
	var s Spider
	var t Task
	var sch Schedule
	var u User
	var st Setting
	var tk Token
	var v Variable

	// find result
	fr := svc.find(query, opts)

	// bind
	switch svc.id {
	case ModelIdNode:
		err = fr.One(&n)
		return n, err
	case ModelIdProject:
		err = fr.One(&p)
		return p, err
	case ModelIdSpider:
		err = fr.One(&s)
		return s, err
	case ModelIdTask:
		err = fr.One(&t)
		return t, err
	case ModelIdSchedule:
		err = fr.One(&sch)
		return sch, err
	case ModelIdUser:
		err = fr.One(&u)
		return u, err
	case ModelIdSetting:
		err = fr.One(&st)
		return st, err
	case ModelIdToken:
		err = fr.One(&tk)
		return tk, err
	case ModelIdVariable:
		err = fr.One(&v)
		return v, err
	default:
		return nil, errors.ErrorModelInvalidModelId
	}
}

func (svc *Service) GetList(query bson.M, opts *mongo.FindOptions) (res interface{}, err error) {
	// declare
	var n []Node
	var p []Project
	var s []Spider
	var t []Task
	var sch []Schedule
	var u []User
	var st []Setting
	var tk []Token
	var v []Variable

	// find result
	fr := svc.find(query, opts)

	// bind
	switch svc.id {
	case ModelIdNode:
		err = fr.All(&n)
		return n, err
	case ModelIdProject:
		err = fr.All(&p)
		return p, err
	case ModelIdSpider:
		err = fr.All(&s)
		return s, err
	case ModelIdTask:
		err = fr.All(&t)
		return t, err
	case ModelIdSchedule:
		err = fr.All(&sch)
		return sch, err
	case ModelIdUser:
		err = fr.All(&u)
		return u, err
	case ModelIdSetting:
		err = fr.All(&st)
		return st, err
	case ModelIdToken:
		err = fr.All(&tk)
		return tk, err
	case ModelIdVariable:
		err = fr.All(&v)
		return v, err
	default:
		return nil, errors.ErrorModelInvalidModelId
	}
}

func (svc *Service) DeleteById(id primitive.ObjectID) (err error) {
	return svc.deleteId(id)
}

func (svc *Service) Delete(query bson.M) (err error) {
	return svc.delete(query)
}

func (svc *Service) DeleteList(query bson.M) (err error) {
	return svc.delete(query)
}

func (svc *Service) UpdateById(id primitive.ObjectID, update interface{}) (err error) {
	return svc.updateId(id, update)
}

func (svc *Service) Update(query bson.M, update interface{}) (err error) {
	return svc.update(query, update)
}

func (svc *Service) Count(query bson.M) (total int, err error) {
	return svc.count(query)
}
