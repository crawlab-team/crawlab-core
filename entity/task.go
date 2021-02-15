package entity

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskMessage struct {
	Id  primitive.ObjectID
	Cmd string
}

func (m *TaskMessage) ToString() (string, error) {
	data, err := json.Marshal(&m)
	if err != nil {
		return "", err
	}
	return string(data), err
}
