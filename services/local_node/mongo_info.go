package local_node

import (
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/crawlab-core/models"
	"go.uber.org/atomic"
	"sync"
	"time"
)

var locker atomic.Int32

type mongo struct {
	node *models.Node
	sync.RWMutex
}

func (n *mongo) load(retry bool) (err error) {
	n.Lock()
	defer n.Unlock()
	var node models.Node
	if retry {
		b := backoff.NewConstantBackOff(1 * time.Second)
		err = backoff.Retry(func() error {
			node, err = models.GetNodeByKey(GetLocalNode().Identify)
			if err != nil {
				log.WithError(err).Warnf("Get current node info from database failed.  Will after %f seconds, try again.", b.NextBackOff().Seconds())
			}
			return err
		}, b)
	} else {
		node, err = models.GetNodeByKey(localNode.Identify)
	}

	if err != nil {
		return
	}
	n.node = &node
	return nil
}
func (n *mongo) watch() {
	timer := time.NewTicker(time.Second * 5)
	for range timer.C {
		if locker.CAS(0, 1) {

			err := n.load(false)

			if err != nil {
				log.WithError(err).Errorf("load current node from database failed")
			}
			locker.Store(0)
		}
		continue
	}
}

func (n *mongo) Current() *models.Node {
	n.RLock()
	defer n.RUnlock()
	return n.node
}
