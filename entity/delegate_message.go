package entity

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/go-trace"
)

type DelegateMessage struct {
	ModelId interfaces.ModelId             `json:"id"`
	Method  interfaces.ModelDelegateMethod `json:"m"`
	Data    []byte                         `json:"d"`
}

func (msg *DelegateMessage) GetModelId() interfaces.ModelId {
	return msg.ModelId
}

func (msg *DelegateMessage) GetMethod() interfaces.ModelDelegateMethod {
	return msg.Method
}

func (msg *DelegateMessage) GetData() []byte {
	return msg.Data
}

func (msg *DelegateMessage) ToBytes() (data []byte) {
	data, err := json.Marshal(*msg)
	if err != nil {
		_ = trace.TraceError(err)
		return data
	}
	return data
}
