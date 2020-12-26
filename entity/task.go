package entity

import (
	"encoding/json"
)

type TaskMessage struct {
	Id  string
	Cmd string
}

func (m *TaskMessage) ToString() (string, error) {
	data, err := json.Marshal(&m)
	if err != nil {
		return "", err
	}
	return string(data), err
}
