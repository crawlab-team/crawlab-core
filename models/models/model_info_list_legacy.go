package models

import (
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/interfaces"
)

var ModelInfoList = []entity.ModelInfo{
	{Id: interfaces.ModelIdArtifact, ColName: interfaces.ModelColNameArtifact},
	{Id: interfaces.ModelIdTag, ColName: interfaces.ModelColNameTag},
	{Id: interfaces.ModelIdNode, ColName: interfaces.ModelColNameNode},
	{Id: interfaces.ModelIdProject, ColName: interfaces.ModelColNameProject},
	{Id: interfaces.ModelIdSpider, ColName: interfaces.ModelColNameSpider},
	{Id: interfaces.ModelIdTask, ColName: interfaces.ModelColNameTask},
	{Id: interfaces.ModelIdJob, ColName: interfaces.ModelColNameJob},
	{Id: interfaces.ModelIdSchedule, ColName: interfaces.ModelColNameSchedule},
	{Id: interfaces.ModelIdUser, ColName: interfaces.ModelColNameUser},
	{Id: interfaces.ModelIdSetting, ColName: interfaces.ModelColNameSetting},
	{Id: interfaces.ModelIdToken, ColName: interfaces.ModelColNameToken},
	{Id: interfaces.ModelIdVariable, ColName: interfaces.ModelColNameVariable},
}
