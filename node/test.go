package node

import (
	"os"
	"testing"
)

var TestService *service

func setupTest(t *testing.T) {
	var err error
	TestService, err = NewService(nil)
	if err != nil {
		panic(err)
	}
	t.Cleanup(cleanupTest)
}

func cleanupTest() {
	_ = os.RemoveAll(defaultConfigDirPath)
}
