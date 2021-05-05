package client

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClient(t *testing.T) {
	var err error
	c, err := NewClient()
	require.Nil(t, err)
	require.NotNil(t, c)
}
