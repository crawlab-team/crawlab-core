package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/model"
	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func cleanupTestProject() {
	_ = model.ProjectService.DeleteList(nil)
}

func TestProjectController_Put(t *testing.T) {
	setupTest(t)
	app := gin.New()
	app.PUT("/projects", ProjectController.Put)
	s := httptest.NewServer(app)
	e := httpexpect.New(t, s.URL)
	defer s.Close()

	p := model.Project{
		Name:        "test project",
		Description: "this is a test project",
		Tags:        []string{"tag 1", "tag 2"},
	}

	res := e.PUT("/projects").WithJSON(p).Expect().Status(http.StatusOK).JSON().Object()
	res.Path("$.data._id").NotNull()
	res.Path("$.data.name").Equal("test project")
	res.Path("$.data.description").Equal("this is a test project")
	res.Path("$.data.tags").Array().Element(0).Equal("tag 1")
	res.Path("$.data.tags").Array().Element(1).Equal("tag 2")

	cleanupTestProject()
}

func TestProjectController_PutList(t *testing.T) {
	setupTest(t)
	app := gin.New()
	app.GET("/projects", ProjectController.GetList)
	app.PUT("/projects/batch", ProjectController.PutList)
	s := httptest.NewServer(app)
	e := httpexpect.New(t, s.URL)
	defer s.Close()

	n := 10
	var docs []model.Project
	for i := 0; i < n; i++ {
		docs = append(docs, model.Project{
			Name:        fmt.Sprintf("project %d", i+1),
			Description: "this is a project",
			Tags:        []string{strconv.Itoa(i % 2)},
		})
	}

	e.PUT("/projects/batch").WithJSON(docs).Expect().Status(http.StatusOK)

	res := e.GET("/projects").
		WithQueryObject(entity.Pagination{Page: 1, Size: 100}).
		Expect().Status(http.StatusOK).
		JSON().Object()
	res.Path("$.data").Array().Length().Equal(n)

	cond := []entity.Condition{
		{Key: "tags", Op: constants.FilterOpIn, Value: []string{"0"}},
	}
	condBytes, err := json.Marshal(&cond)
	require.Nil(t, err)

	res = e.GET("/projects").
		WithQueryObject(entity.Pagination{Page: 1, Size: 100}).
		WithQuery("conditions", string(condBytes)).
		Expect().Status(http.StatusOK).
		JSON().Object()
	res.Path("$.data").Array().Length().Equal(n / 2)
	data := res.Path("$.data").Array()
	for i := 0; i < n/2; i++ {
		obj := data.Element(i)
		obj.Path("$.name").Equal(fmt.Sprintf("project %d", i*2+1))
		obj.Path("$.tags").Array().First().Equal("0")
	}

	cleanupTestProject()
}

func TestProjectController_Get(t *testing.T) {
	setupTest(t)
	app := gin.New()
	app.GET("/projects/:id", ProjectController.Get)
	app.PUT("/projects", ProjectController.Put)
	s := httptest.NewServer(app)
	e := httpexpect.New(t, s.URL)
	defer s.Close()

	p := model.Project{
		Name: "test project",
	}
	res := e.PUT("/projects").WithJSON(p).Expect().Status(http.StatusOK).JSON().Object()
	res.Path("$.data._id").NotNull()
	id := res.Path("$.data._id").String().Raw()
	oid, err := primitive.ObjectIDFromHex(id)
	require.Nil(t, err)
	require.False(t, oid.IsZero())

	res = e.GET("/projects/" + id).WithJSON(p).Expect().Status(http.StatusOK).JSON().Object()
	res.Path("$.data._id").NotNull()
	res.Path("$.data.name").Equal("test project")

	cleanupTestProject()
}

