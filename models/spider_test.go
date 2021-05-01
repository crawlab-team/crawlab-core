package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSpider_Add(t *testing.T) {
	setupTest(t)

	s := Spider{}

	err := s.Add()
	require.Nil(t, err)
	require.NotNil(t, s.Id)
}

func TestSpider_Save(t *testing.T) {
	setupTest(t)

	s := Spider{}

	err := s.Add()
	require.Nil(t, err)

	name := "test_spider"
	s.Name = name
	err = s.Save()
	require.Nil(t, err)

	err = mongo.GetMongoCol(interfaces.ModelColNameSpider).FindId(s.Id).One(&s)
	require.Nil(t, err)
	require.Equal(t, name, s.Name)
}

func TestSpider_Delete(t *testing.T) {
	setupTest(t)

	s := Spider{
		Name: "test_spider",
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

func TestSpider_DeleteList(t *testing.T) {
	setupTest(t)

	doc := Spider{
		Name: "test_Spider",
	}

	err := doc.Add()
	require.Nil(t, err)

	err = MustGetService(interfaces.ModelIdSpider).DeleteList(nil)
	require.Nil(t, err)

	total, err := MustGetService(interfaces.ModelIdSpider).Count(nil)
	require.Equal(t, 0, total)
}
