package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/stretchr/testify/require"
	"testing"
)

func setupTokenTest() (err error) {
	return mongo.InitMongo()
}

func cleanupTokenTest() {
	_ = mongo.GetMongoCol(interfaces.ModelColNameToken).Delete(nil)
	_ = mongo.GetMongoCol(interfaces.ModelColNameArtifact).Delete(nil)
}

func TestToken_Add(t *testing.T) {
	var err error
	setupTest(t)

	token := Token{}

	err = token.Add()
	require.Nil(t, err)
	require.NotNil(t, token.Id)
}

func TestToken_Save(t *testing.T) {
	var err error
	setupTest(t)

	token := Token{}

	err = token.Add()
	require.Nil(t, err)

	tokenValue := "test_token"
	token.Token = tokenValue
	err = token.Save()
	require.Nil(t, err)

	err = mongo.GetMongoCol(interfaces.ModelColNameToken).FindId(token.Id).One(&token)
	require.Nil(t, err)
	require.Equal(t, tokenValue, token.Token)
}

func TestToken_Delete(t *testing.T) {
	var err error
	setupTest(t)

	s := Token{
		Token: "test_token",
	}

	err = s.Add()
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

func TestToken_DeleteList(t *testing.T) {
	var err error
	setupTest(t)

	doc := Token{
		Token: "test_Token",
	}

	err = doc.Add()
	require.Nil(t, err)

	err = MustGetService(interfaces.ModelIdToken).DeleteList(nil)
	require.Nil(t, err)

	total, err := MustGetService(interfaces.ModelIdToken).Count(nil)
	require.Equal(t, 0, total)
}
