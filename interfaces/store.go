package interfaces

import (
	"github.com/goava/di"
	"sync"
)

type Store interface {
	GetMap() *sync.Map
	Set(key string, c *di.Container) (err error)
	Get(key string) (c *di.Container, err error)
}
