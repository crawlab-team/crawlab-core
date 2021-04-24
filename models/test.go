package models

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"testing"
)

func setupTest(t *testing.T) {
	_ = mongo.InitMongo()
	_ = InitModelServices()
	t.Cleanup(cleanupTest)
}

func cleanupTest() {
	// system model services
	_ = ArtifactService.Delete(nil)
	_ = TagService.Delete(nil)

	// operation model services
	_ = NodeService.Delete(nil)
	_ = ProjectService.Delete(nil)
	_ = ScheduleService.Delete(nil)
	_ = SettingService.Delete(nil)
	_ = SpiderService.Delete(nil)
	_ = TaskService.Delete(nil)
	_ = JobService.Delete(nil)
	_ = TokenService.Delete(nil)
	_ = UserService.Delete(nil)
	_ = VariableService.Delete(nil)
}
