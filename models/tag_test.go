package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTag_Add(t *testing.T) {
	setupTest(t)

	s := Tag{}

	err := s.Add()
	require.Nil(t, err)
	require.NotNil(t, s.Id)

	a, err := s.GetArtifact()
	require.Nil(t, err)
	require.Equal(t, s.Id, a.Id)
	require.NotNil(t, a.CreateTs)
	require.NotNil(t, a.UpdateTs)

	col := mongo.GetMongoCol(interfaces.ModelColNameTag)
	col.GetContext()
}

func TestTag_Save(t *testing.T) {
	setupTest(t)

	s := Tag{}

	err := s.Add()
	require.Nil(t, err)

	name := "test_schedule"
	s.Name = name
	err = s.Save()
	require.Nil(t, err)

	err = mongo.GetMongoCol(interfaces.ModelColNameTag).FindId(s.Id).One(&s)
	require.Nil(t, err)
	require.Equal(t, name, s.Name)
}

func TestTag_Delete(t *testing.T) {
	setupTest(t)

	s := Tag{
		Name: "test_schedule",
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

func TestTag_DeleteList(t *testing.T) {
	setupTest(t)

	doc := Tag{
		Name: "test_Tag",
	}

	err := doc.Add()
	require.Nil(t, err)

	err = TagService.DeleteList(nil)
	require.Nil(t, err)

	total, err := TagService.Count(nil)
	require.Equal(t, 0, total)
}
