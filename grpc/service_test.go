package grpc

import (
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewService(t *testing.T) {
	svc, err := NewService(nil)
	require.Nil(t, err)
	require.NotNil(t, svc)
	err = svc.Stop()
	require.Nil(t, err)
}

func TestService_AddClient(t *testing.T) {
	svc, err := NewService(nil)
	require.Nil(t, err)
	require.NotNil(t, svc)

	// add client
	err = svc.AddClient(nil)
	require.Nil(t, err)

	// get client
	client, err := svc.GetClient(entity.NewAddress(nil))
	require.Nil(t, err)
	require.NotNil(t, client)

	err = svc.Stop()
	require.Nil(t, err)
}

func TestService_DeleteClient(t *testing.T) {
	svc, err := NewService(nil)
	require.Nil(t, err)
	require.NotNil(t, svc)

	// add client
	err = svc.AddClient(nil)
	require.Nil(t, err)

	// delete client
	err = svc.DeleteClient(entity.NewAddress(nil))
	require.Nil(t, err)

	// get client
	c, err := svc.GetClient(entity.NewAddress(nil))
	require.Nil(t, c)
	require.Equal(t, errors.ErrorGrpcClientNotExists, err)

	err = svc.Stop()
	require.Nil(t, err)
}
