package grpc

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewService(t *testing.T) {
	svc, err := NewService(nil)
	require.Nil(t, err)
	require.NotNil(t, svc)
}

func TestService_AddClient(t *testing.T) {
	svc, err := NewService(nil)
	require.Nil(t, err)
	require.NotNil(t, svc)

	// add client
	err = svc.AddClient(nil)
	require.Nil(t, err)

	// get client
	client, err := svc.GetClient(NewAddress(nil))
	require.Nil(t, err)
	require.NotNil(t, client)
}

func TestService_DeleteClient(t *testing.T) {
	svc, err := NewService(nil)
	require.Nil(t, err)
	require.NotNil(t, svc)

	// add client
	err = svc.AddClient(nil)
	require.Nil(t, err)

	// remove client
	err = svc.DeleteClient(NewAddress(nil))
	require.Nil(t, err)

	// get client
	_, err = svc.GetClient(NewAddress(nil))
	require.Equal(t, errors.ErrorGrpcClientNotExists, err)
}
