package controllers

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

var FilterController ActionController

func getFilterActions() []Action {
	ctx := newFilterContext()
	return []Action{
		{
			Method:      http.MethodGet,
			Path:        "/:col/:field",
			HandlerFunc: ctx.getColFieldOptions,
		},
	}
}

type filterContext struct {
}

func (ctx *filterContext) getColFieldOptions(c *gin.Context) {
	colName := c.Param("col")
	field := c.Param("field")
	query := MustGetFilterQuery(c)
	pipelines := mongo2.Pipeline{}
	if query != nil {
		pipelines = append(pipelines, bson.D{{"$match", query}})
	}
	pipelines = append(pipelines, bson.D{{"$group", bson.D{{"_id", "$" + field}}}})
	var results []struct {
		Id interface{} `bson:"_id"`
	}
	if err := mongo.GetMongoCol(colName).Aggregate(pipelines, nil).All(&results); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	var values []interface{}
	for _, result := range results {
		values = append(values, result.Id)
	}
	HandleSuccessWithData(c, values)
}

func newFilterContext() *filterContext {
	return &filterContext{}
}
