package model

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/stretchr/testify/require"
	"testing"
)

func setupVariableTest() (err error) {
	return mongo.InitMongo()
}

func cleanupVariableTest() {
	_ = mongo.GetMongoCol(VariableColName).Delete(nil)
	_ = mongo.GetMongoCol(ArtifactColName).Delete(nil)
}

func TestVariable_Add(t *testing.T) {
	err := setupVariableTest()
	require.Nil(t, err)

	s := Variable{}

	err = s.Add()
	require.Nil(t, err)
	require.NotNil(t, s.Id)

	a, err := s.GetArtifact()
	require.Nil(t, err)
	require.Equal(t, s.Id, a.Id)
	require.NotNil(t, a.CreateTs)
	require.NotNil(t, a.UpdateTs)

	col := mongo.GetMongoCol(VariableColName)
	col.GetContext()

	cleanupVariableTest()
}

func TestVariable_Save(t *testing.T) {
	err := setupVariableTest()
	require.Nil(t, err)

	s := Variable{}

	err = s.Add()
	require.Nil(t, err)

	key := "test_variable"
	s.Key = key
	err = s.Save()
	require.Nil(t, err)

	err = mongo.GetMongoCol(VariableColName).FindId(s.Id).One(&s)
	require.Nil(t, err)
	require.Equal(t, key, s.Key)

	cleanupVariableTest()
}

func TestVariable_Delete(t *testing.T) {
	err := setupVariableTest()
	require.Nil(t, err)

	s := Variable{
		Key: "test_variable",
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

	cleanupVariableTest()
}
