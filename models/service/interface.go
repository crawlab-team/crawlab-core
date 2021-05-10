package service

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	models2 "github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModelService interface {
	interfaces.ModelService
	GetNodeById(id primitive.ObjectID) (res *models2.Node, err error)
	GetNode(query bson.M, opts *mongo.FindOptions) (res *models2.Node, err error)
	GetNodeList(query bson.M, opts *mongo.FindOptions) (res []models2.Node, err error)
	GetNodeByKey(key string, opts *mongo.FindOptions) (res *models2.Node, err error)
	GetProjectById(id primitive.ObjectID) (res *models2.Project, err error)
	GetProject(query bson.M, opts *mongo.FindOptions) (res *models2.Project, err error)
	GetProjectList(query bson.M, opts *mongo.FindOptions) (res []models2.Project, err error)
	GetArtifactById(id primitive.ObjectID) (res *models2.Artifact, err error)
	GetArtifact(query bson.M, opts *mongo.FindOptions) (res *models2.Artifact, err error)
	GetArtifactList(query bson.M, opts *mongo.FindOptions) (res []models2.Artifact, err error)
	GetTagById(id primitive.ObjectID) (res *models2.Tag, err error)
	GetTag(query bson.M, opts *mongo.FindOptions) (res *models2.Tag, err error)
	GetTagList(query bson.M, opts *mongo.FindOptions) (res []models2.Tag, err error)
	GetTagIds(colName string, tags []interfaces.Tag) (tagIds []primitive.ObjectID, err error)
	UpdateTagsById(colName string, id primitive.ObjectID, tags []interfaces.Tag) (tagIds []primitive.ObjectID, err error)
	UpdateTags(colName string, query bson.M, tags []interfaces.Tag) (tagIds []primitive.ObjectID, err error)
	GetJobById(id primitive.ObjectID) (res *models2.Job, err error)
	GetJob(query bson.M, opts *mongo.FindOptions) (res *models2.Job, err error)
	GetJobList(query bson.M, opts *mongo.FindOptions) (res []models2.Job, err error)
	GetScheduleById(id primitive.ObjectID) (res *models2.Schedule, err error)
	GetSchedule(query bson.M, opts *mongo.FindOptions) (res *models2.Schedule, err error)
	GetScheduleList(query bson.M, opts *mongo.FindOptions) (res []models2.Schedule, err error)
	GetUserById(id primitive.ObjectID) (res *models2.User, err error)
	GetUser(query bson.M, opts *mongo.FindOptions) (res *models2.User, err error)
	GetUserList(query bson.M, opts *mongo.FindOptions) (res []models2.User, err error)
	GetUserByUsername(username string, opts *mongo.FindOptions) (res *models2.User, err error)
	GetSettingById(id primitive.ObjectID) (res *models2.Setting, err error)
	GetSetting(query bson.M, opts *mongo.FindOptions) (res *models2.Setting, err error)
	GetSettingList(query bson.M, opts *mongo.FindOptions) (res []models2.Setting, err error)
	GetSettingByKey(key string, opts *mongo.FindOptions) (res *models2.Setting, err error)
	GetSpiderById(id primitive.ObjectID) (res *models2.Spider, err error)
	GetSpider(query bson.M, opts *mongo.FindOptions) (res *models2.Spider, err error)
	GetSpiderList(query bson.M, opts *mongo.FindOptions) (res []models2.Spider, err error)
	GetTaskById(id primitive.ObjectID) (res *models2.Task, err error)
	GetTask(query bson.M, opts *mongo.FindOptions) (res *models2.Task, err error)
	GetTaskList(query bson.M, opts *mongo.FindOptions) (res []models2.Task, err error)
	GetTokenById(id primitive.ObjectID) (res *models2.Token, err error)
	GetToken(query bson.M, opts *mongo.FindOptions) (res *models2.Token, err error)
	GetTokenList(query bson.M, opts *mongo.FindOptions) (res []models2.Token, err error)
	GetVariableById(id primitive.ObjectID) (res *models2.Variable, err error)
	GetVariable(query bson.M, opts *mongo.FindOptions) (res *models2.Variable, err error)
	GetVariableList(query bson.M, opts *mongo.FindOptions) (res []models2.Variable, err error)
	GetVariableByKey(key string, opts *mongo.FindOptions) (res *models2.Variable, err error)
	DropAll() (err error)
}
