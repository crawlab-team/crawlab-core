package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetting_Add(t *testing.T) {
	setupTest(t)

	s := Setting{}

	err := s.Add()
	require.Nil(t, err)
	require.NotNil(t, s.Id)
}

func TestSetting_Save(t *testing.T) {
	setupTest(t)

	s := Setting{}

	err := s.Add()
	require.Nil(t, err)

	key := "test_setting"
	s.Key = key
	err = s.Save()
	require.Nil(t, err)

	err = mongo.GetMongoCol(interfaces.ModelColNameSetting).FindId(s.Id).One(&s)
	require.Nil(t, err)
	require.Equal(t, key, s.Key)
}

func TestSetting_Delete(t *testing.T) {
	setupTest(t)

	s := Setting{
		Key: "test_setting",
	}

	err := s.Add()
	require.Nil(t, err)

	err = s.Delete()
	require.Nil(t, err)

	var a Artifact
	col := mongo.GetMongoCol(interfaces.ModelColNameArtifact)
	err = col.FindId(s.Id).One(&a)
	require.Nil(t, err)
	require.NotNil(t, a.Obj)
	require.True(t, a.Del)
}

func TestSetting_DeleteList(t *testing.T) {
	setupTest(t)

	doc := Setting{
		Key: "test_Setting",
	}

	err := doc.Add()
	require.Nil(t, err)

	err = MustGetService(interfaces.ModelIdSetting).DeleteList(nil)
	require.Nil(t, err)

	total, err := MustGetService(interfaces.ModelIdSetting).Count(nil)
	require.Equal(t, 0, total)
}
