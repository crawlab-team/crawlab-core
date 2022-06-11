package apps

import (
	"github.com/imroc/req"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestServer_Start(t *testing.T) {
	svr := GetServer()

	// start
	go Start(svr)

	res, err := req.Get("http://localhost:8000/system-info")
	require.Nil(t, err)
	resStr, err := res.ToString()
	require.Nil(t, err)
	require.Contains(t, resStr, "success")
}