func TestProjectController_GetList(t *testing.T) {
	setupTest(t)
	app := gin.New()
	app.GET("/projects", ProjectController.GetList)
	app.PUT("/projects", ProjectController.Put)
	s := httptest.NewServer(app)
	e := httpexpect.New(t, s.URL)
	defer s.Close()

	n := 10

	for i := 0; i < n; i++ {
		p := model.Project{
			Name: fmt.Sprintf("test name %d", i+1),
			Tags: []string{"test tag"},
		}
		obj := e.PUT("/projects").WithJSON(p).Expect().Status(http.StatusOK).JSON().Object()
		obj.Path("$.data._id").NotNull()
	}

	f := entity.Filter{
		//IsOr: false,
		Conditions: []entity.Condition{
			{Key: "tags", Op: constants.FilterOpIn, Value: []string{"test tag"}},
		},
	}
	condBytes, err := json.Marshal(&f.Conditions)
	require.Nil(t, err)

	pagination := entity.Pagination{
		Page: 1,
		Size: 100,
	}

	res := e.GET("/projects").
		WithQuery("conditions", string(condBytes)).
		WithQueryObject(pagination).
		Expect().Status(http.StatusOK).JSON().Object()
	res.Path("$.data").Array().Length().Equal(n)
	res.Path("$.total").Number().Equal(n)

	data := res.Path("$.data").Array()
	for i := 0; i < n; i++ {
		obj := data.Element(i)
		obj.Path("$.name").Equal(fmt.Sprintf("test name %d", i+1))
		obj.Path("$.tags").Array().First().Equal("test tag")
	}

	cleanupTestProject()
}

func TestProjectController_Delete(t *testing.T) {
	setupTest(t)
	app := gin.New()
	app.GET("/projects/:id", ProjectController.Get)
	app.PUT("/projects", ProjectController.Put)
	app.DELETE("/projects/:id", ProjectController.Delete)
	s := httptest.NewServer(app)
	e := httpexpect.New(t, s.URL)
	defer s.Close()

	p := model.Project{
		Name:        "test project",
		Description: "this is a test project",
		Tags:        []string{"test tag"},
	}

	// add
	res := e.PUT("/projects").
		WithJSON(p).
		Expect().Status(http.StatusOK).
		JSON().Object()
	res.Path("$.data._id").NotNull()
	id := res.Path("$.data._id").String().Raw()
	oid, err := primitive.ObjectIDFromHex(id)
	require.Nil(t, err)
	require.False(t, oid.IsZero())

	// get
	res = e.GET("/projects/" + id).
		Expect().Status(http.StatusOK).
		JSON().Object()
	res.Path("$.data._id").NotNull()
	id = res.Path("$.data._id").String().Raw()
	oid, err = primitive.ObjectIDFromHex(id)
	require.Nil(t, err)
	require.False(t, oid.IsZero())

	// delete
	e.DELETE("/projects/" + id).
		Expect().Status(http.StatusOK).
		JSON().Object()

	// get
	e.GET("/projects/" + id).
		Expect().Status(http.StatusNotFound)

	cleanupTestProject()
}

func TestProjectController_DeleteList(t *testing.T) {
	setupTest(t)
	app := gin.New()
	app.GET("/projects", ProjectController.GetList)
	app.PUT("/projects/batch", ProjectController.PutList)
	app.DELETE("/projects", ProjectController.DeleteList)
	s := httptest.NewServer(app)
	e := httpexpect.New(t, s.URL)
	defer s.Close()

	n := 10
	var docs []model.Project
	for i := 0; i < n; i++ {
		docs = append(docs, model.Project{
			Name:        fmt.Sprintf("project %d", i+1),
			Description: "this is a project",
			Tags:        []string{strconv.Itoa(i % 2)},
		})
	}

	// add
	res := e.PUT("/projects/batch").WithJSON(docs).Expect().Status(http.StatusOK).
		JSON().Object()
	var ids []primitive.ObjectID
	data := res.Path("$.data").Array()
	for i := 0; i < n; i++ {
		obj := data.Element(i)
		id := obj.Path("$._id").String().Raw()
		oid, err := primitive.ObjectIDFromHex(id)
		require.Nil(t, err)
		require.False(t, oid.IsZero())
		ids = append(ids, oid)
	}

	// delete
	payload := entity.BatchRequestPayload{
		Ids: ids,
	}
	e.DELETE("/projects").
		WithJSON(payload).
		Expect().Status(http.StatusOK)

	// check
	res = e.GET("/projects").
		Expect().Status(http.StatusOK).JSON().Object()
	res.Path("$.data").Null()
	res.Path("$.total").Number().Equal(0)

	cleanupTestProject()
}

