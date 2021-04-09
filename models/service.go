package models

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"time"
)

type ServiceInterface interface {
	findId(id primitive.ObjectID) (fr *mongo.FindResult)
	find(query bson.M, opts *mongo.FindOptions) (fr *mongo.FindResult)
	deleteId(id primitive.ObjectID) (err error)
	delete(query bson.M) (err error)
	count(query bson.M) (total int, err error)
	updateId(id primitive.ObjectID, update interface{}) (err error)
	update(query bson.M, update interface{}, fields []string) (err error)
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
	Update(query bson.M, update interface{}, fields []string) (err error)
	Count(query bson.M) (total int, err error)
}

func NewCommonService(id ModelId) (svc *CommonService) {
	if mongo.Client == nil {
		_ = mongo.InitMongo()
	}
	colName := getModelColName(id)
	col := mongo.GetMongoCol(colName)
	return &CommonService{
		id:  id,
		col: col,
	}
}

func getModelColName(id ModelId) (colName string) {
	return NewColNameBinder(id).MustBindString()
}

type CommonService struct {
	id  ModelId
	col *mongo.Col
	ServiceInterface
	PublicServiceInterface
}

func (svc *CommonService) findId(id primitive.ObjectID) (fr *mongo.FindResult) {
	if svc.col == nil {
		return mongo.NewFindResultWithError(constants.ErrMissingCol)
	}
	return svc.col.FindId(id)
}

func (svc *CommonService) find(query bson.M, opts *mongo.FindOptions) (fr *mongo.FindResult) {
	if svc.col == nil {
		return mongo.NewFindResultWithError(constants.ErrMissingCol)
	}
	return svc.col.Find(query, opts)
}

func (svc *CommonService) deleteId(id primitive.ObjectID) (err error) {
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

func (svc *CommonService) delete(query bson.M) (err error) {
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

func (svc *CommonService) count(query bson.M) (total int, err error) {
	if svc.col == nil {
		return total, trace.TraceError(constants.ErrMissingCol)
	}
	return svc.col.Count(query)
}

func (svc *CommonService) update(query bson.M, update interface{}, fields []string) (err error) {
	v := reflect.ValueOf(update)
	switch reflect.TypeOf(update).Kind() {
	case reflect.Struct:
		// ids of query
		var ids []primitive.ObjectID
		docs := NewListBinder(svc.id, NewModelListMap(), svc.find(query, nil)).MustBind()
		vDocs := reflect.ValueOf(docs)
		for i := 0; i < vDocs.Len(); i++ {
			item := vDocs.Index(i)
			fId := item.FieldByName("Id")
			if !fId.CanInterface() {
				return errors.ErrorModelInvalidType
			}
			objId := fId.Interface()
			id, ok := objId.(primitive.ObjectID)
			if !ok {
				return errors.ErrorModelInvalidType
			}
			ids = append(ids, id)
		}

		// convert to bson.M
		var updateBsonM bson.M
		bytes, err := json.Marshal(&update)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(bytes, &updateBsonM); err != nil {
			return err
		}

		// fields map
		fieldsMap := map[string]bool{}
		for _, f := range fields {
			fieldsMap[f] = true
		}

		// remove unselected fields
		for k := range updateBsonM {
			if _, ok := fieldsMap[k]; !ok {
				delete(updateBsonM, k)
			}
		}

		// update model objects
		if err := svc.col.Update(query, bson.M{"$set": updateBsonM}); err != nil {
			return err
		}

		// update artifacts
		colA := mongo.GetMongoCol(ArtifactColName)
		if err := colA.Update(query, bson.M{
			"$set": bson.M{
				"_sys.update_ts": time.Now(),
				// TODO: update_uid
			},
		}); err != nil {
			return err
		}

		return nil
	case reflect.Ptr:
		return svc.update(query, v.Elem().Interface(), fields)
	default:
		return errors.ErrorModelInvalidType
	}
}

func (svc *CommonService) updateId(id primitive.ObjectID, update interface{}) (err error) {
	if svc.col == nil {
		return trace.TraceError(constants.ErrMissingCol)
	}
	return svc.col.UpdateId(id, update)
}

func (svc *CommonService) GetById(id primitive.ObjectID) (res interface{}, err error) {
	// declare
	m := NewModelMap()

	// find result
	fr := svc.findId(id)

	// bind
	return NewBasicBinder(svc.id, m, fr).Bind()
}

func (svc *CommonService) Get(query bson.M, opts *mongo.FindOptions) (res interface{}, err error) {
	// declare
	m := NewModelMap()

	// find result
	fr := svc.find(query, opts)

	// bind
	return NewBasicBinder(svc.id, m, fr).Bind()
}

func (svc *CommonService) GetList(query bson.M, opts *mongo.FindOptions) (res interface{}, err error) {
	// declare
	m := NewModelListMap()

	// find result
	fr := svc.find(query, opts)

	// bind
	return NewListBinder(svc.id, m, fr).Bind()
}

func (svc *CommonService) DeleteById(id primitive.ObjectID) (err error) {
	return svc.deleteId(id)
}

func (svc *CommonService) Delete(query bson.M) (err error) {
	return svc.delete(query)
}

func (svc *CommonService) DeleteList(query bson.M) (err error) {
	return svc.delete(query)
}

func (svc *CommonService) UpdateById(id primitive.ObjectID, update interface{}) (err error) {
	return svc.updateId(id, update)
}

func (svc *CommonService) Update(query bson.M, update interface{}, fields []string) (err error) {
	return svc.update(query, update, fields)
}

func (svc *CommonService) Count(query bson.M) (total int, err error) {
	return svc.count(query)
}
