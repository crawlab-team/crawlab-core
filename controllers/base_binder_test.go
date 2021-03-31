package controllers

import (
	"github.com/crawlab-team/crawlab-core/models"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestJsonBinder_getDoc(t *testing.T) {
	doc, err := NewJsonBinder(ControllerIdProject).getDoc()
	require.Nil(t, err)

	id := primitive.NewObjectID()
	doc = &models.Project{
		Id: id,
	}
	require.Equal(t, id, doc.GetId())
}

func TestJsonBinder_getDocList(t *testing.T) {
	docs, err := NewJsonBinder(ControllerIdProject).getDocList()
	require.Nil(t, err)
	require.Equal(t, 0, len(docs))

	var ids []primitive.ObjectID
	n := 10
	for i := 0; i < n; i++ {
		id := primitive.NewObjectID()
		docs = append(docs, &models.Project{
			Id: id,
		})
		ids = append(ids, id)
	}

	for i, doc := range docs {
		id := ids[i]
		require.Equal(t, id, doc.GetId())
	}
}
