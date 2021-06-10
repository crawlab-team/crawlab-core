package service

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/emirpasic/gods/lists/arraylist"
	"go.mongodb.org/mongo-driver/bson"
)

func (svc *Service) serializeList(list arraylist.List, target interface{}) (err error) {
	// bytes
	bytes, err := list.ToJSON()
	if err != nil {
		return err
	}

	// unmarshal
	if err := json.Unmarshal(bytes, target); err != nil {
		return err
	}

	return nil
}

func (svc *Service) getListSerializeTarget(id interfaces.ModelId, query bson.M, opts *mongo.FindOptions, target interface{}) (err error) {
	list, err := svc.GetBaseService(id).GetList(query, opts)
	if err != nil {
		return err
	}
	return svc.serializeList(list, target)
}
