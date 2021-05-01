package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func setupTaskTest() (err error) {
	return mongo.InitMongo()
}

func cleanupTaskTest() {
	_ = mongo.GetMongoCol(interfaces.ModelColNameTask).Delete(nil)
	_ = mongo.GetMongoCol(interfaces.ModelColNameArtifact).Delete(nil)
}

func TestTask_Add(t *testing.T) {
	var err error
	setupTest(t)

	task := Task{}

	err = task.Add()
	require.Nil(t, err)
	require.NotNil(t, task.Id)
}

func TestTask_Save(t *testing.T) {
	var err error
	setupTest(t)

	task := Task{}
	spider := Spider{
		Name: "test_task",
	}
	err = spider.Add()
	require.Nil(t, err)

	err = task.Add()
	require.Nil(t, err)

	task.SpiderId = spider.Id
	err = task.Save()
	require.Nil(t, err)

	err = mongo.GetMongoCol(interfaces.ModelColNameTask).FindId(task.Id).One(&task)
	require.Nil(t, err)
	require.Equal(t, spider.Id, task.SpiderId)

	err = mongo.GetMongoCol(interfaces.ModelColNameSpider).FindId(task.SpiderId).One(&spider)
	require.Nil(t, err)
	require.Equal(t, spider.Id, task.SpiderId)
	require.Equal(t, "test_task", spider.Name)
}

func TestTask_Delete(t *testing.T) {
	var err error
	setupTest(t)

	id := primitive.NewObjectID()
	task := Task{
		SpiderId: id,
	}

	err = task.Add()
	require.Nil(t, err)

	err = task.Delete()
	require.Nil(t, err)

	var a Artifact
	col := mongo.GetMongoCol(interfaces.ModelColNameArtifact)
	err = col.FindId(task.Id).One(&a)
	require.Nil(t, err)
	require.NotNil(t, a.Obj)
	require.True(t, a.Del)

	data, err := bson.Marshal(&a.Obj)
	require.Nil(t, err)
	err = bson.Unmarshal(data, &task)
	require.Nil(t, err)
	require.Equal(t, id, task.SpiderId)
}

func TestTask_DeleteList(t *testing.T) {
	var err error
	setupTest(t)

	doc := Task{}

	err = doc.Add()
	require.Nil(t, err)

	err = MustGetService(interfaces.ModelIdTask).DeleteList(nil)
	require.Nil(t, err)

	total, err := MustGetService(interfaces.ModelIdTask).Count(nil)
	require.Equal(t, 0, total)
}
