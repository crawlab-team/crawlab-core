package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Result map[string]interface{}

func (r *Result) GetId() (id primitive.ObjectID) {
	res, ok := r.Value()["_id"]
	if ok {
		id, ok = res.(primitive.ObjectID)
		if ok {
			return id
		}
	}
	return id
}

func (r *Result) SetId(id primitive.ObjectID) {
	(*r)["_id"] = id
}

func (r *Result) Value() map[string]interface{} {
	return *r
}

func (r *Result) SetValue(key string, value interface{}) {
	(*r)[key] = value
}

func (r *Result) GetValue(key string) (value interface{}) {
	return (*r)[key]
}
