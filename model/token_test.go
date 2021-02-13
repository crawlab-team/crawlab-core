package model

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/stretchr/testify/require"
	"testing"
)

func setupTokenTest() (err error) {
	return mongo.InitMongo()
}

func cleanupTokenTest() {
	_ = mongo.GetMongoCol(TokenColName).Delete(nil)
	_ = mongo.GetMongoCol(ArtifactColName).Delete(nil)
}

func TestToken_Add(t *testing.T) {
	err := setupTokenTest()
	require.Nil(t, err)

	token := Token{}

	err = token.Add()
	require.Nil(t, err)
	require.NotNil(t, token.Id)

	a, err := token.GetArtifact()
	require.Nil(t, err)
	require.Equal(t, token.Id, a.Id)
	require.NotNil(t, a.CreateTs)
	require.NotNil(t, a.UpdateTs)

	col := mongo.GetMongoCol(TokenColName)
	col.GetContext()

	cleanupTokenTest()
}

func TestToken_Save(t *testing.T) {
	err := setupTokenTest()
	require.Nil(t, err)

	token := Token{}

	err = token.Add()
	require.Nil(t, err)

	tokenValue := "test_token"
	token.Token = tokenValue
	err = token.Save()
	require.Nil(t, err)

	err = mongo.GetMongoCol(TokenColName).FindId(token.Id).One(&token)
	require.Nil(t, err)
	require.Equal(t, tokenValue, token.Token)

	cleanupTokenTest()
}

func TestToken_Delete(t *testing.T) {
	err := setupTokenTest()
	require.Nil(t, err)

	s := Token{
		Token: "test_token",
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

	cleanupTokenTest()
}
