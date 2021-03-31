package models

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestConvertToBaseModelInterface_Struct(t *testing.T) {
	id := primitive.NewObjectID()
	obj := Project{
		Id: id,
	}
	res, err := ConvertToBaseModelInterface(obj)
	require.Nil(t, res)
	require.NotNil(t, err)
	require.Equal(t, errors.ErrorModelNotImplemented, err)
}

func TestConvertToBaseModelInterface_Pointer(t *testing.T) {
	id := primitive.NewObjectID()
	ptr := &Project{
		Id: id,
	}
	res, err := ConvertToBaseModelInterface(ptr)
	require.Nil(t, err)
	require.Equal(t, id, res.GetId())
}

func TestGetBaseModelInterfaceList(t *testing.T) {
	n := 10
	var list []Project
	for i := 0; i < n; i++ {
		id := primitive.NewObjectID()
		p := Project{
			Id:          id,
			Name:        "test name",
			Description: "test description",
			Tags:        []string{"test tag"},
		}
		list = append(list, p)
	}
	docs, err := GetBaseModelInterfaceList(list)
	require.Nil(t, err)
	require.Equal(t, n, len(docs))
	for i, doc := range docs {
		require.Less(t, i, n)
		p := list[i]
		require.Equal(t, p.Id, doc.GetId())
	}
}
