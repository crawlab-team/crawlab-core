package model

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func setupJobTest() (err error) {
	return nil
}

func cleanupJobTest() {
}

func TestJobAdd(t *testing.T) {
	err := setupJobTest()
	require.Nil(t, err)

	j := Job{}

	err = j.Add()
	require.Nil(t, err)
}
