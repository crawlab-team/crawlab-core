package controllers

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

// GetFilter Get entity.Filter from gin.Context
func GetFilter(c *gin.Context) (f *entity.Filter, err error) {
	// bind
	condStr := c.Query(constants.FilterQueryFieldConditions)
	var conditions []entity.Condition
	if err := json.Unmarshal([]byte(condStr), &conditions); err != nil {
		return nil, err
	}

	// attempt to convert object id
	for i, cond := range conditions {
		switch cond.Value.(type) {
		case string:
			id, err := primitive.ObjectIDFromHex(cond.Value.(string))
			if err == nil {
				conditions[i].Value = id
			}
		}
	}

	return &entity.Filter{
		IsOr:       false,
		Conditions: conditions,
	}, nil
}

// GetFilterQuery Get bson.M from gin.Context
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

// FilterToQuery Translate entity.Filter to bson.M
func FilterToQuery(f *entity.Filter) (q bson.M, err error) {
	q = bson.M{}
	for _, cond := range f.Conditions {
		switch cond.Op {
		case constants.FilterOpNotSet:
			// do nothing
		case constants.FilterOpEqual:
			q[cond.Key] = cond.Value
		case constants.FilterOpNotEqual:
			q[cond.Key] = bson.M{"$ne": cond.Value}
		case constants.FilterOpContains, constants.FilterOpRegex, constants.FilterOpSearch:
			q[cond.Key] = bson.M{"$regex": cond.Value}
		case constants.FilterOpNotContains:
			q[cond.Key] = bson.M{"$not": bson.M{"$regex": cond.Value}}
		case constants.FilterOpIn:
			q[cond.Key] = bson.M{"$in": cond.Value}
		case constants.FilterOpNotIn:
			q[cond.Key] = bson.M{"$nin": cond.Value}
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

// GetFilterAll Get all from gin.Context
func GetFilterAll(c *gin.Context) (res bool, err error) {
	resStr := c.Query(constants.FilterQueryFieldAll)
	switch strings.ToUpper(resStr) {
	case "1":
		return true, nil
	case "0":
		return false, nil
	case "Y":
		return true, nil
	case "N":
		return false, nil
	case "T":
		return true, nil
	case "F":
		return false, nil
	case "TRUE":
		return true, nil
	case "FALSE":
		return false, nil
	default:
		return false, errors.ErrorFilterInvalidOperation
	}
}

func MustGetFilterAll(c *gin.Context) (res bool) {
	res, err := GetFilterAll(c)
	if err != nil {
		return false
	}
	return res
}
