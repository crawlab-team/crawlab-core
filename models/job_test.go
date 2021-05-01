package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestJob_Add(t *testing.T) {
	setupTest(t)

	j := Job{}

	err := j.Add()
	require.Nil(t, err)
	require.NotNil(t, j.Id)

	a, err := j.GetArtifact()
	require.Nil(t, err)
	require.Equal(t, j.Id, a.GetId())
	require.NotNil(t, a.GetSys().GetCreateTs())
	require.NotNil(t, a.GetSys().GetUpdateTs())

	col := mongo.GetMongoCol(interfaces.ModelColNameJob)
	col.GetContext()
}

func TestJob_Save(t *testing.T) {
	setupTest(t)

	j := Job{}

	err := j.Add()
	require.Nil(t, err)

	id := primitive.NewObjectID()
	j.TaskId = id
	err = j.Save()
	require.Nil(t, err)

	err = mongo.GetMongoCol(interfaces.ModelColNameJob).FindId(j.Id).One(&j)
	require.Nil(t, err)
	require.Equal(t, id, j.TaskId)
}

func TestJob_Delete(t *testing.T) {
	setupTest(t)

	j := Job{
		TaskId: primitive.NewObjectID(),
	}

	err := j.Add()
	require.Nil(t, err)

	err = j.Delete()
	require.Nil(t, err)

	var a Artifact
	col := mongo.GetMongoCol(interfaces.ModelColNameArtifact)
	err = col.FindId(j.Id).One(&a)
	require.Nil(t, err)
	require.NotNil(t, a.Obj)
	require.True(t, a.Del)
}

func TestJob_DeleteList(t *testing.T) {
	setupTest(t)

	doc := Job{
		TaskId: primitive.NewObjectID(),
	}

	err := doc.Add()
	require.Nil(t, err)

	err = MustGetService(interfaces.ModelIdJob).DeleteList(nil)
	require.Nil(t, err)

	total, err := MustGetService(interfaces.ModelIdJob).Count(nil)
	require.Equal(t, 0, total)
}
