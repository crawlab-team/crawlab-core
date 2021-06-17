package test

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/interfaces"
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

func TestUserService_Create_Login_CheckToken(t *testing.T) {
	var err error
	T.Setup(t)

	err = T.userSvc.Create(&interfaces.UserCreateOptions{
		Username: T.TestUsername,
		Password: T.TestPassword,
	})
	require.Nil(t, err)

	u, err := T.modelSvc.GetUserByUsername(T.TestUsername, nil)
	require.Nil(t, err)
	require.Equal(t, T.TestUsername, u.Username)
	require.Equal(t, utils.EncryptPassword(T.TestPassword), u.Password)

	token, u2, err := T.userSvc.Login(&interfaces.UserLoginOptions{
		Username: T.TestUsername,
		Password: T.TestPassword,
	})
	require.Nil(t, err)
	require.Greater(t, len(token), 10)
	require.Equal(t, u.Username, u2.GetUsername())

	u3, err := T.userSvc.CheckToken(token)
	require.Nil(t, err)
	require.Equal(t, u.Username, u3.GetUsername())
}

func TestUserService_ChangePassword(t *testing.T) {
	var err error
	T.Setup(t)

	u, err := T.modelSvc.GetUserByUsername(constants.DefaultAdminUsername, nil)
	require.Nil(t, err)
	err = T.userSvc.ChangePassword(u.Id, T.TestNewPassword)
	require.Nil(t, err)

	u2, err := T.modelSvc.GetUserByUsername(constants.DefaultAdminUsername, nil)
	require.Nil(t, err)
	require.Equal(t, utils.EncryptPassword(T.TestNewPassword), u2.Password)
}
