package models

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func setupJobTest() (err error) {
	return mongo.InitMongo()
}

func cleanupJobTest() {
	_ = mongo.GetMongoCol(ModelColNameJob).Delete(nil)
	_ = mongo.GetMongoCol(ArtifactColName).Delete(nil)
}

func TestJob_Add(t *testing.T) {
	err := setupJobTest()
	require.Nil(t, err)

	j := Job{}

	err = j.Add()
	require.Nil(t, err)
	require.NotNil(t, j.Id)

	a, err := j.GetArtifact()
	require.Nil(t, err)
	require.Equal(t, j.Id, a.Id)
	require.NotNil(t, a.CreateTs)
	require.NotNil(t, a.UpdateTs)

	col := mongo.GetMongoCol(ModelColNameJob)
	col.GetContext()

	cleanupJobTest()
}

func TestJob_Save(t *testing.T) {
	err := setupJobTest()
	require.Nil(t, err)

	j := Job{}

	err = j.Add()
	require.Nil(t, err)

	id := primitive.NewObjectID()
	j.TaskId = id
	err = j.Save()
	require.Nil(t, err)

	err = mongo.GetMongoCol(ModelColNameJob).FindId(j.Id).One(&j)
	require.Nil(t, err)
	require.Equal(t, id, j.TaskId)

	cleanupJobTest()
}

func TestJob_Delete(t *testing.T) {
	err := setupJobTest()
	require.Nil(t, err)

	j := Job{
		TaskId: primitive.NewObjectID(),
	}

	err = j.Add()
	require.Nil(t, err)

	err = j.Delete()
	require.Nil(t, err)

	var a Artifact
	col := mongo.GetMongoCol(ArtifactColName)
	err = col.FindId(j.Id).One(&a)
	require.Nil(t, err)
	require.NotNil(t, a.Obj)
	require.True(t, a.Del)

	cleanupJobTest()
}

func TestJob_DeleteList(t *testing.T) {
	err := setupJobTest()
	require.Nil(t, err)

	doc := Job{
		TaskId: primitive.NewObjectID(),
	}

	err = doc.Add()
	require.Nil(t, err)

	err = JobService.DeleteList(nil)
	require.Nil(t, err)

	total, err := JobService.Count(nil)
	require.Equal(t, 0, total)

	cleanupJobTest()
}
