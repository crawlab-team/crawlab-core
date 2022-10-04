package admin

import (
	mock_interfaces "github.com/crawlab-team/crawlab-core/mock"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestService_Schedule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//mockModelSvc := mock_interfaces.

	mockTaskSchedulerSvc := mock_interfaces.NewMockTaskSchedulerService(ctrl)
	mockTaskSchedulerSvc.EXPECT().Schedule(gomock.Any())
}
