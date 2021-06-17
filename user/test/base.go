package test

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/user"
	"go.uber.org/dig"
	"testing"
)

func init() {
	var err error
	T, err = NewTest()
	if err != nil {
		panic(err)
	}
}

var T *Test

type Test struct {
	// dependencies
	modelSvc service.ModelService
	userSvc  interfaces.UserService

	// test data
	TestUsername string
	TestPassword string
}

func (t *Test) Setup(t2 *testing.T) {
	t2.Cleanup(t.Cleanup)
}

func (t *Test) Cleanup() {
	_ = t.modelSvc.GetBaseService(interfaces.ModelIdTask).Delete(nil)
}

func NewTest() (t *Test, err error) {
	// test
	t = &Test{
		TestUsername: "test_username",
		TestPassword: "test_password",
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.GetService); err != nil {
		return nil, err
	}
	if err := c.Provide(user.GetUserService); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(modelSvc service.ModelService, userSvc interfaces.UserService) {
		t.modelSvc = modelSvc
		t.userSvc = userSvc
	}); err != nil {
		return nil, err
	}

	return t, nil
}
