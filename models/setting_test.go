package models

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/stretchr/testify/require"
	"testing"
)

func setupSettingTest() (err error) {
	return mongo.InitMongo()
}

func cleanupSettingTest() {
	_ = mongo.GetMongoCol(ModelColNameSetting).Delete(nil)
	_ = mongo.GetMongoCol(ArtifactColName).Delete(nil)
}

func TestSetting_Add(t *testing.T) {
	err := setupSettingTest()
	require.Nil(t, err)

	s := Setting{}

	err = s.Add()
	require.Nil(t, err)
	require.NotNil(t, s.Id)

	a, err := s.GetArtifact()
	require.Nil(t, err)
	require.Equal(t, s.Id, a.Id)
	require.NotNil(t, a.CreateTs)
	require.NotNil(t, a.UpdateTs)

	col := mongo.GetMongoCol(ModelColNameSetting)
	col.GetContext()

	cleanupSettingTest()
}

func TestSetting_Save(t *testing.T) {
	err := setupSettingTest()
	require.Nil(t, err)

	s := Setting{}

	err = s.Add()
	require.Nil(t, err)

	key := "test_setting"
	s.Key = key
	err = s.Save()
	require.Nil(t, err)

	err = mongo.GetMongoCol(ModelColNameSetting).FindId(s.Id).One(&s)
	require.Nil(t, err)
	require.Equal(t, key, s.Key)

	cleanupSettingTest()
}

func TestSetting_Delete(t *testing.T) {
	err := setupSettingTest()
	require.Nil(t, err)

	s := Setting{
		Key: "test_setting",
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

	cleanupSettingTest()
}

func TestSetting_DeleteList(t *testing.T) {
	err := setupSettingTest()
	require.Nil(t, err)

	doc := Setting{
		Key: "test_Setting",
	}

	err = doc.Add()
	require.Nil(t, err)

	err = SettingService.DeleteList(nil)
	require.Nil(t, err)

	total, err := SettingService.Count(nil)
	require.Equal(t, 0, total)

	cleanupSettingTest()
}
