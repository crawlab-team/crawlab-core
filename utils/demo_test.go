package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestImportDemo(t *testing.T) {
	err := ImportDemo()
	require.Nil(t, err)
}
