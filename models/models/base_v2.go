package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ModelV2 interface {
	GetCreatedAt() time.Time
	SetCreatedAt(time.Time)
	GetCreatedBy() primitive.ObjectID
	SetCreatedBy(primitive.ObjectID)
	SetCreated(primitive.ObjectID)
	GetUpdatedAt() time.Time
	SetUpdatedAt(time.Time)
	GetUpdatedBy() primitive.ObjectID
	SetUpdatedBy(primitive.ObjectID)
	SetUpdated(primitive.ObjectID)
}

type BaseModelV2[T any] struct {
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	CreatedBy primitive.ObjectID `json:"created_by" bson:"created_by"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	UpdatedBy primitive.ObjectID `json:"updated_by" bson:"updated_by"`
}

func (m *BaseModelV2[T]) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func (m *BaseModelV2[T]) SetCreatedAt(t time.Time) {
	m.CreatedAt = t
}

func (m *BaseModelV2[T]) GetCreatedBy() primitive.ObjectID {
	return m.CreatedBy
}

func (m *BaseModelV2[T]) SetCreatedBy(id primitive.ObjectID) {
	m.CreatedBy = id
}

func (m *BaseModelV2[T]) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}

func (m *BaseModelV2[T]) SetUpdatedAt(t time.Time) {
	m.UpdatedAt = t
}

func (m *BaseModelV2[T]) GetUpdatedBy() primitive.ObjectID {
	return m.UpdatedBy
}

func (m *BaseModelV2[T]) SetUpdatedBy(id primitive.ObjectID) {
	m.UpdatedBy = id
}

func (m *BaseModelV2[T]) SetCreated(id primitive.ObjectID) {
	m.SetCreatedAt(time.Now())
	m.SetCreatedBy(id)
}

func (m *BaseModelV2[T]) SetUpdated(id primitive.ObjectID) {
	m.SetUpdatedAt(time.Now())
	m.SetUpdatedBy(id)
}
