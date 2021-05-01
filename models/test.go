package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"testing"
)

func setupTest(t *testing.T) {
	if err := mongo.InitMongo(); err != nil {
		panic(err)
	}
	if err := InitModels(); err != nil {
		panic(err)
	}
	t.Cleanup(cleanupTest)
}

func cleanupTest() {
	// system model services
	_ = MustGetService(interfaces.ModelIdArtifact).Delete(nil)
	_ = MustGetService(interfaces.ModelIdTag).Delete(nil)

	// operation model services
	_ = MustGetService(interfaces.ModelIdNode).Delete(nil)
	_ = MustGetService(interfaces.ModelIdProject).Delete(nil)
	_ = MustGetService(interfaces.ModelIdSchedule).Delete(nil)
	_ = MustGetService(interfaces.ModelIdSetting).Delete(nil)
	_ = MustGetService(interfaces.ModelIdSpider).Delete(nil)
	_ = MustGetService(interfaces.ModelIdTask).Delete(nil)
	_ = MustGetService(interfaces.ModelIdJob).Delete(nil)
	_ = MustGetService(interfaces.ModelIdToken).Delete(nil)
	_ = MustGetService(interfaces.ModelIdUser).Delete(nil)
	_ = MustGetService(interfaces.ModelIdVariable).Delete(nil)
}

func CleanupTest() {
	cleanupTest()
}
