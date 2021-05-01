package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNode_Add(t *testing.T) {
	setupTest(t)

	n := Node{
		Tags: []Tag{
			{Name: "tag 1"},
		},
	}

	err := n.Add()
	require.Nil(t, err)
	require.NotNil(t, n.Id)

	a, err := n.GetArtifact()
	require.Nil(t, err)
	require.Equal(t, n.Id, a.GetId())
	require.NotNil(t, a.GetSys().GetCreateTs())
	require.NotNil(t, a.GetSys().GetUpdateTs())
}

func TestNode_Save(t *testing.T) {
	setupTest(t)

	n := Node{}

	err := n.Add()

	name := "test_node"
	n.Name = name
	err = n.Save()
	require.Nil(t, err)

	err = mongo.GetMongoCol(interfaces.ModelColNameNode).FindId(n.Id).One(&n)
	require.Nil(t, err)
	require.Equal(t, name, n.Name)
}

func TestNode_Delete(t *testing.T) {
	setupTest(t)

	n := Node{
		Name: "test_node",
	}

	err := n.Add()
	require.Nil(t, err)

	err = n.Delete()
	require.Nil(t, err)

	var a Artifact
	col := mongo.GetMongoCol(interfaces.ModelColNameArtifact)
	err = col.FindId(n.Id).One(&a)
	require.Nil(t, err)
	require.NotNil(t, a.Obj)
	require.True(t, a.Del)
}

func TestNode_DeleteList(t *testing.T) {
	setupTest(t)

	doc := Node{
		Name: "test_node",
	}

	err := doc.Add()
	require.Nil(t, err)

	err = MustGetService(interfaces.ModelIdNode).DeleteList(nil)
	require.Nil(t, err)

	total, err := MustGetService(interfaces.ModelIdNode).Count(nil)
	require.Equal(t, 0, total)
}
