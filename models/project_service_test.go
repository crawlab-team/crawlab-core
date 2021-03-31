package models

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func cleanupTestService() {
	_ = ProjectService.Delete(nil)
}

func TestService_GetList(t *testing.T) {
	setupTest(t, cleanupTestService)

	n := 10
	for i := 0; i < n; i++ {
		p := Project{
			Name:        "test name",
			Description: "test description",
			Tags:        []string{"test tag"},
		}
		_ = p.Add()
	}

	data, err := ProjectService.GetList(nil, nil)
	require.Nil(t, err)
	require.Equal(t, n, len(data))
}