func TestProjectController_Post(t *testing.T) {
	setupTest(t)
	app := gin.New()
	app.GET("/projects/:id", ProjectController.Get)
	app.PUT("/projects", ProjectController.Put)
	app.POST("/projects/:id", ProjectController.Post)
	s := httptest.NewServer(app)
	e := httpexpect.New(t, s.URL)
	defer s.Close()

	p := model.Project{
		Name:        "old name",
		Description: "old description",
		Tags:        []string{"old tag"},
	}

	// add
	res := e.PUT("/projects").
		WithJSON(p).
		Expect().Status(http.StatusOK).
		JSON().Object()
	res.Path("$.data._id").NotNull()
	id := res.Path("$.data._id").String().Raw()
	oid, err := primitive.ObjectIDFromHex(id)
	require.Nil(t, err)
	require.False(t, oid.IsZero())

	// change object
	p.Id = oid
	p.Name = "new name"
	p.Description = "new description"
	p.Tags = []string{"new tag"}

	// update
	e.POST("/projects/" + id).
		WithJSON(p).
		Expect().Status(http.StatusOK)

	// check
	res = e.GET("/projects/" + id).Expect().Status(http.StatusOK).JSON().Object()
	res.Path("$.data._id").Equal(id)
	res.Path("$.data.name").Equal("new name")
	res.Path("$.data.description").Equal("new description")
	res.Path("$.data.tags").Equal([]string{"new tag"})

	cleanupTestProject()
}

func TestProjectController_PostList(t *testing.T) {
	setupTest(t)
	app := gin.New()
	app.GET("/projects", ProjectController.GetList)
	app.PUT("/projects/batch", ProjectController.PutList)
	app.POST("/projects", ProjectController.PostList)
	s := httptest.NewServer(app)
	e := httpexpect.New(t, s.URL)
	defer s.Close()

	n := 10
	var docs []model.Project
	for i := 0; i < n; i++ {
		docs = append(docs, model.Project{
			Name:        "old name",
			Description: "old description",
			Tags:        []string{"old tag"},
		})
	}

	// add
	res := e.PUT("/projects/batch").WithJSON(docs).Expect().Status(http.StatusOK).
		JSON().Object()
	var ids []primitive.ObjectID
	data := res.Path("$.data").Array()
	for i := 0; i < n; i++ {
		obj := data.Element(i)
		id := obj.Path("$._id").String().Raw()
		oid, err := primitive.ObjectIDFromHex(id)
		require.Nil(t, err)
		require.False(t, oid.IsZero())
		ids = append(ids, oid)
	}

	// update
	p := model.Project{
		Name:        "new name",
		Description: "new description",
		Tags:        []string{"new tag"},
	}
	dataBytes, err := json.Marshal(&p)
	require.Nil(t, err)
	payload := entity.BatchRequestPayloadWithStringData{
		Ids:  ids,
		Data: string(dataBytes),
	}
	e.POST("/projects").
		WithJSON(payload).
		Expect().Status(http.StatusOK)

	// check
	res = e.GET("/projects").
		WithQueryObject(entity.Pagination{
			Page: 1,
			Size: 100,
		}).
		Expect().Status(http.StatusOK).JSON().Object()
	data = res.Path("$.data").Array()
	for i := 0; i < n; i++ {
		obj := data.Element(i)
		obj.Path("$.name").Equal("new name")
		obj.Path("$.description").Equal("new description")
		obj.Path("$.tags").Equal([]string{"new tag"})
	}

	cleanupTestProject()
}
