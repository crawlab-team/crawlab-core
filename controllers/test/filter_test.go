package test

import (
	"encoding/json"
	"fmt"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"testing"
)

func init() {
	viper.Set("mongo.db", "crawlab_test")
}

func TestFilterController_GetColFieldOptions(t *testing.T) {
	T.Setup(t)
	e := T.NewExpect(t)

	// mongo collection
	colName := "test_collection_for_filter"
	field1 := "field1"
	field2 := "field2"
	value1 := "value1"
	col := mongo.GetMongoCol(colName)
	n := 10
	for i := 0; i < n; i++ {
		_, err := col.Insert(bson.M{field1: value1, field2: i % 2})
		require.Nil(t, err)
	}

	// validate filter options field 1
	res := T.WithAuth(e.GET(fmt.Sprintf("/filters/%s/%s", colName, field1))).
		Expect().Status(http.StatusOK).JSON().Object()
	res.Path("$.data").NotNull()
	res.Path("$.data").Array().Length().Equal(1)

	// validate filter options field 2
	res = T.WithAuth(e.GET(fmt.Sprintf("/filters/%s/%s", colName, field2))).
		Expect().Status(http.StatusOK).JSON().Object()
	res.Path("$.data").NotNull()
	res.Path("$.data").Array().Length().Equal(2)

	// validate filter options with query
	conditions := []entity.Condition{{field2, constants.FilterOpEqual, 0}}
	conditionsJson, err := json.Marshal(conditions)
	conditionsJsonStr := string(conditionsJson)
	require.Nil(t, err)
	res = T.WithAuth(e.GET(fmt.Sprintf("/filters/%s/%s", colName, field2))).
		WithQuery(constants.FilterQueryFieldConditions, conditionsJsonStr).
		Expect().Status(http.StatusOK).JSON().Object()
	res.Path("$.data").NotNull()
	res.Path("$.data").Array().Length().Equal(1)
}
