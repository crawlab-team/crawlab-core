package model

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/stretchr/testify/require"
	"testing"
)

func setupUserTest() (err error) {
	return mongo.InitMongo()
}

func cleanupUserTest() {
	_ = mongo.GetMongoCol(UserColName).Delete(nil)
	_ = mongo.GetMongoCol(ArtifactColName).Delete(nil)
}

func TestUser_Add(t *testing.T) {
	err := setupUserTest()
	require.Nil(t, err)

	u := User{}

	err = u.Add()
	require.Nil(t, err)
	require.NotNil(t, u.Id)

	a, err := u.GetArtifact()
	require.Nil(t, err)
	require.Equal(t, u.Id, a.Id)
	require.NotNil(t, a.CreateTs)
	require.NotNil(t, a.UpdateTs)

	col := mongo.GetMongoCol(UserColName)
	col.GetContext()

	cleanupUserTest()
}

func TestUser_Save(t *testing.T) {
	err := setupUserTest()
	require.Nil(t, err)

	u := User{}

	err = u.Add()
	require.Nil(t, err)

	name := "test_user"
	u.Username = name
	err = u.Save()
	require.Nil(t, err)

	err = mongo.GetMongoCol(UserColName).FindId(u.Id).One(&u)
	require.Nil(t, err)
	require.Equal(t, name, u.Username)

	cleanupUserTest()
}

func TestUser_Delete(t *testing.T) {
	err := setupUserTest()
	require.Nil(t, err)

	s := User{
		Username: "test_user",
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

	cleanupUserTest()
}

func TestUser_DeleteList(t *testing.T) {
	err := setupUserTest()
	require.Nil(t, err)

	doc := User{
		Username: "test_User",
	}

	err = doc.Add()
	require.Nil(t, err)

	err = UserService.DeleteList(nil)
	require.Nil(t, err)

	total, err := UserService.Count(nil)
	require.Equal(t, 0, total)

	cleanupUserTest()
}
