package models

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/crawlab-team/go-trace"
	"github.com/emirpasic/gods/lists/arraylist"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"time"
)

type BaseService struct {
	id  interfaces.ModelId
	col *mongo.Col
}

func (svc *BaseService) GetById(id primitive.ObjectID) (res interface{}, err error) {
	// declare
	m := NewModelMap()

	// find result
	fr := svc.findId(id)

	// bind
	return NewBasicBinder(svc.id, m, fr).Bind()
}

func (svc *BaseService) Get(query bson.M, opts *mongo.FindOptions) (res interface{}, err error) {
	// declare
	m := NewModelMap()

	// find result
	fr := svc.find(query, opts)

	// bind
	return NewBasicBinder(svc.id, m, fr).Bind()
}

func (svc *BaseService) GetList(query bson.M, opts *mongo.FindOptions) (res arraylist.List, err error) {
	// declare
	m := NewModelListMap()

	// find result
	fr := svc.find(query, opts)

	// bind
	return NewListBinder(svc.id, m, fr).BindList()
}

func (svc *BaseService) GetListSerializeTarget(query bson.M, opts *mongo.FindOptions, target interface{}) (err error) {
	list, err := svc.GetList(query, opts)
	if err != nil {
		return err
	}
	return serializeList(list, target)
}

func (svc *BaseService) DeleteById(id primitive.ObjectID) (err error) {
	return svc.deleteId(id)
}

func (svc *BaseService) Delete(query bson.M) (err error) {
	return svc.delete(query)
}

func (svc *BaseService) DeleteList(query bson.M) (err error) {
	return svc.delete(query)
}

func (svc *BaseService) UpdateById(id primitive.ObjectID, update interface{}) (err error) {
	return svc.updateId(id, update)
}

func (svc *BaseService) Update(query bson.M, update interface{}, fields []string) (err error) {
	return svc.update(query, update, fields)
}

func (svc *BaseService) Count(query bson.M) (total int, err error) {
	return svc.count(query)
}

func (svc *BaseService) findId(id primitive.ObjectID) (fr *mongo.FindResult) {
	if svc.col == nil {
		return mongo.NewFindResultWithError(constants.ErrMissingCol)
	}
	return svc.col.FindId(id)
}

func (svc *BaseService) find(query bson.M, opts *mongo.FindOptions) (fr *mongo.FindResult) {
	if svc.col == nil {
		return mongo.NewFindResultWithError(constants.ErrMissingCol)
	}
	return svc.col.Find(query, opts)
}

func (svc *BaseService) deleteId(id primitive.ObjectID) (err error) {
	if svc.col == nil {
		return trace.TraceError(constants.ErrMissingCol)
	}
	var doc BaseModel
	if err := svc.findId(id).One(&doc); err != nil {
		return err
	}
	d := NewDelegate(svc.id, &doc)
	return d.Delete()
}

func (svc *BaseService) delete(query bson.M) (err error) {
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

func (svc *BaseService) count(query bson.M) (total int, err error) {
	if svc.col == nil {
		return total, trace.TraceError(constants.ErrMissingCol)
	}
	return svc.col.Count(query)
}

func (svc *BaseService) update(query bson.M, update interface{}, fields []string) (err error) {
	vUpdate := reflect.ValueOf(update)
	switch reflect.TypeOf(update).Kind() {
	case reflect.Struct:
		// ids of query
		var ids []primitive.ObjectID
		list := NewListBinder(svc.id, NewModelListMap(), svc.find(query, nil)).MustBindListAsPtr()
		for _, value := range list.Values() {
			item, ok := value.(interfaces.BaseModelInterface)
			if !ok {
				return errors.ErrorModelInvalidType
			}
			ids = append(ids, item.GetId())
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
		colA := mongo.GetMongoCol(interfaces.ModelColNameArtifact)
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
		return svc.update(query, vUpdate.Elem().Interface(), fields)
	default:
		return errors.ErrorModelInvalidType
	}
}

func (svc *BaseService) updateId(id primitive.ObjectID, update interface{}) (err error) {
	if svc.col == nil {
		return trace.TraceError(constants.ErrMissingCol)
	}
	return svc.col.UpdateId(id, update)
}

func NewBaseService(id interfaces.ModelId) (svc *BaseService) {
	if mongo.Client == nil {
		_ = mongo.InitMongo()
	}
	colName := getModelColName(id)
	col := mongo.GetMongoCol(colName)
	return &BaseService{
		id:  id,
		col: col,
	}
}
