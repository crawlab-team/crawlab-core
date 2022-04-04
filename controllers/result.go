package controllers

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/result"
	"github.com/crawlab-team/crawlab-db/generic"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

var ResultController ActionController

func getResultActions() []Action {
	var resultCtx = newResultContext()
	return []Action{
		{
			Method:      http.MethodGet,
			Path:        "/:id",
			HandlerFunc: resultCtx.getList,
		},
	}
}

type resultContext struct {
}

func (ctx *resultContext) getList(c *gin.Context) {
	// id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// service
	svc, err := result.GetResultService(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// params
	pagination := MustGetPagination(c)
	query := generic.ListQuery{} // TODO: implement query

	// get results
	data, err := svc.List(query, &generic.ListOptions{
		Sort:  []generic.ListSort{{"_id", generic.SortDirectionDesc}},
		Skip:  pagination.Size * (pagination.Page - 1),
		Limit: pagination.Size,
	})
	if err != nil {
		if err.Error() == mongo2.ErrNoDocuments.Error() {
			HandleSuccessWithListData(c, nil, 0)
			return
		}
		HandleErrorInternalServerError(c, err)
		return
	}

	// validate results
	if len(data) == 0 {
		HandleSuccessWithListData(c, nil, 0)
		return
	}

	// total count
	total, err := svc.Count(query)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// response
	HandleSuccessWithListData(c, data, total)
}

func (ctx *resultContext) _getSvc(id primitive.ObjectID) (svc interfaces.ResultService, err error) {
	return result.GetResultService(id)
}

func newResultContext() *resultContext {
	// context
	ctx := &resultContext{}

	return ctx
}
