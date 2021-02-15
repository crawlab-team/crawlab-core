package model

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceInterface interface {
	findId(id primitive.ObjectID) (fr *mongo.FindResult)
	find(query bson.M, opts *mongo.FindOptions) (fr *mongo.FindResult)
	deleteId(id primitive.ObjectID) (err error)
	delete(query bson.M) (err error)
}

func NewService(colName string) (svc *Service) {
	col := mongo.GetMongoCol(colName)
	return &Service{
		col: col,
	}
}

type Service struct {
	col *mongo.Col
}

func (s *Service) findId(id primitive.ObjectID) (fr *mongo.FindResult) {
	if s.col == nil {
		return mongo.NewFindResultWithError(constants.ErrMissingCol)
	}
	return s.col.FindId(id)
}

func (s *Service) find(query bson.M, opts *mongo.FindOptions) (fr *mongo.FindResult) {
	if s.col == nil {
		return mongo.NewFindResultWithError(constants.ErrMissingCol)
	}
	return s.col.Find(query, opts)
}

func (s *Service) deleteId(id primitive.ObjectID) (err error) {
	if s.col == nil {
		return constants.ErrMissingCol
	}
	var doc BaseModel
	if err := s.findId(id).One(&doc); err != nil {
		return err
	}
	d := NewDelegate(s.col.GetName(), doc)
	return d.Delete()
}

func (s *Service) delete(query bson.M) (err error) {
	if s.col == nil {
		return constants.ErrMissingCol
	}
	var docs []BaseModel
	if err := s.find(query, nil).All(&docs); err != nil {
		return err
	}
	for _, doc := range docs {
		if err := s.deleteId(doc.Id); err != nil {
			return err
		}
	}
	return nil
}
