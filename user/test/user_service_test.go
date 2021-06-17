package test

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserService_Init(t *testing.T) {
	var err error
	T.Setup(t)

	u, err := T.modelSvc.GetUserByUsername(constants.DefaultAdminUsername, nil)
	require.Nil(t, err)
	require.Equal(t, constants.DefaultAdminUsername, u.Username)
	require.Equal(t, utils.EncryptPassword(constants.DefaultAdminPassword), u.Password)
}
