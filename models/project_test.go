package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProject_Add(t *testing.T) {
	setupTest(t)

	p := Project{}

	err := p.Add()
	require.Nil(t, err)
	require.NotNil(t, p.Id)

	a, err := p.GetArtifact()
	require.Nil(t, err)
	require.Equal(t, p.Id, a.Id)
	require.NotNil(t, a.CreateTs)
	require.NotNil(t, a.UpdateTs)

	col := mongo.GetMongoCol(interfaces.ModelColNameProject)
	col.GetContext()
}

func TestProject_Save(t *testing.T) {
	setupTest(t)

	p := Project{}

	err := p.Add()
	require.Nil(t, err)

	name := "test_project"
	p.Name = name
	err = p.Save()
	require.Nil(t, err)

	err = mongo.GetMongoCol(interfaces.ModelColNameProject).FindId(p.Id).One(&p)
	require.Nil(t, err)
	require.Equal(t, name, p.Name)
}

func TestProject_Delete(t *testing.T) {
	setupTest(t)

	p := Project{
		Name: "test_project",
	}

	err := p.Add()
	require.Nil(t, err)

	err = p.Delete()
	require.Nil(t, err)

	var a Artifact
	col := mongo.GetMongoCol(interfaces.ModelColNameArtifact)
	err = col.FindId(p.Id).One(&a)
	require.Nil(t, err)
	require.NotNil(t, a.Obj)
	require.True(t, a.Del)
}

func TestProject_DeleteList(t *testing.T) {
	setupTest(t)

	doc := Project{
		Name: "test_Project",
	}

	err := doc.Add()
	require.Nil(t, err)

	err = ProjectService.DeleteList(nil)
	require.Nil(t, err)

	total, err := ProjectService.Count(nil)
	require.Equal(t, 0, total)
}
