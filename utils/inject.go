package utils

import (
	"github.com/crawlab-team/crawlab-core/inject"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/goava/di"
	"os"
	"strings"
)

func GetContainer(env string) (c *di.Container, err error) {
	return inject.Store.Get(env)
}

func resolveModule(c *di.Container, target di.Pointer, opts ...di.ResolveOption) (err error) {
	return c.Resolve(target, opts...)
}

func ResolveModule(env string, target di.Pointer, opts ...di.ResolveOption) (err error) {
	c, err := GetContainer(env)
	if err != nil {
		return err
	}
	return resolveModule(c, target, opts...)
}

func MustResolveModule(env string, target di.Pointer, opts ...di.ResolveOption) {
	if err := ResolveModule(env, target, opts...); err != nil {
		panic(err)
	}
}

func ProvideModule(provide interfaces.Provide) {
	// default
	provide("")

	// test
	testEnv := os.Getenv("TEST_ENV")
	if testEnv == "" {
		return
	}
	envs := strings.Split(testEnv, ",")
	for _, env := range envs {
		provide(env)
	}
}
