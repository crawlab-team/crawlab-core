package test

import "testing"

type Test struct {
}

func (t *Test) Setup(t2 *testing.T) {
	t2.Cleanup(t.Cleanup)
}

func (t *Test) Cleanup() {
}
