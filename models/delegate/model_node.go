package delegate

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"time"
)

type ModelNodeDelegate struct {
	n interfaces.Node
	interfaces.ModelDelegate
}

func (d *ModelNodeDelegate) UpdateStatus(active bool, activeTs time.Time, status string) (err error) {
	d.n.SetActive(active)
	d.n.SetActiveTs(activeTs)
	d.n.SetStatus(status)
	return d.Save()
}

func NewModelNodeDelegate(n interfaces.Node) interfaces.ModelNodeDelegate {
	return &ModelNodeDelegate{
		n:             n,
		ModelDelegate: NewModelDelegate(n),
	}
}
