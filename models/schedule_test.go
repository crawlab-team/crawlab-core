package models

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/stretchr/testify/require"
	"testing"
)

func setupScheduleTest() (err error) {
	return mongo.InitMongo()
}

func cleanupScheduleTest() {
	_ = mongo.GetMongoCol(ModelColNameSchedule).Delete(nil)
	_ = mongo.GetMongoCol(ArtifactColName).Delete(nil)
}

func TestSchedule_Add(t *testing.T) {
	err := setupScheduleTest()
	require.Nil(t, err)

	s := Schedule{}

	err = s.Add()
	require.Nil(t, err)
	require.NotNil(t, s.Id)

	a, err := s.GetArtifact()
	require.Nil(t, err)
	require.Equal(t, s.Id, a.Id)
	require.NotNil(t, a.CreateTs)
	require.NotNil(t, a.UpdateTs)

	col := mongo.GetMongoCol(ModelColNameSchedule)
	col.GetContext()

	cleanupScheduleTest()
}

func TestSchedule_Save(t *testing.T) {
	err := setupScheduleTest()
	require.Nil(t, err)

	s := Schedule{}

	err = s.Add()
	require.Nil(t, err)

	name := "test_schedule"
	s.Name = name
	err = s.Save()
	require.Nil(t, err)

	err = mongo.GetMongoCol(ModelColNameSchedule).FindId(s.Id).One(&s)
	require.Nil(t, err)
	require.Equal(t, name, s.Name)

	cleanupScheduleTest()
}

func TestSchedule_Delete(t *testing.T) {
	err := setupScheduleTest()
	require.Nil(t, err)

	s := Schedule{
		Name: "test_schedule",
	}

	err = s.Add()
	require.Nil(t, err)

	err = s.Delete()
	require.Nil(t, err)

	var a Artifact
	col := mongo.GetMongoCol(ArtifactColName)
	err = col.FindId(s.Id).One(&a)
	require.Nil(t, err)
	require.NotNil(t, a.Obj)
	require.True(t, a.Del)

	cleanupScheduleTest()
}

func TestSchedule_DeleteList(t *testing.T) {
	err := setupScheduleTest()
	require.Nil(t, err)

	doc := Schedule{
		Name: "test_Schedule",
	}

	err = doc.Add()
	require.Nil(t, err)

	err = ScheduleService.DeleteList(nil)
	require.Nil(t, err)

	total, err := ScheduleService.Count(nil)
	require.Equal(t, 0, total)

	cleanupScheduleTest()
}
