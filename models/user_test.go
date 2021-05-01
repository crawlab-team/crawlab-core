package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/stretchr/testify/require"
	"testing"
)

func setupUserTest() (err error) {
	return mongo.InitMongo()
}

func cleanupUserTest() {
	_ = mongo.GetMongoCol(interfaces.ModelColNameUser).Delete(nil)
	_ = mongo.GetMongoCol(interfaces.ModelColNameArtifact).Delete(nil)
}

func TestUser_Add(t *testing.T) {
	var err error
	setupTest(t)

	u := User{}

	err = u.Add()
	require.Nil(t, err)
	require.NotNil(t, u.Id)
}

func TestUser_Save(t *testing.T) {
	var err error
	setupTest(t)

	u := User{}

	err = u.Add()
	require.Nil(t, err)

	name := "test_user"
	u.Username = name
	err = u.Save()
	require.Nil(t, err)

	err = mongo.GetMongoCol(interfaces.ModelColNameUser).FindId(u.Id).One(&u)
	require.Nil(t, err)
	require.Equal(t, name, u.Username)
}

func TestUser_Delete(t *testing.T) {
	var err error
	setupTest(t)

	s := User{
		Username: "test_user",
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

func TestUser_DeleteList(t *testing.T) {
	var err error
	setupTest(t)

	doc := User{
		Username: "test_User",
	}

	err = doc.Add()
	require.Nil(t, err)

	err = MustGetService(interfaces.ModelIdUser).DeleteList(nil)
	require.Nil(t, err)

	total, err := MustGetService(interfaces.ModelIdUser).Count(nil)
	require.Equal(t, 0, total)
}
