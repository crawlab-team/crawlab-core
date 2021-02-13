package model

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/stretchr/testify/require"
	"testing"
)

func setupProjectTest() (err error) {
	return mongo.InitMongo()
}

func cleanupProjectTest() {
	_ = mongo.GetMongoCol(ProjectColName).Delete(nil)
	_ = mongo.GetMongoCol(ArtifactColName).Delete(nil)
}

func TestProject_Add(t *testing.T) {
	err := setupProjectTest()
	require.Nil(t, err)

	p := Project{}

	err = p.Add()
	require.Nil(t, err)
	require.NotNil(t, p.Id)

	a, err := p.GetArtifact()
	require.Nil(t, err)
	require.Equal(t, p.Id, a.Id)
	require.NotNil(t, a.CreateTs)
	require.NotNil(t, a.UpdateTs)

	col := mongo.GetMongoCol(ProjectColName)
	col.GetContext()

	cleanupProjectTest()
}

func TestProject_Save(t *testing.T) {
	err := setupProjectTest()
	require.Nil(t, err)

	p := Project{}

	err = p.Add()
	require.Nil(t, err)

	name := "test_project"
	p.Name = name
	err = p.Save()
	require.Nil(t, err)

	err = mongo.GetMongoCol(ProjectColName).FindId(p.Id).One(&p)
	require.Nil(t, err)
	require.Equal(t, name, p.Name)

	cleanupProjectTest()
}

func TestProject_Delete(t *testing.T) {
	err := setupProjectTest()
	require.Nil(t, err)

	p := Project{
		Name: "test_project",
	}

	err = p.Add()
	require.Nil(t, err)

	err = p.Delete()
	require.Nil(t, err)

	var a Artifact
	col := mongo.GetMongoCol(ArtifactColName)
	err = col.FindId(p.Id).One(&a)
	require.Nil(t, err)
	require.NotNil(t, a.Obj)
	require.True(t, a.Del)

	cleanupProjectTest()
}
