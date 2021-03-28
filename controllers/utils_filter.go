package controllers

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// Get entity.Filter from gin.Context
func GetFilter(c *gin.Context) (f *entity.Filter, err error) {
	condStr := c.Query("conditions")
	var conditions []entity.Condition
	if err := json.Unmarshal([]byte(condStr), &conditions); err != nil {
		return nil, err
	}
	return &entity.Filter{
		IsOr:       false,
		Conditions: conditions,
	}, nil
}

// Get bson.M from gin.Context
func GetFilterQuery(c *gin.Context) (q bson.M, err error) {
	f, err := GetFilter(c)
	if err != nil {
		return nil, err
	}

	if f == nil {
		return nil, nil
	}

	// TODO: implement logic OR

	return FilterToQuery(f)
}

func MustGetFilterQuery(c *gin.Context) (q bson.M) {
	q, err := GetFilterQuery(c)
	if err != nil {
		return nil
	}
	return q
}

// Translate entity.Filter to bson.M
func FilterToQuery(f *entity.Filter) (q bson.M, err error) {
	q = bson.M{}
	for _, cond := range f.Conditions {
		switch cond.Op {
		case constants.FilterOpEqual:
			q[cond.Key] = cond.Value
		case constants.FilterOpContains, constants.FilterOpRegex:
			q[cond.Key] = bson.M{"$regex": cond.Value}
		case constants.FilterOpIn:
			q[cond.Key] = bson.M{"$in": cond.Value}
		case constants.FilterOpGreaterThan:
			q[cond.Key] = bson.M{"$gt": cond.Value}
		case constants.FilterOpGreaterThanEqual:
			q[cond.Key] = bson.M{"$gte": cond.Value}
		case constants.FilterOpLessThan:
			q[cond.Key] = bson.M{"$lt": cond.Value}
		case constants.FilterOpLessThanEqual:
			q[cond.Key] = bson.M{"$lte": cond.Value}
		default:
			return nil, errors.ErrorFilterInvalidOperation
		}
	}
	return q, nil
}
